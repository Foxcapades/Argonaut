package argo

type CommandGroupBuilder interface {
	Parent(node CommandNode)

	WithDescription(desc string) CommandGroupBuilder

	HasDescription() bool

	GetDescription() string

	WithBranch(branch CommandBranchBuilder) CommandGroupBuilder

	GetBranches() []CommandBranchBuilder

	WithLeaf(leaf CommandLeafBuilder) CommandGroupBuilder

	GetLeaves() []CommandLeafBuilder

	HasSubcommands() bool

	Build() (CommandGroup, error)
}
