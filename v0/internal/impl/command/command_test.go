package command_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"testing"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/command"
	"github.com/smartystreets/goconvey/convey"
)

func TestCommand_Description(t *testing.T) {
	convey.Convey("Command.Description", t, func() {
		str := "knee deep in the hoopla"
		convey.So((&command.Command{Described: trait.Described{DescTxt: str}}).
			Description(), convey.ShouldEqual, str)
	})
}

func TestCommand_Name(t *testing.T) {
	convey.Convey("Command.Name", t, func() {
		com := command.Command{}
		convey.So(com.Name(), convey.ShouldEqual, com.String())
	})
}

func TestCommand_Passthroughs(t *testing.T) {
	convey.Convey("Command.Passthroughs", t, func() {
		convey.Convey("Should return the self contained passthrough values", func() {
			convey.So((&command.Command{Passthrough: []string{"foo"}}).Passthroughs(),
				convey.ShouldResemble, []string{"foo"})
		})
	})
}

func TestCommand_UnmappedInput(t *testing.T) {
	convey.Convey("Command.UnmappedInput", t, func() {
		convey.Convey("Should return the self contained unmapped input values", func() {
			convey.So((&command.Command{Unmapped: []string{"foo"}}).UnmappedInput(),
				convey.ShouldResemble, []string{"foo"})
		})
	})
}

func TestCommand_AppendPassthrough(t *testing.T) {
	convey.Convey("Command.AppendPassthrough", t, func() {
		convey.Convey("Should update the self contained passthrough values", func() {
			test := &command.Command{}
			convey.So(test.Passthrough, convey.ShouldResemble, []string(nil))

			test.AppendPassthrough("bar")
			convey.So(test.Passthrough, convey.ShouldResemble, []string{"bar"})
		})
	})
}

func TestCommand_AppendUnmapped(t *testing.T) {
	convey.Convey("Command.AppendUnmapped", t, func() {
		convey.Convey("Should update the self contained passthrough values", func() {
			test := &command.Command{}
			convey.So(test.Unmapped, convey.ShouldResemble, []string(nil))

			test.AppendUnmapped("bar")
			convey.So(test.Unmapped, convey.ShouldResemble, []string{"bar"})
		})
	})
}
