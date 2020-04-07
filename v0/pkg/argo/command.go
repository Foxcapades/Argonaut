package argo

type Command interface {
	Name() string
	Description() string
	HasDescription() bool
	Arguments() []Argument
	FlagGroups() []FlagGroup
	UnmappedInput() []string
	Unmarshaler() ValueUnmarshaler
}
