{{ $vars := (map "ImplType" "Branch" "OutputType" "argo.BranchCommand") -}}
package tree

import (
  "github.com/foxcapades/argonaut/v3/internal/argo/command/common"
  "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Branch struct {
  common.CommandBase[argo.BranchCommand]
  {{ template "tree-parent-impl-props" $vars }}
  {{ template "tree-child-impl-props" $vars }}
}

{{ template "tree-parent-impl-funcs" $vars }}

{{ template "tree-child-impl-funcs" $vars }}
