package command

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Command struct {
	helpDisabled bool

	description   string
	unmappedLabel string

	options argo.CommandOptions

	flagGroups []argo.FlagGroup
	arguments  []argo.Argument

	unmapped []string

	callback argo.CommandCallback[argo.Command]
}

func (c *Command) Name() string {
	return filepath.Base(os.Args[0])
}

func (c *Command) HasDescription() bool {
	return len(c.description) > 0
}

func (c *Command) Description() string {
	return c.description
}

func (c *Command) IsHelpDisabled() bool {
	return c.helpDisabled
}

func (c *Command) HasFlagGroups() bool {
	return len(c.flagGroups) > 0
}

func (c *Command) FlagGroups() []argo.FlagGroup {
	return c.flagGroups
}

func (c *Command) hasDefaultFlagGroup() bool {
	if c.flagGroups[0].Name() == flag.DefaultFlagGroupName {
		return true
	}

	return false
}

func (c *Command) FindShortFlag(b byte) argo.Flag {
	for _, group := range c.flagGroups {
		if shortFlag := group.FindShortFlag(b); shortFlag != nil {
			return shortFlag
		}
	}

	return nil
}

func (c *Command) FindLongFlag(name string) argo.Flag {
	for _, group := range c.flagGroups {
		if longFlag := group.FindLongFlag(name); longFlag != nil {
			return longFlag
		}
	}

	return nil
}

func (c *Command) HasArguments() bool {
	return len(c.arguments) > 0
}

func (c *Command) Arguments() []argo.Argument {
	return c.arguments
}

func (c *Command) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *Command) UnmappedInputs() []string {
	return c.unmapped
}

func (c *Command) AppendUnmappedInput(input string) {
	c.unmapped = append(c.unmapped, input)
}

func (c *Command) HasUnmappedInputLabel() bool {
	return len(c.unmappedLabel) > 0
}

func (c *Command) UnmappedInputLabel() string {
	return c.unmappedLabel
}

func (c *Command) HasCallback() bool {
	return c.callback != nil
}

func (c *Command) Callback() argo.CommandCallback[argo.Command] {
	return c.callback
}

func (c *Command) Options() argo.CommandOptions {
	return c.options
}
