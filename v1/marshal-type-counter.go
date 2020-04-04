package argo

type UseCounter int

func (u *UseCounter) Unmarshal(string) error {
	*u++
	return nil
}

