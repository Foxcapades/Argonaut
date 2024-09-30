package argo

type ArgumentSpec interface {
	// Name returns the name configured on this ArgumentSpec.
	Name() string

	// Summary returns a short summary text describing this ArgumentSpec.
	//
	// If this ArgumentSpec has no help text available, this method will return an
	// empty string.  This may be tested for using HasHelpText.
	Summary() string

	// Description returns the longform description of this ArgumentSpec.
	//
	// If this ArgumentSpec has no description available, this method will return
	// an empty string.  This may be tested for using HasHelpText.
	Description() string

	// HasHelpText returns an indicator as to whether this ArgumentSpec has help
	// text available.
	//
	// If the ArgumentSpec has any help text, both Summary and Description will
	// return a value, though the value may be the same for both.
	HasHelpText() bool

	// IsRequired returns an indicator as to whether this ArgumentSpec has been
	// marked as required.
	IsRequired() bool

	// HasValue returns whether this ArgumentSpec has had a value provided either
	// from the CLI call or from a default.
	HasValue() bool

	// Value returns the value set on this ArgumentSpec.
	//
	// If no value has been set on this ArgumentSpec, calls to this method will
	// panic.  Whether a value is available should be tested by using HasValue.
	Value() any

	// PreValidate performs any configured pre-parse validation on a raw CLI input
	// string before any attempt is made to parse the value.
	//
	// This method is separate from ProcessInput as it may be used to decide what
	// ArgumentSpec instance should be given a CLI input value when there is
	// ambiguity in the CLI definition.
	PreValidate(rawInput string) error

	// ProcessInput parses the raw CLI input string and performs any post-parse
	// validation on it before setting it as the value for this ArgumentSpec.
	//
	// If this method returns an error, a value parsed from the given input MUST
	// NOT be set on the ArgumentSpec instance.  For example, if HasValue returns
	// `false` before this method is called, and this method returns an error,
	// HasValue MUST still return `false` after this method returns.
	//
	// If this method does not return an error, the value parsed from the given
	// input MUST be used as the value for the ArgumentSpec instance.  Meaning,
	// after this method is called and returns `nil`, HasValue MUST return `true`.
	ProcessInput(rawInput string) error

	ToArgument() Argument
}
