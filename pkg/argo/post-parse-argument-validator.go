package argo

type PostParseArgumentValidator[T any] interface {
	Validate(parsedValue T, rawInput string, prev error) error
}

type PostParseArgumentValidatorFn[T any] func(parsedValue T, rawInput string, prev error) error

func (p PostParseArgumentValidatorFn[T]) Validate(parsedValue T, rawInput string, prev error) error {
	return p(parsedValue, rawInput, prev)
}

func SimplePostParseArgumentValidatorFn[T any](fn func(parsedValue T) error) PostParseArgumentValidator[T] {
	return PostParseArgumentValidatorFn[T](func(parsedValue T, _ string, _ error) error { return fn(parsedValue) })
}
