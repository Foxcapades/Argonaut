package impl

import "github.com/Foxcapades/Argonaut/v1/pkg/argo"

type flagProps = uint8

const (
	flagHasShort flagProps = 1 << iota
	flagHasLong
	flagIsReq
)

func NewFlag() *Flag {
	return new(Flag)
}

type Flag struct {
	props flagProps
	short byte
	arg   argo.Argument
	long  string
	desc  string
}

func (f Flag) Short() byte {
	panic("implement me")
}

func (f Flag) HasShort() bool {
	panic("implement me")
}

func (f Flag) Long() string {
	panic("implement me")
}

func (f Flag) HasLong() bool {
	panic("implement me")
}

func (f Flag) Required() bool {
	panic("implement me")
}

func (f Flag) Argument() argo.Argument {
	panic("implement me")
}

func (f Flag) HasArgument() bool {
	panic("implement me")
}

func (f Flag) Description() string {
	panic("implement me")
}

func (f Flag) HasDescription() bool {
	panic("implement me")
}

func (f Flag) Hits() int {
	panic("implement me")
}

func (f Flag) hit() {
	panic("implement me")
}
