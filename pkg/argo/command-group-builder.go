package argo

import "errors"

type CommandGroupBuilder interface {
	parent(node CommandNode)

	// WithDescription sets a description value for this CommandGroupBuilder.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandGroupBuilder

	// WithBranch appends the given branch builder to this command group builder.
	WithBranch(branch CommandBranchBuilder) CommandGroupBuilder

	getBranches() []CommandBranchBuilder

	// WithLeaf appends the given leaf builder to this command group builder.
	WithLeaf(leaf CommandLeafBuilder) CommandGroupBuilder

	getLeaves() []CommandLeafBuilder

	hasSubcommands() bool

	Build(warnings *WarningContext) (CommandGroup, error)
}

func NewCommandGroupBuilder(name string) CommandGroupBuilder {
	return &commandGroupBuilder{name: name}
}

type commandGroupBuilder struct {
	name        string
	description string
	parentNode  CommandNode
	branches    []CommandBranchBuilder
	leaves      []CommandLeafBuilder
}

func (g *commandGroupBuilder) parent(node CommandNode) {
	g.parentNode = node
}

func (g *commandGroupBuilder) WithDescription(description string) CommandGroupBuilder {
	g.description = description
	return g
}

func (g *commandGroupBuilder) WithBranch(branch CommandBranchBuilder) CommandGroupBuilder {
	g.branches = append(g.branches, branch)
	return g
}

func (g commandGroupBuilder) getBranches() []CommandBranchBuilder {
	return g.branches
}

func (g *commandGroupBuilder) WithLeaf(leaf CommandLeafBuilder) CommandGroupBuilder {
	g.leaves = append(g.leaves, leaf)
	return g
}

func (g commandGroupBuilder) getLeaves() []CommandLeafBuilder {
	return g.leaves
}

func (g commandGroupBuilder) hasSubcommands() bool {
	return len(g.leaves) > 0 || len(g.branches) > 0
}

func (g *commandGroupBuilder) Build(ctx *WarningContext) (CommandGroup, error) {
	errs := newMultiError()

	// Ensure that the name is not blank
	if isBlank(g.name) {
		errs.AppendError(errors.New("command group names must not be blank"))
	}

	// Require a parent value to be set.  If it is not set, it is a developer
	// error.
	if g.parentNode == nil {
		panic("illegal state: attempted to build a command group with no parent set")
	}

	// Ensure the command names and aliases are unique across the group.
	uniqueCommandNames(g.branches, g.leaves, errs)

	branches := make([]CommandBranch, 0, len(g.branches))
	for _, builder := range g.branches {
		builder.parent(g.parentNode)
		if branch, err := builder.Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			branches = append(branches, branch)
		}
	}

	leaves := make([]CommandLeaf, 0, len(g.leaves))
	for _, builder := range g.leaves {
		builder.parent(g.parentNode)
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
