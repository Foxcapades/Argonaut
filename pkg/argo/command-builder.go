package argo

import (
	"fmt"
	"os"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/util"
)

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

	WithHelpDisabled() CommandBuilder

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

	// WithCallback sets a callback function that will be executed immediately
	// after CLI parsing has completed successfully.
	WithCallback(cb CommandCallback) CommandBuilder

	Build(ctx *WarningContext) (Command, error)

	// Parse reads the given arguments and attempts to populate the built Command
	// instance based on the values parsed from the given inputs.
	Parse(args []string) (Command, error)

	// MustParse is the same as Parse, however if an error is encountered while
	// building the Command or parsing the input arguments, this method will
	// panic.
	MustParse(args []string) Command
}

func NewCommandBuilder() CommandBuilder {
	return &commandBuilder{
		flagGroups: []FlagGroupBuilder{NewFlagGroupBuilder(chars.DefaultGroupName)},
	}
}

type commandBuilder struct {
	description string
	unmapLabel  string
	flagGroups  []FlagGroupBuilder
	arguments   []ArgumentBuilder
	disableHelp bool
	callback    CommandCallback
}

func (b *commandBuilder) WithDescription(desc string) CommandBuilder {
	b.description = desc
	return b
}

func (b *commandBuilder) WithHelpDisabled() CommandBuilder {
	b.disableHelp = true
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

func (b *commandBuilder) WithCallback(cb CommandCallback) CommandBuilder {
	b.callback = cb
	return b
}

func (b commandBuilder) Parse(args []string) (Command, error) {
	ctx := new(WarningContext)
	if cmd, err := b.Build(ctx); err != nil {
		return nil, err
	} else {
		if err = newCommandInterpreter(args, cmd).Run(); err != nil {
			return nil, err
		}

		return cmd, nil
	}
}

func (b commandBuilder) MustParse(args []string) Command {
	return util.MustReturn(b.Parse(args))
}

func (b commandBuilder) Build(ctx *WarningContext) (Command, error) {
	errs := newMultiError()
	com := new(command)

	com.warnings = ctx

	if !b.disableHelp {
		group := b.flagGroups[0]

		if len(b.flagGroups) > 1 || b.flagGroups[0].size() > 5 {
			group = NewFlagGroupBuilder("Help Flags")
			b.flagGroups = append(b.flagGroups, group)
		}

		useLongH := true
		useShortH := true

		for _, group := range b.flagGroups {
			for _, flag := range group.getFlags() {
				if flag.getShortForm() == 'h' {
					useShortH = false
				}
				if flag.getLongForm() == "help" {
					useLongH = false
				}
				if !(useShortH || useLongH) {
					break
				}
			}
		}

		if useShortH || useLongH {
			group.WithFlag(makeCommandHelpFlag(useShortH, useLongH, com))
		}

	}

	com.flagGroups = make([]FlagGroup, 0, len(b.flagGroups))
	uniqueFlagNames(b.flagGroups, errs)
	for _, builder := range b.flagGroups {
		if builder.hasFlags() {
			if group, err := builder.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				com.flagGroups = append(com.flagGroups, group)
			}
		}
	}

	forceRequiredUntil := 0
	for i, builder := range b.arguments {
		if builder.isRequired() {
			forceRequiredUntil = i
		}
	}

	com.arguments = make([]Argument, 0, len(b.arguments))
	for i, builder := range b.arguments {

		if i < forceRequiredUntil && !builder.isRequired() {
			builder.Require()
			ctx.appendWarning(fmt.Sprintf("argument %d was not marked as required, but preceded required argument %d", i+1, forceRequiredUntil+1))
		}

		if arg, err := builder.Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			com.arguments = append(com.arguments, arg)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	com.description = b.description
	com.unmappedLabel = b.unmapLabel
	com.callback = b.callback

	return com, nil
}

func makeCommandHelpFlag(short, long bool, com Command) FlagBuilder {
	out := NewFlagBuilder().
		setIsHelpFlag().
		WithCallback(func(flag Flag) {
			util.Must(comRenderer{}.RenderHelp(com, os.Stdout))
			os.Exit(0)
		}).
		WithDescription("Prints this help text.")

	if short {
		out.WithShortForm('h')
	}

	if long {
		out.WithLongForm("help")
	}

	return out
}
