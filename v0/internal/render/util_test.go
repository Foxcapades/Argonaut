package render_test

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/Foxcapades/Argonaut/v0/internal/render"
)

func TestIsBreakChar(t *testing.T) {
	convey.Convey("render.IsBreakChar()", t, func() {
		ws := []byte{' ', '\t', '\r', '\n'}
		convey.Convey("Should return true for whitespace characters", func() {

			for _, test := range ws {
				convey.Convey(fmt.Sprintf("Character %d", test), func() {
					convey.So(render.IsBreakChar(test), convey.ShouldBeTrue)
				})
			}
		})
		convey.Convey("Should return false for other characters", func() {
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
				convey.Convey(fmt.Sprintf("Character %d", test), func() {
					convey.So(render.IsBreakChar(test), convey.ShouldBeFalse)
				})
			}
		})
	})
}
