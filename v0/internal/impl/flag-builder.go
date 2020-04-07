package impl

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewFlagBuilder() A.FlagBuilder {
	return new(FlagBuilder)
}

type FlagBuilder struct {
	err   error
	short byte
	long  string
	desc  string

	arg A.ArgumentBuilder

	shortSet bool
	longSet  bool

	onHit  A.FlagEventHandler
	parent A.FlagGroup
}

func (f *FlagBuilder) GetShort() byte            { return f.short }
func (f *FlagBuilder) HasShort() bool            { return f.shortSet }
func (f *FlagBuilder) GetLong() string           { return f.long }
func (f *FlagBuilder) HasLong() bool             { return f.longSet }
func (f *FlagBuilder) GetDescription() string    { return f.desc }
func (f *FlagBuilder) HasDescription() bool      { return len(f.desc) > 0 }
func (f *FlagBuilder) GetArg() A.ArgumentBuilder { return f.arg }
func (f *FlagBuilder) HasArg() bool              { return f.arg != nil }

func (f *FlagBuilder) Short(flag byte) A.FlagBuilder {
	f.shortSet = true
	f.short = flag
	return f
}

func (f *FlagBuilder) OnHit(fn A.FlagEventHandler) A.FlagBuilder {
	f.onHit = fn
	return f
}

func (f *FlagBuilder) Long(flag string) A.FlagBuilder {
	f.longSet = true
	f.long = flag
	return f
}

func (f *FlagBuilder) Description(desc string) A.FlagBuilder {
	f.desc = desc
	return f
}

func (f *FlagBuilder) Arg(arg A.ArgumentBuilder) A.FlagBuilder {
	f.arg = arg
	return f
}

func (f *FlagBuilder) Parent(fg A.FlagGroup) A.FlagBuilder {
	f.parent = fg
	return f
}

func (f *FlagBuilder) Build() (out A.Flag, err error) {
	if !(f.shortSet || f.longSet) {
		return nil, A.NewInvalidFlagError(A.InvalidFlagNoFlags)
	}

	if f.longSet && !util.IsValidLongFlag(f.long) {
		return nil, A.NewInvalidFlagError(A.InvalidFlagBadLongFlag)
	}

	if f.shortSet && !util.IsValidShortFlag(f.short) {
		return nil, A.NewInvalidFlagError(A.InvalidFlagBadShortFlag)
	}

	var arg A.Argument

	if f.arg != nil {
		arg, err = f.arg.Build()
		if err != nil {
			return nil, err
		}
	}

	return &Flag{
		arg:    arg,
		long:   f.long,
		desc:   f.desc,
		short:  f.short,
		onHit:  f.onHit,
		parent: f.parent,
	}, nil
}

func (f *FlagBuilder) MustBuild() A.Flag {
	if f, e := f.Build(); e != nil {
		panic(e)
	} else {
		return f
	}
}

func (f *FlagBuilder) Bind(ptr interface{}, required bool) A.FlagBuilder {
	if f.arg == nil {
		f.arg = GetProvider().NewArg()
	}

	f.arg.Bind(ptr).Required(required)

	return f
}

func (f *FlagBuilder) Default(val interface{}) A.FlagBuilder {
	if f.arg == nil {
		f.arg = GetProvider().NewArg()
	}
	f.arg.Default(val)
	return f
}
