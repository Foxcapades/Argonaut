package parse_test

import (
	"errors"
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	com2 "github.com/Foxcapades/Argonaut/v0/internal/impl/command"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/parse"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestParser_Parse(t *T) {
	Convey("Parser.Parse", t, func() {
		tests := []ParserTest{
			&stringParserTest{
				name:   "Long flag with value",
				expect: "fizz",
				params: []string{"--foo=fizz"},
				long:   "foo",
			},
			&stringParserTest{
				name:   "Required long flag with value",
				expect: "fizz",
				params: []string{"--foo=fizz"},
				long:   "foo",
				reqArg: true,
			},
			&stringParserTest{
				name:   "Short disconnected flag with value",
				short:  'a',
				expect: "fizz",
				params: []string{"-a", "fizz"},
			},
			&stringParserTest{
				name:   "Short disconnected flag with required value",
				short:  'a',
				expect: "fizz",
				params: []string{"-a", "fizz"},
				reqArg: true,
			},
			&stringParserTest{
				name:   "Short connected flag with value",
				short:  'a',
				expect: "fizz",
				params: []string{"-afizz"},
			},
			&stringParserTest{
				name:   "Short connected flag with required value",
				short:  'a',
				expect: "fizz",
				params: []string{"-afizz"},
				reqArg: true,
			},
		}

		for _, t := range tests {
			Convey(t.Name(), func() {
				ctx := parserTestCtx{parser: parse.NewParser()}
				com := com2.NewBuilder(impl.NewProvider())
				t.Setup(com)
				ctx.error = ctx.parser.Parse(append([]string{"cmd"}, t.Params()...),
					com.MustBuild())
				t.Test(&ctx)
			})
		}

		Convey("StringSlice", func() {
			var derp []string
			com := com2.NewBuilder(impl.NewProvider()).
				Flag(flag.NewBuilder(impl.NewProvider()).
					Short('a').
					Long("aa").
					Bind(&derp, true)).
				MustBuild()
			input := []string{"bar", "-a", "fizz", "-abuzz", "--aa=pong"}
			parser := parse.NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldBeNil)
			So(parser.Unrecognized(), ShouldBeEmpty)
			So(parser.Passthroughs(), ShouldBeEmpty)
			So(derp, ShouldResemble, []string{"fizz", "buzz", "pong"})
		})

		Convey("Short flag with required arg with no value", func() {
			Convey("Type is bool", func() {
				derp := false
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-a"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldBeTrue)
			})

			Convey("Type is *bool", func() {
				var derp *bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-a"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeNil)
				So(*derp, ShouldBeTrue)
			})

			Convey("Type is []bool", func() {
				var derp []bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-aa"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(len(derp), ShouldEqual, 2)
				So(derp[0], ShouldBeTrue)
				So(derp[1], ShouldBeTrue)
			})

			Convey("Type is []*bool, required is true", func() {
				var derp []*bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-aa", "-a"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(len(derp), ShouldEqual, 3)
				So(derp[0], ShouldNotBeNil)
				So(*derp[0], ShouldBeTrue)
				So(derp[1], ShouldNotBeNil)
				So(*derp[1], ShouldBeTrue)
				So(derp[2], ShouldNotBeNil)
				So(*derp[2], ShouldBeTrue)
			})

			Convey("Type is []*bool, required is false", func() {
				var derp []*bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Short('a').
						Bind(&derp, false)).
					MustBuild()
				input := []string{"bar", "-aa", "-a"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(len(derp), ShouldEqual, 3)
				So(derp[0], ShouldNotBeNil)
				So(*derp[0], ShouldBeTrue)
				So(derp[1], ShouldNotBeNil)
				So(*derp[1], ShouldBeTrue)
				So(derp[2], ShouldNotBeNil)
				So(*derp[2], ShouldBeTrue)
			})
		})

		Convey("Long flag with required arg with no value", func() {
			Convey("Type is bool", func() {
				derp := false
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldBeTrue)
			})

			Convey("Type is *bool", func() {
				var derp *bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeNil)
				So(*derp, ShouldBeTrue)
			})

			Convey("Type is []bool", func() {
				var derp []bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeEmpty)
				So(derp[0], ShouldBeTrue)
			})

			Convey("Type is []*bool", func() {
				var derp []*bool
				com := com2.NewBuilder(impl.NewProvider()).
					Flag(flag.NewBuilder(impl.NewProvider()).
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies", "--applies"}
				parser := parse.NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(len(derp), ShouldEqual, 2)
				So(derp[0], ShouldNotBeNil)
				So(*derp[0], ShouldBeTrue)
				So(derp[1], ShouldNotBeNil)
				So(*derp[1], ShouldBeTrue)
			})
		})

		Convey("Required arg not provided", func() {
			var str string
			com := com2.NewBuilder(impl.NewProvider()).
				Arg(argument.NewBuilder(impl.NewProvider()).Bind(&str).Require()).
				MustBuild()
			input := []string{"bar"}
			parser := parse.NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldResemble, errors.New("missing required params"))
		})

		Convey("Required arg provided with leading flags", func() {
			var str string
			var no bool
			com := com2.NewBuilder(impl.NewProvider()).
				Arg(argument.NewBuilder(impl.NewProvider()).Bind(&str).Require()).
				Flag(flag.NewBuilder(impl.NewProvider()).Short('v').Bind(&no, false)).
				MustBuild()
			input := []string{"bar", "-vv", "value"}
			parser := parse.NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "value")
			So(no, ShouldBeTrue)
		})

	})
}

type parserTestCtx struct {
	parser A.Parser
	error  error
}

type ParserTest interface {
	Name() string
	Params() []string
	Setup(com A.CommandBuilder)
	Test(ctx *parserTestCtx)
}

type stringParserTest struct {
	name    string
	text    string
	expect  string
	reqArg  bool
	reqFlag bool
	short   byte
	long    string
	params  []string
}

func (s *stringParserTest) Name() string {
	return s.name
}

func (s *stringParserTest) Params() []string {
	return s.params
}

func (s *stringParserTest) Setup(com A.CommandBuilder) {
	bld := flag.NewBuilder(impl.NewProvider()).Bind(&s.text, s.reqArg)
	if s.short != 0 {
		bld.Short(s.short)
	}
	if s.long != "" {
		bld.Long(s.long)
	}
	com.Flag(bld)
}

func (s *stringParserTest) Test(ctx *parserTestCtx) {
	So(ctx.error, ShouldBeNil)
	So(ctx.parser.Unrecognized(), ShouldBeEmpty)
	So(ctx.parser.Passthroughs(), ShouldBeEmpty)
	So(s.text, ShouldEqual, s.expect)
}
