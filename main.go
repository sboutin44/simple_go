package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Ask for name
	fmt.Print("What's your name? ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Ask for password
	fmt.Print("Enter a password to encrypt your name: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Derive AES key from password using SHA-256
	key := sha256.Sum256([]byte(password))

	// Encrypt name
	encrypted, nonce, err := encrypt([]byte(name), key[:])
	if err != nil {
		panic(err)
	}
	encCombined := append(nonce, encrypted...)
	encBase64 := base64.StdEncoding.EncodeToString(encCombined)

	fmt.Println("🔐 Encrypted name (base64):", encBase64)

	// Decrypt for verification
	decoded, _ := base64.StdEncoding.DecodeString(encBase64)
	decrypted, err := decrypt(decoded, key[:])
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ Decrypted name: %s\n", string(decrypted))
	fmt.Printf("👋 Hello again, %s! Your name has %d characters.\n", decrypted, len(decrypted))
}

// encrypt encrypts plaintext using AES-256-GCM
func encrypt(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// decrypt decrypts ciphertext using AES-256-GCM
func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
