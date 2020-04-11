package argo_test

import (
	. "testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func TestHex_Unmarshal(t *T) {
	Convey("argo.Hex.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Hex

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t Hex

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestHex8_Unmarshal(t *T) {
	Convey("argo.Hex8.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Hex8

			So(t.Unmarshal("66"), ShouldBeNil)
			So(t, ShouldEqual, 102)
		})

		Convey("Invalid values", func() {
			var t Hex8

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestHex16_Unmarshal(t *T) {
	Convey("argo.Hex16.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Hex16

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t Hex16

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestHex32_Unmarshal(t *T) {
	Convey("argo.Hex32.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Hex32

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t Hex32

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestHex64_Unmarshal(t *T) {
	Convey("argo.Hex64.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Hex64

			So(t.Unmarshal("ff"), ShouldBeNil)
			So(t, ShouldEqual, 255)
		})

		Convey("Invalid values", func() {
			var t Hex64

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}
