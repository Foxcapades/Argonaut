package util_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	. "testing"
)

func TestToUnmarshalable(t *T) {
	Convey("IsUnmarshalable", t, func() {
		Convey("unaddressable input", func() {
			_, b := util.ToUnmarshalable("", reflect.ValueOf(3), false)
			So(b, ShouldNotBeNil)

			_, c := util.ToUnmarshalable("", reflect.ValueOf("foo"), false)
			So(c, ShouldNotBeNil)
		})

		Convey("slices", toUnmarshalableSliceTests)
		Convey("maps", toUnmarshalableMapTests)
	})
}

func toUnmarshalableSliceTests() {
	Convey("[]string", func() {
		var foo []string

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[][]byte", func() {
		var foo [][]byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]byte", func() {
		var foo []byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]uint", func() {
		var foo []uint

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]uint8", func() {
		var foo []uint8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]uint16", func() {
		var foo []uint16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]uint32", func() {
		var foo []uint32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]uint64", func() {
		var foo []uint64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]int", func() {
		var foo []int

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]int8", func() {
		var foo []int8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]int16", func() {
		var foo []int16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]int32", func() {
		var foo []int32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]int64", func() {
		var foo []int64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]float32", func() {
		var foo []float32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]float64", func() {
		var foo []float64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]bool", func() {
		var foo []bool

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*string", func() {
		var foo []*string

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*[]byte", func() {
		var foo []*[]byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*byte", func() {
		var foo []*byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*uint", func() {
		var foo []*uint

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*uint8", func() {
		var foo []*uint8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*uint16", func() {
		var foo []*uint16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*uint32", func() {
		var foo []*uint32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*uint64", func() {
		var foo []*uint64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*int", func() {
		var foo []*int

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*int8", func() {
		var foo []*int8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*int16", func() {
		var foo []*int16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*int32", func() {
		var foo []*int32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*int64", func() {
		var foo []*int64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*float32", func() {
		var foo []*float32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*float64", func() {
		var foo []*float64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*bool", func() {
		var foo []*bool

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]argo.Unmarshaler", func() {
		var foo []argo.Unmarshaler

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*argo.Unmarshaler", func() {
		var foo []*argo.Unmarshaler

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]interface{}", func() {
		var foo []interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*interface{}", func() {
		var foo []*interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]map[string]string", func() {
		var foo []map[string]string

		_, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldNotBeNil)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("[]*map[string]string", func() {
		var foo []*map[string]string

		_, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldNotBeNil)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
}
func toUnmarshalableMapTests() {
	Convey("map[string]string", func() {
		var foo map[string]string

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string][]byte", func() {
		var foo map[string][]byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]byte", func() {
		var foo map[string]byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]uint", func() {
		var foo map[string]uint

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]uint8", func() {
		var foo map[string]uint8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]uint16", func() {
		var foo map[string]uint16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]uint32", func() {
		var foo map[string]uint32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]uint64", func() {
		var foo map[string]uint64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]int", func() {
		var foo map[string]int

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]int8", func() {
		var foo map[string]int8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]int16", func() {
		var foo map[string]int16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]int32", func() {
		var foo map[string]int32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]int64", func() {
		var foo map[string]int64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]float32", func() {
		var foo map[string]float32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]float64", func() {
		var foo map[string]float64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]bool", func() {
		var foo map[string]bool

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*string", func() {
		var foo map[string]*string

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*[]byte", func() {
		var foo map[string]*[]byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*byte", func() {
		var foo map[string]*byte

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*uint", func() {
		var foo map[string]*uint

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*uint8", func() {
		var foo map[string]*uint8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*uint16", func() {
		var foo map[string]*uint16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*uint32", func() {
		var foo map[string]*uint32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*uint64", func() {
		var foo map[string]*uint64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*int", func() {
		var foo map[string]*int

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*int8", func() {
		var foo map[string]*int8

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*int16", func() {
		var foo map[string]*int16

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*int32", func() {
		var foo map[string]*int32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*int64", func() {
		var foo map[string]*int64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*float32", func() {
		var foo map[string]*float32

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*float64", func() {
		var foo map[string]*float64

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*bool", func() {
		var foo map[string]*bool

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]argo.Unmarshaler", func() {
		var foo map[string]argo.Unmarshaler

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*argo.Unmarshaler", func() {
		var foo map[string]*argo.Unmarshaler

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]interface{}", func() {
		var foo map[string]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*interface{}", func() {
		var foo map[string]*interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]map[string]string", func() {
		var foo map[string]map[string]string

		_, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldNotBeNil)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[string]*map[string]string", func() {
		var foo map[string]*map[string]string

		_, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldNotBeNil)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[int]interface{}", func() {
		var foo map[int]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[int8]interface{}", func() {
		var foo map[int8]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[int16]interface{}", func() {
		var foo map[int16]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[int32]interface{}", func() {
		var foo map[int32]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[int64]interface{}", func() {
		var foo map[int64]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[uint]interface{}", func() {
		var foo map[uint]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[uint8]interface{}", func() {
		var foo map[uint8]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[uint16]interface{}", func() {
		var foo map[uint16]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[uint32]interface{}", func() {
		var foo map[uint32]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[uint64]interface{}", func() {
		var foo map[uint64]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[float32]interface{}", func() {
		var foo map[float32]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[float64]interface{}", func() {
		var foo map[float64]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[bool]interface{}", func() {
		var foo map[bool]interface{}

		a, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldBeNil)
		So(a.Interface(), ShouldResemble, foo)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
	Convey("map[*bool]interface{}", func() {
		var foo map[*bool]interface{}

		_, b := util.ToUnmarshalable("", reflect.ValueOf(foo), true)
		So(b, ShouldNotBeNil)

		_, d := util.ToUnmarshalable("", reflect.ValueOf(foo), false)
		So(d, ShouldNotBeNil)
	})
}
