package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrRepository       = errors.New("repository error")
	ErrCache            = errors.New("cache error")
	ErrPackSizesEmpty   = errors.New("pack sizes cannot be empty")
	ErrItemsInvalid     = errors.New("items must be greater than 0")
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
