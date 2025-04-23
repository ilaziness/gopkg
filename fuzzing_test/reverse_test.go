// reverse_test.go
package main

import (
	"testing"
	"unicode/utf8"
)

// 普通单元测试
func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"Hello", "olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	for _, tc := range testcases {
		rev := Reverse(tc.in)
		if rev != tc.want {
			t.Errorf("Reverse(%q) = %q, want %q", tc.in, rev, tc.want)
		}
	}
}

// 模糊测试（Go 1.18+）
func FuzzReverse(f *testing.F) {
	// 添加种子语料库（可选）
	f.Add("hello")
	f.Add("世界")

	// 模糊测试函数
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)

		// 断言1：两次反转应恢复原字符串
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}

		// 断言2：反转后的字符串应为有效UTF-8
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8: %q", rev)
		}
	})
}
