package argument_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	com "github.com/Foxcapades/Argonaut/v0/internal/impl/command"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
	. "testing"
)

func TestArgument_Name(t *T) {
	Convey("Argument.Name", t, func() {
		So(NewBuilder(NewProvider()).Name("foo").MustBuild().Name(), ShouldEqual, "foo")
	})
}

func TestArgument_HasName(t *T) {
	Convey("Argument.HasName", t, func() {
		So(NewBuilder(NewProvider()).Name("straps").MustBuild().HasName(), ShouldBeTrue)
		So(NewBuilder(NewProvider()).MustBuild().HasName(), ShouldBeFalse)
	})
}

func TestArgument_Binding(t *T) {
	Convey("Argument.Binding", t, func() {
		test := "hi"

		So(NewBuilder(NewProvider()).Bind(&test).MustBuild().Binding(),
			ShouldPointTo, &test)
	})
}

func TestArgument_BindingType(t *T) {
	Convey("Argument.BindingType", t, func() {
		test1 := "foo"
		test2 := &test1
		test3 := &test2

		So(NewBuilder(NewProvider()).Bind(&test1).MustBuild().BindingType().Kind(),
			ShouldEqual, R.String)
		So(NewBuilder(NewProvider()).Bind(&test2).MustBuild().BindingType().Kind(),
			ShouldEqual, R.String)
		So(NewBuilder(NewProvider()).Bind(&test3).MustBuild().BindingType().Kind(),
			ShouldEqual, R.String)
		So(NewBuilder(NewProvider()).MustBuild().BindingType(), ShouldBeNil)
	})
}

func TestArgument_HasBinding(t *T) {
	Convey("Argument.HasBinding", t, func() {
		test1 := "foo"

		So(NewBuilder(NewProvider()).Bind(&test1).MustBuild().HasBinding(),
			ShouldBeTrue)
		So(NewBuilder(NewProvider()).MustBuild().HasBinding(), ShouldBeFalse)
	})
}

func TestArgument_Parent(t *T) {
	Convey("Argument.Parent", t, func() {
		test1 := new(com.Command)
		test2 := new(flag.Flag)

		So(NewBuilder(NewProvider()).Parent(test1).MustBuild().Parent(),
			ShouldPointTo, test1)
		So(NewBuilder(NewProvider()).Parent(test2).MustBuild().Parent(),
			ShouldPointTo, test2)
	})
}

func TestArgument_RawValue(t *T) {
	Convey("Argument.RawValue", t, func() {
		test1 := "foo"

		So(NewBuilder(NewProvider()).MustBuild().RawValue(), ShouldEqual, "")

		arg := NewBuilder(NewProvider()).MustBuild()
		arg.SetRawValue(test1)
		So(arg.RawValue(), ShouldEqual, test1)
	})
}

func TestArgument_IsFlagArg(t *T) {
	Convey("Argument.IsFlagArg", t, func() {
		test1 := A.Command(new(com.Command))
		test2 := A.Flag(new(flag.Flag))

		So(NewBuilder(NewProvider()).Parent(test1).MustBuild().IsFlagArg(),
			ShouldBeFalse)
		So(NewBuilder(NewProvider()).Parent(test2).MustBuild().IsFlagArg(),
			ShouldBeTrue)
	})
}

func TestArgument_IsPositionalArg(t *T) {
	Convey("Argument.IsPositionalArg", t, func() {
		test1 := A.Command(new(com.Command))
		test2 := A.Flag(new(flag.Flag))

		So(NewBuilder(NewProvider()).Parent(test1).MustBuild().IsPositionalArg(),
			ShouldBeTrue)
		So(NewBuilder(NewProvider()).Parent(test2).MustBuild().IsPositionalArg(),
			ShouldBeFalse)
	})
}

func TestArgument_String(t *T) {
	Convey("Argument.String", t, func() {
		So(NewBuilder(NewProvider()).MustBuild().String(), ShouldEqual, "arg")
		So(NewBuilder(NewProvider()).Name("foo").MustBuild().String(), ShouldEqual,
			"foo")
	})
}

func TestArgument_HasDefault(t *T) {
	Convey("Argument.HasDefault", t, func() {
		bind := ""

		So(NewBuilder(NewProvider()).Bind(&bind).Default("foo").MustBuild().HasDefault(),
			ShouldBeTrue)
		So(NewBuilder(NewProvider()).Bind(&bind).MustBuild().HasDefault(),
			ShouldBeFalse)
	})
}

func TestArgument_DefaultType(t *T) {
	Convey("Argument.DefaultType", t, func() {
		bind := ""

		So(NewBuilder(
			NewProvider()).
			Bind(&bind).
			Default("foo").
			MustBuild().
			DefaultType().
			Kind(),
			ShouldEqual, R.String)
		So(NewBuilder(NewProvider()).Bind(&bind).MustBuild().DefaultType(),
			ShouldBeNil)
	})
}
