package flag

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type Group struct {
	trait.Named
	trait.Described

	parent A.Command

	flags []A.Flag
}

func (f *Group) Flags() []A.Flag   { return f.flags }
func (f *Group) HasFlags() bool    { return len(f.flags) > 0 }
func (f *Group) Parent() A.Command { return f.parent }
