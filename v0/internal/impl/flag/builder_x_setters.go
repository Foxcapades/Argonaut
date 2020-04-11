package flag

import "github.com/Foxcapades/Argonaut/v0/pkg/argo"

func (f *Builder) Default(val interface{}) argo.FlagBuilder {
	if f.ArgBuilder == nil {
		f.ArgBuilder = f.Provider.NewArg()
	}
	f.ArgBuilder.Default(val)
	return f
}

func (f *Builder) Bind(ptr interface{}, required bool) argo.FlagBuilder {
	if f.ArgBuilder == nil {
		f.ArgBuilder = f.Provider.NewArg()
	}

	f.ArgBuilder.Bind(ptr).Required(required)

	return f
}

func (f *Builder) Short(flag byte) argo.FlagBuilder {
	f.IsShortSet = true
	f.ShortFlag = flag
	return f
}

func (f *Builder) OnHit(fn argo.FlagEventHandler) argo.FlagBuilder {
	f.OnHitCallback = fn
	return f
}

func (f *Builder) Long(flag string) argo.FlagBuilder {
	f.IsLongSet = true
	f.LongFlag = flag
	return f
}

func (f *Builder) Description(desc string) argo.FlagBuilder {
	f.DescriptionText.DescriptionValue = desc
	return f
}

func (f *Builder) Arg(arg argo.ArgumentBuilder) argo.FlagBuilder {
	f.ArgBuilder = arg
	return f
}

func (f *Builder) Parent(fg argo.FlagGroup) argo.FlagBuilder {
	f.ParentElement = fg
	return f
}

func (f *Builder) BindUseCount(ptr *int) argo.FlagBuilder {
	f.UseCountBinding = ptr
	return f
}
