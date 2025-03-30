package command_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/foxcapades/argonaut/v3"
	command "github.com/foxcapades/argonaut/v3/internal/argo/command/simple"
	opts2 "github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const commandHelp001 = `Usage:
  %s -m [options] [argument] [arg2] [asteroids...]

    A command description.

Meta Flags
    Meat flags.

  -m
      Enables so much meat.

General Flags
  -h | --help
      Prints this help text.

Arguments
  [argument]
      poo

  [arg2]
`

func TestCommandHelpRenderer001(t *testing.T) {
	com := utils.MustReturn(command.Build(command.NewBuilder().
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
		WithUnmappedInputLabel("asteroids..."), argo.Options{}))

	// cli.MustParseCommand(com)

	renderOutputCheck(t, commandHelp001, com)
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

General Flags
  -h | --help
      Prints this help text.
`

// https://github.com/Foxcapades/Argonaut/issues/51
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

	com, err := command.Build(command.NewBuilder().
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
				utils.Exit(0)
			})).
		WithUnmappedInputLabel("FILE..."), argo.Options{
		MaxDefaultFlagGroupSizeForMetaGroup: 5,
	})

	if err != nil {
		t.Error("expected err to be nil but was", err)
	} else {
		renderOutputCheck(t, regression51Expected, com)
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
	com, err := command.Build(command.NewBuilder().
		WithFlag(cli.ComboFlag('a', "all").WithBinding(&bind, false)), argo.Options{})

	if err != nil {
		t.Error("expected err to be nil but was", err)
	} else {
		renderOutputCheck(t, commandHelpRendererExpectOptionalArgs, com)
	}
}

func renderOutputCheck(
	t *testing.T,
	pattern string,
	com argo.Command,
) {
	sb := new(strings.Builder)
	opts := argo.Options{}
	opts2.FixOptions(&opts)

	utils.Must(command.RenderHelp(com, opts, sb))

	expected := fmt.Sprintf(pattern, filepath.Base(os.Args[0]))

	if sb.String() != expected {
		t.Errorf("expected: '%s'\n\ngot: '%s'", expected, sb.String())
	}
}
