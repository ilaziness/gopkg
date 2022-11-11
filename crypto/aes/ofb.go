package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

// OFB模式

var ofbKey = []byte("4589652631457889")
var ofbRawData = []byte("测试数据test!abc")

func testOFB() {
	log.Println("OFB:")
	encrypted := aesEncryptOFB()
	log.Println(hex.EncodeToString(encrypted))
	log.Printf("'%s'\n", string(aesDecryptOFB(encrypted)))
}

func aesEncryptOFB() []byte {
	// key的长度，16、24、32byte对应aes-128、aes-192、aes-256
	block, err := aes.NewCipher(ofbKey)
	if err != nil {
		log.Fatalln(err)
	}
	//IV初始向量放在已经加密数据的开头，最终密文 = IV值 + 已加密原始数据
	ciphertext := make([]byte, aes.BlockSize+len(ofbRawData))
	iv := ciphertext[:aes.BlockSize]
	// 填充iv值
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalln(err)
	}
	// iv必须等于block size长度
	stream := cipher.NewOFB(block, iv)
	// 加密后的数据长度和加密前的数据长度是一致的
	// 加密ofbRawData,填充到ciphertext对应位置
	stream.XORKeyStream(ciphertext[aes.BlockSize:], ofbRawData)

	return ciphertext
}

func aesDecryptOFB(ciphertext []byte) []byte {
	block, err := aes.NewCipher(ofbKey)
	if err != nil {
		log.Fatalln(err)
	}
	iv = ciphertext[:aes.BlockSize]
	plaintext := ciphertext[aes.BlockSize:]
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return plaintext
}
