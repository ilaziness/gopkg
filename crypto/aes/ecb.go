package main

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"log"
)

// ECB加密模式
// ECB模式不安全，内置库没有直接支持

var ecbKey = []byte("4589652631457856")
var ecbRawData = []byte("测试待加密数据，test!")

func testECB() {
	log.Println("ECB：")
	encrypt := aesEncryptECB()
	log.Println(hex.EncodeToString(encrypt))
	log.Println(base64.StdEncoding.EncodeToString(encrypt))
	log.Println(string(aesDecryptECB(encrypt)))
}

func aesEncryptECB() []byte {
	// key的长度，16、24、32byte对应aes-128、aes-192、aes-256
	cipher, _ := aes.NewCipher(ecbKey)
	plain := pkcs7Padding(ecbRawData, aes.BlockSize)
	encrypted := make([]byte, len(plain))

	for bs, be := 0, cipher.BlockSize(); bs <= len(ecbRawData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

func aesDecryptECB(encrypted []byte) []byte {
	cipher, _ := aes.NewCipher(ecbKey)
	decrypted := make([]byte, len(encrypted))

	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	return pkcs7UnPadding(decrypted)
}
