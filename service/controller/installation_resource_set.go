package controller

import (
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/resource/wrapper/retryresource"
	"github.com/giantswarm/terraform-controller/pkg/generated/clientset/versioned"
	"github.com/hashicorp/vault/api"
	
	"github.com/giantswarm/installation-operator/pkg/project"
	"github.com/giantswarm/installation-operator/service/controller/key"
	"github.com/giantswarm/installation-operator/service/controller/resource/bastionsubnets"
	"github.com/giantswarm/installation-operator/service/controller/resource/bootstrap"
	"github.com/giantswarm/installation-operator/service/controller/resource/ipsec"
	"github.com/giantswarm/installation-operator/service/controller/resource/terraform"
	"github.com/giantswarm/installation-operator/service/controller/resource/tfenv"
	"github.com/giantswarm/installation-operator/service/controller/resource/tfmodule"
	"github.com/giantswarm/installation-operator/service/controller/resource/tfsecret"
	"github.com/giantswarm/installation-operator/service/controller/resource/vault"
)

type installationResourceSetConfig struct {
	K8sClient   k8sclient.Interface
	Logger      micrologger.Logger
	TFClient    versioned.Interface
	VaultClient *api.Client
}

func newInstallationResourceSet(config installationResourceSetConfig) (*controller.ResourceSet, error) {
	var err error

	var tfenvResource resource.Interface
	{
		c := tfenv.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		tfenvResource, err = tfenv.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var tfsecretResource resource.Interface
	{
		c := tfsecret.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		tfsecretResource, err = tfsecret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var bastionsubnetsResource resource.Interface
	{
		c := bastionsubnets.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		bastionsubnetsResource, err = bastionsubnets.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var tfmoduleResource resource.Interface
	{
		c := tfmodule.Config{
			TFClient: config.TFClient,
			Logger:   config.Logger,
		}

		tfmoduleResource, err = tfmodule.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

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

	var ipsecResource resource.Interface
	{
		c := ipsec.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		ipsecResource, err = ipsec.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultResource resource.Interface
	{
		c := vault.Config{
			K8sClient:   config.K8sClient,
			Logger:      config.Logger,
			VaultClient: config.VaultClient,
		}

		vaultResource, err = vault.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var bootstrapResource resource.Interface
	{
		c := bootstrap.Config{
			K8sClient:   config.K8sClient,
			Logger:      config.Logger,
		}

		bootstrapResource, err = bootstrap.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		tfenvResource,
		tfsecretResource,
		bastionsubnetsResource,
		tfmoduleResource,
		terraformResource,
		ipsecResource,
		vaultResource,
		bootstrapResource,
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
