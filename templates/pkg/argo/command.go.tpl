{{ $vars := (types "Command" "CommandBuilder" "Command" false) -}}
package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {
	{{ template "CommandCommon" $vars }}
	{{- template "CommandEnd" $vars }}
}
