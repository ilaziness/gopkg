package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

func main() {

	//---------- 默认格式
	slog.Info("hello", "count", 3)
	// output: 2023/10/12 18:52:59 INFO hello count=3
	slog.Warn("hello2", "list", []int{1, 2, 3})
	// output: 2023/10/12 18:52:59 WARN hello2 list="[1 2 3]"

	//-------- 设置日志公共输出属性
	logger := slog.With("url", "/abc/sdfsdf", "test", "qw3e")
	logger.Info("hello3", "sfd", 1)
	// output: 2023/10/12 19:03:43 INFO hello3 url=/abc/sdfsdf test=qw3e sfd=1
	logger.Warn("hello4", "sfd", 2)
	// output: 2023/10/12 19:03:43 WARN hello4 url=/abc/sdfsdf test=qw3e sfd=2

	//--------- 文本格式
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	slog.Info("hello", "count", 3)
	// output: time=2023-10-12T18:54:38.800+08:00 level=INFO msg=hello count=3
	slog.Warn("hello2", "list", []int{1, 2, 3})
	// output: time=2023-10-12T18:54:38.800+08:00 level=WARN msg=hello2 list="[1 2 3]"

	//--------- json格式
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("hello", "count", 3)
	// output: {"time":"2023-10-12T18:55:18.4247282+08:00","level":"INFO","msg":"hello","count":3
	slog.Warn("hello2", "list", []int{1, 2, 3})
	// output: {"time":"2023-10-12T18:55:18.4252459+08:00","level":"WARN","msg":"hello2","list":[1,2,3]}

	//------ 设置日志输出级别
	var programLevel = new(slog.LevelVar) //线程安全
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	slog.Info("set level info hello", "count", 3)   //有输出
	slog.Debug("set level debug hello", "count", 3) //没输出

	// 设置日志输出级别是debug
	programLevel.Set(slog.LevelDebug)
	slog.Info("set level info hello", "count", 3)   //有输出
	slog.Debug("set level debug hello", "count", 3) //有输出

	//-------- 分组
	// 分组的作用是把多个属性聚合到一个key名称下面

	// 分组方式1 slog.Group
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("GET req", slog.Group("request", "method", "GET", "url", "/ab/c"))
	// text output: time=2023-10-13T11:21:05.580+08:00 level=INFO msg="GET req" request.method=GET request.url=/ab/c
	// json output: {"time":"2023-10-13T11:22:05.6329972+08:00","level":"INFO","msg":"GET req","request":{"method":"GET","url":"/ab/c"}}

	// 分组方式2 WithGroup
	l := slog.Default().WithGroup("request")
	l.Info("POST req", "method", "POST", "url", "/cd/a")
	// text output: time=2023-10-13T11:21:05.580+08:00 level=INFO msg="POST req" request.method=POST request.url=/cd/a
	// json output: {"time":"2023-10-13T11:22:05.6329972+08:00","level":"INFO","msg":"POST req","request":{"method":"POST","url":"/cd/a"}}

	//-------- context，携带ctx
	ctx := context.WithValue(context.Background(), "username", "testuser")
	slog.InfoContext(ctx, "message with ctx")

	//-------- Attrs
	slog.Info("hello", "count", 3)
	slog.Info("hello Attr", slog.Int("count", 3))
	slog.LogAttrs(ctx, slog.LevelInfo, "hello Attr2", slog.Int("count", 3)) //更高效
	// {"time":"2024-07-17T14:56:11.4407555+08:00","level":"INFO","msg":"hello","count":3}
	// {"time":"2024-07-17T14:56:11.4407555+08:00","level":"INFO","msg":"hello Attr","count":3}
	// {"time":"2024-07-17T14:56:11.4407555+08:00","level":"INFO","msg":"hello Attr2","count":3}

	//-------- 自定义类型记录行为
	// 实现LogValuer接口
	slog.Info("LogValuer", "test", TestValue("234"))

	testMyHandler()
}

type TestValue string

// LogValue 返回的值将作为实际记录的值
func (TestValue) LogValue() slog.Value {
	return slog.StringValue("TestValue")
}

// testMyHandler 自定义后端
func testMyHandler() {
	fmt.Println("-------------MyHandler-----------------")
	logger := slog.New(NewMyHandler(&Options{OutputConsole: true, OutputFile: "log.log", HandlerOptions: slog.HandlerOptions{AddSource: true}}))

	logger.Info("info msg", "a", 123)
	logger2 := logger.With("common", 789)
	logger2.Warn("warn msg", "b", "b")
	time.Sleep(time.Millisecond * 100)

	// 分组
	logger.Info("--------------测试分组--------------")
	logger3 := logger.WithGroup("request").With("url", "https://example.com")
	logger3.Error("error msg", "c", "error")
	logger3.Info("info msg", "c2", "error2")
	logger.Info("loggger msg", "c3", "d3")

	logger3.Info("info msg", slog.Group("group2", "c3", "error3"))

	logger4 := logger3.WithGroup("group3")
	logger4.Error("error msg", "c", "error", "p1", 1)
	logger4.Info("info msg", "c2", "error2", "p2", 2)
}
