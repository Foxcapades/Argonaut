{{ $vars := (types "argo.BranchCommandBuilder" "CommandBranchBuilder" "argo.BranchCommand") }}
package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBranchBuilder(name string) argo.BranchCommandBuilder {
	return &CommandBranchBuilder{
		name:       name,
		flagGroups: []argo.FlagGroupBuilder{flag.NewGroupBuilder(chars.DefaultGroupName)},
		comGroups:  []argo.CommandGroupBuilder{NewGroupBuilder(chars.DefaultGroupName)},
	}
}

type CommandBranchBuilder struct {
	name      string
	comGroups []argo.CommandGroupBuilder
	aliases   []string
	parent    argo.ParentNodeBuilder

	onIncompleteHandler argo.IncompleteCommandHandler[argo.BranchCommand]

	{{ template "CommandBuilderBaseProps" $vars }}
}

func (i *CommandBranchBuilder) Name() string {
	return i.name
}

func (i *CommandBranchBuilder) WithAlias(alias string) argo.BranchCommandBuilder {
	i.aliases = append(i.aliases, alias)
	return i
}

func (i *CommandBranchBuilder) WithAliases(aliases ...string) argo.BranchCommandBuilder {
	i.aliases = append(i.aliases, aliases...)
	return i
}

func (i *CommandBranchBuilder) Aliases() []string {
	return i.aliases
}

func (i *CommandBranchBuilder) HasAliases() bool {
	return len(i.aliases) > 0
}

func (i *CommandBranchBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.BranchCommandBuilder {
	i.comGroups = append(i.comGroups, group)
	return i
}

func (i *CommandBranchBuilder) WithCommandGroups(groups ...argo.CommandGroupBuilder) argo.BranchCommandBuilder {
	i.comGroups = append(i.comGroups, groups...)
	return i
}

func (i *CommandBranchBuilder) HasCommandGroups(includeDefault bool) bool {
	return (includeDefault && len(i.comGroups) > 0) ||
		(i.hasDefaultCommandGroup() && len(i.comGroups) > 1) ||
		len(i.comGroups) > 0
}

func (i *CommandBranchBuilder) CommandGroups(includeDefault bool) []argo.CommandGroupBuilder {
	if !includeDefault && i.hasDefaultCommandGroup() {
		return i.comGroups[1:]
	}

	return i.comGroups
}

func (i *CommandBranchBuilder) hasDefaultCommandGroup() bool {
	return len(i.comGroups) > 0 && i.comGroups[0].Name() == DefaultCommandGroupName
}

func (i *CommandBranchBuilder) WithBranch(branch argo.BranchCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranch(branch)
	return i
}

func (i *CommandBranchBuilder) WithBranches(branches ...argo.BranchCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranches(branches...)
	return i
}

func (i *CommandBranchBuilder) WithLeaf(leaf argo.LeafCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaf(leaf)
	return i
}

func (i *CommandBranchBuilder) WithLeaves(leaves ...argo.LeafCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaves(leaves...)
	return i
}

func (i *CommandBranchBuilder) HasSubcommands() bool {
	for _, group := range i.comGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *CommandBranchBuilder) WithIncompleteHandler(handler argo.IncompleteCommandHandler[argo.BranchCommand]) argo.BranchCommandBuilder {
	i.onIncompleteHandler = handler
	return i
}

func (i *CommandBranchBuilder) HasIncompleteHandler() bool {
	return i.onIncompleteHandler != nil
}

func (i *CommandBranchBuilder) IncompleteHandler() argo.IncompleteCommandHandler[argo.BranchCommand] {
	return i.onIncompleteHandler
}

func (i *CommandBranchBuilder) SetParentNode(parent argo.ParentNodeBuilder) argo.BranchCommandBuilder {
	i.parent = parent
	return i
}

func (i *CommandBranchBuilder) ParentNode() argo.ParentNodeBuilder {
	return i.parent
}

{{ template "CommandBuilderBaseFuncs" $vars }}
