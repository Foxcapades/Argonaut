package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFlag_String(t *testing.T) {
	Convey("Flag.String", t, func() {
		Convey("Long flag, no arg", func() {
			So(NewFlagBuilder().Long("fail").MustBuild().String(),
				ShouldEqual, "--fail")
		})
		Convey("Long flag, optional arg", func() {
			So(NewFlagBuilder().Long("fail").Arg(NewArgBuilder().Name("flump")).
				MustBuild().String(), ShouldEqual, "--fail=[flump]")
		})
		Convey("Long flag, required arg", func() {
			So(NewFlagBuilder().Long("fail").Arg(NewArgBuilder().Name("flump").
				Require()).MustBuild().String(), ShouldEqual, "--fail=flump")
		})
		Convey("Short flag, no arg", func() {
			So(NewFlagBuilder().Short('f').MustBuild().String(),
				ShouldEqual, "-f")
		})
		Convey("Short flag, optional arg", func() {
			So(NewFlagBuilder().Short('f').Arg(NewArgBuilder().Name("flump")).
				MustBuild().String(), ShouldEqual, "-f [flump]")
		})
		Convey("Short flag, required arg", func() {
			So(NewFlagBuilder().Short('f').Arg(NewArgBuilder().Name("flump").
				Require()).MustBuild().String(), ShouldEqual, "-f flump")
		})
		Convey("Short & Long flag, no arg", func() {
			So(NewFlagBuilder().Short('f').Long("fail").MustBuild().String(),
				ShouldEqual, "-f | --fail")
		})
		Convey("Short & Long flag, optional arg", func() {
			So(NewFlagBuilder().Short('f').Long("fail").Arg(NewArgBuilder().
				Name("flump")).MustBuild().String(), ShouldEqual,
				"-f [flump] | --fail=[flump]")
		})
		Convey("Short & Long flag, required arg", func() {
			So(NewFlagBuilder().Short('f').Long("fail").Arg(NewArgBuilder().
				Name("flump").Require()).MustBuild().String(), ShouldEqual,
				"-f flump | --fail=flump")
		})
	})
}