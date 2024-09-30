package argo

type CommandLeafSpecBuilder interface {
	SubCommandSpecBuilder[CommandLeafSpecBuilder]

	WithPositionalArgument(arg ArgumentSpecBuilder) CommandLeafSpecBuilder

	WithUnmappedInputLabel(label string) CommandLeafSpecBuilder

	WithUnmappedInputDescription(description string) CommandLeafSpecBuilder

	WithPassthroughsDisabled() CommandLeafSpecBuilder

	Build(config Config) (CommandLeafSpec, error)
}
