package impl

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"os"
)

type Command struct {
	description string
	unmarshal   A.ValueUnmarshaler

	groups      []A.FlagGroup
	arguments   []A.Argument
}

func (c *Command) LookupFlag(key interface{}) (match A.Flag, found bool) {
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

func (c *Command) FlagGroups() []A.FlagGroup {
	return c.groups
}

func (c *Command) Description() string {
	return c.description
}

func (c *Command) Arguments() []A.Argument {
	return c.arguments
}

func (c *Command) UnmappedInput() []string {
	panic("implement me")
}

func (c *Command) Unmarshaler() A.ValueUnmarshaler {
	return c.unmarshal
}

func (c *Command) Name() string {
	return os.Args[0]
}

func (c *Command) String() string {
	return c.Name()
}
