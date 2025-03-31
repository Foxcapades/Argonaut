package flag

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Index struct {
	byShort map[byte]argo.Flag
	byLong  map[string]argo.Flag
}

func (i Index) ByShortName(f byte) argo.Flag {
	return i.byShort[f]
}

func (i Index) ByLongName(f string) argo.Flag {
	return i.byLong[f]
}

func (i Index) Overlay(c GroupContainer) {
	for _, f := range Stream(c) {
		if f.HasShortForm() {
			i.byShort[f.ShortForm()] = f
		}

		if f.HasLongForm() {
			i.byLong[f.LongForm()] = f
		}
	}
}

func BuildFlagIndex(c GroupContainer) Index {
	out := Index{
		byShort: make(map[byte]argo.Flag, 8),
		byLong:  make(map[string]argo.Flag, 8),
	}

	out.Overlay(c)

	return out
}
