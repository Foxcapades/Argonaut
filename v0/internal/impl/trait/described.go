package trait

type Described struct {
	DescriptionText string
}

func (d *Described) Description() string {
	return d.DescriptionText
}

func (d *Described) HasDescription() bool {
	return len(d.DescriptionText) > 0
}
