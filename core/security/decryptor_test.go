package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"
	"time"
)

// Helper to encrypt AES-CBC for testing
func encryptAES_CBC(plain []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// PKCS7 pad
	padding := aes.BlockSize - len(plain)%aes.BlockSize
	padtext := append(plain, bytesRepeat(byte(padding), padding)...)
	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padtext)
	return ciphertext, nil
}

func bytesRepeat(b byte, count int) []byte {
	res := make([]byte, count)
	for i := range res {
		res[i] = b
	}
	return res
}

func TestDecryptAES_CBC(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes (AES-128)
	iv := []byte("abcdefghijklmnop")  // 16 bytes
	original := "Kerahasiaan data terjamin dengan AES-128"

	encrypted, err := encryptAES_CBC([]byte(original), key, iv)
	if err != nil {
		t.Fatalf("Gagal mengenkripsi sampel: %v", err)
	}

	decrypted, err := DecryptAES_CBC(encrypted, key, iv)
	if err != nil {
		t.Fatalf("DecryptAES_CBC error: %v", err)
	}

	if string(decrypted) != original {
		t.Errorf("Hasil dekripsi '%s' tidak sama dengan aslinya '%s'", string(decrypted), original)
	}
}

func TestDecryptXOR(t *testing.T) {
	plain := []byte("TopSecretPayload2026")
	key := []byte("KUNCI_RAHASIA")

	encrypted := DecryptXOR(plain, key)
	decrypted := DecryptXOR(encrypted, key)

	if string(decrypted) != string(plain) {
		t.Errorf("XOR decryption gagal: %s vs %s", string(decrypted), string(plain))
	}
}

func TestDecryptVigenere(t *testing.T) {
	plain := "ATTACK AT DAWN!"
	key := "LEMON"
	// L=11, E=4, M=12, O=14, N=13
	// ATTACK -> LXFOPV
	cipher := "LXFOPV EF RNHR!"

	decrypted := DecryptVigenere(cipher, key)
	if decrypted != plain {
		t.Errorf("Vigenere decryption gagal: %s != %s", decrypted, plain)
	}
}

func TestDecryptWithKey_Unified(t *testing.T) {
	key := "mysecretkey12345" // 16 bytes
	iv := "myinitialvector1"  // 16 bytes
	plain := "Halo Dunia Cyber!"

	encBytes, _ := encryptAES_CBC([]byte(plain), []byte(key), []byte(iv))
	b64Cipher := base64.StdEncoding.EncodeToString(encBytes)

	res, err := DecryptWithKey("AES-128 (CBC PKCS7)", "base64", b64Cipher, key, iv)
	if err != nil {
		t.Fatalf("DecryptWithKey error: %v", err)
	}

	if res != plain {
		t.Errorf("Hasil tidak cocok: %s != %s", res, plain)
	}
}

func TestMD5LocalLookup(t *testing.T) {
	// MD5("password") = 5f4dcc3b5aa765d61d8327deb882cf99
	targetHash := "5f4dcc3b5aa765d61d8327deb882cf99"
	plain, found := LookupMD5Local(targetHash)
	if !found {
		t.Fatalf("LookupMD5Local gagal menemukan hash untuk 'password'")
	}
	if plain != "password" {
		t.Errorf("Plaintext salah: %s", plain)
	}
}

func TestMD5CustomWordlist(t *testing.T) {
	word := "kuncirahasiaku999"
	h := ComputeMD5(word)

	candidates := []string{"salah", "coba1", "kuncirahasiaku999", "admin"}
	plain, found := LookupMD5Custom(h, candidates)
	if !found {
		t.Fatalf("LookupMD5Custom gagal menemukan kata dari wordlist")
	}
	if plain != word {
		t.Errorf("Plaintext salah: %s != %s", plain, word)
	}
}

func TestReverseMD5Hybrid_Local(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// MD5("admin") = 21232f297a57a5a743894a0e4a801fc3
	res, err := ReverseMD5Hybrid(ctx, "21232f297a57a5a743894a0e4a801fc3", nil)
	if err != nil {
		t.Fatalf("ReverseMD5Hybrid error: %v", err)
	}
	if res.Plaintext != "admin" {
		t.Errorf("Hasil salah: %s", res.Plaintext)
	}
}
