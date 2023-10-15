package interpret

import (
	"github.com/Foxcapades/Argonaut/v1/internal/event"
	"github.com/Foxcapades/Argonaut/v1/internal/parse"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func CommandInterpreter(args []string, command argo.Command) Interpreter {
	return &commandInterpreter{parse.NewParser(event.NewEmitter(args)), command}
}

func CommandTreeInterpreter(args []string, command argo.CommandTree) Interpreter {
	return &commandTreeInterpreter{parse.NewParser(event.NewEmitter(args)), command, false}
}

type Interpreter interface {
	Run() error
}
