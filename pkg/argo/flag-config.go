package argo

import "github.com/foxcapades/argonaut/internal/flag/fextras"

type FlagConfig struct {
	// ShortFormValidator defines a function that will be used to validate flag
	// short-forms.
	//
	// If the function returns `false` for a given value, an error will be
	// returned from the attempt to build the target FlagSpec.
	//
	// This function is only expected to validate the flag short-form with the
	// given context, additional validations such as short-form flag uniqueness is
	// checked by Argonaut.
	ShortFormValidator func(byte, Config) bool

	// ShortFormPrefix defines the leader character that must precede single
	// short-form flags in a CLI call.
	ShortFormPrefix byte

	ShortFormValueSeparator byte

	// EnableShortFormNoSeparator configures whether CLI parsing should allow for
	// short-form flags with argument values that follow immediately with no
	// separation, e.g. `-aValue`.
	//
	// If set to `true`, short-form flags may have argument values that
	// immediately follow the short-form name.
	//
	// If set to `false`, short-form flags may only take arguments that are first
	// preceded by one of the characters defined in the ShortFormValueSeparators
	// string.
	//
	// Examples:
	//
	//   config.ShortFormValueSeparators = "="
	//
	//   Form     | Enabled | Disabled
	//   -aValue  | OK      | Error
	//   -a=Value | OK      | OK
	EnableShortFormNoSeparator bool

	// TODO: enables "-f Value"
	EnableShortFormSpaceSeparation bool

	EnableShortFormGrouping bool

	// LongFormValidator defines a function that will be used to validate flag
	// long-forms.
	//
	// If the function returns `false` for a given value, an error will be
	// returned from the attempt to build the target FlagSpec.
	//
	// This function is only expected to validate the flag long-form with the
	// given context, additional validations such as long-form flag uniqueness is
	// checked by Argonaut.
	LongFormValidator func(string, Config) bool

	// TODO: handle single-char long-form prefix that matches short-form prefix.
	LongFormPrefix string

	LongFormValueSeparator byte

	EnableLongFormSpaceSeparation bool
}

func DefaultFlagConfig() FlagConfig {
	return FlagConfig{
		ShortFormValidator:             fextras.ShortFormValidator,
		ShortFormPrefix:                '-',
		ShortFormValueSeparator:        '=',
		EnableShortFormNoSeparator:     true,
		EnableShortFormSpaceSeparation: true,
		EnableShortFormGrouping:        true,

		LongFormValidator:             fextras.LongFormValidator,
		LongFormPrefix:                "--",
		LongFormValueSeparator:        '=',
		EnableLongFormSpaceSeparation: true,
	}
}
