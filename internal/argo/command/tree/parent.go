package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type parent[T any] struct {
	node[T]

	groups []argo.CommandGroup

	selectedChild argo.ChildNode[any]

	incompleteFn argo.IncompleteCommandHandler[T]
}

func (p *parent[T]) HasCommandGroups(includeDefault bool) bool {
	return len(p.groups) > 1 || (includeDefault && p.hasDefaultCommandGroup())
}

func (p *parent[T]) CommandGroups(includeDefault bool) []argo.CommandGroup {
	if !p.HasCommandGroups(includeDefault) {
		return nil
	}

	if includeDefault && p.hasDefaultCommandGroup() {
		return p.groups
	}

	return p.groups[1:]
}

func (p *parent[T]) hasDefaultCommandGroup() bool {
	return p.groups[0] != nil
}

func (p *parent[T]) FindChild(name string) argo.ChildNode[any] {
	for _, group := range p.groups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

func (p *parent[T]) HasSelectedChild() bool {
	return p.selectedChild != nil
}

func (p *parent[T]) SelectChild(name string) bool {
	for _, group := range p.groups {
		if child := group.FindChild(name); child != nil {
			p.selectedChild = child
			return true
		}
	}

	return false
}

func (p *parent[T]) SelectedChild() argo.ChildNode[any] {
	return p.selectedChild
}

func (p *parent[T]) HasIncompleteHandler() bool {
	return p.incompleteFn != nil
}

func (p *parent[T]) IncompleteHandler() argo.IncompleteCommandHandler[argo.CommandTree] {
	return p.incompleteFn
}
