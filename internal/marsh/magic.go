package marsh

import (
	"reflect"
	"strconv"
	"time"

	"github.com/Foxcapades/Argonaut/internal/xraw"
	"github.com/Foxcapades/Argonaut/internal/xref"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

var unmarshalerType = reflect.TypeOf((*argo.Unmarshaler)(nil)).Elem()

func NewDefaultedValueUnmarshaler() ValueUnmarshaler {
	tmp := argo.DefaultUnmarshalProps()
	return NewValueUnmarshaler(&tmp)
}

func NewValueUnmarshaler(props *argo.UnmarshalProps) ValueUnmarshaler {
	return valueUnmarshaler{props: props}
}

type valueUnmarshaler struct {
	props *argo.UnmarshalProps
}

func (v valueUnmarshaler) Unmarshal(raw string, val interface{}) (err error) {
	// If input string was empty, there's nothing to do
	if len(raw) == 0 {
		return nil
	}

	ptrVal := reflect.ValueOf(val)

	if ptrVal, err = ToUnmarshalable(raw, ptrVal, false); err != nil {
		return err
	}

	ptrDef := ptrVal.Type()

	if ptrDef.AssignableTo(unmarshalerType) {
		return ptrVal.Interface().(Unmarshaler).Unmarshal(raw)
	}

	switch ptrVal.Kind() {

	case reflect.String:
		ptrVal.SetString(raw)
	case reflect.Int:
		return unmarshalInt(ptrVal, raw, strconv.IntSize, &v.props.Integers)
	case reflect.Float32:
		return unmarshalFloat(ptrVal, raw, 32)
	case reflect.Int64:
		return unmarshalInt64(ptrVal, raw, &v.props.Integers)
	case reflect.Float64:
		return unmarshalFloat(ptrVal, raw, 64)
	case reflect.Uint64:
		return unmarshalUInt(ptrVal, raw, 64, &v.props.Integers)
	case reflect.Uint:
		return unmarshalUInt(ptrVal, raw, strconv.IntSize, &v.props.Integers)
	case reflect.Int32:
		return unmarshalInt(ptrVal, raw, 32, &v.props.Integers)
	case reflect.Uint32:
		return unmarshalUInt(ptrVal, raw, 32, &v.props.Integers)
	case reflect.Uint8:
		return unmarshalUInt(ptrVal, raw, 8, &v.props.Integers)
	case reflect.Slice:
		return v.unmarshalSlice(ptrVal, raw)
	case reflect.Map:
		return v.unmarshalMap(ptrVal, raw)
	case reflect.Int8:
		return unmarshalInt(ptrVal, raw, 8, &v.props.Integers)
	case reflect.Int16:
		return unmarshalInt(ptrVal, raw, 16, &v.props.Integers)
	case reflect.Uint16:
		return unmarshalUInt(ptrVal, raw, 16, &v.props.Integers)
	case reflect.Bool:
		return unmarshalBool(ptrVal, raw)
	case reflect.Interface:
		ptrVal.Set(reflect.ValueOf(raw))
	case reflect.Struct:
		// TODO: handle time
		fallthrough

	default:
		panic("invalid unmarshal state")
	}

	return
}

func (v valueUnmarshaler) unmarshalSlice(val reflect.Value, raw string) error {
	if xref.IsByteSlice(val.Type()) {
		val.Set(reflect.ValueOf([]byte(raw)))
		return nil
	}

	if tmp, err := v.unmarshalValue(val.Type().Elem(), raw); err != nil {
		return err
	} else {
		val.Set(reflect.Append(val, tmp))
	}
	return nil
}

func (v valueUnmarshaler) unmarshalMap(m reflect.Value, raw string) error {
	key, val, err := xraw.ParseMapEntry(raw, &v.props.Maps)

	if err != nil {
		return err
	}

	mt := m.Type()
	kt := mt.Key()
	vt := mt.Elem()

	kv := reflect.New(kt).Interface()
	if err := v.Unmarshal(key, kv); err != nil {
		return err
	}

	vv, err := v.unmarshalValue(vt, val)
	if err != nil {
		return err
	}

	if m.IsNil() {
		m.Set(reflect.MakeMap(mt))
	}
	m.SetMapIndex(reflect.ValueOf(kv).Elem(), vv)
	return nil
}

func (v valueUnmarshaler) unmarshalValue(vt reflect.Type, raw string) (reflect.Value, error) {
	if xref.IsByteSlice(vt) {
		return reflect.ValueOf([]byte(raw)), nil
	}

	if xref.IsBasicKind(vt.Kind()) {
		vv := reflect.New(vt).Interface()
		if err := v.Unmarshal(raw, vv); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(vv).Elem(), nil
	}

	if vt.Kind() == reflect.Ptr {
		vp := vt.Elem()
		if xref.IsByteSlice(vp) {
			tmp := []byte(raw)
			return reflect.ValueOf(&tmp), nil
		}

		if xref.IsBasicKind(vp.Kind()) {
			vv := reflect.New(vp).Interface()
			if err := v.Unmarshal(raw, vv); err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(vv), nil
		}
	}

	panic("invalid state")
}

func unmarshalInt(v reflect.Value, raw string, size int, props *argo.UnmarshalIntegerProps) error {
	if tmp, e := xraw.ParseInt(raw, size, props); e != nil {
		return &argo.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetInt(tmp)
	}
	return nil
}

func unmarshalInt64(v reflect.Value, raw string, props *argo.UnmarshalIntegerProps) error {
	if tmp, e := xraw.ParseInt(raw, 64, props); e != nil {

		// Because durations are just a wrapped int64 value, attempt to parse the
		// value as a duration.
		if v.Type().AssignableTo(reflect.TypeOf(time.Duration(0))) {
			if dur, e := time.ParseDuration(raw); e == nil {
				v.SetInt(int64(dur))
				return nil
			}
		}

		return &argo.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetInt(tmp)
		return nil
	}
}

func unmarshalUInt(v reflect.Value, raw string, size int, props *argo.UnmarshalIntegerProps) error {
	if tmp, e := xraw.ParseUInt(raw, size, props); e != nil {
		return &argo.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetUint(tmp)
	}
	return nil
}

func unmarshalFloat(v reflect.Value, raw string, size int) error {
	if tmp, e := strconv.ParseFloat(raw, size); e != nil {
		return &argo.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetFloat(tmp)
	}
	return nil
}

func unmarshalBool(v reflect.Value, raw string) error {
	if tmp, e := xraw.ParseBool(raw); e != nil {
		return &argo.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetBool(tmp)
	}
	return nil
}

// ////////////////////////////////////////////////////////////////////////// //

func ToUnmarshalable(arg string, ov reflect.Value, skipPtr bool) (reflect.Value, error) {

	if !skipPtr && (ov.Kind() != reflect.Ptr || IsNil(&ov)) {
		return reflect.Value{}, &argo.InvalidUnmarshalError{Value: ov, Argument: arg}
	}

	v := xref.GetRootValue(ov)

	kind := v.Kind()

	if xref.IsBasicKind(kind) {
		return v, nil
	}

	if v.Type().AssignableTo(unmarshalerType) {
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

	if !xref.IsBasicKind(vt.Key().Kind()) {
		return reflect.Value{}, &argo.InvalidTypeError{Value: ov}
	}

	if err := validateContainerValue(vt.Elem(), ov); err != nil {
		return reflect.Value{}, err
	}

	return v, nil
}

func validateContainerValue(t reflect.Type, ov reflect.Value) error {
	sk := t.Kind()

	if xref.IsBasicKind(sk) {
		return nil
	}

	if sk == reflect.Ptr {
		if xref.IsBasicKind(t.Elem().Kind()) {
			return nil
		}
		if xref.IsByteSlice(t.Elem()) {
			return nil
		}
		if xref.IsUnmarshaler(t.Elem()) {
			return nil
		}
		if xref.IsInterface(t.Elem()) {
			return nil
		}

		return &argo.InvalidTypeError{Value: ov}
	}

	if xref.IsByteSlice(t) {
		return nil
	}
	if xref.IsUnmarshaler(t) {
		return nil
	}
	if xref.IsInterface(t) {
		return nil
	}

	return &argo.InvalidTypeError{Value: ov}
}

func IsNil(v *reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
