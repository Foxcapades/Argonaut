package impl

import (
	"github.com/Foxcapades/Argonaut/v1/internal/util"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func NewFlagBuilder() *FlagBuilder {
	return new(FlagBuilder)
}

type FlagBuilder struct {
	err   error
	short byte
	long  string
	desc  string

	arg argo.ArgumentBuilder
}

func (f *FlagBuilder) Short(flag byte) argo.FlagBuilder {
	if f.err == nil {
		if util.IsValidShortFlag(flag) {
			f.short = flag
		} else {
			f.err = &argo.InvalidFlagCharError{Flag: string([]byte{flag}), Hint: argo.ErrHintShortInvalidChar}
		}
	}
	return f
}

func (f *FlagBuilder) Long(flag string) argo.FlagBuilder {
	if f.err == nil {
		if util.IsValidLongFlag(flag) {
			f.long = flag
		} else {
			f.err = &argo.InvalidFlagCharError{Flag: flag, Hint: argo.ErrHintLongInvalidChar}
		}
	}
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
	if f.short == 0 && len(f.long) == 0 {
		return nil, &argo.InvalidFlagError{Hint: argo.ErrHintNoFlag}
	}

	var arg argo.Argument

	if f.arg != nil {
		arg, err = f.arg.Build()

		if err != nil {
			return nil, err
		}
	}

	return &Flag{
		short: f.short,
		arg:   arg,
		long:  f.long,
		desc:  f.desc,
	}, nil
}

func (f *FlagBuilder) MustBuild() argo.Flag {
	panic("implement me")
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
