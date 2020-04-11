package argument_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
	. "testing"
)

func TestArgumentBuilder_Bind(t *T) {
	Convey("ArgumentBuilder.Bind", t, func() {
		t := "flarf"
		a := &argument.Builder{}

		a.Bind(&t)

		So(a.BindValue, ShouldPointTo, &t)
		So(a.IsBindingSet, ShouldBeTrue)
	})
}

func TestArgumentBuilder_Default(t *T) {
	Convey("ArgumentBuilder.Default", t, func() {
		Convey("Using a direct value", func() {
			t := "flumps"
			a := argument.Builder{}

			a.Default(t)

			So(a.DefaultValue, ShouldResemble, t)
			So(a.IsDefaultSet, ShouldBeTrue)
		})

		Convey("Using an indirect value", func() {
			t := "gampus"
			a := argument.Builder{}

			a.Default(&t)

			So(a.DefaultValue, ShouldPointTo, &t)
			So(a.IsDefaultSet, ShouldBeTrue)
		})
	})
}

func TestArgumentBuilder_Description(t *T) {
	Convey("ArgumentBuilder.Description", t, func() {
		v := "interior crocodile, alligator"
		a := argument.Builder{}

		a.Description(v)

		So(a.DescriptionValue.DescriptionText, ShouldResemble, v)
	})
}

func TestArgumentBuilder_Require(t *T) {
	Convey("ArgumentBuilder.Require", t, func() {
		a := &argument.Builder{}
		a.Require()
		So(a.IsArgRequired, ShouldBeTrue)
	})
}

func TestArgumentBuilder_Required(t *T) {
	Convey("ArgumentBuilder.Required", t, func() {
		a := argument.Builder{}
		a.Required(true)
		So(a.IsArgRequired, ShouldBeTrue)

		b := argument.Builder{}
		b.Required(false)
		So(b.IsArgRequired, ShouldBeFalse)
	})
}

func TestArgumentBuilder_Build(t *T) {
	Convey("ArgumentBuilder.Build", t, func() {
		Convey("nil binding", func() {
			a, b := argument.NewBuilder(impl.NewProvider()).Bind(nil).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})

		Convey("unaddressable binding", func() {
			a, b := argument.NewBuilder(impl.NewProvider()).Bind(3).Default(3).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
			So(R.TypeOf(c.Builder().GetBinding()).Kind(), ShouldEqual, R.Int)
			So(R.TypeOf(c.Builder().GetDefault()).Kind(), ShouldEqual, R.Int)
		})

		Convey("type mismatch", func() {
			e := ""
			f := 3
			a, b := argument.NewBuilder(impl.NewProvider()).Default(e).Bind(&f).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidDefaultVal)
			So(R.TypeOf(c.Builder().GetBinding()).Elem().Kind(), ShouldEqual, R.Int)
			So(R.TypeOf(c.Builder().GetDefault()).Kind(), ShouldEqual, R.String)
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
			fn, b := prep(argument.NewBuilder(impl.NewProvider()).Bind(nil))
			So(fn, ShouldPanic)
			c, d := b().(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})

		Convey("unaddressable binding", func() {
			fn, b := prep(argument.NewBuilder(impl.NewProvider()).Bind(3).Default(3))
			So(fn, ShouldPanic)
			c, d := b().(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
			So(R.TypeOf(c.Builder().GetBinding()).Kind(), ShouldEqual, R.Int)
			So(R.TypeOf(c.Builder().GetDefault()).Kind(), ShouldEqual, R.Int)
		})

		Convey("type mismatch", func() {
			e := ""
			f := 3
			fn, b := prep(argument.NewBuilder(impl.NewProvider()).Default(e).Bind(&f))
			So(fn, ShouldPanic)
			c, d := b().(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidDefaultVal)
			So(R.TypeOf(c.Builder().GetBinding()).Elem().Kind(), ShouldEqual, R.Int)
			So(R.TypeOf(c.Builder().GetDefault()).Kind(), ShouldEqual, R.String)
		})
	})
}
