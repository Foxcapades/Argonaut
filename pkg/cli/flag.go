package cli

import (
	"github.com/foxcapades/argonaut/internal/flag"
	"github.com/foxcapades/argonaut/pkg/argo"
)

func Flag(short byte, long string) argo.FlagSpecBuilder {
	return flag.NewBuilder().
		WithShortForm(short).
		WithLongForm(long)
}

func ShortFlag(name byte) argo.FlagSpecBuilder {
	return flag.NewBuilder().WithShortForm(name)
}

func LongFlag(name string) argo.FlagSpecBuilder {
	return flag.NewBuilder().WithLongForm(name)
}

func FlagGroup(name string) argo.FlagGroupSpecBuilder {}
