package testmod

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/dsnet/try"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/jtolio/dicekeys/entropy"
	"github.com/jtolio/dicekeys/keygen"
)

func checkEthMatch(t testing.TB, source []byte) {
	defer try.F(t.Fatal)

	ent := entropy.New(source[:])
	key := try.E1(ecdsa.GenerateKey(crypto.S256(), ent))
	wallet1 := hex.EncodeToString(crypto.FromECDSA(key))
	address1 := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())

	ent = entropy.New(source[:])
	wallet2, address2 := try.E2(keygen.DeterministicEthereumKey(ent))

	require.Equal(t, wallet1, wallet2)
	require.Equal(t, address1, address2)
}

func TestDeterministicEthereumKey(t *testing.T) {
	defer try.F(t.Fatal)

	for i := 0; i < 100; i++ {
		var source [keygen.BytesForEthKey]byte
		try.E1(rand.Read(source[:]))
		checkEthMatch(t, source[:])
	}
}

func TestDeterministicEthereumKeyPin(t *testing.T) {
	defer try.F(t.Fatal)
	source := try.E1(hex.DecodeString("7a69a8d19b57f6b164c19baeeade2f4d78824a37809ff2920986b366a4f393c038bc5906c74e63c4"))
	wallet, address := try.E2(keygen.DeterministicEthereumKey(entropy.New(source)))
	require.Equal(t, wallet, "64c19baeeade2f4e141139c092279e5f3a75cbc37c89f56c3324c864c257c685")
	require.Equal(t, address, "0x0197a4a46412aad778ab07c775b8109966152688")
}

func FuzzDeterministicEthereumKey(f *testing.F) {
	var initial [keygen.BytesForEthKey]byte
	f.Add(initial[:])
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) != keygen.BytesForEthKey {
			t.Skip()
			return
		}
		checkEthMatch(t, data)
	})
}

func TestNewFromBase(t *testing.T) {
	defer try.F(t.Fatal)
	ent := try.E1(entropy.NewFromBase("7a69a8d19b57f6b164c19baeeade2f4d78824a37809ff2920986b366a4f393c038bc5906c74e63c4", 16))
	wallet, address := try.E2(keygen.DeterministicEthereumKey(ent))
	require.Equal(t, wallet, "64c19baeeade2f4e141139c092279e5f3a75cbc37c89f56c3324c864c257c685")
	require.Equal(t, address, "0x0197a4a46412aad778ab07c775b8109966152688")
}
