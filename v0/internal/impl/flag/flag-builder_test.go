package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/arg"
	. "testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func TestFlagBuilder_Arg(t *T) {
	prov := NewProvider()

	Convey("FlagBuilder.Arg", t, func() {
		a := arg.NewBuilder(prov).TypeHint("diced bagels")
		b := flag.NewBuilder(prov).Arg(a).Short('a').MustBuild()
		So(b.Argument(), ShouldResemble, a.MustBuild())
	})
}

func TestFlagBuilder_Bind(t *T) {
	prov := NewProvider()

	Convey("FlagBuilder.Bind", t, func() {
		p := "is this even a good game?"

		Convey("required", func() {
			b := flag.NewBuilder(prov).Bind(&p, true).Short('a').MustBuild()
			So(b.Argument().Required(), ShouldBeTrue)
			So(b.Argument().(*arg.Argument).Binding(), ShouldPointTo, &p)
		})

		Convey("not required", func() {
			b := flag.NewBuilder(prov).Bind(&p, false).Short('a').MustBuild()
			So(b.Argument().Required(), ShouldBeFalse)
			So(b.Argument().(*arg.Argument).Binding(), ShouldPointTo, &p)
		})
	})
}

func TestFlagBuilder_Short(t *T) {
	Convey("FlagBuilder.Short", t, func() {
		f := flag.NewBuilder(NewProvider()).Short('z').MustBuild()
		So(f.Short(), ShouldEqual, 'z')
		So(f.HasShort(), ShouldBeTrue)
	})
}

func TestFlagBuilder_Long(t *T) {
	Convey("FlagBuilder.Long", t, func() {
		f := flag.NewBuilder(NewProvider()).Long("smerty").MustBuild()
		So(f.Long(), ShouldEqual, "smerty")
		So(f.HasLong(), ShouldBeTrue)
	})
}

func TestFlagBuilder_Description(t *T) {
	Convey("FlagBuilder.Description", t, func() {
		f := flag.NewBuilder(NewProvider()).Short('a').Description("bananas are superior to mangos").MustBuild()
		So(f.Description(), ShouldEqual, "bananas are superior to mangos")
		So(f.HasDescription(), ShouldBeTrue)
	})
}

func TestFlagBuilder_Default(t *T) {
	Convey("FlagBuilder.Default", t, func() {
		p := "i'm pretty sure this is a bad game"

		Convey("required", func() {
			b := flag.NewBuilder(NewProvider()).Default(&p).Short('a').MustBuild()
			So(b.Argument().Required(), ShouldBeFalse)
			So(b.Argument().(*arg.Argument).Default(), ShouldPointTo, &p)
		})
	})
}

func TestFlagBuilder_Build(t *T) {
	Convey("FlagBuilder.Build", t, func() {
		Convey("No Flags", func() {
			a, b := flag.NewBuilder(NewProvider()).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.InvalidFlagError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.InvalidFlagNoFlags)
		})

		Convey("Invalid Short Flag", func() {
			a, b := flag.NewBuilder(NewProvider()).Short(0).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.InvalidFlagError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.InvalidFlagBadShortFlag)
		})

		Convey("Invalid Long Flag", func() {
			a, b := flag.NewBuilder(NewProvider()).Long(" ").Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.InvalidFlagError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.InvalidFlagBadLongFlag)
		})

		Convey("Invalid Argument", func() {
			a, b := flag.NewBuilder(NewProvider()).Short('3').Bind(nil, false).Build()
			So(a, ShouldBeNil)
			So(b, ShouldNotBeNil)
			e, o := b.(A.ArgumentError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})
	})
}

func TestFlagBuilder_MustBuild(t *T) {
	Convey("FlagBuilder.Build", t, func() {
		Convey("No Flags", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(NewProvider()).MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.InvalidFlagError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.InvalidFlagNoFlags)
		})

		Convey("Invalid Short Flag", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(NewProvider()).Short(0).MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.InvalidFlagError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.InvalidFlagBadShortFlag)
		})

		Convey("Invalid Long Flag", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(NewProvider()).Long(" ").MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.InvalidFlagError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.InvalidFlagBadLongFlag)
		})

		Convey("Invalid Argument", func() {
			var err interface{}
			fn := func() {
				defer func() { err = recover(); panic(err) }()
				flag.NewBuilder(NewProvider()).Short('3').Bind(nil, false).MustBuild()
			}

			So(fn, ShouldPanic)
			e, o := err.(A.ArgumentError)
			So(o, ShouldBeTrue)
			So(e.Type(), ShouldEqual, A.ArgErrInvalidBindingBadType)
		})
	})
}
