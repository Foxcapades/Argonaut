package util

import (
	C "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSameTypeOrPointerTo(t *testing.T) {
	C.Convey("Compatible", t, func() {
		v := ""
		C.So(Compatible(v, ""), C.ShouldBeTrue)
		C.So(Compatible("", &v), C.ShouldBeTrue)
	})
}
