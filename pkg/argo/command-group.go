package argo

// A CommandGroup is an organizational category containing one or more commands.
type CommandGroup interface {

	// Name returns the custom name for the CommandGroup.
	Name() string

	// Description returns the description value set on this CommandGroup.
	//
	// If no description was set on this CommandGroup, this method will return an
	// empty string.
	Description() string

	// HasDescription indicates whether a description value was set on this
	// CommandGroup.
	HasDescription() bool

	// Branches returns the CommandBranch nodes attached to this CommandGroup.
	Branches() []CommandBranch

	// HasBranches indicates whether this CommandGroup contains any branch nodes.
	HasBranches() bool

	// Leaves returns the CommandLeaf nodes attached to this CommandGroup.
	Leaves() []CommandLeaf

	// HasLeaves indicates whether this CommandGroup contains any leaf nodes.
	HasLeaves() bool

	// FindChild searches this CommandGroup instance for a CommandBranch or
	// CommandLeaf node that matches the given string.
	//
	// Commands may match on either their name or one of their aliases.
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
