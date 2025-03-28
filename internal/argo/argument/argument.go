package argument

import (
	"fmt"
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/unmarshal"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type argument struct {
	name string
	desc string
	raw  string

	required bool
	isUsed   bool

	binding Binding
	defVal  Default

	unmarshal argo.ValueUnmarshaler

	preParseValidators  []any
	postParseValidators []any
}

func (a *argument) Name() string {
	return a.name
}

func (a *argument) HasName() bool {
	return len(a.name) > 0
}

func (a *argument) Description() string {
	return a.desc
}

func (a *argument) HasDescription() bool {
	return len(a.desc) > 0
}

func (a *argument) HasBinding() bool {
	return a.binding.bType != argo.BindingTypeNone
}

func (a *argument) Binding() argo.ArgumentBinding {
	return &a.binding
}

func (a *argument) HasDefault() bool {
	return a.defVal.dType != argo.DefaultTypeNone
}

func (a *argument) Default() any {
	return a.defVal.value
}

func (a *argument) WasHit() bool {
	return a.isUsed
}

func (a *argument) RawValue() string {
	return a.raw
}

func (a *argument) IsRequired() bool {
	return a.required
}

func (a *argument) SetToDefault() error {
	// If there is no binding set, what are we going to set to the default value?
	if !a.HasBinding() {
		return nil
	}

	// If there is no default set, what are we going to do here?
	if !a.HasDefault() {
		return nil
	}

	a.isUsed = true

	rootDef := unmarshal.GetRootValue(reflect.ValueOf(a.defVal.value))
	defType := rootDef.Type()

	rootBinding := unmarshal.GetRootValue(reflect.ValueOf(a.binding.raw))

	if defType.Kind() == reflect.Func {
		defFn := reflect.ValueOf(a.defVal.value)

		switch defType.NumOut() {

		// Function returns (value)
		case 1:
			ret := defFn.Call(nil)

			rootBinding.Set(ret[0])
			a.raw = ret[0].Type().String()

			return nil

		// Function returns (value, error)
		case 2:
			ret := defFn.Call(nil)

			// If err != nil
			if !ret[1].IsNil() {
				return ret[1].Interface().(error)
			}

			if unmarshal.IsUnmarshaler(rootBinding.Type()) {
				rootBinding.Elem().Set(ret[0])
			} else {
				rootBinding.Set(ret[0])
			}

			a.raw = ret[0].Type().String()

			return nil

		default:
			panic(fmt.Errorf("given default value provider returns an invalid number of arguments (%d), expected 1 or 2", defType.NumOut()))
		}
	}

	if defType.Kind() == reflect.String {
		strVal := rootDef.String()

		if rootBinding.Type().Kind() == reflect.String {
			rootBinding.Set(rootDef)
			a.raw = strVal
			return nil
		}

		return a.unmarshal.Unmarshal(strVal, a.binding.raw)
	}

	if rootBinding.Kind() == reflect.Ptr {
		rootBinding.Elem().Set(rootDef)
	} else {
		rootBinding.Set(rootDef)
	}

	return nil
}

func (a *argument) SetValue(rawString string) error {
	a.isUsed = true
	a.raw = rawString

	for _, fn := range a.preParseValidators {
		if err := a.callPreArgFunc(fn, rawString); err != nil {
			return err
		}
	}

	if !a.HasBinding() {
		return nil
	}

	rootBinding := unmarshal.GetRootValue(reflect.ValueOf(a.binding.BoundTo()))

	if err := a.unmarshal.Unmarshal(rawString, a.binding.raw); err != nil {
		return err
	}

	for _, fn := range a.postParseValidators {
		if err := a.callPostArgFunc(fn, rootBinding, rawString); err != nil {
			return err
		}
	}

	return nil
}

func (a *argument) callPreArgFunc(fn any, raw string) error {
	errs := reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(raw)})

	if !errs[0].IsNil() {
		return errs[0].Interface().(error)
	}

	return nil
}

func (a *argument) callPostArgFunc(fn any, root reflect.Value, raw string) error {
	errs := reflect.ValueOf(fn).Call([]reflect.Value{root, reflect.ValueOf(raw)})

	if !errs[0].IsNil() {
		return errs[0].Interface().(error)
	}

	return nil
}
