package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

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

	WithOptions(opts CommandOptions) CommandBuilder

	Options() CommandOptions
}
