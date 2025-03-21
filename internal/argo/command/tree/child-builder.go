package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

func newChildBuilder[T any](name string) childBuilder[T] {
	return childBuilder[T]{
		nodeBuilder: newNodeBuilder[T](),
		name:        name,
	}
}

type childBuilder[T any] struct {
	nodeBuilder[T]

	name    string
	aliases []string
	parent  argo.ParentBuilder[any]
}

func (c *childBuilder[T]) Name() string {
	return c.name
}

func (c *childBuilder[T]) WithAlias(alias string) T {
	c.aliases = append(c.aliases, alias)
	return T(any(c))
}

func (c *childBuilder[T]) WithAliases(aliases ...string) T {
	c.aliases = append(c.aliases, aliases...)
	return T(any(c))
}

func (c *childBuilder[T]) HasAliases() bool {
	return len(c.aliases) > 0
}

func (c *childBuilder[T]) Aliases() []string {
	return c.aliases
}

func (c *childBuilder[T]) ParentNode() argo.ParentBuilder[any] {
	return c.parent
}

func (c *childBuilder[T]) SetParentNode(parent argo.ParentBuilder[any]) T {
	c.parent = parent
	return T(any(c))
}
