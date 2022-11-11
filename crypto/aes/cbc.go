package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
)

//CBC模式

var iv = []byte("1234567891011121")

func testCBC() {
	txt := "123456测试"
	key := []byte("12b57eb210a6bf257297797e93a04dfa")

	// CBC模式
	log.Println("CBC:")
	encryptData := aesEncryptCBC([]byte(txt), key)
	log.Println(
		hex.EncodeToString(encryptData),
	)
	log.Println(string(aesDecryptCBD(encryptData, key)))
}

// CBC模式
func aesEncryptCBC(rawData, key []byte) (encryptData []byte) {
	// key的长度，16、24、32byte对应aes-128、aes-192、aes-256
	block, _ := aes.NewCipher(key)
	blokSize := block.BlockSize()
	rawData = pkcs7Padding(rawData, blokSize)
	// iv的长度必须和blokSize的大小一样
	// iv需要唯一，不需要保密，所以可以和密文拼接在一起，最终密文 = iv + 密文，这样解密的时候就可以拿到iv的值
	// 这里用的固定的
	blockModel := cipher.NewCBCEncrypter(block, iv)
	// encryptData 的长度必须大于等于rawData的长度
	// 加密后的数据长度和加密前的数据长度是一致的
	encryptData = make([]byte, len(rawData))
	blockModel.CryptBlocks(encryptData, rawData)
	return
}

func aesDecryptCBD(ciphertext, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockModel := cipher.NewCBCDecrypter(block, iv)
	encryptData := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(encryptData, ciphertext)
	return pkcs7UnPadding(encryptData)
}

// 填充模式：
// PKCS5填充blockSize = 8 byte
// PKCS7填充blockSize = 1-255 byte
// ZeroPadding填充，固定填充0
// 上面三种模式，不管能不能对齐都需要填充
// 还有一种Nopadding，就是自己要对分组对齐负责，想怎么填充就怎么填充只要块分组能够对齐

// 原始数据长度需要按照分块大小对齐，不足的地方需要填充上去
// 数据长度刚好是blocksize的倍数也需要填充
func pkcs7Padding(rawData []byte, blockSize int) []byte {
	// 需要填充的长度
	padding := blockSize - len(rawData)%blockSize
	// 填充长度作为填充内容
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(rawData, padText...)
}

func pkcs7UnPadding(rawData []byte) []byte {
	length := len(rawData)
	// 填充长度
	unpadding := int(rawData[length-1])
	// 原始数据 = 第一个字节 到 （总长度 - 填充长度）的位置
	return rawData[:(length - unpadding)]
}
