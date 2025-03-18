package parse

import (
	"fmt"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
	"strconv"
	"strings"
)

func Bool(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "true", "t", "yes", "y", "on", "1":
		return true, nil
	case "false", "f", "no", "n", "off", "0":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %s as bool", raw)
	}
}

func Int64(v string, bits int, opt *argo.UnmarshalIntegerProps) (int64, error) {
	var neg string

	if v[0] == '-' {
		neg = "-"
		v = v[1:]
	}

	for i := range opt.HexLeaders {
		if strings.HasPrefix(v, opt.HexLeaders[i]) {
			return strconv.ParseInt(neg+v[len(opt.HexLeaders[i]):], 16, bits)
		}
	}

	for i := range opt.OctalLeaders {
		if strings.HasPrefix(v, opt.OctalLeaders[i]) {
			return strconv.ParseInt(neg+v[len(opt.OctalLeaders[i]):], 8, bits)
		}
	}

	return strconv.ParseInt(neg+v, 10, bits)
}

func UInt64(v string, bits int, opt *argo.UnmarshalIntegerProps) (uint64, error) {
	// TODO: Wrap this error

	for i := range opt.HexLeaders {
		if strings.HasPrefix(v, opt.HexLeaders[i]) {
			return strconv.ParseUint(v[len(opt.HexLeaders[i]):], 16, bits)
		}
	}

	for i := range opt.OctalLeaders {
		if strings.HasPrefix(v, opt.OctalLeaders[i]) {
			return strconv.ParseUint(v[len(opt.OctalLeaders[i]):], 8, bits)
		}
	}

	return strconv.ParseUint(v, 10, bits)
}
