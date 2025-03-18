package flag

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type flag struct {
	hits uint16

	short    byte
	required bool
	isHelp   bool

	arg argo.Argument

	long string
	desc string

	callback argo.FlagCallback
	warnings *argo.WarningContext
}

//

func (f *flag) ShortForm() byte {
	return f.short
}

func (f *flag) HasShortForm() bool {
	return f.short != 0
}

//

func (f *flag) LongForm() string {
	return f.long
}

func (f *flag) HasLongForm() bool {
	return len(f.long) > 0
}

//

func (f *flag) Description() string {
	return f.desc
}

func (f *flag) HasDescription() bool {
	return len(f.desc) > 0
}

//

func (f *flag) Argument() argo.Argument {
	return f.arg
}

func (f *flag) HasArgument() bool {
	return f.arg != nil
}

//

func (f *flag) IsRequired() bool {
	return f.required
}

func (f *flag) RequiresArgument() bool {
	return f.arg != nil && f.arg.IsRequired()
}

//

func (f *flag) WasHit() bool {
	return f.hits > 0
}

func (f *flag) HitCount() int {
	return int(f.hits)
}

//

func (f *flag) HasCallback() bool {
	return f.callback != nil
}

func (f *flag) Callback() argo.FlagCallback {
	return f.callback
}

//

func (f *flag) IsHelpFlag() bool {
	return f.isHelp
}

//

func (f *flag) AppendWarning(warning string) {
	f.warnings.AppendWarning(warning)
}

//

func (f *flag) String() string {
	return PrintFlagNames(f)
}

func (f *flag) hit() error {
	f.hits++
	if f.HasArgument() && f.arg.IsRequired() {
		return fmt.Errorf("flag %s requires an input", PrintFlagNames(f))
	}

	if HasBooleanArgument(f) {
		return f.arg.SetValue("true")
	}

	return nil
}

func (f *flag) hitWithArg(rawArg string) error {
	f.hits++

	if f.arg != nil {
		return f.arg.SetValue(rawArg)
	} else {
		return nil // TODO: warning for this
	}
}
