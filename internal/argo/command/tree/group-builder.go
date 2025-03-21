package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewGroupBuilder(name string) argo.CommandGroupBuilder {
	return &commandGroupBuilder{name: name}
}

type commandGroupBuilder struct {
	name        string
	description string
	parentNode  argo.ParentNode
	branches    []argo.BranchCommandBuilder
	leaves      []argo.LeafCommandBuilder
}

func (g *commandGroupBuilder) parent(node argo.ParentNode) {
	g.parentNode = node
}

func (g *commandGroupBuilder) WithDescription(description string) argo.CommandGroupBuilder {
	g.description = description
	return g
}

func (g *commandGroupBuilder) WithBranch(branch argo.BranchCommandBuilder) argo.CommandGroupBuilder {
	g.branches = append(g.branches, branch)
	return g
}

func (g *commandGroupBuilder) getBranches() []argo.BranchCommandBuilder {
	return g.branches
}

func (g *commandGroupBuilder) WithLeaf(leaf argo.LeafCommandBuilder) argo.CommandGroupBuilder {
	g.leaves = append(g.leaves, leaf)
	return g
}

func (g *commandGroupBuilder) getLeaves() []argo.LeafCommandBuilder {
	return g.leaves
}

func (g *commandGroupBuilder) hasSubcommands() bool {
	return len(g.leaves) > 0 || len(g.branches) > 0
}

func (g *commandGroupBuilder) Build(ctx *argo.WarningContext) (argo.CommandGroup, error) {
	errs := argo.NewMultiError()

	// Ensure that the name is not blank
	if chars.IsBlank(g.name) {
		errs.AppendError(errors.New("command group names must not be blank"))
	}

	// Require a parent value to be set.  If it is not set, it is a developer
	// error.
	if g.parentNode == nil {
		panic("illegal state: attempted to build a command group with no parent set")
	}

	// Ensure the command names and aliases are unique across the group.
	uniqueCommandNames(g.branches, g.leaves, errs)

	branches := make([]argo.BranchCommand, 0, len(g.branches))
	for _, builder := range g.branches {
		builder.ParentNode(g.parentNode)
		if branch, err := builder.Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			branches = append(branches, branch)
		}
	}

	leaves := make([]argo.LeafCommand, 0, len(g.leaves))
	for _, builder := range g.leaves {
		builder.ParentNode(g.parentNode)
		if leaf, err := builder.Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			leaves = append(leaves, leaf)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return commandGroup{g.name, g.description, branches, leaves}, nil
}
