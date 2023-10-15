package flagutil

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func UniqueFlagNames(groups []argo.FlagGroupBuilder, errs argo.MultiError) {
	longs := make(map[string]uint8, len(groups))
	shorts := make(map[byte]uint8, len(groups))

	for _, group := range groups {
		for _, flag := range group.GetFlags() {
			if flag.HasLongForm() {
				longs[flag.GetLongForm()]++
				if longs[flag.GetLongForm()] == 2 {
					errs.AppendError(fmt.Errorf("conflicting flag longform name %s", flag.GetLongForm()))
				}
			}

			if flag.HasShortForm() {
				shorts[flag.GetShortForm()]++
				if shorts[flag.GetShortForm()] == 2 {
					errs.AppendError(fmt.Errorf("conflicting flag shortform character %c", flag.GetShortForm()))
				}
			}
		}
	}
}
