package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {
	
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

  Callback() CommandCallback[Command]

  IsHelpDisabled() bool

  // Arguments returns the positional Argument instances attached to this
  // Command.
  Arguments() []Argument

  // HasArguments indicates whether this Command has any positional
	// Arguments attached.
  //
  // This method does not indicate whether those Arguments were present on the
  // command line, it only indicates whether any Argument instances were
  // attached to the command with the CommandBuilder.
  //
  // To determine whether an Argument was present on the command line, test the
  // Argument itself by using the Argument.WasHit method.
  HasArguments() bool

  // UnmappedInputs returns a collection of inputs that were passed to this
  // Command that do not match any registered flag or argument.
  //
  // Unmapped inputs may be used to collect slices of positional arguments when
  // singular arguments can't be used.  For these situations, consider using
  // CommandBuilder.WithUnmappedLabel to set a help-text label indicating
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

// A CommandBuilder provides an API to configure the construction of a new
// Command instance.
//
// Example Usage:
//
//	cli.Command().
//	    WithDescription("This is my command that does something.").
//	    WithFlag(cli.Flag().
//	        WithShortForm('v').
//	        WithLongForm("verbose").
//	        WithDescription("Enable verbose logging.")
//	        WithBinding(&config.verbose)).
//	    WithArgument(cli.Argument().
//	        WithName("file").
//	        WithDescription("File path.").
//	        WithBinding(&config.file)).
type CommandBuilder interface {
  
  // WithDescription sets the description value that will be used for the built
  // cli command or subcommand.
  //
  // Descriptions are used when rendering help text.
  WithDescription(desc string) CommandBuilder

  HasDescription() bool

  Description() string

  // WithHelpDisabled disables the automatic `-h` and `--help` flags for
  // rendering help text.
  WithHelpDisabled() CommandBuilder

  IsHelpDisabled() bool

  // WithFlagGroup appends the given FlagGroupBuilder to this CLI component
  // builder.
  WithFlagGroup(group FlagGroupBuilder) CommandBuilder

  WithFlagGroups(groups ...FlagGroupBuilder) CommandBuilder

  HasFlagGroups(includeDefault bool) bool

  FlagGroups(includeDefault bool) []FlagGroupBuilder

  // WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
  // instance attached to this CLI component builder.
  WithFlag(flag FlagBuilder) CommandBuilder

  WithFlags(flags ...FlagBuilder) CommandBuilder

  HasFlags() bool

  WithCallback(callback CommandCallback[Command]) CommandBuilder

  Callback() CommandCallback[Command]

  HasCallback() bool

  
  // WithArgument appends the given ArgumentBuilder instance to this
  // CommandBuilder's list of positional arguments.
  WithArgument(arg ArgumentBuilder) CommandBuilder

  // WithArguments appends the given ArgumentBuilder instances to this
  // CommandBuilder's list of positional arguments.
  WithArguments(args ...ArgumentBuilder) CommandBuilder

  // HasArguments indicates whether any positional arguments have been set on
  // this CommandBuilder instance.
  HasArguments() bool

  // Arguments returns a slice containing the ArgumentBuilders that have been
  // set on this CommandBuilder instance.
  Arguments() []ArgumentBuilder

  // WithUnmappedInputLabel sets the help-text label for unmapped arguments.
  //
  // This is useful when your command takes an arbitrary number of argument
  // inputs, and you would like the help text to indicate as such.
  //
  // Example Config:
  //     cli.Command().
  //         WithUnmappedInputLabel("[FILE...]")
  //
  // Example Result:
  //     Usage:
  //       my-command [FILE...]
  WithUnmappedInputLabel(label string) CommandBuilder

  // HasUnmappedInputLabel indicates whether an unmapped input label has been
  // set on this CommandBuilder instance.
  HasUnmappedInputLabel() bool

  // UnmappedInputLabel returns the unmapped input label set on this
  // CommandBuilder instance.
  //
  // If no unmapped input label has been set, this method returns an empty
  // string.
  UnmappedInputLabel() string

}