package argo

type CommandNodeCallback[T any] = func(node T)

// A Node is a base type common to all elements in a CommandTree.
type Node[T any] interface {
	// Name returns the name of the command or subcommand.
	Name() string

	// Description returns the description value assigned to this node.
	//
	// Description values are used when rendering help text.
	Description() string

	// HasDescription indicates whether this Node has a description value
	// set.
	HasDescription() bool

	// FlagGroups returns the flag groups assigned to this Node.
	//
	// This method will only return flag groups that had flags assigned to them,
	// the rest of the flag groups will have been filtered out when the node was
	// built.
	FlagGroups() []FlagGroup

	// HasFlagGroups indicates whether this Node has at least one populated
	// flag group.
	HasFlagGroups() bool

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists on this Node or any of its parents, this
	// method will return nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a target Flag instance by its long-form name.
	//
	// If no such flag exists on this Node or any of its parents, this
	// method will return nil.
	FindLongFlag(name string) Flag

	HasCallback() bool

	Callback() CommandNodeCallback[T]

	IsHelpDisabled() bool
}

//

// NodeBuilder defines the base methods for common CommandTree Node setup.
type NodeBuilder[T any] interface {
	// WithDescription sets the description for the command tree node being built.
	WithDescription(description string) T

	// HasDescription indicates whether the command tree node being built has a
	// description value assigned.
	HasDescription() bool

	// Description returns the description value attached to the command tree node
	// builder.
	Description() string

	// WithFlagGroup attaches a flag group to the command tree node being built.
	//
	// See FlagGroup
	WithFlagGroup(flagGroup FlagGroupBuilder) T

	// WithFlagGroups allows attaching multiple flag groups at a time to the
	// command tree node being built.
	//
	// See FlagGroup
	WithFlagGroups(flagGroups ...FlagGroupBuilder) T

	// HasFlagGroups indicates whether the command tree node being built has one
	// or more flag groups attached.
	//
	// The `includeDefault` argument determines whether the default flag group
	// should be counted when testing for flag groups.
	//
	// See FlagGroup for more information about the default flag group.
	HasFlagGroups(includeDefault bool) bool

	// FlagGroups returns the flag groups attached to the command tree node
	// builder.
	//
	// The `includeDefault` argument determines whether the default flag group
	// should be included in the returned slice.
	//
	// See FlagGroup for more information about the default flag group.
	FlagGroups(includeDefault bool) []FlagGroupBuilder

	// WithFlag attaches the given flag to the command tree node being built.
	//
	// The given flag will be attached to the default flag group.
	//
	// See FlagGroup for more information about the default flag group.
	WithFlag(flag FlagBuilder) T

	// WithFlags allows attaching multiple flags at a time to the command tree
	// node being built.
	//
	// The given flags will be attached to the default flag group.
	//
	// See FlagGroup for more information about the default flag group.
	WithFlags(flags ...FlagBuilder) T

	// HasFlags indicates whether there are any flags attached to the command tree
	// node being built.
	//
	// This test does not include default `--help | -h` flags that may optionally
	// be provided by Argonaut.
	HasFlags() bool

	WithCallback(callback CommandNodeCallback[T]) T

	HasCallback() bool

	Callback() CommandNodeCallback[T]

	// WithHelpDisabled disables the automatic `-h` and `--help` flags for
	// rendering help text.
	WithHelpDisabled() T

	IsHelpDisabled() bool
}
