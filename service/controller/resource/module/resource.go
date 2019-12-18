package module

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/rancher/terraform-controller/pkg/generated/clientset/versioned"
)

const (
	Name = "terraform"
)

type Config struct {
	TFClient versioned.Interface
	Logger   micrologger.Logger
}

type Resource struct {
	logger   micrologger.Logger
	tfClient versioned.Interface
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.TFClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.TFClient must not be empty", config)
	}

	r := &Resource{
		logger:   config.Logger,
		tfClient: config.TFClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
