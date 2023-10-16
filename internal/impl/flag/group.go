package flag

import (
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

// implements argo.FlagGroup
type group struct {
	name  string
	desc  string
	flags []argo.Flag
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

func (f group) TryFlag(ref argo.FlagRef) (bool, error) {
	for _, flag := range f.flags {
		if ref.HasLongForm() && flag.HasLongForm() && ref.LongForm() == flag.LongForm() {
			if ref.HasValue() {
				return true, flag.HitWithArg(ref.Value())
			} else {
				flag.Hit()
				return true, nil
			}
		}
		if ref.HasShortForm() && flag.HasShortForm() && ref.ShortForm() == flag.ShortForm() {
			if ref.HasValue() {
				return true, flag.HitWithArg(ref.Value())
			} else {
				flag.Hit()
				return true, nil
			}
		}
	}

	return false, nil
}
