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
			if f.IsRequired() {
				log.DebugLn2("flag %s is marked as required in flag group %s", func() { PrintFlagNames(f) }, group.Name)

				if !f.WasHit() {
					errs.AppendError(NewMissingFlagError(f))
				} else if f.Argument().IsRequired() && !f.Argument().WasHit() {
					errs.AppendError(fmt.Errorf("flag %s requires an argument", PrintFlagNames(f)))
				}
			} else if !f.WasHit() && f.HasArgument() && f.Argument().HasDefault() {
				log.DebugLn2("applying default value to flag %s's argument in flag group %s", func() { PrintFlagNames(f) }, group.Name)

				if err := f.Argument().SetToDefault(); err != nil {
					errs.AppendError(err)
				}
			}
		}
	}
}
