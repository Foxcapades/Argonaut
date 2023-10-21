package argo

import (
	"errors"
	"os"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/util"
)

// A CommandTreeBuilder is a builder type used to construct a CommandTree
// instance.
//
// A CommandTree is a command that consists of branching subcommands.  Examples
// of such commands include the `go` command, `docker`, or `kubectl`.
//
// To use the Docker command example we have a command tree that includes the
// following:
//     docker
//      |- compose
//      |   |- build
//      |   |- down
//      |   |- ...
//      |- container
//      |   |- exec
//      |   |- ls
//      |   |- ...
//      |- ...
type CommandTreeBuilder interface {

	// WithDescription sets a description value for the root of this command tree.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandTreeBuilder

	// WithCallback sets a callback for the CommandTree being built.
	//
	// If set, this callback will be executed on parsing success.  Each level of
	// a command tree may have a callback.  The callbacks are called in the order
	// the command segments appear in the CLI call.
	WithCallback(cb CommandTreeCallback) CommandTreeBuilder

	// WithHelpDisabled disables the automatic `-h` and `--help` flags for
	// rendering help text.
	WithHelpDisabled() CommandTreeBuilder

	// WithBranch appends the given branch builder to be built with this command
	// tree.
	//
	// The built branch will be available as a subcommand directly under the root
	// command call.
	WithBranch(branch CommandBranchBuilder) CommandTreeBuilder

	// WithLeaf appends the given leaf builder to be built with this command tree.
	//
	// The built leaf will be available as a subcommand directly under the root
	// command call.
	WithLeaf(leaf CommandLeafBuilder) CommandTreeBuilder

	// WithCommandGroup appends the given command group builder to be built with
	// this command tree.
	//
	// Command groups are used for organizing subcommands into named groups that
	// are primarily used for rendering help text.
	WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder

	// WithFlag appends the given flag builder to the default flag group attached
	// to this command tree builder.  The default flag group is an automatically
	// created group for containing otherwise ungrouped flags.
	WithFlag(flag FlagBuilder) CommandTreeBuilder

	// WithFlagGroup appends the given flag group builder to this command tree
	// builder.  Flag groups are for organizing flags into named categories that
	// are primarily used for rendering help text.
	WithFlagGroup(flagGroup FlagGroupBuilder) CommandTreeBuilder

	Build(warnings *WarningContext) (CommandTree, error)

	// Parse builds the command tree and attempts to parse the given CLI arguments
	// into that command tree's components.
	Parse(args []string) (CommandTree, error)

	// MustParse calls Parse and panics if an error is returned.
	MustParse(args []string) CommandTree

	hasSubCommands() bool
}

func NewCommandTreeBuilder() CommandTreeBuilder {
	return &commandTreeBuilder{
		commandGroups: []CommandGroupBuilder{NewCommandGroupBuilder(chars.DefaultGroupName)},
		flagGroups:    []FlagGroupBuilder{NewFlagGroupBuilder(chars.DefaultGroupName)},
	}
}

type commandTreeBuilder struct {
	desc          string
	helpDisabled  bool
	commandGroups []CommandGroupBuilder
	flagGroups    []FlagGroupBuilder
	callback      CommandTreeCallback
}

func (t *commandTreeBuilder) WithDescription(desc string) CommandTreeBuilder {
	t.desc = desc
	return t
}

func (t *commandTreeBuilder) WithCallback(cb CommandTreeCallback) CommandTreeBuilder {
	t.callback = cb
	return t
}

func (t *commandTreeBuilder) WithHelpDisabled() CommandTreeBuilder {
	t.helpDisabled = true
	return t
}

func (t *commandTreeBuilder) WithBranch(branch CommandBranchBuilder) CommandTreeBuilder {
	t.commandGroups[0].WithBranch(branch)
	return t
}

func (t *commandTreeBuilder) WithLeaf(leaf CommandLeafBuilder) CommandTreeBuilder {
	t.commandGroups[0].WithLeaf(leaf)
	return t
}

func (t *commandTreeBuilder) WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder {
	t.commandGroups = append(t.commandGroups, group)
	return t
}

func (t *commandTreeBuilder) WithFlag(flag FlagBuilder) CommandTreeBuilder {
	t.flagGroups[0].WithFlag(flag)
	return t
}

func (t *commandTreeBuilder) WithFlagGroup(flagGroup FlagGroupBuilder) CommandTreeBuilder {
	t.flagGroups = append(t.flagGroups, flagGroup)
	return t
}

func (t commandTreeBuilder) Parse(args []string) (CommandTree, error) {
	ctx := new(WarningContext)
	ct, err := t.Build(ctx)
	if err != nil {
		return nil, err
	}

	err = newCommandTreeInterpreter(args, ct).Run()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (t commandTreeBuilder) MustParse(args []string) CommandTree {
	ctx := new(WarningContext)
	ct := util.MustReturn(t.Build(ctx))
	util.Must(newCommandTreeInterpreter(args, ct).Run())
	return ct
}

func (t *commandTreeBuilder) hasSubCommands() bool {
	for _, group := range t.commandGroups {
		if group.hasSubcommands() {
			return true
		}
	}

	return false
}

func (t *commandTreeBuilder) Build(warnings *WarningContext) (CommandTree, error) {
	errs := newMultiError()

	tree := new(commandTree)

	if !t.hasSubCommands() {
		errs.AppendError(errors.New("command tree has no subcommands"))
	}

	if !t.helpDisabled {
		group := t.flagGroups[0]

		if len(t.flagGroups) > 1 || t.flagGroups[0].size() > 5 {
			group = NewFlagGroupBuilder("Help Flags")
			t.flagGroups = append(t.flagGroups, group)
		}

		useLongH := true
		useShortH := true

		for _, group := range t.flagGroups {
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
			group.WithFlag(makeCommandTreeHelpFlag(useLongH, useLongH, tree))
		}
	}

	flagGroups := make([]FlagGroup, 0, len(t.flagGroups))
	uniqueFlagNames(t.flagGroups, errs)
	for _, builder := range t.flagGroups {
		if builder.hasFlags() {
			if group, err := builder.Build(warnings); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	commandGroups := make([]CommandGroup, 0, len(t.commandGroups))
	massUniqueCommandNames(t.commandGroups, errs)
	for _, builder := range t.commandGroups {
		builder.parent(tree)
		if builder.hasSubcommands() {
			if group, err := builder.Build(warnings); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	tree.warnings = warnings
	tree.description = t.desc
	tree.flagGroups = flagGroups
	tree.commandGroups = commandGroups
	tree.callback = t.callback

	return tree, nil
}

func makeCommandTreeHelpFlag(short, long bool, tree CommandTree) FlagBuilder {
	out := NewFlagBuilder().
		setIsHelpFlag().
		WithCallback(func(flag Flag) {
			util.Must(comTreeRenderer{}.RenderHelp(tree, os.Stdout))
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
