package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	Name string `json:"name"`
}

func main() {
	fx.New(
		// module, // 注册模块，模块里面的provide会被注册到fx.New

		// 使用fx.Supply提供非接口值（具体的类型）
		fx.Supply(&Config{
			Name: "my-app",
		}),
		// Provide 用于注册构造函数，构造函数可以依赖其他类型，返回多个对象或者错误
		// 只有需要时才会调用构造函数创建单例对象，Provide提供的值在其他构造函数里面都是可用的
		fx.Provide(
			NewHTTPServer,
			zap.NewExample,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			AsRoute(NewEchoHandler),
			AsRoute(NewHelloHandler),
		),
		// Invoke 用于注册应用启动的时候执行的函数
		// 执行函数的参数使用Provide已注册的构造函数来构造
		// 多个Invoke会按顺序执行
		// 通常用来启动服务器或主循环，启动后台任务，设置全局配置等
		fx.Invoke(
			GlobalConfig,
			func(srv *http.Server) {
				http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "Hello, World!")
				}))
			}),
	).Run()
}

// module 定义了应用的模块，包含了所有的构造函数和路由，可以打包到一起注册到fx.New
// 可以避免main膨胀，每个模块可以自己管理自己的provide
// fx.Invoke,fx.Decorate 都可以放到module里面
var module = fx.Module(
	"app",
	fx.Provide(
		fx.Private, // 私有提供，只能在模块内部使用，外部无法依赖注入，构造函数的结果只能在模块内部使用
		NewHTTPServer,
		zap.NewExample,
		fx.Annotate(
			NewServeMux,
			fx.ParamTags(`group:"routes"`),
		),
		AsRoute(NewEchoHandler),
		AsRoute(NewHelloHandler),
	),
)

func GlobalConfig(cfg *Config) {
	fmt.Printf("GlobalConfig: %+v\n", cfg)
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("HTTP server is listening on", zap.String("addr", srv.Addr))
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("HTTP server is shutting down")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

type Route interface {
	http.Handler

	// Pattern 返回路由的路径
	Pattern() string
}

// EchoHandler 处理 /echo 路由
// curl -X POST -d 'hello' http://localhost:8080/echo
type EchoHandler struct {
	log *zap.Logger
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log: log}
}

func (h *EchoHandler) Pattern() string {
	return "/echo"
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		h.log.Error("Failed to handle request", zap.Error(err))
	}
}

type HelloHandler struct {
	log *zap.Logger
}

func NewHelloHandler(log *zap.Logger) *HelloHandler {
	return &HelloHandler{log: log}
}

func (h *HelloHandler) Pattern() string {
	return "/hello"
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HelloHandler is handling request")
	fmt.Fprintf(w, "Hello, World!")
}

// NewServeMux 注册路由
func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
