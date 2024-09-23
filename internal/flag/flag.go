package flag

import (
	"github.com/foxcapades/argonaut/internal/arg"
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Flag struct {
	hasArgument bool
	isRequired  bool
	hits        uint16

	argument arg.ArgumentSpec

	activeCallbacks []argo.FlagCallback // how to ref argo without reffing argo here?
	lazyCallbacks   []argo.FlagCallback
}

// region argo.Flag API

func (f Flag) IsRequired() bool { return f.isRequired }

func (f Flag) HasArgument() bool { return f.hasArgument }

func (f Flag) WasUsed() bool { return f.hits > 0 }

func (f Flag) UsageCount() int { return int(f.hits) }

func (f Flag) Argument() argo.Argument {
	// Don't return the argument unless it was explicitly defined by the user.
	//
	// We always have an argument, if not one defined by the user, then a default
	// boolean toggle arg.
	if f.hasArgument {
		return f.argument
	}

	return nil
}

// endregion argo.Flag API

// region internal API

// IArgument returns the internal form of the argument contained by this Flag
// instance.
func (f Flag) IArgument() arg.ArgumentSpec {
	return f.argument
}

// Hit marks this Flag instance as having been used in the CLI call.
//
// This method should be called _after_ any argument value has been parsed.
func (f *Flag) Hit() {
	f.hits++

	for i := range f.activeCallbacks {
		f.activeCallbacks[i](f)
	}
}

// Finalize performs any final actions on the flag after the CLI input parsing
// has completed.
func (f *Flag) Finalize() {
	for i := range f.lazyCallbacks {
		f.lazyCallbacks[i](f)
	}
}

// endregion internal API
