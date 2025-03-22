{{ $vars := (map "ImplType" "Branch" "OutputType" "argo.BranchCommand") -}}
package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Branch struct {
  {{ template "CommandBaseProps" $vars }}
  {{ template "ChildCommonProps" $vars }}
  {{ template "ParentCommonProps" $vars }}
}

{{ template "CommandBaseFuncs" $vars }}

{{ template "ChildCommonFuncs" $vars }}

{{ template "ParentCommonFuncs" $vars }}
