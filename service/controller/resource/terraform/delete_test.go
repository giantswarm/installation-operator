package terraform

import (
	"context"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/pkg/label"
)

func Test_Resource_DynamoDBTable_newDelete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		obj                interface{}
		currentState       []TableState
		desiredState       []TableState
		expectedBucketName string
		description        string
	}{
		{
			description: "current and desired state empty, expected empty",
			obj: &v1alpha1.Installation{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						label.Installation: "5xchu",
					},
				},
			},
			currentState:       []TableState{},
			desiredState:       []TableState{},
			expectedBucketName: "",
		},
		{
			description: "current state empty, desired state not empty, expected empty",
			obj: &v1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						label.Installation: "5xchu",
					},
				},
			},
			currentState: []TableState{},
			desiredState: []TableState{
				{
					Name: "desired",
				},
			},
			expectedBucketName: "",
		},
		{
			description: "current state not empty, desired state not empty but equal, expected desired state avoiding delivery log bucket",
			obj: &v1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						label.Installation: "5xchu",
					},
				},
			},
			currentState: []TableState{
				{
					Name: "current",
				},
			},
			desiredState: []TableState{
				{
					Name: "current",
				},
			},
			expectedBucketName: "current",
		},
	}

	var err error

	var newResource *Resource
	{
		c := Config{
			Logger:           microloggertest.New(),
			InstallationName: "test-install",
		}

		newResource, err = New(c)
		if err != nil {
			t.Error("expected", nil, "got", err)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := newResource.newDeleteChange(context.TODO(), tc.obj, tc.currentState, tc.desiredState)
			if err != nil {
				t.Errorf("expected '%v' got '%#v'", nil, err)
			}

			deleteChanges, ok := result.([]TableState)
			if !ok {
				t.Errorf("expected '%T', got '%T'", deleteChanges, result)
			}

			for _, deleteChange := range deleteChanges {
				if deleteChange.Name != tc.expectedBucketName {
					t.Errorf("expected %s, got %s", tc.expectedBucketName, deleteChange.Name)
				}
			}
		})
	}
}
