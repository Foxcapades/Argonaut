package util

import (
	"reflect"

	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func ToUnmarshalable(arg string, ov reflect.Value) (reflect.Value, error) {

	if ov.Kind() != reflect.Ptr || ov.IsNil() {
		return reflect.Value{}, &argo.InvalidUnmarshalError{Value: ov, Argument: arg}
	}

	v := GetRootValue(ov)

	kind := v.Kind()

	if IsBasicKind(kind) {
		return v, nil
	}

	if kind == reflect.Slice {
		return toValidSlice(v, ov)
	}
	if kind == reflect.Map {
		return toValidMap(v, ov)
	}
	if kind == reflect.Interface {
		return v, nil
	}

	return reflect.Value{}, &argo.InvalidTypeError{Value: v}
}

// Valid slice types:
//   []<basic>
//   []<*basic>
//   [][]byte
//   []*[]byte
func toValidSlice(v, ov reflect.Value) (out reflect.Value, err error) {
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
func toValidMap(v, ov reflect.Value) (reflect.Value, error) {
	vt := v.Type()

	if !IsBasicKind(vt.Key().Kind()) {
		return reflect.Value{}, &argo.InvalidTypeError{Value: ov}
	}

	if err := validateContainerValue(vt.Elem(), ov); err != nil {
		return reflect.Value{}, err
	}

	return v, nil
}

func validateContainerValue(t reflect.Type, ov reflect.Value) error {
	sk := t.Kind()

	if IsBasicKind(sk) {
		return nil
	}

	if sk == reflect.Ptr {
		if IsBasicKind(t.Elem().Kind()) {
			return nil
		}
		if IsByteSlice(t.Elem()) {
			return nil
		}

		return &argo.InvalidTypeError{Value: ov}
	}

	if IsByteSlice(t) {
		return nil
	}

	return &argo.InvalidTypeError{Value: ov}
}
