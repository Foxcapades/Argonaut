package tree_test

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	opts2 "github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
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

Inherited Flags
  -h | --help
      Prints this help text.

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
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
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
					WithDescription("Hail Satan")))), opt))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "branch1", "cruise", "-a"}))

	renderBranchOutputCheck(t, branchHelp001, com.SelectedCommand().Parent().(argo.BranchCommand))
}

func TestCommandBranchHelpRendererFail01(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
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
					WithArgument(argument.NewBuilder()))).
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
				WithDescription("Hail Satan"))), opt))

	utils.MustReturn(tree.Parse(com, []string{"command", "branch1", "cruise", "-a"}))

	for p := 1; p <= len(branchHelp001); p++ {
		wri := FailingWriter{FailAfter: p}
		buf := bufio.NewWriterSize(&wri, 1)

		err := tree.RenderBranchHelp(com.SelectedCommand().Parent().(argo.BranchCommand), opt, buf)
		if err == nil {
			t.Error("expected err to not be nil but it was")
		}
	}
}

func renderBranchOutputCheck(
	t *testing.T,
	pattern string,
	com argo.BranchCommand,
) {
	sb := new(strings.Builder)
	opts := argo.Options{}
	opts2.FixOptions(&opts)

	utils.Must(tree.RenderBranchHelp(com, opts, sb))

	expected := fmt.Sprintf(pattern, commandName)

	if sb.String() != expected {
		t.Errorf("expected: '%s'\n\ngot: '%s'", expected, sb.String())
	}
}

type FailingWriter struct {
	FailAfter int
	current   int
}

func (f *FailingWriter) Write(p []byte) (n int, err error) {
	if f.current < f.FailAfter {
		f.current++
		return len(p), nil
	} else {
		return 0, errors.New("fake error")
	}
}
