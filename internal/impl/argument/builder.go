package argument

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Foxcapades/Argonaut/internal/marsh"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/internal/xref"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func NewBuilder() argo.ArgumentBuilder {
	return &argumentBuilder{
		marsh: marsh.NewDefaultedValueUnmarshaler(),
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

	marsh marsh.ValueUnmarshaler
}

func (a *argumentBuilder) WithName(name string) argo.ArgumentBuilder {
	a.name = name
	return a
}

func (a argumentBuilder) GetName() string {
	return a.name
}

func (a argumentBuilder) HasName() bool {
	return len(a.name) > 0
}

func (a *argumentBuilder) WithDescription(desc string) argo.ArgumentBuilder {
	a.desc = desc
	return a
}

func (a argumentBuilder) GetDescription() string {
	return a.desc
}

func (a argumentBuilder) HasDescription() bool {
	return len(a.desc) > 0
}

func (a *argumentBuilder) Require() argo.ArgumentBuilder {
	a.required = true
	return a
}

func (a argumentBuilder) IsRequired() bool {
	return a.required
}

func (a *argumentBuilder) WithBinding(binding any) argo.ArgumentBuilder {
	a.hasBind = true
	a.bind = binding
	return a
}

func (a argumentBuilder) HasBinding() bool {
	return a.hasBind
}

func (a *argumentBuilder) GetBinding() any {
	return a.bind
}

func (a *argumentBuilder) WithDefault(def any) argo.ArgumentBuilder {
	a.hasDef = true
	a.def = def
	return a
}

func (a argumentBuilder) HasDefault() bool {
	return a.hasDef
}

func (a argumentBuilder) GetDefault() any {
	return a.def
}

func (a *argumentBuilder) WithUnmarshaler(fn marsh.ValueUnmarshaler) argo.ArgumentBuilder {
	a.marsh = fn
	return a
}

func (a *argumentBuilder) Build() (argo.Argument, error) {
	errs := xerr.NewMultiError()

	if err := a.validateBinding(); err != nil {
		errs.AppendError(err)
	}

	if err := a.validateDefault(); err != nil {
		errs.AppendError(err)
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &argument{
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

	if tmp, err := marsh.ToUnmarshalable("", reflect.ValueOf(a.bind), false); err != nil {
		return argo.NewInvalidArgError(argo.ArgErrInvalidBindingBadType, a, "")
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

	if a.hasDef && xref.GetRootValue(reflect.ValueOf(a.def)).Kind() == reflect.Func {
		a.rootDef = xref.GetRootValue(reflect.ValueOf(a.def))
		return a.validateDefaultProvider()
	}

	if tmp, err := marsh.ToUnmarshalable("", reflect.ValueOf(a.def), true); err != nil {
		// TODO: This is not necessarily the correct error type
		return invalidDefaultValError(a)
	} else {
		a.rootDef = tmp
	}

	if a.rootDef.Kind() != reflect.String && !xref.Compatible(&a.rootDef, &a.rootBind) {
		return invalidDefaultValError(a)
	}

	return nil
}

func (a *argumentBuilder) validateDefaultProvider() error {
	root := &a.rootDef
	rType := root.Type()

	oLen := rType.NumOut()
	if oLen == 0 || oLen > 2 {
		return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultFn, a, errDefFnOutNum)
	}

	if !rType.Out(0).AssignableTo(a.rootBind.Type()) {
		// Second chance for Unmarshalable short circuit logic
		// GetRootValue
		if xref.IsUnmarshaler(a.rootBind.Type()) && rType.Out(0).AssignableTo(a.rootBind.Type().Elem()) {
			return nil
		}

		return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultVal, a,
			fmt.Sprintf(errBadType, rType.Out(0), reflect.TypeOf(a.bind)))
	}

	if oLen == 2 && !rType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultFn, a, err2ndOut)
	}

	return nil
}

func invalidDefaultValError(b *argumentBuilder) error {
	return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultVal, b,
		fmt.Sprintf(errBadType, b.rootDef.Type(), b.rootBind.Type()))
}
