package render_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/Foxcapades/Argonaut/v0/internal/render"
)

func TestIsBreakChar(t *testing.T) {
	Convey("render.IsBreakChar()", t, func() {
		ws := []byte{' ', '\t', '\r', '\n'}
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

	})
}
