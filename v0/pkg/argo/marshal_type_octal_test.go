package argo_test

import (
	. "testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func TestOctal_Unmarshal(t *T) {
	Convey("argo.Octal.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Octal

			So(t.Unmarshal("177"), ShouldBeNil)
			So(t, ShouldEqual, 127)
		})

		Convey("Invalid values", func() {
			var t Octal

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestOctal8_Unmarshal(t *T) {
	Convey("argo.Octal8.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Octal8

			So(t.Unmarshal("177"), ShouldBeNil)
			So(t, ShouldEqual, 127)
		})

		Convey("Invalid values", func() {
			var t Octal8

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestOctal16_Unmarshal(t *T) {
	Convey("argo.Octal16.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Octal16

			So(t.Unmarshal("177"), ShouldBeNil)
			So(t, ShouldEqual, 127)
		})

		Convey("Invalid values", func() {
			var t Octal16

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestOctal32_Unmarshal(t *T) {
	Convey("argo.Octal32.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Octal32

			So(t.Unmarshal("177"), ShouldBeNil)
			So(t, ShouldEqual, 127)
		})

		Convey("Invalid values", func() {
			var t Octal32

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}

func TestOctal64_Unmarshal(t *T) {
	Convey("argo.Octal64.Unmarshal", t, func() {
		Convey("Valid values", func() {
			var t Octal64

			So(t.Unmarshal("177"), ShouldBeNil)
			So(t, ShouldEqual, 127)
		})

		Convey("Invalid values", func() {
			var t Octal64

			So(t.Unmarshal("apple"), ShouldNotBeNil)
		})
	})
}
