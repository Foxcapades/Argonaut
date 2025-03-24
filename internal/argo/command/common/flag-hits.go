package common

import (
	"fmt"
	"iter"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/log"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// TODO: iterate backwards!!!
func ExecuteHelpFlagCallbacks(flags iter.Seq[argo.Flag]) {
	for flag := range flags {
		if flag.IsHelpFlag() {
			if flag.HasCallback() {
				flag.Callback()(flag)
			}
		}
	}
}

func ExecuteFlagCallbacks(flags iter.Seq[argo.Flag]) {
	// Iterate through the rest of the flags and execute any callbacks.
	for flag := range flags {
		if !flag.IsHelpFlag() && flag.HasCallback() {
			flag.Callback()(flag)
		}
	}
}

func CheckRequiredArguments(arguments []argo.Argument, errs argo.MultiError) {
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

func CheckRequiredFlags(flagGroups []argo.FlagGroup, errs argo.MultiError) {
	for _, group := range flagGroups {
		log.DebugLn1("testing if required flags were hit for flag group %s", group.Name)

		for _, f := range group.Flags() {
			if f.IsRequired() {
				log.DebugLn2("flag %s is marked as required in flag group %s", func() { flag.PrintFlagNames(f) }, group.Name)

				if !f.WasHit() {
					errs.AppendError(fmt.Errorf("required flag %s was not used", flag.PrintFlagNames(f)))
				} else if f.Argument().IsRequired() && !f.Argument().WasHit() {
					errs.AppendError(fmt.Errorf("flag %s requires an argument", flag.PrintFlagNames(f)))
				}
			} else if !f.WasHit() && f.HasArgument() && f.Argument().HasDefault() {
				log.DebugLn2("applying default value to flag %s's argument in flag group %s", func() { flag.PrintFlagNames(f) }, group.Name)

				if err := f.Argument().SetToDefault(); err != nil {
					errs.AppendError(err)
				}
			}
		}
	}
}
