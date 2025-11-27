# 基于依赖注入的框架Fix的demo

## 链接
- github: [uber-go/fx](https://github.com/uber-go/fx)
- go doc: [uber-go/fx](https://pkg.go.dev/go.uber.org/fx)

## 概念

### 参数结构体

```go
type HandlerParams struct {
	fx.In

	Users    *UserGateway
	Comments *CommentGateway
	Posts    *PostGateway
	Votes    *VoteGateway
	AuthZ    *AuthZGateway

    // 如果新增字段，可以添加optional:"true"标签标记为可选，使其保持向后兼容
    Logger *zap.Logger `optional:"true"`

    // `name:"..."`命名值，用来区别相通类型的不同值，将会根据名字来注入
    WriteToConn  *sql.DB `name:"rw"`
	ReadFromConn *sql.DB `name:"ro"`
}

// 多个参数合并到一个结构体里面，结构体里面的每一个字段依赖都会自动注入
func NewHandler(p HandlerParams) *Handler {
	// ...
}
```

命名一般使用`Params`做后缀，嵌入`fx.In`，结构体里面的每一个字段都是一个依赖，会被框架自动注入。

### 结果结构体

```go
type GatewaysResult struct {
	fx.Out

	Users    *UserGateway
	Comments *CommentGateway
	Posts    *PostGateway

    // 同样结果结构体也有命名值
    ReadWrite *sql.DB `name:"rw"`
	ReadOnly  *sql.DB `name:"ro"`
}

func SetupGateways(conn *sql.DB) (GatewaysResult, error) {
	// ...
}
```

命名一般使用`Result`做后缀，嵌入`fx.Out`，其他构造函数可以直接依赖结果结构体。

### 注解（Annotation）

可以使用`fx.Annotate`来包装参数结构体和结果结构体：

```go
// 原始
fx.Provide(
    NewHTTPClient,
)

// 使用Annotate包装，还可以给结果结构体命名
fx.Provide(
    fx.Annotate(
        NewHTTPClient,
        fx.ResultTags(`name:"client"`),
    ),
),

func NewEmitter(watchers []Watcher) (*Emitter, error) {}
// 原始
fx.Provide(
    NewEmitter,
),

// 通过注解提供，还可以指明参数是哪一个值组，这样就不需要像下面声明值组那样写一堆
fx.Annotate(
    NewEmitter,
    fx.ParamTags(`group:"watchers"`),
),
```

将结构体转换为接口：

```go
// 原始
func NewHTTPClient(Config) (*http.Client, error) {}
func NewGitHubClient(client *http.Client) *github.Client {}
fx.Provide(
    NewHTTPClient,
    NewGitHubClient,
),


// 使用注解函数将结构体转换为接口, 使用者与实现解耦，可以独立测试
type HTTPClient interface {
   Do(*http.Request) (*http.Response, error)
}

// This is a compile-time check that verifies
// that our interface matches the API of http.Client.
var _ HTTPClient = (*http.Client)(nil)

func NewGitHubClient(client HTTPClient) *github.Client {}

fx.Provide(
    fx.Annotate(
        NewHTTPClient,
        fx.As(new(HTTPClient)),
    ),
    NewGitHubClient,
),
```

### 值组

值组是相同类型的值的集合value groups。

```go
// Handler 是一个处理请求的函数类型
type Handler func(ctx context.Context, req string) string

// 结果结构体
type HandlerResult struct {
	fx.Out

    // `group:"..."`标签，用来将值分组，相同组名的值会被注入到同一个切片中
	Handler Handler `group:"server"`
}

// 使用构造函数给值组提供值
func NewHelloHandler() HandlerResult {
	return HandlerResult{
		Handler: &HelloHandler{},
	}
}

func NewEchoHandler() HandlerResult {
	return HandlerResult{
		Handler: &EchoHandler{},
	}
}

// 参数结构体
type ServerParams struct {
	fx.In

    // 接收上面值组得值，它会执行所有为该组提供值的构造函数，顺序未定，然后将所有结果汇聚到单一切片
	Handlers []Handler `group:"server"`
}

func NewServer(p ServerParams) *Server {
	server := newServer()
	for _, h := range p.Handlers {
		server.Register(h)
	}
	return server
}
```

#### 软值组

默认情况下，当构造函数声明对某个值组的依赖时，该值组提供的所有值都会立即实例化。

声明软值组，若构造函数的输出类型仅被软值组调用，则该构造函数将不会被执行。

```go
type Params struct {
	fx.In

    // 在值组标签里面增加soft标记，软值组只能在输入参数里面使用，不能在结果结构体里面使用
    // 通过这样的声明，只有当存在另一个已实例化的组件使用该构造函数的结果时，才会调用向“server”值组提供值的构造函数。
    // 只有当该值组中的某个元素被“硬依赖”（即被某个非 soft 的输入实际使用）时，提供该值组成员的构造函数才会被调用。
	Handlers []Handler `group:"server,soft"` // 软依赖
	Logger   *zap.Logger // 硬依赖
}

// NewServer 依赖*zap.Logger
func NewServer(p Params) *Server {
    // 因为NewHandler没用被调用，所以Params.Handles只会有NewHandlerAndLogger提供那一个Handler
	// ...
}

// NewServer 依赖*zap.Logger，NewHandlerAndLogger会被调用
func NewHandlerAndLogger() (Handler, *zap.Logger) {
	// ...
}

// NewHandler只提供了server组的值，则不会被调用
func NewHandler() Handler {
	// ...
}

// 由于 Logger 会被应用程序使用，因此会调用 NewHandlerAndLogger，但不会调用 NewHandler，因为它仅被软值组使用。
fx.Provide(
	fx.Annotate(NewHandlerAndLogger, fx.ResultTags(`group:"server"`)),
	fx.Annotate(NewHandler, fx.ResultTags(`group:"server"`)),
)
```

- 软值组（soft group）不会主动触发 provider；
- 如果一个 provider 同时提供 soft group 成员 + 其他被硬依赖的值，它仍会被调用；
- 这个机制可用于按需初始化：只有当某个功能真正被使用时，才初始化相关组件。