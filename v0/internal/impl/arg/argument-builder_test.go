package arg_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/arg"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
	. "testing"
)

func TestArgumentBuilder_Bind(t *T) {
	Convey("ArgumentBuilder.Bind", t, func() {
		t := "flarf"
		a := arg.NewBuilder(impl.NewProvider()).Bind(&t).MustBuild().(*arg.Argument)
		So(a.Binding(), ShouldPointTo, &t)
		So(a.HasBinding(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Default(t *T) {
	Convey("ArgumentBuilder.Default", t, func() {
		Convey("Using a direct value", func() {
			t := "flumps"
			a := arg.NewBuilder(impl.NewProvider()).Default(t).MustBuild().(*arg.Argument)
			So(a.Default(), ShouldResemble, t)
			So(a.HasDefault(), ShouldBeTrue)
		})
		Convey("Using an indirect value", func() {
			t := "gampus"
			a := arg.NewBuilder(impl.NewProvider()).Default(&t).MustBuild().(*arg.Argument)
			So(a.Default(), ShouldPointTo, &t)
			So(a.HasDefault(), ShouldBeTrue)
		})
	})
}

func TestArgumentBuilder_Hint(t *T) {
	Convey("ArgumentBuilder.TypeHint", t, func() {
		v := "boots with the fur"
		a := arg.NewBuilder(impl.NewProvider()).TypeHint(v).MustBuild()
		So(a.Hint(), ShouldResemble, v)
		So(a.HasHint(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Description(t *T) {
	Convey("ArgumentBuilder.Description", t, func() {
		v := "interior crocodile, alligator"
		a := arg.NewBuilder(impl.NewProvider()).Description(v).MustBuild()
		So(a.Description(), ShouldResemble, v)
		So(a.HasDescription(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Require(t *T) {
	Convey("ArgumentBuilder.Require", t, func() {
		a := arg.NewBuilder(impl.NewProvider()).Require().MustBuild()
		So(a.Required(), ShouldBeTrue)
	})
}

func TestArgumentBuilder_Required(t *T) {
	Convey("ArgumentBuilder.Required", t, func() {
		a := arg.NewBuilder(impl.NewProvider()).Required(true).MustBuild()
		So(a.Required(), ShouldBeTrue)
		b := arg.NewBuilder(impl.NewProvider()).Required(false).MustBuild()
		So(b.Required(), ShouldBeFalse)
	})
}

func TestArgumentBuilder_Build(t *T) {
	Convey("ArgumentBuilder.Build", t, func() {
		Convey("nil binding", func() {
			a, b := arg.NewBuilder(impl.NewProvider()).Bind(nil).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			c, d := b.(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})

		Convey("unaddressable binding", func() {
			a, b := arg.NewBuilder(impl.NewProvider()).Bind(3).Default(3).Build()
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
			a, b := arg.NewBuilder(impl.NewProvider()).Default(e).Bind(&f).Build()
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
			fn, b := prep(arg.NewBuilder(impl.NewProvider()).Bind(nil))
			So(fn, ShouldPanic)
			c, d := b().(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})

		Convey("unaddressable binding", func() {
			fn, b := prep(arg.NewBuilder(impl.NewProvider()).Bind(3).Default(3))
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
			fn, b := prep(arg.NewBuilder(impl.NewProvider()).Default(e).Bind(&f))
			So(fn, ShouldPanic)
			c, d := b().(A.ArgumentError)
			So(d, ShouldBeTrue)
			So(c.Type(), ShouldEqual, A.ArgErrInvalidDefaultVal)
			So(R.TypeOf(c.Builder().GetBinding()).Elem().Kind(), ShouldEqual, R.Int)
			So(R.TypeOf(c.Builder().GetDefault()).Kind(), ShouldEqual, R.String)
		})
	})
}
