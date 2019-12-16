package controllercontext

import (
	"github.com/giantswarm/installation-operator/client/aws"
)

type ContextClient struct {
	AWS aws.Clients
}
