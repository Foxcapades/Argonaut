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

// RootValue returns the root of the given value.  The root will be the
// dereferenced type of the given reflect.Value instance.
//
// Parameters:
//   1. v  = reflect.Value of the type whose root value should be determined.
func RootValue(v reflect.Value) reflect.Value {
	// Used for recursion detection
	c := v

	haveAddr := false

	// see json.Unmarshaler indirect()
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}

	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			tmp := v.Elem()
			if tmp.Kind() == reflect.Ptr && !tmp.IsNil() {
				haveAddr = false
				v = tmp
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}

		if v.IsNil() {
			if v.CanAddr() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}

		if haveAddr {
			v = c
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}

	return v
}
