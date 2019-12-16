package terraform

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	createModulesState, err := toModulesState(createChange)
	if err != nil {
		return microerror.Mask(err)
	}
	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, moduleInput := range createModulesState {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating terraform module %#q", moduleInput.Name))

		{
			_, err = r.k8sClient.CreateModule(i)
			if IsModuleAlreadyExists(err) {
				// Fall through.
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created dynamodb module %#q", moduleInput.Name))
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentModules, err := toModulesState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredModules, err := toModulesState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var createState []ModuleState
	for _, module := range desiredModules {
		if !containsModuleState(module.Name, currentModules) {
			createState = append(createState, module)
		}
	}

	return createState, nil
}
