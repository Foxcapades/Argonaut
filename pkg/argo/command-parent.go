package argo

type CommandParent interface {
	CommandGroups() []CommandGroup

	HasCustomCommandGroups() bool

	FindChild(name string) CommandNode
}
