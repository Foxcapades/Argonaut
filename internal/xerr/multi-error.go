package xerr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func NewMultiError() argo.MultiError {
	return multiError{
		errs:    make([]error, 0, 10),
		strings: make(map[string]bool, 10),
	}
}

type multiError struct {
	errs    []error
	strings map[string]bool
}

func (m multiError) Error() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("encountered %d unique errors:", len(m.strings)))

	for k := range m.strings {
		sb.WriteString("\n  ")
		sb.WriteString(k)
	}

	return sb.String()
}

func (m multiError) Errors() []error {
	return m.errs
}

func (m multiError) AppendError(err error) {
	var e argo.MultiError
	if errors.As(err, &e) {
		for _, err := range e.Errors() {
			m.AppendError(err)
		}
	} else {
		m.errs = append(m.errs, err)
		m.strings[err.Error()] = true
	}
}
