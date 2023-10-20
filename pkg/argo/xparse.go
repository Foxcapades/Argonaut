package argo

import (
	"fmt"
	"strconv"
	"strings"
)

func parseBool(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "true", "t", "yes", "y", "on", "1":
		return true, nil
	case "false", "f", "no", "n", "off", "0":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %s as bool", raw)
	}
}

func parseInt(v string, bits int, opt *UnmarshalIntegerProps) (int64, error) {
	var neg string
	// TODO: Wrap this error

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

func parseUInt(v string, bits int, opt *UnmarshalIntegerProps) (uint64, error) {
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

const errMapEntry = "cannot parse \"%s\" as a map entry. " +
	"Must contain a valid key value separator (one of %s)"

func parseMapEntry(s string, props *UnmarshalMapProps) (k string, v string, err error) {
	i := strings.IndexAny(s, props.KeyValSeparatorChars)
	if i < 0 {
		err = fmt.Errorf(errMapEntry, s, props.KeyValSeparatorChars)
	} else {
		k = s[:i]
		v = s[i+1:]
	}
	return
}
