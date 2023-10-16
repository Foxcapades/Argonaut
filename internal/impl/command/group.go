package command

import "github.com/Foxcapades/Argonaut/pkg/argo"

// implements argo.CommandGroup
type group struct {
	name        string
	description string
	branches    []argo.CommandBranch
	leaves      []argo.CommandLeaf
}

func (g group) Description() string {
	return g.description
}

func (g group) HasDescription() bool {
	return len(g.description) > 0
}

func (g group) Name() string {
	return g.name
}

func (g group) Branches() []argo.CommandBranch {
	return g.branches
}

func (g group) HasBranches() bool {
	return len(g.branches) > 0
}

func (g group) Leaves() []argo.CommandLeaf {
	return g.leaves
}

func (g group) HasLeaves() bool {
	return len(g.branches) > 0
}

func (g group) FindChild(name string) argo.CommandNode {
	for _, leaf := range g.leaves {
		if leaf.Matches(name) {
			return leaf
		}
	}

	for _, branch := range g.branches {
		if branch.Matches(name) {
			return branch
		}
	}

	return nil
}
