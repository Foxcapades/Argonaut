package argo

// FlagSpec represents a successfully built CLI call flag specification.
//
// The FlagSpec is used by Argonaut when parsing CLI input or rendering help
// text.  After that point, non input data such as summaries or names will be
// released.
type FlagSpec interface {
	// LongForm returns the long-form name assigned to the current FlagSpec.  If
	// the FlagSpec does not have a long-form name, the return value will be an
	// empty string.
	//
	// Callers may test whether the FlagSpec has a long-form name by using
	// HasLongForm.
	LongForm() string

	// HasLongForm returns an indicator as to whether the current FlagSpec has a
	// long-form name assigned.
	HasLongForm() bool

	// ShortForm returns the short-form name assigned to the current FlagSpec.  If
	// the FlagSpec does not have a short-form name, the return value will be a
	// NULL byte (value `0`).
	//
	// Callers may test whether the FlagSpec has a short-form name by using
	// HasShortForm.
	ShortForm() byte

	// HasShortForm returns an indicator as to whether the current FlagSpec has a
	// short-form name assigned.
	HasShortForm() bool

	// Summary returns a short summary text describing this FlagSpec.
	//
	// If this FlagSpec has no help text available, this method will return an
	// empty string.
	//
	// Callers may test whether the FlagSpec has help text available by calling
	// HasHelpText.
	Summary() string

	// Description returns the longform description of this FlagSpec.
	//
	// If this FlagSpec has no help text available, this method will return an
	// empty string.
	//
	// Callers may test whether the FlagSpec has help text available by calling
	// HasHelpText.
	Description() string

	// HasHelpText returns an indicator as to whether this FlagSpec has help text
	// available.
	//
	// If the FlagSpec has any help text, both Summary and Description will return
	// a value, though the value may be the same for both.
	HasHelpText() bool

	// IsRequired returns an indicator as to whether this FlagSpec has been marked
	// as required.
	IsRequired() bool

	// WasUsed returns an indicator as to whether this FlagSpec has been marked as
	// having been used at least once in the parsed CLI call.
	WasUsed() bool

	// UsageCount returns a count of the number of times that this FlagSpec was
	// seen in the parsed CLI call.
	UsageCount() uint32

	// MarkUsed marks the current FlagSpec as having been encountered in the CLI
	// call.
	//
	// This method will result in any immediate FlagCallback instances configured
	// on the current FlagSpec to be called.
	//
	// MarkUsed should be called AFTER providing an argument value when parsing
	// as the value of the argument may be required by immediate FlagCallback
	// instances.
	//
	// This method may be called multiple times if the flag is used more than once
	// in the CLI call.
	MarkUsed()

	// HasExplicitArgument return an indicator as to whether the FlagSpec was
	// configured with an ArgumentSpec or if a fallback was used.
	//
	// A value of `true` indicates that the FlagSpec was given a custom
	// ArgumentSpec.  A value of `false` indicates that the FlagSpec is using a
	// fallback ArgumentSpec.
	HasExplicitArgument() bool

	// Argument returns the ArgumentSpec instance configured on this FlagSpec.
	//
	// If no Argument was explicitly set for the FlagSpec when it was built, the
	// ArgumentSpec returned here will be a fallback instance.
	//
	// Callers can test whether the FlagSpec has a custom ArgumentSpec by using
	// HasExplicitArgument.
	Argument() ArgumentSpec
}
