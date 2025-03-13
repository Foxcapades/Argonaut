package comp

import (
	"fmt"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

type filterMode uint8

const (
	filterAll filterMode = iota
	filterLong
	filterShort
)

func AutocompleteSimple(words []string, current int, command any) ([]string, error) {
	// ignore the command name
	words = words[1:]
	current--

	if current < 0 || current >= len(words) {
		return nil, fmt.Errorf("current word index %d is out of bounds for the input word count", current+1)
	}

	if com, ok := command.(argo.CommandTree); ok {
		return simpleAutocompleteTree(words, current, com)
	} else if com, ok := command.(argo.Command); ok {
		return simpleAutocompleteSingle(words, current, com)
	} else {
		return nil, fmt.Errorf("unrecognized command type %T", command)
	}
}

type AutocompleteLine struct {
	Words    []string
	HelpText string
}

func AutocompleteFull(words []string, current int, command any) ([]AutocompleteLine, error) {
	// ignore the command name
	words = words[1:]
	current--

}
