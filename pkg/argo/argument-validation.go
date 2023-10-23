package argo

import (
	"errors"
)

// OneOfPreParseArgumentValidator builds an argument pre-parse validator
// function that ensures the raw value passed on the command line call matches
// at least one of the values in the given values slice.
//
// If the incoming raw value does not match one of the values in the given
// values slice, the given message will be used to build a new error which will
// be returned.
func OneOfPreParseArgumentValidator(values []string, message string) any {
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
func OneOfPostParseArgumentValidator[T comparable](values []T, message string) any {
	return func(a T, s string) error {
		for i := range values {
			if a == values[i] {
				return nil
			}
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
func NoneOfPreParseArgumentValidator(values []string, message string) any {
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
func NoneOfPostParseArgumentValidator[T comparable](values []T, message string) any {
	return func(a T, s string) error {
		for i := range values {
			if a == values[i] {
				return errors.New(message)
			}
		}

		return nil
	}
}

// NumericValue defines the numeric types that Argonaut can parse.
type NumericValue interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// NumericRangePostParseArgumentValidator builds an argument post-parse
// validator function that ensures the value parsed from the raw command line
// input falls within the given inclusive range.
//
// If the value falls outside the given inclusive range, the given message will
// be used to build a new error which will be returned.
func NumericRangePostParseArgumentValidator[T NumericValue](minimum, maximum T, message string) any {
	return func(a T, s string) error {
		if a < minimum || a > maximum {
			return errors.New(message)
		}

		return nil
	}
}
