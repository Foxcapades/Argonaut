package argo_test

import (
	. "testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func TestUOctal_Unmarshal(t *T) {
	Convey("argo.UOctal.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UOctal

			So(t.Unmarshal("377"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UOctal

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUOctal8_Unmarshal(t *T) {
	Convey("argo.UOctal8.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UOctal8

			So(t.Unmarshal("377"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UOctal8

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUOctal16_Unmarshal(t *T) {
	Convey("argo.UOctal16.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UOctal16

			So(t.Unmarshal("377"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UOctal16

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUOctal32_Unmarshal(t *T) {
	Convey("argo.UOctal32.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UOctal32

			So(t.Unmarshal("377"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UOctal32

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUOctal64_Unmarshal(t *T) {
	Convey("argo.UOctal64.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UOctal64

			So(t.Unmarshal("377"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UOctal64

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}
