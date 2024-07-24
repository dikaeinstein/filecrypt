package crypto

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	testMessage   = []byte("Do not go gentle into that good night.")
	testPassword1 = []byte("correct horse battery staple")
	// testPassword2 = []byte("incorrect horse battery staple")
)

// generateKey creates a new random secret key.
func generateKey(t *testing.T) *[keySize]byte {
	t.Helper()

	key := new([keySize]byte)

	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		t.Fatal(err)
	}

	return key
}

func TestEncrypt(t *testing.T) {
	testKey := generateKey(t)
	ct, err := encrypt(testKey, testMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}

	pt, err := decrypt(testKey, ct)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !bytes.Equal(testMessage, pt) {
		t.Fatalf("messages don't match")
	}
}

func TestDecryptFailures(t *testing.T) {
	targetLength := 24 + secretbox.Overhead
	testKey := generateKey(t)

	for i := 0; i < targetLength; i++ {
		buf := make([]byte, i)
		if _, err := decrypt(testKey, buf); err == nil {
			t.Fatal("expected decryption failure with bad message length")
		}
	}

	otherKey := generateKey(t)

	ct, err := encrypt(testKey, testMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if _, err = decrypt(otherKey, ct); err == nil {
		t.Fatal("decrypt should fail with wrong key")
	}
}

func TestEncryptCycle(t *testing.T) {
	out, err := Seal(testPassword1, testMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}

	out, err = Open(testPassword1, out)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !bytes.Equal(testMessage, out) {
		t.Fatal("recovered plaintext doesn't match original")
	}
}
