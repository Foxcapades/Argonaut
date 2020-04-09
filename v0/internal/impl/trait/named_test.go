package trait_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestNamed_Name(t *T) {
	Convey("Named.Name", t, func() {
		var test Named
		So(test.Name(), ShouldEqual, "")
		test.NameValue = "test"
		So(test.Name(), ShouldEqual, "test")
	})
}

func TestNamed_HasName(t *T) {
	Convey("Named.HasName", t, func() {
		var test Named
		So(test.HasName(), ShouldBeFalse)
		test.NameValue = "foo"
		So(test.HasName(), ShouldBeTrue)
	})
}
