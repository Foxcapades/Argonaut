package arg

import (
	"fmt"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	R "reflect"
)

func checkDefault(a A.ArgumentBuilder) error {
	return checkDefVal(a, a.GetDefault(), a.GetBinding(), false)
}

const (
	errDefFnOutNum = "default value providers must return either 1 or 2 values"
	err2ndOut      = "the second output type of a default value provider must " +
		"be compatible with error"
	errBadType = "default value type %s is not compatible with binding type %s"
)

var errType = R.TypeOf((*error)(nil)).Elem()

func checkDefVal(a A.ArgumentBuilder, val, test interface{}, inFn bool) error {
	if vt := R.TypeOf(val); vt.Kind() == R.Func && !inFn {
		return checkDefFn(a, vt, test)
	} else if !util.Compatible(val, test) {
		return A.NewInvalidArgError(A.ArgErrInvalidDefaultVal, a,
			fmt.Sprintf(errBadType, vt, R.TypeOf(test)))
	}
	return nil
}

func checkDefFn(a A.ArgumentBuilder, fnT R.Type, test interface{}) error {
	outNum := fnT.NumOut()

	switch outNum {
	case 2:
		if !fnT.Out(1).AssignableTo(errType) {
			return A.NewInvalidArgError(A.ArgErrInvalidDefaultFn, a, err2ndOut)
		}
		fallthrough
	case 1:
		return checkDefVal(a, fnT.Out(0), test, true)
	}

	return A.NewInvalidArgError(A.ArgErrInvalidDefaultFn, a, errDefFnOutNum)
}
