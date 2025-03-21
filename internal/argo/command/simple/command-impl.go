package command

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type command struct {
	description   string
	unmappedLabel string

	flagGroups []argo.FlagGroup
	arguments  []argo.Argument

	unmapped    []string
	passthrough []string

	callback argo.CommandCallback
}

func (c *command) Name() string {
	return filepath.Base(os.Args[0])
}

//

func (c *command) HasDescription() bool {
	return len(c.description) > 0
}

func (c *command) Description() string {
	return c.description
}

//

func (c *command) HasFlagGroups(includeDefault bool) bool {
	if len(c.flagGroups) == 0 {
		return false
	}

	if c.hasDefaultFlagGroup() {
		if includeDefault {
			return true
		}

		return len(c.flagGroups) > 1
	}

	return len(c.flagGroups) > 0
}

func (c *command) FlagGroups(includeDefault bool) []argo.FlagGroup {
	if c.hasDefaultFlagGroup() {
		if includeDefault {
			return c.flagGroups
		}

		return c.flagGroups[1:]
	}

	return c.flagGroups
}

func (c *command) hasDefaultFlagGroup() bool {
	if c.flagGroups[0].Name() == flag.DefaultFlagGroupName {
		return true
	}

	return false
}

//

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

//

func (c *command) HasArguments() bool {
	return len(c.arguments) > 0
}

func (c *command) Arguments() []argo.Argument {
	return c.arguments
}

//

func (c *command) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *command) UnmappedInputs() []string {
	return c.unmapped
}

func (c *command) AppendUnmappedInput(input string) {
	c.unmapped = append(c.unmapped, input)
}

//

func (c *command) HasUnmappedInputLabel() bool {
	return len(c.unmappedLabel) > 0
}

func (c *command) UnmappedInputLabel() string {
	return c.unmappedLabel
}

//

func (c *command) HasCallback() bool {
	return c.callback != nil
}

func (c *command) Callback() argo.CommandCallback {
	return c.callback
}
