package util_test

import (
	"errors"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMust(t *testing.T) {
	Convey("util.Must()", t, func() {
		So(func() {util.Must(errors.New("derp"))}, ShouldPanic)
		So(func() {util.Must(nil)}, ShouldNotPanic)
	})
}
