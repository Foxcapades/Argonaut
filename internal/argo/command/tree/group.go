package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type commandGroup struct {
	name        string
	description string
	branches    []argo.BranchCommand
	leaves      []argo.LeafCommand
}

func (g commandGroup) Description() string {
	return g.description
}

func (g commandGroup) HasDescription() bool {
	return len(g.description) > 0
}

func (g commandGroup) Name() string {
	return g.name
}

func (g commandGroup) Branches() []argo.BranchCommand {
	return g.branches
}

func (g commandGroup) HasBranches() bool {
	return len(g.branches) > 0
}

func (g commandGroup) Leaves() []argo.LeafCommand {
	return g.leaves
}

func (g commandGroup) HasLeaves() bool {
	return len(g.branches) > 0
}

func (g commandGroup) FindChild(name string) argo.ChildNode {
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
