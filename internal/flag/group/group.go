package group

import (
	"github.com/foxcapades/argonaut/pkg/argo"
)

type FlagGroup struct {
	name  string
	flags []argo.Flag
}

func (f FlagGroup) Name() string {
	return f.name
}

func (f FlagGroup) Flags() []argo.Flag {
	return f.flags
}

func (f FlagGroup) FindFlagByShortForm(name byte) argo.Flag {
	if name == 0 {
		return nil
	}

	for i := range f.flags {
		if f.flags[i].ShortForm() == name {
			return f.flags[i]
		}
	}

	return nil
}

func (f FlagGroup) FindFlagByLongForm(name string) argo.Flag {
	if len(name) == 0 {
		return nil
	}

	for i := range f.flags {
		if f.flags[i].LongForm() == name {
			return f.flags[i]
		}
	}

	return nil
}
