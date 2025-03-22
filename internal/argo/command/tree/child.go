package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type child[T any] struct {
	common.CommandBase[T]

	name    string
	parent  argo.ParentNode[any]
	aliases []string
}

func (c child[T]) Name() string {
	return c.name
}

func (c child[T]) Parent() argo.ParentNode[any] {
	return c.parent
}

func (c child[T]) HasAliases() bool {
	return len(c.aliases) > 0
}

func (c child[T]) Aliases() []string {
	return c.aliases
}

func (c child[T]) Matches(name string) bool {
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
