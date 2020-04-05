package util_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestIsUnmarshalable(t *T) {
	Convey("IsUnmarshalable", t, func() {
		Convey("nil input", func() {
			So(IsUnmarshalable(nil), ShouldBeFalse)
		})
		Convey("unaddressable input", func() {
			So(IsUnmarshalable(3), ShouldBeFalse)
			So(IsUnmarshalable("ford"), ShouldBeFalse)
		})
	})
}
