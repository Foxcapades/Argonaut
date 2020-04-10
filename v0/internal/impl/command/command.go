package command

import (
	"os"
	"path"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
)

type Command struct {
	trait.Described
	unmarshal AVU

	groups    []AFG
	arguments []AA
	unmapped  []string
}

func (c *Command) FlagGroups() []AFG       { return c.groups }
func (c *Command) Arguments() []AA         { return c.arguments }
func (c *Command) UnmappedInput() []string { return c.unmapped }
func (c *Command) Unmarshaler() AVU        { return c.unmarshal }
func (c *Command) String() string          { return c.Name() }
func (c *Command) Name() string            { return path.Base(os.Args[0]) }

func (c *Command) LookupFlag(key interface{}) (match AF, found bool) {
	for _, g := range c.groups {
		for _, f := range g.Flags() {
			if f.HasShort() && f.Short() == key {
				return f, true
			}
			if f.HasLong() && f.Long() == key {
				return f, true
			}
		}
	}
	return nil, false
}
