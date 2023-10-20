package argo

import (
	"fmt"
)

func printFlagNames(flag Flag) string {
	if flag.HasLongForm() {
		if flag.HasShortForm() {
			return fmt.Sprintf("-%c | --%s", flag.ShortForm(), flag.LongForm())
		}

		return fmt.Sprintf("--%s", flag.LongForm())
	}

	return fmt.Sprintf("-%c", flag.ShortForm())
}

func uniqueFlagNames(groups []FlagGroupBuilder, errs MultiError) {
	longs := make(map[string]uint8, len(groups))
	shorts := make(map[byte]uint8, len(groups))

	for _, group := range groups {
		for _, flag := range group.getFlags() {
			if flag.hasLongForm() {
				longs[flag.getLongForm()]++
				if longs[flag.getLongForm()] == 2 {
					errs.AppendError(fmt.Errorf("conflicting flag longform name %s", flag.getLongForm()))
				}
			}

			if flag.hasShortForm() {
				shorts[flag.getShortForm()]++
				if shorts[flag.getShortForm()] == 2 {
					errs.AppendError(fmt.Errorf("conflicting flag shortform character %c", flag.getShortForm()))
				}
			}
		}
	}
}
