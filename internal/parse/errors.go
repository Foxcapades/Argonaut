package parse

import (
	"fmt"
	"reflect"
)

const unmarshalErrFormat = "Format error in input text, could not unmarshal to type %s"

type formatError struct {
	Value    reflect.Value
	Argument string
	Kind     reflect.Kind
	Root     error
}

func (f formatError) Error() string {
	return fmt.Sprintf(unmarshalErrFormat, f.Kind)
}
