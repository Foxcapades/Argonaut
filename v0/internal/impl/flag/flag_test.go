package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFlag_String(t *testing.T) {
	Convey("Flag.String", t, func() {
		Convey("Long flag, no arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Long("fail").MustBuild().String(),
				ShouldEqual, "--fail")
		})
		Convey("Long flag, optional arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Long("fail").Arg(argument.NewBuilder(impl.NewProvider()).Name("flump")).
				MustBuild().String(), ShouldEqual, "--fail=flump")
		})
		Convey("Long flag, required arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Long("fail").Arg(argument.NewBuilder(impl.NewProvider()).Name("flump").
				Require()).MustBuild().String(), ShouldEqual, "--fail=flump")
		})
		Convey("Short flag, no arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').MustBuild().String(),
				ShouldEqual, "-f")
		})
		Convey("Short flag, optional arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Arg(argument.NewBuilder(impl.NewProvider()).Name("flump")).
				MustBuild().String(), ShouldEqual, "-f flump")
		})
		Convey("Short flag, required arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Arg(argument.NewBuilder(impl.NewProvider()).Name("flump").
				Require()).MustBuild().String(), ShouldEqual, "-f flump")
		})
		Convey("Short & Long flag, no arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Long("fail").MustBuild().String(),
				ShouldEqual, "-f | --fail")
		})
		Convey("Short & Long flag, optional arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Long("fail").Arg(argument.NewBuilder(impl.NewProvider()).
				Name("flump")).MustBuild().String(), ShouldEqual,
				"-f flump | --fail=flump")
		})
		Convey("Short & Long flag, required arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Long("fail").Arg(argument.NewBuilder(impl.NewProvider()).
				Name("flump").Require()).MustBuild().String(), ShouldEqual,
				"-f flump | --fail=flump")
		})
	})
}

func TestFlag_Argument(t *testing.T) {
	Convey("flag.Flag.Argument", t, func() {
		arg := &argument.Argument{}
		So((&flag.Flag{ArgumentElement: arg}).Argument(), ShouldPointTo, arg)
	})
}

func TestFlag_HasArgument(t *testing.T) {
	Convey("flag.Flag.HasArgument", t, func() {
		arg := &argument.Argument{}
		So((&flag.Flag{ArgumentElement: arg}).HasArgument(), ShouldBeTrue)
		So((&flag.Flag{}).HasArgument(), ShouldBeFalse)
	})
}

func TestFlag_Long(t *testing.T) {
	Convey("flag.Flag.Long", t, func() {
		So((&flag.Flag{LongForm: "chewer"}).Long(), ShouldEqual, "chewer")
	})
}

func TestFlag_HasLong(t *testing.T) {
	Convey("flag.Flag.HasLong", t, func() {
		So((&flag.Flag{LongForm: "height"}).HasLong(), ShouldBeTrue)
		So((&flag.Flag{LongForm: ""}).HasLong(), ShouldBeFalse)
	})
}

func TestFlag_Short(t *testing.T) {
	Convey("flag.Flag.Short", t, func() {
		So((&flag.Flag{ShortForm: 'q'}).Short(), ShouldEqual, 'q')
	})
}

func TestFlag_HasShort(t *testing.T) {
	Convey("flag.Flag.HasShort", t, func() {
		So((&flag.Flag{ShortForm: 'q'}).HasShort(), ShouldBeTrue)
		So((&flag.Flag{ShortForm: 0}).HasShort(), ShouldBeFalse)
	})
}

func TestFlag_Required(t *testing.T) {
	Convey("flag.Flag.Required", t, func() {
		So((&flag.Flag{IsRequired: true}).Required(), ShouldBeTrue)
		So((&flag.Flag{}).Required(), ShouldBeFalse)
	})
}

func TestFlag_Parent(t *testing.T) {
	Convey("flag.Flag.Parent", t, func() {
		fg := &flag.Group{}
		So((&flag.Flag{ParentElement: fg}).Parent(), ShouldPointTo, fg)
	})
}

func TestFlag_Hits(t *testing.T) {
	Convey("flag.Flag.Hits", t, func() {
		So((&flag.Flag{}).Hits(), ShouldEqual, 0)
		So((&flag.Flag{HitCount: 42}).Hits(), ShouldEqual, 42)
	})
}

func TestFlag_IncrementHits(t *testing.T) {
	Convey("flag.Flag.IncrementHits", t, func() {
		Convey("Should increment internal use counter", func() {
			f := &flag.Flag{}
			f.IncrementHits()
			So(f.HitCount, ShouldEqual, 1)
		})

		Convey("Should call OnHit callback", func() {
			val := 0
			fn := func(argo.Flag) { val++ }
			(&flag.Flag{OnHitCallback: fn}).IncrementHits()
			So(val, ShouldEqual, 1)
		})

		Convey("Should increment hit binding", func() {
			val := 0
			(&flag.Flag{HitCountBinding: &val}).IncrementHits()
			So(val, ShouldEqual, 1)
		})
	})
}
