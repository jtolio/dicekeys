package keygen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/sha3"

	"github.com/jtolio/dicekeys/entropy"
)

const BytesForEthKey = 40

func DeterministicEthereumKey(source *entropy.Entropy) (private, public string, err error) {
	// the ethereum curve
	curve := secp256k1.S256()

	// we generate the key manually. we want perfect certainty that the data in
	// the entropy source is always used in exactly the same way to generate the
	// key, and the standard library's crypto/ecdsa code has some TODOs right
	// near the part where the randomness is read. I'd hate for a refactor to
	// break this.

	// the following code was refactored from the crypto/elliptic package
	// from the standard library
	params := curve.Params()
	buf := make([]byte, params.BitSize/8+8)
	if len(buf) != BytesForEthKey {
		return "", "", fmt.Errorf("unexpected curve bitsize")
	}
	_, err = io.ReadFull(source, buf)
	if err != nil {
		return "", "", err
	}

	one := new(big.Int).SetInt64(1)
	k := new(big.Int).SetBytes(buf)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)

	key := new(ecdsa.PrivateKey)
	key.PublicKey.Curve = curve
	key.D = k
	key.PublicKey.X, key.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())

	// ethereum specific things
	private = PrivateKeyToWallet(key)
	public = PubkeyToAddress(key.PublicKey)
	return private, public, nil
}

func PrivateKeyToWallet(k *ecdsa.PrivateKey) string {
	return hex.EncodeToString(k.D.FillBytes(make([]byte, k.Params().BitSize/8)))
}

func PubkeyToAddress(pub ecdsa.PublicKey) string {
	pubBytes := elliptic.Marshal(secp256k1.S256(), pub.X, pub.Y)
	if pubBytes[0] != 0x04 {
		// 0x04 is the prefix that defines that this is an uncompressed elliptic curve point
		panic("unexpected marshal output")
	}
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(pubBytes[1:])
	if err != nil {
		panic(err)
	}
	data := hash.Sum(nil)
	if len(data) != 32 {
		panic("unexpected keccak output")
	}
	address := data[len(data)-20:] // ethereum only uses the 20 least significant bytes
	return "0x" + hex.EncodeToString(address)
}
