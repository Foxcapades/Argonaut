package argo

import (
	"reflect"
)

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

func reflectGetRootValue(v reflect.Value) reflect.Value {
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
			v.Set(reflect.New(v.Type().Elem()))
		}

		if v.Type().AssignableTo(unmarshalerType) {
			break
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

func reflectIsUnmarshaler(t reflect.Type) bool {
	return t.AssignableTo(unmarshalerType)
}

func reflectIsInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}

func reflectIsBasicKind(k reflect.Kind) bool {
	return k == reflect.String || k == reflect.Bool || reflectIsNumericKind(k)
}

func reflectIsNumericKind(k reflect.Kind) bool {
	return numericKinds[k]
}

func reflectIsByteSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		t.Elem().Kind() == reflect.Uint8
}

func reflectCompatible(val, test *reflect.Value) bool {
	return val.Type().AssignableTo(test.Type())
}