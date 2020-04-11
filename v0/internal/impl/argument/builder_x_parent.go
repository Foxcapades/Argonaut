package argument

import (
	"errors"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func (a *Builder) ValidateParent() error {
	if a.ParentElement == nil {
		// TODO: Make this a real error
		return errors.New("no parent set for argument")
	}

	if _, ok := a.ParentElement.(A.Flag); ok {
		return nil
	}

	if _, ok := a.ParentElement.(A.Command); ok {
		return nil
	}

	// TODO: make this a real error
	return errors.New("argument parent must be either Command or Flag")
}
