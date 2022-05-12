package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jtolio/dicekeys/entropy"
	"github.com/jtolio/dicekeys/keygen"
)

var (
	flagBase = flag.Int("base", 6, "the base for the entropy input")
)

func main() {
	flag.Parse()
	neededDigits := entropy.DigitsForBytes(keygen.BytesForEthKey, *flagBase)
	fmt.Printf("need %d digits\n", neededDigits)
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	digits := strings.TrimSpace(string(data))
	if len(digits) != neededDigits {
		panic(fmt.Sprintf("%d is not the right amount of digits", len(digits)))
	}
	source, err := entropy.NewFromBase(digits, *flagBase)
	if err != nil {
		panic(err)
	}
	private, public, err := keygen.DeterministicEthereumKey(source)
	if err != nil {
		panic(err)
	}
	fmt.Println("address:", public)
	fmt.Println("wallet :", private)
}
