// Package cli provides a convenience methods for constructing command line
// interfaces.
//
// Commands may be either singular commands or command trees.  A singular
// command has no subcommands and may take any number of flags and/or arguments.
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
// CLI interface construction starts with either the `cli.Command` function or
// the `cli.Tree` function.
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
	return BuildCommandCustom(builder, argo.DefaultOptions())
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
	argo.FixOptions(&options)
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
func MustBuildCommandCustom(builder argo.CommandBuilder, options argo.Options) argo.Command {
	if out, err := BuildCommandCustom(builder, options); err != nil {
		panic(err)
	} else {
		return out
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
func Tree() argo.TreeCommandBuilder {
	return tree.NewBuilder()
}

// Branch returns a new BranchCommandBuilder instance which can be used to
// construct an argo.BranchCommand instance as part of a root argo.TreeCommand.
//
// The input value is the primary name of the branch command, additional aliases
// may optionally be added via the argo.BranchCommandBuilder methods.
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
	return tree.Build(builder, argo.DefaultOptions())
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
func BuildTreeCustom(builder argo.TreeCommandBuilder, options argo.Options) (argo.TreeCommand, error) {
	argo.FixOptions(&options)
	return tree.Build(builder, options)
}

// MustBuildTreeCustom calls BuildTreeCustom and panics if an error is
// returned.
//
// See BuildTreeCustom for more information about the inputs and output.
func MustBuildTreeCustom(builder argo.TreeCommandBuilder, options argo.Options) argo.TreeCommand {
	if out, err := BuildTreeCustom(builder, options); err != nil {
		panic(err)
	} else {
		return out
	}
}

// endregion Command Tree

// region Parse

// Parse attempts to parse the CLI call inputs as options and arguments to the
// given command type instance.
//
// The Parse method accepts as its argument one of the following types:
// * argo.Command
// * argo.TreeCommand
// * argo.CommandBuilder
// * argo.TreeCommandBuilder
//
// If a value of any other type is passed to Parse, it will panic.
//
// For the builder type cases, the appropriate MustBuild* function will be
// called, meaning an invalid builder config will cause a panic.
//
// Parse returns two values, an argo.ParseResult struct containing any warnings
// and possibly an error from the attempt to parse the result, and an error
// value if an error occurred during the parse attempt.
//
// If the returned error value is not nil, that same error will also be set as
// the value of the result's Error field.  This enables callers to ignore one of
// the outputs if desired and have access to any error either way.
//
// Example 1:
//
//	// Use parse result
//	result, _ := cli.Parse(myCommand)
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
func Parse(com any) (argo.ParseResult, error) {
	if typed, ok := com.(argo.CommandBuilder); ok {
		return command.Parse(MustBuildCommand(typed), os.Args)
	}

	if typed, ok := com.(argo.TreeCommandBuilder); ok {
		return tree.Parse(MustBuildTree(typed), os.Args)
	}

	if typed, ok := com.(argo.Command); ok {
		return command.Parse(typed, os.Args)
	}

	if typed, ok := com.(argo.TreeCommand); ok {
		return tree.Parse(typed, os.Args)
	}

	panic(fmt.Sprintf("given value %T does not implement any supported Argonaut command or command builder type", com))
}

// MustParse calls Parse and panics if an error is returned.
//
// See Parse for more information about the input and output.
func MustParse(com any) argo.ParseResult {
	if res, err := Parse(com); err != nil {
		panic(err)
	} else {
		return res
	}
}

// endregion Parse

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
// This function is a shortcut for:
//
//	cli.Flag().WithLongForm(...)
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
