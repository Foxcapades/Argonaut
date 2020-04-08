package flag

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewFlagGroupBuilder(A.Provider) A.FlagGroupBuilder {
	return new(gBuilder)
}

type iFgb = A.FlagGroupBuilder

type gBuilder struct {
	parent   A.Command
	name     string
	desc     string
	flags    []A.FlagBuilder
	warnings []string
}

//
// Getters
//

func (f *gBuilder) GetName() string           { return f.name }
func (f *gBuilder) GetDescription() string    { return f.desc }
func (f *gBuilder) GetFlags() []A.FlagBuilder { return f.flags }

//
// Setters
//

func (f *gBuilder) Parent(com A.Command) iFgb    { f.parent = com; return f }
func (f *gBuilder) Name(name string) iFgb        { f.name = name; return f }
func (f *gBuilder) Description(desc string) iFgb { f.desc = desc; return f }

//
// Operations
//

func (f *gBuilder) Flag(flag A.FlagBuilder) iFgb {
	if flag == nil {
		f.warnings = append(f.warnings, "FlagGroupBuilder: nil value passed to Flag()")
	} else {
		f.flags = append(f.flags, flag)
	}
	return f
}

func (f *gBuilder) Build() (out A.FlagGroup, err error) {
	flags := make([]A.Flag, len(f.flags))

	out = &FlagGroup{parent: f.parent, desc: f.desc, name: f.name, flags: flags}

	for i, fb := range f.flags {
		fb.Parent(out)
		if flag, err := fb.Build(); err != nil {
			return nil, err
		} else {
			flags[i] = flag
		}
	}

	return
}

func (f *gBuilder) MustBuild() A.FlagGroup {
	if out, err := f.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
