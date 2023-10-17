package argo

type CommandParent interface {
	CommandGroups() []CommandGroup

	FindChild(name string) CommandNode

	onIncomplete()
}
