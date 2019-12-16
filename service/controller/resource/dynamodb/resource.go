package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/installation-operator/pkg/awstags"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

const (
	// Name is the identifier of the resource.
	Name = "s3bucket"
	// LifecycleLoggingBucketID is the Lifecycle ID for the logging bucket
	LifecycleLoggingBucketID = "ExpirationLogs"
)

// Config represents the configuration used to create a new s3bucket resource.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	IncludeTags          bool
}

// Resource implements the s3bucket resource.
type Resource struct {
	// Dependencies.
	logger micrologger.Logger

	// Settings.
	accessLogsExpiration int
	deleteLoggingBucket  bool
	includeTags          bool
	installationName     string
}

// New creates a new configured s3bucket resource.
func New(config Config) (*Resource, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// Settings.

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

func toTablesState(v interface{}) ([]TableState, error) {
	if v == nil {
		return []TableState{}, nil
	}

	tablesState, ok := v.([]TableState)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", []TableState{}, v)
	}

	return tablesState, nil
}

func containsTableState(tableStateName string, tableStateList []TableState) bool {
	for _, b := range tableStateList {
		if b.Name == tableStateName {
			return true
		}
	}

	return false
}

func (r *Resource) getS3BucketTags(customObject v1alpha1.Installation) []*s3.Tag {
	tags := key.AWSTags(&customObject)
	return awstags.NewS3(tags)
}

func (r *Resource) canBeDeleted(table TableState) bool {
	return true
}
