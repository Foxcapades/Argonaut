package impl_test

import (
	"errors"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
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
				ctx := parserTestCtx{parser: NewParser()}
				com := NewCommandBuilder()
				t.Setup(com)
				ctx.error = ctx.parser.Parse(append([]string{"cmd"}, t.Params()...),
					com.MustBuild())
				t.Test(&ctx)
			})
		}

		Convey("StringSlice", func() {
			var derp []string
			com := NewCommandBuilder().
				Flag(NewFlagBuilder().
					Short('a').
					Long("aa").
					Bind(&derp, true)).
				MustBuild()
			input := []string{"bar", "-a", "fizz", "-abuzz", "--aa=pong"}
			parser := NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldBeNil)
			So(parser.Unrecognized(), ShouldBeEmpty)
			So(parser.Passthroughs(), ShouldBeEmpty)
			So(derp, ShouldResemble, []string{"fizz", "buzz", "pong"})
		})

		Convey("Short flag with required arg with no value", func() {
			Convey("Type is bool", func() {
				derp := false
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-a"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldBeTrue)
			})

			Convey("Type is *bool", func() {
				var derp *bool
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-a"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeNil)
				So(*derp, ShouldBeTrue)
			})

			Convey("Type is []bool", func() {
				var derp []bool
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-a"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeEmpty)
				So(derp[0], ShouldBeTrue)
			})

			Convey("Type is []*bool", func() {
				var derp []*bool
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Short('a').
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "-a"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeEmpty)
				So(derp[0], ShouldNotBeNil)
				So(*derp[0], ShouldBeTrue)
			})
		})

		Convey("Long flag with required arg with no value", func() {
			Convey("Type is bool", func() {
				derp := false
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldBeTrue)
			})

			Convey("Type is *bool", func() {
				var derp *bool
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeNil)
				So(*derp, ShouldBeTrue)
			})

			Convey("Type is []bool", func() {
				var derp []bool
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeEmpty)
				So(derp[0], ShouldBeTrue)
			})

			Convey("Type is []*bool", func() {
				var derp []*bool
				com := NewCommandBuilder().
					Flag(NewFlagBuilder().
						Long("applies").
						Bind(&derp, true)).
					MustBuild()
				input := []string{"bar", "--applies"}
				parser := NewParser()
				err := parser.Parse(input, com)
				So(err, ShouldBeNil)
				So(parser.Unrecognized(), ShouldBeEmpty)
				So(parser.Passthroughs(), ShouldBeEmpty)
				So(derp, ShouldNotBeEmpty)
				So(derp[0], ShouldNotBeNil)
				So(*derp[0], ShouldBeTrue)
			})
		})

		Convey("Required arg not provided", func() {
			var str string
			com := NewCommandBuilder().
				Arg(NewArgBuilder().Bind(&str).Require()).
				MustBuild()
			input := []string{"bar"}
			parser := NewParser()
			err := parser.Parse(input, com)
			So(err, ShouldResemble, errors.New("missing required params"))
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
	bld := NewFlagBuilder().Bind(&s.text, s.reqArg)
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
