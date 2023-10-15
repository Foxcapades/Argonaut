package flag

import (
	"errors"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func GroupBuilder(name string) argo.FlagGroupBuilder {
	return &groupBuilder{
		name: name,
	}
}

type groupBuilder struct {
	name  string
	desc  string
	flags []argo.FlagBuilder
}

func (g groupBuilder) GetName() string {
	return g.name
}

func (g *groupBuilder) WithDescription(desc string) argo.FlagGroupBuilder {
	g.desc = desc
	return g
}

func (g groupBuilder) HasDescription() bool {
	return len(g.desc) > 0
}

func (g groupBuilder) GetDescription() string {
	return g.desc
}

func (g *groupBuilder) WithFlag(flag argo.FlagBuilder) argo.FlagGroupBuilder {
	g.flags = append(g.flags, flag)
	return g
}

func (g groupBuilder) HasFlags() bool {
	return len(g.flags) > 0
}

func (g groupBuilder) GetFlags() []argo.FlagBuilder {
	return g.flags
}

func (g *groupBuilder) Build() (argo.FlagGroup, error) {
	errs := xerr.NewMultiError()
	flags := make([]argo.Flag, len(g.flags))

	// Ensure the group name is not blank
	if chars.IsBlank(g.name) {
		errs.AppendError(errors.New("flag group names must not be blank"))
	}

	for i := range g.flags {
		if flag, err := g.flags[i].Build(); err != nil {
			errs.AppendError(err)
		} else {
			flags = append(flags, flag)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &group{
		name:  g.name,
		desc:  g.desc,
		flags: flags,
	}, nil
}
