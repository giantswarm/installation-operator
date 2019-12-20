package tfmodule

import (
	"context"

	"github.com/giantswarm/microerror"
	tfv1 "github.com/giantswarm/terraform-controller/pkg/apis/terraformcontroller.cattle.io/v1"
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
	_, err = r.tfClient.TerraformcontrollerV1().Modules(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		module := tfv1.Module{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Name,
				Namespace: cr.Namespace,
			},
			Spec: tfv1.ModuleSpec{
				ModuleContent: tfv1.ModuleContent{
					Git: tfv1.GitLocation{
						URL:    "https://github.com/dramich/domodule",
						Branch: "master",
					},
				},
			},
		}
		_, err = r.tfClient.TerraformcontrollerV1().Modules(cr.Namespace).Create(&module)
		if err != nil {
			return microerror.Mask(err)
		}
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
