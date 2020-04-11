package flag_test

import (
	. "testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func TestFlagBuilder_Build(t *T) {
	provider := NewProvider()

	Convey("FlagBuilder.Build", t, func() {
		Convey("No Flags", func() {
			a, b := flag.NewBuilder(provider).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)

			e, o := b.(A.FlagBuilderError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.FlagBuilderErrNoFlags)
		})

		Convey("Invalid Short Flag", func() {
			a, b := flag.NewBuilder(provider).Short(0).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.FlagBuilderError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.FlagBuilderErrBadShortFlag)
		})

		Convey("Invalid Long Flag", func() {
			a, b := flag.NewBuilder(provider).Long(" ").Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.FlagBuilderError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.FlagBuilderErrBadLongFlag)
		})

		Convey("Invalid Argument", func() {
			a, b := flag.NewBuilder(provider).Short('3').Bind(nil, false).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.ArgumentError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})
	})
}

func TestFlagBuilder_MustBuild(t *T) {
	provider := NewProvider()

	Convey("FlagBuilder.Build", t, func() {
		Convey("No Flags", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(provider).MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.FlagBuilderError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.FlagBuilderErrNoFlags)
		})

		Convey("Invalid Short Flag", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(provider).Short(0).MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.FlagBuilderError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.FlagBuilderErrBadShortFlag)
		})

		Convey("Invalid Long Flag", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(provider).Long(" ").MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.FlagBuilderError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.FlagBuilderErrBadLongFlag)
		})

		Convey("Invalid Argument", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(provider).Short('3').Bind(nil, false).MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.ArgumentError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})
	})
}
