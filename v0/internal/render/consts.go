package render

import "github.com/Foxcapades/Argonaut/v0/pkg/argo"

type (
	ac  = argo.Command
	aa  = argo.Argument
	af  = argo.Flag
	afg = argo.FlagGroup
)

const (
	dblLineBreak = "\n\n"
	sngLineBreak = '\n'
	sngSpace     = ' '
	maxWidth     = 100
)
