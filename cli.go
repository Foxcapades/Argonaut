// Package cli provides a convenience methods for constructing command line
// interfaces.
//
// Commands may be either singular commands or command trees.  A singular
// command has no subcommand and may take any number of flags and/or arguments.
// A command tree is a root command that may have an arbitrary depth of
// branching subcommands which each may take their own flags, with the leaf
// nodes of the tree also accepting arguments.
//
// An example single command:
//
//	tar -xf foo.tgz
//
// Here the single command `tar` accepts the flags `-x` and -f`, with the `-f`
// flag taking the argument `foo.tgz`.
//
// An example command tree:
//
//	docker compose -f my-docker-compose.yml up my-service
//
// Here the command tree is constructed of 3 levels, the root of the tree
// (docker), the intermediary branch (compose) and the leaf command (up).  The
// branch is taking a flag (-f) which is itself taking an argument
// (my-docker-compose.yml).  The leaf command (up) is accepting an optional
// argument (my-service).
//
// Command construction starts with either the `cli.Command` function or the
// `cli.Tree` function.
package cli

import (
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/simple"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// region Single Command

// Command returns a new CommandBuilder instance which can be used to construct
// a Command instance.
//
// This function and Tree are the two entrypoints into the Argonaut library.
// This function returns a value that may be used in a call chain to construct
// a full-featured command line interface.
func Command() argo.CommandBuilder {
	return command.NewBuilder()
}

func BuildCommand(builder argo.CommandBuilder) (argo.Command, error) {
	return BuildCommandCustom(builder, argo.DefaultOptions())
}

func BuildCommandCustom(builder argo.CommandBuilder, options argo.Options) (argo.Command, error) {
	argo.FixOptions(&options)
	return command.Build(builder, options)
}

func MustBuildCommand(builder argo.CommandBuilder) argo.Command {
	return utils.MustReturn(BuildCommand(builder))
}

// ParseCommand attempts to parse the CLI call inputs as options and arguments
// to the given command type instance.
//
// This method returns two values, an argo.ParseResult struct containing any
// warnings and possibly an error from the attempt to parse the result, and an
// error value if an error occurred during the parse attempt.
//
// If the returned error value is not nil, it will also be set on the result
// struct, enabling callers to ignore one of the outputs if desired and have
// access to any error either way.
//
// Example 1:
//
//	// Use parse result
//	result, _ := cli.ParseCommand(myCommand)
//	if result.Error != nil {
//	    panic(result.Error)
//	}
//
//	for _, warning := result.Warnings {
//	    _, _ = fmt.Fprintf(os.Stderr, "parse warning: %s", warning.Message)
//	}
//
// Example 2:
//
//	// Disregard parse result
//	_, err := cli.ParseCommand(myCommand)
//	if err != nil {
//	    panic(err)
//	}
func ParseCommand(com argo.Command) (argo.ParseResult, error) {
	return command.Parse(com)
}

func MustParseCommand(com argo.Command) argo.ParseResult {
	return utils.MustReturn(ParseCommand(com))
}

// endregion Single Command

// region Command Tree

// Tree returns a new argo.CommandTreeBuilder instance which can be used to
// construct an argo.TreeCommand instance.
//
// A command tree is a tree of nested subcommands of arbitrary depth.  The tree
// consists of branch and leaf nodes, with the leaf nodes being the selectable
// final commands.
func Tree() argo.TreeCommandBuilder {
	return tree.NewBuilder()
}

// Branch returns a new BranchCommandBuilder instance which can be used to
// construct a Branch instance.
func Branch(name string) argo.BranchCommandBuilder {
	return tree.NewBranchBuilder(name)
}

// Leaf returns a new LeafCommandBuilder instance which can be used to construct
// a LeafCommand instance.
func Leaf(name string) argo.LeafCommandBuilder {
	return tree.NewLeafBuilder(name)
}

// CommandGroup returns a new CommandGroupBuilder in stance which can be used to
// construct a CommandGroup instance.
func CommandGroup(name string) argo.CommandGroupBuilder {
	return tree.NewGroupBuilder(name)
}

func BuildTree(builder argo.TreeCommandBuilder) (argo.TreeCommand, error) {
	return tree.Build(builder, argo.DefaultOptions())
}

func ParseTreeCLI(com argo.TreeCommand) (argo.ParseResult, error) {
	return tree.Parse(com, os.Args)
}

// endregion Command Tree

// region Flags

// FlagGroup returns a new FlagGroupBuilder instance which can be used to
// construct an FlagGroup instance.
func FlagGroup(name string) argo.FlagGroupBuilder {
	return flag.NewGroupBuilder(name)
}

// ShortFlag returns a new FlagBuilder instance with the short form already set
// to the given value.
//
// This function is a shortcut for:
//
//	cli.Render().WithShortForm(...)
func ShortFlag(f byte) argo.FlagBuilder {
	return flag.NewBuilder().WithShortForm(f)
}

// LongFlag returns a new FlagBuilder instance with the long form already set
// to the given value.
//
// This function is a shortcut for:
//
//	cli.Render().WithLongForm(...)
func LongFlag(name string) argo.FlagBuilder {
	return flag.NewBuilder().WithLongForm(name)
}

// Flag returns a new FlagBuilder instance which can be used to construct a Flag
// instance.
func Flag() argo.FlagBuilder {
	return flag.NewBuilder()
}

// ComboFlag returns a new FlagBuilder instance with the short and long forms
// already set to the given values.
//
// This function is a shortcut for:
//
//	cli.Render().WithShortForm(...).WithLongForm(...)
func ComboFlag(short byte, long string) argo.FlagBuilder {
	return flag.NewBuilder().WithShortForm(short).WithLongForm(long)
}

// endregion Flags

// Argument returns a new ArgumentBuilder instance which can be used to
// construct an Argument instance.
func Argument() argo.ArgumentBuilder {
	return argument.NewBuilder()
}
