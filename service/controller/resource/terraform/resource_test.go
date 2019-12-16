package terraform

import (
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_ContainsModuleState(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description      string
		installation     string
		moduleNameToFind string
		moduleStateList  []ModuleState
		expectedValue    bool
	}{
		{
			description:      "basic match",
			installation:     "test-install",
			moduleNameToFind: "bck1",
			moduleStateList:  []ModuleState{},
			expectedValue:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := containsModuleState(tc.moduleNameToFind, tc.moduleStateList)

			if result != tc.expectedValue {
				t.Fatalf("Expected %t tags, found %t", tc.expectedValue, result)
			}
		})
	}
}

func Test_ModuleCanBeDeleted(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description         string
		installation        string
		deleteLoggingModule bool
		moduleState         ModuleState
		expectedValue       bool
	}{
		{
			description:         "test env true",
			installation:        "test-install",
			deleteLoggingModule: true,
			moduleState:         ModuleState{},
			expectedValue:       true,
		},
		{
			description:         "test env false",
			installation:        "test-install",
			deleteLoggingModule: false,
			moduleState:         ModuleState{},
			expectedValue:       true,
		},
		{
			description:         "test env true no logging module",
			installation:        "test-install",
			deleteLoggingModule: true,
			moduleState:         ModuleState{},
			expectedValue:       true,
		},
		{
			description:         "test env true logging module",
			installation:        "test-install",
			deleteLoggingModule: true,
			moduleState:         ModuleState{},
			expectedValue:       true,
		},
		{
			description:         "test env false no logging module",
			installation:        "test-install",
			deleteLoggingModule: true,
			moduleState:         ModuleState{},
			expectedValue:       true,
		},
		{
			description:         "test env false logging module",
			installation:        "test-install",
			deleteLoggingModule: false,
			moduleState:         ModuleState{},
			expectedValue:       false,
		},
	}

	c := Config{}

	c.Logger = microloggertest.New()
	c.AccessLogsExpiration = 0

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c.DeleteLoggingModule = tc.deleteLoggingModule
			r, err := New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			result := r.canBeDeleted(tc.moduleState)

			if result != tc.expectedValue {
				t.Fatalf("Expected %t, found %t", tc.expectedValue, result)
			}
		})
	}
}
