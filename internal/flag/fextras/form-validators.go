package fextras

import "github.com/foxcapades/argonaut/pkg/argo"

func LongFormValidator(name string, config argo.Config) bool {
	for _, c := range name {
		if c == rune(config.Flags.LongFormValueSeparator) {
			return false
		}

		if c >= '-' && c <= ':' {
			continue
		}

		if c >= '@' && c <= 'Z' {
			continue
		}

		if c >= 'a' && c <= 'z' {
			continue
		}

		switch c {
		case '%', '+', '=', '^', '_':
			continue
		}

		return false
	}

	return true
}

func ShortFormValidator(name byte, config argo.Config) bool {
	if name == config.Flags.ShortFormValueSeparator {
		return false
	}

	// avoid issues with ambiguous '--'
	switch string([]byte{config.Flags.ShortFormPrefix, name}) {
	case config.EndOfOptionsMarker:
		return false
	case config.Flags.LongFormPrefix:
		return false
	}

	if name >= '-' && name <= ':' {
		return true
	}

	if name >= '@' && name <= 'Z' {
		return true
	}

	if name >= 'a' && name <= 'z' {
		return true
	}

	switch name {
	case '%', '+', '=', '^', '_':
		return true
	}

	return false
}
