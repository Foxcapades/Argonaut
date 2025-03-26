package tree_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
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
	renderTreeOutputCheck(t, output001, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")))
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
	renderTreeOutputCheck(t, output002, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithShortForm('c')))
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
	renderTreeOutputCheck(t, output003, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().
			WithLongForm("hello").
			WithDescription("Some description of the flag.")))
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
	renderTreeOutputCheck(t, output004, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithDescription("Help text about the leaf.")))
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
	renderTreeOutputCheck(t, output005, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf1")).
		WithLeaf(tree.NewLeafBuilder("leaf2")))
}

const treeHelp006 = `Usage:
  %s -c=<arg> [options] <command>

Flags
  -c <arg>

  -h | --help
      Prints this help text.

Commands
  leaf
`

func TestCommandTreeHelpRenderer006(t *testing.T) {
	renderTreeOutputCheck(t, treeHelp006, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().
			WithShortForm('c').
			Require().
			WithArgument(argument.NewBuilder().Require())))
}

const treeHelp007 = `Usage:
  %s [options] <command>

Something
  -c

Help Flags
  -h | --help
      Prints this help text.

Commands
  leaf
`

func TestCommandTreeHelpRenderer007(t *testing.T) {
	renderTreeOutputCheck(t, treeHelp007, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlagGroup(flag.NewGroupBuilder("Something").
			WithFlag(flag.NewBuilder().WithShortForm('c'))))
}

const treeHelp008 = `Usage:
  %s [options] <command>

Something
    A description of something.

  -c

Help Flags
  -h | --help
      Prints this help text.

Commands
  leaf
`

func TestCommandTreeHelpRenderer008(t *testing.T) {
	renderTreeOutputCheck(t, treeHelp008, tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlagGroup(flag.NewGroupBuilder("Something").
			WithDescription("A description of something.").
			WithFlag(flag.NewBuilder().WithShortForm('c'))))
}

func renderTreeOutputCheck(
	t *testing.T,
	pattern string,
	com argo.TreeCommandBuilder,
) {
	sb := new(strings.Builder)
	opts := argo.DefaultOptions()

	utils.Must(tree.RenderHelp(utils.MustReturn(tree.Build(com, opts)), opts, sb))

	expected := fmt.Sprintf(pattern, commandName)

	if sb.String() != expected {
		t.Errorf("expected: '%s'\n\ngot: '%s'", expected, sb.String())
	}
}
