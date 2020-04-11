package flag

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"strings"
)

type Flag struct {
	trait.Described

	ParentElement   A.FlagGroup
	ArgumentElement A.Argument

	HitCount   uint
	LongForm   string
	ShortForm  byte
	IsRequired bool

	HitCountBinding *int
	OnHitCallback   A.FlagEventHandler
}

func (f *Flag) Short() byte          { return f.ShortForm }
func (f *Flag) HasShort() bool       { return f.ShortForm > 0 }
func (f *Flag) Long() string         { return f.LongForm }
func (f *Flag) HasLong() bool        { return len(f.LongForm) > 0 }
func (f *Flag) Required() bool       { return f.IsRequired }
func (f *Flag) HasArgument() bool    { return f.ArgumentElement != nil }
func (f *Flag) Argument() A.Argument { return f.ArgumentElement }
func (f *Flag) Hits() int            { return int(f.HitCount) }
func (f *Flag) IncrementHits() {
	f.HitCount++

	if f.HitCountBinding != nil {
		*f.HitCountBinding++
	}

	if f.OnHitCallback != nil {
		f.OnHitCallback(f)
	}
}
func (f *Flag) Parent() A.FlagGroup { return f.ParentElement }

func (f *Flag) String() (out string) {
	var bld strings.Builder

	if f.HasShort() {
		bld.WriteByte('-')
		bld.WriteByte(f.ShortForm)

		if f.HasArgument() {
			bld.WriteByte(' ')
			bld.WriteString(f.ArgumentElement.String())
		}
	}

	if f.HasLong() {
		if bld.Len() > 0 {
			bld.WriteString(" | ")
		}

		bld.WriteString("--")
		bld.WriteString(f.LongForm)

		if f.HasArgument() {
			bld.WriteByte('=')
			bld.WriteString(f.ArgumentElement.String())
		}
	}

	return bld.String()
}
