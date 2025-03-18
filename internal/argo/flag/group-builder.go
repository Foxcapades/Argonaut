package flag

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewGroupBuilder(name string) argo.FlagGroupBuilder {
	return &groupBuilder{
		name: name,
	}
}

type groupBuilder struct {
	name  string
	desc  string
	flags []argo.FlagBuilder
}

func (g *groupBuilder) Name() string {
	return g.name
}

func (g *groupBuilder) WithDescription(desc string) argo.FlagGroupBuilder {
	g.desc = desc
	return g
}

func (g *groupBuilder) HasDescription() bool {
	return len(g.desc) > 0
}

func (g *groupBuilder) Description() string {
	return g.desc
}

func (g *groupBuilder) WithFlag(flag argo.FlagBuilder) argo.FlagGroupBuilder {
	g.flags = append(g.flags, flag)
	return g
}

func (g *groupBuilder) HasFlags() bool {
	return len(g.flags) > 0
}

func (g *groupBuilder) Flags() []argo.FlagBuilder {
	return g.flags
}

func (g *groupBuilder) Size() int {
	return len(g.flags)
}

func (g *groupBuilder) Build(ctx *argo.WarningContext) (argo.FlagGroup, error) {
	errs := xerr.NewMultiError()
	flags := make([]argo.Flag, 0, len(g.flags))

	// Ensure the group name is not blank
	if chars.IsBlank(g.name) {
		errs.AppendError(errors.New("flag group names must not be blank"))
	}

	for i := range g.flags {
		if flag, err := g.flags[i].Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			flags = append(flags, flag)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &group{
		warnings: ctx,
		name:     g.name,
		desc:     g.desc,
		flags:    flags,
	}, nil
}
