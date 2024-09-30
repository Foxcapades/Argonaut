package argo

type SubCommandSpec interface {
	CommandNodeSpec
	Parent() ParentCommandSpec
}
