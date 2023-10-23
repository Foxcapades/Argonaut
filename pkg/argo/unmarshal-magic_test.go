package argo_test

import (
	"testing"
	"time"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestMagicUnmarshaler001(t *testing.T) {
	un := argo.NewDefaultMagicUnmarshaler()

	v001 := 0
	must(un.Unmarshal("1", &v001))
	if v001 != 1 {
		t.Fail()
	}

	v002 := "hello"
	must(un.Unmarshal("goodbye", &v002))
	if v002 != "goodbye" {
		t.Fail()
	}

	v003 := float32(0)
	must(un.Unmarshal("3.3", &v003))
	if v003 != 3.3 {
		t.Fail()
	}

	v004 := int64(0)
	must(un.Unmarshal("12345", &v004))
	if v004 != 12345 {
		t.Fail()
	}

	v005 := 0.0
	must(un.Unmarshal("4.4", &v005))
	if v005 != 4.4 {
		t.Fail()
	}

	v006 := uint64(0)
	must(un.Unmarshal("123456789123456789", &v006))
	if v006 != 123456789123456789 {
		t.Fail()
	}

	v007 := uint(0)
	must(un.Unmarshal("123456789", &v007))
	if v007 != 123456789 {
		t.Fail()
	}

	v008 := int32(0)
	must(un.Unmarshal("-1234", &v008))
	if v008 != -1234 {
		t.Fail()
	}

	v009 := uint32(0)
	must(un.Unmarshal("01234", &v009))
	if v009 != 01234 {
		t.Fail()
	}

	v010 := uint16(0)
	must(un.Unmarshal("x0000123", &v010))
	if v010 != 0x0000123 {
		t.Fail()
	}

	v011 := uint8(0)
	must(un.Unmarshal("255", &v011))
	if v011 != 255 {
		t.Fail()
	}

	v012 := int8(0)
	must(un.Unmarshal("-128", &v012))
	if v012 != -128 {
		t.Fail()
	}

	v013 := int16(0)
	must(un.Unmarshal("4321", &v013))
	if v013 != 4321 {
		t.Fail()
	}

	for _, val := range []string{"true", "t", "yes", "y", "on", "1"} {
		v014 := false
		must(un.Unmarshal(val, &v014))
		if !v014 {
			t.Fail()
		}
	}

	for _, val := range []string{"false", "f", "no", "n", "off", "0"} {
		v015 := true
		must(un.Unmarshal(val, &v015))
		if v015 {
			t.Fail()
		}
	}

	v016 := time.Duration(0)
	must(un.Unmarshal("10m", &v016))
	if v016.String() != "10m0s" {
		t.Fail()
	}

	var v017 time.Time
	must(un.Unmarshal("2023-10-23T15:13:43.1234-04:00", &v017))
	if v017.Format(time.RFC3339Nano) != "2023-10-23T15:13:43.1234-04:00" {
		t.Fail()
	}
}

func TestMagicUnmarshaler_MapOfSlice(t *testing.T) {
	un := argo.NewDefaultMagicUnmarshaler()
	var foo map[string][]string

	must(un.Unmarshal("foo:bar,foo:fizz,foo:buzz,fizz:buzz", &foo))

	if slice, ok := foo["foo"]; !ok {
		t.Error("expected key was not found in unmarshalled map")
	} else {
		if len(slice) != 3 {
			t.Error("expected slice to have 3 elements but it didn't")
		} else {
			if slice[0] != "bar" || slice[1] != "fizz" || slice[2] != "buzz" {
				t.Error("expected slice to contain input values but it didn't")
			}
		}
	}

	if slice, ok := foo["fizz"]; !ok {
		t.Error("expected key was not found in unmarshalled map")
	} else {
		if len(slice) != 1 {
			t.Error("expected slice to have 1 element but it didn't")
		} else {
			if slice[0] != "buzz" {
				t.Error("expected slice to contain input value but it didn't")
			}
		}
	}
}

func TestMagicUnmarshaler_MapOfByteSlice(t *testing.T) {
	un := argo.NewDefaultMagicUnmarshaler()
	var foo map[string][]byte

	must(un.Unmarshal("foo:bar,fizz:buzz", &foo))

	if bytes, ok := foo["foo"]; !ok {
		t.Error("expected key was not found in unmarshalled map")
	} else {
		if string(bytes) != "bar" {
			t.Error("expected byte slice to match input value but it didn't")
		}
	}

	if bytes, ok := foo["fizz"]; !ok {
		t.Error("expected key was not found in unmarshalled map")
	} else {
		if string(bytes) != "buzz" {
			t.Error("expected byte slice to match input value but it didn't")
		}
	}
}

func TestMagicUnmarshaler_MapOfByteSlicePointer(t *testing.T) {
	un := argo.NewDefaultMagicUnmarshaler()
	var foo map[string]*[]byte

	must(un.Unmarshal("foo:bar,fizz:buzz", &foo))

	if bytes, ok := foo["foo"]; !ok {
		t.Error("expected key was not found in unmarshalled map")
	} else {
		if string(*bytes) != "bar" {
			t.Error("expected byte slice to match input value but it didn't")
		}
	}

	if bytes, ok := foo["fizz"]; !ok {
		t.Error("expected key was not found in unmarshalled map")
	} else {
		if string(*bytes) != "buzz" {
			t.Error("expected byte slice to match input value but it didn't")
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
