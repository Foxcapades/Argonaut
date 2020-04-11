package flag

import (
	"errors"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewFlagGroupBuilder(A.Provider) A.FlagGroupBuilder {
	return new(GBuilder)
}

type GBuilder struct {
	ParentNode  A.Command
	NameTxt     trait.Named
	DescTxt     trait.Described
	FlagNodes   []A.FlagBuilder
	WarningVals []string
}

//
// Getters
//

func (f *GBuilder) GetName() string {
	return f.NameTxt.NameTxt
}

func (f *GBuilder) GetDescription() string {
	return f.DescTxt.DescTxt
}

func (f *GBuilder) GetFlags() []A.FlagBuilder {
	return f.FlagNodes
}

//
// Setters
//

func (f *GBuilder) Parent(com A.Command) A.FlagGroupBuilder {
	f.ParentNode = com
	return f
}

func (f *GBuilder) Name(name string) A.FlagGroupBuilder {
	f.NameTxt.NameTxt = name
	return f
}

func (f *GBuilder) Description(desc string) A.FlagGroupBuilder {
	f.DescTxt.DescTxt = desc
	return f
}

//
// Operations
//

func (f *GBuilder) Flag(flag A.FlagBuilder) A.FlagGroupBuilder {
	if flag == nil {
		f.WarningVals = append(f.WarningVals, "FlagGroupBuilder: nil value passed to Flag()")
	} else {
		f.FlagNodes = append(f.FlagNodes, flag)
	}
	return f
}

func (f *GBuilder) Build() (out A.FlagGroup, err error) {
	if len(f.FlagNodes) == 0 {
		// TODO: make this a real error
		return nil, errors.New("no flags in group")
	}

	flags := make([]A.Flag, len(f.FlagNodes))

	out = &Group{ParentElement: f.ParentNode, Described: f.DescTxt, Named: f.NameTxt, FlagNodes: flags}

	for i, fb := range f.FlagNodes {
		fb.Parent(out)
		if flag, err := fb.Build(); err != nil {
			return nil, err
		} else {
			flags[i] = flag
		}
	}

	return
}

func (f *GBuilder) MustBuild() A.FlagGroup {
	if out, err := f.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
