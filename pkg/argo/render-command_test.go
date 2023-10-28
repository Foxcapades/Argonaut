package argo_test

import (
	"bufio"
	"fmt"
	"os"
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
		WithArgument(cli.Argument()).
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

const regression51Expected = `Usage:
  %s [options] [FILE...]

    Concatenate FILE(s) to standard output.

General Flags
  -A | --show-all
      equivalent to -vET
  -b | --number-nonblank
      number nonempty output lines, overrides -n
  -e
      equivalent to -vE
  -E | --show-ends
      display $ at end of each line
  -n | --number
      number all output lines
  -s | --squeeze-blank
      suppress repeated empty output lines
  -t
      equivalent to -vT
  -T | --show-tabs
      display TAB characters as ^I
  -v | --show-nonprinting
      use ^ and M- notation, except for LFD and TAB
  --version
      output version information and exit

Help Flags
  -h | --help
      Prints this help text.
`

func TestCommandHelpRenderer_regression51(t *testing.T) {
	type Config struct {
		NumberNonBlank  bool
		NumberLines     bool
		SqueezeBlank    bool
		ShowTabs        bool
		ShowEnds        bool
		ShowNonPrinting bool
	}

	var config Config

	com, err := cli.Command().
		WithDescription("Concatenate FILE(s) to standard output.").
		WithFlag(cli.ComboFlag('A', "show-all").
			WithDescription("equivalent to -vET").
			WithCallback(func(_ argo.Flag) {
				config.ShowNonPrinting = true
				config.ShowEnds = true
				config.ShowTabs = true
			})).
		WithFlag(cli.ComboFlag('b', "number-nonblank").
			WithDescription("number nonempty output lines, overrides -n").
			WithBinding(&config.NumberNonBlank, false)).
		WithFlag(cli.ShortFlag('e').
			WithDescription("equivalent to -vE").
			WithCallback(func(_ argo.Flag) {
				config.ShowNonPrinting = true
				config.ShowEnds = true
			})).
		WithFlag(cli.ComboFlag('E', "show-ends").
			WithDescription("display $ at end of each line").
			WithBinding(&config.ShowEnds, false)).
		WithFlag(cli.ComboFlag('n', "number").
			WithDescription("number all output lines").
			WithBinding(&config.NumberLines, false)).
		WithFlag(cli.ComboFlag('s', "squeeze-blank").
			WithDescription("suppress repeated empty output lines").
			WithBinding(&config.SqueezeBlank, false)).
		WithFlag(cli.ShortFlag('t').
			WithDescription("equivalent to -vT").
			WithCallback(func(_ argo.Flag) {
				config.ShowNonPrinting = true
				config.ShowTabs = true
			})).
		WithFlag(cli.ComboFlag('T', "show-tabs").
			WithDescription("display TAB characters as ^I").
			WithBinding(&config.ShowTabs, false)).
		WithFlag(cli.ComboFlag('v', "show-nonprinting").
			WithDescription("use ^ and M- notation, except for LFD and TAB").
			WithBinding(&config.ShowNonPrinting, false)).
		WithFlag(cli.LongFlag("version").
			WithDescription("output version information and exit").
			WithCallback(func(_ argo.Flag) {
				fmt.Println("<version information>")
				os.Exit(0)
			})).
		WithUnmappedLabel("FILE...").
		Build(nil)

	if err != nil {
		t.Error("expected err to be nil but was", err)
	} else {
		renderOutputCheck(t, regression51Expected, com, argo.CommandHelpRenderer())
	}
}

const commandHelpRendererExpectOptionalArgs = `Usage:
  %s [options]

Flags
  -a [arg] | --all=[arg]

  -h | --help
      Prints this help text.
`

func TestCommandHelpRenderer_optionalArgs(t *testing.T) {
	var bind string
	com, err := cli.Command().
		WithFlag(cli.ComboFlag('a', "all").WithBinding(&bind, false)).
		Build(nil)

	if err != nil {
		t.Error("expected err to be nil but was", err)
	} else {
		renderOutputCheck(t, commandHelpRendererExpectOptionalArgs, com, argo.CommandHelpRenderer())
	}
}
