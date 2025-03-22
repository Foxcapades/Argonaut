package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func newParentBuilder[T, O any](root T) parentBuilder[T, O] {
	return parentBuilder[T, O]{
		nodeBuilder: newNodeBuilder[T, O](root),
	}
}

type parentBuilder[T, O any] struct {
	nodeBuilder[T, O]

	groups []argo.CommandGroupBuilder

	incompleteFn argo.IncompleteCommandHandler[O]
}

// Branches

func (p *parentBuilder[T, O]) WithBranch(branch argo.BranchCommandBuilder) T {
	p.defaultCommandGroup().WithBranch(branch)
	branch.SetParentNode(utils.Cast[argo.ParentBuilder[any, any]](p))
	return p.root
}

func (p *parentBuilder[T, O]) WithBranches(branches ...argo.BranchCommandBuilder) T {
	p.defaultCommandGroup().WithBranches(branches...)
	for _, branch := range branches {
		branch.SetParentNode(utils.Cast[argo.ParentBuilder[any, any]](p))
	}
	return p.root
}

// Leaves

func (p *parentBuilder[T, O]) WithLeaf(leaf argo.LeafCommandBuilder) T {
	p.defaultCommandGroup().WithLeaf(leaf)
	leaf.SetParentNode(utils.Cast[argo.ParentBuilder[any, any]](p))
	return p.root
}

func (p *parentBuilder[T, O]) WithLeaves(leaves ...argo.LeafCommandBuilder) T {
	p.defaultCommandGroup().WithLeaves(leaves...)
	for _, leaf := range leaves {
		leaf.SetParentNode(utils.Cast[argo.ParentBuilder[any, any]](p))
	}
	return p.root
}

// Command Groups

func (p *parentBuilder[T, O]) WithCommandGroup(group argo.CommandGroupBuilder) T {
	p.groups = append(p.groups, group)
	return p.root
}

func (p *parentBuilder[T, O]) WithCommandGroups(groups ...argo.CommandGroupBuilder) T {
	p.groups = append(p.groups, groups...)
	return p.root
}

func (p *parentBuilder[T, O]) HasCommandGroups() bool {
	return len(p.groups) > 0
}

func (p *parentBuilder[T, O]) CommandGroups() []argo.CommandGroupBuilder {
	return p.groups
}

func (p *parentBuilder[T, O]) defaultCommandGroup() argo.CommandGroupBuilder {
	p.ensureDefaultCommandGroup()
	return p.groups[0]
}

func (p *parentBuilder[T, O]) ensureDefaultCommandGroup() {
	if len(p.groups) == 0 {
		p.groups = append(p.groups, NewGroupBuilder(DefaultCommandGroupName))
		return
	}

	if p.groups[0].Name() != DefaultCommandGroupName {
		p.groups = append(
			append(
				make([]argo.CommandGroupBuilder, 0, len(p.groups)+1),
				NewGroupBuilder(DefaultCommandGroupName),
			),
			p.groups...,
		)
	}
}

// Subcommands

func (p *parentBuilder[T, O]) HasSubcommands() bool {
	if p.hasDefaultFlagGroup() {
		return true
	}

	for _, group := range p.groups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

// Incomplete Handlers

func (p *parentBuilder[T, O]) WithIncompleteHandler(handler argo.IncompleteCommandHandler[O]) T {
	p.incompleteFn = handler
	return p.root
}

func (p *parentBuilder[T, O]) HasIncompleteHandler() bool {
	return p.incompleteFn != nil
}

func (p *parentBuilder[T, O]) IncompleteHandler() argo.IncompleteCommandHandler[O] {
	return p.incompleteFn
}
