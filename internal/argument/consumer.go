package argument

import (
	"fmt"
	"reflect"

	"github.com/foxcapades/argonaut/internal/errs"
	"github.com/foxcapades/argonaut/internal/util/xreflect"
	"github.com/foxcapades/argonaut/pkg/argo"
)

// region Pointer Consumer

func NewPointerConsumer[T any](pointer *T) argo.ArgumentValueConsumer[T] {
	return pointerConsumer[T]{pointer}
}

type pointerConsumer[T any] struct {
	pointer *T
}

func (p pointerConsumer[T]) Accept(value T) error {
	*p.pointer = value
	return nil
}

// endregion Pointer Consumer

// region Magic Pointer Consumer

func NewMagicPointerConsumer[T any](pointer any) (argo.ArgumentValueConsumer[T], error) {
	var tmp T
	tt := reflect.TypeOf(tmp)
	pt := reflect.TypeOf(pointer)

	if tt.AssignableTo(xreflect.RootType(pt)) {
		return magicPointerConsumer[T]{xreflect.RootValue(reflect.ValueOf(pointer))}, nil
	}

	// TODO: newtype this as a binding error
	return nil, fmt.Errorf("cannot unmarshal a value of type %s into a pointer of type %s", tt.Name(), pt.Name())
}

type magicPointerConsumer[T any] struct {
	root reflect.Value
}

func (m magicPointerConsumer[T]) Accept(value T) error {
	m.root.Set(reflect.ValueOf(value))
	return nil
}

// endregion Magic Pointer Consumer

// region Multi-Consumer

func NewMultiConsumer[T any](consumers []argo.ArgumentValueConsumer[T]) argo.ArgumentValueConsumer[T] {
	return multiConsumer[T]{consumers}
}

type multiConsumer[T any] struct {
	consumers []argo.ArgumentValueConsumer[T]
}

func (m multiConsumer[T]) Accept(value T) error {
	errors := errs.NewMultiError()

	for i := range m.consumers {
		errors.AppendIfNotNil(m.consumers[i].Accept(value))
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

// endregion Multi-Consumer

// region Void Consumer

type VoidConsumer[T any] struct{}

func (VoidConsumer[T]) Accept(T) error { return nil }

// endregion Void Consumer
