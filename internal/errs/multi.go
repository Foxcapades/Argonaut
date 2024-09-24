package errs

import (
	"fmt"
	"strings"

	"github.com/foxcapades/argonaut/pkg/argo"
)

type MultiError interface {
	argo.MultiError

	Append(e error)
	AppendIfNotNil(e error)
}

func NewMultiError() MultiError {
	return new(multiError)
}

type multiError struct {
	errors []error
}

func (m *multiError) Error() string {
	switch len(m.errors) {
	case 0:
		return "unknown error"
	case 1:
		return m.errors[0].Error()
	}

	unique := make(map[string]bool, len(m.errors))
	messages := make([]string, 0, len(m.errors))

	for i := range m.errors {
		msg := m.errors[i].Error()

		if ok := unique[msg]; !ok {
			unique[msg] = true
			messages = append(messages, msg)
		}
	}

	sb := strings.Builder{}

	sb.WriteString(argo.ErrMsgMultiErrorHeaderLine(len(messages), len(m.errors)))
	sb.WriteString(fmt.Sprintln())

	for i := range messages {
		sb.WriteString("    ")
		sb.WriteString(messages[i])
		sb.WriteString(fmt.Sprintln())
	}

	return sb.String()
}

func (m *multiError) Size() int {
	return len(m.errors)
}

func (m *multiError) Get(idx int) error {
	return m.errors[idx]
}

func (m *multiError) Append(e error) {
	m.errors = append(m.errors, e)
}

func (m *multiError) AppendIfNotNil(e error) {
	if e != nil {
		m.errors = append(m.errors, e)
	}
}

func (m *multiError) IsEmpty() bool {
	return len(m.errors) == 0
}
