package cli

import (
	"errors"

	"github.com/Foxcapades/Argonaut/internal/impl/argument"
	"github.com/Foxcapades/Argonaut/internal/impl/command"
	"github.com/Foxcapades/Argonaut/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/internal/interpret"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func Command() argo.CommandBuilder {
	return command.Builder()
}

// CommandTree returns a new argo.CommandTreeBuilder instance which can be used
// to construct an argo.CommandTree instance.
//
// A command tree is a tree of nested subcommands of arbitrary depth.  The tree
// consists of branch and leaf nodes, with the leaf nodes being the selectable
// final commands.
func CommandTree() argo.CommandTreeBuilder {
	return command.TreeBuilder()
}

func Parse(args []string, command any) error {
	if com, ok := command.(argo.CommandTree); ok {
		return interpret.CommandTreeInterpreter(args, com).Run()
	} else if com, ok := command.(argo.Command); ok {
		return interpret.CommandInterpreter(args, com).Run()
	} else {
		return errors.New("invalid command type passed to cli.Parse")
	}
}

// CommandBranch returns a new argo.CommandBranchBuilder instance which can be
// used to construct an argo.CommandBranch instance.
func CommandBranch(name string) argo.CommandBranchBuilder {
	return command.NewBranchBuilder(name)
}

// CommandLeaf returns a new argo.CommandLeafBuilder instance which can be used
// to construct an argo.CommandLeaf instance.
func CommandLeaf(name string) argo.CommandLeafBuilder {
	return command.NewLeafBuilder(name)
}

// FlagGroup returns a new argo.FlagGroupBuilder instance which can be used to
// construct an argo.FlagGroup instance.
func FlagGroup(name string) argo.FlagGroupBuilder {
	return flag.GroupBuilder(name)
}

// Flag returns a new argo.FlagBuilder instance which can be used to construct
// an argo.Flag instance.
func Flag() argo.FlagBuilder {
	return flag.NewBuilder()
}

// Argument returns a new argo.ArgumentBuilder instance which can be used to
// construct an argo.Argument instance.
func Argument() argo.ArgumentBuilder {
	return argument.NewBuilder()
}
