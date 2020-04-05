package impl

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewFlagGroupBuilder() A.FlagGroupBuilder {
	return new(FlagGroupBuilder)
}

type FlagGroupBuilder struct {
	name     string
	desc     string
	flags    []A.FlagBuilder
	warnings []string
}

func (f *FlagGroupBuilder) Name(name string) (this A.FlagGroupBuilder) {
	f.name = name
	return f
}

func (f *FlagGroupBuilder) GetName() string {
	return f.name
}

func (f *FlagGroupBuilder) Description(desc string) (this A.FlagGroupBuilder) {
	f.desc = desc
	return f
}

func (f *FlagGroupBuilder) GetDescription() string {
	return f.desc
}

func (f *FlagGroupBuilder) Flag(flag A.FlagBuilder) (this A.FlagGroupBuilder) {
	if flag == nil {
		f.warnings = append(f.warnings, "FlagGroupBuilder: nil value passed to Flag()")
	}
	f.flags = append(f.flags, flag)
	return f
}

func (f *FlagGroupBuilder) GetFlags() []A.FlagBuilder {
	return f.flags
}

func (f *FlagGroupBuilder) Build() (A.FlagGroup, error) {
	flags := make([]A.Flag, len(f.flags))
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

func (f *FlagGroupBuilder) MustBuild() A.FlagGroup {
	if out, err := f.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
