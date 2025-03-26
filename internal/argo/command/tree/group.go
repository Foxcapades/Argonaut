package tree

import (
	"iter"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const DefaultCommandGroupName = "__DEFAULT_COMMAND_GROUP__"

type CommandGroup struct {
	name        string
	description string
	branches    []argo.BranchCommand
	leaves      []argo.LeafCommand
}

func (g CommandGroup) Description() string {
	return g.description
}

func (g CommandGroup) HasDescription() bool {
	return len(g.description) > 0
}

func (g CommandGroup) Name() string {
	return g.name
}

func (g CommandGroup) Branches() []argo.BranchCommand {
	return g.branches
}

func (g CommandGroup) HasBranches() bool {
	return len(g.branches) > 0
}

func (g CommandGroup) Leaves() []argo.LeafCommand {
	return g.leaves
}

func (g CommandGroup) HasLeaves() bool {
	return len(g.branches) > 0
}

func (g CommandGroup) HasSubcommands() bool {
	return len(g.branches)+len(g.leaves) > 0
}

func (g CommandGroup) Subcommands() iter.Seq[argo.ChildNode] {
	return func(y func(argo.ChildNode) bool) {
		for _, b := range g.branches {
			y(b)
		}
		for _, l := range g.leaves {
			y(l)
		}
	}
}

func (g CommandGroup) FindChild(name string) argo.ChildNode {
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
