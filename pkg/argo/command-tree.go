package argo

type CommandTree interface {
	CommandNode
	ParentCommand
}
