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
  return len(i.description) > 0
}

func (i *{{ .ImplType }}) Description() string {
  return i.description
}

func (i *{{ .ImplType }}) HasFlagGroups() bool {
  return len(i.flagGroups) > 0
}

func (i *{{ .ImplType }}) FlagGroups() []argo.FlagGroup {
  return i.flagGroups
}

func (i *{{ .ImplType }}) HasCallback() bool {
  return i.callback != nil
}

func (i *{{ .ImplType }}) Callback() argo.CommandCallback[{{ .OutputType }}] {
  return i.callback
}

func (i *{{ .ImplType }}) FindShortFlag(b byte) argo.Flag {
  for _, group := range i.flagGroups {
    if flag := group.FindShortFlag(b); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *{{ .ImplType }}) FindLongFlag(name string) argo.Flag {
  for _, group := range i.flagGroups {
    if flag := group.FindLongFlag(name); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *{{ .ImplType }}) IsHelpDisabled() bool {
  return i.disableHelp
}
{{- end }}
