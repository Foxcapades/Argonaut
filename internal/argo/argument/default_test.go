package argument_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestDetermineDefaultType01(t *testing.T) {
	var bind int
	var def int

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType02(t *testing.T) {
	var bind int
	var def string

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeRaw {
		t.Error("expected default kind to be raw but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType03(t *testing.T) {
	var bind int
	var def float32

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultType04(t *testing.T) {
	var bind int
	def := func() int { return 0 }

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeProviderPlain {
		t.Error("expected default kind to be plain provider but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType05(t *testing.T) {
	var bind int
	def := func() (int, error) { return 0, nil }

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeProviderWithErr {
		t.Error("expected default kind to be provider with error but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType06(t *testing.T) {
	var bind int
	def := func() string { return "" }

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultType07(t *testing.T) {
	var bind int
	def := func() (int, string) { return 0, "" }

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultType08(t *testing.T) {
	var bind int
	def := func() {}

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultType09(t *testing.T) {
	var bind []int
	var def []int

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType10(t *testing.T) {
	var bind int
	var def []int

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultType11(t *testing.T) {
	var bind []string
	var def []int

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultType12(t *testing.T) {
	bind := func([]int) {}
	var def []int

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType13(t *testing.T) {
	var bind map[string]string
	var def map[string]string

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultType14(t *testing.T) {
	var bind int
	var def map[int]int

	kind, err := argument.DetermineDefaultType(bind, def)

	if kind != argo.DefaultTypeInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}
