package entropy

import (
	"fmt"
)

var (
	ErrEntropyExhausted = fmt.Errorf("fatal error: read past entropy end")
)

type Entropy struct {
	data []byte
}

func New(data []byte) *Entropy {
	return &Entropy{data: data}
}

func (e *Entropy) Read(p []byte) (n int, err error) {
	if len(e.data) == 0 {
		return 0, ErrEntropyExhausted
	}
	n = len(p)
	if n > len(e.data) {
		n = len(e.data)
	}
	n = copy(p, e.data[:n])
	e.data = e.data[n:]
	return n, nil
}
