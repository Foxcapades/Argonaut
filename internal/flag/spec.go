package flag

import (
	"github.com/foxcapades/argonaut/internal/argument"
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Spec struct {
	shortForm   byte
	isRequired  bool
	usageCount  uint32
	longForm    string
	summary     string
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

func (s Spec) Summary() string {
	return s.summary
}

func (s Spec) Description() string {
	return s.description
}

func (s Spec) HasHelpText() bool {
	return len(s.summary) > 0
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
	_, ok := s.argument.(*argument.FallbackArgument)
	return !ok
}

func (s Spec) Argument() argo.ArgumentSpec {
	return s.argument
}
