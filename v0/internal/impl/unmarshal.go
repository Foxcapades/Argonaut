package impl

import (
	R "reflect"
	S "strconv"

	U "github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type iProps = A.UnmarshalIntegerProps

var unmarshalerType = R.TypeOf((*A.Unmarshaler)(nil)).Elem()

func NewDefaultedValueUnmarshaler() A.ValueUnmarshaler {
	return NewValueUnmarshaler(&defaultUnmarshalProps)
}

func NewValueUnmarshaler(props *A.UnmarshalProps) A.ValueUnmarshaler {
	return &ValueUnmarshaler{props: props}
}

type ValueUnmarshaler struct {
	props *A.UnmarshalProps
}

func (v *ValueUnmarshaler) Unmarshal(raw string, val interface{}) (err error) {
	// If input string was empty, there's nothing to do
	if len(raw) == 0 {
		return nil
	}

	ptrVal := R.ValueOf(val)

	if ptrVal, err = U.ToUnmarshalable(raw, ptrVal); err != nil {
		return err
	}

	ptrDef := ptrVal.Type()

	if ptrDef.AssignableTo(unmarshalerType) {
		return ptrVal.Interface().(A.Unmarshaler).Unmarshal(raw)
	}

	switch ptrVal.Kind() {

	case R.String:
		ptrVal.SetString(raw)
	case R.Int:
		return unmarshalInt(ptrVal, raw, S.IntSize, &v.props.Integers)
	case R.Float32:
		return unmarshalFloat(ptrVal, raw, 32)
	case R.Int64:
		return unmarshalInt(ptrVal, raw, 64, &v.props.Integers)
	case R.Float64:
		return unmarshalFloat(ptrVal, raw, 64)
	case R.Uint64:
		return unmarshalUInt(ptrVal, raw, 64, &v.props.Integers)
	case R.Uint:
		return unmarshalUInt(ptrVal, raw, S.IntSize, &v.props.Integers)
	case R.Int32:
		return unmarshalInt(ptrVal, raw, 32, &v.props.Integers)
	case R.Uint32:
		return unmarshalUInt(ptrVal, raw, 32, &v.props.Integers)
	case R.Uint8:
		return unmarshalUInt(ptrVal, raw, 8, &v.props.Integers)
	case R.Slice:
		return v.unmarshalSlice(ptrVal, raw)
	case R.Map:
		return v.unmarshalMap(ptrVal, raw)
	case R.Int8:
		return unmarshalInt(ptrVal, raw, 8, &v.props.Integers)
	case R.Int16:
		return unmarshalInt(ptrVal, raw, 16, &v.props.Integers)
	case R.Uint16:
		return unmarshalUInt(ptrVal, raw, 16, &v.props.Integers)
	case R.Bool:
		return unmarshalBool(ptrVal, raw)
	case R.Interface:
		ptrVal.Set(R.ValueOf(raw))

	default:
		panic("invalid unmarshal state")
	}

	return
}

func unmarshalInt(v R.Value, raw string, size int, props *iProps) error {
	if tmp, e := U.ParseInt(raw, size, props); e != nil {
		return &A.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetInt(tmp)
	}
	return nil
}

func unmarshalUInt(v R.Value, raw string, size int, props *iProps) error {
	if tmp, e := U.ParseUInt(raw, size, props); e != nil {
		return &A.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetUint(tmp)
	}
	return nil
}

func unmarshalFloat(v R.Value, raw string, size int) error {
	if tmp, e := S.ParseFloat(raw, size); e != nil {
		return &A.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetFloat(tmp)
	}
	return nil
}

func unmarshalBool(v R.Value, raw string) error {
	if tmp, e := U.ParseBool(raw); e != nil {
		return &A.FormatError{Value: v, Argument: raw, Kind: v.Kind(), Root: e}
	} else {
		v.SetBool(tmp)
	}
	return nil
}

func (v *ValueUnmarshaler) unmarshalSlice(val R.Value, raw string) error {
	if U.IsByteSlice(val.Type()) {
		val.Set(R.ValueOf([]byte(raw)))
		return nil
	}

	if tmp, err := v.unmarshalValue(val.Type().Elem(), raw); err != nil {
		return err
	} else {
		val.Set(R.Append(val, tmp))
	}
	return nil
}

func (v *ValueUnmarshaler) unmarshalMap(m R.Value, raw string) error {
	key, val, err := U.ParseMapEntry(raw, &v.props.Maps)

	if err != nil {
		return err
	}

	mt := m.Type()
	kt := mt.Key()
	vt := mt.Elem()

	kv := R.New(kt).Interface()
	if err := v.Unmarshal(key, kv); err != nil {
		return err
	}

	vv, err := v.unmarshalValue(vt, val)
	if err != nil {
		return err
	}

	if m.IsNil() {
		m.Set(R.MakeMap(mt))
	}
	m.SetMapIndex(R.ValueOf(kv).Elem(), vv)
	return nil
}

func (v *ValueUnmarshaler) unmarshalValue(vt R.Type, raw string) (R.Value, error) {
	if U.IsByteSlice(vt) {
		return R.ValueOf([]byte(raw)), nil
	}

	if U.IsBasicKind(vt.Kind()) {
		vv := R.New(vt).Interface()
		if err := v.Unmarshal(raw, vv); err != nil {
			return R.Value{}, err
		}
		return R.ValueOf(vv).Elem(), nil
	}

	if vt.Kind() == R.Ptr {
		vp := vt.Elem()
		if U.IsByteSlice(vp) {
			tmp := []byte(raw)
			return R.ValueOf(&tmp), nil
		}

		if U.IsBasicKind(vp.Kind()) {
			vv := R.New(vp).Interface()
			if err := v.Unmarshal(raw, vv); err != nil {
				return R.Value{}, err
			}
			return R.ValueOf(vv), nil
		}
	}

	panic("invalid state")
}
