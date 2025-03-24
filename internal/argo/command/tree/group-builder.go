package tree

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewGroupBuilder(name string) argo.CommandGroupBuilder {
	return &CommandGroupBuilder{name: name}
}

type CommandGroupBuilder struct {
	name        string
	description string
	branches    []argo.BranchCommandBuilder
	leaves      []argo.LeafCommandBuilder
}

func (g *CommandGroupBuilder) Name() string {
	return g.name
}

func (g *CommandGroupBuilder) WithDescription(description string) argo.CommandGroupBuilder {
	g.description = description
	return g
}

func (g *CommandGroupBuilder) HasDescription() bool {
	return len(g.description) > 0
}

func (g *CommandGroupBuilder) Description() string {
	return g.description
}

func (g *CommandGroupBuilder) WithBranch(branch argo.BranchCommandBuilder) argo.CommandGroupBuilder {
	g.branches = append(g.branches, branch)
	return g
}

func (g *CommandGroupBuilder) WithBranches(branches ...argo.BranchCommandBuilder) argo.CommandGroupBuilder {
	// TODO implement me
	panic("implement me")
}

func (g *CommandGroupBuilder) HasBranches() bool {
	// TODO implement me
	panic("implement me")
}

func (g *CommandGroupBuilder) Branches() []argo.BranchCommandBuilder {
	// TODO implement me
	panic("implement me")
}

func (g *CommandGroupBuilder) WithLeaf(leaf argo.LeafCommandBuilder) argo.CommandGroupBuilder {
	g.leaves = append(g.leaves, leaf)
	return g
}

func (g *CommandGroupBuilder) WithLeaves(leaves ...argo.LeafCommandBuilder) argo.CommandGroupBuilder {
	// TODO implement me
	panic("implement me")
}

func (g *CommandGroupBuilder) HasLeaves() bool {
	// TODO implement me
	panic("implement me")
}

func (g *CommandGroupBuilder) Leaves() []argo.LeafCommandBuilder {
	// TODO implement me
	panic("implement me")
}

func (g *CommandGroupBuilder) HasSubcommands() bool {
	// TODO implement me
	panic("implement me")
}
