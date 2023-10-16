package argo

// CommandLeafBuilder defines a builder type that is used to construct
// CommandLeaf instances.
type CommandLeafBuilder interface {

	// GetName returns the name assigned to this CommandLeafBuilder.
	GetName() string

	// WithAliases attaches the given aliases to this CommandLeafBuilder.
	//
	// Command aliases must be unique per level in a command tree.  This means
	// that for a given step in the tree, no alias may conflict with another
	// CommandNode's name or aliases.
	//
	// If a conflict is found, an error will be returned when attempting to build
	// the CommandLeaf.
	WithAliases(aliases ...string) CommandLeafBuilder

	// GetAliases returns the aliases that have been attached to this
	// CommandLeafBuilder.
	GetAliases() []string

	// Parent sets the parent node for this CommandLeafBuilder.  Values set with
	// this method before build time will be disregarded.
	Parent(node CommandNode)

	// WithDescription sets a description for the CommandLeaf to be built.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandLeafBuilder

	WithUnmappedLabel(label string) CommandLeafBuilder

	GetUnmappedLabel() string

	HasUnmappedLabel() bool

	// WithArgument adds a positional argument to the CommandLeaf being built.
	WithArgument(argument ArgumentBuilder) CommandLeafBuilder

	// WithFlagGroup adds a new FlagGroup to this CommandLeaf being built.
	WithFlagGroup(flagGroup FlagGroupBuilder) CommandLeafBuilder

	// WithFlag adds the given flag to the default FlagGroup for the CommandLeaf
	// being built.
	WithFlag(flag FlagBuilder) CommandLeafBuilder

	// Build attempts to build a CommandLeaf from the values set on this builder.
	Build() (CommandLeaf, error)
}
