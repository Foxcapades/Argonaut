package util

import (
	R "reflect"

	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

var numericKinds = map[R.Kind]bool{
	R.Int:     true,
	R.Int8:    true,
	R.Int16:   true,
	R.Int32:   true,
	R.Int64:   true,
	R.Uint:    true,
	R.Uint8:   true,
	R.Uint16:  true,
	R.Uint32:  true,
	R.Uint64:  true,
	R.Float32: true,
	R.Float64: true,
}

func GetRootValue(v R.Value) R.Value {
	// Used for recursion detection
	c := v

	haveAddr := false

	// see json.Unmarshaler indirect()
	if v.Kind() != R.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}

	for {
		if v.Kind() == R.Interface && !v.IsNil() {
			tmp := v.Elem()
			if tmp.Kind() == R.Ptr && !tmp.IsNil() {
				haveAddr = false
				v = tmp
				continue
			}
		}

		if v.Kind() != R.Ptr {
			break
		}

		if v.Elem().Kind() == R.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}

		if v.IsNil() {
			v.Set(R.New(v.Type().Elem()))
		}

		if v.Type().AssignableTo(R.TypeOf((*argo.Unmarshaler)(nil)).Elem()) {
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

func IsBasicKind(k R.Kind) bool {
	return k == R.String || k == R.Bool || IsNumericKind(k)
}

func IsNumericKind(k R.Kind) bool {
	return numericKinds[k]
}

func IsByteSlice(t R.Type) bool {
	return t.Kind() == R.Slice &&
		t.Elem().Kind() == R.Uint8
}

func Compatible(val, test interface{}) bool {
	vt := GetRootValue(R.ValueOf(val)).Type()
	tt := R.TypeOf(test)

	if tt.Kind() == R.Ptr {
		return tt.Elem().AssignableTo(vt)
	}

	return tt.AssignableTo(vt)
}
