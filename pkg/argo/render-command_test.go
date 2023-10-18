package argo_test

import (
	"testing"

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
