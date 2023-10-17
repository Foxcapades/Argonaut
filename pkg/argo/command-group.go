package argo

type CommandGroup interface {
	Name() string

	Description() string

	HasDescription() bool

	Branches() []CommandBranch

	HasBranches() bool

	Leaves() []CommandLeaf

	HasLeaves() bool

	FindChild(name string) CommandNode
}

type commandGroup struct {
	name        string
	description string
	branches    []CommandBranch
	leaves      []CommandLeaf
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

func (g commandGroup) Branches() []CommandBranch {
	return g.branches
}

func (g commandGroup) HasBranches() bool {
	return len(g.branches) > 0
}

func (g commandGroup) Leaves() []CommandLeaf {
	return g.leaves
}

func (g commandGroup) HasLeaves() bool {
	return len(g.branches) > 0
}

func (g commandGroup) FindChild(name string) CommandNode {
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
