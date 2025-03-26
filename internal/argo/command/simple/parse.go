package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/emit"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Parse(command argo.Command, args []string) (argo.ParseResult, error) {
	i := Interpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		command:  command,
		elements: utils.NewDeque[parse.Element](2),
		flagHits: flag.NewQueue(),
	}

	return i.Run()
}
