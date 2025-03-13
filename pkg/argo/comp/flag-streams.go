package comp

import (
	"github.com/Foxcapades/Argonaut/pkg/argo"
	"iter"
)

func shortFlags(groups []argo.FlagGroup) iter.Seq[argo.Flag] {
	return func(yield func(argo.Flag) bool) {
		for _, group := range groups {
			for _, flag := range group.Flags() {
				if flag.HasShortForm() {
					yield(flag)
				}
			}
		}
	}
}

func longFlags(groups []argo.FlagGroup) iter.Seq[argo.Flag] {
	return func(yield func(argo.Flag) bool) {
		for _, group := range groups {
			for _, flag := range group.Flags() {
				if flag.HasLongForm() {
					yield(flag)
				}
			}
		}
	}
}

func allFlags(groups []argo.FlagGroup) iter.Seq[argo.Flag] {
	return func(yield func(argo.Flag) bool) {
		for _, group := range groups {
			for _, flag := range group.Flags() {
				yield(flag)
			}
		}
	}
}
