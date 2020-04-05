package impl

import (
	"github.com/Foxcapades/Argonaut/v1/internal/util"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func NewFlagBuilder() argo.FlagBuilder {
	return new(FlagBuilder)
}

type FlagBuilder struct {
	err   error
	short byte
	long  string
	desc  string

	arg argo.ArgumentBuilder

	shortSet bool
	longSet  bool
}

func (f *FlagBuilder) Short(flag byte) argo.FlagBuilder {
	f.shortSet = true
	f.short = flag
	return f
}

func (f *FlagBuilder) Long(flag string) argo.FlagBuilder {
	f.longSet = true
	f.long = flag
	return f
}

func (f *FlagBuilder) Description(desc string) argo.FlagBuilder {
	f.desc = desc
	return f
}

func (f *FlagBuilder) Arg(arg argo.ArgumentBuilder) argo.FlagBuilder {
	f.arg = arg
	return f
}

func (f *FlagBuilder) Build() (out argo.Flag, err error) {
	if !(f.shortSet || f.longSet) {
		return nil, argo.NewInvalidFlagError(argo.InvalidFlagNoFlags)
	}

	if f.longSet && !util.IsValidLongFlag(f.long) {
		return nil, argo.NewInvalidFlagError(argo.InvalidFlagBadLongFlagCharacter)
	}

	if f.shortSet && !util.IsValidShortFlag(f.short) {
		return nil, argo.NewInvalidFlagError(argo.InvalidFlagBadShortFlagCharacter)
	}

	var arg argo.Argument
	var cnt argo.UseCounter

	cArg := &Argument{bind: &cnt}

	if f.arg != nil {
		arg, err = f.arg.Build()

		if err != nil {
			return nil, err
		}
	} else {
		arg = cArg
	}

	return &Flag{
		short: f.short,
		arg:   arg,
		cArg:  cArg,
		hits:  &cnt,
		long:  f.long,
		desc:  f.desc,
		isReq: arg.Required(),
	}, nil
}

func (f *FlagBuilder) MustBuild() argo.Flag {
	if f, e := f.Build(); e != nil {
		panic(e)
	} else {
		return f
	}
}

func (f *FlagBuilder) Bind(ptr interface{}, required bool) argo.FlagBuilder {
	if f.arg != nil {
		f.arg.Bind(ptr)
		f.arg.Required(required)
	} else {
		f.arg = &ArgumentBuilder{required: required, binding: ptr}
	}
	return f
}

func (f *FlagBuilder) Default(val interface{}) argo.FlagBuilder {
	if f.arg != nil {
		f.arg.Default(val)
	} else {
		f.arg = &ArgumentBuilder{defVal: val}
	}
	return f
}
