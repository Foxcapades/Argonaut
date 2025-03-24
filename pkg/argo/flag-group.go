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
	// This value will be used when rendering help text.
	Name() string

	// Description for this FlagGroup.
	//
	// This optional value will be used when rendering help text.
	Description() string

	// HasDescription indicates whether this FlagGroup has a description attached.
	HasDescription() bool

	// Flags returns the Flag instances contained by this FlagGroup.
	Flags() []Flag

	// FindShortFlag checks this FlagGroup's contents for a flag with the given
	// short flag character.  If one could not be found, this method returns nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag checks this FlagGroup's contents for a flag with the given
	// long flag name.  If one could not be found, this method returns nil.
	FindLongFlag(name string) Flag

	Size() int
}

// A FlagGroupBuilder is used to build a group or category of flags that go
// together.  This is primarily used for rendering help text.
//
// FlagGroups are not required but allow for organizing the help text when there
// are many flags attached to a command.
type FlagGroupBuilder interface {
	Name() string

	// WithDescription sets a description value on this FlagGroupBuilder.
	WithDescription(desc string) FlagGroupBuilder

	HasDescription() bool

	Description() string

	// WithFlag appends the given FlagBuilder instance to this FlagGroupBuilder.
	WithFlag(flag FlagBuilder) FlagGroupBuilder

	WithFlags(flag ...FlagBuilder) FlagGroupBuilder

	HasFlags() bool

	Flags() []FlagBuilder

	Size() int
}
