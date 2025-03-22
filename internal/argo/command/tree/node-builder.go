package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func newNodeBuilder[T, O any](root T) nodeBuilder[T, O] {
	return nodeBuilder[T, O]{root: root}
}

type nodeBuilder[T, O any] struct {
	root T

	disableHelp bool
	description string
	flagGroups  []argo.FlagGroupBuilder
	callback    argo.CommandCallback[O]
}

//

func (n *nodeBuilder[T, O]) WithDescription(description string) T {
	n.description = description
	return n.root
}

func (n *nodeBuilder[T, O]) HasDescription() bool {
	return len(n.description) > 0
}

func (n *nodeBuilder[T, O]) Description() string {
	return n.description
}

//

func (n *nodeBuilder[T, O]) WithFlagGroup(flagGroup argo.FlagGroupBuilder) T {
	n.flagGroups = append(n.flagGroups, flagGroup)
	return n.root
}

func (n *nodeBuilder[T, O]) WithFlagGroups(flagGroups ...argo.FlagGroupBuilder) T {
	n.flagGroups = append(n.flagGroups, flagGroups...)
	return n.root
}

func (n *nodeBuilder[T, O]) HasFlagGroups(includeDefault bool) bool {
	return (includeDefault && len(n.flagGroups) > 0) ||
		(n.hasDefaultFlagGroup() && len(n.flagGroups) > 1) ||
		len(n.flagGroups) > 0
}

func (n *nodeBuilder[T, O]) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if includeDefault && n.hasDefaultFlagGroup() {
		return n.flagGroups
	}

	return n.flagGroups[1:]
}

func (n *nodeBuilder[T, O]) hasDefaultFlagGroup() bool {
	return n.flagGroups[0] != nil && n.flagGroups[0].Size() > 0
}

//

func (n *nodeBuilder[T, O]) WithFlag(flag argo.FlagBuilder) T {
	n.flagGroups = common.EnsureDefaultFlagGroup(n.flagGroups)
	n.flagGroups[0].WithFlag(flag)
	return n.root
}

func (n *nodeBuilder[T, O]) WithFlags(flags ...argo.FlagBuilder) T {
	n.flagGroups = common.EnsureDefaultFlagGroup(n.flagGroups)
	n.flagGroups[0].WithFlags(flags...)
	return n.root
}

func (n *nodeBuilder[T, O]) HasFlags() bool {
	for _, group := range n.flagGroups {
		if group != nil && group.HasFlags() {
			return true
		}
	}

	return false
}

//

func (n *nodeBuilder[T, O]) WithCallback(callback argo.CommandCallback[O]) T {
	n.callback = callback
	return n.root
}

func (n *nodeBuilder[T, O]) HasCallback() bool {
	return n.callback != nil
}

func (n *nodeBuilder[T, O]) Callback() argo.CommandCallback[O] {
	return n.callback
}

//

func (n *nodeBuilder[T, O]) WithHelpDisabled() T {
	n.disableHelp = true
	return n.root
}

func (n *nodeBuilder[T, O]) IsHelpDisabled() bool {
	return n.disableHelp
}
