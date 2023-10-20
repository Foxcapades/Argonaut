package argo

// A FlagGroup is a named group or category for a collection of one or more
// Flag instances.
//
// Flag groups are primarily used to categorize CLI flags when rendering help
// text.
//
// Flag groups that do not contain any flags will be filtered out at build time.
type FlagGroup interface {

	// Name of the FlagGroup.
	//
	// This value will be used when rendering help text.
	Name() string

	// Description for this FlagGroup.
	//
	// This optional value will be used when rendering help text.
	Description() string

	// HasDescription indicates whether this FlagGroup has a description attached.
	HasDescription() bool

	// Flags returns the Flag instances contained by this FlagGroup.
	Flags() []Flag

	// FindShortFlag checks this FlagGroup's contents for a flag with the given
	// short flag character.  If one could not be found, this method returns nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag checks this FlagGroup's contents for a flag with the given
	// long flag name.  If one could not be found, this method returns nil.
	FindLongFlag(name string) Flag

	size() int
}

type flagGroup struct {
	warnings *WarningContext
	name     string
	desc     string
	flags    []Flag
}

func (f flagGroup) Name() string {
	return f.name
}

func (f flagGroup) Description() string {
	return f.desc
}

func (f flagGroup) HasDescription() bool {
	return len(f.desc) > 0
}

func (f flagGroup) Flags() []Flag {
	return f.flags
}

func (f flagGroup) FindShortFlag(c byte) Flag {
	for _, flag := range f.flags {
		if flag.HasShortForm() && flag.ShortForm() == c {
			return flag
		}
	}

	return nil
}

func (f flagGroup) FindLongFlag(name string) Flag {
	for _, flag := range f.flags {
		if flag.HasLongForm() && flag.LongForm() == name {
			return flag
		}
	}

	return nil
}

func (f flagGroup) size() int {
	return len(f.flags)
}
