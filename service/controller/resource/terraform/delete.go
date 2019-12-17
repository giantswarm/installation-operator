package terraform

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource/crud"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting terraform module %#q", cr.Name))

	err = r.tfClient.TerraformcontrollerV1().Modules(cr.Namespace).Delete(cr.Name, &metav1.DeleteOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted S3 module %#q", cr.Name))

	return nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj interface{}, currentState interface{}, desiredState interface{}) (*crud.Patch, error) {
	return nil, nil
}
