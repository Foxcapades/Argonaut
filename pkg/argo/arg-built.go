package argo

type Argument interface {
	// IsRequired returns an indicator for whether the Argument was marked as
	// required when configured.
	IsRequired() bool

	// UsedDefault returns an indicator for whether the Argument's default value
	// was used.
	//
	// If the returned value is `true`, then no value was provided for this
	// Argument in the CLI call.
	//
	// If the returned value is `false`, then the CLI call provided a value for
	// this Argument.
	UsedDefault() bool

	// HasValue indicates whether this Argument has a value set.
	//
	// If the Argument was configured with a default value, this method will
	// always return true.
	//
	// If this method returns `false`, calls to Value will panic.
	HasValue() bool

	// Value returns the value set for this Argument, if one was set.
	//
	// If no value is set, this method will panic.  Test whether a value is
	// available by calling HasValue before calling this method.
	//
	// Also see ValueOrNil.
	Value() any

	// ValueOrNil returns either the value set for this Argument, or `nil` if this
	// argument does not have a value.
	//
	// If HasValue returns `false`, this method will return `nil`.  If HasValue
	// returns `true`, this method will return the set value, which may also be
	// `nil`.
	//
	// Also see Value.
	ValueOrNil() any
}
