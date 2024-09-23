package argo

type FlagSpec interface {
	WithLongForm(name string) FlagSpec

	WithShortForm(name byte) FlagSpec

	WithArgument(arg ArgumentSpecBuilder) FlagSpec
}
