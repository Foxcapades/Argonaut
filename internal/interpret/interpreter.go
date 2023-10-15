package interpret

import (
	"github.com/Foxcapades/Argonaut/internal/event"
	"github.com/Foxcapades/Argonaut/internal/parse"
	"github.com/Foxcapades/Argonaut/pkg/argo"
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
