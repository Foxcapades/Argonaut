package cli

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Provider Passthroughs                                               ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func Provider() argo.Provider {
	return impl.GetProvider()
}

func SetProvider(pro argo.Provider) {
	impl.SetProvider(pro)
}


//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Type Passthroughs                                                   ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func NewCommand() argo.CommandBuilder {
	return Provider().NewCommand()
}

func NewFlag() argo.FlagBuilder {
	return Provider().NewFlag()
}

func NewFlagGroup() argo.FlagGroupBuilder {
	return Provider().NewFlagGroup()
}

func NewArg() argo.ArgumentBuilder {
	return Provider().NewArg()
}


//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Value Unmarshaler Passthroughs                                      ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//


func DefaultUnmarshalProps() argo.UnmarshalProps {
	return impl.DefaultUnmarshalProps()
}

func UnmarshalDefault(raw string, val interface{}) (err error) {
	return impl.NewDefaultedValueUnmarshaler().Unmarshal(raw, val)
}

func Unmarshal(raw string, val interface{}, props argo.UnmarshalProps) (err error) {
	return impl.NewValueUnmarshaler(&props).Unmarshal(raw, val)
}
