package argo

type SubCommandContainer interface {
	SubCommands() []SubCommand

	FindSubCommand(name string) SubCommand

	Branches() []CommandBranch

	FindBranch(name string) CommandBranch

	HasBranches() bool

	Leaves() []CommandLeaf

	FindLeaf(name string) CommandLeaf

	HasLeaves() bool
}
