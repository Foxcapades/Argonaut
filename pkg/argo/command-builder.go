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

	// WithFlagGroup appends the given FlagGroupBuilder to this CommandBuilder
	// instance.
	WithFlagGroup(group FlagGroupBuilder) CommandBuilder

	// WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
	// instance attached to this CommandBuilder.
	WithFlag(flag FlagBuilder) CommandBuilder

	// WithArgument appends the given ArgumentBuilder to this CommandBuilder's
	// list of positional arguments.
	WithArgument(arg ArgumentBuilder) CommandBuilder

	// WithUnmappedLabel sets the help-text label for unmapped arguments.
	//
	// This is useful when your command takes an arbitrary number of argument
	// inputs, and you would like the help text to indicate as such.
	//
	// Example Config:
	//     cli.Command().
	//         WithUnmappedLabel("[FILE...]")
	//
	// Example Result:
	//     Usage:
	//       my-command [FILE...]
	WithUnmappedLabel(label string) CommandBuilder

	Parse(args []string) (Command, error)

	MustParse(args []string) Command
}

func NewCommandBuilder() CommandBuilder {
	return &commandBuilder{
		flagGroups: []FlagGroupBuilder{NewFlagGroupBuilder(defaultGroupName)},
	}
}

type commandBuilder struct {
	description string
	unmapLabel  string
	flagGroups  []FlagGroupBuilder
	arguments   []ArgumentBuilder
}

func (b *commandBuilder) WithDescription(desc string) CommandBuilder {
	b.description = desc
	return b
}

func (b *commandBuilder) WithFlagGroup(group FlagGroupBuilder) CommandBuilder {
	b.flagGroups = append(b.flagGroups, group)
	return b
}

func (b *commandBuilder) WithUnmappedLabel(label string) CommandBuilder {
	b.unmapLabel = label
	return b
}

func (b *commandBuilder) WithFlag(flag FlagBuilder) CommandBuilder {
	b.flagGroups[0].WithFlag(flag)
	return b
}

func (b *commandBuilder) WithArgument(arg ArgumentBuilder) CommandBuilder {
	b.arguments = append(b.arguments, arg)
	return b
}

func (b commandBuilder) Parse(args []string) (Command, error) {
	if cmd, err := b.build(); err != nil {
		return nil, err
	} else {
		if err = CommandInterpreter(args, cmd).Run(); err != nil {
			return nil, err
		}

		return cmd, nil
	}
}

func (b commandBuilder) MustParse(args []string) Command {
	return mustReturn(b.Parse(args))
}

func (b commandBuilder) build() (Command, error) {
	errs := newMultiError()

	flagGroups := make([]FlagGroup, 0, len(b.flagGroups))
	for _, builder := range b.flagGroups {
		if builder.hasFlags() {
			if group, err := builder.build(); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	arguments := make([]Argument, 0, len(b.arguments))
	for _, builder := range b.arguments {
		if arg, err := builder.build(); err != nil {
			errs.AppendError(err)
		} else {
			arguments = append(arguments, arg)
		}
	}

	return &command{
		description:   b.description,
		flagGroups:    flagGroups,
		arguments:     arguments,
		unmappedLabel: b.unmapLabel,
	}, nil
}
