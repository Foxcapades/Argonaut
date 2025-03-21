//go:build !dev

package log

func DebugLn1[T any](_ string, _ T) {}

func DebugLn2[A, B any](_ string, _ A, _ B) {}
