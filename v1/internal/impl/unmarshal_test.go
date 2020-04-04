package impl_test

import (
	"github.com/Foxcapades/Argonaut/v1/internal/impl"
	"reflect"
	"testing"

	C "github.com/smartystreets/goconvey/convey"
)

type unmarshalerValueTest struct {
	Name   string
	Input  string
	Output interface{}
	Temp   interface{}
}

func TestUnmarshal(t *testing.T) {
	C.Convey("Unmarshal", t, func() {
		C.Convey("Valid int values", func() {

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

			validIntParseTest := []unmarshalerValueTest{
				//
				// Base 10
				//

				// to `int`
				{"Positive base 10 int", "123", 123, &root.i},
				{"Negative base 10 int", "-123", -123, &root.i},
				// to `int8`
				{"Positive base 10 int8", "123", int8(123), &root.i8},
				{"Negative base 10 int8", "-123", int8(-123), &root.i8},
				//// to `int16`
				{"Positive base 10 int16", "123", int16(123), &root.i16},
				{"Negative base 10 int16", "-123", int16(-123), &root.i16},
				//// to `int32`
				{"Positive base 10 int32", "123", int32(123), &root.i32},
				{"Negative base 10 int32", "-123", int32(-123), &root.i32},
				//// to `int64`
				{"Positive base 10 int64", "123", int64(123), &root.i64},
				{"Negative base 10 int64", "-123", int64(-123), &root.i64},

				//
				// Base 16
				//

				// to `int`
				{"Positive base 16 with prefix `0x` to untyped int", "0x123", 291, &root.i},
				{"Negative base 16 with prefix `0x` to untyped int", "-0x123", -291, &root.i},
				{"Positive base 16 with prefix `0X` to untyped int", "0X123", 291, &root.i},
				{"Negative base 16 with prefix `0X` to untyped int", "-0X123", -291, &root.i},
				{"Positive base 16 with prefix `x` to untyped int", "x123", 291, &root.i},
				{"Negative base 16 with prefix `x` to untyped int", "-x123", -291, &root.i},
				{"Positive base 16 with prefix `X` to untyped int", "X123", 291, &root.i},
				{"Negative base 16 with prefix `X` to untyped int", "-X123", -291, &root.i},

				// to `int8`
				{"Positive base 16 with prefix `0x` to int8", "0x7E", int8(126), &root.i8},
				{"Negative base 16 with prefix `0x` to int8", "-0x7E", int8(-126), &root.i8},
				{"Positive base 16 with prefix `0X` to int8", "0X7E", int8(126), &root.i8},
				{"Negative base 16 with prefix `0X` to int8", "-0X7E", int8(-126), &root.i8},
				{"Positive base 16 with prefix `x` to int8", "x7E", int8(126), &root.i8},
				{"Negative base 16 with prefix `x` to int8", "-x7E", int8(-126), &root.i8},
				{"Positive base 16 with prefix `X` to int8", "X7E", int8(126), &root.i8},
				{"Negative base 16 with prefix `X` to int8", "-X7E", int8(-126), &root.i8},

				// to `int16`
				{"Positive base 16 with prefix `0x` to int16", "0x7E", int16(126), &root.i16},
				{"Negative base 16 with prefix `0x` to int16", "-0x7E", int16(-126), &root.i16},
				{"Positive base 16 with prefix `0X` to int16", "0X7E", int16(126), &root.i16},
				{"Negative base 16 with prefix `0X` to int16", "-0X7E", int16(-126), &root.i16},
				{"Positive base 16 with prefix `x` to int16", "x7E", int16(126), &root.i16},
				{"Negative base 16 with prefix `x` to int16", "-x7E", int16(-126), &root.i16},
				{"Positive base 16 with prefix `X` to int16", "X7E", int16(126), &root.i16},
				{"Negative base 16 with prefix `X` to int16", "-X7E", int16(-126), &root.i16},

				// to `int32`
				{"Positive base 16 with prefix `0x` to int32", "0x7E", int32(126), &root.i32},
				{"Negative base 16 with prefix `0x` to int32", "-0x7E", int32(-126), &root.i32},
				{"Positive base 16 with prefix `0X` to int32", "0X7E", int32(126), &root.i32},
				{"Negative base 16 with prefix `0X` to int32", "-0X7E", int32(-126), &root.i32},
				{"Positive base 16 with prefix `x` to int32", "x7E", int32(126), &root.i32},
				{"Negative base 16 with prefix `x` to int32", "-x7E", int32(-126), &root.i32},
				{"Positive base 16 with prefix `X` to int32", "X7E", int32(126), &root.i32},
				{"Negative base 16 with prefix `X` to int32", "-X7E", int32(-126), &root.i32},

				// to `int64`
				{"Positive base 16 with prefix `0x` to int64", "0x7E", int64(126), &root.i64},
				{"Negative base 16 with prefix `0x` to int64", "-0x7E", int64(-126), &root.i64},
				{"Positive base 16 with prefix `0X` to int64", "0X7E", int64(126), &root.i64},
				{"Negative base 16 with prefix `0X` to int64", "-0X7E", int64(-126), &root.i64},
				{"Positive base 16 with prefix `x` to int64", "x7E", int64(126), &root.i64},
				{"Negative base 16 with prefix `x` to int64", "-x7E", int64(-126), &root.i64},
				{"Positive base 16 with prefix `X` to int64", "X7E", int64(126), &root.i64},
				{"Negative base 16 with prefix `X` to int64", "-X7E", int64(-126), &root.i64},

				//
				// Base 8
				//

				// to `int`
				{"Positive base 8 with prefix `0o` to untyped int", "0o123", 83, &root.i},
				{"Negative base 8 with prefix `0o` to untyped int", "-0o123", -83, &root.i},
				{"Positive base 8 with prefix `0O` to untyped int", "0O123", 83, &root.i},
				{"Negative base 8 with prefix `0O` to untyped int", "-0O123", -83, &root.i},
				{"Positive base 8 with prefix `o` to untyped int", "o123", 83, &root.i},
				{"Negative base 8 with prefix `o` to untyped int", "-o123", -83, &root.i},
				{"Positive base 8 with prefix `O` to untyped int", "O123", 83, &root.i},
				{"Negative base 8 with prefix `O` to untyped int", "-O123", -83, &root.i},

				// to `int8`
				{"Positive base 8 with prefix `0o` to int8", "0o133", int8(91), &root.i8},
				{"Negative base 8 with prefix `0o` to int8", "-0o133", int8(-91), &root.i8},
				{"Positive base 8 with prefix `0O` to int8", "0O133", int8(91), &root.i8},
				{"Negative base 8 with prefix `0O` to int8", "-0O133", int8(-91), &root.i8},
				{"Positive base 8 with prefix `o` to int8", "o133", int8(91), &root.i8},
				{"Negative base 8 with prefix `o` to int8", "-o133", int8(-91), &root.i8},
				{"Positive base 8 with prefix `O` to int8", "O133", int8(91), &root.i8},
				{"Negative base 8 with prefix `O` to int8", "-O133", int8(-91), &root.i8},

				// to `int16`
				{"Positive base 8 with prefix `0o` to int16", "0o1337", int16(735), &root.i16},
				{"Negative base 8 with prefix `0o` to int16", "-0o1337", int16(-735), &root.i16},
				{"Positive base 8 with prefix `0O` to int16", "0O1337", int16(735), &root.i16},
				{"Negative base 8 with prefix `0O` to int16", "-0O1337", int16(-735), &root.i16},
				{"Positive base 8 with prefix `o` to int16", "o1337", int16(735), &root.i16},
				{"Negative base 8 with prefix `o` to int16", "-o1337", int16(-735), &root.i16},
				{"Positive base 8 with prefix `O` to int16", "O1337", int16(735), &root.i16},
				{"Negative base 8 with prefix `O` to int16", "-O1337", int16(-735), &root.i16},

				// to `int32`
				{"Positive base 8 with prefix `0o` to int32", "0o1337", int32(735), &root.i32},
				{"Negative base 8 with prefix `0o` to int32", "-0o1337", int32(-735), &root.i32},
				{"Positive base 8 with prefix `0O` to int32", "0O1337", int32(735), &root.i32},
				{"Negative base 8 with prefix `0O` to int32", "-0O1337", int32(-735), &root.i32},
				{"Positive base 8 with prefix `o` to int32", "o1337", int32(735), &root.i32},
				{"Negative base 8 with prefix `o` to int32", "-o1337", int32(-735), &root.i32},
				{"Positive base 8 with prefix `O` to int32", "O1337", int32(735), &root.i32},
				{"Negative base 8 with prefix `O` to int32", "-O1337", int32(-735), &root.i32},

				// to `int64`
				{"Positive base 8 with prefix `0o` to int64", "0o1337", int64(735), &root.i64},
				{"Negative base 8 with prefix `0o` to int64", "-0o1337", int64(-735), &root.i64},
				{"Positive base 8 with prefix `0O` to int64", "0O1337", int64(735), &root.i64},
				{"Negative base 8 with prefix `0O` to int64", "-0O1337", int64(-735), &root.i64},
				{"Positive base 8 with prefix `o` to int64", "o1337", int64(735), &root.i64},
				{"Negative base 8 with prefix `o` to int64", "-o1337", int64(-735), &root.i64},
				{"Positive base 8 with prefix `O` to int64", "O1337", int64(735), &root.i64},
				{"Negative base 8 with prefix `O` to int64", "-O1337", int64(-735), &root.i64},
			}

			for i := range validIntParseTest {
				tmp := &validIntParseTest[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Int(), C.ShouldEqual, tmp.Output)
				})

				// Reset test container
				zero()
			}
		})

		C.Convey("Valid uint values", func() {

			root := struct {
				u   uint
				u8  uint8
				u16 uint16
				u32 uint32
				u64 uint64
			}{}

			zero := func() {
				root.u = 0
				root.u8 = 0
				root.u16 = 0
				root.u32 = 0
				root.u64 = 0
			}

			validIntParseTest := []unmarshalerValueTest{
				//
				// Base 10
				//

				// to `uint`
				{"base 10 uint", "123", uint(123), &root.u},
				// to `uint8`
				{"base 10 uint8", "123", uint8(123), &root.u8},
				// to `uint16`
				{"base 10 uint16", "123", uint16(123), &root.u16},
				// to `uint32`
				{"base 10 uint32", "123", uint32(123), &root.u32},
				// to `uint64`
				{"base 10 uint64", "123", uint64(123), &root.u64},

				//
				// Base 16
				//

				// to `uint`
				{"base 16 with prefix `0x` to untyped uint", "0x123", uint(291), &root.u},
				{"base 16 with prefix `0X` to untyped uint", "0X123", uint(291), &root.u},
				{"base 16 with prefix `x` to untyped uint", "x123", uint(291), &root.u},
				{"base 16 with prefix `X` to untyped uint", "X123", uint(291), &root.u},

				// to `uint8`
				{"base 16 with prefix `0x` to uint8", "0x7E", uint8(126), &root.u8},
				{"base 16 with prefix `0X` to uint8", "0X7E", uint8(126), &root.u8},
				{"base 16 with prefix `x` to uint8", "x7E", uint8(126), &root.u8},
				{"base 16 with prefix `X` to uint8", "X7E", uint8(126), &root.u8},

				// to `uint16`
				{"base 16 with prefix `0x` to uint16", "0x7E", uint16(126), &root.u16},
				{"base 16 with prefix `0X` to uint16", "0X7E", uint16(126), &root.u16},
				{"base 16 with prefix `x` to uint16", "x7E", uint16(126), &root.u16},
				{"base 16 with prefix `X` to uint16", "X7E", uint16(126), &root.u16},

				// to `uint32`
				{"base 16 with prefix `0x` to uint32", "0x7E", uint32(126), &root.u32},
				{"base 16 with prefix `0X` to uint32", "0X7E", uint32(126), &root.u32},
				{"base 16 with prefix `x` to uint32", "x7E", uint32(126), &root.u32},
				{"base 16 with prefix `X` to uint32", "X7E", uint32(126), &root.u32},

				// to `uint64`
				{"base 16 with prefix `0x` to uint64", "0x7E", uint64(126), &root.u64},
				{"base 16 with prefix `0X` to uint64", "0X7E", uint64(126), &root.u64},
				{"base 16 with prefix `x` to uint64", "x7E", uint64(126), &root.u64},
				{"base 16 with prefix `X` to uint64", "X7E", uint64(126), &root.u64},

				//
				// Base 8
				//

				// to `uint`
				{"base 8 with prefix `0o` to untyped uint", "0o123", uint(83), &root.u},
				{"base 8 with prefix `0O` to untyped uint", "0O123", uint(83), &root.u},
				{"base 8 with prefix `o` to untyped uint", "o123", uint(83), &root.u},
				{"base 8 with prefix `O` to untyped uint", "O123", uint(83), &root.u},

				// to `uint8`
				{"base 8 with prefix `0o` to uint8", "0o133", uint8(91), &root.u8},
				{"base 8 with prefix `0O` to uint8", "0O133", uint8(91), &root.u8},
				{"base 8 with prefix `o` to uint8", "o133", uint8(91), &root.u8},
				{"base 8 with prefix `O` to uint8", "O133", uint8(91), &root.u8},

				// to `uint16`
				{"base 8 with prefix `0o` to uint16", "0o1337", uint16(735), &root.u16},
				{"base 8 with prefix `0O` to uint16", "0O1337", uint16(735), &root.u16},
				{"base 8 with prefix `o` to uint16", "o1337", uint16(735), &root.u16},
				{"base 8 with prefix `O` to uint16", "O1337", uint16(735), &root.u16},

				// to `uint32`
				{"base 8 with prefix `0o` to uint32", "0o1337", uint32(735), &root.u32},
				{"base 8 with prefix `0O` to uint32", "0O1337", uint32(735), &root.u32},
				{"base 8 with prefix `o` to uint32", "o1337", uint32(735), &root.u32},
				{"base 8 with prefix `O` to uint32", "O1337", uint32(735), &root.u32},

				// to `uint64`
				{"base 8 with prefix `0o` to uint64", "0o1337", uint64(735), &root.u64},
				{"base 8 with prefix `0O` to uint64", "0O1337", uint64(735), &root.u64},
				{"base 8 with prefix `o` to uint64", "o1337", uint64(735), &root.u64},
				{"base 8 with prefix `O` to uint64", "O1337", uint64(735), &root.u64},
			}

			for i := range validIntParseTest {
				tmp := &validIntParseTest[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Uint(), C.ShouldEqual, tmp.Output)
				})

				// Reset test container
				zero()
			}
		})

		C.Convey("Valid float values", func() {
			root := struct {
				f32 float32
				f64 float64
			}{}

			zero := func() {
				root.f32 = 0
				root.f64 = 0
			}

			parseTests := []unmarshalerValueTest{
				{"float32", "123", float32(123), &root.f32},
				{"float64", "123", float64(123), &root.f64},
			}

			for i := range parseTests {
				tmp := &parseTests[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Float(), C.ShouldEqual, tmp.Output)
				})

				// Reset test container
				zero()
			}

		})

		C.Convey("Valid bool values", func() {

			root := struct {
				// Values expected to be false
				f bool

				// Values expected to be true
				t bool
			}{}

			// Set the root to the opposite of the expectations
			zero := func() {
				root.f = true
				root.t = false
			}

			parseTests := []unmarshalerValueTest{
				{"bool false `f`", "f", false, &root.f},
				{"bool false `false`", "false", false, &root.f},
				{"bool false `0`", "0", false, &root.f},
				{"bool false `n`", "n", false, &root.f},
				{"bool false `no`", "no", false, &root.f},

				{"bool true `t`", "t", true, &root.t},
				{"bool true `true`", "true", true, &root.t},
				{"bool true `1`", "1", true, &root.t},
				{"bool true `y`", "y", true, &root.t},
				{"bool true `yes`", "yes", true, &root.t},
			}

			for i := range parseTests {
				tmp := &parseTests[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Bool(), C.ShouldEqual, tmp.Output)
				})

				// Reset test container
				zero()
			}

		})

		C.Convey("Valid string values", func() {
			root := struct {
				s string
			}{}

			zero := func() {
				root.s = ""
			}

			parseTests := []unmarshalerValueTest{
				{"string", "123", "123", &root.s},
			}

			for i := range parseTests {
				tmp := &parseTests[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.String(), C.ShouldEqual, tmp.Output)
				})

				// Reset test container
				zero()
			}

		})

		C.Convey("Valid slice values", func() {
			root := struct {
				s   []string
				sp  []*string
				i   []int
				ip  []*int
				b   []byte
				sb  [][]byte
				sbp []*[]byte
			}{}

			zero := func() {
				root.s = make([]string, 0, 1)
				root.sp = make([]*string, 0, 1)
				root.i = make([]int, 0, 1)
				root.ip = make([]*int, 0, 1)
				root.b = make([]byte, 0, 3)
				root.sb = make([][]byte, 0, 1)
				root.sbp = make([]*[]byte, 0, 1)
			}

			sp := "123"
			ip := 123
			bp := []byte{'1', '2', '3'}

			parseTests := []unmarshalerValueTest{
				{"string slice", sp, []string{sp}, &root.s},
				{"string pointer slice", sp, []*string{&sp}, &root.sp},
				{"int slice", sp, []int{ip}, &root.i},
				{"int pointer slice", sp, []*int{&ip}, &root.ip},
				{"byte slice", sp, bp, &root.b},
				{"byte slice slice", sp, [][]byte{bp}, &root.sb},
				{"byte slice pointer slice", sp, []*[]byte{&bp}, &root.sbp},
			}

			for i := range parseTests {
				tmp := &parseTests[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Interface(), C.ShouldResemble, tmp.Output)
				})

				// Reset test container
				zero()
			}
		})

		C.Convey("Valid map values", func() {
			root := struct {
				ss   map[string]string
				ssp  map[string]*string
				si   map[string]int
				sip  map[string]*int
				sb   map[string]byte
				ssb  map[string][]byte
				ssbp map[string]*[]byte
				ii   map[int]int
			}{}

			zero := func() {
				root.ss = make(map[string]string)
				root.ssp = make(map[string]*string)
				root.si = make(map[string]int)
				root.sip = make(map[string]*int)
				root.sb = make(map[string]byte)
				root.ssb = make(map[string][]byte)
				root.ssbp = make(map[string]*[]byte)
			}

			sKey  := "key"
			sVal  := "value"
			ssRaw := "key=value"
			//sp := "123"
			//ip := 123
			//bp := []byte{'1', '2', '3'}

			parseTests := []unmarshalerValueTest{
				{"map[string]string", ssRaw, map[string]string{sKey: sVal}, &root.ss},
				{"map[string]*string", ssRaw, map[string]*string{sKey: &sVal}, &root.ssp},
				//{"int slice", sp, []int{ip}, &root.si},
				//{"int pointer slice", sp, []*int{&ip}, &root.sip},
				//{"byte slice", sp, bp, &root.sb},
				//{"byte slice slice", sp, [][]byte{bp}, &root.ssb},
				//{"byte slice pointer slice", sp, []*[]byte{&bp}, &root.ssbp},
			}

			for i := range parseTests {
				tmp := &parseTests[i]

				C.Convey(tmp.Name, func() {
					C.So(impl.UnmarshalDefault(tmp.Input, &tmp.Temp), C.ShouldBeNil)

					eVal := reflect.ValueOf(tmp.Temp).Elem()

					C.So(eVal.Kind(), C.ShouldEqual, reflect.ValueOf(tmp.Output).Kind())
					C.So(eVal.Interface(), C.ShouldResemble, tmp.Output)
				})

				// Reset test container
				zero()
			}
		})
	})
}
