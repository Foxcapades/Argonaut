package command

import (
	"errors"
	"fmt"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/props"
	"os"

	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/marsh"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/parse"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/internal/render"
)

func NewBuilder(provider AP) ACB {
	return &Builder{
		provider:    provider,
		fGroups:     []AFGB{provider.NewFlagGroup().Name("Options")},
		parser:      parse.NewParser(),
		unmarshaler: marsh.NewDefaultedValueUnmarshaler(),
		options:     props.DefaultCommandOptions(),
	}
}

type Builder struct {
	name trait.Named
	desc trait.Described

	provider    AP
	parser      A.Parser
	unmarshaler AVU

	fGroups  []AFGB
	args     []AAB
	warnings []string
	examples []string

	options props.CommandOptions
}

func (c *Builder) DisableHelp() ACB      { c.options.IncludeHelp = false; return c }
func (c *Builder) GetFlagGroups() []AFGB { return c.fGroups }
func (c *Builder) GetArgs() []AAB        { return c.args }

func (c *Builder) Examples(examples ...string) ACB { c.examples = examples; return c }

func (c *Builder) Description(desc string) ACB { c.desc.DescriptionText = desc; return c }
func (c *Builder) HasDescription() bool        { return len(c.desc.DescriptionText) > 0 }
func (c *Builder) GetDescription() string      { return c.desc.Description() }

func (c *Builder) Warnings() []string { return c.warnings }

func (c *Builder) Flag(flag A.FlagBuilder) ACB {
	if flag == nil {
		return c.warn("nil value passed to Flag()")
	}
	c.fGroups[0].Flag(flag)
	return c
}

func (c *Builder) FlagGroup(builder AFGB) (this ACB) {
	if builder == nil {
		return c.warn("nil value passed to FlagGroup()")
	}
	c.fGroups = append(c.fGroups, builder)
	return c
}

func (c *Builder) Unmarshaler(un AVU) ACB {
	if un == nil {
		return c.warn("nil value passed to Unmarshaler()")
	}
	c.unmarshaler = un
	return c
}

func (c *Builder) Arg(arg AAB) ACB {
	if arg == nil {
		return c.warn("nil value passed to Arg()")
	}
	c.args = append(c.args, arg)
	return c
}

func (c *Builder) Build() (AC, error) {
	short := make(map[byte]bool)
	long := make(map[string]bool)
	out := new(Command)

	// Check for overlapping flags
	for _, g := range c.fGroups {
		for _, f := range g.GetFlags() {
			if f.HasShort() {
				if short[f.GetShort()] {
					// TODO: make this a real error
					return nil, errors.New("more than one flag with the identifier -" + string(f.GetShort()))
				}
				short[f.GetShort()] = true
			}
			if f.HasLong() {
				if long[f.GetLong()] {
					// TODO: make this a real error
					return nil, errors.New("more than one flag with the identifier --" + f.GetLong())
				}
				long[f.GetLong()] = true
			}
		}
	}

	if c.options.IncludeHelp {
		c.fGroups = append(c.fGroups, c.provider.
			NewFlagGroup().
			Name("Help & Info").
			Flag(c.provider.
				NewFlag().
				Short('h').
				Long("help").
				Description("Prints this help text").
				OnHit(func(f AF) {
					fmt.Println(render.Command(out))
					os.Exit(0)
				})))
	}

	// Build groups
	groups := make([]AFG, len(c.fGroups))

	for i, fg := range c.fGroups {
		fg.Parent(out)
		if g, err := fg.Build(); err != nil {
			return nil, err
		} else {
			groups[i] = g
		}
	}

	args := make([]AA, len(c.args))
	for i, arg := range c.args {
		arg.Parent(out)
		if a, err := arg.Build(); err != nil {
			return nil, err
		} else {
			args[i] = a
		}
	}

	out.Described = c.desc
	out.Groups = groups
	out.PositionalArgs = args
	out.ValueUnmarshaler = c.unmarshaler

	return out, nil
}

func (c *Builder) MustBuild() AC {
	com, err := c.Build()
	if err != nil {
		panic(err)
	}
	return com
}

func (c *Builder) Parse() (extra []string, err error) {
	com, err := c.Build()
	if err != nil {
		return nil, err
	}

	err = c.parser.Parse(os.Args, com)

	return c.parser.Passthroughs(), err
}

func (c *Builder) MustParse() []string {
	if a, b := c.Parse(); b != nil {
		panic(b)
	} else {
		return a
	}
}

func (c *Builder) warn(txt string) ACB {
	c.warnings = append(c.warnings, "CommandBuilder: "+txt)
	return c
}
