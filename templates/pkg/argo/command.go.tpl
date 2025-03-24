{{ $vars := (types "Command" "CommandBuilder" "Command" false) -}}
package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

{{ if false -}}
import (
	. "github.com/foxcapades/argonaut/v3/pkg/argo"
)
{{- end -}}

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {
	{{ template "CommandCommon" $vars }}
	{{- template "CommandEnd" $vars }}
}

// A CommandBuilder provides an API to configure the construction of a new
// Command instance.
//
// Example Usage:
//
//	cli.Command().
//	    WithDescription("This is my command that does something.").
//	    WithFlag(cli.Flag().
//	        WithShortForm('v').
//	        WithLongForm("verbose").
//	        WithDescription("Enable verbose logging.")
//	        WithBinding(&config.verbose)).
//	    WithArgument(cli.Argument().
//	        WithName("file").
//	        WithDescription("File path.").
//	        WithBinding(&config.file)).
type CommandBuilder interface {
  {{ template "CommandBuilderCommon" $vars }}
  {{ template "CommandEndBuilder" $vars }}
}