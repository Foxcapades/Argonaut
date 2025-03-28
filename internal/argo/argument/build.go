package argument

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(a argo.ArgumentBuilder) (argo.Argument, error) {
	errs := xerr.NewMultiError()

	var err error

	var bind Binding
	if a.HasBinding() {
		bind.raw = a.Binding().BoundTo()
		bind.bType, err = DetermineBindType(bind.raw)

		if err != nil {
			errs.AppendError(NewBindingError(err, a))
		}
	}

	var def Default
	if a.HasDefault() {
		if bind.bType == argo.BindingTypeNone {
			errs.AppendError(errors.New("default value set with no binding"))
		} else if bind.bType != argo.BindingTypeInvalid {
			def.dType, err = DetermineDefaultType(a.Binding().BoundTo(), a.Default().Value())

			if err != nil {
				errs.AppendError(NewBindingError(err, a))
			} else {
				def.value = a.Default().Value()
			}
		}
	}

	var pre, post []any
	pre, post, err = SiftValidators(a.Validators(), &def, &bind)
	if err != nil {
		errs.AppendError(err)
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &argument{
		name:                a.Name(),
		desc:                a.Description(),
		required:            a.IsRequired(),
		binding:             bind,
		defVal:              def,
		unmarshal:           a.Unmarshaler(),
		preParseValidators:  pre,
		postParseValidators: post,
	}, nil
}

func BuildMulti(builders []argo.ArgumentBuilder, errs argo.MultiError) []argo.Argument {
	forceRequiredUntil := 0
	for i, b := range builders {
		if b.IsRequired() {
			forceRequiredUntil = i
		}
	}

	out := make([]argo.Argument, 0, len(builders))

	for i, builder := range builders {
		if i < forceRequiredUntil && !builder.IsRequired() {
			builder.Require()
		}

		if arg, ok := xerr.TryBuild(builder, Build, errs); ok {
			out = append(out, arg)
		}
	}

	return out
}
