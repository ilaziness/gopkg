package main

// AES加密
// 加密模式选择：如果数据通过非对称签名进行身份验证,则使用CBC,否则使用GCM
// 除了GCM之外，其他模式建议加上crypto/hmac来验证数据的完整性，GCM模式属于AEAD加密，自带数据完整性
func main() {
	//CBC
	testCBC()
	//ECB，不安全
	testECB()
	//CFB
	testCFB()
	//OFB
	testOFB()
	//CTR
	testCTR()
	//GCM
	testGCM()
}
