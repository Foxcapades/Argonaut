package argument

import (
	"errors"
	"fmt"
	R "reflect"

	"github.com/Foxcapades/Argonaut/v0/internal/util"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

const (
	errDefFnOutNum = "default value providers must return either 1 or 2 values"
	err2ndOut      = "the second output type of a default value provider must " +
		"be compatible with error"
	errBadType = "default value type %s is not compatible with binding type %s"
)

func (a *Builder) ValidateDefault() error {
	if !a.IsDefaultSet {
		return nil
	}

	if !a.IsBindingSet {
		// TODO: this should be a real error
		return errors.New("default set with no binding")
	}

	if a.HasDefaultProvider() {
		a.RootDefault = util.GetRootValue(R.ValueOf(a.DefaultValue))
		return a.ValidateDefaultProvider()
	}

	if tmp, err := util.ToUnmarshalable("", R.ValueOf(a.DefaultValue), true); err != nil {
		// TODO: This is not necessarily the correct error type
		return InvalidDefaultValError(a)
	} else {
		a.RootDefault = tmp
	}

	if a.RootDefault.Kind() != R.String && !util.Compatible(&a.RootDefault, &a.RootBinding) {
		return InvalidDefaultValError(a)
	}

	return nil
}

func (a *Builder) ValidateDefaultProvider() error {
	root := &a.RootDefault
	rType := root.Type()

	oLen := rType.NumOut()
	if oLen == 0 || oLen > 2 {
		return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultFn, a, errDefFnOutNum)
	}

	if !rType.Out(0).AssignableTo(a.RootBinding.Type()) {
		// Second chance for Unmarshalable short circuit logic
		// GetRootValue
		if util.IsUnmarshaler(a.RootBinding.Type()) && rType.Out(0).AssignableTo(a.RootBinding.Type().Elem()) {
			return nil
		}
		return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultVal, a,
			fmt.Sprintf(errBadType, rType.Out(0), R.TypeOf(a.BindValue)))
	}

	if oLen == 2 && !rType.Out(1).AssignableTo(R.TypeOf((*error)(nil)).Elem()) {
		return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultFn, a, err2ndOut)
	}

	return nil
}

func InvalidDefaultValError(b *Builder) error {
	return argo.NewInvalidArgError(argo.ArgErrInvalidDefaultVal, b,
		fmt.Sprintf(errBadType, b.RootDefault.Type(), b.RootBinding.Type()))
}
