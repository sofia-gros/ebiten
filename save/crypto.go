package save

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// encryptAES は指定されたキー文字列 (SHA256ハッシュで32byte化) を用いてデータを AES-GCM 暗号化します。
func encryptAES(plainData []byte, keyStr string) ([]byte, error) {
	if keyStr == "" {
		return plainData, nil
	}

	key := sha256.Sum256([]byte(keyStr))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plainData, nil)
	return ciphertext, nil
}

// decryptAES は AES-GCM で暗号化されたデータを復号します。
func decryptAES(cipherData []byte, keyStr string) ([]byte, error) {
	if keyStr == "" {
		return cipherData, nil
	}

	key := sha256.Sum256([]byte(keyStr))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
	plainData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (invalid key or corrupted data): %w", err)
	}

	return plainData, nil
}

// attachChecksum はデータ末尾に SHA-256 チェックサム (32bytes) を付加します。
func attachChecksum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return append(data, hash[:]...)
}

// verifyAndStripChecksum はデータ末尾の SHA-256 チェックサムを検証し、本体データを抽出します。
func verifyAndStripChecksum(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, errors.New("data too short for checksum verification")
	}

	content := data[:len(data)-32]
	expectedHash := data[len(data)-32:]
	actualHash := sha256.Sum256(content)

	if !bytes.Equal(expectedHash, actualHash[:]) {
		return nil, errors.New("checksum verification failed: data tampering or corruption detected")
	}

	return content, nil
}
