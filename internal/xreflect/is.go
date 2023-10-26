package xreflect

import (
	"reflect"
)

var ErrorType = reflect.TypeOf((*error)(nil)).Elem()

var numericKinds = map[reflect.Kind]bool{
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Float32: true,
	reflect.Float64: true,
}

func IsPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}

func IsString(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

func IsNumeric(t reflect.Type) bool {
	return IsNumericKind(t.Kind())
}

func IsUnmarshaler(t reflect.Type, ut reflect.Type) bool {
	return t.AssignableTo(ut)
}

func IsInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}

func IsFunction(t reflect.Type) bool {
	return t.Kind() == reflect.Func
}

func IsSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

func IsBasic(t reflect.Type) bool {
	return IsBasicKind(t.Kind())
}

func IsBasicKind(k reflect.Kind) bool {
	return k == reflect.String || k == reflect.Bool || IsNumericKind(k)
}

func IsNumericKind(k reflect.Kind) bool {
	return numericKinds[k]
}

func IsByteSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		t.Elem().Kind() == reflect.Uint8
}

func IsBasicSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		IsBasicKind(t.Elem().Kind())
}

func IsUnmarshalerSlice(t, ut reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		IsUnmarshaler(t.Elem(), ut)
}

func FuncHasReturn(t reflect.Type) bool {
	return t.NumOut() > 0
}

// IsBasicPointer tests whether the given type is a pointer to a basic built-in
// type.
func IsBasicPointer(vt reflect.Type) bool {
	if vt.Kind() != reflect.Ptr {
		return false
	}

	return IsBasicKind(RootType(vt).Kind())
}

// IsBasicMap tests whether the given type represents a map with a key of a
// basic built-in type to a value of one of a basic type, a pointer to a basic
// type, or a slice of basic type values.
func IsBasicMap(t reflect.Type) bool {
	if t.Kind() != reflect.Map {
		return false
	}

	if !IsBasicKind(t.Key().Kind()) {
		return false
	}

	vt := t.Elem()

	return IsBasicKind(vt.Kind()) || IsBasicSlice(vt) || IsBasicPointer(t)
}

func IsNil(v *reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
