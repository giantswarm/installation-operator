package controller

import (
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/resource/wrapper/retryresource"
	"github.com/rancher/terraform-controller/pkg/generated/clientset/versioned"

	"github.com/giantswarm/installation-operator/pkg/project"
	"github.com/giantswarm/installation-operator/service/controller/key"
	"github.com/giantswarm/installation-operator/service/controller/resource/env"
	"github.com/giantswarm/installation-operator/service/controller/resource/module"
	"github.com/giantswarm/installation-operator/service/controller/resource/secret"
	"github.com/giantswarm/installation-operator/service/controller/resource/state"
)

type installationResourceSetConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	TFClient  versioned.Interface
}

func newInstallationResourceSet(config installationResourceSetConfig) (*controller.ResourceSet, error) {
	var err error

	var envResource resource.Interface
	{
		c := env.Config{
			K8sClient: config.K8sClient,
			Logger:   config.Logger,
		}

		envResource, err = env.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var secretResource resource.Interface
	{
		c := secret.Config{
			K8sClient: config.K8sClient,
			Logger:   config.Logger,
		}

		secretResource, err = secret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var moduleResource resource.Interface
	{
		c := module.Config{
			TFClient: config.TFClient,
			Logger:   config.Logger,
		}

		moduleResource, err = module.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var stateResource resource.Interface
	{
		c := state.Config{
			TFClient: config.TFClient,
			Logger:   config.Logger,
		}

		stateResource, err = state.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		envResource,
		secretResource,
		moduleResource,
		stateResource,
	}

	{
		c := retryresource.WrapConfig{
			Logger: config.Logger,
		}

		resources, err = retryresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	{
		c := metricsresource.WrapConfig{}

		resources, err = metricsresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	handlesFunc := func(obj interface{}) bool {
		cr, err := key.ToInstallation(obj)
		if err != nil {
			return false
		}

		if key.OperatorVersion(&cr) == project.BundleVersion() {
			return true
		}

		return false
	}

	var resourceSet *controller.ResourceSet
	{
		c := controller.ResourceSetConfig{
			Handles:   handlesFunc,
			Logger:    config.Logger,
			Resources: resources,
		}

		resourceSet, err = controller.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resourceSet, nil
}
