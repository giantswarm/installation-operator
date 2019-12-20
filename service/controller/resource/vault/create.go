package vault

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/hashicorp/vault/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

const (
	defaultTTLHours = 4380
	maxTTLHours     = 8760
)

func (r *Resource) transitEnabled() (bool, error) {
	mounts, err := r.vaultClient.Sys().ListMounts()
	if err != nil {
		return false, microerror.Mask(err)
	}
	for _, mount := range mounts {
		if mount.Type == "transit" {
			return true, nil
		}
	}
	return false, nil
}

func (r *Resource) ensureTransitEnabled() error {
	enabled, err := r.transitEnabled()
	if err != nil {
		return microerror.Mask(err)
	}
	if !enabled {
		err := r.vaultClient.Sys().Mount("transit", &api.MountInput{Type: "transit"})
		if err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}

func (r *Resource) autounsealKeyExists() (bool, error) {
	_, err := r.vaultClient.Logical().Read("transit/keys/autounseal")
	if err != nil {
		return false, microerror.Mask(err)
	}
	return true, nil
}

func (r *Resource) ensureAutounsealKey() error {
	exists, err := r.autounsealKeyExists()
	if err != nil {
		return microerror.Mask(err)
	}
	if !exists {
		_, err := r.vaultClient.Logical().Write("transit/keys/autounseal", nil)
		if err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}

func (r *Resource) policyAutounsealCreated() (bool, error) {
	policies, err := r.vaultClient.Sys().ListPolicies()
	if err != nil {
		return false, microerror.Mask(err)
	}
	for _, policy := range policies {
		if policy == "autounseal" {
			return true, nil
		}
	}
	return false, nil
}

func (r *Resource) ensureAutounsealPolicy() error {
	created, err := r.policyAutounsealCreated()
	if err != nil {
		return microerror.Mask(err)
	}
	if !created {
		rules := `
path "transit/encrypt/autounseal" {
  capabilities = [ "update" ]
}

path "transit/decrypt/autounseal" {
  capabilities = [ "update" ]
}`
		err := r.vaultClient.Sys().PutPolicy("autounseal", rules)
		if err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}

func (r *Resource) authTokenTTLCorrect() (bool, error) {
	mount, err := r.vaultClient.Sys().MountConfig("auth/token")
	if err != nil {
		return false, microerror.Mask(err)
	}
	correct := mount.DefaultLeaseTTL == defaultTTLHours && mount.MaxLeaseTTL == maxTTLHours
	return correct, nil
}

func (r *Resource) ensureAuthTokenTTL() error {
	correct, err := r.authTokenTTLCorrect()
	if err != nil {
		return microerror.Mask(err)
	}
	if !correct {
		err := r.vaultClient.Sys().TuneMount("auth/token", api.MountConfigInput{
			DefaultLeaseTTL: fmt.Sprintf("%dh", defaultTTLHours),
			MaxLeaseTTL:     fmt.Sprintf("%dh", maxTTLHours),
		})
		if err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}

func (r *Resource) createVaultToken(name string) (string, error) {
	err := r.ensureTransitEnabled()
	if err != nil {
		return "", microerror.Mask(err)
	}
	err = r.ensureAutounsealKey()
	if err != nil {
		return "", microerror.Mask(err)
	}
	err = r.ensureAutounsealPolicy()
	if err != nil {
		return "", microerror.Mask(err)
	}
	err = r.ensureAuthTokenTTL()
	if err != nil {
		return "", microerror.Mask(err)
	}
	token, err := r.vaultClient.Auth().Token().Create(&api.TokenCreateRequest{
		ID:       name,
		Policies: []string{"autounseal"},
		Period:   fmt.Sprintf("%dh", defaultTTLHours),
	})
	tokenID, err := token.TokenID()
	if err != nil {
		return "", microerror.Mask(err)
	}
	return tokenID, nil
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if cr.Spec.Provider != "aws" {
		r.logger.LogCtx(ctx, "message", "provider not supported", "provider", cr.Spec.Provider)
		return nil
	}

	secret, err := r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	if len(secret.Data["node_vault_token"]) > 0 {
		return nil
	}

	token, err := r.createVaultToken(cr.Name)
	if err != nil {
		return microerror.Mask(err)
	}
	secret.Data["node_vault_token"] = []byte(token)
	_, err = r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Update(secret)

	return nil
}
