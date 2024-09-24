package argo

type MagicUnmarshaler interface {
	Unmarshal(raw string, into any) error
}
