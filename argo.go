package cli

import (
	"errors"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func Command() argo.CommandBuilder {
	return argo.NewCommandBuilder()
}

// Tree returns a new argo.CommandTreeBuilder instance which can be used
// to construct an argo.CommandTree instance.
//
// A command tree is a tree of nested subcommands of arbitrary depth.  The tree
// consists of branch and leaf nodes, with the leaf nodes being the selectable
// final commands.
func Tree() argo.CommandTreeBuilder {
	return argo.NewCommandTreeBuilder()
}

func Parse(args []string, command any) error {
	if com, ok := command.(argo.CommandTree); ok {
		return argo.CommandTreeInterpreter(args, com).Run()
	} else if com, ok := command.(argo.Command); ok {
		return argo.CommandInterpreter(args, com).Run()
	} else {
		return errors.New("invalid command type passed to cli.Parse")
	}
}

// Branch returns a new argo.CommandBranchBuilder instance which can be
// used to construct an argo.CommandBranch instance.
func Branch(name string) argo.CommandBranchBuilder {
	return argo.NewCommandBranchBuilder(name)
}

// Leaf returns a new argo.CommandLeafBuilder instance which can be used
// to construct an argo.CommandLeaf instance.
func Leaf(name string) argo.CommandLeafBuilder {
	return argo.NewCommandLeafBuilder(name)
}

// FlagGroup returns a new argo.FlagGroupBuilder instance which can be used to
// construct an argo.FlagGroup instance.
func FlagGroup(name string) argo.FlagGroupBuilder {
	return argo.NewFlagGroupBuilder(name)
}

// Flag returns a new argo.FlagBuilder instance which can be used to construct
// an argo.Flag instance.
func Flag() argo.FlagBuilder {
	return argo.NewFlagBuilder()
}

// Argument returns a new argo.ArgumentBuilder instance which can be used to
// construct an argo.Argument instance.
func Argument() argo.ArgumentBuilder {
	return argo.NewArgumentBuilder()
}
