package terraform

import (
	"context"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/pkg/label"
	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
)

func Test_Resource_DynamoDBModule_GetDesiredState(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		obj           interface{}
		expectedNames []string
		description   string
	}{
		{
			description: "Get module name from custom object.",
			obj: &v1alpha1.Installation{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						label.Installation: "5xchu",
					},
				},
			},
			expectedNames: []string{
				"5xchu-g8s-access-logs",
				"myaccountid-g8s-5xchu",
			},
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
			t.Fatal("expected", nil, "got", err)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := controllercontext.NewContext(context.Background(), testContextWithAccountID("myaccountid"))

			result, err := newResource.GetDesiredState(ctx, tc.obj)
			if err != nil {
				t.Fatalf("expected '%v' got '%#v'", nil, err)
			}

			desiredModules, ok := result.([]ModuleState)
			if !ok {
				t.Fatalf("case expected '%T', got '%T'", desiredModules, result)
			}

			// Order should be respected in the slice returned (always delivery log module first)
			for key, desiredModule := range desiredModules {
				if tc.expectedNames[key] != desiredModule.Name {
					t.Fatalf("expected module name %q got %q", tc.expectedNames[key], desiredModule.Name)
				}
			}
		})
	}
}

func testContextWithAccountID(id string) controllercontext.Context {
	return controllercontext.Context{
		Status: controllercontext.ContextStatus{
			AWSAccountID: id,
		},
	}
}
