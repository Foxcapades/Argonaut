package argument

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/unmarshal"
	"github.com/foxcapades/argonaut/v3/internal/xreflect"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBinding(value any) Binding {
	return Binding{value, argo.BindingTypeUnknown}
}

type Binding struct {
	Raw   any
	BType argo.ArgumentBindingType
}

func (b *Binding) Type() argo.ArgumentBindingType {
	return b.BType
}

func (b *Binding) BoundTo() any {
	return b.Raw
}

func (b *Binding) IsUsable() bool {
	// This method is only called internally by the default implementations, so we
	// know that the "Unknown" type isn't a possible value.
	return !(b.BType == argo.BindingTypeNone || b.BType == argo.BindingTypeInvalid)
}

func DetermineBindType(bind any) (kind argo.ArgumentBindingType, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			kind = argo.BindingTypeInvalid
			err = errors.New("binding must be a pointer, a consumer func, or an argo.Unmarshaler instance")
		}
	}()

	rt := reflect.TypeOf(bind)
	rk := rt.Kind()

	switch rk {
	case reflect.Ptr:
		if unmarshal.IsUnmarshalable(rt) {
			if rt.Elem().AssignableTo(unmarshal.UnmarshalerType) {
				return argo.BindingTypeUnmarshaler, nil
			}

			return argo.BindingTypePointer, nil
		}

		return argo.BindingTypeInvalid, fmt.Errorf("binding is a pointer to a type that cannot be unmarshalled: %T", bind)

	case reflect.Func:
		if unmarshal.IsUnmarshalable(rt) {
			if xreflect.FuncHasReturn(rt) {
				return argo.BindingTypeErrorFunc, nil
			}

			return argo.BindingTypeSimpleFunc, nil
		}

		return argo.BindingTypeInvalid, fmt.Errorf("binding is invalid function type %s", rt)

	default:
		return argo.BindingTypeInvalid, fmt.Errorf("invalid binding kind: %s", rk)
	}
}

// ///////////////////////
