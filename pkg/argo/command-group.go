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
