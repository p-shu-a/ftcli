package encryption

import (
	"crypto/aes"
	"ftcli/config"
	"golang.org/x/crypto/chacha20"
	"testing"
)

func TestGenerateIV(t *testing.T) {
	iv, err := GenerateIV()
	if err != nil {
		t.Fatalf("GenerateIV returned error: %v", err)
	}
	if len(iv) != aes.BlockSize {
		t.Errorf("expected IV length %d, got %d", aes.BlockSize, len(iv))
	}
}

func TestGenerateSalt(t *testing.T) {
	s, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt returned error: %v", err)
	}
	if len(s) != 16 {
		t.Errorf("expected salt length 16, got %d", len(s))
	}
}

func TestGenerateNonce(t *testing.T) {
	n, err := GenerateNonce()
	if err != nil {
		t.Fatalf("GenerateNonce returned error: %v", err)
	}
	if len(n) != chacha20.NonceSize {
		t.Errorf("expected nonce length %d, got %d", chacha20.NonceSize, len(n))
	}
}

func TestGenerateMasterKey(t *testing.T) {
	salt := make([]byte, 16)
	key := GenerateMasterKey(salt, "password")
	if len(key) != int(config.KeyLength) {
		t.Errorf("expected key length %d, got %d", config.KeyLength, len(key))
	}
}
