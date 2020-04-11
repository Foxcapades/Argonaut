package trait

type Described struct {
	DescTxt string
}

func (d *Described) Description() string {
	return d.DescTxt
}

func (d *Described) HasDescription() bool {
	return len(d.DescTxt) > 0
}
