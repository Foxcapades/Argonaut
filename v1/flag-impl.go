package argo

type flagProps = uint8

const (
	flagHasShort flagProps = 1 << iota
	flagHasLong
	flagIsReq
)

type flag struct {
	props props
	short byte
	arg   Argument
	long  string
	desc  string
}

func (f *flag) Hits() int {
	panic("implement me")
}

func (f *flag) hit() {
	panic("implement me")
}

func (f *flag) Short() byte {
	return f.short
}

func (f *flag) HasShort() bool {
	return flagHasShort == flagHasShort & f.props
}

func (f *flag) Long() string {
	return f.long
}

func (f *flag) HasLong() bool {
	return flagHasLong == flagHasLong & f.props
}

func (f *flag) Required() bool {
	return flagIsReq == flagIsReq & f.props
}

func (f *flag) Argument() Argument {
	return f.arg
}

func (f *flag) HasArgument() bool {
	return f.arg != nil
}

func (f *flag) IsArgumentRequired() bool {
	//return f.arg != nil && f.arg
	return false
}

func (f *flag) Description() string {
	return f.desc
}

func (f *flag) HasDescription() bool {
	return len(f.desc) > 0
}

