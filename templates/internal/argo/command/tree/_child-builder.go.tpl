package tree

type ChildCommandBuilder struct {
{{ define "ChildCommandBuilderProps" -}}
	name    string
	aliases []string
	parent  argo.ParentNodeBuilder
{{- end }}
}

{{ define "ChildCommandBuilderFuncs" -}}
func (i *{{ .ImplType }}) Name() string {
	return i.name
}

func (i *{{ .ImplType }}) WithAlias(alias string) {{ .BuilderType }} {
	i.aliases = append(i.aliases, alias)
	return i
}

func (i *{{ .ImplType }}) WithAliases(aliases ...string) {{ .BuilderType }} {
	i.aliases = append(i.aliases, aliases...)
	return i
}

func (i *{{ .ImplType }}) Aliases() []string {
	return i.aliases
}

func (i *{{ .ImplType }}) HasAliases() bool {
	return len(i.aliases) > 0
}

func (i *{{ .ImplType }}) SetParentNode(parent argo.ParentNodeBuilder) {{ .BuilderType }} {
	i.parent = parent
	return i
}

func (i *{{ .ImplType }}) ParentNode() argo.ParentNodeBuilder {
	return i.parent
}
{{- end }}
