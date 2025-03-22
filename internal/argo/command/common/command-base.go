package common

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type CommandBase[T any] struct {
	disableHelp bool
	description string
	flagGroups  []argo.FlagGroup
	callback    argo.CommandCallback[T]
}

func (c *CommandBase[T]) HasDescription() bool {
	return len(c.description) > 0
}

func (c *CommandBase[T]) Description() string {
	return c.description
}

func (c *CommandBase[T]) HasFlagGroups(includeDefault bool) bool {

}

func (c *CommandBase[T]) FlagGroups(includeDefault bool) []argo.FlagGroup {

}

func (c *CommandBase[T]) HasCallback() bool {
	return c.callback != nil
}

func (c *CommandBase[T]) Callback() argo.CommandCallback[T] {
	return c.callback
}

func (c *CommandBase[T]) FindShortFlag(b byte) argo.Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (c *CommandBase[T]) FindLongFlag(name string) argo.Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (c *CommandBase[T]) IsHelpDisabled() bool {
	return c.disableHelp
}
