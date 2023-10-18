package argo_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

var commandName = filepath.Base(os.Args[0])

const output001 = `Usage:
  %s [options] <command>

Flags
  -h | --help
      Prints this help text.

Commands
  leaf
`

func TestCommandTreeHelpRenderer001(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("leaf")).
		MustParse([]string{"command", "leaf"})

	renderOutputCheck(t, output001, com, argo.CommandTreeHelpRenderer())
}

const output002 = `Usage:
  %s [options] <command>

Flags
  -c

  -h | --help
      Prints this help text.

Commands
  leaf
`

func TestCommandTreeHelpRenderer002(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("leaf")).
		WithFlag(argo.NewFlagBuilder().WithShortForm('c')).
		MustParse([]string{"command", "leaf"})

	renderOutputCheck(t, output002, com, argo.CommandTreeHelpRenderer())
}

const output003 = `Usage:
  %s [options] <command>

Flags
  --hello
      Some description of the flag.
  -h | --help
      Prints this help text.

Commands
  leaf
`

func TestCommandTreeHelpRenderer003(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("leaf")).
		WithFlag(argo.NewFlagBuilder().
			WithLongForm("hello").
			WithDescription("Some description of the flag.")).
		MustParse([]string{"command", "leaf"})

	renderOutputCheck(t, output003, com, argo.CommandTreeHelpRenderer())
}

const output004 = `Usage:
  %s [options] <command>

Flags
  -h | --help
      Prints this help text.

Commands
  leaf
      Help text about the leaf.
`

func TestCommandTreeHelpRenderer004(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("leaf").
			WithDescription("Help text about the leaf.")).
		MustParse([]string{"command", "leaf"})

	renderOutputCheck(t, output004, com, argo.CommandTreeHelpRenderer())
}

const output005 = `Usage:
  %s [options] <command>

Flags
  -h | --help
      Prints this help text.

Commands
  leaf1
  leaf2
`

func TestCommandTreeHelpRenderer005(t *testing.T) {
	com := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("leaf1")).
		WithLeaf(argo.NewCommandLeafBuilder("leaf2")).
		MustParse([]string{"command", "leaf1"})

	renderOutputCheck(t, output005, com, argo.CommandTreeHelpRenderer())
}

func renderOutputCheck[T any](
	t *testing.T,
	pattern string,
	command T,
	renderer argo.HelpRenderer[T],
) {
	buf := new(strings.Builder)

	err := renderer.RenderHelp(command, buf)

	if err != nil {
		t.Fail()
	}

	expected := fmt.Sprintf(pattern, commandName)

	if buf.String() != expected {
		t.Errorf("expected: '%s'\n\ngot: '%s'", expected, buf.String())
		t.Fail()
	}
}
