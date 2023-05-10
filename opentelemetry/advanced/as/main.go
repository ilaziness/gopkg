package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ilaziness/gopkg/opentelemetry/advanced/httpclient"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// serviceId
const serviceId = "AsService"

// tracer Tracer
var tracer trace.Tracer

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

func main() {
	l := log.New(os.Stdout, "", 0)
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
		propagation.Baggage{},
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

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "hello")
	defer span.End()

	// 获取自定义值
	bagg := baggage.FromContext(ctx)
	log.Println(bagg.String(), bagg.Member("memberId").Value())

	data, err := httpclient.Get(ctx, httpclient.ServiceNameBs, "/hello")
	if err != nil {
		log.Println("as get error:", err)
	} else {
		log.Println("as get data:", string(data))
	}
	fmt.Fprintf(w, "hello as service")
}

func runApp() {
	http.Handle("/hello", otelhttp.NewHandler(
		http.HandlerFunc(HelloHandler), "as service hello",
		// 下面两行可以不设置，和默认值一致
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	))
	log.Fatal(http.ListenAndServe(":8088", nil))
}
