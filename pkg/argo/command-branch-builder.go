package argo

import (
	"errors"
	"fmt"
	"os"
)

// A CommandBranchBuilder instance may be used to configure a new CommandBranch
// instance to be built.
//
// CommandBranches are intermediate steps between the root of the CommandTree
// and the CommandLeaf instances.
//
// For example, given the following command example, the tree is "foo", the
// branch is "bar", and the leaf is "fizz":
//     ./foo bar fizz
type CommandBranchBuilder interface {

	// GetName returns the name assigned to this CommandBranchBuilder.
	//
	// CommandBranch names are required and thus are set at the time that the
	// CommandBranchBuilder instance is constructed.
	GetName() string

	// WithAliases appends the given alias strings to this CommandBranchBuilder's
	// alias list.
	//
	// Aliases must be unique and must not conflict with any other command branch
	// or leaf names or aliases at a given command tree level.
	//
	// Example:
	//     cli.CommandBranch("service").
	//         WithAliases("svc")
	WithAliases(aliases ...string) CommandBranchBuilder

	// GetAliases returns the aliases assigned to this CommandBranchBuilder.
	GetAliases() []string

	// Parent sets the parent CommandNode for the CommandBranch being built.
	//
	// Values set using this method before build time will be disregarded.
	Parent(node CommandNode)

	// WithDescription sets a description value for the CommandBranch being built.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandBranchBuilder

	HasDescription() bool

	GetDescription() string

	// WithHelpDisabled disables the automatic '-h | --help' flag that is enabled
	// by default.
	WithHelpDisabled() CommandBranchBuilder

	// WithBranch appends a child branch to the default CommandGroup for this
	// CommandBranchBuilder.
	WithBranch(branch CommandBranchBuilder) CommandBranchBuilder

	// WithLeaf appends a child leaf to the default CommandGroup for this
	// CommandBranchBuilder.
	WithLeaf(leaf CommandLeafBuilder) CommandBranchBuilder

	// WithCommandGroup appends a custom CommandGroup to this
	// CommandBranchBuilder.
	//
	// CommandGroups are used to organize subcommands into named categories that
	// are primarily used for rendering help text.
	WithCommandGroup(group CommandGroupBuilder) CommandBranchBuilder

	// WithFlag appends the given FlagBuilder to the default FlagGroup attached to
	// this CommandBranchBuilder.
	WithFlag(flag FlagBuilder) CommandBranchBuilder

	// WithFlagGroup appends the given custom FlagGroupBuilder to this
	// CommandBranchBuilder instance.
	//
	// Custom flag groups are primarily used for categorizing flags in the
	// rendered help text.
	WithFlagGroup(flagGroup FlagGroupBuilder) CommandBranchBuilder

	WithCallback(cb CommandBranchCallback) CommandBranchBuilder

	HasCallback() bool

	GetCallback() CommandBranchCallback

	build() (CommandBranch, error)
}

func NewCommandBranchBuilder(name string) CommandBranchBuilder {
	return &commandBranchBuilder{
		name:       name,
		flagGroups: []FlagGroupBuilder{NewFlagGroupBuilder(defaultGroupName)},
		comGroups:  []CommandGroupBuilder{NewCommandGroupBuilder(defaultGroupName)},
	}
}

type commandBranchBuilder struct {
	name         string
	desc         string
	helpDisabled bool
	comGroups    []CommandGroupBuilder
	flagGroups   []FlagGroupBuilder
	aliases      []string
	parent       CommandNode
	callback     CommandBranchCallback
}

func (c commandBranchBuilder) GetName() string {
	return c.name
}

func (c *commandBranchBuilder) Parent(node CommandNode) {
	c.parent = node
}

// Aliases /////////////////////////////////////////////////////////////////////

func (c *commandBranchBuilder) WithAliases(aliases ...string) CommandBranchBuilder {
	c.aliases = aliases
	return c
}

func (c commandBranchBuilder) GetAliases() []string {
	return c.aliases
}

func (c commandBranchBuilder) HasAliases() bool {
	return len(c.aliases) > 0
}

// Description /////////////////////////////////////////////////////////////////

