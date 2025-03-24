package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// A LeafCommand is the final node in a CommandTree branch.
//
// Command leaves may be children of either a CommandTree directly, or of a
// BranchCommand.
type LeafCommand interface {
	ChildNode
  
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

  Callback() CommandCallback[LeafCommand]

  IsHelpDisabled() bool

  // Arguments returns the positional Argument instances attached to this
  // LeafCommand.
  Arguments() []Argument

  // HasArguments indicates whether this LeafCommand has any positional
	// Arguments attached.
  //
  // This method does not indicate whether those Arguments were present on the
  // command line, it only indicates whether any Argument instances were
  // attached to the command with the LeafCommandBuilder.
  //
  // To determine whether an Argument was present on the command line, test the
  // Argument itself by using the Argument.WasHit method.
  HasArguments() bool

  // UnmappedInputs returns a collection of inputs that were passed to this
  // LeafCommand that do not match any registered flag or argument.
  //
  // Unmapped inputs may be used to collect slices of positional arguments when
  // singular arguments can't be used.  For these situations, consider using
  // LeafCommandBuilder.WithUnmappedLabel to set a help-text label indicating
	// that the command expects an arbitrary number of positional arguments.
  //
  // Defined positional arguments will always be hit before a value is added to
  // a command's unmapped inputs.
  UnmappedInputs() []string

  // HasUnmappedInputs indicates whether the command has collected any inputs
  // that were not mapped to any registered flag or argument.
  HasUnmappedInputs() bool

  AppendUnmappedInput(val string)

  // UnmappedInputLabel returns the label used when generating help text to
  // indicate the shape or purpose of unmapped inputs.
  UnmappedInputLabel() string

  // HasUnmappedInputLabel indicates whether an unmapped label has been set on this
  // command.
  HasUnmappedInputLabel() bool

}

// LeafCommandBuilder defines a builder type that is used to construct
// LeafCommand instances.
type LeafCommandBuilder interface {
	ChildNodeBuilder
	
  // WithDescription sets the description value that will be used for the built
  // cli command or subcommand.
  //
  // Descriptions are used when rendering help text.
  WithDescription(desc string) LeafCommandBuilder

  HasDescription() bool

  Description() string

  // WithHelpDisabled disables the automatic `-h` and `--help` flags for
  // rendering help text.
  WithHelpDisabled() LeafCommandBuilder

  IsHelpDisabled() bool

  // WithFlagGroup appends the given FlagGroupBuilder to this CLI component
  // builder.
  WithFlagGroup(group FlagGroupBuilder) LeafCommandBuilder

  WithFlagGroups(groups ...FlagGroupBuilder) LeafCommandBuilder

  HasFlagGroups(includeDefault bool) bool

  FlagGroups(includeDefault bool) []FlagGroupBuilder

  // WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
  // instance attached to this CLI component builder.
  WithFlag(flag FlagBuilder) LeafCommandBuilder

  WithFlags(flags ...FlagBuilder) LeafCommandBuilder

  HasFlags() bool

  WithCallback(callback CommandCallback[LeafCommand]) LeafCommandBuilder

  Callback() CommandCallback[LeafCommand]

  HasCallback() bool

  // WithAlias assigns the given alias to the target command node.
  //
  // Command aliases must be unique per level in a command tree.  This means
  // that for any given step in the tree, no alias may conflict with another
  // branch or leaf subcommand's name or aliases.
  //
  // This also applies if a subcommand node is reused at multiple levels of the
  // command tree.
  //
  // If a conflict is found between subcommand names and/or aliases, an error
  // will be returned when attempting to build the command tree.
  //
  // Example:
  //   cli.Leaf("list").WithAlias("ls")
  WithAlias(alias string) LeafCommandBuilder

  // WithAliases assigns the given aliases to the target command node.
  //
  // Command aliases must be unique per level in a command tree.  This means
  // that for any given step in the tree, no alias may conflict with another
  // branch or leaf subcommand's name or aliases.
  //
  // This also applies if a subcommand node is reused at multiple levels of the
  // command tree.
  //
  // If a conflict is found between subcommand names and/or aliases, an error
  // will be returned when attempting to build the command tree.
  WithAliases(aliases ...string) LeafCommandBuilder

	// SetParentNode is used by the builder process to link command tree nodes
	// together.
	//
	// Setting this value before it is passed to the parent command builder will
	// have no effect as the value will be overwritten.
	//
	// Setting this value _after_ it is passed to the parent command builder will
	// cause undefined behavior.
  SetParentNode(parent ParentNodeBuilder) LeafCommandBuilder

  // Arguments returns the positional Argument instances attached to this
  // LeafCommand.
  Arguments() []Argument

  // HasArguments indicates whether this LeafCommand has any positional
	// Arguments attached.
  //
  // This method does not indicate whether those Arguments were present on the
  // command line, it only indicates whether any Argument instances were
  // attached to the command with the LeafCommandBuilder.
  //
  // To determine whether an Argument was present on the command line, test the
  // Argument itself by using the Argument.WasHit method.
  HasArguments() bool

  // UnmappedInputs returns a collection of inputs that were passed to this
  // LeafCommand that do not match any registered flag or argument.
  //
  // Unmapped inputs may be used to collect slices of positional arguments when
  // singular arguments can't be used.  For these situations, consider using
  // LeafCommandBuilder.WithUnmappedLabel to set a help-text label indicating
	// that the command expects an arbitrary number of positional arguments.
  //
  // Defined positional arguments will always be hit before a value is added to
  // a command's unmapped inputs.
  UnmappedInputs() []string

  // HasUnmappedInputs indicates whether the command has collected any inputs
  // that were not mapped to any registered flag or argument.
  HasUnmappedInputs() bool

  AppendUnmappedInput(val string)

  // UnmappedInputLabel returns the label used when generating help text to
  // indicate the shape or purpose of unmapped inputs.
  UnmappedInputLabel() string

  // HasUnmappedInputLabel indicates whether an unmapped label has been set on this
  // command.
  HasUnmappedInputLabel() bool

}
