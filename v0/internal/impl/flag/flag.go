package flag

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"strings"
)

type flag struct {
	parent A.FlagGroup
	arg    A.Argument
	hits   uint
	long   string
	desc   string
	short  byte
	isReq  bool

	hitBinding *int
	onHit A.FlagEventHandler
}

func (f *flag) Short() byte          { return f.short }
func (f *flag) HasShort() bool       { return f.short > 0 }
func (f *flag) Long() string         { return f.long }
func (f *flag) HasLong() bool        { return len(f.long) > 0 }
func (f *flag) Required() bool       { return f.isReq }
func (f *flag) HasArgument() bool    { return f.arg != nil }
func (f *flag) Argument() A.Argument { return f.arg }
func (f *flag) Description() string  { return f.desc }
func (f *flag) HasDescription() bool { return len(f.desc) > 0 }
func (f *flag) Hits() int            { return int(f.hits) }
func (f *flag) IncrementHits() {
	f.hits++

	if f.hitBinding != nil {
		*f.hitBinding++
	}

	if f.onHit != nil {
		f.onHit(f)
	}
}
func (f *flag) Parent() A.FlagGroup { return f.parent }

func (f *flag) String() (out string) {
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
