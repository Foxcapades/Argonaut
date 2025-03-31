package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func WrapOptions(o *argo.CommandOptions) opts.Options {
	return cOptions{o}
}

type cOptions struct{ o *argo.CommandOptions }

func (c cOptions) MaxDefaultFlagGroupSizeForMetaGroup() int {
	return c.o.MaxDefaultFlagGroupSizeForMetaGroup
}

func (c cOptions) MetaFlagGroupName() string {
	return c.o.MetaFlagGroupName
}

func (c cOptions) HelpTextMaxWidth() int {
	return c.o.HelpTextMaxWidth
}
