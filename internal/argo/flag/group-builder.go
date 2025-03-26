package flag

import "github.com/foxcapades/argonaut/v3/pkg/argo"

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

func (g *groupBuilder) WithFlags(flags ...argo.FlagBuilder) argo.FlagGroupBuilder {
	g.flags = append(g.flags, flags...)
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
