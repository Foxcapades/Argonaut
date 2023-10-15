package flagutil

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func PrintFlagNames(flag argo.Flag) string {
	if flag.HasLongForm() {
		if flag.HasShortForm() {
			return fmt.Sprintf("-%c | --%s", flag.ShortForm(), flag.LongForm())
		}

		return fmt.Sprintf("--%s", flag.LongForm())
	}

	return fmt.Sprintf("-%c", flag.ShortForm())
}
