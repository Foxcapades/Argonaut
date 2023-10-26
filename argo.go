// Package cli provides a convenience methods for constructing command line
// interfaces using the argo package.
//
// Commands may be either singular commands or command trees.  A singular
// command has no subcommand and may take any number of flags and/or arguments.
// A command tree is a root command that may have an arbitrary depth of
// branching subcommands which each may take their own flags, with the leaf
// nodes of the tree also accepting arguments.
//
// An example single command:
//     tar -xf foo.tgz
//
// Here the single command `tar` accepts the flags `-x` and -f`, with the `-f`
// flag taking the argument `foo.tgz`.
//
// An example command tree:
//     docker compose -f my-docker-compose.yml up my-service
//
// Here the command tree is constructed of 3 levels, the root of the tree
// (docker), the intermediary branch (compose) and the leaf command (up).  The
// branch is taking an argument (-f) which is itself taking an argument
// (my-docker-compose.yml).  The leaf command (up) is accepting an optional
// argument (my-service).
//
// Command construction starts with either the `cli.Command` function or the
// `cli.Tree` function.
package cli

import (
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

// Command returns a new CommandBuilder instance which can be used to construct
// a Command instance.
//
// This function and Tree are the two entrypoints into the Argonaut library.
// This function returns a value that may be used in a call chain to construct
// a full-featured command line interface.
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

// CommandGroup returns a new CommandGroupBuilder in stance which can be used to
// construct a CommandGroup instance.
func CommandGroup(name string) argo.CommandGroupBuilder {
	return argo.NewCommandGroupBuilder(name)
}

// FlagGroup returns a new FlagGroupBuilder instance which can be used to
// construct an FlagGroup instance.
func FlagGroup(name string) argo.FlagGroupBuilder {
	return argo.NewFlagGroupBuilder(name)
}

// ShortFlag returns a new FlagBuilder instance with the short form already set
// to the given value.
//
// This function is a shortcut for:
//     cli.Flag().WithShortForm(...)
func ShortFlag(f byte) argo.FlagBuilder {
	return argo.NewFlagBuilder().WithShortForm(f)
}

// LongFlag returns a new FlagBuilder instance with the long form already set
// to the given value.
//
// This function is a shortcut for:
//     cli.Flag().WithLongForm(...)
func LongFlag(name string) argo.FlagBuilder {
	return argo.NewFlagBuilder().WithLongForm(name)
}

// Flag returns a new FlagBuilder instance which can be used to construct a Flag
// instance.
func Flag() argo.FlagBuilder {
	return argo.NewFlagBuilder()
}

// Argument returns a new ArgumentBuilder instance which can be used to
// construct an Argument instance.
func Argument() argo.ArgumentBuilder {
	return argo.NewArgumentBuilder()
}
