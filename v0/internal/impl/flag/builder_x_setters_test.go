package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"testing"

	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	"github.com/smartystreets/goconvey/convey"
)

func TestFlagBuilder_Arg(t *testing.T) {
	convey.Convey("FlagBuilder.Arg", t, func() {
		a := argument.NewBuilder(impl.NewProvider())
		b := &flag.Builder{}

		convey.So(b.Arg(a), convey.ShouldNotBeNil)
		convey.So(b.ArgBuilder, convey.ShouldEqual, a)
	})
}

func TestFlagBuilder_Bind(t *testing.T) {
	provider := impl.NewProvider()

	convey.Convey("FlagBuilder.Bind", t, func() {
		p := "is this even a good game?"

		convey.Convey("required", func() {
			b := flag.Builder{Provider: provider}

			convey.So(b.Bind(&p, true), convey.ShouldNotBeNil)
			convey.So(b.ArgBuilder.IsRequired(), convey.ShouldBeTrue)
			convey.So(b.ArgBuilder.GetBinding(), convey.ShouldPointTo, &p)
		})

		convey.Convey("not required", func() {
			b := flag.Builder{Provider: provider}

			convey.So(b.Bind(&p, false), convey.ShouldNotBeNil)
			convey.So(b.ArgBuilder.IsRequired(), convey.ShouldBeFalse)
			convey.So(b.ArgBuilder.GetBinding(), convey.ShouldPointTo, &p)
		})
	})
}

func TestFlagBuilder_Short(t *testing.T) {
	convey.Convey("FlagBuilder.Short", t, func() {
		f := flag.Builder{}

		convey.So(f.Short('z'), convey.ShouldNotBeNil)
		convey.So(f.ShortFlag, convey.ShouldEqual, 'z')
		convey.So(f.IsShortSet, convey.ShouldBeTrue)
	})
}

func TestFlagBuilder_Long(t *testing.T) {
	convey.Convey("FlagBuilder.Long", t, func() {
		f := flag.Builder{}

		convey.So(f.Long("smerty"), convey.ShouldNotBeNil)
		convey.So(f.LongFlag, convey.ShouldEqual, "smerty")
		convey.So(f.IsLongSet, convey.ShouldBeTrue)
	})
}

func TestFlagBuilder_Description(t *testing.T) {
	convey.Convey("FlagBuilder.Description", t, func() {
		f := flag.Builder{}

		convey.So(f.Description("bananas are superior to mangos"),
			convey.ShouldNotBeNil)
		convey.So(f.DescriptionText.DescriptionText, convey.ShouldEqual,
			"bananas are superior to mangos")
	})
}

func TestFlagBuilder_Default(t *testing.T) {
	provider := impl.NewProvider()

	convey.Convey("FlagBuilder.Default", t, func() {
		p := "i'm pretty sure this is a bad game"

		convey.Convey("required", func() {
			b := flag.Builder{Provider: provider}

			convey.So(b.Default(&p), convey.ShouldNotBeNil)
			convey.So(b.ArgBuilder.IsRequired(), convey.ShouldBeFalse)
			convey.So(b.ArgBuilder.GetDefault(), convey.ShouldPointTo, &p)
		})
	})
}

func TestBuilder_BindUseCount(t *testing.T) {
	convey.Convey("FlagBuilder.Long", t, func() {
		f := flag.Builder{}
		u := 0

		convey.So(f.BindUseCount(&u), convey.ShouldNotBeNil)
		convey.So(f.UseCountBinding, convey.ShouldPointTo, &u)
	})
}

func TestBuilder_OnHit(t *testing.T) {
	convey.Convey("FlagBuilder.Long", t, func() {
		f := flag.Builder{}
		o := func(argo.Flag) {}

		convey.So(f.OnHit(o), convey.ShouldNotBeNil)
		convey.So(f.OnHitCallback, convey.ShouldEqual, o)
	})
}

func TestBuilder_Parent(t *testing.T) {
	convey.Convey("FlagBuilder.Long", t, func() {
		f := flag.Builder{}
		g := &flag.Group{}

		convey.So(f.Parent(g), convey.ShouldNotBeNil)
		convey.So(f.ParentElement, convey.ShouldEqual, g)
	})
}
