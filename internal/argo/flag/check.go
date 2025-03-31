package flag

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/log"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func CheckRequired(flagGroups []argo.FlagGroup, errs argo.MultiError) {
	for _, group := range flagGroups {
		log.DebugLn1("testing if required flags were hit for flag group %s", group.Name)

		for _, f := range group.Flags() {
			if f.WasHit() {
				if f.HasArgument() && !f.Argument().WasHit() {
					if f.Argument().HasDefault() {
						if err := f.Argument().SetToDefault(); err != nil {
							errs.AppendError(err)
						}
					} else if f.Argument().IsRequired() {
						errs.AppendError(fmt.Errorf("argument for flag %s is required", PrintFlagNames(f)))
					}
				}
			} else if f.IsRequired() {
				errs.AppendError(NewMissingFlagError(f))
			} else if f.HasArgument() && f.Argument().HasDefault() {
				if err := f.Argument().SetToDefault(); err != nil {
					errs.AppendError(err)
				}
			}
		}
	}
}
