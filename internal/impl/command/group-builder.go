package command

import (
	"errors"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/impl/command/comutil"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func GroupBuilder(name string) argo.CommandGroupBuilder {
	return &groupBuilder{name: name}
}

type groupBuilder struct {
	name     string
	parent   argo.CommandNode
	branches []argo.CommandBranchBuilder
	leaves   []argo.CommandLeafBuilder
}

func (g *groupBuilder) Parent(node argo.CommandNode) {
	g.parent = node
}

func (g *groupBuilder) AddBranch(branch argo.CommandBranchBuilder) argo.CommandGroupBuilder {
	g.branches = append(g.branches, branch)
	return g
}

func (g groupBuilder) GetBranches() []argo.CommandBranchBuilder {
	return g.branches
}

func (g *groupBuilder) AddLeaf(leaf argo.CommandLeafBuilder) argo.CommandGroupBuilder {
	g.leaves = append(g.leaves, leaf)
	return g
}

func (g groupBuilder) GetLeaves() []argo.CommandLeafBuilder {
	return g.leaves
}

func (g groupBuilder) HasSubcommands() bool {
	return len(g.leaves) > 0 || len(g.branches) > 0
}

func (g *groupBuilder) Build() (argo.CommandGroup, error) {
	errs := xerr.NewMultiError()

	// Ensure that the name is not blank
	if chars.IsBlank(g.name) {
		errs.AppendError(errors.New("command group names must not be blank"))
	}

	// Require a parent value to be set.  If it is not set, it is a developer
	// error.
	if g.parent == nil {
		panic("illegal state: attempted to build a command group with no parent set")
	}

	// Ensure the command names and aliases are unique across the group.
	comutil.UniqueNames(g.branches, g.leaves, errs)

	branches := make([]argo.CommandBranch, 0, len(g.branches))
	for _, builder := range g.branches {
		builder.Parent(g.parent)
		if branch, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			branches = append(branches, branch)
		}
	}

	leaves := make([]argo.CommandLeaf, 0, len(g.leaves))
	for _, builder := range g.leaves {
		builder.Parent(g.parent)
		if leaf, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			leaves = append(leaves, leaf)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return group{g.name, branches, leaves}, nil
}
