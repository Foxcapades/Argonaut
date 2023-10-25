package argo_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestOctal_Unmarshal_prefix0(t *testing.T) {
	var foo argo.Octal
	err := foo.Unmarshal("01234")

	if err != nil {
		t.Error(err)
	}

	if foo != 668 {
		t.Fail()
	}
}

func TestOctal_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Octal
	err := foo.Unmarshal("1234")

	if err != nil {
		t.Error(err)
	}

	if foo != 668 {
		t.Fail()
	}
}

func TestOctal8_Unmarshal_prefix0(t *testing.T) {
	var foo argo.Octal8
	err := foo.Unmarshal("0177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal8_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Octal8
	err := foo.Unmarshal("177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal16_Unmarshal_prefix0(t *testing.T) {
	var foo argo.Octal16
	err := foo.Unmarshal("0177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal16_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Octal16
	err := foo.Unmarshal("177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal32_Unmarshal_prefix0(t *testing.T) {
	var foo argo.Octal32
	err := foo.Unmarshal("0177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal32_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Octal32
	err := foo.Unmarshal("177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal64_Unmarshal_prefix0(t *testing.T) {
	var foo argo.Octal64
	err := foo.Unmarshal("0177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestOctal64_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Octal64
	err := foo.Unmarshal("177")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}
