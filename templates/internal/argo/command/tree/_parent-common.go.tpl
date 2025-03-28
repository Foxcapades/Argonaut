{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.ImplementationFields */ -}}
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
func (i *{{ .ImplType }}) HasCommandGroups() bool {
	return len(i.commandGroups) > 0
}

func (i *{{ .ImplType }}) CommandGroups() []argo.CommandGroup {
	return i.commandGroups
}

func (i *{{ .ImplType }}) hasDefaultCommandGroup() bool {
	return i.commandGroups[0] != nil
}

func (i *{{ .ImplType }}) HasSubcommands() bool {
	for _, group := range i.commandGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *{{ .ImplType }}) FindChild(name string) argo.ChildNode {
	for _, group := range i.commandGroups {
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
	for _, group := range i.commandGroups {
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

func (i *{{ .ImplType }}) IncompleteHandler() argo.IncompleteCommandHandler[{{ .OutputType }}] {
	return i.incompleteFn
}
{{- end }}
