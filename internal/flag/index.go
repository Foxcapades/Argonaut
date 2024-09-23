package flag

type Index struct {
	byLongName  map[string]int
	byShortName map[byte]int
	flags       []Flag
}

func (i Index) LongFlag(name string) (Flag, bool) {
	if idx, ok := i.byLongName[name]; ok {
		return i.flags[idx], true
	}

	return Flag{}, false
}

func (i Index) ShortFlag(name byte) (Flag, bool) {
	if idx, ok := i.byShortName[name]; ok {
		return i.flags[idx], true
	}

	return Flag{}, false
}
