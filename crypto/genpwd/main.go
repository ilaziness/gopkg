package main

import (
	"crypto/rand"
	"log"
)

// 生成密码代码片段

func main() {
	log.Println(5, getPassword(5))
	log.Println(10, getPassword(10))
	log.Println(15, getPassword(15))
	log.Println(20, getPassword(20))
}

// Generates password of length n
func getPassword(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-/.+?=&"

	rbuf := make([]byte, n)
	if _, err := rand.Read(rbuf); err != nil {
		log.Fatalln("Unable to generate password", err)
	}

	passwd := make([]byte, n)
	for i, r := range rbuf {
		passwd[i] = letters[int(r)%len(letters)]
	}

	return string(passwd)
}
