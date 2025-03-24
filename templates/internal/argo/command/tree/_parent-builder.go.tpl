package tree

type ParentCommandBuilder struct {
{{ define "ParentCommandBuilderProps" -}}
	comGroups    []argo.CommandGroupBuilder
	incompleteFn argo.IncompleteCommandHandler[argo.BranchCommand]
{{- end }}
}

{{ define "ParentCommandBuilderFuncs" -}}
func (i *{{ .ImplType }}) WithCommandGroup(group argo.CommandGroupBuilder) {{ .BuilderType }} {
	i.comGroups = append(i.comGroups, group)
	return i
}

func (i *{{ .ImplType }}) WithCommandGroups(groups ...argo.CommandGroupBuilder) {{ .BuilderType }} {
	i.comGroups = append(i.comGroups, groups...)
	return i
}

func (i *{{ .ImplType }}) HasCommandGroups(includeDefault bool) bool {
	return (includeDefault && len(i.comGroups) > 0) ||
		(i.hasDefaultCommandGroup() && len(i.comGroups) > 1) ||
		len(i.comGroups) > 0
}

func (i *{{ .ImplType }}) CommandGroups(includeDefault bool) []argo.CommandGroupBuilder {
	if !includeDefault && i.hasDefaultCommandGroup() {
		return i.comGroups[1:]
	}

	return i.comGroups
}

func (i *{{ .ImplType }}) hasDefaultCommandGroup() bool {
	return len(i.comGroups) > 0 && i.comGroups[0].Name() == DefaultCommandGroupName
}

func (i *{{ .ImplType }}) WithBranch(branch argo.BranchCommandBuilder) {{ .BuilderType }} {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranch(branch)
	return i
}

func (i *{{ .ImplType }}) WithBranches(branches ...argo.BranchCommandBuilder) {{ .BuilderType }} {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranches(branches...)
	return i
}

func (i *{{ .ImplType }}) WithLeaf(leaf argo.LeafCommandBuilder) {{ .BuilderType }} {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaf(leaf)
	return i
}

func (i *{{ .ImplType }}) WithLeaves(leaves ...argo.LeafCommandBuilder) {{ .BuilderType }} {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaves(leaves...)
	return i
}

func (i *{{ .ImplType }}) HasSubcommands() bool {
	for _, group := range i.comGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *{{ .ImplType }}) WithIncompleteHandler(handler argo.IncompleteCommandHandler[argo.BranchCommand]) {{ .BuilderType }} {
	i.incompleteFn = handler
	return i
}

func (i *{{ .ImplType }}) HasIncompleteHandler() bool {
	return i.incompleteFn != nil
}

func (i *{{ .ImplType }}) IncompleteHandler() argo.IncompleteCommandHandler[argo.BranchCommand] {
	return i.incompleteFn
}
{{- end }}
