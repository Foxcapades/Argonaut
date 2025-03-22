package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

func newChildBuilder[T, O any](name string, root T) childBuilder[T, O] {
	return childBuilder[T, O]{
		nodeBuilder: newNodeBuilder[T, O](root),
		name:        name,
	}
}

type childBuilder[T, O any] struct {
	nodeBuilder[T, O]

	name    string
	aliases []string
	parent  argo.ParentBuilder[any, any]
}

func (c *childBuilder[T, O]) Name() string {
	return c.name
}

func (c *childBuilder[T, O]) WithAlias(alias string) T {
	c.aliases = append(c.aliases, alias)
	return T(any(c))
}

func (c *childBuilder[T, O]) WithAliases(aliases ...string) T {
	c.aliases = append(c.aliases, aliases...)
	return T(any(c))
}

func (c *childBuilder[T, O]) HasAliases() bool {
	return len(c.aliases) > 0
}

func (c *childBuilder[T, O]) Aliases() []string {
	return c.aliases
}

func (c *childBuilder[T, O]) ParentNode() argo.ParentBuilder[any, any] {
	return c.parent
}

func (c *childBuilder[T, O]) SetParentNode(parent argo.ParentBuilder[any, any]) T {
	c.parent = parent
	return T(any(c))
}
