package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v1/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestArgumentBuilder_Bind(t *T) {
	Convey("ArgumentBuilder.Bind", t, func() {
		t := "flarf"
		a := NewArgBuilder().Bind(&t).MustBuild().(*Argument)
		So(a.Binding(), ShouldPointTo, &t)
		So(a.HasBinding(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Default(t *T) {
	Convey("ArgumentBuilder.Default", t, func() {
		Convey("Using a direct value", func() {
			t := "flumps"
			a := NewArgBuilder().Default(t).MustBuild().(*Argument)
			So(a.Default(), ShouldResemble, t)
			So(a.HasDefault(), ShouldBeTrue)
		})
		Convey("Using an indirect value", func() {
			t := "gampus"
			a := NewArgBuilder().Default(&t).MustBuild().(*Argument)
			So(a.Default(), ShouldPointTo, &t)
			So(a.HasDefault(), ShouldBeTrue)
		})
	})
}

func TestArgumentBuilder_Hint(t *T) {
	Convey("ArgumentBuilder.Hint", t, func() {
		v := "boots with the fur"
		a := NewArgBuilder().Hint(v).MustBuild()
		So(a.Hint(), ShouldResemble, v)
		So(a.HasHint(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Description(t *T) {
	Convey("ArgumentBuilder.Description", t, func() {
		v := "interior crocodile, alligator"
		a := NewArgBuilder().Description(v).MustBuild()
		So(a.Description(), ShouldResemble, v)
		So(a.HasDescription(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Require(t *T) {
	Convey("ArgumentBuilder.Require", t, func() {
		a := NewArgBuilder().Require().MustBuild()
		So(a.Required(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Required(t *T) {
	Convey("ArgumentBuilder.Required", t, func() {
		a := NewArgBuilder().Required(true).MustBuild()
		So(a.Required(), ShouldBeTrue)
		b := NewArgBuilder().Required(false).MustBuild()
		So(b.Required(), ShouldBeFalse)
	})
}
