package tree

import (
  "github.com/foxcapades/argonaut/v3/internal/argo/command/common"
  "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type child struct {
  common.CommandBase[{{ .OutputType }}]
  {{ define "tree-child-impl-props" -}}
  name    string
  parent  argo.ParentNode[any]
  aliases []string
  {{- end }}
}

{{ define "tree-child-impl-funcs" -}}
func (i *{{ .ImplType }}) Name() string {
  return i.name
}

func (i *{{ .ImplType }}) Parent() argo.ParentNode[any] {
  return i.parent
}

func (i *{{ .ImplType }}) HasAliases() bool {
  return len(i.aliases) > 0
}

func (i *{{ .ImplType }}) Aliases() []string {
  return i.aliases
}

func (i *{{ .ImplType }}) Matches(name string) bool {
  if i.name == name {
    return true
  }

  for _, alias := range i.aliases {
    if alias == name {
      return true
    }
  }

  return false
}
{{- end }}