package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/command"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
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

func TestGBuilder_GetFlags(t *testing.T) {
	convey.Convey("flag.GBuilder.GetFlags()", t, func() {
		f := &flag.Builder{}
		t := &flag.GBuilder{FlagNodes: []argo.FlagBuilder{f}}

		convey.So(t.GetFlags(), convey.ShouldResemble, []argo.FlagBuilder{f})
	})
}

func TestGBuilder_Parent(t *testing.T) {
	convey.Convey("flag.GBuilder.Parent()", t, func() {
		c := &command.Command{}
		t := &flag.GBuilder{}

		convey.So(t.Parent(c), convey.ShouldNotBeNil)
		convey.So(t.ParentNode, convey.ShouldEqual, c)
	})
}

func TestGBuilder_Flag(t *testing.T) {
	convey.Convey("flag.GBuilder.Flag()", t, func() {
		convey.Convey("Nil input", func() {
			t := &flag.GBuilder{}

			convey.So(t.WarningVals, convey.ShouldBeEmpty)
			convey.So(t.Flag(nil), convey.ShouldNotBeNil)
			convey.So(t.FlagNodes, convey.ShouldBeEmpty)
			convey.So(t.WarningVals, convey.ShouldNotBeEmpty)
		})
		convey.Convey("Non nil input", func() {
			t := &flag.GBuilder{}
			f := &flag.Builder{}

			convey.So(t.Flag(f), convey.ShouldNotBeNil)
			convey.So(t.FlagNodes, convey.ShouldNotBeEmpty)
			convey.So(t.WarningVals, convey.ShouldBeEmpty)
		})
		c := &command.Command{}
		t := &flag.GBuilder{}

		convey.So(t.Parent(c), convey.ShouldNotBeNil)
		convey.So(t.ParentNode, convey.ShouldEqual, c)
	})
}

func TestGBuilder_Build(t *testing.T) {
	convey.Convey("flag.GBuilder.Build()", t, func() {
		convey.Convey("0 flags", func() {
			b := flag.GBuilder{}

			res, err := b.Build()
			convey.So(res, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)

			// TODO: Group builder errors
			//cst, ok := err.(argo.FlagBuilderError)
			//convey.So(ok, convey.ShouldBeTrue)
			//convey.So(cst, convey.ShouldEqual, argo.FlagBuilderErrNoFlags)
		})

		convey.Convey("Invalid Flag", func() {
			f := &flag.Builder{}
			b := flag.GBuilder{FlagNodes: []argo.FlagBuilder{f}}

			res, err := b.Build()
			convey.So(res, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)

			cst, ok := err.(argo.FlagBuilderError)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(cst.Type(), convey.ShouldEqual, argo.FlagBuilderErrNoFlags)
		})

		convey.Convey("Valid Flag", func() {
			f := &flag.Builder{ShortFlag: 'c', IsShortSet: true}
			b := flag.GBuilder{FlagNodes: []argo.FlagBuilder{f}}

			res, err := b.Build()
			convey.So(res, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			g := res.(*flag.Group)
			convey.So(len(g.FlagNodes), convey.ShouldEqual, 1)

			cf := g.FlagNodes[0].(*flag.Flag)
			convey.So(cf.ShortForm, convey.ShouldEqual, 'c')
		})
	})
}
