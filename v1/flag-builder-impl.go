package argo

type flagBuilder struct {
	err error
	short byte
	long string
	desc string

	arg ArgumentBuilder
}

func (f *flagBuilder) Short(flag byte) FlagBuilder {
	if f.err == nil {
		if isValidShortFlag(flag) {
			f.short = flag
		} else {
			f.err = newInvalidShortFlagErr(flag, ErrHintShortInvalidChar)
		}
	}
	return f
}

func (f *flagBuilder) Long(flag string) FlagBuilder {
	if f.err == nil {
		if isValidLongFlag(flag) {
			f.long = flag
		} else {
			f.err = newInvalidLongFlagErr(flag, ErrHintLongInvalidChar)
		}
	}
	return f
}

func (f *flagBuilder) Description(desc string) FlagBuilder {
	f.desc = desc
	return f
}

func (f *flagBuilder) Arg(arg ArgumentBuilder) FlagBuilder {
	f.arg = arg
	return f
}

func (f *flagBuilder) build() (out Flag, err error) {
	if f.short == 0 && len(f.long) == 0 {
		return nil, newFlagBuildErr(ErrHintNoFlag)
	}

	var arg Argument

	if f.arg != nil {
		arg, err = f.arg.build()

		if err != nil {
			return nil, err
		}
	}

	return &flag{
		short: f.short,
		arg:   arg,
		long:  f.long,
		desc:  f.desc,
	}, nil
}

func (f *flagBuilder) mustBuild() Flag {
	panic("implement me")
}

func (f *flagBuilder) Bind(ptr interface{}, required bool) FlagBuilder {
	if f.arg != nil {
		f.arg.Bind(ptr)
		f.arg.Required(required)
	} else {
		f.arg = &argumentBuilder{required: required, binding: ptr}
	}
	return f
}

func (f *flagBuilder) Default(val interface{}) FlagBuilder {
	if f.arg != nil {
		f.arg.Default(val)
	} else {
		f.arg = &argumentBuilder{defVal: val}
	}
	return f
}
