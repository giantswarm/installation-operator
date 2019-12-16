package terraform

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	modulesState := []ModuleState{
		{
			Name: key.ModuleName(&cr),
		},
	}

	return modulesState, nil
}
