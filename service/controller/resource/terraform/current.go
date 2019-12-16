package terraform

import (
	"context"
	"fmt"
	"sync"

	"github.com/giantswarm/microerror"
	"golang.org/x/sync/errgroup"

	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	moduleStateNames := []string{
		key.ModuleName(&cr),
	}

	var currentModuleState []ModuleState
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "finding the terraform modules")

		g := &errgroup.Group{}
		m := sync.Mutex{}

		for _, inputModuleName := range moduleStateNames {
			moduleName := inputModuleName

			g.Go(func() error {
				inputModule := ModuleState{
					Name: moduleName,
				}

				// TODO this check should not be done here. Here we only fetch the
				// current state. We have to make a request anyway so fetching what we
				// want and handling the not found errors as usual should be the way to
				// go.
				//
				//
				//     https://github.com/giantswarm/giantswarm/issues/5246
				//
				isCreated, err := r.isModuleCreated(ctx, moduleName)
				if err != nil {
					return microerror.Mask(err)
				}
				if !isCreated {
					return nil
				}

				r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding the S3 module %#q", moduleName))

				m.Lock()
				currentModuleState = append(currentModuleState, inputModule)
				m.Unlock()

				r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found the S3 module %#q", moduleName))

				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			return nil, microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "found the S3 modules")
	}

	return currentModuleState, nil
}

func (r *Resource) isModuleCreated(ctx context.Context, name string) (bool, error) {
	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return false, microerror.Mask(err)
	}

	_, err = r.k8sClient.DescribeModule(headInput)
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}
