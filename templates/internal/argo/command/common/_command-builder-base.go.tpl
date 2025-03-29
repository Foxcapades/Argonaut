package common

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type CommandBuilderBase struct {
{{ define "CommandBuilderBaseProps" -}}
	disableHelp bool
	description string
	flagGroups  []argo.FlagGroupBuilder
	callback    argo.CommandCallback[{{ .OutputType }}]
{{- end }}
}

{{ define "CommandBuilderBaseFuncs" -}}
func (i *{{ .ImplType }}) WithDescription(desc string) {{ .BuilderType }} {
	i.description = desc
	return i
}

func (i *{{ .ImplType }}) HasDescription() bool {
	return len(i.description) > 0
}

func (i *{{ .ImplType }}) Description() string {
	return i.description
}

func (i *{{ .ImplType }}) WithHelpDisabled() {{ .BuilderType }} {
	i.disableHelp = true
	return i
}

func (i *{{ .ImplType }}) IsHelpDisabled() bool {
	return i.disableHelp
}

func (i *{{ .ImplType }}) WithFlagGroup(group argo.FlagGroupBuilder) {{ .BuilderType }} {
	i.flagGroups = append(i.flagGroups, group)
	return i
}

func (i *{{ .ImplType }}) WithFlagGroups(groups ...argo.FlagGroupBuilder) {{ .BuilderType }} {
	i.flagGroups = append(i.flagGroups, groups...)
	return i
}

func (i *{{ .ImplType }}) HasFlagGroups(includeDefault bool) bool {
	return (includeDefault && len(i.flagGroups) > 0) ||
		(i.hasDefaultFlagGroup() && len(i.flagGroups) > 1) ||
		len(i.flagGroups) > 0
}

func (i *{{ .ImplType }}) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if !includeDefault && i.hasDefaultFlagGroup() {
		return i.flagGroups[1:]
	}

	return i.flagGroups
}

func (i *{{ .ImplType }}) hasDefaultFlagGroup() bool {
	return len(i.flagGroups) > 0 && flag.IsDefaultGroup(i.flagGroups[0])
}

func (i *{{ .ImplType }}) WithFlag(fb argo.FlagBuilder) {{ .BuilderType }} {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlag(fb)
	return i
}

func (i *{{ .ImplType }}) WithFlags(flags ...argo.FlagBuilder) {{ .BuilderType }} {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlags(flags...)
	return i
}

func (i *{{ .ImplType }}) HasFlags() bool {
	for _, group := range i.flagGroups {
		if group != nil && group.HasFlags() {
			return true
		}
	}

	return false
}

func (i *{{ .ImplType }}) WithCallback(callback argo.CommandCallback[{{ .OutputType }}]) {{ .BuilderType }} {
	i.callback = callback
	return i
}

func (i *{{ .ImplType }}) HasCallback() bool {
	return i.callback != nil
}

func (i *{{ .ImplType }}) Callback() argo.CommandCallback[{{ .OutputType }}] {
	return i.callback
}
{{- end }}