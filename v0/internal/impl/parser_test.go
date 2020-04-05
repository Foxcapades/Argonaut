package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestParser_Parse(t *T) {
	Convey("Parser.Parse", t, func() {
		Convey("Long flag with value", func() {
			derp := ""
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Long("foo").
					Bind(&derp, false)).
				MustBuild()
			input := []string{"bar", "--foo=fizz"}
			err := NewParser().Parse(input, com)
			So(err, ShouldBeNil)
			So(derp, ShouldEqual, "fizz")
		})

		Convey("Required long flag with value", func() {
			derp := ""
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Long("foo").
					Bind(&derp, true)).
				MustBuild()
			input := []string{"bar", "--foo=fizz"}
			err := NewParser().Parse(input, com)
			So(err, ShouldBeNil)
			So(derp, ShouldEqual, "fizz")
		})

		Convey("Short disconnected flag with value", func() {
			derp := ""
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Short('a').
					Bind(&derp, false)).
				MustBuild()
			input := []string{"bar", "-a", "fizz"}
			err := NewParser().Parse(input, com)
			So(err, ShouldBeNil)
			So(derp, ShouldEqual, "fizz")
		})

		Convey("Required short disconnected flag with value", func() {
			derp := ""
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Short('a').
					Bind(&derp, true)).
				MustBuild()
			input := []string{"bar", "-a", "fizz"}
			parser := NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldBeNil)
			So(derp, ShouldEqual, "fizz")
		})

		Convey("Short connected flag with value", func() {
			derp := ""
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Short('a').
					Bind(&derp, false)).
				MustBuild()
			input := []string{"bar", "-afizz"}
			err := NewParser().Parse(input, com)
			So(err, ShouldBeNil)
			So(derp, ShouldEqual, "fizz")
		})

		Convey("Required short connected flag with value", func() {
			derp := ""
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Short('a').
					Bind(&derp, true)).
				MustBuild()
			input := []string{"bar", "-afizz"}
			parser := NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldBeNil)
			So(derp, ShouldEqual, "fizz")
		})
	})
}
