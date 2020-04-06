package argo

type Parser interface {
	Parse(args []string, command Command) error
	Unrecognized() []string
	Passthroughs() []string
}
