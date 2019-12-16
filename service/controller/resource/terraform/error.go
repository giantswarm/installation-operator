package terraform

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var tableNotFoundError = &microerror.Error{
	Kind: "tableNotFoundError",
}

// IsTableNotFound asserts bucket not found error from upstream's API code.
func IsTableNotFound(err error) bool {
	c := microerror.Cause(err)
	aerr, ok := c.(awserr.Error)
	if !ok {
		return false
	}
	if aerr.Code() == dynamodb.ErrCodeTableNotFoundException {
		return true
	}
	if c == tableNotFoundError {
		return true
	}

	return false
}

// IsTableAlreadyExists asserts bucket already exists error from upstream's
// API code.
func IsTableAlreadyExists(err error) bool {
	aerr, ok := err.(awserr.Error)
	if !ok {
		return false
	}
	if aerr.Code() == dynamodb.ErrCodeTableAlreadyExistsException {
		return true
	}

	return false
}

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongType asserts wrongTypeError.
func IsWrongType(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
