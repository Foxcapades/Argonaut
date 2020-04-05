package cli

import (
	"github.com/Foxcapades/Argonaut/v1/internal/impl"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Command                                                             ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

var cProvider = CommandProvider(impl.NewCommandBuilder)

type CommandProvider func() argo.CommandBuilder

func SetCommandProvider(pro CommandProvider) {
	if pro == nil {
		panic("cannot set a nil command provider")
	}
	cProvider = pro
}

func NewCommand() argo.CommandBuilder {
	return cProvider()
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Flag                                                                ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

var fProvider = FlagProvider(impl.NewFlagBuilder)

type FlagProvider func() argo.FlagBuilder

func SetFlagProvider(pro FlagProvider) {
	if pro == nil {
		panic("cannot set a nil flag provider")
	}
	fProvider = pro
}

func NewFlag() argo.FlagBuilder {
	return fProvider()
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Flag Group                                                          ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

var gProvider = FlagGroupProvider(impl.NewFlagGroupBuilder)

type FlagGroupProvider func() argo.FlagGroupBuilder

func SetFlagGroupProvider(pro FlagGroupProvider) {
	if pro == nil {
		panic("cannot set a nil flag group provider")
	}
	gProvider = pro
}

func NewFlagGroup() argo.FlagGroupBuilder {
	return gProvider()
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Arguments                                                           ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

var aProvider = ArgumentProvider(impl.NewArgBuilder)

type ArgumentProvider func() argo.ArgumentBuilder

func SetArgumentProvider(pro ArgumentProvider) {
	if pro == nil {
		panic("cannot set a nil argument provider")
	}
	aProvider = pro
}

func NewArg() argo.ArgumentBuilder {
	return aProvider()
}

func DefaultUnmarshalProps() argo.UnmarshalProps {
	return impl.DefaultUnmarshalProps()
}

func UnmarshalDefault(raw string, val interface{}) (err error) {
	return impl.UnmarshalDefault(raw, val)
}

func Unmarshal(raw string, val interface{}, props argo.UnmarshalProps) (err error) {
	return impl.Unmarshal(raw, val, &props)
}
