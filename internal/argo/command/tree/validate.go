package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/text"
)

func ValidateNodeName(name string) error {
	if len(name) == 0 {
		return errors.New("command names must not be blank")
	}

	if !(text.IsWord(name[0])) {
		return errors.New("command names must begin with an alphanumeric character or an underscore")
	}

	for i := 1; i < len(name); i++ {
		if !text.IsWord(name[i]) || name[i] == text.DashByte {
			return errors.New("command names may only contain alphanumeric characters, dashes, and/or underscores")
		}
	}

	return nil
}
