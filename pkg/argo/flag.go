package argo

// Flag represents a single CLI flag which may have an argument.
type Flag interface {

	// ShortForm returns the short-form character representing this Flag.
	ShortForm() byte

	// HasShortForm indicates whether this Flag has a short-form character.
	HasShortForm() bool

	// LongForm returns the long-form string representing this Flag.
	LongForm() string

	// HasLongForm indicates whether this Flag has a long-form string.
	HasLongForm() bool

	// Description returns the help-text description of this flag.
	Description() string

	// HasDescription indicates whether this Flag has a help-text description.
	HasDescription() bool

	// Argument returns the argument value attached to this Flag.
	//
	// If this Flag does not have an Argument attached, this method will return
	// nil.
	Argument() Argument

	// HasArgument indicates whether this Flag accepts an argument.
	HasArgument() bool

	// IsRequired indicates whether this Flag is required.
	//
	// Required flags must be present in the CLI call.
	IsRequired() bool

	// WasHit indicates whether this flag was used in the CLI call.
	WasHit() bool

	// HitCount returns the number of times that this flag was used in the CLI
	// call.
	HitCount() int

	IncrementHitCount()

	IsHelpFlag() bool

	HasCallback() bool

	Callback() FlagCallback
}
