package argo

type ParentCommand interface {
	CommandNode

	SubCommandContainer
}
