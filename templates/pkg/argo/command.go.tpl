{{ $vars := (types "Command" "CommandBuilder" "Command" false) -}}
package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {
	{{ template "CommandCommon" $vars }}
	{{- template "CommandEnd" $vars }}

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists, this method will return nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a target Flag instance by its long-form name.
	//
	// If no such flag exists, this method will return nil.
	FindLongFlag(name string) Flag

	Options() CommandOptions
}
