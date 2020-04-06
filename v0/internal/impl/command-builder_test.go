package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
	. "testing"
)

func TestCommandBuilder_Arg(t *T) {
	Convey("CommandBuilder.Arg", t, func() {
		Convey("Non-nil arg", func() {
			arg := NewArgBuilder()
			com := NewCommandBuilder()
			So(com.GetArgs(), ShouldBeEmpty)
			So(R.ValueOf(com.Arg(arg)).Pointer(), ShouldEqual, R.ValueOf(com).Pointer())
			So(com.GetArgs(), ShouldNotBeEmpty)
			So(R.ValueOf(com.GetArgs()[0]).Pointer(), ShouldEqual, R.ValueOf(arg).Pointer())
		})

		Convey("Nil arg", func() {
			com := NewCommandBuilder().Arg(nil)
			So(com.GetArgs(), ShouldBeEmpty)
			So(len(com.Warnings()), ShouldEqual, 1)
		})
	})
}

func TestCommandBuilder_Description(t *T) {
	Convey("CommandBuilder.Description", t, func() {
		com := NewCommandBuilder()
		So(com.HasDescription(), ShouldBeFalse)
		So(R.ValueOf(com.Description("foo")).Pointer(), ShouldEqual, R.ValueOf(com).Pointer())
		So(com.HasDescription(), ShouldBeTrue)
		So(com.GetDescription(), ShouldEqual, "foo")
	})
}