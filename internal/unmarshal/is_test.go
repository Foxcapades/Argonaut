package unmarshal_test

import (
	"reflect"
	"testing"

	"github.com/Foxcapades/Argonaut/internal/unmarshal"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

var unmarshalerType = reflect.TypeOf((*argo.Unmarshaler)(nil)).Elem()

func TestIsUnmarshalable01(t *testing.T) {
	var foo int
	if unmarshal.IsUnmarshalable(reflect.TypeOf(foo), unmarshalerType) {
		t.Fail()
	}
}

func TestIsUnmarshalable02(t *testing.T) {
	var foo int
	if !unmarshal.IsUnmarshalable(reflect.TypeOf(&foo), unmarshalerType) {
		t.Fail()
	}
}

func TestIsUnmarshalable03(t *testing.T) {
	if !unmarshal.IsUnmarshalable(unmarshalerType, unmarshalerType) {
		t.Fail()
	}
}

func TestIsUnmarshalable04(t *testing.T) {
	var foo = func(int) {}
	if !unmarshal.IsUnmarshalable(reflect.TypeOf(foo), unmarshalerType) {
		t.Fail()
	}
}

func TestIsUnmarshalable05(t *testing.T) {
	if unmarshal.IsUnmarshalable(reflect.TypeOf(any(nil)), unmarshalerType) {
		t.Fail()
	}
}

func TestIsUnmarshalable06(t *testing.T) {
	var foo map[string]string
	if !unmarshal.IsUnmarshalable(reflect.TypeOf(&foo), unmarshalerType) {
		t.Fail()
	}
}

func TestIsUnmarshalable07(t *testing.T) {
	var foo []string
	if !unmarshal.IsUnmarshalable(reflect.TypeOf(&foo), unmarshalerType) {
		t.Fail()
	}
}

func TestIsConsumerFunc01(t *testing.T) {
	foo := func(string) error { return nil }
	if !unmarshal.IsConsumerFunc(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsConsumerFunc02(t *testing.T) {
	foo := func(string) {}
	if !unmarshal.IsConsumerFunc(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsConsumerFunc03(t *testing.T) {
	foo := func(string) int { return 0 }
	if unmarshal.IsConsumerFunc(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsConsumerFunc04(t *testing.T) {
	foo := func() {}
	if unmarshal.IsConsumerFunc(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsConsumerFunc05(t *testing.T) {
	foo := func(string) (error, error) { return nil, nil }
	if unmarshal.IsConsumerFunc(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsConsumerFunc06(t *testing.T) {
	foo := func(complex128) {}
	if unmarshal.IsConsumerFunc(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsUnmarshalableValue(t *testing.T) {
	var foo map[string]string
	if !unmarshal.IsUnmarshalableValue(reflect.TypeOf(foo)) {
		t.Fail()
	}
}
