package terraform

import (
	"context"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/pkg/label"
)

func Test_Resource_DynamoDBModule_newDelete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		obj                interface{}
		currentState       []ModuleState
		desiredState       []ModuleState
		expectedModuleName string
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
			currentState:       []ModuleState{},
			desiredState:       []ModuleState{},
			expectedModuleName: "",
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
			currentState: []ModuleState{},
			desiredState: []ModuleState{
				{
					Name: "desired",
				},
			},
			expectedModuleName: "",
		},
		{
			description: "current state not empty, desired state not empty but equal, expected desired state avoiding delivery log module",
			obj: &v1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						label.Installation: "5xchu",
					},
				},
			},
			currentState: []ModuleState{
				{
					Name: "current",
				},
			},
			desiredState: []ModuleState{
				{
					Name: "current",
				},
			},
			expectedModuleName: "current",
		},
	}

	var err error

	var newResource *Resource
	{
		c := Config{
			Logger:           microloggertest.New(),
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

			deleteChanges, ok := result.([]ModuleState)
			if !ok {
				t.Errorf("expected '%T', got '%T'", deleteChanges, result)
			}

			for _, deleteChange := range deleteChanges {
				if deleteChange.Name != tc.expectedModuleName {
					t.Errorf("expected %s, got %s", tc.expectedModuleName, deleteChange.Name)
				}
			}
		})
	}
}
