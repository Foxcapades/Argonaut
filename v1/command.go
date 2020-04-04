package argo

type Command interface {
	Description() string
	Arguments() []Argument
	Flags() []Flag
	UnmappedInput() []string

	appendUnmappedInput(string)
}

