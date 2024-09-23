package arg

import (
	"github.com/foxcapades/argonaut/internal/arg/argo"
	"github.com/foxcapades/argonaut/internal/errs"
)

type ArgumentSpec[T any] struct {
	isRequired bool

	// 0 = nothing, 1 = parse failed, 2 = has value
	parseState uint8

	parsedValue any

	PreValidators  []argo.PreParseArgumentValidator
	PostValidators []argo.PostParseArgumentValidator
	Deserializer   argo.ValueUnmarshaler[T]
}

func (a ArgumentSpec[T]) IsRequired() bool {
	return a.isRequired
}

func (a ArgumentSpec[T]) HasValue() bool {
	return a.parseState == 2
}

func (a ArgumentSpec[T]) Value() any {
	if a.parseState == 2 {
		return a.parsedValue
	}

	panic("argument has no value, yet Value() was called")
}

func (a ArgumentSpec[T]) PreValidate(input string) error {
	var err error

	errors := errs.NewMultiError()

	for i := range a.PreValidators {
		if err = a.PreValidators[i].Validate(input, err); err != nil {
			errors.Append(err)
		}
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (a *ArgumentSpec[T]) ProcessInput(value string) error {
	errors := errs.NewMultiError()
	parsed, err := a.Deserializer.Unmarshal(value)
	a.parseState = 1

	if err != nil {
		errors.Append(err)
		return errors
	}

	for i := range a.PostValidators {
		if err = a.PostValidators[i].Validate(parsed, value, err); err != nil {
			errors.Append(err)
		}
	}

	if errors.IsEmpty() {
		a.parseState = 2
		a.parsedValue = parsed
		return nil
	}

	return errors
}
