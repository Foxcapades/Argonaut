package impl

import (
	R "reflect"
	S "strconv"

	U "github.com/Foxcapades/Argonaut/v1/internal/util"
	A "github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

type aProps = A.UnmarshalProps
type iProps = A.UnmarshalIntegerProps

var unmarshalerType = R.TypeOf((*A.Unmarshaler)(nil)).Elem()

func UnmarshalDefault(raw string, val interface{}) (err error) {
	return Unmarshal(raw, val, &defaultUnmarshalProps)
}

func Unmarshal(raw string, val interface{}, props *aProps) (err error) {
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
		return unmarshalInt(ptrVal, raw, S.IntSize, &props.Integers)
	case R.Float32:
		return unmarshalFloat(ptrVal, raw, 32)
	case R.Int64:
		return unmarshalInt(ptrVal, raw, 64, &props.Integers)
	case R.Float64:
		return unmarshalFloat(ptrVal, raw, 64)
	case R.Uint64:
		return unmarshalUInt(ptrVal, raw, 64, &props.Integers)
	case R.Uint:
		return unmarshalUInt(ptrVal, raw, S.IntSize, &props.Integers)
	case R.Int32:
		return unmarshalInt(ptrVal, raw, 32, &props.Integers)
	case R.Uint32:
		return unmarshalUInt(ptrVal, raw, 32, &props.Integers)
	case R.Uint8:
		return unmarshalUInt(ptrVal, raw, 8, &props.Integers)
	case R.Slice:
		return unmarshalSlice(ptrVal, raw, props)
	case R.Map:
		return unmarshalMap(ptrVal, raw, props)
	case R.Int8:
		return unmarshalInt(ptrVal, raw, 8, &props.Integers)
	case R.Int16:
		return unmarshalInt(ptrVal, raw, 16, &props.Integers)
	case R.Uint16:
		return unmarshalUInt(ptrVal, raw, 16, &props.Integers)
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

func unmarshalSlice(v R.Value, raw string, props *aProps) error {
	if U.IsByteSlice(v.Type()) {
		v.Set(R.ValueOf([]byte(raw)))
		return nil
	}

	if tmp, err := unmarshalValue(v.Type().Elem(), raw, props); err != nil {
		return err
	} else {
		v.Set(R.Append(v, tmp))
	}
	return nil
}

func unmarshalMap(m R.Value, raw string, props *aProps) error {
	key, val, err := U.ParseMapEntry(raw, &props.Maps)

	if err != nil {
		return err
	}

	mt := m.Type()
	kt := mt.Key()
	vt := mt.Elem()

	kv := R.New(kt).Interface()
	if err := Unmarshal(key, kv, props); err != nil {
		return err
	}

	vv, err := unmarshalValue(vt, val, props)
	if err != nil {
		return err
	}

	m.Set(R.MakeMap(mt))
	m.SetMapIndex(R.ValueOf(kv).Elem(), vv)
	return nil
}

func unmarshalValue(vt R.Type, raw string, props *aProps) (R.Value, error) {
	if U.IsByteSlice(vt) {
		return R.ValueOf([]byte(raw)), nil
	}

	if U.IsBasicKind(vt.Kind()) {
		vv := R.New(vt).Interface()
		if err := Unmarshal(raw, vv, props); err != nil {
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
			if err := Unmarshal(raw, vv, props); err != nil {
				return R.Value{}, err
			}
			return R.ValueOf(vv), nil
		}
	}

	panic("invalid state")
}
