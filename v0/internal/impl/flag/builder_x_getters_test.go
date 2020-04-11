package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestBuilder_GetDescription(t *T) {
	Convey("flag.Builder.GetDescription", t, func() {
		So((&Builder{DescriptionText: Described{DescriptionText: "foo"}}).
			GetDescription(), ShouldEqual, "foo")
	})
}

func TestBuilder_HasDescription(t *T) {
	Convey("flag.Builder.HasDescription", t, func() {
		So((&Builder{DescriptionText: Described{DescriptionText: "foo"}}).
			HasDescription(), ShouldBeTrue)
	})
}

func TestBuilder_GetShort(t *T) {
	Convey("flag.Builder.GetShort", t, func() {
		So((&Builder{ShortFlag: 'c'}).GetShort(), ShouldEqual, 'c')
	})
}

func TestBuilder_HasShort(t *T) {
	Convey("flag.Builder.HasShort", t, func() {
		So((&Builder{IsShortSet: true}).HasShort(), ShouldBeTrue)
		So((&Builder{}).HasShort(), ShouldBeFalse)
	})
}

func TestBuilder_GetLong(t *T) {
	Convey("flag.Builder.GetLong", t, func() {
		So((&Builder{LongFlag: "fudge"}).GetLong(), ShouldEqual, "fudge")
	})
}

func TestBuilder_HasLong(t *T) {
	Convey("flag.Builder.HasLong", t, func() {
		So((&Builder{IsLongSet: true}).HasLong(), ShouldBeTrue)
		So((&Builder{}).HasLong(), ShouldBeFalse)
	})
}

func TestBuilder_GetArg(t *T) {
	Convey("flag.Builder.GetArg", t, func() {
		arg := argument.NewBuilder(impl.NewProvider())
		So((&Builder{ArgBuilder: arg}).GetArg(), ShouldEqual, arg)
	})
}

func TestBuilder_HasArg(t *T) {
	Convey("flag.Builder.HasArg", t, func() {
		arg := argument.NewBuilder(impl.NewProvider())
		So((&Builder{ArgBuilder: arg}).HasArg(), ShouldBeTrue)
	})
}
