package argo

// A CommandBuilder provides an API to configure the construction of a new
// Command instance.
//
// Example Usage:
//     cli.Command().
//         WithDescription("This is my command that does something.").
//         WithFlag(cli.Flag().
//             WithShortForm('v').
//             WithLongForm("verbose").
//             WithDescription("Enable verbose logging.")
//             WithBinding(&config.verbose)).
//         WithArgument(cli.Argument().
//             WithName("file").
//             WithDescription("File path.").
//             WithBinding(&config.file)).
//         Build()
type CommandBuilder interface {

	// WithDescription sets the description value that will be used for the built
	// Command instance.
	//
	// Command descriptions are used when rendering help text.
	WithDescription(desc string) CommandBuilder

	// HasDescription indicates whether a description value has been set on this
	// CommandBuilder instance.
	HasDescription() bool

	// GetDescription returns the description value set on this CommandBuilder.
	//
	// If no description has been set on this CommandBuilder, this method will
	// return an empty string.
	GetDescription() string

	// WithFlagGroup appends the given FlagGroupBuilder to this CommandBuilder
	// instance.
	WithFlagGroup(group FlagGroupBuilder) CommandBuilder

	// HasFlagGroups indicates whether this CommandBuilder has at least one custom
	// FlagGroup added to it.
	//
	// While all commands have a default flag group, this method will only return
	// true if there exists at least one additional FlagGroupBuilder instances.
	HasFlagGroups() bool

	// GetFlagGroups returns the FlagGroupBuilder instances that have been
	// attached to this CommandBuilder.
	GetFlagGroups() []FlagGroupBuilder

	// WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
	// instance attached to this CommandBuilder.
	WithFlag(flag FlagBuilder) CommandBuilder

	// HasFlags indicates whether this CommandBuilder instance has any FlagBuilder
	// instances attached.
	HasFlags() bool

	// WithArgument appends the given ArgumentBuilder to this CommandBuilder's
	// list of positional arguments.
	WithArgument(arg ArgumentBuilder) CommandBuilder

	WithUnmappedLabel(label string) CommandBuilder

	GetUnmappedLabel() string

	HasUnmappedLabel() bool

	// HasArguments indicates whether this CommandBuilder has any ArgumentBuilder
	// instances attached to it.
	HasArguments() bool

	// GetArguments returns a list of the arguments attached to this
	// CommandBuilder.
	GetArguments() []ArgumentBuilder

	Parse(args []string) (Command, error)

	MustParse(args []string) Command
}
