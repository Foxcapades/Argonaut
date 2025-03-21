package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type node[T any] struct {
	disableHelp bool
	description string
	flagGroups  []argo.FlagGroup
	callback    argo.CommandNodeCallback[T]
}

func (n node[T]) HasDescription() bool {
	return len(n.description) > 0
}

func (n node[T]) Description() string {
	return n.description
}

func (n node[T]) HasFlagGroups() bool {
	return len(n.flagGroups) > 0
}

func (n node[T]) FlagGroups() []argo.FlagGroup {
	return n.flagGroups
}

func (n node[T]) HasCallback() bool {
	return n.callback != nil
}

func (n node[T]) Callback() argo.CommandNodeCallback[T] {
	return n.callback
}

//

func (n node[T]) FindShortFlag(b byte) argo.Flag {
	for _, group := range n.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (n node[T]) FindLongFlag(name string) argo.Flag {
	for _, group := range n.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (n node[T]) IsHelpDisabled() bool {
	return n.disableHelp
}
