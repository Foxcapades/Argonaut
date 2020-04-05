package impl

import "github.com/Foxcapades/Argonaut/v1/pkg/argo"

type Flag struct {
	arg  argo.Argument
	cArg argo.Argument
	hits *argo.UseCounter

	long string
	desc string

	short byte
	isReq bool
}

func (f *Flag) Short() byte {
	return f.short
}

func (f *Flag) HasShort() bool {
	return f.short > 0
}

func (f *Flag) Long() string {
	return f.long
}

func (f *Flag) HasLong() bool {
	return len(f.long) > 0
}

func (f *Flag) Required() bool {
	return f.isReq
}

func (f *Flag) Argument() argo.Argument {
	if f.arg == nil {
		return f.cArg
	}
	return f.arg
}

func (f *Flag) Description() string {
	return f.desc
}

func (f Flag) HasDescription() bool {
	return len(f.desc) > 0
}

func (f Flag) Hits() int {
	return int(*f.hits)
}
