package impl

import "github.com/Foxcapades/Argonaut/v0/pkg/argo"

func DefaultUnmarshalProps() argo.UnmarshalProps {
	return defaultUnmarshalProps
}

var defaultUnmarshalProps = argo.UnmarshalProps{
	Integers: argo.UnmarshalIntegerProps{
		OctalLeaders: []string{"0o", "0O", "o", "O"},
		HexLeaders:   []string{"0x", "0X", "x", "X"},
		DefaultBase:  10,
	},
	Maps: argo.UnmarshalMapProps{
		KeyValSeparatorChars: "=:",
		EntrySeparatorChars:  ",; ",
	},
	Slices: argo.UnmarshalSliceProps{},
}
