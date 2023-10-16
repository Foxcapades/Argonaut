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

	// WithUnmappedLabel provides a label for unmapped inputs.
	//
	// The unmapped label value is used when rendering the command usage line of
	// the auto-generated help text.  If a command expects an unknown number of
	// positional argument values, it is best to capture them as unmapped inputs
	// with a label.
	//
	// Example configuration:
	//     cli.CommandLeaf("my-leaf").
	//         WithUnmappedLabel("ITEMS...")
	//
	// Example usage line:
	//     Usage:
	//       my-leaf [ITEMS...]
	WithUnmappedLabel(label string) CommandLeafBuilder

	// GetUnmappedLabel returns the label configured on this builder for unmapped
	// cli inputs.
	//
	// If no unmapped value label has been set, this method returns an empty
	// string.
	GetUnmappedLabel() string

	// HasUnmappedLabel indicates whether an unmapped value label has been set on
	// this CommandLeafBuilder instance.
	HasUnmappedLabel() bool

	// WithCallback sets a callback on the CommandLeaf to be built that will be
	// executed when the command leaf is used, after any tree or branch callbacks
	// that were set on parent nodes.
	WithCallback(cb CommandLeafCallback) CommandLeafBuilder

	// GetCallback returns the callback set on this CommandLeafBuilder instance.
	//
	// If no callback has been set, this method returns nil.
	GetCallback() CommandLeafCallback

	// HasCallback indicates whether a CommandLeafCallback has been set on this
	// CommandLeafBuilder instance.
	HasCallback() bool

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
