package argument

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Container interface {
	HasArguments() bool
	Arguments() []argo.Argument
	HasUnmappedInputLabel() bool
	UnmappedInputLabel() string
	HasUnmappedInputs() bool
	UnmappedInputs() []string
}
