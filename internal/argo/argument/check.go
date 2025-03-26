package argument

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/log"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func CheckRequired(arguments []argo.Argument, errs argo.MultiError) {
	for i, arg := range arguments {
		if arg.IsRequired() {
			log.DebugLn2("testing if required argument %d (<%s>) has received a value", i+1, arg.Name())

			if !arg.WasHit() {
				if arg.HasName() {
					errs.AppendError(fmt.Errorf("argument %d (<%s>) is required", i+1, arg.Name()))
				} else {
					errs.AppendError(fmt.Errorf("argument %d is required", i+1))
				}
			}
		} else if !arg.WasHit() && arg.HasDefault() {
			log.DebugLn2("applying default value to optional argument %d (<%s>)", i+1, arg.Name())

			if err := arg.SetToDefault(); err != nil {
				errs.AppendError(err)
			}
		}
	}
}
