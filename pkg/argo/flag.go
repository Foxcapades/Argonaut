package argo

type FlagCallback = func(flag FlagSpec)

func SimpleFlagCallback(fn func()) FlagCallback {
	return func(FlagSpec) { fn() }
}

// region Flag

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

// endregion Flag

// region FlagSpec

// FlagSpec represents a successfully built CLI call flag specification.
//
// This intermediate form between FlagSpecBuilder and Flag is used internally by
// Argonaut when parsing CLI input, or printing help text.
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

	// Description returns the longform description of this FlagSpec.
	//
	// Callers may test whether the FlagSpec has a description by calling
	// HasDescription.
	Description() string

	// HasDescription returns an indicator as to whether this FlagSpec has a
	// description set.
	HasDescription() bool

	// IsRequired returns an indicator as to whether this FlagSpec has been marked
	// as required.
	IsRequired() bool

	// WasUsed returns an indicator as to whether this FlagSpec has been marked as
	// having been used at least once in the parsed CLI call.
	WasUsed() bool

	// UsageCount returns a count of the number of times that this FlagSpec was
	// encountered in the parsed CLI call.
	UsageCount() uint32

	// MarkUsed marks the current FlagSpec as having been encountered in the CLI
	// call.
	//
	// This method will result in any immediate FlagCallback instances configured
	// on the current FlagSpec to be called.
	//
	// MarkUsed should be called AFTER providing an argument value when parsing as
	// the value of the argument may be required by immediate FlagCallback
	// instances.
	//
	// This method may be called multiple times if the flag is used more than once
	// in the CLI call.
	MarkUsed()

	// HasExplicitArgument return an indicator as to whether the FlagSpec was
	// configured with an ArgumentSpec or if a fallback was used.
	//
	// A value of `true` indicates that the FlagSpec was given a custom
	// ArgumentSpec.
	//
	// A value of `false` indicates that the FlagSpec is using a fallback
	// ArgumentSpec.
	HasExplicitArgument() bool

	// Argument returns the ArgumentSpec instance configured on this FlagSpec.
	//
	// If no Argument was explicitly set for the FlagSpec when it was built, the
	// ArgumentSpec returned here will be a fallback instance.
	//
	// Callers can test whether the FlagSpec has a custom ArgumentSpec by using
	// HasExplicitArgument.
	Argument() ArgumentSpec

	// ToFlag drops all unnecessary meta and descriptive data returning a Flag
	// instance containing data parsed from the CLI call.
	ToFlag() Flag
}

// endregion FlagSpec

// region FlagSpecBuilder

// TODO: Define an expected behavior for flags without explicitly defined
//       arguments.  Perhaps they should have a "magic" catch-all argument that
//       can be used as either a toggle or a container for raw values (when the
//       cli context makes it unequivocally clear that the argument belongs to
//       the flag, such as "--foo=bar")

type FlagSpecBuilder interface {

	// WithLongForm sets the long-form flag name that the flag may be referenced
	// by on the CLI.
	//
	// Long-form flags consist of one or more characters preceded immediately by
	// two prefix/leader characters (typically dashes).
	//
	// Long-form flags must start with an alphanumeric character and may only
	// consist of alphanumeric characters, dashes, and/or underscores.
	WithLongForm(name string) FlagSpecBuilder

	// WithShortForm sets the short-form flag character that the flag may be
	// referenced by on the CLI.
	//
	// Short-form flags consist of a single character preceded by either a
	// prefix/leader character, or by one or more other short-form flags.
	WithShortForm(name byte) FlagSpecBuilder

	// WithDescription sets a description for the flag that will be used when
	// printing command help text.
	WithDescription(description string) FlagSpecBuilder

	// WithArgument sets a custom argument for the flag that will handle arguments
	// passed for the flag in CLI calls.
	WithArgument(arg ArgumentSpecBuilder) FlagSpecBuilder

	// WithLazyCallback sets a callback on the FlagSpec being built that will be
	// executed on successful parsing of a valid CLI call that included this flag.
	//
	// Callbacks are executed in the order they are configured.
	WithLazyCallback(callback FlagCallback) FlagSpecBuilder

	// WithImmediateCallback sets a callback on the FlagSpec being built that will
	// be executed as soon as the flag is encountered, every time it is
	// encountered while parsing the CLI call, regardless of whether the CLI call
	// was valid.
	//
	// Immediate callbacks will be given the flag state as it was when the flag
	// was encountered, which may not include the final usage count if more flag
	// usages appear later in the CLI call.
	//
	// Callbacks are executed in the order they are configured.
	//
	// If an unrecoverable error is encountered while parsing a CLI call, Argonaut
	// will halt parsing and no further callbacks will be executed regardless of
	// whether there are more flag instances in the CLI call text.
	WithImmediateCallback(callback FlagCallback) FlagSpecBuilder

	// Require marks the flag being built as being required.
	//
	// If the flag is marked as required and is not present in a CLI call, an
	// error will be generated when parsing the CLI input.
	Require() FlagSpecBuilder

	// Build builds a new FlagSpec instance from the configuration set on this
	// FlagSpecBuilder instance.
	Build(config Config) (FlagSpec, error)
}

