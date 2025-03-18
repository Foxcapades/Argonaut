package command

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type command struct {
	warnings      *argo.WarningContext
	description   string
	unmappedLabel string
	flagGroups    []argo.FlagGroup
	arguments     []argo.Argument
	unmapped      []string
	passthrough   []string
	callback      argo.CommandCallback
}

func (c *command) Name() string {
	return filepath.Base(os.Args[0])
}

func (c *command) Description() string {
	return c.description
}

func (c *command) HasDescription() bool {
	return len(c.description) > 0
}

func (c *command) FlagGroups() []argo.FlagGroup {
	return c.flagGroups
}

func (c *command) HasFlagGroups() bool {
	return len(c.flagGroups) > 0
}

func (c *command) HasUnmappedLabel() bool {
	return len(c.unmappedLabel) > 0
}

func (c *command) GetUnmappedLabel() string {
	return c.unmappedLabel
}

func (c *command) FindShortFlag(b byte) argo.Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (c *command) FindLongFlag(name string) argo.Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (c *command) Arguments() []argo.Argument {
	return c.arguments
}

func (c *command) HasArguments() bool {
	return len(c.arguments) > 0
}

func (c *command) appendArgument(rawArgument string) error {
	for _, arg := range c.arguments {
		if !arg.WasHit() {
			return arg.SetValue(rawArgument)
		}
	}

	c.unmapped = append(c.unmapped, rawArgument)
	return nil
}

func (c *command) UnmappedInputs() []string {
	return c.unmapped
}

func (c *command) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *command) appendUnmapped(val string) {
	c.unmapped = append(c.unmapped, val)
}

func (c *command) PassthroughInputs() []string {
	return c.passthrough
}

func (c *command) HasPassthroughInputs() bool {
	return len(c.passthrough) > 0
}

func (c *command) executeCallback() {
	if c.callback != nil {
		c.callback(c)
	}
}

func (c *command) appendPassthrough(val string) {
	c.passthrough = append(c.passthrough, val)
}

func (c *command) Warnings() []string {
	return c.warnings.GetWarnings()
}

func (c *command) AppendWarning(warning string) {
	c.warnings.AppendWarning(warning)
}
