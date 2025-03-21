package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func newNodeBuilder[T any](root T) nodeBuilder[T] {
	return nodeBuilder[T]{root: root}
}

type nodeBuilder[T any] struct {
	root T

	disableHelp bool
	description string
	flagGroups  []argo.FlagGroupBuilder
	callback    argo.CommandNodeCallback[T]
}

//

func (n *nodeBuilder[T]) WithDescription(description string) T {
	n.description = description
	return n.root
}

func (n *nodeBuilder[T]) HasDescription() bool {
	return len(n.description) > 0
}

func (n *nodeBuilder[T]) Description() string {
	return n.description
}

//

func (n *nodeBuilder[T]) WithFlagGroup(flagGroup argo.FlagGroupBuilder) T {
	n.flagGroups = append(n.flagGroups, flagGroup)
	return n.root
}

func (n *nodeBuilder[T]) WithFlagGroups(flagGroups ...argo.FlagGroupBuilder) T {
	n.flagGroups = append(n.flagGroups, flagGroups...)
	return n.root
}

func (n *nodeBuilder[T]) HasFlagGroups(includeDefault bool) bool {
	return len(n.flagGroups) > 1 || (includeDefault && n.hasDefaultFlagGroup())
}

func (n *nodeBuilder[T]) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if includeDefault && n.hasDefaultFlagGroup() {
		return n.flagGroups
	}

	return n.flagGroups[1:]
}

func (n *nodeBuilder[T]) hasDefaultFlagGroup() bool {
	return n.flagGroups[0] != nil && n.flagGroups[0].Size() > 0
}

func (n *nodeBuilder[T]) initDefaultFlagGroup() {
	if n.flagGroups[0] == nil {
		n.flagGroups[0] = flag.NewGroupBuilder("Ungrouped") // TODO: make this configurable
	}
}

//

func (n *nodeBuilder[T]) WithFlag(flag argo.FlagBuilder) T {
	n.initDefaultFlagGroup()
	n.flagGroups[0].WithFlag(flag)
	return n.root
}

func (n *nodeBuilder[T]) WithFlags(flags ...argo.FlagBuilder) T {
	n.initDefaultFlagGroup()
	n.flagGroups[0].WithFlags(flags...)
	return n.root
}

func (n *nodeBuilder[T]) HasFlags() bool {
	for _, group := range n.flagGroups {
		if group != nil && group.HasFlags() {
			return true
		}
	}

	return false
}

//

func (n *nodeBuilder[T]) WithCallback(callback argo.CommandNodeCallback[T]) T {
	n.callback = callback
	return n.root
}

func (n *nodeBuilder[T]) HasCallback() bool {
	return n.callback != nil
}

func (n *nodeBuilder[T]) Callback() argo.CommandNodeCallback[T] {
	return n.callback
}

//

func (n *nodeBuilder[T]) WithHelpDisabled() T {
	n.disableHelp = true
	return n.root
}

func (n *nodeBuilder[T]) IsHelpDisabled() bool {
	return n.disableHelp
}
