package errors

import "fmt"

type DuplicateError struct {
	message string
}

func (e *DuplicateError) Error() string {
	return e.message
}

func NewDuplicateError(message string) *DuplicateError {
	return &DuplicateError{
		message: message,
	}
}

type NotFoundError struct {
	message string
}

func (e *NotFoundError) Error() string {
	return e.message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		message: message,
	}
}

func NewNotFoundErrorf(message string, args ...interface{}) *NotFoundError {
	message = fmt.Sprintf(message, args)
	return NewNotFoundError(message)
}

type NotAcceptableError struct {
	message string
}

func (e *NotAcceptableError) Error() string {
	return e.message
}

func NewNotAcceptableError(message string) *NotAcceptableError {
	return &NotAcceptableError{
		message: message,
	}
}

type IllegalArgumentError struct {
	message string
}

func (e *IllegalArgumentError) Error() string {
	return e.message
}

func NewIllegalArgumentError(message string) *IllegalArgumentError {
	return &IllegalArgumentError{
		message: message,
	}
}

func NewIllegalArgumentErrorf(message string, args ...interface{}) *IllegalArgumentError {
	message = fmt.Sprintf(message, args)
	return NewIllegalArgumentError(message)
}

type AccessDeniedError struct {
	message string
}

func (e *AccessDeniedError) Error() string {
	return e.message
}

func NewAccessDeniedError(message string) *AccessDeniedError {
	return &AccessDeniedError{
		message: message,
	}
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		message: message,
	}
}

func NewConflictErrorf(message string, args ...interface{}) *ConflictError {
	message = fmt.Sprintf(message, args)
	return NewConflictError(message)
}

type ConflictError struct {
	message string
}

func (e *ConflictError) Error() string {
	return e.message
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
