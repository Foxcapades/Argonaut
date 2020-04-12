package cli

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/marsh"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/props"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Type Passthroughs                                                   ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// NewCommand returns a new instance of a CommandBuilder
// from the current provider.
//
// Same as calling `Provider().NewCommand()`
func NewCommand() argo.CommandBuilder {
	return Provider().NewCommand()
}

// NewFlag returns a new instance of a FlagBuilder from the
// current provider.
//
// Same as calling `Provider().NewFlag()`
func NewFlag() argo.FlagBuilder {
	return Provider().NewFlag()
}

// NewFlagGroup returns a new instance of a FlagGroupBuilder
// from the current provider
//
// Same as calling `Provider().NewFlagGroup()`
func NewFlagGroup() argo.FlagGroupBuilder {
	return Provider().NewFlagGroup()
}

// NewArg returns a new instance of an ArgumentBuilder from
// the current provider.
//
// Same as calling `Provider().NewArg()`
func NewArg() argo.ArgumentBuilder {
	return Provider().NewArg()
}

//┏━━━━━━━━━━━━━━━━━━━━━`━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Value Unmarsha`ler Passthroughs                                      ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// DefaultUnmarshalProps returns a defaulted unmarshaling
// configuration.
func DefaultUnmarshalProps() argo.UnmarshalProps {
	return props.DefaultUnmarshalProps()
}

// UnmarshalDefault attempts to unmarshal the given string
// into the given pointer using the default unmarshaling
// configuration.
func UnmarshalDefault(raw string, ptr interface{}) (err error) {
	return marsh.NewDefaultedValueUnmarshaler().Unmarshal(raw, ptr)
}

// Unmarshal attempts to unmarshal the given string into the
// given pointer using the provided unmarshaling
// configuration.
func Unmarshal(raw string, ptr interface{}, props argo.UnmarshalProps) (err error) {
	return marsh.NewValueUnmarshaler(&props).Unmarshal(raw, ptr)
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Provider Passthroughs                                               ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// Provider returns an instance of the default Provider
// implementation.
//
// This method is provided to allow the use of custom CLI
// element implementations.
func Provider() argo.Provider {
	return impl.GetProvider()
}

// SetProvider allows the override of the provider type.
//
// Intended to allow completely changing the behavior of
// of the way CLI elements are constructed.
func SetProvider(pro argo.Provider) {
	impl.SetProvider(pro)
}
