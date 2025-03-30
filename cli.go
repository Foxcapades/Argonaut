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
	"fmt"
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

// BuildCommand attempts to build a new argo.Command instance based on the
// configuration applied to the given builder and the default argo.Options
// values.
//
// If the builder is misconfigured, this method will return nil, with an error
// containing information about the reason the command could not be built.
//
// The error will likely be an instance of argo.MultiError which will contain a
// list of all the configuration errors that were encountered while attempting
// to build the argo.Command.
func BuildCommand(builder argo.CommandBuilder) (argo.Command, error) {
	return BuildCommandCustom(builder, argo.Options{})
}

// BuildCommandCustom attempts to build a new argo.Command instance based on the
// configuration applied to the given builder and the given argo.Options value.
//
// If the builder is misconfigured, this method will return nil, with an error
// containing information about the reason the command could not be built.
//
// The error will likely be an instance of argo.MultiError which will contain a
// list of all the configuration errors that were encountered while attempting
// to build the argo.Command.
func BuildCommandCustom(builder argo.CommandBuilder, options argo.Options) (argo.Command, error) {
	return command.Build(builder, options)
}

// MustBuildCommand calls BuildCommand and panics if an error is returned.
//
// See BuildCommand for more information about the input and output.
func MustBuildCommand(builder argo.CommandBuilder) argo.Command {
	if out, err := BuildCommand(builder); err != nil {
		panic(err)
	} else {
		return out
	}
}

// MustBuildCommandCustom calls BuildCommandCustom and panics if an error is
// returned.
//
// See BuildCommandCustom for more information about the inputs and output.
func MustBuildCommandCustom(options argo.Options, builder argo.CommandBuilder) argo.Command {
	if out, err := BuildCommandCustom(builder, options); err != nil {
		panic(err)
	} else {
		return out
	}
}

// ParseCommand attempts to parse the CLI call inputs as options and arguments
// to an argo.Command.
//
// The ParseCommand method accepts as its argument one of the following types:
// * argo.Command
// * argo.CommandBuilder
//
// If a value of any other type is passed to ParseCommand, it will panic.
//
// For the argo.CommandBuilder type, MustBuildCommand will be called, to build
// an argo.Command instance, meaning an invalid builder config will cause a
// panic.
//
// ParseCommand returns three values, an argo.Command instance, an
// argo.ParseResult struct, and an error value if an error occurred during the
// parse attempt.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the argo.ParseResult's Error field.
//
// If the given argument value is an argo.Command instance, that same value will
// be the first return output.
func ParseCommand(com any) (out argo.Command, res argo.ParseResult, err error) {
	switch t := com.(type) {
	case argo.CommandBuilder:
		out = MustBuildCommand(t)
	case argo.Command:
		out = t
	default:
		panic(fmt.Sprintf("%T does not implement argo.Command or argo.CommandBuilder", com))
	}

	res, err = command.Parse(out, os.Args)
	return
}

// ParseCommandCustom attempts to parse the CLI call inputs as options and
// arguments to an argo.Command built from the given input, using the given
// argo.Options as configuration.
//
// MustBuildCommandCustom will be called to build the argo.Command instance,
// meaning an invalid builder config will cause a panic.
//
// ParseCommandCustom returns three values, an argo.Command instance, an
// argo.ParseResult struct, and an error value if an error occurred during the
// parse attempt.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the argo.ParseResult's Error field.
func ParseCommandCustom(options argo.Options, com argo.CommandBuilder) (out argo.Command, res argo.ParseResult, err error) {
	out = MustBuildCommandCustom(options, com)
	res, err = command.Parse(out, os.Args)
	return
}

// MustParseCommand calls Parse and panics if an error is returned.
//
// See ParseCommand for more information about the input and output.
func MustParseCommand(com any) (argo.Command, argo.ParseResult) {
	if out, res, err := ParseCommand(com); err != nil {
		panic(err)
	} else {
		return out, res
	}
}

