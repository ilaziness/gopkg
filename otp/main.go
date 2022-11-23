package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

// otp全称One Time Password，一次性密码，通常用来做为账户的二次认证
// 分两种算法，TOTP基于时间的, HTOP基于次数，TOTP比较常用
// 手机客户端：Google身份验证器、微软的Authenticator
//
// TOTP使用流程:
// 注册：
// 1、生成TOTP的key，输出一个二维码给用户用扫码
// 2、扫码成功后再输入客户端生成的一次性密码验证
// 3、验证通过之后保存key的密钥secret到用户帐号资料里面
//
// 登录：
// 输入一次性密钥，提交验证即可

func main() {
	log.Println("HTOP:")
	testTOTP()
	log.Println("HTOP:")
	testHOTP()
}

func testTOTP() {
	// 生成密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Google",
		AccountName: "test@gmail.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		log.Fatalln(err)
	}
	png.Encode(&buf, img)

	//buf写入到http响应的body里面，网页展示
	// ...

	//二维码写入图片，用Google身份验证器扫描二维码

	err = os.WriteFile("qr.png", buf.Bytes(), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	//keyStr 密钥在验证一次性密码后保存到用户信息里面，二次认证时使用
	keyStr := key.Secret()

	validateCode(keyStr)

}

// 验证一次性验证码
func validateCode(secretKey string) {
	// passcode Google身份验证器获取，再输入
	passcode := ""
	fmt.Print("输入验证码：")
	fmt.Scan(&passcode)
	valid := totp.Validate(passcode, secretKey)
	if valid {
		log.Println("test@gmail.com", secretKey)
	} else {
		log.Println("passcode error.")
	}
}

// 基于次数
func testHOTP() {
	// 生成密钥
	key, err := hotp.Generate(hotp.GenerateOpts{
		Issuer:      "Google",
		AccountName: "test@gmail.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		log.Fatalln(err)
	}
	png.Encode(&buf, img)

	//buf写入到http响应的body里面，网页展示
	// ...

	//二维码写入图片，用Google身份验证器扫描二维码
	err = os.WriteFile("qr.png", buf.Bytes(), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	//keyStr 密钥在验证一次性密码后保存到用户信息里面，二次认证时使用
	keyStr := key.Secret()

	// 验证验证码
	// passcode Google身份验证器获取，再用户输入
	var counter uint64 = 10
	// 客户端生成一次性密码的时候需要提供次数
	passcode, err := hotp.GenerateCode(keyStr, counter)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(passcode, key.Secret(), key.URL())
	// 验证
	valid := hotp.Validate(passcode, counter, keyStr)
	if valid {
		log.Println("test@gmail.com", keyStr)
	} else {
		log.Println("passcode error.")
	}

	// 输入验证
	for {
		passcode = ""
		inputCounter := "0"
		fmt.Print("输入验证码：")
		fmt.Scan(&passcode)

		fmt.Print("输入计数：")
		fmt.Scan(&inputCounter)
		number, err := strconv.Atoi(inputCounter)
		if err != nil {
			log.Println("错误的计数")
			continue
		}
		counter = uint64(number)

		// 验证
		valid := hotp.Validate(passcode, counter, keyStr)
		if valid {
			log.Println("test@gmail.com", keyStr)
		} else {
			log.Println("passcode error.")
		}
	}
}
