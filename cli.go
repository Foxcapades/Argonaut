// Package cli is the entrypoint into the Argonaut library functionality. For
// most simple cases, no other package need be imported to configure a CLI
// interface.
//
// This package provide methods for creating simple commands, command trees,
// flag options, arguments, and additional organization functionality.
//
// CLI interfaces may be constructed with a root of either an argo.Command
// (using the Command function), or an argo.TreeCommand (using the Tree
// function).  Flags, arguments, and subcommands (of command trees) may then be
// appended via builder methods to define the CLI call API.
//
// Simple (or singular) commands have no subcommands and may take any number of
// flags and/or positional arguments.
//
// Tree commands have a 'root' which is the callable binary, which may be
// followed by any number of flags, and must be followed by one or subcommands
// in the form of branches or leaves.  Each branch may also have its own flags
// and also must have at least one subcommand which may be another branch or a
// leaf.  Leaf nodes may accept option flags as well as positional arguments.
package cli

import (
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/simple"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// region Single Command

// Command returns a new argo.CommandBuilder instance which can be used to
// construct an argo.Command instance.
//
// An example single command:
//
//	tar -xf foo.tgz
//
// Here the single command 'tar' accepts the flags '-x' and '-f', with the '-f'
// flag taking the argument 'foo.tgz'.
func Command() argo.CommandBuilder {
	return command.NewBuilder()
}

// ParseCommand attempts to parse the CLI call inputs as options and arguments
// to a new argo.Command instance built from the given argo.CommandBuilder.
//
// ParseCommand returns three values, an argo.Command instance (which may be nil
// on error), an argo.ParseResult struct, and an error value.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the argo.ParseResult's Error field.
func ParseCommand(builder argo.CommandBuilder) (out argo.Command, res argo.ParseResult, err error) {
	if out, err = command.Build(builder); err != nil {
		res.Error = err
		res.ErrorType = argo.ConfigurationError
		return
	}

	if res.InputWarnings, err = command.Parse(out, os.Args); err != nil {
		res.Error = err
		res.ErrorType = argo.InputError
	}

	return
}

// MustParseCommand calls Parse and panics if an error is returned.
//
// See ParseCommand for more information about the input and output.
func MustParseCommand(com argo.CommandBuilder) (argo.Command, argo.ParseResult) {
	if out, res, err := ParseCommand(com); err != nil {
		panic(err)
	} else {
		return out, res
	}
}

// endregion Single Command

// region Command Tree

// Tree returns a new argo.CommandTreeBuilder instance which can be used to
// construct an argo.TreeCommand instance.
//
// A command tree is a tree of nested subcommands of arbitrary depth.  The tree
// consists of branch and leaf nodes, with the leaf nodes being the selectable
// final commands.
//
// An example command tree:
//
//	docker compose -f my-docker-compose.yml up my-service
//
// Here the command tree is constructed of 3 levels, the root of the tree
// 'docker', the intermediary branch 'compose' and the leaf command 'up'.  The
// branch command is taking a flag '-f' which is itself taking the argument
// 'my-docker-compose.yml'.  The leaf command 'up' is accepting the optional,
// positional argument 'my-service'.
func Tree() argo.TreeCommandBuilder {
	return tree.NewBuilder()
}

// Branch returns a new BranchCommandBuilder instance which can be used to
// construct an argo.BranchCommand instance as part of a root argo.TreeCommand.
//
// The input value is the primary name of the branch command, additional aliases
// may optionally be added via the argo.BranchCommandBuilder methods.
//
// Branch names must begin with an alphanumeric character or underscore, and be
// followed by zero or more characters that are alphanumeric, underscores, or
// hyphens.
//
// Argonaut provides no public functions for building a branch alone, they may
// only be built as part of an argo.TreeCommandBuilder build via BuildTree.
func Branch(name string) argo.BranchCommandBuilder {
	return tree.NewBranchBuilder(name)
}

// Leaf returns a new LeafCommandBuilder instance which can be used to construct
// a LeafCommand instance.
//
// The input value is the primary name of the leaf command, additional aliases
// may optionally be added via the argo.LeafCommandBuilder methods.
//
// Leaf names must begin with an alphanumeric character or underscore, and be
// followed by zero or more characters that are alphanumeric, underscores, or
// hyphens.
//
// Argonaut provides no public functions for building a leaf alone, they may
// only be built as part of an argo.TreeCommandBuilder build via BuildTree.
func Leaf(name string) argo.LeafCommandBuilder {
	return tree.NewLeafBuilder(name)
}

// CommandGroup returns a new argo.CommandGroupBuilder in stance which can be
// used to construct an argo.CommandGroup instance.
//
// The input value is the display name of the command group.
func CommandGroup(name string) argo.CommandGroupBuilder {
	return tree.NewGroupBuilder(name)
}

// ParseTree attempts to parse the CLI call inputs as options and arguments
// to an argo.TreeCommand instance built from the given argo.TreeCommandBuilder.
//
// ParseTree returns three values, an argo.TreeCommand instance (which may be
// nil on error), an argo.ParseResult struct, and an error value if an error
// occurred during the parse attempt.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the argo.ParseResult's Error field.
func ParseTree(builder argo.TreeCommandBuilder) (out argo.TreeCommand, res argo.ParseResult, err error) {
	if out, err = tree.Build(builder); err != nil {
		res.Error = err
		res.ErrorType = argo.ConfigurationError
		return
	}

	if res.InputWarnings, err = tree.Parse(out, os.Args); err != nil {
		res.Error = err
		res.ErrorType = argo.InputError
	}

	return
}

// MustParseTree calls ParseTree and panics if an error is returned.
//
// See ParseTree for more information about the input and outputs.
func MustParseTree(com argo.TreeCommandBuilder) (argo.TreeCommand, argo.ParseResult) {
	if out, res, err := ParseTree(com); err != nil {
		panic(err)
	} else {
		return out, res
	}
}

// endregion Command Tree

// region Flags

// FlagGroup returns a new FlagGroupBuilder instance which can be used to
// construct an FlagGroup instance.
//
// The input value sets the display name of the flag group.
func FlagGroup(name string) argo.FlagGroupBuilder {
	return flag.NewGroupBuilder(name)
}

// ShortFlag returns a new argo.FlagBuilder instance with the short form set to
// the given value.
//
// This function is a shortcut for:
//
//	cli.Flag().WithShortForm(...)
func ShortFlag(f byte) argo.FlagBuilder {
	return flag.NewBuilder().WithShortForm(f)
}

// LongFlag returns a new argo.FlagBuilder instance with the long form set to
// the given value.
//
// This function is a convenience shortcut for:
//
//	cli.Flag().WithLongForm(...)
//
// Input value must be a valid flag name as defined by the argo.FlagBuilder
// WithLongForm method.
func LongFlag(name string) argo.FlagBuilder {
	return flag.NewBuilder().WithLongForm(name)
}

// Flag returns a new argo.FlagBuilder instance which can be used to construct
// an argo.Flag instance.
func Flag() argo.FlagBuilder {
	return flag.NewBuilder()
}

// ComboFlag returns a new argo.FlagBuilder instance with the short and long
// forms set to the given values.
//
// This function is a shortcut for:
//
//	cli.Flag().WithShortForm(...).WithLongForm(...)
func ComboFlag(short byte, long string) argo.FlagBuilder {
	return flag.NewBuilder().WithShortForm(short).WithLongForm(long)
}

// endregion Flags

// Argument returns a new argo.ArgumentBuilder instance which can be used to
// construct an argo.Argument instance.
func Argument() argo.ArgumentBuilder {
	return argument.NewBuilder()
}
