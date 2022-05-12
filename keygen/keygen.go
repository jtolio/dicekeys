package keygen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/crypto"
	csecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"

	"github.com/jtolio/dicekeys/entropy"
)

const BytesForEthKey = 40

func DeterministicEthereumKey(source *entropy.Entropy) (private, public string, err error) {
	// we're going to use two separate implementations of ScalarBaseMult on
	// ethereum's curve to make sure the output is correct.
	curve1 := secp256k1.S256()
	params1 := curve1.Params()
	curve2 := csecp256k1.S256()
	params2 := curve2.Params()

	if params1.BitSize != params2.BitSize {
		return "", "", fmt.Errorf("differing curve params")
	}
	buf := make([]byte, params1.BitSize/8+8)
	if len(buf) != BytesForEthKey {
		return "", "", fmt.Errorf("unexpected curve bitsize")
	}
	_, err = io.ReadFull(source, buf)
	if err != nil {
		return "", "", err
	}

	priv1, pub1, err := deterministicEthereumKey(curve1, params1, buf)
	if err != nil {
		return "", "", err
	}
	priv2, pub2, err := deterministicEthereumKey(curve2, params2, buf)
	if err != nil {
		return "", "", err
	}

	if priv1 != priv2 {
		return "", "", fmt.Errorf("differing results")
	}
	if pub1 != pub2 {
		return "", "", fmt.Errorf("differing results")
	}

	// make sure we haven't drifted from geth
	gethKey, err := crypto.HexToECDSA(priv1)
	if err != nil {
		return "", "", err
	}
	if pub1 != strings.ToLower(crypto.PubkeyToAddress(gethKey.PublicKey).Hex()) {
		return "", "", err
	}

	return priv1, pub1, nil
}

func deterministicEthereumKey(curve elliptic.Curve, params *elliptic.CurveParams, buf []byte) (private, public string, err error) {
	// we generate the key manually. we want perfect certainty that the data in
	// the entropy source is always used in exactly the same way to generate the
	// key, and the standard library's crypto/ecdsa code has some TODOs right
	// near the part where the randomness is read. I'd hate for a refactor to
	// break this.

	// the following code was refactored from the crypto/elliptic package
	// from the standard library
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
