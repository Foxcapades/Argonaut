package render_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	"github.com/Foxcapades/Argonaut/v0/internal/render"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	. "testing"
)

func TestArgName(t *T) {
	Convey("ArgName", t, func() {
		var (
			bindInt   int
			bindUint  uint
			bindFloat float32
			bindMap   map[string]uint
			bindSlice []float32
			bindBytes []byte
			bindPtr   *string
		)

		tests := []*struct {
			name string
			arg  A.ArgumentBuilder
			val  string
		}{
			{
				"With Name",
				NewBuilder(NewProvider()).Name("Hi"),
				"Hi",
			},
			{
				"With Int Type",
				NewBuilder(NewProvider()).Bind(&bindInt),
				"int",
			},
			{
				"With Unsigned Int Type",
				NewBuilder(NewProvider()).Bind(&bindUint),
				"uint",
			},
			{
				"With Float32 Type",
				NewBuilder(NewProvider()).Bind(&bindFloat),
				"float",
			},
			{
				"With Map Type",
				NewBuilder(NewProvider()).Bind(&bindMap),
				"string=uint",
			},
			{
				"With Slice Type",
				NewBuilder(NewProvider()).Bind(&bindSlice),
				"float",
			},
			{
				"With Bytes Type",
				NewBuilder(NewProvider()).Bind(&bindBytes),
				"bytes",
			},
			{
				"With String Pointer Type",
				NewBuilder(NewProvider()).Bind(&bindPtr),
				"string",
			},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				arg, err := test.arg.Build()
				So(err, ShouldBeNil)
				So(render.ArgName(arg), ShouldEqual, test.val)
			})
		}
	})
}

func TestFormatArgName(t *T) {
	var bindType map[string][]uint8

	Convey("FormatArgName", t, func() {
		tests := []*struct {
			name string
			arg  A.ArgumentBuilder
			val  string
		}{
			{
				"Required With Name",
				NewBuilder(NewProvider()).Name("Hi").Require(),
				"<Hi>",
			},
			{
				"Optional With Name",
				NewBuilder(NewProvider()).Name("Hi"),
				"[Hi]",
			},
			{
				"Required With Type",
				NewBuilder(NewProvider()).Bind(&bindType).Require(),
				"<string=bytes>",
			},
			{
				"Optional With Type",
				NewBuilder(NewProvider()).Bind(&bindType),
				"[string=bytes]",
			},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				arg, err := test.arg.Build()
				So(err, ShouldBeNil)
				tmp := strings.Builder{}
				render.FormattedArgName(arg, &tmp)
				So(tmp.String(), ShouldEqual, test.val)
			})
		}
	})
}
