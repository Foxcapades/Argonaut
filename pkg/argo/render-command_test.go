package argo_test

import (
	"bufio"
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const commandHelp001 = `Usage:
  %s -m [options] [argument] [arg2] [asteroids...]

    A command description.

Meta Flags
    Meat flags.

  -m
      Enables so much meat.

Help Flags
  -h | --help
      Prints this help text.

Arguments
  [argument]
      poo

  [arg2]
`

func TestCommandHelpRenderer001(t *testing.T) {
	com := argo.NewCommandBuilder().
		WithDescription("A command description.").
		WithFlagGroup(cli.FlagGroup("Meta Flags").
			WithDescription("Meat flags.").
			WithFlag(cli.ShortFlag('m').
				WithDescription("Enables so much meat.").
				Require())).
		WithArgument(cli.Argument().
			WithName("argument").
			WithDescription("poo")).
		WithArgument(cli.Argument()).
		WithUnmappedLabel("asteroids...").
		MustParse([]string{"command", "-m"})

	renderOutputCheck(t, commandHelp001, com, argo.CommandHelpRenderer())
}

func TestCommandHelpRendererFail01(t *testing.T) {
	com := argo.NewCommandBuilder().
		WithDescription("A command description.").
		WithFlagGroup(cli.FlagGroup("Meta Flags").
			WithDescription("Meat flags.").
			WithFlag(cli.ShortFlag('m').
				WithDescription("Enables so much meat.").
				Require())).
		WithArgument(cli.Argument().
			WithName("argument").
			WithDescription("poo")).
		WithUnmappedLabel("asteroids...").
		MustParse([]string{"command", "-m"})
	ren := argo.CommandHelpRenderer()

	for p := 1; p <= len(commandHelp001); p++ {
		wri := FailingWriter{FailAfter: p}
		buf := bufio.NewWriterSize(&wri, 1)

		err := ren.RenderHelp(com, buf)
		if err == nil {
			t.Error("expected err to not be nil but it was")
		}
	}
}
