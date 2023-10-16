package argo

type Command interface {
	Name() string

	Description() string
	HasDescription() bool

	FlagGroups() []FlagGroup
	HasFlagGroups() bool

	FindShortFlag(c byte) Flag
	FindLongFlag(name string) Flag

	TryFlag(ref FlagRef) (bool, error)

	Arguments() []Argument
	HasArguments() bool
	AppendArgument(rawArgument string) error

	UnmappedInputs() []string
	HasUnmappedInputs() bool
	AppendUnmapped(val string)

	PassthroughInputs() []string
	HasPassthroughInputs() bool
	AppendPassthrough(val string)

	GetUnmappedLabel() string

	HasUnmappedLabel() bool
}
