package com_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	arg2 "github.com/Foxcapades/Argonaut/v0/internal/impl/arg"
	com2 "github.com/Foxcapades/Argonaut/v0/internal/impl/com"
	. "github.com/smartystreets/goconvey/convey"
	R "reflect"
	. "testing"
)

func TestCommandBuilder_Arg(t *T) {
	Convey("CommandBuilder.Arg", t, func() {
		Convey("Non-nil arg", func() {
			arg := arg2.NewBuilder(impl.NewProvider())
			com := com2.NewBuilder(impl.NewProvider())
			So(com.GetArgs(), ShouldBeEmpty)
			So(R.ValueOf(com.Arg(arg)).Pointer(), ShouldEqual, R.ValueOf(com).Pointer())
			So(com.GetArgs(), ShouldNotBeEmpty)
			So(R.ValueOf(com.GetArgs()[0]).Pointer(), ShouldEqual, R.ValueOf(arg).Pointer())
		})

		Convey("Nil arg", func() {
			com := com2.NewBuilder(impl.NewProvider()).Arg(nil)
			So(com.GetArgs(), ShouldBeEmpty)
			So(len(com.Warnings()), ShouldEqual, 1)
		})
	})
}

func TestCommandBuilder_Description(t *T) {
	Convey("CommandBuilder.Description", t, func() {
		com := com2.NewBuilder(impl.NewProvider())
		So(com.HasDescription(), ShouldBeFalse)
		So(R.ValueOf(com.Description("foo")).Pointer(), ShouldEqual, R.ValueOf(com).Pointer())
		So(com.HasDescription(), ShouldBeTrue)
		So(com.GetDescription(), ShouldEqual, "foo")
	})
}
