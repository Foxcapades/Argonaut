package argo

type CommandCallback = func(command Command)

func SimpleCommandCallback(callback func()) CommandCallback {
	return func(Command) { callback() }
}
