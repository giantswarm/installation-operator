package terraform

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	tfv1 "github.com/rancher/terraform-controller/pkg/apis/terraformcontroller.cattle.io/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	for _, moduleInput := range createModulesState {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating terraform module %#q", moduleInput.Name))

		{
			module := tfv1.Module{
				ObjectMeta: metav1.ObjectMeta{
					Name:                       cr.Name,
					Namespace:                  cr.Namespace,
				},
				Spec:       tfv1.ModuleSpec{
						ModuleContent: tfv1.ModuleContent{
							Content: nil,
							Git:     tfv1.GitLocation{},
						},
				},
			}
			_, err := r.tfClient.TerraformcontrollerV1().Modules(cr.Namespace).Create(&module)
			if apierrors.IsAlreadyExists(err) {
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
