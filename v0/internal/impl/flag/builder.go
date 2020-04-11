package flag

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewBuilder(provider A.Provider) A.FlagBuilder {
	return &Builder{Provider: provider}
}

type Builder struct {
	Provider A.Provider

	Error error

	ShortFlag  byte
	IsShortSet bool

	LongFlag  string
	IsLongSet bool

	DescriptionText trait.Described

	ArgBuilder A.ArgumentBuilder

	UseCountBinding *int
	OnHitCallback   A.FlagEventHandler

	ParentElement A.FlagGroup
}

func (f *Builder) Build() (out A.Flag, err error) {
	if !(f.IsShortSet || f.IsLongSet) {
		return nil, A.NewInvalidFlagError(A.InvalidFlagNoFlags)
	}

	if f.IsLongSet && !util.IsValidLongFlag(f.LongFlag) {
		return nil, A.NewInvalidFlagError(A.InvalidFlagBadLongFlag)
	}

	if f.IsShortSet && !util.IsValidShortFlag(f.ShortFlag) {
		return nil, A.NewInvalidFlagError(A.InvalidFlagBadShortFlag)
	}

	var arg A.Argument
	tmp := new(Flag)

	if f.ArgBuilder != nil {
		f.ArgBuilder.Parent(tmp)
		arg, err = f.ArgBuilder.Build()
		if err != nil {
			return nil, err
		}
	}

	tmp.ArgumentElement = arg
	tmp.LongForm = f.LongFlag
	tmp.Described = f.DescriptionText
	tmp.ShortForm = f.ShortFlag
	tmp.OnHitCallback = f.OnHitCallback
	tmp.ParentElement = f.ParentElement
	tmp.HitCountBinding = f.UseCountBinding
	return tmp, nil
}

func (f *Builder) MustBuild() A.Flag {
	if f, e := f.Build(); e != nil {
		panic(e)
	} else {
		return f
	}
}
