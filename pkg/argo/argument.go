package argo

import (
	"fmt"
	"reflect"
)

// Argument represents a positional or flag argument that may be attached
// directly to a Command or CommandLeaf, or may be attached to a Flag.
type Argument interface {

	// Name returns the custom name assigned to this Argument.
	//
	// If no custom name was assigned to this Argument when it was built, this
	// method will return an empty string.
	Name() string

	// HasName tests whether this Argument has a custom name assigned.
	HasName() bool

	// Default returns the default value or value provider attached to this
	// Argument, if such a value exists.
	//
	// If this Argument does not have a default value or provider set, this method
	// will return nil.
	Default() any

	// HasDefault indicates whether a default value has been set on this
	// Argument.
	HasDefault() bool

	// DefaultType returns the reflect.Type value for the configured default
	// value.
	//
	// If no default value has been set on this Argument, this method will return
	// nil.
	DefaultType() reflect.Type

	// Description returns the description attached to this Argument.
	//
	// If no description was attached to this Argument when it was built, this
	// method will return an empty string.
	Description() string

	// HasDescription tests whether this Argument has a description attached.
	HasDescription() bool

	// WasHit tests whether this Argument was hit in a CLI call.
	//
	// This does not necessarily indicate that there is no value available for
	// this argument, just that it wasn't hit in the CLI call.  If the argument
	// had a default value provided, it will have been set in that case.
	WasHit() bool

	// RawValue returns the raw text value that was assigned to this Argument in
	// the CLI call.
	//
	// If this Argument was not hit during the CLI call, this method will return
	// an empty string.  This empty string IS NOT an indicator whether this
	// Argument was hit, as it may have been intentionally assigned an empty
	// value.  To test whether the Argument was hit, use WasHit.
	RawValue() string

	// IsRequired returns whether this Argument is required by its parent CLI
	// component.
	//
	// When parsing the CLI, if this argument is not found, an error will be
	// returned.
	IsRequired() bool

	binding() any

	// HasBinding indicates whether this Argument has a value binding.
	HasBinding() bool

	// BindingType returns the reflect.Type value for the configured binding.
	//
	// If this argument has no binding, this method will return nil.
	BindingType() reflect.Type
	setValue(rawValue string) error
	setToDefault() error
}

type argument struct {
	name string
	desc string
	raw  string

	required  bool
	isBindSet bool
	isDefSet  bool
	isUsed    bool

	bindVal any
	defVal  any

	rootBind reflect.Value
	rootDef  reflect.Value

	unmarshal ValueUnmarshaler
}

func (a argument) Name() string {
	return a.name
}

func (a argument) HasName() bool {
	return len(a.name) > 0
}

func (a argument) Description() string {
	return a.desc
}

func (a argument) HasDescription() bool {
	return len(a.desc) > 0
}

func (a argument) binding() any {
	return a.bindVal
}

func (a argument) HasBinding() bool {
	return a.isBindSet
}

func (a argument) BindingType() reflect.Type {
	if !a.isBindSet {
		return nil
	} else {
		return a.rootBind.Type()
	}
}

func (a argument) Default() any {
	return a.defVal
}

func (a argument) HasDefault() bool {
	return a.isDefSet
}

func (a argument) DefaultType() reflect.Type {
	if a.isDefSet {
		return a.rootDef.Type()
	} else {
		return nil
	}
}

func (a argument) WasHit() bool {
	return a.isUsed
}

func (a argument) RawValue() string {
	return a.raw
}

func (a argument) IsRequired() bool {
	return a.required
}

func (a *argument) setToDefault() error {
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

			if reflectIsUnmarshaler(a.rootBind.Type()) {
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

func (a *argument) setValue(rawString string) error {
	a.isUsed = true
	a.raw = rawString

	if !a.isBindSet {
		return nil
	}

	// TODO: why the heck is this here? what did past me know that present me doesn't?
	if a.isBoolArg() {
		if _, err := parseBool(rawString); err != nil {
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
