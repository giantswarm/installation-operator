package terraform

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource/crud"
	"golang.org/x/sync/errgroup"

	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
)

const (
	// loopLimit is the maximum amount of delete actions we want to allow per
	// S3 module. Reason here is to execute resources fast and prevent
	// them blocking other resources for too long. In case a S3 module has more
	// than 3000 objects, we delete 3 batches of 1000 objects and leave the rest
	// for the next reconciliation loop.
	loopLimit = 3
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	modulesInput, err := toModulesState(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	g := &errgroup.Group{}

	for _, b := range modulesInput {
		moduleName := b.Name

		g.Go(func() error {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting terraform module %#q", moduleName))

			{
				_, err = r.k8sClient.DeleteModule(i)
				if err != nil {
					return microerror.Mask(err)
				}
			}

			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted S3 module %#q", moduleName))

			return nil
		})
	}

	err = g.Wait()
	if IsNotFound(err) {
		// Fall through.
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	delete, err := r.newDeleteChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := crud.NewPatch()
	patch.SetDeleteChange(delete)

	return patch, nil
}

func (r *Resource) newDeleteChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentModules, err := toModulesState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredModules, err := toModulesState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var modulesToDelete []ModuleState

	for _, module := range currentModules {
		// Destination Logs Module should not be deleted because it has to keep logs
		// even when cluster is removed (rotation of these logs are managed externally).
		if r.canBeDeleted(module) && containsModuleState(module.Name, desiredModules) {
			modulesToDelete = append(modulesToDelete, module)
		}
	}

	return modulesToDelete, nil
}
