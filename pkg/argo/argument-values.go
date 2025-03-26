package argo

type ArgumentBindingType uint8

const (
	BindingTypeNone ArgumentBindingType = iota
	BindingTypePointer
	BindingTypeUnmarshaler
	BindingTypeSimpleFunc
	BindingTypeErrorFunc

	BindingTypeUnknown ArgumentBindingType = 254
	BindingTypeInvalid ArgumentBindingType = 255
)

type ArgumentBinding interface {
	Type() ArgumentBindingType

	BoundTo() any
}

type ArgumentDefaultType uint8

const (
	DefaultTypeNone ArgumentDefaultType = iota
	DefaultTypeRaw
	DefaultTypeParsed
	DefaultTypeProviderPlain
	DefaultTypeProviderWithErr

	DefaultTypeUnknown ArgumentDefaultType = 254
	DefaultTypeInvalid ArgumentDefaultType = 255
)

type ArgumentDefault interface {
	Type() ArgumentDefaultType

	Value() any
}

type ArgumentBindingError interface {
	error

	ArgumentBuilder() ArgumentBuilder

	Unwrap() error
}

type MissingRequiredArgumentError interface {
	error
	Argument() Argument
	Flag() Flag
	HasFlag() bool
}
