package command

import (
	"os"
	"path"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type Command struct {
	trait.Described

	ValueUnmarshaler argo.ValueUnmarshaler

	// flag groups
	Groups []argo.FlagGroup

	// positional arguments
	PositionalArgs []argo.Argument

	// unknown inputs
	Unmapped []string

	// passthrough values (values appearing after "--" on the
	// cli)
	Passthrough []string
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Interface Implementation                                            ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// See argo.Command#FlagGroups()
func (c *Command) FlagGroups() []AFG {
	return c.Groups
}

// See argo.Command#Arguments()
func (c *Command) Arguments() []AA {
	return c.PositionalArgs
}

// See argo.Command#UnmappedInput()
func (c *Command) UnmappedInput() []string {
	return c.Unmapped
}

// See argo.Command#Passthroughs()
func (c *Command) Passthroughs() []string {
	return c.Passthrough
}

// See argo.Command#Unmarshaller()
func (c *Command) Unmarshaler() AVU {
	return c.ValueUnmarshaler
}

func (c *Command) String() string {
	return c.Name()
}

// See argo.Command#Name()
func (c *Command) Name() string {
	return path.Base(os.Args[0])
}

// See argo.Command#AppendUnmapped()
func (c *Command) AppendUnmapped(val string) {
	c.Unmapped = append(c.Unmapped, val)
}

// See argo.Command#AppendPassthrough()
func (c *Command) AppendPassthrough(val string) {
	c.Passthrough = append(c.Passthrough, val)
}
