package flag

import "github.com/Foxcapades/Argonaut/pkg/argo"

func NewLongFlagRef(name, value string, hasValue bool) argo.FlagRef {
	return ref{long: name, value: value, hasValue: hasValue}
}

func NewShortFlagRef(c byte, value string, hasValue bool) argo.FlagRef {
	return ref{short: c, value: value, hasValue: hasValue}
}

type ref struct {
	long     string
	short    byte
	value    string
	hasValue bool
}

func (r ref) LongForm() string   { return r.long }
func (r ref) HasLongForm() bool  { return len(r.long) > 0 }
func (r ref) ShortForm() byte    { return r.short }
func (r ref) HasShortForm() bool { return r.short != 0 }
func (r ref) Value() string      { return r.value }
func (r ref) HasValue() bool     { return r.hasValue }
