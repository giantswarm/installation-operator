package terraform

import (
	"context"

	"github.com/giantswarm/microerror"
	tfv1 "github.com/rancher/terraform-controller/pkg/apis/terraformcontroller.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
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
					Content: nil,
					Git:     tfv1.GitLocation{},
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

	r.logger.LogCtx(ctx, "ensured created")
	return nil
}
