package opts

type Options interface {
	MaxDefaultFlagGroupSizeForMetaGroup() int
	MetaFlagGroupName() string
	HelpTextMaxWidth() int
}
