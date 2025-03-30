package opts

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
	"github.com/foxcapades/argonaut/v3/pkg/argoutil"
)

func FixOptions(opts *argo.Options) {
	consoleWidth, _ := argoutil.GetConsoleWidth()
	if opts.HelpTextMaxWidth < 60 {
		opts.HelpTextMaxWidth = max(min(consoleWidth, 120), 60)
	}
	if len(opts.MetaFlagGroupName) == 0 {
		opts.MetaFlagGroupName = "General Flags"
	}
	if len(opts.DefaultCommandGroupName) == 0 {
		opts.DefaultCommandGroupName = "Commands"
	}
}
