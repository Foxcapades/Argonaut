package argo

type bindValidator func(argumentBuilder) error

type argumentBuilder struct {
	bindValidators []bindValidator

	required bool
	defVal   interface{}
	binding  interface{}
	hintTxt  string
}

func (a *argumentBuilder) Hint(string) ArgumentBuilder {
	panic("implement me")
}

func (a *argumentBuilder) Default(interface{}) ArgumentBuilder {
	panic("implement me")
}

func (a *argumentBuilder) Bind(ptr interface{}) ArgumentBuilder {
	panic("implement me")
}

func (a *argumentBuilder) Description(string) ArgumentBuilder {
	panic("implement me")
}

func (a *argumentBuilder) Require() ArgumentBuilder {
	panic("implement me")
}

func (a *argumentBuilder) Required(bool) ArgumentBuilder {
	panic("implement me")
}

func (a *argumentBuilder) build() (Argument, error) {
	panic("implement me")
}

func (a *argumentBuilder) mustBuild() Argument {
	panic("implement me")
}
