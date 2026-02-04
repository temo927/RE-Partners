package errors

import (
	"errors"
	"fmt"
)

const (
	MaxPackSize = 2147483647
	MaxItems    = 2147483647
	MinPackSize = 1
	MinItems    = 1
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrRepository          = errors.New("repository error")
	ErrCache               = errors.New("cache error")
	ErrPackSizesEmpty      = errors.New("pack sizes cannot be empty")
	ErrItemsInvalid        = errors.New("items must be greater than 0")
	ErrPackSizeOutOfRange  = errors.New("pack size is out of range (must be between 1 and 2147483647)")
	ErrItemsOutOfRange     = errors.New("items value is out of range (must be between 1 and 2147483647)")
	ErrDuplicatePackSizes  = errors.New("duplicate pack sizes are not allowed")
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func WrapDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

func WrapWithDomain(err error, domainErr error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w: %w", message, domainErr, err)
}
