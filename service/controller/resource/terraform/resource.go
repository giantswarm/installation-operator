package terraform

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	// Name is the identifier of the resource.
	Name = "terraform"
)

// Config represents the configuration used to create a new s3module resource.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	AccessLogsExpiration int
	DeleteLoggingModule  bool
	IncludeTags          bool
}

// Resource implements the s3module resource.
type Resource struct {
	// Dependencies.
	logger micrologger.Logger

	// Settings.
	includeTags          bool
}

// New creates a new configured s3module resource.
func New(config Config) (*Resource, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// Settings.
	if config.AccessLogsExpiration < 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.AccessLogsExpiration must not be lower than 0", config)
	}

	r := &Resource{
		// Dependencies.
		logger: config.Logger,

		// Settings.
		includeTags:          config.IncludeTags,
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
