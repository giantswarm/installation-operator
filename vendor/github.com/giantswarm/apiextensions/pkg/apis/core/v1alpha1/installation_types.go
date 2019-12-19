package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewInstallationCRD returns a new custom resource definition for
// Installation. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: installations.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: Installation
//         plural: installations
//         singular: installation
//
func NewInstallationCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "installations.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "Installation",
				Plural:   "installations",
				Singular: "installation",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Installation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              InstallationSpec   `json:"spec"`
	Status            InstallationStatus `json:"status"`
}

// InstallationAWSAccountConfig is the structure used to represent AWS credentials for an
// installation.
type InstallationAWSAccountConfig struct {
	Account          string `json:"account"`
	AdminRoleARN     string `json:"adminRoleARN"`
	CloudtrailBucket string `json:"cloudtrailBucket"`
	GuardDuty        bool   `json:"guardDuty"`
}

// InstallationAWSConfig is the structure used to represent AWS specific configuration for
// an installation.
type InstallationAWSConfig struct {
	Region       string                       `json:"region"`
	HostCluster  InstallationAWSAccountConfig `json:"hostCluster"`
	GuestCluster InstallationAWSAccountConfig `json:"guestCluster"`
}

// InstallationStorageConfig is a structure for information about what kind of persistent storage is used for the installation CP.
// this value make sense only for on-prem (KVM) installations.
type InstallationStorageConfig struct {
	StorageType string `json:"storageType,omitempty"`
	Description string `json:"description,omitempty"`
}

// InstallationServiceConfig is the structure used to represent details for each service in
// an installation.
type InstallationServiceConfig struct {
	Protocol  string `json:"protocol"`
	SubDomain string `json:"subDomain"`
	Port      int    `json:"port"`
}

// InstallationJumphostConfig is the data structure used to represent the jumphost config
// for an installation.
type InstallationJumphostConfig struct {
	Domains    []string          `json:"domains"`
	SubDomains []string          `json:"subDomains"`
	Users      map[string]string `json:"users"`
}

// InstallationInfo is a structure holding general information about the
// installation.
type InstallationInfo struct {
	AWS                   InstallationAWSConfig                `json:"aws"`
	Active                bool                                 `json:"active"`
	Base                  string                               `json:"base"`
	Codename              string                               `json:"codename"`
	ContainerLinuxVersion string                               `json:"containerLinuxVersion"`
	Created               DeepCopyTime                         `json:"created"`
	Customer              string                               `json:"customer"`
	Jumphost              InstallationJumphostConfig           `json:"jumphosts"`
	Machines              []string                             `json:"machines"`
	Pipeline              string                               `json:"pipeline"`
	Provider              string                               `json:"provider"`
	SupportSlackChannels  []string                             `json:"supportSlackChannels"`
	SSHTunnelNeeded       bool                                 `json:"sshTunnelNeeded"`
	SSHJumphostUser       string                               `json:"sshJumphostUser"`
	Services              map[string]InstallationServiceConfig `json:"services"`
	SolutionEngineer      string                               `json:"solutionEngineer"`
	Storage               InstallationStorageConfig            `json:"storage"`
	Updated               DeepCopyTime                         `json:"updated"`
}

type InstallationSpec struct {
	InstallationInfo     `json:",inline"` // separate struct so it can be used in opsctl
	DraughtsmanSecret    string           `json:"draughtsmanSecret"`
	DraughtsmanConfigMap string           `json:"draughtsmanConfigMap"`
}

type InstallationStatus struct {
	NodeVaultToken string   `json:"nodeVaultToken"`
	RootDNSZone    string   `json:"rootDNSZone"`
	VPNGateway0    string   `json:"vpnGateway0"`
	VPNGateway1    string   `json:"vpnGateway1"`
	BastionSubnets []string `json:"bastionSubnets"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstallationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Installation `json:"items"`
}
