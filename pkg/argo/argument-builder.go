package argo

import (
	"errors"
	"fmt"
	"reflect"
)

// An ArgumentBuilder instance is used to construct a CLI argument that may be
// attached to a Flag or CommandLeaf.
type ArgumentBuilder interface {

	// WithName sets the name for this argument.
	//
	// The name value is used when rendering help information about this argument.
	WithName(name string) ArgumentBuilder

	// WithDescription sets the description of this argument to be shown in rendered
	// help text.
	WithDescription(desc string) ArgumentBuilder

	// Require marks the output Argument as being required.
	Require() ArgumentBuilder

	// WithBinding sets the bind value for the Argument.
	//
	// The bind value may either be a pointer or an instance of Unmarshaler.
	//
	// If the bind value is a pointer, the Argument's value unmarshaler will be
	// called to unmarshal the raw string value into a value of the type passed
	// to this method.
	//
	// If the bind value is a Unmarshaler instance, that instance's Unmarshal method
	// will be called with the raw input from the CLI.
	//
	// Setting this value to anything other than a pointer or a Unmarshaler instance
	// will result in an error being returned when building the argument is
	// attempted.
	//
	// Example 1:
	//     var myValue time.Duration
	//     cli.Argument.WithBinding(&myValue)
	//
	// Example 2:
	//     cli.Argument.WithBinding(UnmarshalerFunc(func(raw string) error {
	//         fmt.Println(raw)
	//         return nil
	//     }))
	//
	// Example 3 (lets get silly with it):
	//     var myValue map[bool]**string
	//     cli.Argument.WithBinding(&myValue)
	WithBinding(pointer any) ArgumentBuilder

	getBinding() any

	// WithDefault sets the default value for the argument to be used if the
	// argument is not provided on the command line.
	//
	// Setting this value without providing a binding value using `Bind()` will
	// mean that the given default will not be set to anything when the CLI input
	// is parsed.
	//
	// When used, the type of this value must meet one of the following criteria:
	//   1. `val` is compatible with the type of the value used with `Bind()`
	//   2. `val` is a function which returns a type that is compatible with the
	//      type of the value used with `Bind()`
	//   3. `val` is a function which returns a type that is compatible with the
	//      type of the value used with `Bind()` in addition to returning an error
	//      as the second return value.
	//
	// Examples:
	//     arg.WithBinding(&fooString).WithDefault(3)   // Type mismatch
	//
	//     arg.WithBinding(&fooInt).WithDefault(3)      // OK
	//
	//     arg.WithBinding(&fooInt).
	//       WithDefault(func() int {return 3})         // OK
	//
	//     arg.WithBinding(&fooInt).
	//       WithDefault(func() (int, error) {
	//         return 3, nil
	//       })                                         // OK
	//
	// If the value provided to this method is a pointer to the type of the bind
	// value it will be dereferenced to set the bind value.
	WithDefault(def any) ArgumentBuilder

	getDefault() any

	// WithUnmarshaler allows providing a custom ValueUnmarshaler instance that
	// will be used to unmarshal string values into the binding type.
	//
	// If no binding is set on this argument, the provided ValueUnmarshaler will
	// not be used.
	//
	// If a custom unmarshaler is not provided by way of this method, then the
	// internal magic unmarshaler will be used to parse raw arguments.
	WithUnmarshaler(fn ValueUnmarshaler) ArgumentBuilder

	// Build attempts to build an Argument instance out of the configuration given
	// to this ArgumentBuilder instance.
	Build(ctx *WarningContext) (Argument, error)
}

func NewArgumentBuilder() ArgumentBuilder {
	return &argumentBuilder{
		marsh: NewDefaultMagicUnmarshaler(),
	}
}

type argumentBuilder struct {
	name string

	desc string

	required bool
	hasBind  bool
	hasDef   bool

	def  any
	bind any

	rootDef  reflect.Value
	rootBind reflect.Value

	marsh ValueUnmarshaler
}

func (a *argumentBuilder) WithName(name string) ArgumentBuilder {
	a.name = name
	return a
}

func (a *argumentBuilder) WithDescription(desc string) ArgumentBuilder {
	a.desc = desc
	return a
}

func (a *argumentBuilder) Require() ArgumentBuilder {
	a.required = true
	return a
}

func (a *argumentBuilder) WithBinding(binding any) ArgumentBuilder {
	a.hasBind = true
	a.bind = binding
	return a
}

