package xerr

import "github.com/foxcapades/argonaut/v3/pkg/argo"

func AppendIfPresent(err error, errs argo.MultiError) {
	if err != nil {
		errs.AppendError(err)
	}
}

func TryBuild[B, O any](builder B, buildFunc func(B) (O, error), errs argo.MultiError) (O, bool) {
	if out, err := buildFunc(builder); err != nil {
		errs.AppendError(err)
		return out, false
	} else {
		return out, true
	}
}
