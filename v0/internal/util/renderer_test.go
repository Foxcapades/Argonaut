package util_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	. "testing"
)

func TestBreakFmt(t *T) {
	Convey("BreakFmt", t, func() {
		Convey("Long unbroken string", func() {
			out := strings.Builder{}
			BreakFmt("foobarfizzbuzzfoobarfizzbuzzfoobarfizzbuzz", 10, 15, &out)

			So(out.String(), ShouldEqual, `foob-
          arfi-
          zzbu-
          zzfo-
          obar-
          fizz-
          buzz-
          foob-
          arfi-
          zzbu-
          zz`)
		})
		Convey("Broken string", func() {
			out := strings.Builder{}
			BreakFmt("foo bar fizz buzz foo bar fizz buzz foo bar fizz buzz", 10, 18, &out)

			So(out.String(), ShouldEqual, `foo bar
          fizz
          buzz foo
          bar fizz
          buzz foo
          bar fizz
          buzz`)
		})
		Convey("1 char width", func() {
			out := strings.Builder{}
			BreakFmt("foobarfizzbuzzfoobar fizzbuzzfoobarfizzbuzz", 10, 11, &out)

			So(out.String(), ShouldEqual, `f
          o
          o
          b
          a
          r
          f
          i
          z
          z
          b
          u
          z
          z
          f
          o
          o
          b
          a
          r
          f
          i
          z
          z
          b
          u
          z
          z
          f
          o
          o
          b
          a
          r
          f
          i
          z
          z
          b
          u
          z
          z`)
		})
	})
}
