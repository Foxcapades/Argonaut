package flag

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewBuilder(provider A.Provider) A.FlagBuilder {
	return &builder{provider: provider}
}

type iFb = A.FlagBuilder

type builder struct {
	provider A.Provider

	err   error
	short byte
	long  string
	desc  string

	arg A.ArgumentBuilder

	shortSet bool
	longSet  bool

	hitBind *int

	onHit  A.FlagEventHandler
	parent A.FlagGroup
}

func (f *builder) GetShort() byte            { return f.short }
func (f *builder) HasShort() bool            { return f.shortSet }
func (f *builder) GetLong() string           { return f.long }
func (f *builder) HasLong() bool             { return f.longSet }
func (f *builder) GetDescription() string    { return f.desc }
func (f *builder) HasDescription() bool      { return len(f.desc) > 0 }
func (f *builder) GetArg() A.ArgumentBuilder { return f.arg }
func (f *builder) HasArg() bool              { return f.arg != nil }

func (f *builder) Short(flag byte) iFb             { f.shortSet = true; f.short = flag; return f }
func (f *builder) OnHit(fn A.FlagEventHandler) iFb { f.onHit = fn; return f }
func (f *builder) Long(flag string) iFb            { f.longSet = true; f.long = flag; return f }
func (f *builder) Description(desc string) iFb     { f.desc = desc; return f }
func (f *builder) Arg(arg A.ArgumentBuilder) iFb   { f.arg = arg; return f }
func (f *builder) Parent(fg A.FlagGroup) iFb       { f.parent = fg; return f }
func (f *builder) BindUseCount(ptr *int) iFb       { f.hitBind = ptr; return f }

func (f *builder) Build() (out A.Flag, err error) {
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
	tmp := new(flag)

	if f.arg != nil {
		f.arg.Parent(tmp)
		arg, err = f.arg.Build()
		if err != nil {
			return nil, err
		}
	}

	tmp.arg = arg
	tmp.long = f.long
	tmp.desc = f.desc
	tmp.short = f.short
	tmp.onHit = f.onHit
	tmp.parent = f.parent
	tmp.hitBinding = f.hitBind
	return tmp, nil
}

func (f *builder) MustBuild() A.Flag {
	if f, e := f.Build(); e != nil {
		panic(e)
	} else {
		return f
	}
}

func (f *builder) Bind(ptr interface{}, required bool) iFb {
	if f.arg == nil {
		f.arg = f.provider.NewArg()
	}

	f.arg.Bind(ptr).Required(required)

	return f
}

func (f *builder) Default(val interface{}) iFb {
	if f.arg == nil {
		f.arg = f.provider.NewArg()
	}
	f.arg.Default(val)
	return f
}
