package argo

type ArgumentSpecBuilder interface {
	// Build validates the state of this ArgumentSpecBuilder instance and produces
	// a new ArgumentSpec based on the options configured on the builder.
	//
	// If the ArgumentSpecBuilder has been configured incorrectly, an error will
	// be returned.
	Build(config Config) (ArgumentSpec, error)
}
