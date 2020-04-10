package command

import (
	"os"
	"path"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
)

type Command struct {
	trait.Described

	// value unmarshaller
	unmarshal AVU

	// flag groups
	groups []AFG

	// positional arguments
	arguments []AA

	// unknown inputs
	unmapped []string

	// passthrough values (values appearing after "--" on the
	// cli)
	passthroughs []string
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Interface Implementation                                            ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// See argo.Command#FlagGroups()
func (c *Command) FlagGroups() []AFG {
	return c.groups
}

// See argo.Command#Arguments()
func (c *Command) Arguments() []AA {
	return c.arguments
}

// See argo.Command#UnmappedInput()
func (c *Command) UnmappedInput() []string {
	return c.unmapped
}

// See argo.Command#Passthroughs()
func (c *Command) Passthroughs() []string {
	return c.passthroughs
}

// See argo.Command#Unmarshaller()
func (c *Command) Unmarshaler() AVU {
	return c.unmarshal
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
	c.unmapped = append(c.unmapped, val)
}

// See argo.Command#AppendPassthrough()
func (c *Command) AppendPassthrough(val string) {
	c.passthroughs = append(c.passthroughs, val)
}
