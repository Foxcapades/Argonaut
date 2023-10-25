package unmarshal

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Foxcapades/Argonaut/internal/xreflect"
)

type InvalidTypeError struct {
	Value reflect.Value
}

func (i *InvalidTypeError) Error() string {
	return fmt.Sprintf("Cannot unmarshal type %s", i.Value.Type())
}

type InvalidUnmarshalError struct {
	Value    reflect.Value
	Argument string
}

func (i InvalidUnmarshalError) Error() string {
	if xreflect.IsNil(&i.Value) {
		return "Attempted to unmarshal into nil"
	}
	return "Attempted to unmarshal into a non-pointer"
}

func ToUnmarshalable(
	arg string,
	ov reflect.Value,
	skipPtr bool,
	unmarshalerType reflect.Type,
) (reflect.Value, error) {

	if !skipPtr && ((ov.Kind() != reflect.Ptr && ov.Kind() != reflect.Func) || xreflect.IsNil(&ov)) {
		return reflect.Value{}, &InvalidUnmarshalError{Value: ov, Argument: arg}
	}

	v := GetRootValue(ov, unmarshalerType)

	kind := v.Kind()

	if xreflect.IsBasicKind(kind) {
		return v, nil
	}

	if v.Type().AssignableTo(unmarshalerType) {
		return v, nil
	}

	if kind == reflect.Slice {
		return ToValidSlice(v, ov, unmarshalerType)
	}
	if kind == reflect.Map {
		return ToValidMap(v, ov, unmarshalerType)
	}
	if kind == reflect.Struct && ov.Type().AssignableTo(reflect.TypeOf((*time.Time)(nil))) {
		return v, nil
	}
	if kind == reflect.Interface {
		return v, nil
	}

	return reflect.Value{}, &InvalidTypeError{Value: v}
}

// ToValidSlice is an internal method that is not exposed to package consumers.
//
// Valid slice types:
//   []<basic>
//   []<*basic>
//   [][]byte
//   []*[]byte
func ToValidSlice(v, ov reflect.Value, ut reflect.Type) (out reflect.Value, err error) {
	err = ValidateContainerValue(v.Type().Elem(), ov, ut)
	if err != nil {
		return
	}
	return v, nil
}

// ToValidMap is butts and this comment line is meaningless.
//
// Valid map types:
//   map[<basic>]<basic>
//   map[<basic>]<*basic>
//   map[<basic>][]byte
//   map[<basic>]<*[]byte>
func ToValidMap(v, ov reflect.Value, ut reflect.Type) (reflect.Value, error) {
	vt := v.Type()

	if !xreflect.IsBasicKind(vt.Key().Kind()) {
		return reflect.Value{}, &InvalidTypeError{Value: ov}
	}

	if err := ValidateContainerValue(vt.Elem(), ov, ut); err != nil {
		return reflect.Value{}, err
	}

	return v, nil
}

func ValidateContainerValue(t reflect.Type, ov reflect.Value, ut reflect.Type) error {
	sk := t.Kind()

	if xreflect.IsBasicKind(sk) {
		return nil
	}

	if sk == reflect.Ptr {
		if xreflect.IsBasicKind(t.Elem().Kind()) {
			return nil
		}
		if xreflect.IsByteSlice(t.Elem()) {
			return nil
		}
		if xreflect.IsUnmarshaler(t.Elem(), ut) {
			return nil
		}
		if xreflect.IsInterface(t.Elem()) {
			return nil
		}

		return &InvalidTypeError{Value: ov}
	}

	if xreflect.IsByteSlice(t) {
		return nil
	}
	if xreflect.IsBasicSlice(t) {
		return nil
	}
	if xreflect.IsUnmarshaler(t, ut) {
		return nil
	}
	if xreflect.IsInterface(t) {
		return nil
	}

	return &InvalidTypeError{Value: ov}
}

// GetRootValue returns the root of the given value.  The root will be the
// dereferenced type of the given reflect.Value instance.
//
// Parameters:
//   1. v  = reflect.Value of the type whose root value should be determined.
//   2. ut = unmarshaler type to test if the given value is an argo.Unmarshaler
//           instance.
func GetRootValue(v reflect.Value, ut reflect.Type) reflect.Value {
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

		if v.Type().AssignableTo(ut) {
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
