package trait_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestDescribed_Description(t *T) {
	Convey("Described.Description", t, func() {
		var test Described
		So(test.Description(), ShouldEqual, "")
		test.DescTxt = "test"
		So(test.Description(), ShouldEqual, "test")
	})
}

func TestDescribed_HasDescription(t *T) {
	Convey("Described.HasDescription", t, func() {
		var test Described
		So(test.HasDescription(), ShouldBeFalse)
		test.DescTxt = "foo"
		So(test.HasDescription(), ShouldBeTrue)
	})
}
