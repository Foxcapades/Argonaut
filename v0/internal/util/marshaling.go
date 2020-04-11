package util

import (
	R "reflect"

	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

var unmarshalerType = R.TypeOf((*A.Unmarshaler)(nil)).Elem()

func IsUnmarshalable(o interface{}) bool {
	_, err := ToUnmarshalable("", R.ValueOf(o), false)
	return err == nil
}

func ToUnmarshalable(arg string, ov R.Value, skipPtr bool) (R.Value, error) {

	if (!skipPtr && ov.Kind() != R.Ptr) || IsNil(&ov) {
		return R.Value{}, &A.InvalidUnmarshalError{Value: ov, Argument: arg}
	}

	v := GetRootValue(ov)

	kind := v.Kind()

	if IsBasicKind(kind) {
		return v, nil
	}

	if v.Type().AssignableTo(unmarshalerType) {
		return v, nil
	}

	if kind == R.Slice {
		return toValidSlice(v, ov)
	}
	if kind == R.Map {
		return toValidMap(v, ov)
	}
	if kind == R.Interface {
		return v, nil
	}

	return R.Value{}, &A.InvalidTypeError{Value: v}
}

// Valid slice types:
//   []<basic>
//   []<*basic>
//   [][]byte
//   []*[]byte
func toValidSlice(v, ov R.Value) (out R.Value, err error) {
	err = validateContainerValue(v.Type().Elem(), ov)
	if err != nil {
		return
	}
	return v, nil
}

// Valid map types:
//   map[<basic>]<basic>
//   map[<basic>]<*basic>
//   map[<basic>][]byte
//   map[<basic>]<*[]byte>
func toValidMap(v, ov R.Value) (R.Value, error) {
	vt := v.Type()

	if !IsBasicKind(vt.Key().Kind()) {
		return R.Value{}, &A.InvalidTypeError{Value: ov}
	}

	if err := validateContainerValue(vt.Elem(), ov); err != nil {
		return R.Value{}, err
	}

	return v, nil
}

func validateContainerValue(t R.Type, ov R.Value) error {
	sk := t.Kind()

	if IsBasicKind(sk) {
		return nil
	}

	if sk == R.Ptr {
		if IsBasicKind(t.Elem().Kind()) {
			return nil
		}
		if IsByteSlice(t.Elem()) {
			return nil
		}

		return &A.InvalidTypeError{Value: ov}
	}

	if IsByteSlice(t) {
		return nil
	}

	return &A.InvalidTypeError{Value: ov}
}

func IsNil(v *R.Value) bool {
	switch v.Kind() {
	case R.Ptr, R.Chan, R.Func, R.Interface, R.Map, R.Slice:
		return v.IsNil()
	default:
		return false
	}
}
