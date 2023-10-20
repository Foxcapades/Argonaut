package argo_test

import "errors"

type FailingWriter struct {
	FailAfter int
	current   int
}

func (f *FailingWriter) Write(p []byte) (n int, err error) {
	if f.current < f.FailAfter {
		f.current++
		return len(p), nil
	} else {
		return 0, errors.New("fake error")
	}
}
