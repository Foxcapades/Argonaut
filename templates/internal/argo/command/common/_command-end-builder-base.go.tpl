{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.ImplementationFields */ -}}
package common

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type CommandEndBuilderBase struct {
{{ define "CommandEndBuilderBaseProps" -}}
	arguments     []argo.Argument
	unmapped      []string
	unmappedLabel string
{{- end }}
}

{{ define "CommandEndBuilderBaseFuncs" -}}
func (i *{{ .ImplType }}) HasArguments() bool {
	return len(i.arguments) > 0
}

func (i *{{ .ImplType }}) Arguments() []argo.Argument {
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

func (i *{{ .ImplType }}) HasUnmappedInputLabel() bool {
	return len(i.unmappedLabel) > 0
}

func (i *{{ .ImplType }}) UnmappedInputLabel() string {
	return i.unmappedLabel
}
{{- end }}