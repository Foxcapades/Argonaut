package argo

import (
	"fmt"
	R "reflect"
)

type ArgumentErrorType uint8

const (
	ArgErrInvalidDefault ArgumentErrorType = 1 << iota
	ArgErrInvalidBinding

	ArgErrInvalidDefaultFn  = ArgErrInvalidDefault | 1<<iota
	ArgErrInvalidDefaultVal = ArgErrInvalidDefault | 1<<iota

	ArgErrInvalidBindingNil     = ArgErrInvalidBinding | 1<<iota
	ArgErrInvalidBindingBadType = ArgErrInvalidBinding | 1<<iota
)

type ArgumentError interface {
	error
	Error() string
	Type() ArgumentErrorType
	Is(errorType ArgumentErrorType) bool
	Builder() ArgumentBuilder
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Invalid Arg Config      ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// Error message formats
const (
	errArgBinding = `Argument bound to type "%s" which cannot be unmarshaled.`
	errArgDefault = "Argument default type (%s) must be compatible with binding" +
		" type (%s)"
	errArgDefFn = "Invalid default provider func (%s): %s"
)

func NewInvalidArgError(
	kind ArgumentErrorType,
	build ArgumentBuilder,
	reason string,
) error {
	return &invalidArgError{eType: kind, build: build, reason: reason}
}

type invalidArgError struct {
	eType  ArgumentErrorType
	build  ArgumentBuilder
	reason string
}

func (i *invalidArgError) Type() ArgumentErrorType {
	return i.eType
}

func (i *invalidArgError) Is(kind ArgumentErrorType) bool {
	return i.eType&kind == kind
}

func (i *invalidArgError) Builder() ArgumentBuilder {
	return i.build
}

func (i *invalidArgError) Error() string {
	switch i.eType {
	case ArgErrInvalidBindingBadType, ArgErrInvalidBindingNil:
		return fmt.Sprintf(errArgBinding, R.TypeOf(i.build.GetBinding()))
	case ArgErrInvalidDefaultVal:
		return fmt.Sprintf(errArgDefault, R.TypeOf(i.build.GetDefault()),
			R.TypeOf(i.build.GetBinding()))
	case ArgErrInvalidDefaultFn:
		return fmt.Sprintf(errArgDefFn, R.TypeOf(i.build.GetDefault()), i.reason)
	}
	panic(fmt.Errorf("invalid argument error type %d", i.eType))
}
