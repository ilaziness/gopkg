package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ilaziness/gopkg/opentelemetry/advanced/api"
	"github.com/ilaziness/gopkg/opentelemetry/advanced/httpclient"
	"github.com/ilaziness/gopkg/opentelemetry/advanced/rpcall"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	go runRpcServer()

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

	// http call
	data, err := httpclient.Get(ctx, httpclient.ServiceNameBs, "/hello")
	if err != nil {
		log.Println("as get error:", err)
	} else {
		log.Println("as get data:", string(data))
	}

	fmt.Fprintf(w, "hello as service")
}

// runApp 启动运行http服务
func runApp() {
	// rpc client
	addr := "127.0.0.1:7001"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	rpcClient = api.NewBsRpcClient(conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	http.Handle("/hello", otelhttp.NewHandler(
		http.HandlerFunc(HelloHandler), "as service hello",
		// 下面两行可以不设置，和默认值一致
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	))

	log.Println("http server listening at port", "8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}

var rpcClient api.BsRpcClient

// runRpcServer 启动RPC服务
func runRpcServer() {
	lis, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	api.RegisterAsRpcServer(s, &AsRpc{})
	log.Println("rpc server listening at port", "7000")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type AsRpc struct {
	api.UnimplementedAsRpcServer
}

// Test test rpc服务接口
func (a *AsRpc) Test(ctx context.Context, req *api.Req) (*api.Resp, error) {
	ctx, span := tracer.Start(ctx, "as rpc Test")
	defer span.End()
	log.Println("as rpc server receive data:", req.GetA(), req.GetB())

	// rpc call 调用Bs服务的Test
	rpcall.CallBsTest(ctx, rpcClient)

	return &api.Resp{S: "this is as rpc server response"}, nil
}
