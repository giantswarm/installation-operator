package terraform

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var currentModuleState ModuleState
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "finding the terraform modules")

		module, err := r.tfClient.TerraformcontrollerV1().Modules(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
		if IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, microerror.Mask(err)
		}
		currentModuleState.Name = module.Name

		r.logger.LogCtx(ctx, "level", "debug", "message", "found the S3 modules")
	}

	return currentModuleState, nil
}

