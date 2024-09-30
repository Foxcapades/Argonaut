package argument

import (
	"github.com/foxcapades/argonaut/internal/cereal"
	"github.com/foxcapades/argonaut/internal/util/xerr"
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Spec[T any] struct {
	name string

	summary string

	description string

	isRequired bool

	preValidators []argo.PreParseArgumentValidator

	postValidators []argo.PostParseArgumentValidator[T]

	defaultProvider func() (T, error)

	deserializer cereal.Deserializer[T]

	consumer argo.ArgumentValueConsumer[T]

	// 0 = nothing, 1 = parse failed, 2 = has value
	parseState uint8

	parsedValue T
}

func (a Spec[T]) Name() string { return a.name }

func (a Spec[T]) Summary() string { return a.summary }

func (a Spec[T]) Description() string { return a.description }

func (a Spec[T]) HasHelpText() bool { return len(a.summary) > 0 }

func (a Spec[T]) IsRequired() bool { return a.isRequired }

func (a Spec[T]) HasValue() bool { return a.parseState == 2 }

func (a Spec[T]) Value() any {
	if a.parseState == 2 {
		return a.parsedValue
	}

	panic("argument has no value, yet Value() was called")
}

func (a Spec[T]) PreValidate(input string) error {
	var err error

	errors := xerr.NewMultiError()

	for i := range a.preValidators {
		if err = a.preValidators[i].Validate(input, err); err != nil {
			errors.Append(err)
		}
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (a *Spec[T]) ProcessInput(value string) error {
	errors := xerr.NewMultiError()

	// If we've previously encountered an error trying to parse values for this
	// argument, then bail here.
	if a.parseState == 1 {
		return nil
	}

	var prev *T
	if a.parseState == 2 {
		prev = &a.parsedValue
	}

	parsed, err := a.deserializer.Deserialize(value, prev)
	a.parseState = 1

	if err != nil {
		errors.Append(err)
		return errors
	}

	for i := range a.postValidators {
		if err = a.postValidators[i].Validate(parsed, value, err); err != nil {
			errors.Append(err)
		}
	}

	if !errors.IsEmpty() {
		return errors
	}

	a.parseState = 2
	a.parsedValue = parsed

	return a.consumer.Accept(parsed)
}