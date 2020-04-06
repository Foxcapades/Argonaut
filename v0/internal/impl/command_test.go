package impl_test

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/impl"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestCommand_Description(t *T) {
	Convey("Command.Description", t, func() {
		str := "knee deep in the hoopla"
		So(NewCommandBuilder().Description(str).MustBuild().Description(),
			ShouldEqual, str)
	})
}

func TestCommand_Name(t *T) {
	Convey("Command.Name", t, func() {
		com := Command{}
		So(com.Name(), ShouldEqual, com.String())
	})
}