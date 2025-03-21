package argo

type CommandCallback = func(command Command)

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {

	// Name returns the name of the command.
	Name() string

	//

	// HasDescription indicates whether this command has a description value set.
	HasDescription() bool

	// Description returns the custom description for the command.
	//
	// Description values are used internally for rendering held text.
	Description() string

	//

	// FlagGroups returns the flag groups attached to this command.
	//
	// Flag groups are named categories of flags defined when building the
	// command.
	FlagGroups(includeDefault bool) []FlagGroup

	// HasFlagGroups indicates whether this command has any flag groups attached
	// to it.
	HasFlagGroups(includeDefault bool) bool

	//

	// FindShortFlag looks up a Flag instance by its short form.
	//
	// If no such flag could be found on this command, this method will return
	// nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a Flag instance by its long form.
	//
	// If no such flag could be found on this command, this method will return
	// nil.
	FindLongFlag(name string) Flag

	//

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

	// Arguments returns the positional Argument instances attached to this
	// Command.
	Arguments() []Argument

	//

	// HasUnmappedInputs indicates whether the command has collected any inputs
	// that were not mapped to any registered flag or argument.
	HasUnmappedInputs() bool

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

	AppendUnmappedInput(val string)

	//

	// HasUnmappedInputLabel indicates whether an unmapped label has been set on this
	// command instance.
	HasUnmappedInputLabel() bool

	// UnmappedInputLabel returns the label used when generating help text to
	// indicate the shape or purpose of unmapped inputs.
	UnmappedInputLabel() string

	//

	HasCallback() bool

	Callback() CommandCallback
}

//

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
//	    Build()
type CommandBuilder interface {

	// WithDescription sets the description value that will be used for the built
	// Command instance.
	//
	// Command descriptions are used when rendering help text.
	WithDescription(desc string) CommandBuilder

	HasDescription() bool

	Description() string

	//

	WithHelpDisabled() CommandBuilder

	IsHelpDisabled() bool

	//

	// WithFlagGroup appends the given FlagGroupBuilder to this CommandBuilder
	// instance.
	WithFlagGroup(group FlagGroupBuilder) CommandBuilder

	WithFlagGroups(groups ...FlagGroupBuilder) CommandBuilder

	HasFlagGroups(includeDefault bool) bool

	FlagGroups(includeDefault bool) []FlagGroupBuilder

	// WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
	// instance attached to this CommandBuilder.
	WithFlag(flag FlagBuilder) CommandBuilder

	WithFlags(flags ...FlagBuilder) CommandBuilder

	//

	// WithArgument appends the given ArgumentBuilder to this CommandBuilder's
	// list of positional arguments.
	WithArgument(arg ArgumentBuilder) CommandBuilder

	WithArguments(args ...ArgumentBuilder) CommandBuilder

	HasArguments() bool

	Arguments() []ArgumentBuilder

	//

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

	HasUnmappedInputLabel() bool

	UnmappedInputLabel() string

	//

	// WithCallback sets a callback function that will be executed immediately
	// after CLI parsing has completed successfully.
	WithCallback(cb CommandCallback) CommandBuilder

	HasCallback() bool

	Callback() CommandCallback

	//
	//
	// Build(ctx *WarningContext) (Command, error)
	//
	// // Parse reads the given arguments and attempts to populate the built Command
	// // instance based on the values parsed from the given inputs.
	// Parse(args []string) (Command, error)
	//
	// // MustParse is the same as Parse, however if an error is encountered while
	// // building the Command or parsing the input arguments, this method will
	// // panic.
	// MustParse(args []string) Command
}
