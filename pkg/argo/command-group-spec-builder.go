package argo

type CommandGroupSpecBuilder interface {
	WithName(name string) CommandGroupSpecBuilder

	WithDescription(description string) CommandGroupSpecBuilder

	WithBranch(branch CommandBranchSpecBuilder) CommandGroupSpecBuilder

	WithLeaf(leaf CommandLeafSpecBuilder) CommandGroupSpecBuilder

	Build(config Config) (CommandBranchSpec, error)
}
