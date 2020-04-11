package render_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl/argument"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/internal/render"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
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
		)

		tests := []*struct {
			name string
			arg  Argument
			val  string
		}{
			{
				"With Name",
				Argument{Named: trait.Named{NameTxt: "Hi"}},
				"Hi",
			},
			{
				"With Int Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindInt)},
				"int",
			},
			{
				"With Unsigned Int Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindUint)},
				"uint",
			},
			{
				"With Float32 Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindFloat)},
				"float",
			},
			{
				"With Map Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindMap)},
				"string=uint",
			},
			{
				"With Slice Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindSlice)},
				"float",
			},
			{
				"With Bytes Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindBytes)},
				"bytes",
			},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				So(render.ArgName(&test.arg), ShouldEqual, test.val)
			})
		}
	})
}

func TestFormatArgName(t *T) {
	var bindType map[string][]uint8

	Convey("FormatArgName", t, func() {
		tests := []*struct {
			name string
			arg  Argument
			val  string
		}{
			{
				"Required With Name",
				Argument{Named: trait.Named{NameTxt: "Hi"}, IsRequired: true},
				"<Hi>",
			},
			{
				"Optional With Name",
				Argument{Named: trait.Named{NameTxt: "Hi"}},
				"[Hi]",
			},
			{
				"Required With Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindType), IsRequired: true},
				"<string=bytes>",
			},
			{
				"Optional With Type",
				Argument{IsBindingSet: true, RootBinding: reflect.ValueOf(bindType)},
				"[string=bytes]",
			},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				tmp := strings.Builder{}
				render.FormattedArgName(&test.arg, &tmp)
				So(tmp.String(), ShouldEqual, test.val)
			})
		}
	})
}
