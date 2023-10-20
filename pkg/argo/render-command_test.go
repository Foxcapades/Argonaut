package argo_test

import (
	"bufio"
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const commandHelp001 = `Usage:
  %s [options]

Flags
  -h | --help
      Prints this help text.
`

func TestCommandHelpRenderer001(t *testing.T) {
	com := argo.NewCommandBuilder().MustParse([]string{"command"})

	renderOutputCheck(t, commandHelp001, com, argo.CommandHelpRenderer())
}

func TestCommandHelpRendererFail01(t *testing.T) {
	com := cli.Command().MustParse([]string{"command"})
	ren := argo.CommandHelpRenderer()

	for p := 1; p <= 101; p++ {
		wri := FailingWriter{FailAfter: p}
		buf := bufio.NewWriterSize(&wri, 1)

		err := ren.RenderHelp(com, buf)
		if err == nil {
			t.Error("expected err to not be nil but it was")
		}
	}
}
