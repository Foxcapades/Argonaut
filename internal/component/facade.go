package component

import (
	"github.com/foxcapades/argonaut/internal/argument"
	"github.com/foxcapades/argonaut/internal/flag"
)

type Facade interface {
	HasCallback() bool

	HasSubCommands() bool

	IsSubCommand(value string) bool

	Shift(subCommand string)

	LongFlag(name string) (flag.Flag, bool)

	ShortFlag(name byte) (flag.Flag, bool)

	PositionalArguments() []argument.Spec

	AppendUnknown(value string)

	AppendPassthrough(value string)
}