// MustParseCommandCustom calls ParseCommandCustom and panics if an error is returned.
//
// See ParseCommandCustom for more information about the input and outputs.
func MustParseCommandCustom(options argo.Options, com argo.CommandBuilder) (argo.Command, argo.ParseResult) {
	if out, res, err := ParseCommandCustom(options, com); err != nil {
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

// BuildTree attempts to build a new argo.TreeCommand instance based on the
// configuration applied to the given builder and the default argo.Options
// values.
//
// If the builder is misconfigured, this method will return nil, with an error
// containing information about the reason the command could not be built.
//
// The error will likely be an instance of argo.MultiError which will contain a
// list of all the configuration errors that were encountered while attempting
// to build the argo.TreeCommand.
func BuildTree(builder argo.TreeCommandBuilder) (argo.TreeCommand, error) {
	return tree.Build(builder, argo.Options{})
}

// MustBuildTree calls BuildTree and panics if an error is returned.
//
// See BuildTree for more information about the input and output.
func MustBuildTree(builder argo.TreeCommandBuilder) argo.TreeCommand {
	if out, err := BuildTree(builder); err != nil {
		panic(err)
	} else {
		return out
	}
}

// BuildTreeCustom attempts to build a new argo.TreeCommand instance based on
// the configuration applied to the given builder and the given argo.Options
// value.
//
// If the builder is misconfigured, this method will return nil, with an error
// containing information about the reason the command could not be built.
//
// The error will likely be an instance of argo.MultiError which will contain a
// list of all the configuration errors that were encountered while attempting
// to build the argo.TreeCommand.
func BuildTreeCustom(options argo.Options, builder argo.TreeCommandBuilder) (argo.TreeCommand, error) {
	return tree.Build(builder, options)
}

// MustBuildTreeCustom calls BuildTreeCustom and panics if an error is
// returned.
//
// See BuildTreeCustom for more information about the inputs and outputs.
func MustBuildTreeCustom(options argo.Options, builder argo.TreeCommandBuilder) argo.TreeCommand {
	if out, err := BuildTreeCustom(options, builder); err != nil {
		panic(err)
	} else {
		return out
	}
}

// ParseTree attempts to parse the CLI call inputs as options and arguments
// to an argo.TreeCommand.
//
// The ParseTree method accepts as its argument one of the following types:
// * argo.TreeCommand
// * argo.TreeCommandBuilder
//
// If a value of any other type is passed to ParseTree, it will panic.
//
// For the argo.TreeCommandBuilder type, MustBuildTree will be called, to build
// an argo.TreeCommand instance, meaning an invalid builder config will cause a
// panic.
//
// ParseTree returns three values, an argo.TreeCommand instance, an
// argo.ParseResult struct, and an error value if an error occurred during the
// parse attempt.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the argo.ParseResult's Error field.
//
// If the given argument value is an argo.TreeCommand instance, that same value
// will be the first return output.
func ParseTree(com any) (out argo.TreeCommand, res argo.ParseResult, err error) {
	switch t := com.(type) {
	case argo.TreeCommandBuilder:
		out = MustBuildTree(t)
	case argo.TreeCommand:
		out = t
	default:
		panic(fmt.Sprintf("%T does not implement argo.TreeCommand or argo.TreeCommandBuilder", com))
	}

	res, err = tree.Parse(out, os.Args)
	return
}

// ParseTreeCustom attempts to parse the CLI call inputs as options and
// arguments to an argo.TreeCommand built from the given input, using the given
// argo.Options as configuration.
//
// MustBuildTreeCustom will be called to build the argo.TreeCommand instance,
// meaning an invalid builder config will cause a panic.
//
// ParseTreeCustom returns three values, the built argo.TreeCommand instance, an
// argo.ParseResult struct, and an error value if an error occurred during the
// parse attempt.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the argo.ParseResult's Error field.
func ParseTreeCustom(options argo.Options, com argo.TreeCommandBuilder) (out argo.TreeCommand, res argo.ParseResult, err error) {
	out = MustBuildTreeCustom(options, com)
	res, err = tree.Parse(out, os.Args)
	return
}

// MustParseTree calls ParseTree and panics if an error is returned.
//
// See ParseTree for more information about the input and outputs.
func MustParseTree(com any) (argo.TreeCommand, argo.ParseResult) {
	if out, res, err := ParseTree(com); err != nil {
		panic(err)
	} else {
		return out, res
	}
}

// MustParseTreeCustom calls ParseTreeCustom and panics if an error is returned.
//
// See ParseTreeCustom for more information about the input and outputs.
func MustParseTreeCustom(options argo.Options, com argo.TreeCommandBuilder) (argo.TreeCommand, argo.ParseResult) {
	if out, res, err := ParseTreeCustom(options, com); err != nil {
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
