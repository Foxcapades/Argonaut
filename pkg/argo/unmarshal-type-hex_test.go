package argo_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestHex_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Hex
	err := foo.Unmarshal("FF")

	if err != nil {
		t.Error(err)
	}

	if foo != 255 {
		t.Fail()
	}
}

func TestHex8_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Hex8
	err := foo.Unmarshal("7F")

	if err != nil {
		t.Error(err)
	}

	if foo != 127 {
		t.Fail()
	}
}

func TestHex16_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Hex16
	err := foo.Unmarshal("FFFF")

	if err != nil {
		t.Error(err)
	}

	if foo != -1 {
		t.Fail()
	}
}

func TestHex32_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Hex32
	err := foo.Unmarshal("7FFFFFFF")

	if err != nil {
		t.Error(err)
	}

	if foo != 2147483647 {
		t.Fail()
	}
}

func TestHex64_Unmarshal_prefixNone(t *testing.T) {
	var foo argo.Hex64
	err := foo.Unmarshal("FFFFFFFFFFFFFFFF")

	if err != nil {
		t.Error(err)
	}

	if foo != -1 {
		t.Fail()
	}
}
