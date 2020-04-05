package impl

import (
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func NewFlagGroupBuilder() argo.FlagGroupBuilder {
	return new(FlagGroupBuilder)
}

type FlagGroupBuilder struct {
	name  string
	desc  string
	flags []argo.FlagBuilder
}

func (f *FlagGroupBuilder) Name(name string) (this argo.FlagGroupBuilder) {
	f.name = name
	return f
}

func (f *FlagGroupBuilder) Description(desc string) (this argo.FlagGroupBuilder) {
	f.desc = desc
	return f
}

func (f *FlagGroupBuilder) Flag(flag argo.FlagBuilder) (this argo.FlagGroupBuilder) {
	f.flags = append(f.flags, flag)
	return f
}

func (f *FlagGroupBuilder) Build() (argo.FlagGroup, error) {
	flags := make([]argo.Flag, len(f.flags))
	for i, fb := range f.flags {
		if flag, err := fb.Build(); err != nil {
			return nil, err
		} else {
			flags[i] = flag
		}
	}
	return &FlagGroup{
		description: f.desc,
		name:        f.name,
		flags:       flags,
	}, nil
}

func (f *FlagGroupBuilder) MustBuild() argo.FlagGroup {
	if out, err := f.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
