package argo

type SubCommand interface {
	CommandNode
	Parent() ParentCommand
}
