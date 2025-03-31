{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.ImplementationFields */ -}}
package tree

import (
  "fmt"

  "github.com/foxcapades/argonaut/v3/internal/argo/command/common"
  "github.com/foxcapades/argonaut/v3/internal/argo/flag"
  "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ChildCommon struct {
  common.CommandBase[{{ .OutputType }}]
{{ define "ChildCommonProps" -}}
  name    string
  parent  argo.ParentNode
  aliases []string
{{- end }}
}

{{ define "ChildCommonFuncs" -}}
func (i *{{ .ImplType }}) Name() string {
  return i.name
}

func (i *{{ .ImplType }}) Parent() argo.ParentNode {
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

func (i *{{ .ImplType }}) FindShortFlagRecursive(c byte) argo.Flag {
  if f := i.FindShortFlag(c); f != nil {
    return f
  }

  switch t := i.parent.(type) {
  case argo.TreeCommand:
    return t.FindShortFlag(c)
  case argo.BranchCommand:
    return t.FindShortFlagRecursive(c)
  default:
    panic(fmt.Sprintf("illegal state: unknown node parent type: %v", i.parent))
  }
}

func (i *{{ .ImplType }}) FindLongFlagRecursive(name string) argo.Flag {
  if f := i.FindLongFlag(name); f != nil {
    return f
  }

  switch t := i.parent.(type) {
  case argo.TreeCommand:
    return t.FindLongFlag(name)
  case argo.BranchCommand:
    return t.FindLongFlagRecursive(name)
  default:
    panic(fmt.Sprintf("illegal state: unknown node parent type: %v", i.parent))
  }
}

{{- end }}