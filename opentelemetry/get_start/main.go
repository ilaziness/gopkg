package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	// go get go.opentelemetry.io/otel go.opentelemetry.io/otel/trace
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// name 是Tracer名称，用来标识应用
const name = "fib"

// newExporter 返回一个控制台导出器（exporter）
// 用来接收打印遥测数据
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// 使用人类友好输出
		stdouttrace.WithPrettyPrint(),
		// 不打印时间戳
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource 返回一个resource
// resource 表示生成的遥测实体
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("fib"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

// 运行应用
func main() {
	l := log.New(os.Stdout, "", 0)

	// 遥测数据写到文件
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	// 创建导出到文件的控制台导出器
	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	// 全局注册TracerProvider
	otel.SetTracerProvider(tp)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
