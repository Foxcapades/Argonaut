package argo

type PassthroughContainerSpec interface {
	PassthroughContainer

	AppendPassthrough(passthrough string)
}
