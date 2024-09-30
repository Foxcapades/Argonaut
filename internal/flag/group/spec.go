package group

import (
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Spec struct {
	name        string
	description string
	flags       []argo.FlagSpec
}

func (s Spec) Name() string {
	return s.name
}

func (s Spec) Description() string {
	return s.description
}

func (s Spec) HasDescription() bool {
	return len(s.description) > 0
}

func (s Spec) Flags() []argo.FlagSpec {
	return s.flags
}

func (s Spec) FindFlagByShortForm(name byte) argo.FlagSpec {
	if name == 0 {
		return nil
	}

	for i := range s.flags {
		if s.flags[i].ShortForm() == name {
			return s.flags[i]
		}
	}

	return nil
}

func (s Spec) FindFlagByLongForm(name string) argo.FlagSpec {
	if len(name) == 0 {
		return nil
	}

	for i := range s.flags {
		if s.flags[i].LongForm() == name {
			return s.flags[i]
		}
	}

	return nil
}

func (s Spec) ToFlagGroup() argo.FlagGroup {
	out := new(FlagGroup)
	out.name = s.name
	out.flags = make([]argo.Flag, len(s.flags))

	for i := range s.flags {
		out.flags[i] = s.flags[i].ToFlag()
	}

	return out
}
