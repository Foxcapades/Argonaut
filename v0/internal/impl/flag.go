package impl

import (
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"strings"
)

type Flag struct {
	arg  argo.Argument
	hits uint

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

func (f *Flag) HasArgument() bool {
	return f.arg != nil
}

func (f *Flag) Argument() argo.Argument {
	return f.arg
}

func (f *Flag) Description() string {
	return f.desc
}

func (f Flag) HasDescription() bool {
	return len(f.desc) > 0
}

func (f Flag) Hits() int {
	return int(f.hits)
}

func (f *Flag) IncrementHits() {
	f.hits++
}

func (f *Flag) String() (out string) {
	var bld strings.Builder

	if f.HasShort() {
		bld.WriteByte('-')
		bld.WriteByte(f.short)

		if f.HasArgument() {
			bld.WriteByte(' ')
			bld.WriteString(f.arg.String())
		}
	}

	if f.HasLong() {
		if bld.Len() > 0 {
			bld.WriteString(" | ")
		}

		bld.WriteString("--")
		bld.WriteString(f.long)

		if f.HasArgument() {
			bld.WriteByte('=')
			bld.WriteString(f.arg.String())
		}
	}

	return bld.String()
}
