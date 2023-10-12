package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"os"
)

type Test struct {
	A string
	B int
	C *int
}

// openssl 3.0.8 生成rsa密钥，格式PKCS8，PEM格式保存
// 私钥(PKCS8)：openssl genrsa -out prv.pem 2048
// 公钥(x509)：openssl rsa -in prv.pem -out pub.pem -pubout
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	text := "我是hello34"
	log.Println("原文：", text)
	log.Println("--------字符串密钥--------")
	cipher := encryptByStrKey(text)
	log.Println("密文：", base64.StdEncoding.EncodeToString(cipher))
	log.Println("解密：", string(decryptByStrKey(cipher)))

	// PEM 文件里面的保存的数据就是der编码的数据
	log.Println("--------PEM密钥--------")
	cipher = encryptByPemKey(text)
	log.Println("密文：", base64.StdEncoding.EncodeToString(cipher))
	log.Println("解密：", string(decryptByPemKey(cipher)))
}

var (
	// 字符串密钥是pem密钥BEGIN到END之间的内容，去掉换行
	StrKeyPublicFile  = "RSA2048_pub.txt"
	StrKeyPrivateFile = "RSA2048_prv.txt"

	PemKeyPublicFile  = "pub.pem"
	PemKeyPrivateFile = "prv.pem"
)

func encryptByStrKey(text string) (cipher []byte) {
	strKey, err := os.ReadFile(StrKeyPublicFile)
	if err != nil {
		log.Println(err)
		return
	}
	pbkstr, err := base64.StdEncoding.DecodeString(string(strKey))
	if err != nil {
		log.Println(err)
		return
	}
	pubki, err := x509.ParsePKIXPublicKey(pbkstr)
	if err != nil {
		log.Println(err)
		return
	}
	pubk, ok := pubki.(*rsa.PublicKey)
	if !ok {
		log.Println("public key断言错误")
	}
	cipher, err = rsa.EncryptPKCS1v15(rand.Reader, pubk, []byte(text))
	if err != nil {
		log.Println(err)
		return
	}
	return
}
func decryptByStrKey(cipher []byte) (text []byte) {
	strKey, err := os.ReadFile(StrKeyPrivateFile)
	if err != nil {
		log.Println(err)
		return
	}
	prvkstr, err := base64.StdEncoding.DecodeString(string(strKey))
	if err != nil {
		log.Println(err)
		return
	}
	// PKCS8格式密钥
	prvki, err := x509.ParsePKCS8PrivateKey(prvkstr)
	if err != nil {
		log.Println(err)
		return
	}
	prvk, ok := prvki.(*rsa.PrivateKey)
	if !ok {
		log.Println("private key断言错误")
	}
	text, err = rsa.DecryptPKCS1v15(rand.Reader, prvk, cipher)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func encryptByPemKey(text string) (cipher []byte) {
	pemKey, err := os.ReadFile(PemKeyPublicFile)
	if err != nil {
		log.Println(err)
		return
	}
	block, _ := pem.Decode(pemKey)
	if block == nil {
		log.Println("pem public key error")
		return
	}
	pubki, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return
	}
	pubk, ok := pubki.(*rsa.PublicKey)
	if !ok {
		log.Println("public key断言错误")
		return
	}
	// EncryptPKCS1v15函数加密不安全，文档建议用EncryptOAEP
	//cipher, err = rsa.EncryptPKCS1v15(rand.Reader, pubk, []byte(text))
	// label随意，可以为空，label的主要作用是同一个公钥加密可以使用不通的label来确保安全
	label := []byte("test")
	cipher, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, pubk, []byte(text), label)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func decryptByPemKey(cipher []byte) (text []byte) {
	pemKey, err := os.ReadFile(PemKeyPrivateFile)
	if err != nil {
		log.Println(err)
		return
	}
	block, _ := pem.Decode(pemKey)
	if block == nil {
		log.Println("pem private key error")
		return
	}
	// block.Bytes 就是der编码的数据
	prvki, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return
	}
	prvk, ok := prvki.(*rsa.PrivateKey)
	if !ok {
		log.Println("public key断言错误")
		return
	}
	//text, err = rsa.DecryptPKCS1v15(rand.Reader, prvk, cipher)
	label := []byte("test")
	text, err = rsa.DecryptOAEP(sha256.New(), nil, prvk, cipher, label)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
