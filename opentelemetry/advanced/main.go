package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ilaziness/gopkg/opentelemetry/advanced/httpclient"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// serviceId
const serviceId = "HelloService"

// tracer Tracer
var tracer trace.Tracer

// newExport create export
func newExporter(url string) (sdktrace.SpanExporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

// newTracerProvider return TracerProvider
func newTracerProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceId),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

var l *log.Logger

func main() {
	l = log.New(os.Stdout, "", 0)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	ctx := context.Background()

	exp, err := newExporter("http://127.0.0.1:14268/api/traces")
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	tp := newTracerProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)
	//otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{}, // 用这个可以传自定义的值到其他服务
	))
	tracer = tp.Tracer(serviceId)

	go runApp()

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

// runApp 启动http服务，设置了一个handler
func runApp() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(context.Background(), "hello")
		defer span.End()
		span.SetAttributes(attribute.String("path", "/hello"))

		// 传自定义值
		member, _ := baggage.NewMember("memberId", "234324")
		bgg, _ := baggage.New(member)

		data, err := httpclient.Get(baggage.ContextWithBaggage(ctx, bgg), httpclient.ServiceNameAs, "/hello")
		if err != nil {
			log.Println("main get error:", err)
		} else {
			log.Println("main get data:", string(data))
		}
		fmt.Fprintf(w, "hello")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
