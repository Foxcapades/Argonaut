package argument

import "github.com/foxcapades/argonaut/pkg/argo"

type FallbackArgumentSpec struct {
	hasValue bool
	value    string
}

func (m FallbackArgumentSpec) Name() string { return "" }

func (m FallbackArgumentSpec) Summary() string { return "" }

func (m FallbackArgumentSpec) Description() string { return "" }

func (m FallbackArgumentSpec) HasHelpText() bool { return false }

func (m FallbackArgumentSpec) IsRequired() bool { return false }

func (m FallbackArgumentSpec) HasValue() bool { return m.hasValue }

func (m FallbackArgumentSpec) Value() any { return m.value }

func (m FallbackArgumentSpec) PreValidate(string) error { return nil }

func (m *FallbackArgumentSpec) ProcessInput(rawInput string) error {
	m.hasValue = true
	m.value = rawInput
	return nil
}

func (m FallbackArgumentSpec) ToArgument() argo.Argument {
	return FallbackArgument{m.hasValue, m.value}
}

type FallbackArgument struct {
	hasValue bool
	value    string
}

func (f FallbackArgument) IsRequired() bool { return false }

func (f FallbackArgument) UsedDefault() bool { return false }

func (f FallbackArgument) HasValue() bool { return f.hasValue }

func (f FallbackArgument) Value() any {
	if f.hasValue {
		return f.value
	}
	panic("attempted to get the value of an argument that has none")
}

func (f FallbackArgument) ValueOrNil() any {
	if f.hasValue {
		return f.value
	}
	return nil
}
