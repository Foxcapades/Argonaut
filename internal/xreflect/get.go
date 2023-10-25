package xreflect

import "reflect"

func GetFunctionParamType(t reflect.Type, index uint8) reflect.Type {
	return t.In(int(index))
}

// RootType returns the root of the given type, unpacking nested pointer series
// if necessary.
func RootType(vt reflect.Type) reflect.Type {
	if vt.Kind() != reflect.Ptr {
		return vt
	}

	current := vt
	for current.Kind() == reflect.Ptr {
		current = current.Elem()
	}

	return current
}
