package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// App是一个斐波那契计算应用
type App struct {
	r io.Reader
	l *log.Logger
}

// NewApp 创建App对象
func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

// Run 启动轮训，等待用户计算请求和返回计算就结果
func (a *App) Run(ctx context.Context) error {
	for {
		// 检测Run方法，每个循环都需要创建一个新的root span和context
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}

		a.Write(ctx, n)
	}
}

// Poll 询问用户输入，然后返回请求计算数
func (a *App) Poll(ctx context.Context) (uint, error) {
	// 检测Poll方法
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	if err != nil {
		// 记录错误
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	// 添加属性，如果希望在查看遥测时可以看到一些信息可以添加进去
	// n转成字符串，防止移除
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, nil
}

// Write 输出斐波那契数计算结果给用户
func (a *App) Write(ctx context.Context, n uint) {
	// 检测Write方法
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	// 跟踪对核心逻辑Fibonacci的调用
	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		// 调用Fibonacci函数计算
		f, err := Fibonacci(n)
		if err != nil {
			// 跟踪记录错误
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return f, err
	}(ctx)

	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}
