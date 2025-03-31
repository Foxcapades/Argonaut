package tree_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const branchHelp001 = `Usage:
  %s branch1 -a=[arg] [options] <command>
  Aliases: branch2, branch3

    A description of this command.

My Special Little Commands
    A category of commands that are
    special and
    little.

  666
      Hail Satan
  cruise        Aliases: flight
  prescriber    Aliases: doctor, nurse-practitioner
      A description.

Super Flags
    A group of flags that are just super.

  -a [arg] | --apple=[arg]
      A description of the apple flag.
  -b <arg> | --bear=<arg>

Boring Flags
  -d | --diameter

  -e | --ergonomics

Inherited Flags
  -h | --help
      Prints this help text.
`

func TestCommandBranchHelpRenderer001(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithOptions(argo.TreeCommandOptions{InheritParentFlags: argo.FlagInheritanceEnabled}).
		WithBranch(tree.NewBranchBuilder("branch1").
			WithAliases("branch2", "branch3").
			WithDescription("A description of this command.").
			WithHelpDisabled().
			WithFlagGroup(flag.NewGroupBuilder("Super Flags").
				WithDescription("A group of flags that are just super.").
				WithFlag(flag.NewBuilder().
					WithShortForm('a').
					WithLongForm("apple").
					WithDescription("A description of the apple flag.").
					Require().
					WithArgument(argument.NewBuilder())).
				WithFlag(flag.NewBuilder().
					WithShortForm('b').
					WithLongForm("bear").
					WithArgument(argument.NewBuilder().Require()))).
			WithFlagGroup(flag.NewGroupBuilder("Boring Flags").
				WithFlag(flag.NewBuilder().
					WithShortForm('d').
					WithLongForm("diameter")).
				WithFlag(flag.NewBuilder().
					WithShortForm('e').
					WithLongForm("ergonomics"))).
			WithCommandGroup(tree.NewGroupBuilder("My Special Little Commands").
				WithDescription("A category of commands that are\nspecial and\nlittle.").
				WithLeaf(tree.NewLeafBuilder("cruise").
					WithAliases("flight")).
				WithBranch(tree.NewBranchBuilder("prescriber").
					WithAliases("doctor", "nurse-practitioner").
					WithDescription("A description.").
					WithLeaf(tree.NewLeafBuilder("dethrone"))).
				WithLeaf(tree.NewLeafBuilder("666").
					WithDescription("Hail Satan"))))))

	renderBranchOutputCheck(t, branchHelp001, com.FindChild("branch1").(argo.BranchCommand), com.Options())
}

const branchHelp002 = `Usage:
  %s branch1 -a=[arg] cruise [options]
  Aliases: flight

Flags
  -h | --help
      Prints this help text.

Inherited Flags
  -a [arg] | --apple=[arg]
      A description of the apple flag.
  -b <arg> | --bear=<arg>

  -d | --diameter

  -e | --ergonomics
`

// Required parent flags shown in the help for child flags.
func TestCommandBranchHelpRenderer002(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithOptions(argo.TreeCommandOptions{InheritParentFlags: argo.FlagInheritanceEnabled}).
		WithBranch(tree.NewBranchBuilder("branch1").
			WithAliases("branch2", "branch3").
			WithDescription("A description of this command.").
			WithHelpDisabled().
			WithFlagGroup(flag.NewGroupBuilder("Super Flags").
				WithDescription("A group of flags that are just super.").
				WithFlag(flag.NewBuilder().
					WithShortForm('a').
					WithLongForm("apple").
					WithDescription("A description of the apple flag.").
					Require().
					WithArgument(argument.NewBuilder())).
				WithFlag(flag.NewBuilder().
					WithShortForm('b').
					WithLongForm("bear").
					WithArgument(argument.NewBuilder().Require()))).
			WithFlagGroup(flag.NewGroupBuilder("Boring Flags").
				WithFlag(flag.NewBuilder().
					WithShortForm('d').
					WithLongForm("diameter")).
				WithFlag(flag.NewBuilder().
					WithShortForm('e').
					WithLongForm("ergonomics"))).
			WithCommandGroup(tree.NewGroupBuilder("My Special Little Commands").
				WithDescription("A category of commands that are\nspecial and\nlittle.").
				WithLeaf(tree.NewLeafBuilder("cruise").
					WithAliases("flight")).
				WithBranch(tree.NewBranchBuilder("prescriber").
					WithAliases("doctor", "nurse-practitioner").
					WithDescription("A description.").
					WithLeaf(tree.NewLeafBuilder("dethrone"))).
				WithLeaf(tree.NewLeafBuilder("666").
					WithDescription("Hail Satan"))))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "branch1", "cruise", "-a"}))

	renderBranchOutputCheck(t, branchHelp002, com.SelectedCommand().Parent().(argo.BranchCommand), com.Options())
}

func renderBranchOutputCheck(
	t *testing.T,
	pattern string,
	com argo.BranchCommand,
	opts argo.TreeCommandOptions,
) {
	sb := new(strings.Builder)

	utils.Must(tree.RenderBranchHelp(com, tree.WrapOptions(&opts), sb))

	expected := fmt.Sprintf(pattern, commandName)

	if sb.String() != expected {
		t.Errorf("expected: '%s'\n\ngot: '%s'", expected, sb.String())
	}
}
