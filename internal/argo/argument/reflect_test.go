package argument_test

import (
	"reflect"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestSiftValidators01(t *testing.T) {
	validators := []any{
		func(int, string) error { return nil },
		func(string) error { return nil },
	}

	bind := 0
	bt := argument.Binding{Raw: &bind, BType: argo.BindingTypePointer}

	a, b, e := argument.SiftValidators(validators, &bt)

	if len(a) != 1 {
		t.Error("expected pre-parse validator slice to have a length of 1")
	}
	if len(b) != 1 {
		t.Error("expected post-parse validator slice to have a length of 1")
	}
	if e != nil {
		t.Error("expected err to be nil")
		t.Log(e)
	}
}

func TestSiftValidators02(t *testing.T) {
	bind := ""
	validators := []any{
		func(int, string) error { return nil },
		func(string) error { return nil },
	}

	bt := argument.Binding{Raw: &bind, BType: argo.BindingTypePointer}

	a, b, e := argument.SiftValidators(validators, &bt)

	if a != nil {
		t.Error("expected pre-parse validator slice be nil")
	}
	if b != nil {
		t.Error("expected post-parse validator slice to be nil")
	}
	if e == nil {
		t.Error("expected err not to be nil")
	} else {
		t.Log(e)
	}
}

func TestSiftValidators03(t *testing.T) {
	validators := []any{
		func(int, string) error { return nil },
		func(string) error { return nil },
	}
	bind := 0
	bt := argument.Binding{Raw: &bind, BType: argo.BindingTypePointer}

	a, b, e := argument.SiftValidators(validators, &bt)

	if len(a) != 1 {
		t.Error("expected pre-parse validator slice to have a length of 1")
	}
	if len(b) != 1 {
		t.Error("expected post-parse validator slice to have a length of 1")
	}
	if e != nil {
		t.Error("expected err to be nil")
		t.Log(e)
	}
}

func TestValidateValidator01(t *testing.T) {
	validator := func(string) error { return nil }
	bindType := reflect.TypeOf(666)
	count, err := argument.ValidateValidator(10, validator, bindType, true)

	if count != 1 {
		t.Errorf("expected count to equal 1 but was %d", count)
	}
	if err != nil {
		t.Error("expected err to be nil but it wasn't")
		t.Log(err)
	}
}

func TestValidateValidator02(t *testing.T) {
	validator := "hello"
	bindType := reflect.TypeOf(666)
	count, err := argument.ValidateValidator(10, validator, bindType, true)

	if count != 0 {
		t.Errorf("expected count to equal 0 but was %d", count)
	}
	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestValidateValidator03(t *testing.T) {
	validator := func() int { return 0 }
	bindType := reflect.TypeOf(666)
	count, err := argument.ValidateValidator(10, validator, bindType, true)

	if count != 0 {
		t.Errorf("expected count to equal 0 but was %d", count)
	}
	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestValidateValidator04(t *testing.T) {
	validator := func() (int, int) { return 0, 0 }
	bindType := reflect.TypeOf(666)
	count, err := argument.ValidateValidator(10, validator, bindType, true)

	if count != 0 {
		t.Errorf("expected count to equal 0 but was %d", count)
	}
	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestValidateValidator05(t *testing.T) {
	validator := func(a, b, c int) error { return nil }
	bindType := reflect.TypeOf(666)
	count, err := argument.ValidateValidator(10, validator, bindType, true)

	if count != 0 {
		t.Errorf("expected count to equal 0 but was %d", count)
	}
	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestValidateSoloValidator_failsOnNonStringParam(t *testing.T) {
	validator := func(int) {}
	err := argument.ValidateSoloValidator(0, reflect.TypeOf(validator))

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestValidateDoubleValidator_failsOnNonStringSecondParam(t *testing.T) {
	validator := func(int, int) {}
	binding := 0
	err := argument.ValidateDoubleValidator(0, reflect.TypeOf(validator), reflect.TypeOf(binding))

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}
