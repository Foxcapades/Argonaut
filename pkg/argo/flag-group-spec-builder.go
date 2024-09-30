package argo

// FlagGroupSpecBuilder defines a builder type which is used to configure and
// construct FlagGroupSpec instances.
type FlagGroupSpecBuilder interface {
	// WithName configures the name for the output flag group.
	WithName(name string) FlagGroupSpecBuilder

	// WithDescription configures the help text description of the flag group.
	WithDescription(description string) FlagGroupSpecBuilder

	// WithFlag appends a flag to the flag group.
	WithFlag(flag FlagSpecBuilder) FlagGroupSpecBuilder

	// Build attempts to build a new FlagGroupSpec instance from the current
	// configured state of the parent FlagGroupSpecBuilder.
	Build(config Config) (FlagGroupSpec, error)
}
