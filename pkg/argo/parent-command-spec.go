package argo

type ParentCommandSpec interface {
	CommandNodeSpec

	SubCommandSpecContainer
}
