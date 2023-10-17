package argo

import "errors"

// A FlagGroupBuilder is used to build a group or category of flags that go
// together.  This is primarily used for rendering help text.
//
// FlagGroups are not required but allow for organizing the help text when there
// are many flags attached to a command.
type FlagGroupBuilder interface {
	getName() string

	// WithDescription sets a description value on this FlagGroupBuilder.
	WithDescription(desc string) FlagGroupBuilder

	// WithFlag appends the given FlagBuilder instance to this FlagGroupBuilder.
	WithFlag(flag FlagBuilder) FlagGroupBuilder

	hasFlags() bool
	size() int
	getFlags() []FlagBuilder
	build() (FlagGroup, error)
}

func NewFlagGroupBuilder(name string) FlagGroupBuilder {
	return &flagGroupBuilder{
		name: name,
	}
}

type flagGroupBuilder struct {
	name  string
	desc  string
	flags []FlagBuilder
}

func (g flagGroupBuilder) getName() string {
	return g.name
}

func (g *flagGroupBuilder) WithDescription(desc string) FlagGroupBuilder {
	g.desc = desc
	return g
}

func (g *flagGroupBuilder) WithFlag(flag FlagBuilder) FlagGroupBuilder {
	g.flags = append(g.flags, flag)
	return g
}

func (g flagGroupBuilder) hasFlags() bool {
	return len(g.flags) > 0
}

func (g flagGroupBuilder) getFlags() []FlagBuilder {
	return g.flags
}

func (g flagGroupBuilder) size() int {
	return len(g.flags)
}

func (g *flagGroupBuilder) build() (FlagGroup, error) {
	errs := newMultiError()
	flags := make([]Flag, 0, len(g.flags))

	// Ensure the group name is not blank
	if isBlank(g.name) {
		errs.AppendError(errors.New("flag group names must not be blank"))
	}

	for i := range g.flags {
		if flag, err := g.flags[i].build(); err != nil {
			errs.AppendError(err)
		} else {
			flags = append(flags, flag)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &flagGroup{
		name:  g.name,
		desc:  g.desc,
		flags: flags,
	}, nil
}
