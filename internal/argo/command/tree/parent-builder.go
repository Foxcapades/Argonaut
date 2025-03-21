package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func newParentBuilder[T any](root T) parentBuilder[T] {
	return parentBuilder[T]{
		nodeBuilder: newNodeBuilder(root),
	}
}

type parentBuilder[T any] struct {
	nodeBuilder[T]

	groups []argo.CommandGroupBuilder

	incompleteFn argo.IncompleteCommandHandler[T]
}

// ╔════════════════════════════════════════════════════════════════════════╗ //
// ║    API Implementation                                                  ║ //
// ╚════════════════════════════════════════════════════════════════════════╝ //

func (p *parentBuilder[T]) WithBranch(branch argo.BranchCommandBuilder) T {
	p.defaultCommandGroup().WithBranch(branch)
	branch.SetParentNode(utils.Cast[argo.ParentBuilder[any]](p))
	return p.root
}

func (p *parentBuilder[T]) WithBranches(branches ...argo.BranchCommandBuilder) T {
	p.defaultCommandGroup().WithBranches(branches...)
	for _, branch := range branches {
		branch.SetParentNode(utils.Cast[argo.ParentBuilder[any]](p))
	}
	return p.root
}

func (p *parentBuilder[T]) WithLeaf(leaf argo.LeafCommandBuilder) T {
	p.defaultCommandGroup().WithLeaf(leaf)
	leaf.SetParentNode(utils.Cast[argo.ParentBuilder[any]](p))
	return p.root
}

func (p *parentBuilder[T]) WithLeaves(leaves ...argo.LeafCommandBuilder) T {
	p.defaultCommandGroup().WithLeaves(leaves...)
	for _, leaf := range leaves {
		leaf.SetParentNode(utils.Cast[argo.ParentBuilder[any]](p))
	}
	return p.root
}

func (p *parentBuilder[T]) WithCommandGroup(group argo.CommandGroupBuilder) T {
	p.groups = append(p.groups, group)
	return p.root
}

func (p *parentBuilder[T]) WithCommandGroups(groups ...argo.CommandGroupBuilder) T {
	p.groups = append(p.groups, groups...)
	return p.root
}

func (p *parentBuilder[T]) HasCommandGroups() bool {
	return len(p.groups) > 0
}

func (p *parentBuilder[T]) HasSubcommands() bool {
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

func (p *parentBuilder[T]) CommandGroups() []argo.CommandGroupBuilder {
	return p.groups
}

func (p *parentBuilder[T]) WithIncompleteHandler(handler argo.IncompleteCommandHandler[T]) T {
	p.incompleteFn = handler
	return p.root
}

func (p *parentBuilder[T]) HasIncompleteHandler() bool {
	return p.incompleteFn != nil
}

func (p *parentBuilder[T]) IncompleteHandler() argo.IncompleteCommandHandler[T] {
	return p.incompleteFn
}

// ╔════════════════════════════════════════════════════════════════════════╗ //
// ║    Internal Methods                                                    ║ //
// ╚════════════════════════════════════════════════════════════════════════╝ //

func (p *parentBuilder[T]) ensureDefaultCommandGroup() {
	if len(p.groups) == 0 {
		p.groups = append(p.groups, NewGroupBuilder(DefaultCommandGroupName))
		return
	}

	if p.groups[0].Name() != DefaultCommandGroupName {
		tmp := make([]argo.CommandGroupBuilder, 0, len(p.groups)+1)
		tmp = append(tmp, NewGroupBuilder(DefaultCommandGroupName))
		p.groups = append(tmp, p.groups...)
	}
}

func (p *parentBuilder[T]) defaultCommandGroup() argo.CommandGroupBuilder {
	p.ensureDefaultCommandGroup()
	return p.groups[0]
}
