package main

// sha256

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func HexHash256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func AESEncrypt(plaintext, key, iv []byte) (b []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			b = nil
			err = fmt.Errorf("%v", e)
		}
	}()
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypt := cipher.NewCBCEncrypter(block, iv)
	var source = PKCS5Padding(plaintext, 16)
	var dst = make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	return dst, nil
}

func AESDecrypt(ciphertext, key, iv []byte) (b []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			b = nil
			err = fmt.Errorf("%v", e)
		}
	}()
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypt := cipher.NewCBCDecrypter(block, iv)
	var dst = make([]byte, len(ciphertext))
	if len(dst) < len(ciphertext) || len(dst)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid ciphertext")
	}
	decrypt.CryptBlocks(dst, ciphertext)
	dst, err = PKCS5UnPadding(dst, 16)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func AESBase64Encrypt(originData []byte, key []byte, iv []byte) (base64Result string, err error) {
	encrypted, err := AESEncrypt(originData, key, iv)
	base64Result = base64.StdEncoding.EncodeToString(encrypted)
	return
}

func AESBase64Decrypt(encryptData string, key []byte, iv []byte) (originData string, err error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return
	}
	originDataByte, err := AESDecrypt(encrypted, key, iv)
	if err != nil {
		return
	}
	originData = string(originDataByte)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	paddingSize := blockSize - len(ciphertext)%blockSize
	paddingData := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(ciphertext, paddingData...)
}

func PKCS5UnPadding(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padLen := int(data[len(data)-1])
	if padLen > blocklen || padLen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	pad := data[len(data)-padLen:]
	for i := 0; i < padLen; i++ {
		if pad[i] != byte(padLen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padLen], nil
}
