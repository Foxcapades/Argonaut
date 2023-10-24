package argo

import (
	"errors"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

// A FlagBuilder is used to construct a Flag instance which represents the input
// from the CLI call.
type FlagBuilder interface {

	// WithShortForm sets the short-form flag character that the flag may be
	// referenced by on the CLI.
	//
	// Short-form flags consist of a single character preceded by either a
	// prefix/leader character, or by one or more other short-form flags.
	//
	// Short-form flags must be alphanumeric.
	//
	// Examples:
	//     # Single, unchained short flags.
	//     -f
	//     -f bar
	//     -f=bar
	//     # Multiple short flags chained.  In these examples, the short flag '-c'
	//     # takes an optional string argument, which will be "def" in the last
	//     # two examples.
	//     -abc
	//     -abc def
	//     -abc=def
	WithShortForm(char byte) FlagBuilder

	hasShortForm() bool
	getShortForm() byte

	// WithLongForm sets the long-form flag name that the flag may be referenced
	// by on the CLI.
	//
	// Long-form flags consist of one or more characters preceded immediately by
	// two prefix/leader characters (typically dashes).
	//
	// Long-form flags must start with an alphanumeric character and may only
	// consist of alphanumeric characters, dashes, and/or underscores.
	//
	// Example long-form flags:
	//     # The '--foo' flag takes an optional string argument
	//     --foo
	//     --foo bar
	//     --foo=bar
	WithLongForm(form string) FlagBuilder

	hasLongForm() bool
	getLongForm() string

	// WithDescription sets an optional description value for the Flag being
	// built.
	//
	// The description value is used for rendering help text.
	WithDescription(desc string) FlagBuilder

	// WithCallback provides a function that will be called when a Flag is hit
	// while parsing the CLI inputs.
	//
	// The given function will be called after parsing has completed, regardless
	// of whether there were parsing errors.
	//
	// Flag on-hit callbacks will be executed in priority order with the higher
	// priority values executing before lower priority values.  For flags that
	// have the same priority, the callbacks will be called in the order the flags
	// appeared in the CLI call.
	WithCallback(fn FlagCallback) FlagBuilder

	// WithArgument attaches the given argument to the Flag being built.
	//
	// Only one argument may be set on a Flag at a time.
	WithArgument(arg ArgumentBuilder) FlagBuilder

	// WithBinding is a shortcut method for attaching an argument and binding it
	// to the given pointer.
	//
	// Bind is equivalent to calling one of the following:
	//    WithArgument(cli.Argument().Bind(ptr))
	//    // or
	//    WithArgument(cli.Argument().Bind(ptr).Require())
	WithBinding(pointer any, required bool) FlagBuilder

	// WithBindingAndDefault is a shortcut method for attaching an argument,
	// binding it to the given pointer, and setting a default on that argument.
	//
	// BindWithDefault is equivalent to calling one of the following:
	//     WithArgument(cli.Argument().WithBinding(ptr).WithDefault(something))
	//     // or
	//     WithArgument(cli.Argument().WithBinding(ptr).WithDefault(something).Require())
	WithBindingAndDefault(pointer, def any, required bool) FlagBuilder

	setIsHelpFlag() FlagBuilder

	// Require marks this Flag as being required.
	//
	// If this flag is not present in the CLI call, an error will be returned when
	// parsing the CLI input.
	Require() FlagBuilder

	// Build builds a new Flag instance constructed from the components set on
	// this FlagBuilder.
	Build(warnings *WarningContext) (Flag, error)
}

// NewFlagBuilder returns a new FlagBuilder instance.
func NewFlagBuilder() FlagBuilder {
	return &flagBuilder{}
}

type flagBuilder struct {
	short  byte
	req    bool
	isHelp bool
	long   string
	desc   string
	onHit  FlagCallback
	arg    ArgumentBuilder
}

func (b *flagBuilder) WithShortForm(char byte) FlagBuilder {
	b.short = char
	return b
}

func (b flagBuilder) hasShortForm() bool {
	return b.short != 0
}

func (b flagBuilder) getShortForm() byte {
	return b.short
}

func (b *flagBuilder) WithLongForm(form string) FlagBuilder {
	b.long = form
	return b
}

func (b flagBuilder) hasLongForm() bool {
	return len(b.long) > 0
}

func (b flagBuilder) getLongForm() string {
	return b.long
}

func (b *flagBuilder) WithDescription(desc string) FlagBuilder {
	b.desc = desc
	return b
}

func (b *flagBuilder) WithCallback(fn func(Flag)) FlagBuilder {
	b.onHit = fn
	return b
}

func (b *flagBuilder) WithArgument(arg ArgumentBuilder) FlagBuilder {
	b.arg = arg
	return b
}

func (b *flagBuilder) Require() FlagBuilder {
	b.req = true
	return b
}

func (b *flagBuilder) WithBinding(pointer any, required bool) FlagBuilder {
	b.arg = NewArgumentBuilder().WithBinding(pointer)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *flagBuilder) WithBindingAndDefault(pointer, def any, required bool) FlagBuilder {
	b.arg = NewArgumentBuilder().WithBinding(pointer).WithDefault(def)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *flagBuilder) setIsHelpFlag() FlagBuilder {
	b.isHelp = true
	return b
}

func (b *flagBuilder) Build(ctx *WarningContext) (Flag, error) {
	errs := newMultiError()

	if b.short > 0 {
		if err := validateShortForm(b.short); err != nil {
			errs.AppendError(err)
		}
	}

	if len(b.long) > 0 {
		if err := validateLongForm(b.long); err != nil {
			errs.AppendError(err)
		}
	}

	if !b.hasShortForm() && !b.hasLongForm() {
		errs.AppendError(errors.New("flag declared with neither a long or short form"))
	}

	var arg Argument

	if b.arg != nil {
		var err error
		arg, err = b.arg.Build(ctx)
		if err != nil {
			errs.AppendError(err)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &flag{
		warnings: ctx,
		short:    b.short,
		required: b.req,
		arg:      arg,
		long:     b.long,
		desc:     b.desc,
		isHelp:   b.isHelp,
		callback: b.onHit,
	}, nil
}

func validateShortForm(c byte) error {
	if !chars.IsAlphanumeric(c) {
		return errors.New("short-form flags must be alphanumeric")
	}

	return nil
}

func validateLongForm(f string) error {
	if !chars.IsAlphanumeric(f[0]) {
		return errors.New("long-form flags must begin with an alphanumeric character")
	}

	for i := 1; i < len(f); i++ {
		if !chars.IsFlagStringSafe(f[i]) {
			return errors.New("long-form flags must only contain alphanumeric characters, dashes, and/or underscores")
		}
	}

	return nil
}
