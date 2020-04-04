package cli

import (
	"github.com/Foxcapades/Argonaut/v1/internal/impl"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func NewCommand() argo.CommandBuilder {
	return impl.NewCommandBuilder()
}

func NewFlag() argo.FlagBuilder {
	return impl.NewFlagBuilder()
}

func NewArg() argo.ArgumentBuilder {
	return impl.NewArgBuilder()
}

func DefaultUnmarshalProps() argo.UnmarshalProps {
	return impl.DefaultUnmarshalProps()
}

func UnmarshalDefault(raw string, val interface{}) (err error) {
	return impl.UnmarshalDefault(raw, val)
}

func Unmarshal(raw string, val interface{}, props argo.UnmarshalProps) (err error) {
	return impl.Unmarshal(raw, val, props)
}
