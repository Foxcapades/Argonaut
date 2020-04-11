package flag_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/command"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGroup_Flags(t *testing.T) {
	convey.Convey("flag.Group.Flags()", t, func() {
		test := &flag.Group{}
		val1 := &flag.Flag{}

		convey.So(test.Flags(), convey.ShouldResemble, []argo.Flag(nil))

		test.FlagElements = []argo.Flag{val1}
		convey.So(test.Flags(), convey.ShouldResemble, []argo.Flag{val1})
	})
}

func TestGroup_HasFlags(t *testing.T) {
	convey.Convey("flag.Group.HasFlags()", t, func() {
		test := &flag.Group{}
		val1 := &flag.Flag{}

		convey.So(test.HasFlags(), convey.ShouldBeFalse)

		test.FlagElements = []argo.Flag{val1}
		convey.So(test.HasFlags(), convey.ShouldBeTrue)
	})
}

func TestGroup_Parent(t *testing.T) {
	convey.Convey("flag.Group.Parent()", t, func() {
		val1 := &command.Command{}
		test := &flag.Group{ParentElement: val1}

		convey.So(test.Parent(), convey.ShouldEqual, val1)
	})
}
