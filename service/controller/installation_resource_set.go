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
	"github.com/giantswarm/installation-operator/service/controller/resource/terraform"
)

type installationResourceSetConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	TFClient  versioned.Interface
}

func newInstallationResourceSet(config installationResourceSetConfig) (*controller.ResourceSet, error) {
	var err error

	var terraformResource resource.Interface
	{
		c := terraform.Config{
			TFClient: config.TFClient,
			Logger:   config.Logger,
		}

		terraformResource, err = terraform.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		terraformResource,
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
