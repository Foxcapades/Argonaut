package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {

	// Name returns the name of the command or subcommand.
	//
	// For root commands, this value will be the name of the command as it
	// appeared in the command line call.
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

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists, this method will return nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a target Flag instance by its long-form name.
	//
	// If no such flag exists, this method will return nil.
	FindLongFlag(name string) Flag

	Options() CommandOptions
}
