package argument

import (
	"errors"
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/unmarshal"
	"github.com/foxcapades/argonaut/v3/internal/xarg"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.ArgumentBuilder {
	return &argumentBuilder{
		marsh: parse.NewDefaultMagicUnmarshaler(),
	}
}

type argumentBuilder struct {
	name string
	desc string

	required bool

	bindKind    xarg.BindKind
	defaultKind xarg.DefaultKind

	def  any
	bind any

	rootDef  reflect.Value
	rootBind reflect.Value

	marsh argo.ValueUnmarshaler

	validators []any

	errors []error
}

func (a *argumentBuilder) WithName(name string) argo.ArgumentBuilder {
	a.name = name
	return a
}

func (a *argumentBuilder) WithDescription(desc string) argo.ArgumentBuilder {
	a.desc = desc
	return a
}

func (a *argumentBuilder) Require() argo.ArgumentBuilder {
	a.required = true
	return a
}

func (a *argumentBuilder) isRequired() bool {
	return a.required
}

func (a *argumentBuilder) WithBinding(binding any) argo.ArgumentBuilder {
	a.bindKind = xarg.BindKindUnknown
	a.bind = binding
	return a
}

func (a *argumentBuilder) getBinding() any {
	return a.bind
}

func (a *argumentBuilder) WithDefault(def any) argo.ArgumentBuilder {
	a.defaultKind = xarg.DefaultKindUnknown
	a.def = def
	return a
}

func (a *argumentBuilder) getDefault() any {
	return a.def
}

func (a *argumentBuilder) WithUnmarshaler(fn argo.ValueUnmarshaler) argo.ArgumentBuilder {
	a.marsh = fn
	return a
}

func (a *argumentBuilder) WithValidator(fn any) argo.ArgumentBuilder {
	a.validators = append(a.validators, fn)
	return a
}

func (a *argumentBuilder) Build(warnings *argo.WarningContext) (argo.Argument, error) {
	errs := xerr.NewMultiError()

	if a.bindKind != xarg.BindKindNone {
		kind, err := xarg.DetermineBindKind(a.bind, unmarshalerType)
		a.bindKind = kind
		if err != nil {
			errs.AppendError(newArgumentBindingError(err, a))
		} else {
			a.rootBind = unmarshal.GetRootValue(reflect.ValueOf(a.bind), unmarshalerType)
		}
	}

	if a.defaultKind != xarg.DefaultKindNone {
		if a.bindKind == xarg.BindKindNone {
			errs.AppendError(errors.New("default value set with no binding"))
		} else if a.bindKind != xarg.BindKindInvalid {
			kind, err := xarg.DetermineDefaultKind(a.bind, a.def)
			a.defaultKind = kind
			if err != nil {
				errs.AppendError(newArgumentBindingError(err, a))
			} else {
				a.rootDef = reflect.ValueOf(a.def)
			}
		}
	}

	var pre, post []any
	var err error
	pre, post, err = xarg.SiftValidators(a.validators, &a.rootBind, a.bindKind)
	if err != nil {
		errs.AppendError(err)
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &argument{
		warnings:            warnings,
		name:                a.name,
		desc:                a.desc,
		required:            a.required,
		bindingKind:         a.bindKind,
		defaultKind:         a.defaultKind,
		bindVal:             a.bind,
		defVal:              a.def,
		rootBind:            a.rootBind,
		rootDef:             a.rootDef,
		unmarshal:           a.marsh,
		preParseValidators:  pre,
		postParseValidators: post,
	}, nil
}
