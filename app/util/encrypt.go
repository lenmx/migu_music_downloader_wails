package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

func DesEncrypt(origData, key string) string {
	keyBytes := []byte(key)
	origDataBytes := []byte(origData)

	block, _ := des.NewCipher(keyBytes)
	bs := block.BlockSize()

	origDataBytes = PKCS5Padding(origDataBytes, bs)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:bs])
	crypted := make([]byte, len(origDataBytes))
	blockMode.CryptBlocks(crypted, origDataBytes)
	return base64.StdEncoding.EncodeToString(crypted)
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func DesDecrypt(crypted, key string) string {
	cryptedBytes := []byte(crypted)
	keyBytes := []byte(key)

	block, _ := des.NewCipher(keyBytes)
	blockMode := cipher.NewCBCDecrypter(block, keyBytes)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, cryptedBytes)
	origData = PKCS5UnPadding(origData)
	return string(origData)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
