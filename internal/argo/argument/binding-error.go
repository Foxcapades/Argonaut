package argument

import "github.com/foxcapades/argonaut/v3/pkg/argo"

func NewBindingError(root error, builder argo.ArgumentBuilder) argo.ArgumentBindingError {
	return &argumentBindingError{root, builder}
}

type argumentBindingError struct {
	root    error
	builder argo.ArgumentBuilder
}

func (a argumentBindingError) Unwrap() error {
	return a.root
}

func (a argumentBindingError) Error() string {
	return "ArgumentBindingError: " + a.root.Error()
}

func (a argumentBindingError) ArgumentBuilder() argo.ArgumentBuilder {
	return a.builder
}
