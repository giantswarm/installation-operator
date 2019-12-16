package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/installation-operator/pkg/label"
)

// AWS Tags used for cost analysis and general resource tagging.
const (
	TagCluster           = "giantswarm.io/cluster"
	TagInstallation      = "giantswarm.io/installation"
	TagOrganization      = "giantswarm.io/organization"
)

func ToInstallation(v interface{}) (v1alpha1.Installation, error) {
	if v == nil {
		return v1alpha1.Installation{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.Cluster{}, v)
	}

	p, ok := v.(*v1alpha1.Installation)
	if !ok {
		return v1alpha1.Installation{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.Cluster{}, v)
	}

	c := p.DeepCopy()

	return *c, nil
}

func CredentialName(installation v1alpha1.Installation) string {
	return installation.Name
}

func CredentialNamespace(installation v1alpha1.Installation) string {
	return installation.Namespace
}

func InstallationID(getter LabelsGetter) string {
	return getter.GetLabels()[label.Installation]
}

func TargetLogBucketName(installation v1alpha1.Installation) string {
	return fmt.Sprintf("%s-g8s-access-logs", InstallationID(&installation))
}
