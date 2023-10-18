package argo

import (
	"io"
)

// A HelpRenderer is a type that renders help text for the given type to the
// given io.Writer instance.
//
// This interface may be implemented to provide custom help rendering for your
// command.
type HelpRenderer[T any] interface {

	// RenderHelp renders help text for the given command and writes it to the
	// given io.Writer instance.
	RenderHelp(command T, writer io.Writer) error
}
