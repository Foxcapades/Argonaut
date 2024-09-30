package argo

// Flag represents a single CLI flag definition which may or may not have been
// used in the CLI call.
type Flag interface {

	// LongForm returns the long-form name assigned to the current Flag.  If the
	// Flag does not have a long-form name, the return value will be an empty
	// string.
	//
	// Callers may test whether the Flag has a long-form name by using
	// HasLongForm.
	LongForm() string

	// HasLongForm returns an indicator as to whether the current Flag has a
	// long-form name assigned.
	HasLongForm() bool

	// ShortForm returns the short-form name assigned to the current Flag.  If the
	// FlagSpec does not have a short-form name, the return value will be a NULL
	// byte (value `0`).
	//
	// Callers may test whether the FlagSpec has a short-form name by using
	// HasShortForm.
	ShortForm() byte

	// HasShortForm returns an indicator as to whether the current Flag has a
	// short-form name assigned.
	HasShortForm() bool

	// IsRequired returns an indicator as to whether this Flag has been marked as
	// required.
	IsRequired() bool

	// WasUsed returns an indicator as to whether this FlagSpec has been marked as
	// having been used at least once in the parsed CLI call.
	WasUsed() bool

	// UsageCount returns a count of the number of times that this Flag was
	// encountered in the parsed CLI call.
	UsageCount() uint32

	// HasExplicitArgument return an indicator as to whether the Flag was
	// configured with an Argument or if a fallback was used.
	//
	// A value of `true` indicates that the Flag was given a custom Argument.
	//
	// A value of `false` indicates that the Flag is using a fallback Argument.
	HasExplicitArgument() bool

	// Argument returns the Argument instance configured on this Flag.
	//
	// If no Argument was explicitly set for the Flag when it was built, the
	// Argument returned here will be a fallback instance.
	//
	// Callers can test whether the Flag has a custom Argument by using
	// HasExplicitArgument.
	Argument() Argument
}
