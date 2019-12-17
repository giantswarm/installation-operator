package terraform

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/rancher/terraform-controller/pkg/generated/clientset/versioned"
)

const (
	// Name is the identifier of the resource.
	Name = "terraform"
)

// Config represents the configuration used to create a new s3module resource.
type Config struct {
	Logger micrologger.Logger
	TFClient versioned.Interface
}

// Resource implements the terraform resource.
type Resource struct {
	logger micrologger.Logger
	tfClient versioned.Interface
}

// New creates a new configured s3module resource.
func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.TFClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.TFClient must not be empty", config)
	}

	r := &Resource{
		logger: config.Logger,
		tfClient: config.TFClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toModulesState(v interface{}) ([]ModuleState, error) {
	if v == nil {
		return []ModuleState{}, nil
	}

	modulesState, ok := v.([]ModuleState)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", []ModuleState{}, v)
	}

	return modulesState, nil
}

func containsModuleState(moduleStateName string, moduleStateList []ModuleState) bool {
	for _, b := range moduleStateList {
		if b.Name == moduleStateName {
			return true
		}
	}

	return false
}

func (r *Resource) canBeDeleted(module ModuleState) bool {
	return true
}
