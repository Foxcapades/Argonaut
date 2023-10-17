package cli

import (
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

// Command returns a new CommandBuilder instance which can be used to construct
// a Command instance.
func Command() argo.CommandBuilder {
	return argo.NewCommandBuilder()
}

// Tree returns a new CommandTreeBuilder instance which can be used to construct
// a CommandTree instance.
//
// A command tree is a tree of nested subcommands of arbitrary depth.  The tree
// consists of branch and leaf nodes, with the leaf nodes being the selectable
// final commands.
func Tree() argo.CommandTreeBuilder {
	return argo.NewCommandTreeBuilder()
}

// Branch returns a new CommandBranchBuilder instance which can be used to
// construct an CommandBranch instance.
func Branch(name string) argo.CommandBranchBuilder {
	return argo.NewCommandBranchBuilder(name)
}

// Leaf returns a new CommandLeafBuilder instance which can be used to construct
// an CommandLeaf instance.
func Leaf(name string) argo.CommandLeafBuilder {
	return argo.NewCommandLeafBuilder(name)
}

// FlagGroup returns a new FlagGroupBuilder instance which can be used to
// construct an FlagGroup instance.
func FlagGroup(name string) argo.FlagGroupBuilder {
	return argo.NewFlagGroupBuilder(name)
}

// Flag returns a new FlagBuilder instance which can be used to construct
// an Flag instance.
func Flag() argo.FlagBuilder {
	return argo.NewFlagBuilder()
}

// Argument returns a new ArgumentBuilder instance which can be used to
// construct an Argument instance.
func Argument() argo.ArgumentBuilder {
	return argo.NewArgumentBuilder()
}
