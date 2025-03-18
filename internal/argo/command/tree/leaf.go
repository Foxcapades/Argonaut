package tree

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandLeaf struct {
	name        string
	desc        string
	uLabel      string
	parent      argo.Node
	aliases     []string
	flags       []argo.FlagGroup
	args        []argo.Argument
	unmapped    []string
	passthrough []string
	warnings    *argo.WarningContext
	callback    argo.CommandLeafCallback
}

func (c *commandLeaf) Parent() argo.Node { return c.parent }
func (c *commandLeaf) HasParent() bool   { return c.parent != nil }

func (c *commandLeaf) Name() string { return c.name }

func (c *commandLeaf) executeCallback() {
	if c.callback != nil {
		c.callback(&c)
	}
}

func (c *commandLeaf) hasCallback() bool {
	return c.callback != nil
}

func (c *commandLeaf) Aliases() []string {
	return c.aliases
}

func (c *commandLeaf) HasAliases() bool {
	return len(c.aliases) > 0
}

func (c *commandLeaf) Description() string {
	return c.desc
}

func (c *commandLeaf) HasDescription() bool {
	return len(c.desc) > 0
}

func (c *commandLeaf) FlagGroups() []argo.FlagGroup {
	return c.flags
}

func (c *commandLeaf) HasFlagGroups() bool {
	return len(c.flags) > 0
}

func (c *commandLeaf) HasUnmappedLabel() bool {
	return len(c.uLabel) > 0
}

func (c *commandLeaf) GetUnmappedLabel() string {
	return c.uLabel
}

func (c *commandLeaf) Arguments() []argo.Argument {
	return c.args
}

func (c *commandLeaf) HasArguments() bool {
	return len(c.args) > 0
}

func (c *commandLeaf) UnmappedInputs() []string {
	return c.unmapped
}

func (c *commandLeaf) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *commandLeaf) appendUnmapped(val string) {
	c.unmapped = append(c.unmapped, val)
}

func (c *commandLeaf) PassthroughInputs() []string {
	return c.passthrough
}

func (c *commandLeaf) HasPassthroughInputs() bool {
	return len(c.passthrough) > 0
}

func (c *commandLeaf) appendPassthrough(val string) {
	c.passthrough = append(c.passthrough, val)
}

func (c *commandLeaf) Matches(name string) bool {
	if c.name == name {
		return true
	}

	for _, alias := range c.aliases {
		if alias == name {
			return true
		}
	}

	return false
}

func (c *commandLeaf) FindShortFlag(b byte) argo.Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return c.parent.FindShortFlag(b)
}

func (c *commandLeaf) FindLongFlag(name string) argo.Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return c.parent.FindLongFlag(name)
}

func (c *commandLeaf) Warnings() []string {
	return c.warnings.GetWarnings()
}

func (c *commandLeaf) AppendWarning(warning string) {
	c.warnings.AppendWarning(warning)
}