// endregion FlagSpecBuilder

// region FlagGroup

// FlagGroup represents an organizational, named grouping of CLI flags.
type FlagGroup interface {
	// Name returns the name of the FlagGroupSpec.
	Name() string

	// Flags returns a slice containing all the FlagSpec instances contained by this
	// FlagGroupSpec.
	Flags() []Flag

	// FindFlagByShortForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a short-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByShortForm(name byte) Flag

	// FindFlagByLongForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a long-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByLongForm(name string) Flag
}

// endregion FlagGroup

// region FlagGroupSpec

type FlagGroupSpec interface {

	// Name returns the name of the FlagGroupSpec.
	Name() string

	// Description returns the description text for this FlagGroupSpec.
	Description() string

	// HasDescription returns an indicator as to whether this FlagGroupSpec has
	// description text set.
	HasDescription() bool

	// Flags returns a slice containing all the FlagSpec instances contained by this
	// FlagGroupSpec.
	Flags() []FlagSpec

	// FindFlagByShortForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a short-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByShortForm(name byte) FlagSpec

	// FindFlagByLongForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a long-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByLongForm(name string) FlagSpec

	ToFlagGroup() FlagGroup
}

// endregion FlagGroupSpec

//

// FlagGroupSpecBuilder defines a builder type which is used to configure and
// construct FlagGroupSpec instances.
type FlagGroupSpecBuilder interface {
	// WithName configures the name for the output flag group.
	WithName(name string) FlagGroupSpecBuilder

	// WithDescription configures the help text description of the flag group.
	WithDescription(description string) FlagGroupSpecBuilder

	// WithFlag appends a flag to the flag group.
	WithFlag(flag FlagSpecBuilder) FlagGroupSpecBuilder

	// Build attempts to build a new FlagGroupSpec instance from the current
	// configured state of the parent FlagGroupSpecBuilder.
	Build(config Config) (FlagGroupSpec, error)
}

//

// FlagGroupContainer defines the functionality common to all CLI component
// types that may have FlagGroup instances attached.
type FlagGroupContainer interface {
	// FlagGroups returns a slice containing all the FlagGroup instances attached
	// to this CLI component.
	FlagGroups() []FlagGroup

	// FindFlagGroup attempts to find a FlagGroup instance with a name matching
	// the given input.
	FindFlagGroup(name string) FlagGroup

	Flags() []Flag

	// FindFlagByShortForm attempts to find a Flag instance with a short-form name
	// matching the given input.
	FindFlagByShortForm(name byte) Flag

	// FindFlagByLongForm attempts to find a Flag instance with a long-form name
	// matching the given input.
	FindFlagByLongForm(name string) Flag
}

type FlagGroupSpecContainer interface {
	FlagGroups() []FlagGroupSpec

	FindFlagGroup(name string) FlagGroupSpec

	Flags() []FlagSpec

	FindFlagByShortForm(name byte) FlagSpec

	FindFlagByLongForm(name string) FlagSpec
}
