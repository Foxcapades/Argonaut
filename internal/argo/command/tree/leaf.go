package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Leaf struct {
	disableHelp   bool
	description   string
	flagGroups    []argo.FlagGroup
	callback      argo.CommandCallback[argo.LeafCommand]
	name          string
	parent        argo.ParentNode
	aliases       []string
	unmappedLabel string
	args          []argo.Argument
	unmapped      []string
}

func (i *Leaf) HasArguments() bool {
	return len(i.args) > 0
}

func (i *Leaf) Arguments() []argo.Argument {
	return i.args
}

func (i *Leaf) HasUnmappedInputs() bool {
	return len(i.unmapped) > 0
}

func (i *Leaf) UnmappedInputs() []string {
	return i.unmapped
}

func (i *Leaf) AppendUnmappedInput(input string) {
	i.unmapped = append(i.unmapped, input)
}

func (i *Leaf) HasUnmappedInputLabel() bool {
	return len(i.unmappedLabel) > 0
}

func (i *Leaf) UnmappedInputLabel() string {
	return i.unmappedLabel
}

func (i *Leaf) Name() string {
	return i.name
}

func (i *Leaf) Parent() argo.ParentNode {
	return i.parent
}

func (i *Leaf) HasAliases() bool {
	return len(i.aliases) > 0
}

func (i *Leaf) Aliases() []string {
	return i.aliases
}

func (i *Leaf) Matches(name string) bool {
	if i.name == name {
		return true
	}

	for _, alias := range i.aliases {
		if alias == name {
			return true
		}
	}

	return false
}

func (i *Leaf) HasDescription() bool {
	return len(i.description) > 0
}

func (i *Leaf) Description() string {
	return i.description
}

func (i *Leaf) HasFlagGroups() bool {
	return len(i.flagGroups) > 0
}

func (i *Leaf) FlagGroups() []argo.FlagGroup {
	return i.flagGroups
}

func (i *Leaf) HasCallback() bool {
	return i.callback != nil
}

func (i *Leaf) Callback() argo.CommandCallback[argo.LeafCommand] {
	return i.callback
}

func (i *Leaf) FindShortFlag(b byte) argo.Flag {
	for _, group := range i.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (i *Leaf) FindLongFlag(name string) argo.Flag {
	for _, group := range i.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (i *Leaf) IsHelpDisabled() bool {
	return i.disableHelp
}
