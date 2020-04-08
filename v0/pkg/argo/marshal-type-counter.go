package argo

type UseCounter int

func (u *UseCounter) ConsumesArguments() bool {
	return false
}

func (u *UseCounter) Unmarshal(string) error {
	*u++
	return nil
}
