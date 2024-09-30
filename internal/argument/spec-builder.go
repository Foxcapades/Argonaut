package argument

import (
	"errors"
	"reflect"
	"strings"

	"github.com/foxcapades/argonaut/internal/cereal"
	"github.com/foxcapades/argonaut/internal/util/xerr"
	"github.com/foxcapades/argonaut/internal/util/xstr"
	"github.com/foxcapades/argonaut/pkg/argo"
)

func NewBuilder[T any]() *SpecBuilder[T] {
	return new(SpecBuilder[T])
}

type SpecBuilder[T any] struct {
	hasName bool

	hasDescription bool

	hasSummary bool

	required bool

	hasDefault uint8 // 0 = no, 1 = raw, 2 = typed

	name string

	description string

	summary string

	defaultValue any

	preValidators  []argo.PreParseArgumentValidator
	postValidators []argo.PostParseArgumentValidator[T]

	consumers []argo.ArgumentValueConsumer[T]

	unmarshaler argo.ValueUnmarshaler[T]

	errors xerr.MultiError
}

// region Configuration

func (a *SpecBuilder[T]) WithName(name string) argo.TypedArgumentSpecBuilder[T] {
	a.name = strings.TrimSpace(name)
	a.hasName = len(a.name) > 0
	return a
}

func (a *SpecBuilder[T]) WithDescription(desc string) argo.TypedArgumentSpecBuilder[T] {
	a.description = strings.TrimSpace(desc)
	a.hasDescription = len(a.description) > 0
	return a
}

func (a *SpecBuilder[T]) WithSummary(summary string) argo.TypedArgumentSpecBuilder[T] {
	a.summary = strings.TrimSpace(summary)
	a.hasSummary = len(a.summary) > 0
	return a
}

func (a *SpecBuilder[T]) WithRawDefault(value string) argo.TypedArgumentSpecBuilder[T] {
	a.defaultValue = value
	a.hasDefault = 1
	return a
}

func (a *SpecBuilder[T]) WithDefault(value T) argo.TypedArgumentSpecBuilder[T] {
	a.defaultValue = func() T { return value }
	a.hasDefault = 2
	return a
}

func (a *SpecBuilder[T]) WithDefaultProvider(provider func() T) argo.TypedArgumentSpecBuilder[T] {
	a.defaultValue = provider
	a.hasDefault = 2
	return a
}

func (a *SpecBuilder[T]) WithPreValidator(validator argo.PreParseArgumentValidator) argo.TypedArgumentSpecBuilder[T] {
	a.preValidators = append(a.preValidators, validator)
	return a
}

func (a *SpecBuilder[T]) WithPostValidator(validator argo.PostParseArgumentValidator[T]) argo.TypedArgumentSpecBuilder[T] {
	a.postValidators = append(a.postValidators, validator)
	return a
}

func (a *SpecBuilder[T]) WithValueConsumer(consumer argo.ArgumentValueConsumer[T]) argo.TypedArgumentSpecBuilder[T] {
	a.consumers = append(a.consumers, consumer)
	return a
}

func (a *SpecBuilder[T]) WithBinding(binding *T) argo.TypedArgumentSpecBuilder[T] {
	a.consumers = append(a.consumers, NewPointerConsumer(binding))
	return a
}

func (a *SpecBuilder[T]) WithDeepBinding(binding any) argo.TypedArgumentSpecBuilder[T] {
	if con, err := NewMagicPointerConsumer[T](binding); err != nil {
		if a.errors == nil {
			a.errors = xerr.NewMultiError()
		}

		a.errors.Append(err)
	} else {
		a.consumers = append(a.consumers, con)
	}

	return a
}

func (a *SpecBuilder[T]) WithUnmarshaler(unmarshaler argo.ValueUnmarshaler[T]) argo.TypedArgumentSpecBuilder[T] {
	a.unmarshaler = unmarshaler
	return a
}

func (a *SpecBuilder[T]) Required() argo.TypedArgumentSpecBuilder[T] {
	a.required = true
	return a
}

// endregion Configuration

func (a SpecBuilder[T]) Build(config argo.Config) (argo.ArgumentSpec, error) {
	if a.errors != nil && !a.errors.IsEmpty() {
		return nil, a.errors
	}

	if a.required && a.hasDefault != 0 {
		// TODO: newtype this
		return nil, errors.New(argo.ErrMsgArgumentDefaultAndRequired)
	}

	spec := new(Spec[T])

	spec.isRequired = a.required
	spec.preValidators = a.preValidators
	spec.postValidators = a.postValidators

	switch len(a.consumers) {
	case 0:
		spec.consumer = VoidConsumer[T]{}
	case 1:
		spec.consumer = a.consumers[0]
	default:
		spec.consumer = NewMultiConsumer(a.consumers)
	}

	if a.unmarshaler == nil {
		spec.deserializer = cereal.NewMagicWrapper[T](config.DefaultUnmarshaler)
	} else {
		spec.deserializer = cereal.NewCustomWrapper(a.unmarshaler)
	}

	switch a.hasDefault {
	case 1:
		spec.defaultProvider = func() (T, error) { return spec.deserializer.Deserialize(a.defaultValue.(string), nil) }
	case 2:
		spec.defaultProvider = func() (T, error) { return a.defaultValue.(func() T)(), nil }
	}

	if a.hasName {
		spec.name = a.name
	} else {
		var tmp T
		spec.name = reflect.TypeOf(tmp).Name()
	}

	if a.hasDescription {
		spec.description = a.description

		if a.hasSummary {
			spec.summary = a.summary
		} else {
			idx := xstr.IndexOfAnyWithin(a.description, "\r\n", 128)

			if idx == -1 {
				spec.summary = xstr.Truncate(a.description, 128)
			} else {
				spec.summary = a.description[:idx] + "..."
			}
		}
	} else if a.hasSummary {
		spec.summary = a.summary
		spec.description = a.summary
	}

	return spec, nil
}
