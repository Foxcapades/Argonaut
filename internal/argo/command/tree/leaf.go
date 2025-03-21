package tree

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandLeaf struct {
	child[argo.LeafCommand]

	unmappedLabel string
	args          []argo.Argument
	unmapped      []string
}

func (c *commandLeaf) HasArguments() bool {
	return len(c.args) > 0
}

func (c *commandLeaf) Arguments() []argo.Argument {
	return c.args
}

func (c *commandLeaf) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *commandLeaf) UnmappedInputs() []string {
	return c.unmapped
}

func (c *commandLeaf) AppendUnmappedInput(input string) {
	c.unmapped = append(c.unmapped, input)
}

func (c *commandLeaf) HasUnmappedLabel() bool {
	return len(c.unmappedLabel) > 0
}

func (c *commandLeaf) GetUnmappedLabel() string {
	return c.unmappedLabel
}
