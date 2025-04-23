模糊测试（Fuzzing）是一种自动化软件测试技术，通过向目标程序输入大量随机或半随机的异常数据来发现潜在漏洞或异常行为。

```shell
# 运行普通测试
go test

# 运行模糊测试（默认无限运行，Ctrl+C停止）
go test -fuzz=FuzzReverse

# 限制运行时间（10秒）
go test -fuzz=FuzzReverse -fuzztime=10s

# 将崩溃输入保存到testdata目录
go test -fuzz=FuzzReverse -fuzztime=10s -fuzzminimizetime=30s
```

# `f.Add` - 添加种子语料库

作用

    提供初始的测试用例（种子输入），引导模糊测试引擎生成更多变体。

    确保测试覆盖特定的边界值或关键用例。

特点

    参数类型必须匹配：f.Add 的参数类型和数量必须与 f.Fuzz 测试函数的参数完全一致。

    可多次调用：可以添加多个种子用例。

    非必需但推荐：即使不添加种子，模糊引擎也会自动生成随机输入，但种子能提高测试效率。

```go
func FuzzReverse(f *testing.F) {
    // 添加种子用例
    f.Add("hello")      // 普通ASCII
    f.Add("世界")       // Unicode
    f.Add("!@#$%")     // 特殊字符
    f.Add("")          // 空字符串
    
    f.Fuzz(func(t *testing.T, orig string) {
        // 测试逻辑...
    })
}
```

# `f.Fuzz` - 定义模糊测试逻辑

作用

    包含实际的测试逻辑，接收模糊引擎生成的输入并验证程序行为。

    是发现程序漏洞或异常的核心部分。

特点

    参数动态生成：函数的参数（如 orig string）由模糊引擎自动生成。

    需包含断言：通过 t.Errorf 或 t.Fail 标记测试失败。

    幂等性要求：测试逻辑不应依赖外部状态（如全局变量），确保可重复执行。

```go
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
```
