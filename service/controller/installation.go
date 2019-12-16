package controller

import (
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"

	"github.com/giantswarm/installation-operator/pkg/project"
)

type InstallationConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type Installation struct {
	*controller.Controller
}

func NewInstallation(config InstallationConfig) (*Installation, error) {
	var err error

	resourceSets, err := newInstallationResourceSets(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			CRD:          v1alpha1.NewInstallationCRD(),
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			ResourceSets: resourceSets,
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(v1alpha1.Installation)
			},
			Name: project.Name() + "-installation-controller",
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Installation{
		Controller: operatorkitController,
	}

	return c, nil
}

func newInstallationResourceSets(config InstallationConfig) ([]*controller.ResourceSet, error) {
	var err error

	var installationResourceSet *controller.ResourceSet
	{
		c := installationResourceSetConfig{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		installationResourceSet, err = newInstallationResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resourceSets := []*controller.ResourceSet{
		installationResourceSet,
	}

	return resourceSets, nil
}
