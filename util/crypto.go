package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// AesDecrypt AES-CBC解密,PKCS#7,传入密文和密钥，[]byte
func AesDecrypt(src, key []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	dst = make([]byte, len(src))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, src)

	return PKCS7UnPad(dst), nil
}

// PKCS7UnPad PKSC#7解包
func PKCS7UnPad(msg []byte) []byte {
	length := len(msg)
	padlen := int(msg[length-1])
	return msg[:length-padlen]
}

// AesEncrypt AES-CBC加密+PKCS#7打包，传入明文和密钥
func AesEncrypt(src []byte, key []byte) ([]byte, error) {
	k := len(key)
	if len(src)%k != 0 {
		src = PKCS7Pad(src, k)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, src)

	return dst, nil
}

// PKCS7Pad PKCS#7打包
func PKCS7Pad(msg []byte, blockSize int) []byte {
	if blockSize < 1<<1 || blockSize >= 1<<8 {
		panic("unsupported block size")
	}
	padlen := blockSize - len(msg)%blockSize
	padding := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(msg, padding...)
}
