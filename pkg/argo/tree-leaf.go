package argo

// CommandLeafCallback defines the function type for a callback function that
// may be attached to a CommandLeaf.
//
// A CommandLeaf's callback function will be called once if and when a
// CommandLeaf is used in a CLI call.
type CommandLeafCallback = func(leaf CommandLeaf)

// A CommandLeaf is the final node in a CommandTree branch.
//
// Command leaves may be children of either a CommandTree directly, or of a
// Branch.
type CommandLeaf interface {
	ChildNode[CommandLeaf]

	// Arguments returns the positional Argument instances attached to this
	// Command.
	Arguments() []Argument

	// HasArguments indicates whether this Command has any positional arguments
	// attached.
	//
	// This method does not indicate whether those arguments were present on the
	// command line, it simply indicates whether Argument instances were attached
	// to the Command by the CommandBuilder.
	//
	// To determine whether an argument was present on the command line, test the
	// argument itself by using the Argument.WasHit method.
	HasArguments() bool

	// UnmappedInputs returns a collection of inputs that were passed to this
	// command that do not match any registered flag or argument.
	//
	// Unmapped inputs may be used to collect slices of positional arguments when
	// singular arguments can't be used.  For these situations, consider using
	// CommandBuilder.WithUnmappedLabel to set a help-text label indicating that
	// the command expects an arbitrary number of positional arguments.
	//
	// Defined positional arguments will always be hit before a value is added to
	// a command's unmapped inputs.
	UnmappedInputs() []string

	// HasUnmappedInputs indicates whether the command has collected any inputs
	// that were not mapped to any registered flag or argument.
	HasUnmappedInputs() bool

	AppendUnmappedInput(val string)

	// GetUnmappedLabel returns the label used when generating help text to
	// indicate the shape or purpose of unmapped inputs.
	GetUnmappedLabel() string

	// HasUnmappedLabel indicates whether an unmapped label has been set on this
	// command instance.
	HasUnmappedLabel() bool
}

// LeafBuilder defines a builder type that is used to construct
// CommandLeaf instances.
type LeafBuilder interface {
	ChildBuilder[LeafBuilder]

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
	WithUnmappedLabel(label string) LeafBuilder

	// WithArgument adds a positional argument to the CommandLeaf being built.
	WithArgument(argument ArgumentBuilder) LeafBuilder

	WithArguments(arguments ...ArgumentBuilder) LeafBuilder

	HasArguments() bool

	Arguments() []ArgumentBuilder

	Build(warnings *WarningContext) (CommandLeaf, error)
}
