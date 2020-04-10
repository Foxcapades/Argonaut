package argument_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestArgument_Name(t *T) {
	Convey("Argument.Name", t, func() {
		So(NewBuilder(impl.NewProvider()).Name("foo").MustBuild().Name(), ShouldEqual, "foo")
	})
}

func TestArgument_Hint(t *T) {
	Convey("Argument.Hint", t, func() {
		So(NewBuilder(impl.NewProvider()).TypeHint("int").MustBuild().Hint(), ShouldEqual, "int")
	})
}

func TestArgument_HasName(t *T) {
	Convey("Argument.HasName", t, func() {
		So(NewBuilder(impl.NewProvider()).Name("straps").MustBuild().HasName(), ShouldBeTrue)
		So(NewBuilder(impl.NewProvider()).MustBuild().HasName(), ShouldBeFalse)
	})
}

func TestArgument_HasHint(t *T) {
	Convey("Argument.HasHint", t, func() {
		So(NewBuilder(impl.NewProvider()).TypeHint("int").MustBuild().HasHint(), ShouldBeTrue)
		So(NewBuilder(impl.NewProvider()).MustBuild().HasHint(), ShouldBeFalse)
	})
}
