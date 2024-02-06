package adapters

import (
	"reflect"
	"testing"

	clientEntities "reimagined_eureka/internal/client/entities"
)

func TestHashPassword(t *testing.T) {
	c := NewCryptoProvider()
	user := &clientEntities.User{}
	password := "testpassword123"

	err := c.HashPassword(user, password)
	if err != nil {
		t.Errorf("HashPassword() error = %v, wantErr %v", err, false)
	}
	if len(user.PasswordHash) == 0 {
		t.Errorf("Expected password hash to be set, got empty")
	}
}

func TestVerifyPassword(t *testing.T) {
	c := NewCryptoProvider()
	user := &clientEntities.User{}
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	_ = c.HashPassword(user, password)

	if err := c.VerifyPassword(user, password); err != nil {
		t.Errorf("VerifyPassword() with correct password failed: %v", err)
	}

	if err := c.VerifyPassword(user, wrongPassword); err == nil {
		t.Errorf("VerifyPassword() with wrong password succeeded; want failure")
	}
}

func TestCryptoProvider_EncryptDecrypt(t *testing.T) {
	provider := NewCryptoProvider()
	provider.SetMasterKey("some-master-key")

	plaintext := []byte("secret data")

	encryptedResult, err := provider.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	decryptedPlaintext, err := provider.Decrypt(encryptedResult)
	if err != nil {
		t.Fatalf("Decrypt() failed: %v", err)
	}

	if !reflect.DeepEqual(decryptedPlaintext, plaintext) {
		t.Errorf("Decrypted plaintext does not match original, got = %s, want = %s", decryptedPlaintext, plaintext)
	}
}

func TestCryptoProvider_EncryptErrorWithEmptyMasterKey(t *testing.T) {
	provider := NewCryptoProvider()

	_, err := provider.Encrypt([]byte("data"))
	if err == nil {
		t.Error("Encrypt() should fail when master key is empty")
	}
}

func TestCryptoProvider_HashAndVerifyHash(t *testing.T) {
	provider := NewCryptoProvider()

	data := []byte("some data to hash")
	hash, err := provider.Hash(data)
	if err != nil {
		t.Fatalf("Hash() failed: %v", err)
	}

	if err := provider.VerifyHash(data, hash); err != nil {
		t.Errorf("VerifyHash() failed to verify correct hash: %v", err)
	}

	wrongData := []byte("wrong data")
	if err := provider.VerifyHash(wrongData, hash); err == nil {
		t.Error("VerifyHash() should fail when data does not match hash")
	}
}
