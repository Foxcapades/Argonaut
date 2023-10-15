package argo

type CommandGroup interface {
	Name() string

	Branches() []CommandBranch

	Leaves() []CommandLeaf

	FindChild(name string) CommandNode
}

type CommandGroupBuilder interface {
	Parent(node CommandNode)

	AddBranch(branch CommandBranchBuilder) CommandGroupBuilder

	GetBranches() []CommandBranchBuilder

	AddLeaf(leaf CommandLeafBuilder) CommandGroupBuilder

	GetLeaves() []CommandLeafBuilder

	HasSubcommands() bool

	Build() (CommandGroup, error)
}
