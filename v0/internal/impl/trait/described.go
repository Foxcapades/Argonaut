package trait

type Described struct {
	DescriptionValue string
}

func (d *Described) Description() string {
	return d.DescriptionValue
}

func (d *Described) HasDescription() bool {
	return len(d.DescriptionValue) > 0
}
