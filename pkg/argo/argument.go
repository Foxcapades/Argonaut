package argo

// Argument represents a positional or flag argument that may be attached
// directly to a Command or LeafCommand, or may be attached to a Flag.
type Argument interface {

	// HasName tests whether this Argument has a custom name assigned.
	HasName() bool

	// Name returns the custom name assigned to this Argument.
	//
	// If no custom name was assigned to this Argument when it was built, this
	// method will return an empty string.
	Name() string

	// HasDefault indicates whether a default value has been set on this
	// Argument.
	HasDefault() bool

	// Default returns the default value or value provider attached to this
	// Argument, if such a value exists.
	//
	// If this Argument does not have a default value or provider set, this method
	// will return nil.
	Default() any

	// HasDescription tests whether this Argument has a description attached.
	HasDescription() bool

	// Description returns the description attached to this Argument.
	//
	// If no description was attached to this Argument when it was built, this
	// method will return an empty string.
	Description() string

	// WasHit tests whether this Argument was hit in a CLI call.
	//
	// This does not necessarily indicate that there is no value available for
	// this argument, just that it wasn't hit in the CLI call.  If the argument
	// had a default value provided, it will have been set in that case.
	WasHit() bool

	// RawValue returns the raw text value that was assigned to this Argument in
	// the CLI call.
	//
	// If this Argument was not hit during the CLI call, this method will return
	// an empty string.  This empty string IS NOT an indicator whether this
	// Argument was hit, as it may have been intentionally assigned an empty
	// value.  To test whether the Argument was hit, use WasHit.
	RawValue() string

	// IsRequired returns whether this Argument is required by its parent CLI
	// component.
	//
	// When parsing the CLI, if this argument is not found, an error will be
	// returned.
	IsRequired() bool

	// HasBinding indicates whether this Argument has a value binding.
	HasBinding() bool

	Binding() ArgumentBinding

	SetValue(rawValue string) error

	SetToDefault() error
}
