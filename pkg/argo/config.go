package argo

type Config struct {
	Flags FlagConfig

	EndOfOptionsMarker string

	DefaultUnmarshaler MagicUnmarshaler
}
