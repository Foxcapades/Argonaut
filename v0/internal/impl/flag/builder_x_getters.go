package flag

import "github.com/Foxcapades/Argonaut/v0/pkg/argo"

func (f *Builder) GetShort() byte {
	return f.ShortFlag
}

func (f *Builder) HasShort() bool {
	return f.IsShortSet
}

func (f *Builder) GetLong() string {
	return f.LongFlag
}

func (f *Builder) HasLong() bool {
	return f.IsLongSet
}

func (f *Builder) GetDescription() string {
	return f.DescriptionText.Description()
}

func (f *Builder) HasDescription() bool {
	return len(f.DescriptionText.DescriptionText) > 0
}

func (f *Builder) GetArg() argo.ArgumentBuilder {
	return f.ArgBuilder
}

func (f *Builder) HasArg() bool {
	return f.ArgBuilder != nil
}
