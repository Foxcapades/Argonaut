package argo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshalDefault(raw string, val interface{}) (err error) { return nil}

func Unmarshal(raw string, val interface{}, props UnmarshalProps) (err error) {
	// If input string was empty, there's nothing to do
	if len(raw) == 0 {
		return nil
	}

	ptrDef := reflect.TypeOf(val)
	ptrVal := reflect.ValueOf(val)
	subDef := ptrDef.Elem()

	if ptrVal.Kind() != reflect.Ptr || ptrVal.IsNil() {
		return &InvalidUnmarshalError{ptrVal, raw}
	}

	ptrVal = getRootValue(ptrVal)

	if ptrVal.Type().AssignableTo(unmarshalerType) {
		return val.(Unmarshaler).Unmarshal(raw)
	}

	switch ptrVal.Kind() {

	case reflect.String:
		*(val.(*string)) = raw

	case reflect.Int:
		if tmp, e := parseInt(raw, strconv.IntSize, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetInt(tmp)
		}

	case reflect.Float32:
		// TODO: wrap this error
		if tmp, e := strconv.ParseFloat(raw, 32); e != nil {
			err = e
		} else {
			ptrVal.SetFloat(tmp)
		}

	case reflect.Int64:
		if tmp, e := parseInt(raw, 64, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetInt(tmp)
		}

	case reflect.Float64:
		// TODO: wrap this error
		*(val.(*float64)), err = strconv.ParseFloat(raw, 64)

	case reflect.Uint64:
		if tmp, e := parseUInt(raw, 64, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetUint(tmp)
		}

	case reflect.Uint:
		if tmp, e := parseUInt(raw, strconv.IntSize, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetUint(tmp)
		}

	case reflect.Int32:
		if tmp, e := parseInt(raw, 32, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetInt(tmp)
		}

	case reflect.Uint32:
		if tmp, e := parseUInt(raw, 32, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetUint(tmp)
		}

	case reflect.Uint8:
		if tmp, e := parseUInt(raw, 8, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetUint(tmp)
		}

	case reflect.Slice:
		// check slice is unmarshalable
	case reflect.Map:
		// check map is unmarshalable

	case reflect.Int8:
		if tmp, e := parseInt(raw, 8, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetInt(tmp)
		}

	case reflect.Int16:
		if tmp, e := parseInt(raw, 16, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetInt(tmp)
		}

	case reflect.Uint16:
		if tmp, e := parseUInt(raw, 16, &props.Integers); e != nil {
			err = e
		} else {
			ptrVal.SetUint(tmp)
		}

	case reflect.Interface:
		ptrVal.Set(reflect.ValueOf(raw))

	default:
		// TODO: make this an official error
		return fmt.Errorf("cannot unmarshal type %s", subDef)
	}

	return
}

func parseInt(v string, bits int, opt *UnmarshalIntegerProps) (int64, error) {
	var neg string
	// TODO: Wrap this error

	if v[0] == '-' {
		neg = "-"
		v = v[1:]
	}

	for i := range opt.HexLeaders {
		if strings.HasPrefix(v, opt.HexLeaders[i]) {
			return strconv.ParseInt(neg + v[len(opt.HexLeaders[i]):], 16, bits)
		}
	}

	for i := range opt.OctalLeaders {
		if strings.HasPrefix(v, opt.OctalLeaders[i]) {
			return strconv.ParseInt(neg + v[len(opt.OctalLeaders[i]):], 8, bits)
		}
	}

	return strconv.ParseInt(neg + v, 10, bits)
}

func parseUInt(v string, bits int, opt *UnmarshalIntegerProps) (uint64, error) {
	// TODO: Wrap this error

	for i := range opt.HexLeaders {
		if strings.HasPrefix(v, opt.HexLeaders[i]) {
			return strconv.ParseUint(v[len(opt.HexLeaders[i]):], 16, bits)
		}
	}

	for i := range opt.OctalLeaders {
		if strings.HasPrefix(v, opt.OctalLeaders[i]) {
			return strconv.ParseUint(v[len(opt.OctalLeaders[i]):], 8, bits)
		}
	}

	return strconv.ParseUint(v, 10, bits)
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "t", "yes", "y", "1":
		return true, nil
	case "false", "f", "no", "n", "0":
		return false, nil
	default:
		// TODO: return an error here
		return false, nil
	}
}

func getRootValue(v reflect.Value) reflect.Value {
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

		if v.Type().AssignableTo(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
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