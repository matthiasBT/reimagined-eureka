package adapters

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"

	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/common"
)

const iterationCount = 4096
const keyLength = 32

var kdfHashing = sha3.New512
var ErrHashMismatch = errors.New("hash mismatch")

type CryptoProvider struct {
	masterKey string
}

func NewCryptoProvider() *CryptoProvider {
	return &CryptoProvider{}
}

func (c *CryptoProvider) HashPassword(user *clientEntities.User, password string) error {
	pwdHash, err := c.hashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = pwdHash
	return nil
}

func (c *CryptoProvider) VerifyPassword(user *clientEntities.User, password string) error {
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return fmt.Errorf("failed to check password hash: %v", err)
	}
	return nil
}

func (c *CryptoProvider) Hash(data []byte) ([]byte, error) {
	hasher := sha512.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}
	hashBytes := hasher.Sum(nil)
	return hashBytes, nil
}

func (c *CryptoProvider) VerifyHash(data, target []byte) error {
	value, err := c.Hash(data)
	if err != nil {
		return err
	}
	if bytes.Equal(value, target) {
		return nil
	}
	return ErrHashMismatch
}

func (c *CryptoProvider) SetMasterKey(masterKey string) {
	c.masterKey = masterKey
}

func (c *CryptoProvider) Encrypt(what []byte) (*common.EncryptionResult, error) {
	if c.masterKey == "" {
		return nil, fmt.Errorf("empty master key")
	}
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("encryption key generation failed: %v", err)
	}
	key := c.deriveKey(c.masterKey, salt)
	var block cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher generation failed: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM mode encryption failed: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("GCM nonce generation failed: %v", err)
	}
	ciphertext := gcm.Seal(nil, nonce, what, nil) // TODO: use additional data
	return &common.EncryptionResult{Ciphertext: ciphertext, Salt: salt, Nonce: nonce}, nil
}

func (c *CryptoProvider) Decrypt(what *common.EncryptionResult) ([]byte, error) {
	key := c.deriveKey(c.masterKey, what.Salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher block: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating AES-GCM cipher: %v", err)
	}
	decryptedData, err := gcm.Open(nil, what.Nonce, what.Ciphertext, nil) // TODO: use additional data
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %v", err)
	}
	return decryptedData, nil
}

func (c *CryptoProvider) hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	return hashedPassword, nil
}

func (c *CryptoProvider) deriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, iterationCount, keyLength, kdfHashing)
}
