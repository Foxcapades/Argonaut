package argo

type described interface {

	// Description returns this element's description value
	// which is primarily used for rendering help text.
	//
	// If no description value has been set, this method will
	// return an empty string.
	Description() string

	// HasDescription returns whether or not this element has
	// a description set on it.
	HasDescription() bool
}
