package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestCompatible(t *testing.T) {
	Convey("Compatible", t, func() {
		Convey("With compatible values", func() {
			var testVal1 string
			var testVal2 string

			a := reflect.ValueOf(testVal1)
			b := reflect.ValueOf(testVal2)

			So(Compatible(&a, &b), ShouldBeTrue)
		})
		Convey("With incompatible values", func() {
			var testVal1 string
			var testVal2 int

			a := reflect.ValueOf(testVal1)
			b := reflect.ValueOf(testVal2)

			So(Compatible(&a, &b), ShouldBeFalse)

		})
	})
}
