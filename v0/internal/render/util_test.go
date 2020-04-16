package render_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/Foxcapades/Argonaut/v0/internal/render"
)

func TestIsBreakChar(t *testing.T) {
	Convey("render.IsBreakChar()", t, func() {
		ws := []byte{' ', '\t'}
		Convey("Should return true for whitespace characters", func() {

			for _, test := range ws {
				Convey(fmt.Sprintf("Character %d", test), func() {
					So(render.IsBreakChar(test), ShouldBeTrue)
				})
			}
		})
		Convey("Should return false for other characters", func() {
			tests := make([]byte, 0, 256-len(ws))

		outer:
			for i := 0; i < 256; i++ {
				b := byte(i)
				for _, c := range ws {
					if c == b {
						continue outer
					}
				}
				tests = append(tests, b)
			}

			for _, test := range tests {
				Convey(fmt.Sprintf("Character %d", test), func() {
					So(render.IsBreakChar(test), ShouldBeFalse)
				})
			}
		})
	})
}

func TestBreakFmt(t *testing.T) {
	Convey("render.BreakFmt()", t, func() {
		Convey("long unbroken string", func() {
			val := "abcdefghijklmnopqrstuvwxyz"
			exp := "abc-\ndef-\nghi-\njkl-\nmno-\npqr-\nstu-\nvwx-\nyz"
			tmp := strings.Builder{}

			render.BreakFmt(val, 0, 4, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
		Convey("2 long strings", func() {
			val := "abcdefghijklmnopqrstuvwxyz abcdefghijklmnopqrstuvwxyz"
			exp := "abc-\ndef-\nghi-\njkl-\nmno-\npqr-\nstu-\nvwx-\nyz\nabc-\ndef-\nghi-\njkl-\nmno-\npqr-\nstu-\nvwx-\nyz"
			tmp := strings.Builder{}

			render.BreakFmt(val, 0, 4, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
		Convey("Multi words", func() {
			val := "even a small car like that can do some damage"
			exp := "even\n  a\n  sma-\n  ll\n  car\n  like\n  that\n  can\n  do\n  some\n  dam-\n  age"
			tmp := strings.Builder{}

			render.BreakFmt(val, 2, 6, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
		Convey("single char width", func() {
			val := "heckin chonker"
			exp := "h\ne\nc\nk\ni\nn\nc\nh\no\nn\nk\ne\nr"
			tmp := strings.Builder{}

			render.BreakFmt(val, 0, 1, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
		Convey("under width", func() {
			val := "heckin chonker"
			exp := "heckin chonker"
			tmp := strings.Builder{}

			render.BreakFmt(val, 0, 14, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
	})
}

func TestWritePadded(t *testing.T) {
	Convey("render.WritePadded", t, func() {
		Convey("over width", func() {
			tmp := strings.Builder{}
			val := "finally that poor creature has come to rest"
			exp := val

			render.WritePadded(val, 10, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
		Convey("under width", func() {
			tmp := strings.Builder{}
			val := "ramp it"
			exp := "ramp it   "

			render.WritePadded(val, 10, &tmp)
			So(tmp.String(), ShouldEqual, exp)
		})
	})
}
