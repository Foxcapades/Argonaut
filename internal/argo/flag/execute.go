package flag

import (
	"iter"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func ExecuteHelpFlagCallbacks(flags iter.Seq[argo.Flag]) {
	for f := range flags {
		if f.IsHelpFlag() {
			if f.HasCallback() {
				f.Callback()(f)
			}
		}
	}
}

func ExecuteFlagCallbacks(flags iter.Seq[argo.Flag]) {
	// Iterate through the rest of the flags and execute any callbacks.
	for f := range flags {
		if !f.IsHelpFlag() && f.HasCallback() {
			f.Callback()(f)
		}
	}
}
