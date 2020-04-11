package flag_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
)

func TestGBuilder_Description(t *testing.T) {
	convey.Convey("flag.GBuilder.Description()", t, func() {
		t := &flag.GBuilder{}

		convey.So(t.Description("foo"), convey.ShouldNotBeNil)
		convey.So(t.DescTxt.DescTxt, convey.ShouldEqual, "foo")
	})
}

func TestGBuilder_GetDescription(t *testing.T) {
	convey.Convey("flag.GBuilder.GetDescription()", t, func() {
		t := &flag.GBuilder{DescTxt: trait.Described{DescTxt: "bar"}}

		convey.So(t.GetDescription(), convey.ShouldEqual, "bar")
	})
}

func TestGBuilder_Name(t *testing.T) {
	convey.Convey("flag.GBuilder.Name()", t, func() {
		t := &flag.GBuilder{}

		convey.So(t.Name("foo"), convey.ShouldNotBeNil)
		convey.So(t.NameTxt.NameTxt, convey.ShouldEqual, "foo")
	})
}

func TestGBuilder_GetName(t *testing.T) {
	convey.Convey("flag.GBuilder.GetName()", t, func() {
		t := &flag.GBuilder{NameTxt: trait.Named{NameTxt: "bar"}}

		convey.So(t.GetName(), convey.ShouldEqual, "bar")
	})
}
