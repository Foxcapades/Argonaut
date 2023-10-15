package argument

import (
	"fmt"
	"reflect"

	"github.com/Foxcapades/Argonaut/v1/internal/marsh"
	"github.com/Foxcapades/Argonaut/v1/internal/xraw"
	"github.com/Foxcapades/Argonaut/v1/internal/xref"
)

// implements argo.Argument
type argument struct {
	name string
	desc string

	raw string

	required  bool
	isBindSet bool
	isDefSet  bool
	isUsed    bool

	bindVal any
	defVal  any

	rootBind reflect.Value
	rootDef  reflect.Value

	unmarshal marsh.ValueUnmarshaler
}

func (a argument) Name() string  { return a.name }
func (a argument) HasName() bool { return len(a.name) > 0 }

func (a argument) Description() string  { return a.desc }
func (a argument) HasDescription() bool { return len(a.desc) > 0 }

func (a argument) WasHit() bool { return a.isUsed }

func (a argument) RawValue() string { return a.raw }

func (a argument) IsRequired() bool { return a.required }

func (a *argument) SetDefault() error {
	// If there is no binding set, what are we going to set to the default value?
	if !a.isBindSet {
		return nil
	}

	// If there is no default set, what are we going to do here?
	if !a.isDefSet {
		return nil
	}

	defType := a.rootDef.Type()

	if defType.Kind() == reflect.Func {
		defFn := reflect.ValueOf(a.defVal)

		switch defType.NumOut() {

		// Function returns (value)
		case 1:
			ret := defFn.Call(nil)

			a.rootBind.Set(ret[0])
			a.raw = ret[0].Type().String()

			return nil

		// Function returns (value, error)
		case 2:
			ret := defFn.Call(nil)

			// If err != nil
			if !ret[1].IsNil() {
				panic(ret[1].Interface())
			}

			if xref.IsUnmarshaler(a.rootBind.Type()) {
				a.rootBind.Elem().Set(ret[0])
			} else {
				a.rootBind.Set(ret[0])
			}

			a.raw = ret[0].Type().String()

			return nil

		default:
			panic(fmt.Errorf("given default value provider returns an invalid number of arguments (%d), expected 1 or 2", defType.NumOut()))
		}
	}

	if defType.Kind() == reflect.String {
		strVal := a.rootDef.String()

		if a.rootBind.Type().Kind() == reflect.String {
			a.rootBind.Set(a.rootDef)
			a.raw = strVal
			return nil
		}

		return a.unmarshal.Unmarshal(strVal, a.bindVal)
	}

	a.rootBind.Set(a.rootDef)
	return nil
}

func (a *argument) SetValue(rawString string) error {
	a.isUsed = true
	a.raw = rawString

	if !a.isBindSet {
		return nil
	}

	if a.unmarshal != nil {
		var tmp any

		if err := a.unmarshal.Unmarshal(rawString, &tmp); err != nil {
			return err
		}

		a.rootBind.Set(reflect.ValueOf(tmp))

		return nil
	}

	if a.isBoolArg() {
		if _, err := xraw.ParseBool(rawString); err != nil {
			return err
		}

		return a.unmarshal.Unmarshal(rawString, a.bindVal)
	}

	return a.unmarshal.Unmarshal(rawString, a.bindVal)
}

func (a *argument) isBoolArg() bool {
	bt := a.rootBind.Type().String()
	return bt == "bool" || bt == "*bool" || bt == "[]bool" || bt == "[]*bool"
}
