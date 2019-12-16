package dynamodb

import (
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_ContainsTableState(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description      string
		installation     string
		bucketNameToFind string
		bucketStateList  []TableState
		expectedValue    bool
	}{
		{
			description:      "basic match",
			installation:     "test-install",
			bucketNameToFind: "bck1",
			bucketStateList:  []TableState{},
			expectedValue:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := containsTableState(tc.bucketNameToFind, tc.bucketStateList)

			if result != tc.expectedValue {
				t.Fatalf("Expected %t tags, found %t", tc.expectedValue, result)
			}
		})
	}
}

func Test_TableCanBeDeleted(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description         string
		installation        string
		deleteLoggingBucket bool
		bucketState         TableState
		expectedValue       bool
	}{
		{
			description:         "test env true",
			installation:        "test-install",
			deleteLoggingBucket: true,
			bucketState:         TableState{},
			expectedValue:       true,
		},
		{
			description:         "test env false",
			installation:        "test-install",
			deleteLoggingBucket: false,
			bucketState:         TableState{},
			expectedValue:       true,
		},
		{
			description:         "test env true no logging bucket",
			installation:        "test-install",
			deleteLoggingBucket: true,
			bucketState:         TableState{},
			expectedValue:       true,
		},
		{
			description:         "test env true logging bucket",
			installation:        "test-install",
			deleteLoggingBucket: true,
			bucketState:         TableState{},
			expectedValue:       true,
		},
		{
			description:         "test env false no logging bucket",
			installation:        "test-install",
			deleteLoggingBucket: true,
			bucketState:         TableState{},
			expectedValue:       true,
		},
		{
			description:         "test env false logging bucket",
			installation:        "test-install",
			deleteLoggingBucket: false,
			bucketState:         TableState{},
			expectedValue:       false,
		},
	}

	c := Config{}

	c.Logger = microloggertest.New()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			result := r.canBeDeleted(tc.bucketState)

			if result != tc.expectedValue {
				t.Fatalf("Expected %t, found %t", tc.expectedValue, result)
			}
		})
	}
}
