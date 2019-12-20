package tfmodule

import (
	"context"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	_, err = r.tfClient.TerraformcontrollerV1().States(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if err == nil {
		return microerror.New("not deleting module until state is deleted")
	}
	if !errors.IsNotFound(err) {
		return microerror.Mask(err)
	}
	err = r.tfClient.TerraformcontrollerV1().Modules(cr.Namespace).Delete(cr.Name, &metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
