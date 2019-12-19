package secret

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func terraformEncodeSlice(items []string) string {
	encoded := ""
	for i, item := range items {
		encoded += fmt.Sprintf("\"%s\"", item)
		if i < len(items)-1 {
			encoded += ","
		}
	}
	return fmt.Sprintf("[%s]", encoded)
}

func createExpected(cr v1alpha1.Installation) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		StringData: map[string]string{
			"aws_account":               cr.Spec.AWS.HostCluster.Account,
			"aws_region":                cr.Spec.AWS.Region,
			"cluster_name":              cr.Name,
			"nodes_vault_token":         cr.Status.NodeVaultToken,
			"base_domain":               cr.Spec.Base,
			"root_dns_zone_id":          cr.Status.RootDNSZone,
			"aws_customer_gateway_id_0": cr.Status.VPNGateway0,
			"aws_customer_gateway_id_1": cr.Status.VPNGateway1,
			"subnets_bastion":           terraformEncodeSlice(cr.Status.BastionSubnets),
			"container_linux_version":   cr.Spec.ContainerLinuxVersion,
		},
	}
}

func secretDataEqual(a, b *v1.Secret) bool {
	return reflect.DeepEqual(a.Data, b.Data) && reflect.DeepEqual(a.StringData, b.StringData)
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
	expected := createExpected(cr)
	actual, err := r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Create(expected)
		if err != nil {
			return microerror.Mask(err)
		}
	} else if err != nil {
		return microerror.Mask(err)
	} else if secretDataEqual(actual, expected) {
		_, err = r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Update(expected)
	}

	return nil
}
