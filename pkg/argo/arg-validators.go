package argo

type PreParseArgumentValidator interface {
	Validate(rawInput string, prev error) error
}

type PreParseArgumentValidatorFn func(rawInput string, prev error) error

func (p PreParseArgumentValidatorFn) Validate(rawInput string, prev error) error {
	return p(rawInput, prev)
}

func SimplePreParseArgumentValidatorFn(fn func(rawInput string) error) PreParseArgumentValidator {
	return PreParseArgumentValidatorFn(func(rawInput string, _ error) error { return fn(rawInput) })
}

type PostParseArgumentValidator interface {
	Validate(parsedValue any, rawInput string, prev error) error
}

type PostParseArgumentValidatorFn func(parsedValue any, rawInput string, prev error) error

func (p PostParseArgumentValidatorFn) Validate(parsedValue any, rawInput string, prev error) error {
	return p(parsedValue, rawInput, prev)
}

func SimplePostParseArgumentValidatorFn(fn func(parsedValue any) error) PostParseArgumentValidator {
	return PostParseArgumentValidatorFn(func(parsedValue any, _ string, _ error) error { return fn(parsedValue) })
}
