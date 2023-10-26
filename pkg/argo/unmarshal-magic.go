package argo

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/Foxcapades/Argonaut/internal/unmarshal"
	"github.com/Foxcapades/Argonaut/internal/xreflect"
)

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

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

	if ptrVal, err = unmarshal.ToUnmarshalable(raw, ptrVal, false, unmarshalerType); err != nil {
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
		if ptrDef.AssignableTo(reflect.TypeOf(time.Time{})) {
			return v.unmarshalTime(ptrVal, raw)
		}

		fallthrough

	default:
		panic("invalid unmarshal state")
	}

	return
}

func (v valueUnmarshaler) unmarshalTime(val reflect.Value, raw string) error {
	for _, format := range v.props.Time.DateFormats {
		if tim, err := time.Parse(format, raw); err == nil {
			val.Set(reflect.ValueOf(tim))
			return nil
		}
	}

	return errors.New("could not parse input string as a time value")
}

func (v valueUnmarshaler) unmarshalSlice(val reflect.Value, raw string) error {
	if xreflect.IsByteSlice(val.Type()) {
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
	parser := newMapElementParser(v.props.Maps, raw)

	for parser.HasNext() {
		key, val, err := parseMapEntry(&parser)
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

		if (xreflect.IsBasicSlice(vt) && !xreflect.IsByteSlice(vt)) || xreflect.IsUnmarshalerSlice(vt, unmarshalerType) {
			rkv := reflect.ValueOf(kv).Elem()

			tmp := m.MapIndex(rkv)
			if tmp.Kind() == reflect.Invalid {
				slice := reflect.MakeSlice(vt, 1, 10)
				rvv := slice.Index(0)
				rvv.Set(vv)
				m.SetMapIndex(rkv, slice)
			} else {
				m.SetMapIndex(rkv, reflect.Append(tmp, vv))
			}
		} else {
			m.SetMapIndex(reflect.ValueOf(kv).Elem(), vv)
		}
	}

	return nil
}

// unmarshalValue value is used to unmarshal the element contained within a map
// or slice.
func (v valueUnmarshaler) unmarshalValue(vt reflect.Type, raw string) (reflect.Value, error) {
	// If the type of the element is a byte slice, then pass it up in raw string
	// form.
	if xreflect.IsByteSlice(vt) {
		return reflect.ValueOf([]byte(raw)), nil
	}

	//
	if xreflect.IsBasicKind(vt.Kind()) {
		vv := reflect.New(vt).Interface()
		if err := v.Unmarshal(raw, vv); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(vv).Elem(), nil
	}

	if xreflect.IsBasicSlice(vt) {
		return v.unmarshalValue(vt.Elem(), raw)
	}

	if vt.Kind() == reflect.Ptr {
		vp := vt.Elem()
		if xreflect.IsByteSlice(vp) {
			tmp := []byte(raw)
			return reflect.ValueOf(&tmp), nil
		}

		if xreflect.IsBasicKind(vp.Kind()) {
			vv := reflect.New(vp).Interface()
			if err := v.Unmarshal(raw, vv); err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(vv), nil
		}

		if xreflect.IsUnmarshaler(vt, unmarshalerType) {
			vv := reflect.New(vp).Interface()
			if err := vv.(Unmarshaler).Unmarshal(raw); err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(vv), nil
		}
	}

	if xreflect.IsUnmarshalerSlice(vt, unmarshalerType) {
		return v.unmarshalValue(vt.Elem(), raw)
	}

	panic("invalid state: the given type cannot be unmarshalled")
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
