package flag

import "github.com/foxcapades/argonaut/internal/argument"

type Base interface {
	IsRequired() bool

	HasArgument() bool

	Argument() argument.IArgument

	WasUsed() bool

	UsageCount() int
}
