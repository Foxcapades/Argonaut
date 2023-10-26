package argo

import (
	"reflect"

	"github.com/Foxcapades/Argonaut/internal/emit"
	"github.com/Foxcapades/Argonaut/internal/parse"
	"github.com/Foxcapades/Argonaut/internal/util"
)

func newCommandInterpreter(args []string, command Command) interpreter {
	return &commandInterpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		command:  command,
		elements: util.NewDeque[parse.Element](2),
		flagHits: make(map[util.Pair[byte, string]]Flag, 4),
	}
}

func newCommandTreeInterpreter(args []string, command CommandTree) interpreter {
	return &commandTreeInterpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		current:  command,
		tree:     command,
		queue:    util.NewDeque[parse.Element](2),
		flagHits: make(map[util.Pair[byte, string]]Flag, 4),
	}
}

type interpreter interface {
	Run() error
}

func hasBooleanArgument(flag Flag) bool {
	return flag.HasArgument() && flag.Argument().HasBinding() && flag.Argument().BindingType().Kind() == reflect.Bool
}
