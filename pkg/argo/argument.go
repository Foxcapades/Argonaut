package argo

import "reflect"

// Argument represents a positional or flag argument that may be attached
// directly to a Command or CommandLeaf, or may be attached to a Flag.
type Argument interface {

	// Name returns the custom name assigned to this Argument.
	//
	// If no custom name was assigned to this Argument when it was built, this
	// method will return an empty string.
	Name() string

	// HasName tests whether this Argument has a custom name assigned.
	HasName() bool

	Binding() any

	HasBinding() bool

	BindingType() reflect.Type

	Default() any

	HasDefault() bool

	DefaultType() reflect.Type

	// Description returns the description attached to this Argument.
	//
	// If no description was attached to this Argument when it was built, this
	// method will return an empty string.
	Description() string

	// HasDescription tests whether this Argument has a description attached.
	HasDescription() bool

	// WasHit tests whether this Argument was hit in a CLI call.
	//
	// If this method returns true, then a value has been assigned to this
	// Argument.
	//
	// If this method returns false, then no value has been assigned to this
	// Argument.
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

	// SetValue sets the value for this Argument to the given string.  The
	// Argument instance is expected to unmarshal this value into the type
	// expected by an assigned binding if one is present.
	SetValue(rawValue string) error

	// SetDefault sets the value for this Argument to its default value.
	//
	// If this argument has no default value, and/or has no binding, this method
	// does nothing.
	SetDefault() error
}
