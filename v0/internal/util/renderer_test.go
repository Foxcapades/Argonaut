package util_test

import (
	"fmt"
	. "github.com/Foxcapades/Argonaut/v0/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	. "testing"
)

func TestBreakFmt(t *T) {
	Convey("BreakFmt", t, func() {
		out := strings.Builder{}
		BreakFmt("foobarfizzbuzzfoobarfizzbuzzfoobarfizzbuzz", 10, 50, &out)
		fmt.Println(out.String())
	})
}
