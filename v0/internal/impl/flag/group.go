package flag

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type Group struct {
	trait.Named
	trait.Described

	ParentElement A.Command

	FlagElements []A.Flag
}

func (f *Group) Flags() []A.Flag {
	return f.FlagElements
}

func (f *Group) HasFlags() bool {
	return len(f.FlagElements) > 0
}

func (f *Group) Parent() A.Command {
	return f.ParentElement
}
