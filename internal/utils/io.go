package utils

import (
	"bufio"
	"io"
)

func NewBatchWriter(w io.Writer) *BatchWriter {
	if b, ok := w.(*bufio.Writer); ok {
		return &BatchWriter{Buffer: b}
	} else {
		return &BatchWriter{Buffer: bufio.NewWriter(w)}
	}
}

type BatchWriter struct {
	Buffer *bufio.Writer
	Error  error
}

func (w *BatchWriter) WriteByte(b byte) {
	if w.Error == nil {
		w.Error = w.Buffer.WriteByte(b)
	}
}

func (w *BatchWriter) WriteString(s string) {
	if w.Error == nil {
		_, w.Error = w.Buffer.WriteString(s)
	}
}

func (w *BatchWriter) Flush() {
	err := w.Buffer.Flush()
	if w.Error == nil {
		w.Error = err
	}
}

func (w *BatchWriter) HasError() bool {
	return w.Error != nil
}
