package validate

import (
	"encoding/json"
	"errors"
)

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Tag     string `json:"tag"`
	Message string `json:"error"`
	Value   any    `json:"value"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

func (fe FieldErrors) Nil() bool {
	return len(fe) == 0
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

func NewFieldErrors(errs ...FieldError) FieldErrors {
	return errs
}

func IsFieldError(err error) bool {
	v := FieldErrors{}
	return errors.As(err, &v)
}

type InvalidRouteKeyError struct {
	key string
	msg string
}

func (err InvalidRouteKeyError) Error() string {
	msg := "invalid route key: " + err.key
	if err.msg != "" {
		msg += " - " + err.msg
	}
	return msg
}

func NewRouteKeyError(key string) error {
	return &InvalidRouteKeyError{key: key}
}

func NewRouteKeyErrorWithMessage(key, msg string) error {
	return &InvalidRouteKeyError{key: key, msg: msg}
}

func IsInvalidRouteKeyError(err error) bool {
	var re *InvalidRouteKeyError
	return errors.As(err, &re)
}
