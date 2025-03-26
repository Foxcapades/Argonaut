package flag

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Group struct {
	name  string
	desc  string
	flags []argo.Flag
}

func (f Group) Name() string {
	return f.name
}

func (f Group) Description() string {
	return f.desc
}

func (f Group) HasDescription() bool {
	return len(f.desc) > 0
}

func (f Group) Flags() []argo.Flag {
	return f.flags
}

func (f Group) FindShortFlag(c byte) argo.Flag {
	for _, flag := range f.flags {
		if flag.HasShortForm() && flag.ShortForm() == c {
			return flag
		}
	}

	return nil
}

func (f Group) FindLongFlag(name string) argo.Flag {
	for _, flag := range f.flags {
		if flag.HasLongForm() && flag.LongForm() == name {
			return flag
		}
	}

	return nil
}

func (f Group) Size() int {
	return len(f.flags)
}
