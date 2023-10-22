package argo_test

import (
	"bufio"
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const branchHelp001 = `Usage:
  %s branch1 -a=[arg] [options] <command>
  Aliases: branch2, branch3

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
  cruise        Aliases: flight
  prescriber    Aliases: doctor, nurse-practitioner
      A description.
`

func TestCommandBranchHelpRenderer001(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithBranch(argo.NewCommandBranchBuilder("branch1").
			WithAliases("branch2", "branch3").
			WithDescription("A description of this command.").
			WithHelpDisabled().
			WithFlagGroup(argo.NewFlagGroupBuilder("Super Flags").
				WithDescription("A group of flags that are just super.").
				WithFlag(argo.NewFlagBuilder().
					WithShortForm('a').
					WithLongForm("apple").
					WithDescription("A description of the apple flag.").
					Require().
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
				WithLeaf(argo.NewCommandLeafBuilder("cruise").
					WithAliases("flight")).
				WithBranch(argo.NewCommandBranchBuilder("prescriber").
					WithAliases("doctor", "nurse-practitioner").
					WithDescription("A description.").
					WithLeaf(argo.NewCommandLeafBuilder("dethrone"))).
				WithLeaf(argo.NewCommandLeafBuilder("666").
					WithDescription("Hail Satan")))).
		MustParse([]string{"command", "branch1", "cruise", "-a"})

	renderOutputCheck(t, branchHelp001, com.SelectedCommand().Parent().(argo.CommandBranch), argo.CommandBranchHelpRenderer())
}

func TestCommandBranchHelpRendererFail01(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithBranch(argo.NewCommandBranchBuilder("branch1").
			WithAliases("branch2", "branch3").
			WithDescription("A description of this command.").
			WithHelpDisabled().
			WithFlagGroup(argo.NewFlagGroupBuilder("Super Flags").
				WithDescription("A group of flags that are just super.").
				WithFlag(argo.NewFlagBuilder().
					WithShortForm('a').
					WithLongForm("apple").
					WithDescription("A description of the apple flag.").
					Require().
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
				WithLeaf(argo.NewCommandLeafBuilder("cruise").
					WithAliases("flight")).
				WithBranch(argo.NewCommandBranchBuilder("prescriber").
					WithAliases("doctor", "nurse-practitioner").
					WithDescription("A description.").
					WithLeaf(argo.NewCommandLeafBuilder("dethrone"))).
				WithLeaf(argo.NewCommandLeafBuilder("666").
					WithDescription("Hail Satan")))).
		MustParse([]string{"command", "branch1", "cruise", "-a"})
	ren := argo.CommandBranchHelpRenderer()

	for p := 1; p <= len(branchHelp001); p++ {
		wri := FailingWriter{FailAfter: p}
		buf := bufio.NewWriterSize(&wri, 1)

		err := ren.RenderHelp(com.SelectedCommand().Parent().(argo.CommandBranch), buf)
		if err == nil {
			t.Error("expected err to not be nil but it was")
		}
	}
}
