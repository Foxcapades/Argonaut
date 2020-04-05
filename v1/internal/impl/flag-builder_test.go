package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v1/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestFlagBuilder_Arg(t *T) {
	Convey("FlagBuilder.Arg", t, func() {
		a := NewArgBuilder().Hint("diced bagels")
		b := NewFlagBuilder().Arg(a).Short('a').MustBuild()
		So(b.Argument(), ShouldResemble, a.MustBuild())
	})
}

func TestFlagBuilder_Bind(t *T) {
	Convey("FlagBuilder.Bind", t, func() {
		p := "is this even a good game?"

		Convey("required", func() {
			b := NewFlagBuilder().Bind(&p, true).Short('a').MustBuild()
			So(b.Argument().Required(), ShouldBeTrue)
			So(b.Argument().(*Argument).Binding(), ShouldPointTo, &p)
		})

		Convey("not required", func() {
			b := NewFlagBuilder().Bind(&p, false).Short('a').MustBuild()
			So(b.Argument().Required(), ShouldBeFalse)
			So(b.Argument().(*Argument).Binding(), ShouldPointTo, &p)
		})
	})
}

func TestFlagBuilder_Short(t *T) {
	Convey("FlagBuilder.Short", t, func() {
		f := NewFlagBuilder().Short('z').MustBuild()
		So(f.Short(), ShouldEqual, 'z')
		So(f.HasShort(), ShouldBeTrue)
	})
}

func TestFlagBuilder_Long(t *T) {
	Convey("FlagBuilder.Long", t, func() {
		f := NewFlagBuilder().Long("smerty").MustBuild()
		So(f.Long(), ShouldEqual, "smerty")
		So(f.HasLong(), ShouldBeTrue)
	})
}

func TestFlagBuilder_Description(t *T) {
	Convey("FlagBuilder.Description", t, func() {
		f := NewFlagBuilder().Short('a').Description("bananas are superior to mangos").MustBuild()
		So(f.Description(), ShouldEqual, "bananas are superior to mangos")
		So(f.HasDescription(), ShouldBeTrue)
	})
}

func TestFlagBuilder_Default(t *T) {
	Convey("FlagBuilder.Default", t, func() {
		p := "i'm pretty sure this is a bad game"

		Convey("required", func() {
			b := NewFlagBuilder().Default(&p).Short('a').MustBuild()
			So(b.Argument().Required(), ShouldBeFalse)
			So(b.Argument().(*Argument).Default(), ShouldPointTo, &p)
		})
	})
}
