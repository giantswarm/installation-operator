package state

import (
	"context"

	"github.com/giantswarm/microerror"
	tfv1 "github.com/rancher/terraform-controller/pkg/apis/terraformcontroller.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "message", "ensure state created")
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if cr.Spec.Provider != "aws" {
		r.logger.LogCtx(ctx, "message", "provider not supported", "provider", cr.Spec.Provider)
		return nil
	}
	_, err = r.tfClient.TerraformcontrollerV1().States(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		state := tfv1.State{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Name,
				Namespace: cr.Namespace,
			},
			Spec: tfv1.StateSpec{
				Image: "rancher/terraform-controller-executor:v0.0.10-alpha1",
				Variables: tfv1.Variables{
					EnvConfigName: []string{
						"env-config",
					},
					SecretNames: []string{
						cr.Name,
					},
				},
				ModuleName:      cr.Name,
				AutoConfirm:     true,
				DestroyOnDelete: true,
			},
		}
		_, err = r.tfClient.TerraformcontrollerV1().States(cr.Namespace).Create(&state)
		if err != nil {
			return microerror.Mask(err)
		}
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
