package argument_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	com "github.com/Foxcapades/Argonaut/v0/internal/impl/command"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
	. "testing"
)

func TestArgument_Name(t *T) {
	Convey("Argument.Name", t, func() {
		So((&Argument{Named: Named{NameValue: "foo"}}).Name(), ShouldEqual, "foo")
	})
}

func TestArgument_HasName(t *T) {
	Convey("Argument.HasName", t, func() {
		So((&Argument{Named: Named{NameValue: "straps"}}).HasName(), ShouldBeTrue)
		So((&Argument{}).HasName(), ShouldBeFalse)
	})
}

func TestArgument_Binding(t *T) {
	Convey("Argument.Binding", t, func() {
		test := "hi"
		So((&Argument{BindValue: &test}).Binding(), ShouldPointTo, &test)
	})
}

func TestArgument_BindingType(t *T) {
	Convey("Argument.BindingType", t, func() {
		test1 := "foo"

		So((&Argument{RootBinding: R.ValueOf(test1), IsBindingSet: true}).
			BindingType().Kind(), ShouldEqual, R.String)
		So((&Argument{}).BindingType(), ShouldBeNil)
	})
}

func TestArgument_HasBinding(t *T) {
	Convey("Argument.HasBinding", t, func() {
		So((&Argument{IsBindingSet: true}).HasBinding(), ShouldBeTrue)
		So((&Argument{}).HasBinding(), ShouldBeFalse)
	})
}

func TestArgument_Parent(t *T) {
	Convey("Argument.Parent", t, func() {
		test1 := new(com.Command)
		test2 := new(flag.Flag)

		So((&Argument{ParentElement: test1}).Parent(), ShouldPointTo, test1)
		So((&Argument{ParentElement: test2}).Parent(), ShouldPointTo, test2)
	})
}

func TestArgument_RawValue(t *T) {
	Convey("Argument.RawValue", t, func() {
		test1 := "foo"

		So((&Argument{}).RawValue(), ShouldEqual, "")
		So((&Argument{RawInput: test1}).RawValue(), ShouldEqual, test1)
	})
}

func TestArgument_IsFlagArg(t *T) {
	Convey("Argument.IsFlagArg", t, func() {
		test1 := A.Command(new(com.Command))
		test2 := A.Flag(new(flag.Flag))

		So((&Argument{ParentElement: test1}).IsFlagArg(), ShouldBeFalse)
		So((&Argument{ParentElement: test2}).IsFlagArg(), ShouldBeTrue)
	})
}

func TestArgument_IsPositionalArg(t *T) {
	Convey("Argument.IsPositionalArg", t, func() {
		test1 := A.Command(new(com.Command))
		test2 := A.Flag(new(flag.Flag))

		So((&Argument{ParentElement: test1}).IsPositionalArg(), ShouldBeTrue)
		So((&Argument{ParentElement: test2}).IsPositionalArg(), ShouldBeFalse)
	})
}

func TestArgument_String(t *T) {
	Convey("Argument.String", t, func() {
		So((&Argument{}).String(), ShouldEqual, "arg")
		So((&Argument{Named: Named{NameValue: "foo"}}).String(), ShouldEqual,
			"foo")
	})
}

func TestArgument_HasDefault(t *T) {
	Convey("Argument.HasDefault", t, func() {
		So((&Argument{IsDefaultSet: true}).HasDefault(), ShouldBeTrue)
		So((&Argument{}).HasDefault(), ShouldBeFalse)
	})
}

func TestArgument_DefaultType(t *T) {
	Convey("Argument.DefaultType", t, func() {
		bind := ""

		So((&Argument{RootDefault: R.ValueOf(bind), IsDefaultSet: true}).DefaultType().Kind(),
			ShouldEqual, R.String)
		So((&Argument{}).DefaultType(), ShouldBeNil)
	})
}
