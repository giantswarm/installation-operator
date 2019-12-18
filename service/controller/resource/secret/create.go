package secret

import (
	"context"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if cr.Spec.Provider != "aws" {
		r.logger.LogCtx(ctx, "message", "provider not supported", "provider", cr.Spec.Provider)
		return nil
	}
	_, err = r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		secret := v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Name,
				Namespace: cr.Namespace,
			},
			StringData: map[string]string{
				"do_token": "",
				"do_name":  "terraform-controller-test",
			},
		}
		_, err = r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Create(&secret)
		if err != nil {
			return microerror.Mask(err)
		}
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
