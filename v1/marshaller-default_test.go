package argo_test

import (
	"reflect"
	"testing"

	C "github.com/smartystreets/goconvey/convey"

	"github.com/Foxcapades/Argonaut/v1"
)

func TestUnmarshal_ValidInts(t *testing.T) {

	root := struct {
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64
	}{}

	zero := func() {
		root.i = 0
		root.i8 = 0
		root.i16 = 0
		root.i32 = 0
		root.i64 = 0
	}

	type test struct {
		Name   string
		Input  string
		Output interface{}
		Temp   interface{}
	}

	validIntParseTest := []test{
		//
		// Base 10
		//

		// to `int`
		{"Positive base 10 int", "123", 123, &root.i},
		{"Negative base 10 int", "-123", -123, &root.i},
		// to `int8`
		//{"Positive base 10 int8", "123", int8(123), int8(0)},
		//{"Negative base 10 int8", "-123", int8(-123), int8(0)},
		//// to `int16`
		//{"Positive base 10 int16", "123", int16(123), int16(0)},
		//{"Negative base 10 int16", "-123", int16(-123), int16(0)},
		//// to `int32`
		//{"Positive base 10 int32", "123", int32(123), int32(0)},
		//{"Negative base 10 int32", "-123", int32(-123), int32(0)},
		//// to `int64`
		//{"Positive base 10 int64", "123", int64(123), int64(0)},
		//{"Negative base 10 int64", "-123", int64(-123), int64(0)},

		// to `uint`
		//{"Positive base 10 untyped uint", "123", uint(123), uint(0)},
		//// to `uint8`
		//{"Positive base 10 uint8", "123", uint8(123), uint8(0)},
		//// to `uint16`
		//{"Positive base 10 uint16", "123", uint16(123), uint16(0)},
		//// to `uint32`
		//{"Positive base 10 uint32", "123", uint32(123), uint32(0)},
		//// to `uint64`
		//{"Positive base 10 uint64", "123", uint64(123), uint64(0)},

		//
		// Base 16
		//

		// to `int`
		//{"Positive base 16 with prefix `0x` to untyped int", "0x123", 291, 0},
		//{"Negative base 16 with prefix `0x` to untyped int", "-0x123", -291, 0},
		//{"Positive base 16 with prefix `0X` to untyped int", "0X123", 291, 0},
		//{"Negative base 16 with prefix `0X` to untyped int", "-0X123", -291, 0},
		//{"Positive base 16 with prefix `x` to untyped int",  "x123", 291, 0},
		//{"Negative base 16 with prefix `x` to untyped int",  "-x123", -291, 0},
		//{"Positive base 16 with prefix `X` to untyped int",  "X123", 291, 0},
		//{"Negative base 16 with prefix `X` to untyped int",  "-X123", -291, 0},

		// to `int8`
		//{"Positive base 16 with prefix `0x` to int8", "0x7E", int8(126), int8(0)},
		//{"Negative base 16 with prefix `0x` to int8", "-0x7E", int8(-126), int8(0)},
		//{"Positive base 16 with prefix `0X` to int8", "0X7E", int8(126), int8(0)},
		//{"Negative base 16 with prefix `0X` to int8", "-0X7E", int8(-126), int8(0)},
		//{"Positive base 16 with prefix `x` to int8",  "x7E", int8(126), int8(0)},
		//{"Negative base 16 with prefix `x` to int8",  "-x7E", int8(-126), int8(0)},
		//{"Positive base 16 with prefix `X` to int8",  "X7E", int8(126), int8(0)},
		//{"Negative base 16 with prefix `X` to int8",  "-X7E", int8(-126), int8(0)},

		// to `int16`
		//{"Positive base 16 with prefix `0x` to int16", "0x7E", int16(126), int16(0)},
		//{"Negative base 16 with prefix `0x` to int16", "-0x7E", int16(-126), int16(0)},
		//{"Positive base 16 with prefix `0X` to int16", "0X7E", int16(126), int16(0)},
		//{"Negative base 16 with prefix `0X` to int16", "-0X7E", int16(-126), int16(0)},
		//{"Positive base 16 with prefix `x` to int16",  "x7E", int16(126), int16(0)},
		//{"Negative base 16 with prefix `x` to int16",  "-x7E", int16(-126), int16(0)},
		//{"Positive base 16 with prefix `X` to int16",  "X7E", int16(126), int16(0)},
		//{"Negative base 16 with prefix `X` to int16",  "-X7E", int16(-126), int16(0)},

		// to `int32`
		//{"Positive base 16 with prefix `0x` to int32", "0x7E", int32(126), int32(0)},
		//{"Negative base 16 with prefix `0x` to int32", "-0x7E", int32(-126), int32(0)},
		//{"Positive base 16 with prefix `0X` to int32", "0X7E", int32(126), int32(0)},
		//{"Negative base 16 with prefix `0X` to int32", "-0X7E", int32(-126), int32(0)},
		//{"Positive base 16 with prefix `x` to int32",  "x7E", int32(126), int32(0)},
		//{"Negative base 16 with prefix `x` to int32",  "-x7E", int32(-126), int32(0)},
		//{"Positive base 16 with prefix `X` to int32",  "X7E", int32(126), int32(0)},
		//{"Negative base 16 with prefix `X` to int32",  "-X7E", int32(-126), int32(0)},

		// to `int64`
		//{"Positive base 16 with prefix `0x` to int64", "0x7E", int64(126), int64(0)},
		//{"Negative base 16 with prefix `0x` to int64", "-0x7E", int64(-126), int64(0)},
		//{"Positive base 16 with prefix `0X` to int64", "0X7E", int64(126), int64(0)},
		//{"Negative base 16 with prefix `0X` to int64", "-0X7E", int64(-126), int64(0)},
		//{"Positive base 16 with prefix `x` to int64",  "x7E", int64(126), int64(0)},
		//{"Negative base 16 with prefix `x` to int64",  "-x7E", int64(-126), int64(0)},
		//{"Positive base 16 with prefix `X` to int64",  "X7E", int64(126), int64(0)},
		//{"Negative base 16 with prefix `X` to int64",  "-X7E", int64(-126), int64(0)},

		// Base 8

		// Base 2
	}

	C.Convey("Unmarshal", t, func() {
		p := argo.DefaultUnmarshalProps()
		C.Convey("Integer Types", func() {
			for i := range validIntParseTest {
				tmp := &validIntParseTest[i]

				C.Convey(tmp.Name, func() {
					C.So(argo.Unmarshal(tmp.Input, &tmp.Temp, p), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Int(), C.ShouldEqual, tmp.Output)
				})

				// Reset test container
				zero()
			}
		})
	})
}
