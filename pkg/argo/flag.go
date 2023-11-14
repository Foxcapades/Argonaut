package argo

import "fmt"

// Flag represents a single CLI flag which may have an argument.
type Flag interface {

	// ShortForm returns the short-form character representing this Flag.
	ShortForm() byte

	// HasShortForm indicates whether this Flag has a short-form character.
	HasShortForm() bool

	// LongForm returns the long-form string representing this Flag.
	LongForm() string

	// HasLongForm indicates whether this Flag has a long-form string.
	HasLongForm() bool

	// Description returns the help-text description of this flag.
	Description() string

	// HasDescription indicates whether this Flag has a help-text description.
	HasDescription() bool

	// Argument returns the argument value attached to this Flag.
	//
	// If this Flag does not have an Argument attached, this method will return
	// nil.
	Argument() Argument

	// HasArgument indicates whether this Flag accepts an argument.
	HasArgument() bool

	// IsRequired indicates whether this Flag is required.
	//
	// Required flags must be present in the CLI call.
	IsRequired() bool

	// RequiresArgument indicates whether this flag has an argument that is
	// required.
	RequiresArgument() bool

	// WasHit indicates whether this flag was used in the CLI call.
	WasHit() bool

	// HitCount returns the number of times that this flag was used in the CLI
	// call.
	HitCount() int

	AppendWarning(warning string)

	isHelpFlag() bool
	hit() error
	hitWithArg(rawArg string) error
	executeCallback()
}

// A FlagCallback is a function that, if set on a flag, will be called by the
// CLI parsing process if that flag is used in the CLI call.
//
// The flag callback will be called after CLI parsing has completed.
type FlagCallback = func(flag Flag)

type flag struct {
	hits uint16

	short    byte
	required bool
	isHelp   bool

	arg Argument

	long string
	desc string

	callback FlagCallback
	warnings *WarningContext
}

func (f flag) ShortForm() byte {
	return f.short
}

func (f flag) HasShortForm() bool {
	return f.short != 0
}

func (f flag) LongForm() string {
	return f.long
}

func (f flag) HasLongForm() bool {
	return len(f.long) > 0
}

func (f flag) Description() string {
	return f.desc
}

func (f flag) HasDescription() bool {
	return len(f.desc) > 0
}

func (f flag) Argument() Argument {
	return f.arg
}

func (f flag) HasArgument() bool {
	return f.arg != nil
}

func (f flag) IsRequired() bool {
	return f.required
}

func (f flag) RequiresArgument() bool {
	return f.arg != nil && f.arg.IsRequired()
}

func (f flag) isHelpFlag() bool {
	return f.isHelp
}

func (f *flag) hit() error {
	f.hits++
	if f.HasArgument() && f.arg.IsRequired() {
		return fmt.Errorf("flag %s requires an input", printFlagNames(f))
	}

	if hasBooleanArgument(f) {
		return f.arg.setValue("true")
	}

	return nil
}

func (f *flag) hitWithArg(rawArg string) error {
	f.hits++

	if f.arg != nil {
		return f.arg.setValue(rawArg)
	} else {
		return nil // TODO: warning for this
	}
}

func (f flag) WasHit() bool {
	return f.hits > 0
}

func (f flag) HitCount() int {
	return int(f.hits)
}

func (f flag) AppendWarning(warning string) {
	f.warnings.appendWarning(warning)
}

func (f flag) executeCallback() {
	if f.callback != nil {
		f.callback(&f)
	}
}
