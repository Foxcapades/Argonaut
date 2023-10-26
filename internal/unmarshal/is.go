package unmarshal

import (
	"reflect"

	"github.com/Foxcapades/Argonaut/internal/xreflect"
)

// IsUnmarshalable tests the given reflect.Type value to see if it is something
// that Argonaut can unmarshal.
//
// Unmarshalable values include pointers to basic types, maps, slices, or
// consumer functions.
//
// Arguments:
//   1. vt = Type of the value that we are testing to ensure that it is
//      unmarshalable.
//   2. ut = Unmarshaler type.  This is a hacky way of getting around cyclic
//      package imports with the argo package.
func IsUnmarshalable(vt, ut reflect.Type) (out bool) {
	defer func() {
		if err := recover(); err != nil {
			out = false
		}
	}()

	// If it's a pointer, then check that it's something that can actually validly
	// be a pointer.
	if vt.Kind() == reflect.Ptr {
		if xreflect.IsBasicPointer(vt) || xreflect.IsUnmarshaler(vt, ut) || xreflect.IsUnmarshaler(vt.Elem(), ut) {
			return true
		}

		vt = xreflect.RootType(vt)
	}

	// If it's not a pointer, then maybe it's an Unmarshaler instance.
	if xreflect.IsUnmarshaler(vt, ut) {
		out = true
		return
	}

	// If it's not a pointer or unmarshaler, maybe it's a consumer func
	if IsConsumerFunc(vt) {
		out = true
		return
	}

	if xreflect.IsBasicMap(vt) {
		out = true
		return
	}

	if xreflect.IsBasicSlice(vt) {
		out = true
		return
	}

	if xreflect.IsUnmarshalerMap(vt, ut) {
		return true
	}

	if xreflect.IsUnmarshalerSlice(vt, ut) {
		out = true
		return
	}

	if xreflect.IsUnmarshalerSliceMap(vt, ut) {
		return true
	}

	out = false
	return
}

// IsConsumerFunc tests whether the given type represents a function that may be
// used as an argo.Argument's consumer binding.
func IsConsumerFunc(vt reflect.Type) bool {
	if vt.Kind() != reflect.Func {
		return false
	}

	if vt.NumIn() != 1 {
		return false
	}

	if !IsUnmarshalableValue(vt.In(0)) {
		return false
	}

	switch vt.NumOut() {
	case 0:
		return true
	case 1:
		return vt.Out(0).AssignableTo(xreflect.ErrorType)
	default:
		return false
	}
}

// IsUnmarshalableValue tests whether the given type is a value type that may
// be unmarshalled.
func IsUnmarshalableValue(vt reflect.Type) bool {
	rt := xreflect.RootType(vt)

	// If it's a basic built-in type, or a slice of basic built-in type values.
	if xreflect.IsBasicKind(rt.Kind()) || xreflect.IsBasicSlice(rt) {
		return true
	}

	// If it's not basic or a slice, maybe it's a map
	if xreflect.IsBasicMap(rt) {
		return true
	}

	return false
}
