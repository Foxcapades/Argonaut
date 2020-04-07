package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestArgument_Name(t *T) {
	Convey("Argument.Name", t, func() {
		So(NewArgBuilder().Name("foo").MustBuild().Name(), ShouldEqual, "foo")
	})
}

func TestArgument_Hint(t *T) {
	Convey("Argument.Hint", t, func() {
		So(NewArgBuilder().TypeHint("int").MustBuild().Hint(), ShouldEqual, "int")
	})
}

func TestArgument_HasName(t *T) {
	Convey("Argument.HasName", t, func() {
		So(NewArgBuilder().Name("straps").MustBuild().HasName(), ShouldBeTrue)
		So(NewArgBuilder().MustBuild().HasName(), ShouldBeFalse)
	})
}

func TestArgument_HasHint(t *T) {
	Convey("Argument.HasHint", t, func() {
		So(NewArgBuilder().TypeHint("int").MustBuild().HasHint(), ShouldBeTrue)
		So(NewArgBuilder().MustBuild().HasHint(), ShouldBeFalse)
	})
}
