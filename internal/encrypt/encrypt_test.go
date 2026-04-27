package encrypt_test

import (
	"bytes"
	"testing"

	"github.com/user/envdiff/internal/encrypt"
)

var testKey = []byte("0123456789abcdef") // 16-byte AES-128 key

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := []byte("SECRET_VALUE=super-secret-123")

	encoded, err := encrypt.Encrypt(plaintext, testKey)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	got, err := encrypt.Decrypt(encoded, testKey)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if !bytes.Equal(got, plaintext) {
		t.Errorf("Decrypt() = %q, want %q", got, plaintext)
	}
}

func TestEncrypt_ProducesUniqueCiphertexts(t *testing.T) {
	plaintext := []byte("same-value")

	a, err := encrypt.Encrypt(plaintext, testKey)
	if err != nil {
		t.Fatalf("first Encrypt() error: %v", err)
	}
	b, err := encrypt.Encrypt(plaintext, testKey)
	if err != nil {
		t.Fatalf("second Encrypt() error: %v", err)
	}

	if a == b {
		t.Error("expected distinct ciphertexts due to random nonce, got identical output")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := encrypt.Decrypt("not!!valid%%base64", testKey)
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	encoded, _ := encrypt.Encrypt([]byte("value"), testKey)
	tampered := encoded[:len(encoded)-4] + "XXXX"

	_, err := encrypt.Decrypt(tampered, testKey)
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext for tampered data, got %v", err)
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	_, err := encrypt.Encrypt([]byte("value"), []byte("shortkey"))
	if err != encrypt.ErrShortKey {
		t.Errorf("expected ErrShortKey, got %v", err)
	}
}

func TestDecrypt_InvalidKeyLength(t *testing.T) {
	_, err := encrypt.Decrypt("someencoded", []byte("shortkey"))
	if err != encrypt.ErrShortKey {
		t.Errorf("expected ErrShortKey, got %v", err)
	}
}
