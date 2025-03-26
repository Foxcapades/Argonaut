package flag

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Flag struct {
	hits uint16

	short    byte
	required bool
	isHelp   bool

	arg argo.Argument

	long string
	desc string

	callback argo.FlagCallback
}

//

func (f *Flag) ShortForm() byte {
	return f.short
}

func (f *Flag) HasShortForm() bool {
	return f.short != 0
}

//

func (f *Flag) LongForm() string {
	return f.long
}

func (f *Flag) HasLongForm() bool {
	return len(f.long) > 0
}

//

func (f *Flag) Description() string {
	return f.desc
}

func (f *Flag) HasDescription() bool {
	return len(f.desc) > 0
}

//

func (f *Flag) Argument() argo.Argument {
	return f.arg
}

func (f *Flag) HasArgument() bool {
	return f.arg != nil
}

//

func (f *Flag) IsRequired() bool {
	return f.required
}

func (f *Flag) RequiresArgument() bool {
	return f.arg != nil && f.arg.IsRequired()
}

//

func (f *Flag) WasHit() bool {
	return f.hits > 0
}

func (f *Flag) HitCount() int {
	return int(f.hits)
}

//

func (f *Flag) HasCallback() bool {
	return f.callback != nil
}

func (f *Flag) Callback() argo.FlagCallback {
	return f.callback
}

//

func (f *Flag) IsHelpFlag() bool {
	return f.isHelp
}

//

func (f *Flag) String() string {
	return PrintFlagNames(f)
}

func (f *Flag) IncrementHitCount() {
	f.hits++
}
