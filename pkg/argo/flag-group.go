package argo

// A FlagGroup is a named group or category for a collection of one or more
// Flag instances.
//
// Flag groups are primarily used to categorize CLI flags when rendering help
// text.
//
// Flag groups that do not contain any flags will be filtered out at build time.
type FlagGroup interface {

	// Name of the FlagGroup.
	//
	// This value is used when rendering help text.
	Name() string

	// Description for this FlagGroup.
	//
	// This optional value is used when rendering help text.
	Description() string

	// HasDescription indicates whether this FlagGroup has a non-empty description
	// attached.
	HasDescription() bool

	// Flags returns the Flag instances contained by this FlagGroup.
	Flags() []Flag

	// FindShortFlag checks this FlagGroup's contents for a Flag instance with the
	// given value as its short form.
	//
	// If one could not be found, this method returns nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag checks this FlagGroup's contents for a Flag instance with the
	// given value as its long form.
	//
	// If one could not be found, this method returns nil.
	FindLongFlag(name string) Flag

	// Size returns the number of Flag instances attached to this FlagGroup.
	Size() int
}

// A FlagGroupBuilder is used to build a group or category of flags that go
// together.  This is primarily used for rendering help text.
//
// FlagGroup instances are not required but allow for organizing rendered help
// text when there are many flags attached to a command.
//
// Flag groups must have a non-empty Name value and may optionally have
// description text attached.
//
// Flag groups that contain no flags wil not be rendered in help text generated
// by Argonaut.
type FlagGroupBuilder interface {

	// Name returns the name that has been configured on this FlagGroupBuilder
	// instance.
	Name() string

	// WithDescription sets the description text for this FlagGroupBuilder.
	//
	// Description text is optional and is only used by Argonaut when generating
	// help text.
	//
	// An empty string value is treated as if no description value has been set.
	WithDescription(desc string) FlagGroupBuilder

	// HasDescription indicates whether this FlagGroupBuilder has a non-empty
	// description value set.
	HasDescription() bool

	// Description returns the description text set on this FlagGroupBuilder.
	Description() string

	// WithFlag appends the given FlagBuilder instance to this FlagGroupBuilder.
	//
	// When help text is rendered by Argonaut for the built FlagGroup, flags added
	// here will show under the heading of the flag group's name.
	WithFlag(flag FlagBuilder) FlagGroupBuilder

	// WithFlags appends the given FlagBuilder instances to this FlagGroupBuilder.
	//
	// When help text is rendered by Argonaut for the built FlagGroup, flags added
	// here will show under the heading of the flag group's name.
	WithFlags(flag ...FlagBuilder) FlagGroupBuilder

	// HasFlags indicates whether this FlagGroupBuilder contains one or more
	// FlagBuilder instances.
	HasFlags() bool

	// Flags returns the FlagBuilder instances that have been attached to this
	// FlagGroupBuilder via the WithFlag and WithFlags methods.
	Flags() []FlagBuilder

	// Size returns the number of FlagBuilder instances currently attached to this
	// FlagGroupBuilder.
	Size() int
}
