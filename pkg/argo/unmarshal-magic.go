package argo

import (
	"reflect"
	"strconv"
	"time"
)

var unmarshalerType = reflect.TypeOf((*Consumer)(nil)).Elem()

// NewDefaultMagicUnmarshaler creates a new "magic" ValueUnmarshaler instance
// using default UnmarshalProps.
//
// This is the default ValueUnmarshaler used by all arguments if not otherwise
// specified.
func NewDefaultMagicUnmarshaler() ValueUnmarshaler {
	return NewMagicUnmarshaler(DefaultUnmarshalProps())
}

// NewMagicUnmarshaler creates a new "magic" ValueUnmarshaler instance using the
// given UnmarshalProps.
func NewMagicUnmarshaler(props UnmarshalProps) ValueUnmarshaler {
	return valueUnmarshaler{props: props}
}

type valueUnmarshaler struct {
	props UnmarshalProps
}

func (v valueUnmarshaler) Unmarshal(raw string, val interface{}) (err error) {
	// If input string was empty, there's nothing to do
	if len(raw) == 0 {
		return nil
	}

	ptrVal := reflect.ValueOf(val)

	if ptrVal, err = toUnmarshalable(raw, ptrVal, false); err != nil {
		return err
	}

	ptrDef := ptrVal.Type()

	if ptrDef.AssignableTo(unmarshalerType) {
		return ptrVal.Interface().(Consumer).Accept(raw)
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
	if reflectIsByteSlice(val.Type()) {
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
	key, val, err := parseMapEntry(raw, &v.props.Maps)

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
	if reflectIsByteSlice(vt) {
		return reflect.ValueOf([]byte(raw)), nil
	}

	if reflectIsBasicKind(vt.Kind()) {
		vv := reflect.New(vt).Interface()
		if err := v.Unmarshal(raw, vv); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(vv).Elem(), nil
	}

	if vt.Kind() == reflect.Ptr {
		vp := vt.Elem()
		if reflectIsByteSlice(vp) {
			tmp := []byte(raw)
			return reflect.ValueOf(&tmp), nil
		}

		if reflectIsBasicKind(vp.Kind()) {
			vv := reflect.New(vp).Interface()
			if err := v.Unmarshal(raw, vv); err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(vv), nil
		}
	}

	panic("invalid state")
}

func unmarshalInt(v reflect.Value, raw string, size int, props *UnmarshalIntegerProps) error {
	if tmp, e := parseInt(raw, size, props); e != nil {
		return formatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetInt(tmp)
	}
	return nil
}

func unmarshalInt64(v reflect.Value, raw string, props *UnmarshalIntegerProps) error {
	if tmp, e := parseInt(raw, 64, props); e != nil {

		// Because durations are just a wrapped int64 value, attempt to parse the
		// value as a duration.
		if v.Type().AssignableTo(reflect.TypeOf(time.Duration(0))) {
			if dur, e := time.ParseDuration(raw); e == nil {
				v.SetInt(int64(dur))
				return nil
			}
		}

		return formatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetInt(tmp)
		return nil
	}
}

func unmarshalUInt(v reflect.Value, raw string, size int, props *UnmarshalIntegerProps) error {
	if tmp, e := parseUInt(raw, size, props); e != nil {
		return formatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetUint(tmp)
	}
	return nil
}

func unmarshalFloat(v reflect.Value, raw string, size int) error {
	if tmp, e := strconv.ParseFloat(raw, size); e != nil {
		return formatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetFloat(tmp)
	}
	return nil
}

func unmarshalBool(v reflect.Value, raw string) error {
	if tmp, e := parseBool(raw); e != nil {
		return formatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetBool(tmp)
	}
	return nil
}

// ////////////////////////////////////////////////////////////////////////// //

func toUnmarshalable(arg string, ov reflect.Value, skipPtr bool) (reflect.Value, error) {

	if !skipPtr && (ov.Kind() != reflect.Ptr || isNil(&ov)) {
		return reflect.Value{}, &InvalidUnmarshalError{Value: ov, Argument: arg}
	}

	v := reflectGetRootValue(ov)

	kind := v.Kind()

	if reflectIsBasicKind(kind) {
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

	return reflect.Value{}, &InvalidTypeError{Value: v}
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

	if !reflectIsBasicKind(vt.Key().Kind()) {
		return reflect.Value{}, &InvalidTypeError{Value: ov}
	}

	if err := validateContainerValue(vt.Elem(), ov); err != nil {
		return reflect.Value{}, err
	}

	return v, nil
}

func validateContainerValue(t reflect.Type, ov reflect.Value) error {
	sk := t.Kind()

	if reflectIsBasicKind(sk) {
		return nil
	}

	if sk == reflect.Ptr {
		if reflectIsBasicKind(t.Elem().Kind()) {
			return nil
		}
		if reflectIsByteSlice(t.Elem()) {
			return nil
		}
		if reflectIsUnmarshaler(t.Elem()) {
			return nil
		}
		if reflectIsInterface(t.Elem()) {
			return nil
		}

		return &InvalidTypeError{Value: ov}
	}

	if reflectIsByteSlice(t) {
		return nil
	}
	if reflectIsUnmarshaler(t) {
		return nil
	}
	if reflectIsInterface(t) {
		return nil
	}

	return &InvalidTypeError{Value: ov}
}

func isNil(v *reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
