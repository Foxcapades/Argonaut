package flag

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type FlagGroup struct {
	parent A.Command

	desc  string
	name  string
	flags []A.Flag
}

func (f *FlagGroup) Name() string         { return f.name }
func (f *FlagGroup) HasName() bool        { return len(f.name) > 0 }
func (f *FlagGroup) Description() string  { return f.desc }
func (f *FlagGroup) HasDescription() bool { return len(f.desc) > 0 }
func (f *FlagGroup) Flags() []A.Flag      { return f.flags }
func (f *FlagGroup) HasFlags() bool       { return len(f.flags) > 0 }
func (f *FlagGroup) Parent() A.Command    { return f.parent }
