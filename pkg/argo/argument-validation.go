package argo

import (
	"errors"
	"fmt"
	"reflect"
)

// OneOfPreParseArgumentValidator builds an argument pre-parse validator
// function that ensures the raw value passed on the command line call matches
// at least one of the values in the given values slice.
//
// If the incoming raw value does not match one of the values in the given
// values slice, the given message will be used to build a new error which will
// be returned.
func OneOfPreParseArgumentValidator(values []string, message string) ArgumentPreParseValidatorFn {
	return func(s string) error {
		for _, value := range values {
			if s == value {
				return nil
			}
		}

		return errors.New(message)
	}
}

// OneOfPostParseArgumentValidator builds an argument post-parse validator
// function that ensures the value parsed from the raw command line input
// matches at least one of the values in the given values slice.
//
// If the parsed value does not match one of the values in the given values
// slice, the given message will be used to build a new error which will be
// returned.
func OneOfPostParseArgumentValidator[T comparable](values []T, message string) ArgumentPostParseValidatorFn {
	return func(a any, s string) error {
		if val, ok := a.(T); ok {
			for i := range values {
				if val == values[i] {
					return nil
				}
			}
		} else {
			return fmt.Errorf("cannot convert the target argument binding value into a value of type %s", reflect.TypeOf(values).Elem())
		}

		return errors.New(message)
	}
}

// NoneOfPreParseArgumentValidator builds an argument pre-parse validator
// function that ensures the raw value passed on the command line call does not
// match any of the values in the given values slice.
//
// If the parsed value does match one of the values in the given values slice,
// the given message will be used to build a new error which will be returned.
func NoneOfPreParseArgumentValidator(values []string, message string) ArgumentPreParseValidatorFn {
	return func(s string) error {
		for _, value := range values {
			if s == value {
				return errors.New(message)
			}
		}

		return nil
	}
}

// NoneOfPostParseArgumentValidator builds an argument post-parse validator
// function that ensures the value parsed from the raw command line input does
// not match any of the values in the given values slice.
//
// If the parsed value does match one of the values in the given values slice,
// the given message will be used to build a new error which will be returned.
func NoneOfPostParseArgumentValidator[T comparable](values []T, message string) ArgumentPostParseValidatorFn {
	return func(a any, s string) error {
		if val, ok := a.(T); ok {
			for i := range values {
				if val == values[i] {
					return errors.New(message)
				}
			}
		} else {
			return fmt.Errorf("cannot convert the target argument binding value into a value of type %s", reflect.TypeOf(values).Elem())
		}

		return nil
	}
}
