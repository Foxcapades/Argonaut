package argo

type Parser interface {
	Parse(args []string) error
}
