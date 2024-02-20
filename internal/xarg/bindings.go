package xarg

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Foxcapades/Argonaut/internal/unmarshal"
	"github.com/Foxcapades/Argonaut/internal/xreflect"
)

type BindKind uint8

const (
	BindKindNone BindKind = iota
	BindKindPointer
	BindKindUnmarshaler
	BindKindFuncPlain
	BindKindFuncWithErr
	BindKindUnknown = 254
	BindKindInvalid = 255
)

type DefaultKind uint8

const (
	DefaultKindNone DefaultKind = iota
	DefaultKindRaw
	DefaultKindParsed
	DefaultKindProviderPlain
	DefaultKindProviderWithErr
	DefaultKindUnknown = 254
	DefaultKindInvalid = 255
)

func DetermineBindKind(bind any, ut reflect.Type) (kind BindKind, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			kind = BindKindInvalid
			err = errors.New("binding must be a pointer, a consumer func, or an argo.Unmarshaler instance")
		}
	}()

	rt := reflect.TypeOf(bind)
	rk := rt.Kind()

	switch rk {
	case reflect.Ptr:
		if unmarshal.IsUnmarshalable(rt, ut) {
			if rt.Elem().AssignableTo(ut) {
				return BindKindUnmarshaler, nil
			}

			return BindKindPointer, nil
		}

		return BindKindInvalid, errors.New("binding is a pointer to a type that cannot be unmarshalled")

	case reflect.Func:
		if unmarshal.IsUnmarshalable(rt, ut) {
			if xreflect.FuncHasReturn(rt) {
				return BindKindFuncWithErr, nil
			}

			return BindKindFuncPlain, nil
		}

		return BindKindInvalid, fmt.Errorf("binding is invalid function type %s", rt)

	default:
		return BindKindInvalid, fmt.Errorf("invalid binding kind: %s", rk)
	}
}

func DetermineDefaultKind(bind, def any) (kind DefaultKind, err error) {
	// If we encounter a panic when using the "reflect" package in this function,
	// then the most likely scenario is that the default is a zero-valued type.
	defer func() {
		if rec := recover(); rec != nil {
			kind = DefaultKindInvalid
			err = errors.New("default must be a value, a raw string, a provider of the same time as the binding, or a value of the same type as the binding")
		}
	}()

	bt := xreflect.RootType(reflect.TypeOf(bind))
	dt := reflect.TypeOf(def)

	// If the binding type is a consumer function then we want to test against
	// that function's param rather than the function type itself as nothing will
	// match against that.
	if xreflect.IsFunction(bt) {
		bt = xreflect.GetFunctionParamType(bt, 0)
	}

	// If the provided default value is of a basic type
	if xreflect.IsBasic(dt) {
		// If the default value is a string BUT the binding type IS NOT a string,
		// then we have a raw value binding which needs to be parsed before it may
		// be set to the binding pointer.
		if xreflect.IsString(dt) && !xreflect.IsString(bt) {
			kind = DefaultKindRaw
			err = nil
			return
		}

		// If the default value is not a string and is not assignable to the binding
		// type, then we are in an invalid state because we will get a panic if we
		// attempt to set the default value to the binding pointer.
		if !dt.AssignableTo(bt) {
			kind = DefaultKindInvalid
			err = fmt.Errorf("expected default value of type %s but got %s instead", bt.Kind(), dt.Kind())
			return
		}

		// If we've made it here, then the default type is compatible with the
		// binding type.
		kind = DefaultKindParsed
		err = nil
		return
	}

	// If the default value is a slice
	if xreflect.IsBasicSlice(dt) {
		// And it is assignable to the binding value
		if dt.AssignableTo(bt) {
			kind = DefaultKindParsed
			err = nil
			return
		}

		// If it is not assignable then it is invalid.
		kind = DefaultKindInvalid

		// Try and report a helpful error about the situation
		if xreflect.IsSlice(bt) {
			err = fmt.Errorf("expected default value to be a slice of type %s but got a slice of type %s instead", bt.Elem().Kind(), dt.Elem().Kind())
		} else {
			err = fmt.Errorf("expected default value to be of type %s but got a slice of type %s instead", bt.Kind(), dt.Elem().Kind())
		}

		return
	}

	// If the default type is a basic map
	if xreflect.IsBasicMap(dt) {
		// And the binding type is also a basic map, then we can stop here as we
		// know they are compatible.
		if dt.AssignableTo(bt) {
			kind = DefaultKindParsed
			err = nil
			return
		}

		// If the binding type is not a matching basic map, then we fail.
		kind = DefaultKindInvalid
		err = fmt.Errorf("expected default value to be of type %s but got a value of type %s instead", bt, dt)
	}

	// If the default value resembles a provider function
	if ResemblesProviderFunction(dt) {
		// But the first out param is not compatible with the binding type
		if !dt.Out(0).AssignableTo(bt) {
			kind = DefaultKindInvalid
			err = fmt.Errorf("default value provider does not returns type %d which is incompatible with binding type %s", dt, bt)
			return
		}

		// But the second out param is not compatible with error
		if dt.NumOut() == 2 {
			if !dt.Out(1).AssignableTo(xreflect.ErrorType) {
				kind = DefaultKindInvalid
				err = fmt.Errorf("default value provider does not return error as its second return value")
				return
			}

			kind = DefaultKindProviderWithErr
			err = nil
			return
		}

		kind = DefaultKindProviderPlain
		err = nil
		return
	}

	kind = DefaultKindInvalid
	err = fmt.Errorf("invalid default value type %s", dt)
	return
}

func ResemblesProviderFunction(dt reflect.Type) bool {
	if dt.NumIn() != 0 {
		return false
	}

	switch dt.NumOut() {
	case 1:
		return true
	case 2:
		return true
	default:
		return false
	}
}
