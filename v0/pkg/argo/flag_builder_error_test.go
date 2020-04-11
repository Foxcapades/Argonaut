package argo_test

import (
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFlagBuilderError_Type(t *testing.T) {
	errors := []argo.FlagBuilderErrorType{
		argo.FlagBuilderErrNoFlags,
		argo.FlagBuilderErrBadShortFlag,
		argo.FlagBuilderErrBadLongFlag,
	}

	convey.Convey("argo.FlagBuilderError.Type()", t, func() {
		for _, val := range errors {
			convey.Convey(val.String(), func() {
				convey.So(argo.NewFlagBuilderError(val, nil).(argo.FlagBuilderError).Type(),
					convey.ShouldEqual, val)
			})
		}
	})
}
