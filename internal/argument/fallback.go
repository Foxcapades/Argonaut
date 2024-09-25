package argument

type FallbackArgument struct {
	hasValue bool
	value    string
}

func (m FallbackArgument) Name() string { return "" }

func (m FallbackArgument) Summary() string { return "" }

func (m FallbackArgument) Description() string { return "" }

func (m FallbackArgument) HasHelpText() bool { return false }

func (m FallbackArgument) IsRequired() bool { return false }

func (m FallbackArgument) HasValue() bool { return m.hasValue }

func (m FallbackArgument) Value() any { return m.value }

func (m FallbackArgument) PreValidate(string) error { return nil }

func (m *FallbackArgument) ProcessInput(rawInput string) error {
	m.hasValue = true
	m.value = rawInput
	return nil
}
