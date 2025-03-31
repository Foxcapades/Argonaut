package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Options interface {
	opts.Options

	DefaultCommandGroupName() string

	InheritParentFlags() argo.ParentFlagInheritanceMode
}

func WrapOptions(o *argo.TreeCommandOptions) Options {
	return tOptions{o}
}

type tOptions struct{ o *argo.TreeCommandOptions }

func (t tOptions) MaxDefaultFlagGroupSizeForMetaGroup() int {
	return t.o.MaxDefaultFlagGroupSizeForMetaGroup
}

func (t tOptions) MetaFlagGroupName() string {
	return t.o.MetaFlagGroupName
}

func (t tOptions) HelpTextMaxWidth() int {
	return t.o.HelpTextMaxWidth
}

func (t tOptions) DefaultCommandGroupName() string {
	return t.o.DefaultCommandGroupName
}

func (t tOptions) InheritParentFlags() argo.ParentFlagInheritanceMode {
	return t.o.InheritParentFlags
}
