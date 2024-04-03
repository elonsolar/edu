// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsUnknow(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_UNKNOW.String() && e.Code == 500
}

func ErrorUnknow(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_UNKNOW.String(), fmt.Sprintf(format, args...))
}

func IsConcurrencyUpdate(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_CONCURRENCY_UPDATE.String() && e.Code == 506
}

func ErrorConcurrencyUpdate(format string, args ...interface{}) *errors.Error {
	return errors.New(506, ErrorReason_CONCURRENCY_UPDATE.String(), fmt.Sprintf(format, args...))
}
