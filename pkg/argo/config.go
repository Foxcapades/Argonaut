package argo

type Config struct {
	ShortFlagPrefix byte

	ShortFlagValueSeparator byte

	LongFlagPrefix string

	LongFlagValueSeparator byte

	EndOfOptionsMarker string

	DefaultUnmarshaler MagicUnmarshaler
}
