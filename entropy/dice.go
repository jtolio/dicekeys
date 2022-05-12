package entropy

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

var (
	ErrUnsafeBiasRange = errors.New("value outside of safe bias range. please generate new entropy")
	ErrParseFailure    = errors.New("unparseable input")
)

func DigitsForBytes(byteCount, base int) (digitsNeeded int) {
	possibilities := math.Pow(256, float64(byteCount))
	rolls := math.Log2(possibilities) / math.Log2(float64(base))
	return int(math.Ceil(rolls))
}

func NewFromBase(digits string, base int) (*Entropy, error) {
	possibilities := math.Pow(float64(base), float64(len(digits)))
	maxUnbiasedBytes := int(math.Log2(possibilities) / math.Log2(256))
	if maxUnbiasedBytes <= 0 {
		return nil, ErrEntropyExhausted
	}

	two := new(big.Int).SetInt64(2)
	safeValueMax := new(big.Int).Exp(two, new(big.Int).SetInt64(8*int64(maxUnbiasedBytes)), nil)

	width := new(big.Int).SetInt64(int64(len(digits)))
	bigBase := new(big.Int).SetInt64(int64(base))
	bigPossibilities := new(big.Int).Exp(bigBase, width, nil)

	if bigPossibilities.Cmp(safeValueMax) < 0 {
		return nil, fmt.Errorf("unreachable")
	}

	parsed, ok := new(big.Int).SetString(digits, base)
	if !ok {
		return nil, ErrParseFailure
	}

	if parsed.Cmp(safeValueMax) >= 0 {
		return nil, ErrUnsafeBiasRange
	}

	data := parsed.FillBytes(make([]byte, maxUnbiasedBytes))
	return New(data[len(data)-maxUnbiasedBytes:]), nil
}
