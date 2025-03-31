package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/emit"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Parse(command argo.Command, args []string) ([]argo.InputWarning, error) {
	i := Interpreter{
		parser:    parse.NewParser(emit.NewEmitter(args)),
		command:   command,
		warnings:  make([]argo.InputWarning, 0, 2),
		elements:  utils.NewDeque[parse.Element](2),
		flagHits:  flag.NewQueue(),
		flagIndex: flag.BuildFlagIndex(command),
	}

	return i.Run()
}
