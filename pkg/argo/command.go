package argo

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {
	CommandBase[Command]

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
	CommandBuilderBase[CommandBuilder, Command]

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
}
