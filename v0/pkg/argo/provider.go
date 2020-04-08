package argo

type (
	ArgumentProvider  func(Provider) ArgumentBuilder
	CommandProvider   func(Provider) CommandBuilder
	FlagProvider      func(Provider) FlagBuilder
	FlagGroupProvider func(Provider) FlagGroupBuilder
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
