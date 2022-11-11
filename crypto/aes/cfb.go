package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

// CFB加密模式

var cfbKey = []byte("4589652631457889")
var cfbRawData = []byte("测试数据test!abc")

func testCFB() {
	log.Println("CFB:")
	encrypted := aesEncryptCFB()
	log.Println(hex.EncodeToString(encrypted))
	log.Println(string(aesDecryptCFB(encrypted)))
}

func aesEncryptCFB() []byte {
	// key的长度，16、24、32byte对应aes-128、aes-192、aes-256
	block, _ := aes.NewCipher(cfbKey)
	// 最终的密文 = iv + 密文 两段
	encrypted := make([]byte, aes.BlockSize+len(cfbRawData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalln(err)
	}

	// 加密后的数据长度和加密前的数据长度是一致的
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], cfbRawData)
	return encrypted
}

func aesDecryptCFB(encrypted []byte) []byte {
	block, _ := aes.NewCipher(cfbKey)
	if len(encrypted) < aes.BlockSize {
		log.Fatalln("too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
