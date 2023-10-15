package command

import (
	"errors"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/internal/impl/command/comutil"
	"github.com/Foxcapades/Argonaut/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func NewBranchBuilder(name string) argo.CommandBranchBuilder {
	return &commandBranchBuilder{
		name:       name,
		flagGroups: []argo.FlagGroupBuilder{flag.GroupBuilder(consts.DefaultGroupName)},
		comGroups:  []argo.CommandGroupBuilder{GroupBuilder(consts.DefaultGroupName)},
	}
}

type commandBranchBuilder struct {
	name         string
	desc         string
	helpDisabled bool
	comGroups    []argo.CommandGroupBuilder
	flagGroups   []argo.FlagGroupBuilder
	aliases      []string
	parent       argo.CommandNode
}

func (c commandBranchBuilder) GetName() string { return c.name }

func (c commandBranchBuilder) GetAliases() []string { return c.aliases }

func (c *commandBranchBuilder) Parent(node argo.CommandNode) { c.parent = node }

func (c *commandBranchBuilder) WithDescription(desc string) argo.CommandBranchBuilder {
	c.desc = desc
	return c
}

func (c *commandBranchBuilder) WithHelpDisabled() argo.CommandBranchBuilder {
	c.helpDisabled = true
	return c
}

func (c *commandBranchBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.CommandBranchBuilder {
	c.comGroups = append(c.comGroups, group)
	return c
}

func (c *commandBranchBuilder) WithBranch(branch argo.CommandBranchBuilder) argo.CommandBranchBuilder {
	c.comGroups[0].AddBranch(branch)
	return c
}

func (c *commandBranchBuilder) WithAliases(aliases ...string) argo.CommandBranchBuilder {
	c.aliases = aliases
	return c
}

func (c *commandBranchBuilder) WithLeaf(leaf argo.CommandLeafBuilder) argo.CommandBranchBuilder {
	c.comGroups[0].AddLeaf(leaf)
	return c
}

func (c *commandBranchBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandBranchBuilder {
	c.flagGroups[0].WithFlag(flag)
	return c
}

func (c *commandBranchBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.CommandBranchBuilder {
	c.flagGroups = append(c.flagGroups, flagGroup)
	return c
}

func (c *commandBranchBuilder) Build() (argo.CommandBranch, error) {
	errs := xerr.NewMultiError()

	// Ensure name is not blank
	if chars.IsBlank(c.name) {
		errs.AppendError(errors.New("command branch names must not be blank"))
	}

	// Ensure aliases are not blank
	for _, alias := range c.aliases {
		if chars.IsBlank(alias) {
			errs.AppendError(errors.New("command branch aliases must not be blank"))
		}
	}

	// Ensure a parent is set
	if c.parent == nil {
		panic("illegal state: attempted to build a command branch with no parent set")
	}

	// Create the out instance ahead of time so that we can set it as the parent
	// on the command groups we build.
	out := &commandBranch{
		name: c.name,
		desc: c.desc,
	}

	// Process Flag Groups
	flagGroups := make([]argo.FlagGroup, 0, len(c.flagGroups))
	flagutil.UniqueFlagNames(c.flagGroups, errs)
	for _, builder := range c.flagGroups {
		if group, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			flagGroups = append(flagGroups, group)
		}
	}

	// Process Command Groups
	commandGroups := make([]argo.CommandGroup, 0, len(c.comGroups))
	comutil.MassUniqueNames(c.comGroups, errs)
	for _, build := range c.comGroups {
		build.Parent(out)

		if group, err := build.Build(); err != nil {
			errs.AppendError(err)
		} else {
			commandGroups = append(commandGroups, group)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	out.flagGroups = flagGroups
	out.commandGroups = commandGroups
	out.parent = c.parent
	out.aliases = c.aliases

	return out, nil
}
