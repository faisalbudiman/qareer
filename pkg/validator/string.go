package validator

import (
	"errors"
	"fmt"
	"regexp"
)

type stringValidator struct {
	RequiredError  error
	MinLengthError func(int) error
	LatLongError   error
	IPError        error
}

func NewStringValidator() stringValidator {
	return stringValidator{
		RequiredError: ErrRequired,
		MinLengthError: func(l int) error {
			return errors.New(fmt.Sprintf("min length is %d", l))
		},
		LatLongError: ErrLatLong,
		IPError:      ErrIPError,
	}
}

func ValidateString(key string, value string, validators ...func(string) error) error {
	errors := make(Errors)
	for _, validator := range validators {
		if err := validator(value); err != nil {
			errors[key] = err
			return errors
		}
	}

	return nil
}

func (cfg stringValidator) Required(s string) error {
	if s == "" {
		return cfg.RequiredError
	}

	return nil
}

func (cfg stringValidator) StringMinLength(l int) func(string) error {
	return func(s string) error {
		if len(s) < l {
			return cfg.MinLengthError(l)
		}

		return nil
	}
}

func (cfg stringValidator) StringLatLong(s string) error {
	match, _ := regexp.MatchString(`^(-?\d+(\.\d+)?),\s*(-?\d+(\.\d+)?)$`, s)
	if !match {
		return cfg.LatLongError
	}

	return nil
}

func (cfg stringValidator) StringIP(s string) error {
	match, _ := regexp.MatchString(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`, s)
	if !match {
		return cfg.IPError
	}

	return nil
}
