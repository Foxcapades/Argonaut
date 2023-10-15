package command

import "github.com/Foxcapades/Argonaut/pkg/argo"

// implements argo.CommandLeaf
type commandLeaf struct {
	name        string
	desc        string
	parent      argo.CommandNode
	aliases     []string
	flags       []argo.FlagGroup
	args        []argo.Argument
	unmapped    []string
	passthrough []string
	onHit       argo.CommandLeafCallback
}

func (c commandLeaf) Parent() argo.CommandNode { return c.parent }
func (c commandLeaf) HasParent() bool          { return c.parent != nil }

func (c commandLeaf) Name() string { return c.name }

func (c commandLeaf) CallOnHit() {
	if c.onHit != nil {
		c.onHit(&c)
	}
}

func (c commandLeaf) HasOnHit() bool {
	return c.onHit != nil
}

func (c commandLeaf) Aliases() []string { return c.aliases }
func (c commandLeaf) HasAliases() bool  { return len(c.aliases) > 0 }

func (c commandLeaf) Description() string  { return c.desc }
func (c commandLeaf) HasDescription() bool { return len(c.desc) > 0 }

func (c commandLeaf) FlagGroups() []argo.FlagGroup { return c.flags }
func (c commandLeaf) HasFlagGroups() bool          { return len(c.flags) > 0 }

func (c commandLeaf) Arguments() []argo.Argument { return c.args }
func (c commandLeaf) HasArguments() bool         { return len(c.args) > 0 }
func (c *commandLeaf) AppendArgument(val string) error {
	for _, arg := range c.args {
		if !arg.WasHit() {
			return arg.SetValue(val)
		}
	}

	c.AppendUnmapped(val)
	return nil
}

func (c commandLeaf) UnmappedInputs() []string   { return c.unmapped }
func (c commandLeaf) HasUnmappedInputs() bool    { return len(c.unmapped) > 0 }
func (c *commandLeaf) AppendUnmapped(val string) { c.unmapped = append(c.unmapped, val) }

func (c commandLeaf) PassthroughInputs() []string   { return c.passthrough }
func (c commandLeaf) HasPassthroughInputs() bool    { return len(c.passthrough) > 0 }
func (c *commandLeaf) AppendPassthrough(val string) { c.passthrough = append(c.passthrough, val) }

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

func (c commandLeaf) FindShortFlag(b byte) argo.Flag {
	var current argo.CommandNode = c

	for current != nil {
		for _, group := range current.FlagGroups() {
			if flag := group.FindShortFlag(b); flag != nil {
				return flag
			}
		}

		current = current.Parent()
	}

	return nil
}

func (c commandLeaf) FindLongFlag(name string) argo.Flag {
	var current argo.CommandNode = c

	for current != nil {
		for _, group := range current.FlagGroups() {
			if flag := group.FindLongFlag(name); flag != nil {
				return flag
			}
		}

		current = current.Parent()
	}

	return nil
}

func (c commandLeaf) TryFlag(ref argo.FlagRef) (bool, error) {
	var current argo.CommandNode = c

	for current != nil {
		for _, group := range current.FlagGroups() {
			if ok, err := group.TryFlag(ref); ok || err != nil {
				return ok, err
			}
		}

		current = current.Parent()
	}

	return false, nil
}
