package encryption

import (
	"bytes"
	"golang.org/x/crypto/chacha20poly1305"
	"testing"
)

func TestEncryptDecryptAEAD(t *testing.T) {
	key := make([]byte, chacha20poly1305.KeySize)
	nonce := make([]byte, chacha20poly1305.NonceSize)
	ad := []byte("header")
	plaintext := []byte("hello world")

	cipherText, err := EncryptAEAD(nonce, key, plaintext, ad)
	if err != nil {
		t.Fatalf("EncryptAEAD returned error: %v", err)
	}

	decrypted, err := DecryptAEAD(nonce, key, cipherText, ad)
	if err != nil {
		t.Fatalf("DecryptAEAD returned error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted plaintext mismatch: got %q want %q", decrypted, plaintext)
	}
}

func TestDecryptAEADTamper(t *testing.T) {
	key := make([]byte, chacha20poly1305.KeySize)
	nonce := make([]byte, chacha20poly1305.NonceSize)
	ad := []byte("header")
	plaintext := []byte("hello world")

	cipherText, err := EncryptAEAD(nonce, key, plaintext, ad)
	if err != nil {
		t.Fatalf("EncryptAEAD returned error: %v", err)
	}

	// modify additional data to simulate tampering
	badAd := []byte("different")
	if _, err := DecryptAEAD(nonce, key, cipherText, badAd); err == nil {
		t.Fatalf("expected authentication error when additional data mismatches")
	}
}
