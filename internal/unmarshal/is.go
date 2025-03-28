package unmarshal

import (
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/log"
	"github.com/foxcapades/argonaut/v3/internal/xreflect"
)

// IsUnmarshalable tests the given reflect.Type value to see if it is something
// that Argonaut can unmarshal.
//
// Unmarshalable values include pointers to basic types, maps, slices, or
// consumer functions.
//
// Arguments:
//  1. vt = Type of the value that we are testing to ensure that it is
//     unmarshalable.
//  2. ut = Unmarshaler type.  This is a hacky way of getting around cyclic
//     package imports with the argo package.
func IsUnmarshalable(vt reflect.Type) (out bool) {
	defer func() {
		if err := recover(); err != nil {
			out = false
		}
	}()

	log.DebugLn1("IsUnmarshalable(vt = %s)", vt)

	// If it's a pointer, then check that it's something that can actually validly
	// be a pointer.
	if vt.Kind() == reflect.Ptr {
		if xreflect.IsBasicPointer(vt) || IsUnmarshaler(vt) || IsUnmarshaler(vt.Elem()) {
			log.DebugLn("  bail early for easy kind check")
			return true
		}

		vt = xreflect.RootType(vt)
		log.DebugLn1("  root type is %s", vt)
	}

	// If it's not a pointer, then maybe it's an Unmarshaler instance.
	if IsUnmarshaler(vt) {
		log.DebugLn("  root type is unmarshalable")
		return true
	}

	// If it's not a pointer or unmarshaler, maybe it's a consumer func
	if IsConsumerFunc(vt) {
		log.DebugLn("  root type is consumer func")
		return true
	}

	if xreflect.IsBasicMap(vt) {
		log.DebugLn("  root type is is basic map")
		return true
	}

	if xreflect.IsBasicSlice(vt) {
		log.DebugLn("  root type is basic slice")
		return true
	}

	if IsUnmarshalerMap(vt) {
		log.DebugLn("  root type is unmarshaler map")
		return true
	}

	if IsUnmarshalerSlice(vt) {
		log.DebugLn("  root type is unmarshaler slice")
		return true
	}

	if IsUnmarshalerSliceMap(vt) {
		log.DebugLn("  root type is unmarshaler slice map")
		return true
	}

	log.DebugLn("  !!root type is invalid")
	return false
}

func IsUnmarshaler(t reflect.Type) bool {
	return t.AssignableTo(UnmarshalerType)
}

func IsUnmarshalerMap(t reflect.Type) bool {
	return t.Kind() == reflect.Map &&
		xreflect.IsBasic(t.Key()) &&
		IsUnmarshaler(t.Elem())
}

func IsUnmarshalerSliceMap(t reflect.Type) bool {
	return t.Kind() == reflect.Map &&
		xreflect.IsBasic(t.Key()) &&
		t.Elem().Kind() == reflect.Slice &&
		IsUnmarshaler(t.Elem().Elem())
}

func IsUnmarshalerSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		IsUnmarshaler(t.Elem())
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
	if xreflect.IsBasic(rt) || xreflect.IsBasicSlice(rt) {
		return true
	}

	// If it's not basic or a slice, maybe it's a map
	if xreflect.IsBasicMap(rt) {
		return true
	}

	return false
}
