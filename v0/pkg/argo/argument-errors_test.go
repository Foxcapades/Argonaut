package argo_test

import (
	. "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	. "github.com/smartystreets/goconvey/convey"
	. "testing"
)

func TestInvalidArgError_Type(t *T) {
	Convey("InvalidArgError.Type", t, func() {
		tests := []*struct {
			name string
			err  error
			val  ArgumentErrorType
		}{
			{"Invalid Binding", NewInvalidArgError(ArgErrInvalidBinding, nil, ""), ArgErrInvalidBinding},
			{"Invalid Default", NewInvalidArgError(ArgErrInvalidDefault, nil, ""), ArgErrInvalidDefault},
			{"Invalid Default Function", NewInvalidArgError(ArgErrInvalidDefaultFn, nil, ""), ArgErrInvalidDefaultFn},
			{"Invalid Default Value", NewInvalidArgError(ArgErrInvalidDefaultVal, nil, ""), ArgErrInvalidDefaultVal},
			{"Invalid Binding Type", NewInvalidArgError(ArgErrInvalidBindingBadType, nil, ""), ArgErrInvalidBindingBadType},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				iae, ok := test.err.(*InvalidArgError)
				So(ok, ShouldBeTrue)
				So(iae.Type(), ShouldEqual, test.val)
			})
		}
	})
}

func TestInvalidArgError_Is(t *T) {
	Convey("InvalidArgError.Type", t, func() {
		tests := []*struct {
			name string
			err  error
			val  map[ArgumentErrorType]bool
		}{
			{
				"Invalid Binding",
				NewInvalidArgError(ArgErrInvalidBinding, nil, ""),
				map[ArgumentErrorType]bool{
					ArgErrInvalidDefault:        false,
					ArgErrInvalidBinding:        true,
					ArgErrInvalidDefaultFn:      false,
					ArgErrInvalidDefaultVal:     false,
					ArgErrInvalidBindingBadType: false,
				},
			},
			{
				"Invalid Default",
				NewInvalidArgError(ArgErrInvalidDefault, nil, ""),
				map[ArgumentErrorType]bool{
					ArgErrInvalidDefault:        true,
					ArgErrInvalidBinding:        false,
					ArgErrInvalidDefaultFn:      false,
					ArgErrInvalidDefaultVal:     false,
					ArgErrInvalidBindingBadType: false,
				},
			},
			{
				"Invalid Default Function",
				NewInvalidArgError(ArgErrInvalidDefaultFn, nil, ""),
				map[ArgumentErrorType]bool{
					ArgErrInvalidDefault:        true,
					ArgErrInvalidBinding:        false,
					ArgErrInvalidDefaultFn:      true,
					ArgErrInvalidDefaultVal:     false,
					ArgErrInvalidBindingBadType: false,
				},
			},
			{
				"Invalid Default Value",
				NewInvalidArgError(ArgErrInvalidDefaultVal, nil, ""),
				map[ArgumentErrorType]bool{
					ArgErrInvalidDefault:        true,
					ArgErrInvalidBinding:        false,
					ArgErrInvalidDefaultFn:      false,
					ArgErrInvalidDefaultVal:     true,
					ArgErrInvalidBindingBadType: false,
				},
			},
			{
				"Invalid Binding Type",
				NewInvalidArgError(ArgErrInvalidBindingBadType, nil, ""),
				map[ArgumentErrorType]bool{
					ArgErrInvalidDefault:        false,
					ArgErrInvalidBinding:        true,
					ArgErrInvalidDefaultFn:      false,
					ArgErrInvalidDefaultVal:     false,
					ArgErrInvalidBindingBadType: true,
				},
			},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				iae, ok := test.err.(*InvalidArgError)
				So(ok, ShouldBeTrue)
				for kind, expec := range test.val {
					Convey(kind.String(), func() {
						So(iae.Is(kind), ShouldEqual, expec)
					})
				}
			})
		}
	})
}
