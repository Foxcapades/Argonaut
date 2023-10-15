package command

import (
	"os"
	"path/filepath"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

type command struct {
	description string
	flagGroups  []argo.FlagGroup
	arguments   []argo.Argument
	unmapped    []string
	passthrough []string
}

func (c command) Name() string {
	return filepath.Base(os.Args[0])
}

func (c command) Description() string {
	return c.description
}

func (c command) HasDescription() bool {
	return len(c.description) > 0
}

func (c command) FlagGroups() []argo.FlagGroup {
	return c.flagGroups
}

func (c command) HasFlagGroups() bool {
	return len(c.flagGroups) > 0
}

func (c command) FindShortFlag(b byte) argo.Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (c command) FindLongFlag(name string) argo.Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (c command) TryFlag(ref argo.FlagRef) (bool, error) {
	for _, group := range c.flagGroups {
		if ok, err := group.TryFlag(ref); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}

func (c command) Arguments() []argo.Argument {
	return c.arguments
}

func (c command) HasArguments() bool {
	return len(c.arguments) > 0
}

func (c *command) AppendArgument(rawArgument string) error {
	for _, arg := range c.arguments {
		if !arg.WasHit() {
			if err := arg.SetValue(rawArgument); err != nil {
				return err
			}
		}
	}

	c.unmapped = append(c.unmapped, rawArgument)
	return nil
}

func (c command) UnmappedInputs() []string {
	return c.unmapped
}

func (c command) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *command) AppendUnmapped(val string) {
	c.unmapped = append(c.unmapped, val)
}

func (c command) PassthroughInputs() []string {
	return c.passthrough
}

func (c command) HasPassthroughInputs() bool {
	return len(c.passthrough) > 0
}

func (c *command) AppendPassthrough(val string) {
	c.passthrough = append(c.passthrough, val)
}
