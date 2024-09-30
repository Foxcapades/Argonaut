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
