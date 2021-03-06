package algorithm

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
)

func EncryptAesBlockGCM(ctx context.Context, key, plainText string) (*string, *string, error) {
	logger := middleware.ContextData(ctx)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// ถ้าใส่ nonce ใน seal dst แล้ว nonce จะไปอยู่ใน ciphertext ด้วย
	// ciphertext := aesgcm.Seal(nonce, nonce, []byte(plainText), nil)
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plainText), nil)
	ciphertextB64 := bToB64(ciphertext)
	nonceB64 := bToB64(nonce)
	logger.Debug(fmt.Sprintf("Encrypted Message - Base64 %q Hex \"%x\" Nonce - Base64 %q Hex %x\n", ciphertextB64, ciphertext, nonceB64, nonce))
	return &ciphertextB64, &nonceB64, nil
}

func DecryptAesBlockGCM(ctx context.Context, key, nonce, cipherText string) (*string, error) {
	logger := middleware.ContextData(ctx)

	byteText, err := b64ToB(cipherText)
	if err != nil {
		return nil, err
	}
	byteNonce, err := b64ToB(nonce)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(byteText) < nonceSize {
		return nil, err
	}

	// logger.Debug(fmt.Sprintf("Nonce Size - %v | Nonce - %x", nonceSize, nonce))
	// ถ้าตอน encrypt seal dst ไม่ได้ใส่ nonce มา ก็ต้องมาสับแยก nonce ออกมาจาก ciphertext
	// byteNonce, byteText := byteText[:nonceSize], byteText[nonceSize:]
	plaintext, err := aesgcm.Open(nil, byteNonce, byteText, nil)
	if err != nil {
		return nil, err
	}
	text := string(plaintext)
	logger.Debug(fmt.Sprintf("Decrypted Message - Text %q", text))
	return &text, nil
}
