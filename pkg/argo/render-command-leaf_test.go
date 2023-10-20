package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
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
	com := cli.Tree().
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
				WithUnmappedLabel("files..."))).
		MustParse([]string{"command", "branch", "leaf", "argument", "-a"})

	renderOutputCheck(t, leafInput01, com.SelectedCommand(), argo.CommandLeafHelpRenderer())
}