func (c *commandBranchBuilder) WithDescription(desc string) CommandBranchBuilder {
	c.desc = desc
	return c
}

func (c commandBranchBuilder) HasDescription() bool {
	return len(c.desc) > 0
}

func (c commandBranchBuilder) GetDescription() string {
	return c.desc
}

// Callback ////////////////////////////////////////////////////////////////////

func (c *commandBranchBuilder) WithCallback(cb CommandBranchCallback) CommandBranchBuilder {
	c.callback = cb
	return c
}

func (c commandBranchBuilder) HasCallback() bool {
	return c.callback != nil
}

func (c commandBranchBuilder) GetCallback() CommandBranchCallback {
	return c.callback
}

func (c *commandBranchBuilder) WithHelpDisabled() CommandBranchBuilder {
	c.helpDisabled = true
	return c
}

func (c *commandBranchBuilder) WithCommandGroup(group CommandGroupBuilder) CommandBranchBuilder {
	c.comGroups = append(c.comGroups, group)
	return c
}

func (c *commandBranchBuilder) WithBranch(branch CommandBranchBuilder) CommandBranchBuilder {
	c.comGroups[0].WithBranch(branch)
	return c
}

func (c *commandBranchBuilder) WithLeaf(leaf CommandLeafBuilder) CommandBranchBuilder {
	c.comGroups[0].WithLeaf(leaf)
	return c
}

func (c *commandBranchBuilder) WithFlag(flag FlagBuilder) CommandBranchBuilder {
	c.flagGroups[0].WithFlag(flag)
	return c
}

func (c *commandBranchBuilder) WithFlagGroup(flagGroup FlagGroupBuilder) CommandBranchBuilder {
	c.flagGroups = append(c.flagGroups, flagGroup)
	return c
}

func (c *commandBranchBuilder) build() (CommandBranch, error) {
	errs := newMultiError()

	// Ensure name is not blank
	if isBlank(c.name) {
		errs.AppendError(errors.New("command branch names must not be blank"))
	}

	// Ensure aliases are not blank
	for _, alias := range c.aliases {
		if isBlank(alias) {
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

	// If auto-help is not disabled, then...
	if !c.helpDisabled {
		metaGroup := false
		hasShortH := false
		hasLongH := false

		// If the default group name has been changed to a custom name then enable
		// the meta group.
		if c.flagGroups[0].GetName() != defaultGroupName {
			metaGroup = true
		}

		if !metaGroup && len(c.flagGroups[0].getFlags()) > 5 {
			metaGroup = true
		}

		if !metaGroup && len(c.flagGroups) > 1 {
			for i := 1; i < len(c.flagGroups); i++ {
				if c.flagGroups[i].hasFlags() {
					metaGroup = true
				}
			}
		}

	OUTER:
		for _, group := range c.flagGroups {
			for _, flag := range group.getFlags() {
				if flag.hasShortForm() && flag.getShortForm() == 'h' {
					hasShortH = true
				}
				if flag.hasLongForm() && flag.getLongForm() == "help" {
					hasLongH = true
				}
				if hasShortH && hasLongH {
					break OUTER
				}
			}
		}

		if !(hasLongH || hasShortH) {
			var group FlagGroupBuilder

			if metaGroup {
				group = NewFlagGroupBuilder("Meta Flags")
				c.flagGroups = append(c.flagGroups, group)
			} else {
				group = c.flagGroups[0]
			}

			flag := NewFlagBuilder().
				WithDescription("Prints this help text.").
				WithCallback(func(f Flag) {
					fmt.Println(renderCommandBranch(out))
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
	flagGroups := make([]FlagGroup, 0, len(c.flagGroups))
	uniqueFlagNames(c.flagGroups, errs)
	for _, builder := range c.flagGroups {
		if builder.hasFlags() {
			if group, err := builder.build(); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	// Process Command Groups
	commandGroups := make([]CommandGroup, 0, len(c.comGroups))
	massUniqueCommandNames(c.comGroups, errs)
	for _, build := range c.comGroups {
		build.Parent(out)

		if group, err := build.build(); err != nil {
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
	out.callback = c.callback

	return out, nil
}
