package argo

type CommandTreeBuilder interface {
	WithDescription(desc string) CommandTreeBuilder

	HasDescription() bool

	GetDescription() string

	WithCallback(cb CommandTreeCallback) CommandTreeBuilder

	HasCallback() bool

	GetCallback() CommandTreeCallback

	WithHelpDisabled() CommandTreeBuilder

	WithBranch(branch CommandBranchBuilder) CommandTreeBuilder

	WithLeaf(leaf CommandLeafBuilder) CommandTreeBuilder

	WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder

	WithFlag(flag FlagBuilder) CommandTreeBuilder

	WithFlagGroup(flagGroup FlagGroupBuilder) CommandTreeBuilder

	Parse(args []string) (CommandTree, error)

	MustParse(args []string) CommandTree
}
