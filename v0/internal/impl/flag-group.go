package impl

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type FlagGroup struct {
	description string
	name        string
	flags       []A.Flag
}

func (f *FlagGroup) Name() string {
	return f.name
}

func (f *FlagGroup) HasName() bool {
	return len(f.name) > 0
}

func (f *FlagGroup) Description() string {
	return f.description
}

func (f *FlagGroup) HasDescription() bool {
	return len(f.description) > 0
}

func (f *FlagGroup) Flags() []A.Flag {
	return f.flags
}

func (f *FlagGroup) HasFlags() bool {
	return len(f.flags) > 0
}
