package group

import (
	"fmt"
	"strings"

	"github.com/foxcapades/argonaut/internal/util"
	"github.com/foxcapades/argonaut/internal/util/xerr"
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Builder struct {
	name  string
	desc  string
	flags []argo.FlagSpecBuilder
}

func (b *Builder) WithName(name string) argo.FlagGroupSpecBuilder {
	b.name = strings.TrimSpace(name)
	return b
}

func (b *Builder) WithDescription(description string) argo.FlagGroupSpecBuilder {
	b.desc = description
	return b
}

func (b *Builder) WithFlag(flag argo.FlagSpecBuilder) argo.FlagGroupSpecBuilder {
	b.flags = append(b.flags, util.RequireNonNil(flag))
	return b
}

func (b *Builder) Build(config argo.Config) (argo.FlagGroupSpec, error) {
	errs := xerr.NewMultiError()

	out := new(Spec)

	if len(b.name) == 0 {
		errs.AppendMsg("flag group name was empty")
	} else {
		out.name = b.name
	}

	out.flags = make([]argo.FlagSpec, len(b.flags))

	shorts := make(map[byte]bool, len(b.flags))
	longs := make(map[string]bool, len(b.flags))

	for i := range b.flags {
		if flag, err := b.flags[i].Build(config); err == nil {
			out.flags[i] = flag

			if flag.HasLongForm() {
				if longs[flag.LongForm()] {
					errs.AppendMsg(fmt.Sprintf(
						"flag group \"%s\" has more than one flag with the long form \"%s%s\"",
						out.name,
						config.Flags.LongFormPrefix,
						flag.LongForm(),
					))
				} else {
					longs[flag.LongForm()] = true
				}
			}

			if flag.HasShortForm() {
				if shorts[flag.ShortForm()] {
					errs.AppendMsg(fmt.Sprintf(
						"flag group \"%s\" has more than one flag with the short form \"%c%c\"",
						out.name,
						config.Flags.ShortFormPrefix,
						flag.ShortForm(),
					))
				} else {
					shorts[flag.ShortForm()] = true
				}
			}
		} else {
			errs.AppendIfNotNil(err)
		}
	}

	out.description = b.desc

	if errs.IsEmpty() {
		return out, nil
	}

	return nil, errs
}
