package unmarshal_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Foxcapades/Argonaut/internal/unmarshal"
)

func TestGetRootValue_happyPath01(t *testing.T) {
	foo := 69
	bar := &foo
	fizz := &bar
	buzz := &fizz

	val := unmarshal.GetRootValue(reflect.ValueOf(buzz), unmarshalerType)

	if val.Kind() != reflect.Int {
		t.Errorf("expected kind to be int but was %s", val.Kind())
	}

	if val.Interface().(int) != 69 {
		t.Error("expected value to equal input but it didn't")
	}
}

func TestGetRootValue_happyPath02(t *testing.T) {
	foo := 69

	val := unmarshal.GetRootValue(reflect.ValueOf(foo), unmarshalerType)

	if val.Kind() != reflect.Int {
		t.Errorf("expected kind to be int but was %s", val.Kind())
	}

	if val.Interface().(int) != 69 {
		t.Error("expected value to equal input but it didn't")
	}
}

func TestGetRootValue_randomStruct(t *testing.T) {
	foo := struct{ foo string }{foo: "string"}

	val := unmarshal.GetRootValue(reflect.ValueOf(foo), unmarshalerType)

	if val.Kind() != reflect.Struct {
		t.Errorf("expected kind to be int but was %s", val.Kind())
	}

	if val.Interface().(struct{ foo string }).foo != "string" {
		t.Error("expected value to equal input but it didn't")
	}
}

func TestGetRootValue_randomInterface(t *testing.T) {
	var foo any

	val := unmarshal.GetRootValue(reflect.ValueOf(&foo), unmarshalerType)

	if val.Kind() != reflect.Interface {
		t.Errorf("expected kind to be int but was %s", val.Kind())
	}
}

func TestGetRootValue_randomInterfaceString(t *testing.T) {
	var foo any = "foo"

	val := unmarshal.GetRootValue(reflect.ValueOf(&foo), unmarshalerType)

	if val.Kind() != reflect.Interface {
		t.Errorf("expected kind to be int but was %s", val.Kind())
	}
}

type nmrshlr struct{}

func (nmrshlr) Unmarshal(string) error { return nil }

func TestToUnmarshalable_withUnmarshalerStruct(t *testing.T) {
	value := nmrshlr{}
	_, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(value), false, unmarshalerType)

	if err == nil {
		t.Fail()
	}
}

func TestToUnmarshalable_withUnmarshalerPointer(t *testing.T) {
	value := nmrshlr{}
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withIntPointer(t *testing.T) {
	value := 3
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withStringSlice(t *testing.T) {
	var value []string
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withStringPointerSlice(t *testing.T) {
	var value []*string
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withByteSlice(t *testing.T) {
	var value []byte
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withByteSliceSlice(t *testing.T) {
	var value [][]byte
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withByteSlicePointerSlice(t *testing.T) {
	var value []*[]byte
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withStringStringMap(t *testing.T) {
	var value map[string]string
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withStringStringPointerMap(t *testing.T) {
	var value map[string]*string
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withTime(t *testing.T) {
	var value time.Time
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestToUnmarshalable_withInterface(t *testing.T) {
	var value any
	val, err := unmarshal.ToUnmarshalable("hello", reflect.ValueOf(&value), false, unmarshalerType)

	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(val)
	}
}

func TestValidateContainerValue_withMapOfStringToStringPointers(t *testing.T) {
	var value map[string]*string
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateContainerValue_withMapOfStringToStringPointerPointers(t *testing.T) {
	var value map[string]**string
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestValidateContainerValue_withMapOfStringToInterface(t *testing.T) {
	var value map[string]any
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestValidateContainerValue_withMapOfStringToInterfacePointer(t *testing.T) {
	var value map[string]*any
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestValidateContainerValue_withMapOfStringToUnmarshaler(t *testing.T) {
	var value map[string]nmrshlr
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestValidateContainerValue_withMapOfStringToUnmarshalerPointer(t *testing.T) {
	var value map[string]*nmrshlr
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestValidateContainerValue_withMapOfStringToStringSlice(t *testing.T) {
	var value map[string][]string
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestValidateContainerValue_withMapOfStringToStringSlicePointer(t *testing.T) {
	var value map[string]*[]string
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestValidateContainerValue_withMapOfStringToAnonymousStruct(t *testing.T) {
	var value map[string]struct{ foo string }
	var vType = reflect.TypeOf(value).Elem()
	err := unmarshal.ValidateContainerValue(vType, reflect.ValueOf(value), unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestToValidMap_withInvalidKey(t *testing.T) {
	var value map[struct{ foo string }]string
	var vType = reflect.ValueOf(value)
	_, err := unmarshal.ToValidMap(vType, vType, unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestToValidMap_withInvalidValue(t *testing.T) {
	var value map[string]struct{ foo string }
	var vType = reflect.ValueOf(value)
	_, err := unmarshal.ToValidMap(vType, vType, unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestToValidSlice_withInvalidValue(t *testing.T) {
	var value []struct{ foo string }
	var vType = reflect.ValueOf(value)
	_, err := unmarshal.ToValidSlice(vType, vType, unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestToUnmarshalable_withInvalidValue(t *testing.T) {
	var value struct{ foo string }
	var vType = reflect.ValueOf(value)
	_, err := unmarshal.ToUnmarshalable("", vType, true, unmarshalerType)
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestGetRootValue_withNilValue(t *testing.T) {
	var foo *string
	var bar = &foo
	var fv = reflect.ValueOf(&bar)

	val := unmarshal.GetRootValue(fv, unmarshalerType)

	if val.Kind() != reflect.String {
		t.Fail()
	}
}
