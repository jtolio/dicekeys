package entropy

import (
	"io"
	"testing"
)

func RequireNoError(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func RequireEqual(t testing.TB, a, b interface{}) {
	if a != b {
		t.Fatalf("expected equality, %q != %q", a, b)
	}
}

func TestEntropy0(t *testing.T) {
	entropy := New([]byte("abc"))
	dest := make([]byte, 0)
	n, err := entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 0)
}

func TestEntropy1(t *testing.T) {
	entropy := New([]byte("abc"))
	dest := make([]byte, 1)
	n, err := entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 1)
	RequireEqual(t, string(dest[:n]), "a")
	n, err = entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 1)
	RequireEqual(t, string(dest[:n]), "b")
	n, err = entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 1)
	RequireEqual(t, string(dest[:n]), "c")
	_, err = entropy.Read(dest)
	RequireEqual(t, err, ErrEntropyExhausted)
}

func TestEntropy2(t *testing.T) {
	entropy := New([]byte("abc"))
	dest := make([]byte, 2)
	n, err := entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 2)
	RequireEqual(t, string(dest[:n]), "ab")
	n, err = entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 1)
	RequireEqual(t, string(dest[:n]), "c")
	_, err = entropy.Read(dest)
	RequireEqual(t, err, ErrEntropyExhausted)
}

func TestEntropy3(t *testing.T) {
	entropy := New([]byte("abc"))
	dest := make([]byte, 3)
	n, err := entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 3)
	RequireEqual(t, string(dest[:n]), "abc")
	_, err = entropy.Read(dest)
	RequireEqual(t, err, ErrEntropyExhausted)
}

func TestEntropy4(t *testing.T) {
	entropy := New([]byte("abc"))
	dest := make([]byte, 4)
	n, err := entropy.Read(dest)
	RequireNoError(t, err)
	RequireEqual(t, n, 3)
	RequireEqual(t, string(dest[:n]), "abc")
	_, err = entropy.Read(dest)
	RequireEqual(t, err, ErrEntropyExhausted)
}

func TestIOReadLimited(t *testing.T) {
	entropy := New([]byte("abc"))
	data, err := io.ReadAll(io.LimitReader(entropy, 3))
	RequireNoError(t, err)
	RequireEqual(t, string(data), "abc")
}

func TestIORead(t *testing.T) {
	entropy := New([]byte("abc"))
	_, err := io.ReadAll(entropy)
	RequireEqual(t, err, ErrEntropyExhausted)
}
