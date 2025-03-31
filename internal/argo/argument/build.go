package argument

import (
	"errors"
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(a argo.ArgumentBuilder) (argo.Argument, error) {
	errs := xerr.NewMultiError()

	var err error

	var bind Binding
	if a.HasBinding() {
		bind.Raw = a.Binding().BoundTo()
		bind.BType, err = DetermineBindType(bind.Raw)

		if err != nil {
			errs.AppendError(NewBindingError(err, a))
		}
	}

	var def Default
	if a.HasDefault() {
		if bind.BType == argo.BindingTypeNone {
			errs.AppendError(errors.New("default value set with no binding"))
		} else if bind.BType != argo.BindingTypeInvalid {
			def.DType, err = DetermineDefaultType(a.Binding().BoundTo(), a.Default().Value())

			if err != nil {
				errs.AppendError(NewBindingError(err, a))
			} else {
				def.Raw = a.Default().Value()
			}
		}
	}

	var pre, post []any
	pre, post, err = SiftValidators(a.Validators(), &bind)
	if err != nil {
		errs.AppendError(err)
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &Argument{
		name:                a.Name(),
		desc:                a.Description(),
		required:            a.IsRequired(),
		binding:             bind,
		defVal:              def,
		unmarshal:           utils.CallIfElse(a.HasUnmarshaler(), a.Unmarshaler, parse.NewDefaultMagicUnmarshaler),
		preParseValidators:  pre,
		postParseValidators: post,
	}, nil
}

func BuildMulti(builders []argo.ArgumentBuilder, errs argo.MultiError) []argo.Argument {
	lastRequired := 0
	for i, argBuilder := range builders {
		if argBuilder.IsRequired() {
			lastRequired = i
		}
	}

	out := make([]argo.Argument, 0, len(builders))

	for i, builder := range builders {
		if i < lastRequired && !builder.IsRequired() {
			errs.AppendError(fmt.Errorf("argument %d was not marked as required, but preceded required argument %d", i+1, lastRequired+1))
		}

		if arg, ok := xerr.TryBuild(builder, Build, errs); ok {
			out = append(out, arg)
		}
	}

	return out
}
