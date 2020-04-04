package impl

import (
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

type Command struct {
}

func (c *Command) Description() string {
	panic("implement me")
}

func (c *Command) Arguments() []argo.Argument {
	panic("implement me")
}

func (c *Command) Flags() []argo.Flag {
	panic("implement me")
}

func (c *Command) UnmappedInput() []string {
	panic("implement me")
}
