package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ParentCommon struct {
	common.CommandBase[{{ .OutputType }}]
{{ define "ParentCommonProps" -}}
	commandGroups []argo.CommandGroup
	selectedChild argo.ChildNode
	incompleteFn argo.IncompleteCommandHandler[{{ .OutputType }}]
{{- end }}
}

{{ define "ParentCommonFuncs" -}}
func (i *{{ .ImplType }}) HasCommandGroups(includeDefault bool) bool {
	return len(i.groups) > 1 || (includeDefault && i.hasDefaultCommandGroup())
}

func (i *{{ .ImplType }}) CommandGroups(includeDefault bool) []argo.CommandGroup {
	if !i.HasCommandGroups(includeDefault) {
		return nil
	}

	if includeDefault && i.hasDefaultCommandGroup() {
		return i.groups
	}

	return i.groups[1:]
}

func (i *{{ .ImplType }}) hasDefaultCommandGroup() bool {
	return i.groups[0] != nil
}

func (i *{{ .ImplType }}) FindChild(name string) argo.ChildNode {
	for _, group := range i.groups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

func (i *{{ .ImplType }}) HasSelectedChild() bool {
	return i.selectedChild != nil
}

func (i *{{ .ImplType }}) SelectChild(name string) bool {
	for _, group := range i.groups {
		if child := group.FindChild(name); child != nil {
			i.selectedChild = child
			return true
		}
	}

	return false
}

func (i *{{ .ImplType }}) SelectedChild() argo.ChildNode {
	return i.selectedChild
}

func (i *{{ .ImplType }}) HasIncompleteHandler() bool {
	return i.incompleteFn != nil
}

func (i *{{ .ImplType }}) IncompleteHandler() argo.IncompleteCommandHandler[argo.CommandTree] {
	return i.incompleteFn
}
{{- end }}
