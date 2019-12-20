package vault

import (
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/hashicorp/vault/api"
)

const (
	Name = "vault"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	VaultClient *api.Client
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
	vaultClient *api.Client
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.VaultClient must not be empty", config)
	}

	r := &Resource{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		vaultClient: config.VaultClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