func (a *argumentBuilder) getBinding() any {
	return a.bind
}

func (a *argumentBuilder) WithDefault(def any) ArgumentBuilder {
	a.hasDef = true
	a.def = def
	return a
}

func (a argumentBuilder) getDefault() any {
	return a.def
}

func (a *argumentBuilder) WithUnmarshaler(fn ValueUnmarshaler) ArgumentBuilder {
	a.marsh = fn
	return a
}

func (a *argumentBuilder) Build(warnings *WarningContext) (Argument, error) {
	errs := newMultiError()
	valDefault := true

	if err := a.validateBinding(); err != nil {
		errs.AppendError(err)
		valDefault = false
	}

	if valDefault {
		if err := a.validateDefault(); err != nil {
			errs.AppendError(err)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &argument{
		warnings:  warnings,
		name:      a.name,
		desc:      a.desc,
		required:  a.required,
		isBindSet: a.hasBind,
		isDefSet:  a.hasDef,
		bindVal:   a.bind,
		defVal:    a.def,
		rootBind:  a.rootBind,
		rootDef:   a.rootDef,
		unmarshal: a.marsh,
	}, nil
}

func (a *argumentBuilder) validateBinding() error {
	if !a.hasBind {
		return nil
	}

	if tmp, err := toUnmarshalable("", reflect.ValueOf(a.bind), false); err != nil {
		return newInvalidArgError(ArgErrInvalidBindingBadType, a, "")
	} else {
		a.rootBind = tmp
	}

	return nil
}

const (
	errDefFnOutNum = "default value providers must return either 1 or 2 values"
	err2ndOut      = "the second output type of a default value provider must " +
		"be compatible with error"
	errBadType = "default value type %s is not compatible with binding type %s"
)

func (a *argumentBuilder) validateDefault() error {
	if !a.hasDef {
		return nil
	}

	if !a.hasBind {
		// TODO: this should be a real error
		return errors.New("default set with no binding")
	}

	if a.hasDef && reflectGetRootValue(reflect.ValueOf(a.def)).Kind() == reflect.Func {
		a.rootDef = reflectGetRootValue(reflect.ValueOf(a.def))
		return a.validateDefaultProvider()
	}

	if tmp, err := toUnmarshalable("", reflect.ValueOf(a.def), true); err != nil {
		// TODO: This is not necessarily the correct error type
		return invalidDefaultValError(a)
	} else {
		a.rootDef = tmp
	}

	if a.rootDef.Kind() != reflect.String && !reflectCompatible(&a.rootDef, &a.rootBind) {
		return invalidDefaultValError(a)
	}

	return nil
}

func (a *argumentBuilder) validateDefaultProvider() error {
	root := &a.rootDef
	defType := root.Type()

	oLen := defType.NumOut()
	if oLen == 0 || oLen > 2 {
		return newInvalidArgError(ArgErrInvalidDefaultFn, a, errDefFnOutNum)
	}

	outType := defType.Out(0)

	if !defType.Out(0).AssignableTo(a.rootBind.Type()) {

		// Second chance for Unmarshalable short-circuit logic.
		//
		// In this case we can attempt to unmarshal in 2 different scenarios:
		//   1. The bind value is an unmarshaler func AND the default provider
		//      returns a string.
		//   2. The bind value is an unmarshaler instance AND the default provider
		//      returns a value that is of that implementing type.
		if reflectIsUnmarshaler(a.rootBind.Type()) {
			// If the output type is a string, then it is compatible with both the
			// unmarshaler and unmarshaler func.
			if outType.Kind() == reflect.String {
				return nil
			}

			if a.rootBind.Type().Kind() != reflect.Func {
				if outType.AssignableTo(a.rootBind.Type().Elem()) {
					return nil
				}
			}
		}

		return newInvalidArgError(ArgErrInvalidDefaultVal, a,
			fmt.Sprintf(errBadType, defType.Out(0), reflect.TypeOf(a.bind)))
	}

	if oLen == 2 && !defType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return newInvalidArgError(ArgErrInvalidDefaultFn, a, err2ndOut)
	}

	return nil
}

func invalidDefaultValError(b *argumentBuilder) error {
	return newInvalidArgError(ArgErrInvalidDefaultVal, b,
		fmt.Sprintf(errBadType, b.rootDef.Type(), b.rootBind.Type()))
}
