package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/text"
)

func ValidateNodeName(name string) error {
	if len(name) == 0 {
		return errors.New("command names must not be blank")
	}

	if !(text.IsAlphanumeric(name[0]) || name[0] == '_') {
		return errors.New("command names must begin with an alphanumeric character or an underscore")
	}

	for i := 1; i < len(name); i++ {
		if !text.IsFlagStringSafe(name[i]) {
			return errors.New("command names may only contain alphanumeric characters, dashes, and/or underscores")
		}
	}

	return nil
}
