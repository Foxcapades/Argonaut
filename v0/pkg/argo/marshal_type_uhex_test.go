package argo_test

import (
	. "testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func TestUHex_Unmarshal(t *T) {
	Convey("argo.UHex.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UHex

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UHex

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUHex8_Unmarshal(t *T) {
	Convey("argo.UHex8.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UHex8

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UHex8

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUHex16_Unmarshal(t *T) {
	Convey("argo.UHex16.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UHex16

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UHex16

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUHex32_Unmarshal(t *T) {
	Convey("argo.UHex32.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UHex32

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UHex32

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestUHex64_Unmarshal(t *T) {
	Convey("argo.UHex64.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t UHex64

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t UHex64

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}
