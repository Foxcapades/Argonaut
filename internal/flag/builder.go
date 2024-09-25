package flag

import (
	"errors"

	"github.com/foxcapades/argonaut/pkg/argo"
)

type SpecBuilder struct {
	longForm    string
	hasLongForm bool

	shortForm    byte
	hasShortForm bool

	summary     string
	description string

	isRequired bool

	argument argo.ArgumentSpecBuilder

	lazyFns      []argo.FlagCallback
	immediateFns []argo.FlagCallback
}

func (s *SpecBuilder) WithLongForm(name string) argo.FlagSpecBuilder {
	s.longForm = name
	s.hasLongForm = true
	return s
}

func (s *SpecBuilder) WithShortForm(name byte) argo.FlagSpecBuilder {
	s.shortForm = name
	s.hasShortForm = true
	return s
}

func (s *SpecBuilder) WithArgument(arg argo.ArgumentSpecBuilder) argo.FlagSpecBuilder {
	s.argument = arg
	return s
}

func (s *SpecBuilder) WithLazyCallback(callback argo.FlagCallback) argo.FlagSpecBuilder {
	s.lazyFns = append(s.lazyFns, callback)
	return s
}

func (s *SpecBuilder) WithImmediateCallback(callback argo.FlagCallback) argo.FlagSpecBuilder {
	s.immediateFns = append(s.immediateFns, callback)
	return s
}

func (s *SpecBuilder) WithSummary(summary string) argo.FlagSpecBuilder {
	s.summary = summary
	return s
}

func (s *SpecBuilder) WithDescription(description string) argo.FlagSpecBuilder {
	s.description = description
	return s
}

func (s *SpecBuilder) Require() argo.FlagSpecBuilder {
	s.isRequired = true
	return s
}

func (s *SpecBuilder) Build(config argo.Config) (argo.FlagSpec, error) {
	if !s.hasLongForm && !s.hasShortForm {
		return nil, errors.New("flag configured with neither a long-form or short-form name")
	}

	spec := new(Spec)

	spec.

		// Flag must have:
		// * long and/or short form
		//

		// TODO implement me
		panic("implement me")
}
