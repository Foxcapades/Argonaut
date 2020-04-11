package argument

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"reflect"
)

func (a *Builder) ValidateBinding() error {
	if !a.IsBindingSet {
		return nil
	}

	if tmp, err := util.ToUnmarshalable("", reflect.ValueOf(a.BindValue), false); err != nil {
		return argo.NewInvalidArgError(argo.ArgErrInvalidBindingBadType, a, "")
	} else {
		a.RootBinding = tmp
	}

	return nil
}
