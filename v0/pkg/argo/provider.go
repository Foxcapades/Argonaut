package argo

type (
	ArgumentProvider  func() ArgumentBuilder
	CommandProvider   func() CommandBuilder
	FlagProvider      func() FlagBuilder
	FlagGroupProvider func() FlagGroupBuilder
)

type Provider interface {
	NewArg() ArgumentBuilder
	NewCommand() CommandBuilder
	NewFlag() FlagBuilder
	NewFlagGroup() FlagGroupBuilder

	ArgumentProvider(provider ArgumentProvider) (this Provider)
	CommandProvider(provider CommandProvider) (this Provider)
	FlagProvider(provider FlagProvider) (this Provider)
	FlagGroupProvider(provider FlagGroupProvider) (this Provider)
}
