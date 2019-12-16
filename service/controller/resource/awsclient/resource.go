package awsclient

import (
	"context"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/installation-operator/client/aws"
	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
)

const (
	Name = "awsclient"
)

type Config struct {
	K8sClient          kubernetes.Interface
	Logger             micrologger.Logger
	ToInstallationFunc func(v interface{}) (v1alpha1.Installation, error)

	AWSConfig aws.Config
}

type Resource struct {
	k8sClient          kubernetes.Interface
	logger             micrologger.Logger
	toInstallationFunc func(v interface{}) (v1alpha1.Installation, error)

	awsConfig aws.Config
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.ToInstallationFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ToClusterFunc must not be empty", config)
	}

	r := &Resource{
		k8sClient:          config.K8sClient,
		logger:             config.Logger,
		toInstallationFunc: config.ToInstallationFunc,

		awsConfig: config.AWSConfig,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) addAWSClientsToContext(ctx context.Context, cr v1alpha1.Installation) error {
	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	{
		c := r.awsConfig

		clients, err := aws.NewClients(c)
		if err != nil {
			return microerror.Mask(err)
		}

		cc.Client.AWS = clients
	}

	return nil
}
