package flag

import "github.com/foxcapades/argonaut/internal/arg"

type Base interface {
	IsRequired() bool

	HasArgument() bool

	Argument() arg.IArgument

	WasUsed() bool

	UsageCount() int
}
