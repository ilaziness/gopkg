package main

import "unicode/utf8"

// Reverse 字符串反转函数
func Reverse(s string) string {
	if !utf8.ValidString(s) {
		return s // 或返回错误
	}
	r := []rune(s) // 改用rune处理Unicode
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
