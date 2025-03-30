package tree_test

import (
	"fmt"
	"strings"
	"testing"

	cli "github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	opts2 "github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const leafInput01 = `Usage:
  %s branch leaf -a [options] <arg> [files...]
  Aliases: l, le

    A description of this command leaf.

My Flags
    Flags that belong to me.

  -a | --amber
      Description of the amber flag.
  -b | --barn

Your Flags
  -c [arg] | --cantaloupe=[arg]

  -d <name> | --doorknob=<name>

Help Flags
  -h | --help
      Prints this help text.

Arguments
  <arg>
      This argument has a description.
`

func TestRenderCommandLeaf01(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithBranch(cli.Branch("branch").
			WithLeaf(cli.Leaf("leaf").
				WithAliases("l", "le").
				WithDescription("A description of this command leaf.").
				WithArgument(cli.Argument().Require().WithDescription("This argument has a description.")).
				WithFlagGroup(cli.FlagGroup("My Flags").
					WithDescription("Flags that belong to me.").
					WithFlag(cli.ShortFlag('a').
						Require().
						WithLongForm("amber").
						WithDescription("Description of the amber flag.")).
					WithFlag(cli.ShortFlag('b').WithLongForm("barn"))).
				WithFlagGroup(cli.FlagGroup("Your Flags").
					WithFlag(cli.ShortFlag('c').
						WithLongForm("cantaloupe").
						WithArgument(cli.Argument())).
					WithFlag(cli.ShortFlag('d').
						WithLongForm("doorknob").
						WithArgument(cli.Argument().WithName("name").Require()))).
				WithUnmappedInputLabel("files..."))), opt))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "branch", "leaf", "argument", "-a"}))

	renderLeafOutputCheck(t, leafInput01, com.SelectedCommand())
}

func renderLeafOutputCheck(
	t *testing.T,
	pattern string,
	com argo.LeafCommand,
) {
	sb := new(strings.Builder)
	opts := argo.Options{}
	opts2.FixOptions(&opts)

	utils.Must(tree.RenderLeafHelp(com, opts, sb))

	expected := fmt.Sprintf(pattern, commandName)

	if sb.String() != expected {
		t.Errorf("expected: '%s'\n\ngot: '%s'", expected, sb.String())
	}
}
