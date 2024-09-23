package argo

type ArgumentSpec interface {
	Argument

	PreValidate(rawInput string) error

	ProcessInput(rawInput string) error
}
