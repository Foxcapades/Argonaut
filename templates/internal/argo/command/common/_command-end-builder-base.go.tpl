{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.ImplementationFields */ -}}
package common

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type CommandEndBuilderBase struct {
{{ define "CommandEndBuilderBaseProps" -}}
	arguments     []argo.ArgumentBuilder
	unmapped      []string
	unmappedLabel string
{{- end }}
}

{{ define "CommandEndBuilderBaseFuncs" -}}
func (i *{{ .ImplType }}) WithArgument(argument argo.ArgumentBuilder) {{ .BuilderType }} {
	i.arguments = append(i.arguments, argument)
	return i
}

func (i *{{ .ImplType }}) WithArguments(arguments ...argo.ArgumentBuilder) {{ .BuilderType }} {
	i.arguments = append(i.arguments, arguments...)
	return i
}

func (i *{{ .ImplType }}) HasArguments() bool {
	return len(i.arguments) > 0
}

func (i *{{ .ImplType }}) Arguments() []argo.ArgumentBuilder {
	return i.arguments
}

func (i *{{ .ImplType }}) HasUnmappedInputs() bool {
	return len(i.unmapped) > 0
}

func (i *{{ .ImplType }}) UnmappedInputs() []string {
	return i.unmapped
}

func (i *{{ .ImplType }}) AppendUnmappedInput(input string) {
	i.unmapped = append(i.unmapped, input)
}

func (i *{{ .ImplType }}) WithUnmappedInputLabel(label string) {{ .BuilderType }} {
	i.unmappedLabel = label
	return i
}

func (i *{{ .ImplType }}) HasUnmappedInputLabel() bool {
	return len(i.unmappedLabel) > 0
}

func (i *{{ .ImplType }}) UnmappedInputLabel() string {
	return i.unmappedLabel
}
{{- end }}