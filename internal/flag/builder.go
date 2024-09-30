package flag

import (
	"strings"

	"github.com/foxcapades/argonaut/internal/argument"
	"github.com/foxcapades/argonaut/internal/util/xerr"
	"github.com/foxcapades/argonaut/pkg/argo"
)

func NewBuilder() *SpecBuilder {
	return new(SpecBuilder)
}

type SpecBuilder struct {
	longForm    string
	hasLongForm bool

	shortForm    byte
	hasShortForm bool

	description string

	isRequired bool

	hasArg   bool
	argument argo.ArgumentSpecBuilder

	lazyFns      []argo.FlagCallback
	immediateFns []argo.FlagCallback
}

func (s *SpecBuilder) WithLongForm(name string) argo.FlagSpecBuilder {
	s.longForm = strings.TrimSpace(name)
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
	s.hasArg = true
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

func (s *SpecBuilder) WithDescription(description string) argo.FlagSpecBuilder {
	s.description = description
	return s
}

func (s *SpecBuilder) Require() argo.FlagSpecBuilder {
	s.isRequired = true
	return s
}

func (s *SpecBuilder) Build(config argo.Config) (argo.FlagSpec, error) {
	errs := xerr.NewMultiError()

	if !s.hasLongForm && !s.hasShortForm {
		errs.AppendMsg(argo.ErrMsgFlagHasNoNames)
	}

	if s.hasLongForm && !config.Flags.LongFormValidator(s.longForm, config) {
		errs.AppendMsg(argo.ErrMsgInvalidLongFlagName(s.longForm, config))
	}

	if s.hasShortForm && !config.Flags.ShortFormValidator(s.shortForm, config) {
		errs.AppendMsg(argo.ErrMsgInvalidShortFlagName(s.shortForm, config))
	}

	var arg argo.ArgumentSpec

	if s.hasArg {
		if a, e := s.argument.Build(config); e == nil {
			arg = a
		} else {
			errs.Append(e)
		}
	} else {
		arg = new(argument.FallbackArgumentSpec)
	}

	if errs.IsEmpty() {
		return &Spec{
			shortForm:   s.shortForm,
			isRequired:  s.isRequired,
			hasArg:      s.hasArg,
			longForm:    s.longForm,
			description: s.description,
			argument:    arg,
		}, nil
	}

	return nil, errs
}
