package argo

import "reflect"

func newCommandInterpreter(args []string, command Command) interpreter {
	return &commandInterpreter{
		parser:   newParser(newEmitter(args)),
		command:  command,
		elements: newDeque[element](2),
	}
}

func newCommandTreeInterpreter(args []string, command CommandTree) interpreter {
	return &commandTreeInterpreter{
		parser:  newParser(newEmitter(args)),
		current: command,
		tree:    command,
		queue:   newDeque[element](2),
	}
}

type interpreter interface {
	Run() error
}

func hasBooleanArgument(flag Flag) bool {
	return flag.HasArgument() && flag.Argument().HasBinding() && flag.Argument().BindingType().Kind() == reflect.Bool
}
