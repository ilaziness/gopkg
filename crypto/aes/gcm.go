package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
)

// AEAD
// GCM模式

var gcmRawData = []byte("GCM测试文本数据")
var gcmkey = []byte("4561235859621596")
var gcmNonce = []byte("64a9433eae7c")

func testGCM() {
	log.Println("GCM:")
	encrypted := aesEncryptGCM()
	log.Println(hex.EncodeToString(encrypted))
	log.Printf("'%s'\n", string(aesDecryptGCM(encrypted)))
}

func aesEncryptGCM() []byte {
	// key的长度，16、24、32byte对应aes-128、aes-192、aes-256
	block, err := aes.NewCipher(gcmkey)
	if err != nil {
		panic(err.Error())
	}
	//gcmNonce 需要随机生成，也可以和加密后的数据拼接在一起，方便解密的时候取出gcmNonce的值
	/*
		nonce := make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}
	*/
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	// nonce 的长度
	log.Printf("nonce size:%d\n", aesgcm.NonceSize())
	return aesgcm.Seal(nil, gcmNonce, gcmRawData, nil)
}

func aesDecryptGCM(encrypted []byte) []byte {
	block, err := aes.NewCipher(gcmkey)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, gcmNonce, encrypted, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
