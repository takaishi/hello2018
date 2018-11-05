package conn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func aes256key() []byte {
	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")
	return key
}

type RecordCrypt interface {
	Encrypt(dst, plaintext []byte) ([]byte, error)
	Decrypt(dst, cipertext []byte) ([]byte, error)
}

type HelloRecordCrypt struct{}

func (hrc *HelloRecordCrypt) Encrypt(dst, plainText []byte) ([]byte, error) {
	//log.Printf("[DEBUG] Encrypt| plainText = %+v\n", plainText)
	block, err := aes.NewCipher(aes256key())
	if err != nil {
		return nil, err
	}

	dst = make([]byte, aes.BlockSize+len(plainText))
	iv := dst[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(dst[aes.BlockSize:], plainText)
	//log.Printf("[DEBUG] Encrypt| cipherText = %+v\n", dst)
	//copy(dst, cipherText)
	//log.Printf("[DEBUG] Encrypt| dst = %+v\n", dst)

	return dst, nil
}

func (hrc *HelloRecordCrypt) Decrypt(dst, cipherText []byte) ([]byte, error) {
	//log.Printf("[DEBUG] Decrypt| cipherText = %v\n", cipherText)
	block, err := aes.NewCipher(aes256key())
	if err != nil {
		return nil, err
	}

	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
	//log.Printf("[DEBUG] Decrypt| decryptedText = %v\n", decryptedText)

	//copy(dst, decryptedText)

	return decryptedText, nil
}
