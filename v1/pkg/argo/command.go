package argo

type Command interface {
	Description() string
	Arguments() []Argument
	FlagGroups() []FlagGroup
	UnmappedInput() []string
}
