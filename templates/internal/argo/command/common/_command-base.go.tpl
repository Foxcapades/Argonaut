package common

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type CommandBase struct {
{{ define "CommandBaseProps" -}}
  disableHelp bool
  description string
  flagGroups  []argo.FlagGroup
  callback    argo.CommandCallback[{{ .OutputType }}]
{{- end }}
}

{{ define "CommandBaseFuncs" -}}
func (i *{{ .ImplType }}) HasDescription() bool {
  return len(c.description) > 0
}

func (i *{{ .ImplType }}) Description() string {
  return c.description
}

func (i *{{ .ImplType }}) HasFlagGroups(includeDefault bool) bool {

}

func (i *{{ .ImplType }}) FlagGroups(includeDefault bool) []argo.FlagGroup {

}

func (i *{{ .ImplType }}) HasCallback() bool {
  return c.callback != nil
}

func (i *{{ .ImplType }}) Callback() argo.CommandCallback[{{ .OutputType }}] {
  return c.callback
}

func (i *{{ .ImplType }}) FindShortFlag(b byte) argo.Flag {
  for _, group := range c.flagGroups {
    if flag := group.FindShortFlag(b); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *{{ .ImplType }}) FindLongFlag(name string) argo.Flag {
  for _, group := range c.flagGroups {
    if flag := group.FindLongFlag(name); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *{{ .ImplType }}) IsHelpDisabled() bool {
  return c.disableHelp
}
{{- end }}
