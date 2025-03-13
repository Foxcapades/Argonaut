package comp

import (
	"github.com/Foxcapades/Argonaut/pkg/argo"
	"slices"
	"strings"
)

// TODO: don't autocomplete if the previous flag requires an argument!
// TODO: add an option to indicate whether a flag is to be used more than once.
// TODO: add the ability to configure a function to provide autocomplete options
//       for flags and positional arguments.

func simpleAutocompleteSingle(words []string, current int, command argo.Command) ([]string, error) {
	cword := words[current]
	mode := filterAll

	// If they have given us more than one character of input
	if len(cword) > 1 {
		if cword[0] != '-' {
			return nil, nil
		}

		// If the second character of the input string is another prefix char, then
		// we are looking for long flags.
		if cword[1] == '-' {
			mode = filterLong
		} else

		// If the second character of the input string is _not_ another prefix
		// character then we are looking for a combinable short flag
		{
			return simpleAutoCompleteSingleShorts(words, current, cword, command)
		}
	} else

	// If they have given us only one character of input
	if len(cword) == 1 {
		if cword[0] != '-' {
			return nil, nil
		}
	}

	for _, group := range command.FlagGroups() {
		group.Flags()
	}
}

func simpleAutoCompleteSingleShorts(words []string, pos int, cword string, command argo.Command) ([]string, error) {
	cword = cword[1:]

	last := byte(0)

	if len(cword) == 0 {
		if pos > 0 {
			p := pos-1
			if len(words[p]) > 0 {
				last = words[p][len(words[p])-1]
			}
		}
	} else {
		last = cword[len(cword)-1]
	}

	suggestFirst := make([]string, 0, 8)
	suggestLast := make([]string, 0, 2)

	// map out all the possible short flags
	for flag := range shortFlags(command.FlagGroups()) {
		// If the last flag in the group requires an argument, then it would be
		// invalid to suggest another flag to follow it.
		if flag.ShortForm() == last && flag.RequiresArgument() {
			return nil, nil
		}

		if strings.IndexByte(cword, flag.ShortForm()) > -1 {
			suggestLast = append(suggestLast, string(flag.ShortForm()))
		} else {
			suggestFirst = append(suggestFirst, string(flag.ShortForm()))
		}
	}

	slices.Sort(suggestFirst)
	slices.Sort(suggestLast)

	return append(suggestFirst, suggestLast...), nil
}

func simpleAutoCompleteSingleLongs(words []string, pos int, cword string, command argo.Command) ([]string, error) {
	cword = cword[2:]

	suggestFirst := make([]string, 0, 8)
	suggestLast := make([]string, 0, 2)

	for flag := range longFlags(command.FlagGroups()) {
		if ()

		// If the last flag in the group requires an argument, then it would be
		// invalid to suggest another flag to follow it.
		if flag.ShortForm() == last && flag.RequiresArgument() {
			return nil, nil
		}

		if strings.IndexByte(cword, flag.ShortForm()) > -1 {
			suggestLast = append(suggestLast, string(flag.ShortForm()))
		} else {
			suggestFirst = append(suggestFirst, string(flag.ShortForm()))
		}
	}

	slices.Sort(suggestFirst)
	slices.Sort(suggestLast)

	return append(suggestFirst, suggestLast...), nil
}
