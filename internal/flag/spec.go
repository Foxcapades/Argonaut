package flag

import (
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Spec struct {
	shortForm   byte
	isRequired  bool
	hasArg      bool
	usageCount  uint32
	longForm    string
	description string
	argument    argo.ArgumentSpec
}

func (s Spec) LongForm() string {
	return s.longForm
}

func (s Spec) HasLongForm() bool {
	return len(s.longForm) > 0
}

func (s Spec) ShortForm() byte {
	return s.shortForm
}

func (s Spec) HasShortForm() bool {
	return s.shortForm != 0
}

func (s Spec) Description() string {
	return s.description
}

func (s Spec) HasDescription() bool {
	return len(s.description) > 0
}

func (s Spec) IsRequired() bool {
	return s.isRequired
}

func (s Spec) WasUsed() bool {
	return s.usageCount > 0
}

func (s Spec) UsageCount() uint32 {
	return s.usageCount
}

func (s *Spec) MarkUsed() {
	s.usageCount++
}

func (s Spec) HasExplicitArgument() bool {
	return s.hasArg
}

func (s Spec) Argument() argo.ArgumentSpec {
	return s.argument
}

func (s Spec) ToFlag() argo.Flag {
	return &Flag{
		longForm:   s.longForm,
		shortForm:  s.shortForm,
		isRequired: s.isRequired,
		hasArg:     s.hasArg,
		usages:     s.usageCount,
		argument:   s.argument.ToArgument(),
	}
}
