package flag

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type group struct {
	name     string
	desc     string
	flags    []argo.Flag
	warnings *argo.WarningContext
}

func (f group) Name() string {
	return f.name
}

func (f group) Description() string {
	return f.desc
}

func (f group) HasDescription() bool {
	return len(f.desc) > 0
}

func (f group) Flags() []argo.Flag {
	return f.flags
}

func (f group) FindShortFlag(c byte) argo.Flag {
	for _, flag := range f.flags {
		if flag.HasShortForm() && flag.ShortForm() == c {
			return flag
		}
	}

	return nil
}

func (f group) FindLongFlag(name string) argo.Flag {
	for _, flag := range f.flags {
		if flag.HasLongForm() && flag.LongForm() == name {
			return flag
		}
	}

	return nil
}

func (f group) Size() int {
	return len(f.flags)
}
