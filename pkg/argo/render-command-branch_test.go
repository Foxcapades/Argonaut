package argo_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const branchHelp001 = `Usage:
  %s branch1 [options] <command>
    A description of this command.

Super Flags
    A group of flags that are just super.

  -a [arg] | --apple=[arg]
      A description of the apple flag.
  -b <arg> | --bear=<arg>

Boring Flags
  -d | --diameter

  -e | --ergonomics

My Special Little Commands
    A category of commands that are
    special and
    little.

  666
      Hail Satan
  cruise
  physician
      A description.
`

func TestCommandBranchHelpRenderer001(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithBranch(argo.NewCommandBranchBuilder("branch1").
			WithDescription("A description of this command.").
			WithHelpDisabled().
			WithFlagGroup(argo.NewFlagGroupBuilder("Super Flags").
				WithDescription("A group of flags that are just super.").
				WithFlag(argo.NewFlagBuilder().
					WithShortForm('a').
					WithLongForm("apple").
					WithDescription("A description of the apple flag.").
					WithArgument(argo.NewArgumentBuilder())).
				WithFlag(argo.NewFlagBuilder().
					WithShortForm('b').
					WithLongForm("bear").
					WithArgument(argo.NewArgumentBuilder().Require()))).
			WithFlagGroup(argo.NewFlagGroupBuilder("Boring Flags").
				WithFlag(argo.NewFlagBuilder().
					WithShortForm('d').
					WithLongForm("diameter")).
				WithFlag(argo.NewFlagBuilder().
					WithShortForm('e').
					WithLongForm("ergonomics"))).
			WithCommandGroup(argo.NewCommandGroupBuilder("My Special Little Commands").
				WithDescription("A category of commands that are\nspecial and\nlittle.").
				WithLeaf(argo.NewCommandLeafBuilder("cruise")).
				WithBranch(argo.NewCommandBranchBuilder("physician").
					WithDescription("A description.").
					WithLeaf(argo.NewCommandLeafBuilder("dethrone"))).
				WithLeaf(argo.NewCommandLeafBuilder("666").
					WithDescription("Hail Satan")))).
		MustParse([]string{"command", "branch1", "cruise"})

	renderOutputCheck(t, branchHelp001, com.SelectedCommand().Parent().(argo.CommandBranch), argo.CommandBranchHelpRenderer())
}
