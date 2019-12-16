package key

import "fmt"

func AWSTags(getter LabelsGetter, installationName string) map[string]string {
	TagCloudProvider := ClusterCloudProviderTag(getter)

	tags := map[string]string{
		TagCloudProvider: "owned",
		TagCluster:       InstallationID(getter),
		TagInstallation:  installationName,
		TagOrganization:  InstallationID(getter),
	}

	return tags
}

func BucketName(getter LabelsGetter, accountID string) string {
	return fmt.Sprintf("%s-g8s-%s", accountID, InstallationID(getter))
}

func ClusterCloudProviderTag(getter LabelsGetter) string {
	return fmt.Sprintf("kubernetes.io/cluster/%s", InstallationID(getter))
}
