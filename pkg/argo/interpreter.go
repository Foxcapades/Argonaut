package argo

func CommandInterpreter(args []string, command Command) Interpreter {
	return &commandInterpreter{
		parser:   newParser(newEmitter(args)),
		command:  command,
		elements: newDeque[element](2),
	}
}

func CommandTreeInterpreter(args []string, command CommandTree) Interpreter {
	return &commandTreeInterpreter{
		parser:  newParser(newEmitter(args)),
		current: command,
		tree:    command,
		queue:   newDeque[element](2),
	}
}

type Interpreter interface {
	Run() error
}
