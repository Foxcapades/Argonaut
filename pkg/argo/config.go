package argo

type Config struct {
	Flags FlagConfig

	EndOfOptionsMarker string

	Unmarshalling UnmarshallingConfig

	DefaultUnmarshaler MagicUnmarshaler
}
