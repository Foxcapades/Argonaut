package util_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsDecimal(t *testing.T) {
	Convey("util.IsDecimal()", t, func() {
		tests := make(map[byte]bool, 256)
		for i := 0; i < 256; i++ {
			b := byte(i)
			tests[b] = b >= '0' && b <= '9'
		}

		for k, v := range tests {
			So(util.IsDecimal(k), ShouldEqual, v)
		}
	})
}

func TestIsLowerLetter(t *testing.T) {
	Convey("util.IsLowerLetter()", t, func() {
		tests := make(map[byte]bool, 256)
		for i := 0; i < 256; i++ {
			b := byte(i)
			tests[b] = b >= 'a' && b <= 'z'
		}

		for k, v := range tests {
			So(util.IsLowerLetter(k), ShouldEqual, v)
		}
	})
}

func TestIsUpperLetter(t *testing.T) {
	Convey("util.IsUpperLetter()", t, func() {
		tests := make(map[byte]bool, 256)
		for i := 0; i < 256; i++ {
			b := byte(i)
			tests[b] = b >= 'A' && b <= 'Z'
		}

		for k, v := range tests {
			So(util.IsUpperLetter(k), ShouldEqual, v)
		}
	})
}
