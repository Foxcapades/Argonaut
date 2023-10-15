package flag

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/internal/impl/argument/argerr"
	"github.com/Foxcapades/Argonaut/v1/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

// implements argo.Flag
type flag struct {
	hits uint16

	short    byte
	required bool

	arg argo.Argument

	long string
	desc string
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

func (f flag) Argument() argo.Argument {
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

func (f flag) Hit() error {
	f.hits++
	if f.arg.IsRequired() {
		return fmt.Errorf("flag %s requires an input", flagutil.PrintFlagNames(&f))
	}
	return nil
}

func (f flag) HitWithArg(rawArg string) error {
	f.hits++

	if f.arg != nil {
		return f.arg.SetValue(rawArg)
	} else {
		return argerr.NewUnexpectedArgumentError(rawArg, f)
	}
}

func (f flag) WasHit() bool {
	return f.hits > 0
}

func (f flag) HitCount() int {
	return int(f.hits)
}
