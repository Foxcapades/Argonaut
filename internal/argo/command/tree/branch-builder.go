package tree

import (
	"errors"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBranchBuilder(name string) argo.BranchCommandBuilder {
	return &commandBranchBuilder{
		name:       name,
		flagGroups: []argo.FlagGroupBuilder{flag.NewGroupBuilder(chars.DefaultGroupName)},
		comGroups:  []argo.CommandGroupBuilder{NewGroupBuilder(chars.DefaultGroupName)},
	}
}

type commandBranchBuilder struct {
	name         string
	desc         string
	helpDisabled bool
	comGroups    []argo.CommandGroupBuilder
	flagGroups   []argo.FlagGroupBuilder
	aliases      []string
	parentNode   argo.ParentNode
	callback     argo.CommandBranchCallback

	onIncompleteHandler argo.IncompleteCommandHandler
}

func (c *commandBranchBuilder) Name() string {
	return c.name
}

func (c *commandBranchBuilder) parent(node argo.ParentNode) {
	c.parentNode = node
}

func (c *commandBranchBuilder) WithAliases(aliases ...string) argo.BranchCommandBuilder {
	c.aliases = aliases
	return c
}

func (c *commandBranchBuilder) Aliases() []string {
	return c.aliases
}

func (c *commandBranchBuilder) WithDescription(desc string) argo.BranchCommandBuilder {
	c.desc = desc
	return c
}

func (c *commandBranchBuilder) WithCallback(cb argo.CommandBranchCallback) argo.BranchCommandBuilder {
	c.callback = cb
	return c
}

func (c *commandBranchBuilder) WithHelpDisabled() argo.BranchCommandBuilder {
	c.helpDisabled = true
	return c
}

func (c *commandBranchBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.BranchCommandBuilder {
	c.comGroups = append(c.comGroups, group)
	return c
}

func (c *commandBranchBuilder) WithBranch(branch argo.BranchCommandBuilder) argo.BranchCommandBuilder {
	c.comGroups[0].WithBranch(branch)
	return c
}

func (c *commandBranchBuilder) WithLeaf(leaf argo.LeafCommandBuilder) argo.BranchCommandBuilder {
	c.comGroups[0].WithLeaf(leaf)
	return c
}

func (c *commandBranchBuilder) WithFlag(flag argo.FlagBuilder) argo.BranchCommandBuilder {
	c.flagGroups[0].WithFlag(flag)
	return c
}

func (c *commandBranchBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.BranchCommandBuilder {
	c.flagGroups = append(c.flagGroups, flagGroup)
	return c
}

func (c *commandBranchBuilder) OnIncomplete(handler argo.IncompleteCommandHandler) argo.BranchCommandBuilder {
	c.onIncompleteHandler = handler
	return c
}

func (c *commandBranchBuilder) Build(ctx *cli_err.WarningContext) (argo.BranchCommand, error) {
	errs := argo.NewMultiError()

	// Ensure name is not blank
	if err := chars.ValidateCommandNodeName(c.name); err != nil {
		errs.AppendError(err)
	}

	// Ensure aliases are not blank
	for _, alias := range c.aliases {
		if chars.IsBlank(alias) {
			errs.AppendError(errors.New("command branch aliases must not be blank"))
		}
	}

	// Ensure a parent is set
	if c.parentNode == nil {
		panic("illegal state: attempted to build a command branch with no parent set")
	}

	// Create the out instance ahead of time so that we can set it as the parent
	// on the command groups we build.
	out := &commandBranch{
		name: c.name,
		desc: c.desc,
	}

	// If auto-help is not disabled, then...
	if !c.helpDisabled {
		metaGroup := false
		hasShortH := false
		hasLongH := false

		// If the default group name has been changed to a custom name then enable
		// the meta group.
		if c.flagGroups[0].Name() != chars.DefaultGroupName {
			metaGroup = true
		}

		if !metaGroup && c.flagGroups[0].Size() > 5 {
			metaGroup = true
		}

		if !metaGroup && len(c.flagGroups) > 1 {
			for i := 1; i < len(c.flagGroups); i++ {
				if c.flagGroups[i].HasFlags() {
					metaGroup = true
				}
			}
		}

	OUTER:
		for _, group := range c.flagGroups {
			for _, flag := range group.Flags() {
				if flag.HasShortForm() && flag.ShortForm() == 'h' {
					hasShortH = true
				}
				if flag.HasLongForm() && flag.LongForm() == "help" {
					hasLongH = true
				}
				if hasShortH && hasLongH {
					break OUTER
				}
			}
		}

		if !(hasLongH || hasShortH) {
			var group argo.FlagGroupBuilder

			if metaGroup {
				group = flag.NewGroupBuilder("Help Flags")
				c.flagGroups = append(c.flagGroups, group)
			} else {
				group = c.flagGroups[0]
			}

			flag := flag.NewBuilder().
				WithDescription("Prints this help text.").
				WithCallback(func(f argo.Flag) {
					utils.Must(comBranchRenderer{}.RenderHelp(out, os.Stdout))
					os.Exit(0)
				})

			if !hasLongH {
				flag.WithLongForm("help")
			}

			if !hasShortH {
				flag.WithShortForm('h')
			}

			group.WithFlag(flag)
		}
	}

	// Process Flag Groups
	flagGroups := make([]argo.FlagGroup, 0, len(c.flagGroups))
	uniqueFlagNames(c.flagGroups, errs)
	for _, builder := range c.flagGroups {
		if builder.HasFlags() {
			if group, err := builder.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	// Process Command Groups
	commandGroups := make([]argo.CommandGroup, 0, len(c.comGroups))
	massUniqueCommandNames(c.comGroups, errs)
	for _, build := range c.comGroups {
		if build.hasSubcommands() {
			build.parent(out)

			if group, err := build.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	out.warnings = ctx
	out.flagGroups = flagGroups
	out.commandGroups = commandGroups
	out.parent = c.parentNode
	out.aliases = c.aliases
	out.callback = c.callback
	out.onIncompleteHandler = utils.IfElse(c.onIncompleteHandler != nil, c.onIncompleteHandler, nil)

	return out, nil
}
