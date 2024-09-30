package flag

import (
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Flag struct {
	longForm   string
	shortForm  byte
	isRequired bool
	hasArg     bool
	usages     uint32
	argument   argo.Argument
}

func (f Flag) LongForm() string {
	return f.longForm
}

func (f Flag) HasLongForm() bool {
	return len(f.longForm) > 0
}

func (f Flag) ShortForm() byte {
	return f.shortForm
}

func (f Flag) HasShortForm() bool {
	return f.shortForm != 0
}

func (f Flag) IsRequired() bool {
	return f.isRequired
}

func (f Flag) WasUsed() bool {
	return f.usages > 0
}

func (f Flag) UsageCount() uint32 {
	return f.usages
}

func (f Flag) HasExplicitArgument() bool {
	return f.hasArg
}

func (f Flag) Argument() argo.Argument {
	return f.argument
}
