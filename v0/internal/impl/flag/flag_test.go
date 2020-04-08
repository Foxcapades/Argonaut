package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/arg"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
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
			So(flag.NewBuilder(impl.NewProvider()).Long("fail").Arg(arg.NewBuilder(impl.NewProvider()).Name("flump")).
				MustBuild().String(), ShouldEqual, "--fail=[flump]")
		})
		Convey("Long flag, required arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Long("fail").Arg(arg.NewBuilder(impl.NewProvider()).Name("flump").
				Require()).MustBuild().String(), ShouldEqual, "--fail=flump")
		})
		Convey("Short flag, no arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').MustBuild().String(),
				ShouldEqual, "-f")
		})
		Convey("Short flag, optional arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Arg(arg.NewBuilder(impl.NewProvider()).Name("flump")).
				MustBuild().String(), ShouldEqual, "-f [flump]")
		})
		Convey("Short flag, required arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Arg(arg.NewBuilder(impl.NewProvider()).Name("flump").
				Require()).MustBuild().String(), ShouldEqual, "-f flump")
		})
		Convey("Short & Long flag, no arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Long("fail").MustBuild().String(),
				ShouldEqual, "-f | --fail")
		})
		Convey("Short & Long flag, optional arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Long("fail").Arg(arg.NewBuilder(impl.NewProvider()).
				Name("flump")).MustBuild().String(), ShouldEqual,
				"-f [flump] | --fail=[flump]")
		})
		Convey("Short & Long flag, required arg", func() {
			So(flag.NewBuilder(impl.NewProvider()).Short('f').Long("fail").Arg(arg.NewBuilder(impl.NewProvider()).
				Name("flump").Require()).MustBuild().String(), ShouldEqual,
				"-f flump | --fail=flump")
		})
	})
}