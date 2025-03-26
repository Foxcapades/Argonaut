package argument

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/xreflect"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewDefault(value any) Default {
	return Default{value, argo.DefaultTypeUnknown}
}

type Default struct {
	value any
	dType argo.ArgumentDefaultType
}

func (d *Default) Type() argo.ArgumentDefaultType {
	return d.dType
}

func (d *Default) Value() any {
	return d.value
}

func (d *Default) IsUsable() bool {
	// This method is only called internally by the default implementations, so we
	// know that the "Unknown" type isn't a possible value.
	return !(d.dType == argo.DefaultTypeNone || d.dType == argo.DefaultTypeInvalid)
}

func DetermineDefaultType(bind, def any) (kind argo.ArgumentDefaultType, err error) {
	// If we encounter a panic when using the "reflect" package in this function,
	// then the most likely scenario is that the default is a zero-valued type.
	defer func() {
		if rec := recover(); rec != nil {
			kind = argo.DefaultTypeInvalid
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
			kind = argo.DefaultTypeRaw
			err = nil
			return
		}

		// If the default value is not a string and is not assignable to the binding
		// type, then we are in an invalid state because we will get a panic if we
		// attempt to set the default value to the binding pointer.
		if !dt.AssignableTo(bt) {
			kind = argo.DefaultTypeInvalid
			err = fmt.Errorf("expected default value of type %s but got %s instead", bt.Kind(), dt.Kind())
			return
		}

		// If we've made it here, then the default type is compatible with the
		// binding type.
		kind = argo.DefaultTypeParsed
		err = nil
		return
	}

	// If the default value is a slice
	if xreflect.IsBasicSlice(dt) {
		// And it is assignable to the binding value
		if dt.AssignableTo(bt) {
			kind = argo.DefaultTypeParsed
			err = nil
			return
		}

		// If it is not assignable then it is invalid.
		kind = argo.DefaultTypeInvalid

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
			kind = argo.DefaultTypeParsed
			err = nil
			return
		}

		// If the binding type is not a matching basic map, then we fail.
		kind = argo.DefaultTypeInvalid
		err = fmt.Errorf("expected default value to be of type %s but got a value of type %s instead", bt, dt)
	}

	// If the default value resembles a provider function
	if ResemblesProviderFunction(dt) {
		// But the first out param is not compatible with the binding type
		if !dt.Out(0).AssignableTo(bt) {
			kind = argo.DefaultTypeInvalid
			err = fmt.Errorf("default value provider does not returns type %d which is incompatible with binding type %s", dt, bt)
			return
		}

		// But the second out param is not compatible with error
		if dt.NumOut() == 2 {
			if !dt.Out(1).AssignableTo(xreflect.ErrorType) {
				kind = argo.DefaultTypeInvalid
				err = fmt.Errorf("default value provider does not return error as its second return value")
				return
			}

			kind = argo.DefaultTypeProviderWithErr
			err = nil
			return
		}

		kind = argo.DefaultTypeProviderPlain
		err = nil
		return
	}

	kind = argo.DefaultTypeInvalid
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
