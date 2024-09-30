package argo

type SubCommandSpecContainer interface {
	SubCommands() []SubCommandSpec

	FindSubCommand(name string) SubCommandSpec

	Branches() []CommandBranchSpec

	FindBranch(name string) CommandBranchSpec

	HasBranches() bool

	Leaves() []CommandLeafSpec

	FindLeaf(name string) CommandLeafSpec

	HasLeaves() bool
}
