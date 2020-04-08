package com_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	com2 "github.com/Foxcapades/Argonaut/v0/internal/impl/com"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestCommand_Description(t *T) {
	Convey("Command.Description", t, func() {
		str := "knee deep in the hoopla"
		So(com2.NewBuilder(impl.NewProvider()).Description(str).MustBuild().Description(),
			ShouldEqual, str)
	})
}

func TestCommand_Name(t *T) {
	Convey("Command.Name", t, func() {
		com := com2.Command{}
		So(com.Name(), ShouldEqual, com.String())
	})
}