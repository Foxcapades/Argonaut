package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v1/internal/impl"
	A "github.com/Foxcapades/Argonaut/v1/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
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

func TestArgumentBuilder_Build(t *T) {
	Convey("ArgumentBuilder.Build", t, func() {
		Convey("nil binding", func() {
			a, b := NewArgBuilder().Bind(nil).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.InvalidArgError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.InvalidArgBindingError)
		})

		Convey("unaddressable binding", func() {
			a, b := NewArgBuilder().Bind(3).Default(3).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.InvalidArgError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.InvalidArgBindingError)
			So(c.BindingType().Kind(), ShouldEqual, R.Int)
			So(c.DefaultValType().Kind(), ShouldEqual, R.Int)
		})

		Convey("type mismatch", func() {
			e := ""
			f := 3
			a, b := NewArgBuilder().Default(e).Bind(&f).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.InvalidArgError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.InvalidArgDefaultError)
			So(c.BindingType().Elem().Kind(), ShouldEqual, R.Int)
			So(c.DefaultValType().Kind(), ShouldEqual, R.String)
		})
	})
}

func TestArgumentBuilder_MustBuild(t *T) {
	prep := func(a A.ArgumentBuilder) (func(), func() interface{}) {
		var err interface{}
		return func() {
			defer func() { err = recover(); panic(err) }()
			a.MustBuild()
		}, func() interface{} { return err }
	}

	Convey("ArgumentBuilder.MustBuild", t, func() {
		Convey("nil binding", func() {
			fn, b := prep(NewArgBuilder().Bind(nil))
			So(fn, ShouldPanic)
			c, d := b().(A.InvalidArgError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.InvalidArgBindingError)
		})

		Convey("unaddressable binding", func() {
			fn, b := prep(NewArgBuilder().Bind(3).Default(3))
			So(fn, ShouldPanic)
			c, d := b().(A.InvalidArgError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.InvalidArgBindingError)
			So(c.BindingType().Kind(), ShouldEqual, R.Int)
			So(c.DefaultValType().Kind(), ShouldEqual, R.Int)
		})

		Convey("type mismatch", func() {
			e := ""
			f := 3
			fn, b := prep(NewArgBuilder().Default(e).Bind(&f))
			So(fn, ShouldPanic)
			c, d := b().(A.InvalidArgError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.InvalidArgDefaultError)
			So(c.BindingType().Elem().Kind(), ShouldEqual, R.Int)
			So(c.DefaultValType().Kind(), ShouldEqual, R.String)
		})
	})
}
