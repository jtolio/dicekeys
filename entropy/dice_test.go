package entropy

import (
	"fmt"
	"testing"
)

func must(e *Entropy, err error) *Entropy {
	if err != nil {
		panic(err)
	}
	return e
}

func assertEqual(a, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("not equal %q != %q", a, b))
	}
}

func TestNewFromBase(t *testing.T) {
	_, err := NewFromBase("00000", 3)
	assertEqual(err, ErrEntropyExhausted)
	assertEqual(string(must(NewFromBase("000000", 3)).data), string([]byte{0}))
	_, err = NewFromBase("001223", 3)
	assertEqual(err, ErrParseFailure)
	assertEqual(string(must(NewFromBase("001223", 4)).data), string([]byte{0x6b}))
	assertEqual(string(must(NewFromBase("100102", 3)).data), string([]byte{0xfe}))
	assertEqual(string(must(NewFromBase("100110", 3)).data), string([]byte{0xff}))
	_, err = NewFromBase("100111", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	_, err = NewFromBase("100112", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	_, err = NewFromBase("102000", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	_, err = NewFromBase("200000", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	assertEqual(string(must(NewFromBase("0000100102", 3)).data), string([]byte{0xfe}))
	assertEqual(string(must(NewFromBase("0000100110", 3)).data), string([]byte{0xff}))
	_, err = NewFromBase("0000100111", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	_, err = NewFromBase("0000100112", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	assertEqual(string(must(NewFromBase("10022220012", 3)).data), string([]byte{0xff, 0xfe}))
	assertEqual(string(must(NewFromBase("10022220020", 3)).data), string([]byte{0xff, 0xff}))
	_, err = NewFromBase("10022220021", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	_, err = NewFromBase("10022220022", 3)
	assertEqual(err, ErrUnsafeBiasRange)
	assertEqual(string(must(NewFromBase("77777776", 8)).data), string([]byte{0xff, 0xff, 0xfe}))
	assertEqual(string(must(NewFromBase("77777777", 8)).data), string([]byte{0xff, 0xff, 0xff}))
	_, err = NewFromBase("100000000", 8)
	assertEqual(err, ErrUnsafeBiasRange)
	_, err = NewFromBase("100000001", 8)
	assertEqual(err, ErrUnsafeBiasRange)
}

func TestDigitsForBytes(t *testing.T) {
	assertEqual(DigitsForBytes(1, 3), 6)
	assertEqual(DigitsForBytes(2, 3), 11)
	assertEqual(DigitsForBytes(40, 6), 124)
}
