{{ $vars := (map "ImplType" "Leaf") -}}
package tree

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Leaf struct {
  {{ template "CommandBaseProps" $vars }}
  {{ template "ChildCommonProps" $vars }}
  unmappedLabel string
  args          []argo.Argument
  unmapped      []string
}

func (i *Leaf) HasArguments() bool {
  return len(i.args) > 0
}

func (i *Leaf) Arguments() []argo.Argument {
  return i.args
}

func (i *Leaf) HasUnmappedInputs() bool {
  return len(i.unmapped) > 0
}

func (i *Leaf) UnmappedInputs() []string {
  return i.unmapped
}

func (i *Leaf) AppendUnmappedInput(input string) {
  i.unmapped = append(i.unmapped, input)
}

func (i *Leaf) HasUnmappedLabel() bool {
  return len(i.unmappedLabel) > 0
}

func (i *Leaf) GetUnmappedLabel() string {
  return i.unmappedLabel
}

{{ template "ChildCommonFuncs" $vars }}

{{ template "CommandBaseFuncs" $vars }}
