package argo

type CommandLeaf interface {
	CommandNode
	CommandChild
	Command

	hasCallback() bool

	executeCallback()
}

type CommandLeafCallback = func(leaf CommandLeaf)

type commandLeaf struct {
	name        string
	desc        string
	uLabel      string
	parent      CommandNode
	aliases     []string
	flags       []FlagGroup
	args        []Argument
	unmapped    []string
	passthrough []string
	warnings    *WarningContext
	callback    CommandLeafCallback
}

func (c commandLeaf) Parent() CommandNode { return c.parent }
func (c commandLeaf) HasParent() bool     { return c.parent != nil }

func (c commandLeaf) Name() string { return c.name }

func (c commandLeaf) executeCallback() {
	if c.callback != nil {
		c.callback(&c)
	}
}

func (c commandLeaf) hasCallback() bool {
	return c.callback != nil
}

func (c commandLeaf) Aliases() []string {
	return c.aliases
}

func (c commandLeaf) HasAliases() bool {
	return len(c.aliases) > 0
}

func (c commandLeaf) Description() string {
	return c.desc
}

func (c commandLeaf) HasDescription() bool {
	return len(c.desc) > 0
}

func (c commandLeaf) FlagGroups() []FlagGroup {
	return c.flags
}

func (c commandLeaf) HasFlagGroups() bool {
	return len(c.flags) > 0
}

func (c commandLeaf) HasUnmappedLabel() bool {
	return len(c.uLabel) > 0
}

func (c commandLeaf) GetUnmappedLabel() string {
	return c.uLabel
}

func (c commandLeaf) Arguments() []Argument {
	return c.args
}

func (c commandLeaf) HasArguments() bool {
	return len(c.args) > 0
}

func (c *commandLeaf) appendArgument(val string) error {
	for _, arg := range c.args {
		if !arg.WasHit() {
			return arg.setValue(val)
		}
	}

	c.appendUnmapped(val)
	return nil
}

func (c commandLeaf) UnmappedInputs() []string {
	return c.unmapped
}

func (c commandLeaf) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *commandLeaf) appendUnmapped(val string) {
	c.unmapped = append(c.unmapped, val)
}

func (c commandLeaf) PassthroughInputs() []string {
	return c.passthrough
}

func (c commandLeaf) HasPassthroughInputs() bool {
	return len(c.passthrough) > 0
}

func (c *commandLeaf) appendPassthrough(val string) {
	c.passthrough = append(c.passthrough, val)
}

func (c commandLeaf) Matches(name string) bool {
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

func (c commandLeaf) FindShortFlag(b byte) Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return c.parent.FindShortFlag(b)
}

func (c commandLeaf) FindLongFlag(name string) Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return c.parent.FindLongFlag(name)
}

func (c commandLeaf) Warnings() []string {
	return c.warnings.GetWarnings()
}

func (c commandLeaf) AppendWarning(warning string) {
	c.warnings.appendWarning(warning)
}
