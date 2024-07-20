# 编写`slog`处理程序的指南

> 原文[Guide](https://github.com/golang/example/blob/master/slog-handler-guide/README.md)

标准库的 `log/slog` 包采用两部分设计。由 `Logger` 类型实现的“前端”收集经过结构化的日志信息（如消息、级别和属性），并将它们传递给“后端”，即 `Handler` 接口的实现。该软件包带有两个内置处理程序，通常应该足够了。但是您可能需要编写自己的处理程序，这并不总是那么简单。本指南随时为您提供帮助。

# 记录器(logger)及其处理程序(handler)

编写处理程序需要了解 `Logger` 和 `Handler` 类型如何协同工作。

每个记录器都包含一个处理程序。某些 `Logger` 方法会执行一些预备工作，例如将键值对收集到 `Attrs` 中，然后调用一个或多个 `Handler` 方法。这些 `Logger` 方法是 `With`、`WithGroup` 和输出方法。

输出方法履行记录器的主要作用：生成日志输出。下面是对输出方法的调用：

```go
logger.Info("hello", "key", value)
```

有两种常规输出方法：`Log` 和 `LogAttrs`。为方便起见，四个常见级别（`Debug`、`Info`、`Warn` 和 `Error`）中的每一个都有一个输出方法，以及采用上下文的相应方法（`DebugContext`、`InfoContext`、`WarnContext` 和 `ErrorContext`）。

每个 `Logger` 输出方法首先调用其处理程序的 `Enabled` 方法。如果该调用返回 true，则该方法从其参数构造 `Record`，并调用处理程序的 `Handle` 方法。

为了方便和优化，可以通过调用 `With` 方法将属性添加到 `Logger`：

```go
logger = logger.With("k", v)
```

此调用创建一个新的带参数属性的 `Logger` 值;原来的`logger`保持不变。后续所有从新`Logger`的输出都将包含这些属性。一个记录器（logger）的 `With` 方法调用其处理程序的 `WithAttrs` 方法。

`WithGroup` 方法用于通过建立单独的命名空间来避免大型程序中的键冲突。此调用创建一个新的 `Logger` 值，其中包含一个名为“g”的组：

```go
logger = logger.WithGroup("g")
```

`logger` 的所有后续键将由组名“g”限定。“限定”的确切含义取决于记录器的处理程序如何格式化输出。内置的 `TextHandler` 将组视为键的前缀，用点分隔：例如，键 `k` 变成 `g.k`。内置的 `JSONHandler` 使用该组作为嵌套 JSON 对象的键：

```go
{"g": {"k": v}}
```

记录器的 `WithGroup` 方法调用其处理程序的 `WithGroup` 方法。

# 实现 Handler 的方法

现在，我们可以详细讨论 `Handler` 的四种方法了。在此过程中，我们将编写一个处理程序，该处理程序使用类似于 YAML 的格式来格式化日志。它将显示以下日志输出调用：

```go
logger.Info("hello", "key", 23)
```

记录样式：

```yaml
time: 2023-05-15T16:29:00
level: INFO
message: "hello"
key: 23
---
```

尽管此特定输出是有效的 YAML，但我们的实现没有考虑 YAML 语法的微妙之处，因此有时会生成无效的 YAML。例如，它不会引用包含冒号的键。我们将其称为 `IndentHandler`。

我们从`IndentHandler`类型和从`io.Writer`构造`New`函数及选项 options 开始：

```go
type IndentHandler struct {
	opts Options
	// TODO: state for WithGroup and WithAttrs
	mu  *sync.Mutex
	out io.Writer
}

type Options struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler
}

func New(out io.Writer, opts *Options) *IndentHandler {
	h := &IndentHandler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}
```

我们只支持一个选项，即设置最低级别以抑制详细日志输出的能力。处理程序应始终将此选项声明为`slog.Leveler`。`slog.Leveler`接口由 `Level` 和 `LevelVar` 实现。用户很容易提供 `Level` 值，但更改多个处理程序的级别需要跟踪所有处理程序。如果用户改为传递 `LevelVar`，则对该 `LevelVar` 的单个更改将更改包含它的所有处理程序的行为。对 `LevelVar` 的更改是 goroutine 安全的。

还可以考虑将 `ReplaceAttr` 选项添加到处理程序中，例如内置处理程序的选项(https://pkg.go.dev/log/slog#HandlerOptions.ReplaceAttr)。尽管 `ReplaceAttr` 会使实现复杂化，但它也会使处理程序更普遍有用。

互斥锁将用于确保写入 `io.Writer` 以原子方式发生。不同的是，`IndentHandler` 保存指向`sync.Mutex`的指针而不是值。这是有充分理由的，我们将在后面解释。

## `Enabled`方法

`Enabled` 方法是一种优化，可以避免不必要的工作。`Logger` 输出方法将在处理其任何参数之前调用 `Enabled`，以判断是否应继续。

方法签名：

```go
Enabled(context.Context, Level) bool
```

上下文可用于允许基于上下文信息做出决策。例如，自定义 HTTP 请求标头可以指定最低级别，服务器将添加该上下文用于处理请求。处理程序的 `Enabled` 方法可以报告参数级别是否大于或等于上下文的值值，从而可以独立的控制每个请求工作的完成。

我们的 `IndentHandler` 不使用上下文。它只是将参数级别与其配置的最低级别进行比较：

```go
func (h *IndentHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}
```

## `Handle`方法

向 `Handle` 方法传递一个 `Record`，其中包含对 `Logger` 输出方法的单次调用要记录的所有信息。`Handle` 方法应以某种方式处理它。一种方法是以某种格式输出 `Record`，就像 `TextHandler` 和 `JSONHandler` 所做的那样。但其他选项是修改 `Record` 并将其传递给另一个处理程序，将 `Record` 排入队列以供以后处理，或者忽略它。

`Handle`的函数签名：

```go
Handle(context.Context, Record) error
```

提供上下文是为了支持应用可以沿调用链记录日志信息。与通常的 Go 做法不同，`Handle` 方法不应将取消的上下文视为停止工作的信号。

如果 `Handle` 处理`Record`，则应遵循[文档](https://pkg.go.dev/log/slog#Handler.Handle)中的规则。例如，应忽略零时间字段，也应忽略零属性。

`Handle` 方法生成输出应执行以下步骤：

1. 分配一个缓冲区（通常为 []byte）来保存输出。最好先在内存中构造输出，然后通过单次调用`io.Writer.Write`来写入它。以尽量减少与使用相同 writer 的其他 goroutine 的冲突。
2. 格式化特定字段：时间(time)、级别(level)、消息(message)和源码位置 （PC）。作为一般规则，这些字段应首先显示，并且不会嵌套在 `WithGroup` 建立的组中。
3. 格式化调用 `WithGroup` 和 `WithAttrs` 的结果。
4. 格式化`Record`中的属性。
5. 输出缓冲区。

这就是 IndentHandler.Handle 的结构：

```go
func (h *IndentHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)), 0)
	}
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)
	indentLevel := 0
	// TODO: output the Attrs and groups from WithAttrs and WithGroup.
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a, indentLevel)
		return true
	})
	buf = append(buf, "---\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}
```

第一行分配一个 `[]byte`，对于大部分日志输出应该足够大。为缓冲区分配一些初始的、相当大的容量是一个简单但重要的优化：它避免了在初始切片为空或较小时发生的重复复制和分配。我们将在速度(speed)章节回到这一行，并展示我们如何做得更好。

`Handle` 方法的下一部分格式化指定的属性，遵守忽略零时间和零 PC 的规则。

接下来，该方法处理 `WithAttrs` 和 `WithGroup` 调用的结果，我们暂时跳过它。

然后，是时候处理参数记录中的属性了。我们使用 `Record.Attrs` 方法按照用户将属性传递给 `Logger` 输出方法的顺序循环访问属性。处理程序可以自由地对属性进行重新排序或删除重复数据，但我们的处理程序没这么做。

最后，在将行“---”添加到输出以分隔日志记录后，我们的处理程序对缓冲区累积的数据进行一次`h.out.Write`调用。我们持有此写入的锁，以使其相对于可能同时调用 `Handle` 的其他 goroutine 是原子的。

处理程序的核心是 `appendAttr` 方法，它负责格式化单个属性：

```go
func (h *IndentHandler) appendAttr(buf []byte, a slog.Attr, indentLevel int) []byte {
	// Resolve the Attr's value before doing anything else.
	a.Value = a.Value.Resolve()
	// Ignore empty Attrs.
	if a.Equal(slog.Attr{}) {
		return buf
	}
	// Indent 4 spaces per level.
	buf = fmt.Appendf(buf, "%*s", indentLevel*4, "")
	switch a.Value.Kind() {
	case slog.KindString:
		// Quote string values, to make them easy to parse.
		buf = fmt.Appendf(buf, "%s: %q\n", a.Key, a.Value.String())
	case slog.KindTime:
		// Write times in a standard way, without the monotonic time.
		buf = fmt.Appendf(buf, "%s: %s\n", a.Key, a.Value.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		attrs := a.Value.Group()
		// Ignore empty groups.
		if len(attrs) == 0 {
			return buf
		}
		// If the key is non-empty, write it out and indent the rest of the attrs.
		// Otherwise, inline the attrs.
		if a.Key != "" {
			buf = fmt.Appendf(buf, "%s:\n", a.Key)
			indentLevel++
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga, indentLevel)
		}
	default:
		buf = fmt.Appendf(buf, "%s: %s\n", a.Key, a.Value)
	}
	return buf
}
```

它首先解析属性，运行值的 `LogValuer.LogValue` 方法（如果有）。所有处理程序都应解析它们处理的每个属性。

接下来，它遵循处理程序规则，该规则规定应忽略空属性。

然后，用类型断言确定属性要使用的格式。对于大多数类型（switch 的默认 case），它依赖于 `slog.Value` 的 `String` 方法来生成合理的东西。处理字符串和时间：通过引用字符串来处理字符串，通过以标准方式格式化字符串来处理时间。

当 `appendAttr` 看到一个 `Group` 时，它会在应用另外两个处理程序规则后，对该组的属性进行递归调用。首先，忽略没有属性的组，甚至不显示其键。其次，具有空键的组是内联的：组边界不以任何方式标记。在我们的例子中，这意味着组的属性不会缩进。

## `WithAttrs`方法

`slog` 的性能优化之一是支持预格式化 属性。`Logger.With` 方法将键值对转换为 `Attrs` 和 然后调用 `Handler.WithAttrs`。 处理程序存储属性以供后面 `Handle` 方法使用， 或者现在格式化属性，一次， 而不是在每次调用 `Handle` 时重复这样做。

`WithAttrs`的方法签名：

```go
WithAttrs(attrs []Attr) Handler
```

参数属性是传递给`Logger.With`处理过的键值对。返回值应是`handler`的新实例，其中包含 属性，可能的预格式化。

`WithAttrs` 必须返回一个具有附加属性的新`handler`，使 原始处理程序（其接收器 receiver）保持不变。例如，以下调用：

```go
logger2 := logger1.With("k", v)
```

创建一个具有附加属性的新记录器 `logger2`，但是对`logger1`没有影响.

下面讨论`WithGroup`的时候将会展示`WithAttrs`的实现例子。

## `WithGroup`方法

`Logger.WithGroup` 直接调用 `Handler.WithGroup`，具有相同的参数和组名，一个`handler`应记住组名，以便可以使用它来限定所有后续的属性。

`WithGroup`的方法签名：

```go
WithGroup(name string) Handler
```

和`WithAttrs` 一样，`WithGroup` 方法应返回一个新的`handler`，而不修改接收器(receiver)

`WithGroup`和 `WithAttrs` 的实现是相互交织的。 请看以下语句：

```go
logger = logger.WithGroup("g1").With("k1", 1).WithGroup("g2").With("k2", 2)
```

`logger`的输出应使用组“g1”来限定键“k1”， 键“K2”用组组“g1”和“g2”。`Logger.WithGroup` 和 `Logger.With` 调用的顺序必须遵循 `Handler.WithGroup` 和 `Handler.WithAttrs` 的实现。

我们将研究 `WithGroup` 和 `WithAttrs` 的两种实现，一种是预格式化和一个没有。

### 无预格式化

我们的第一个实现将从 `WithGroup` 和 `WithAttrs` 调用以构建包含组名和属性列表的切片， 并在 `Handle` 中循环使用该切片。我们从一个包含有组名称和一些属性的结构体开始：

```go
// groupOrAttrs holds either a group name or a list of slog.Attrs.
type groupOrAttrs struct {
	group string      // group name if non-empty
	attrs []slog.Attr // attrs if non-empty
}
```

然后,我们将 `groupOrAttrs` 的一部分添加到`handler`中:

```go
type IndentHandler struct {
	opts Options
	goas []groupOrAttrs
	mu   *sync.Mutex
	out  io.Writer
}
```

如上所述，`WithGroup` 和`WithAttrs` 方法不应修改其 接收器（receiver）。 为此，我们定义了一个方法，该方法将复制我们的处理程序结构 并将一个 `groupOrAttrs` 附加到副本中：

```go
func (h *IndentHandler) withGroupOrAttrs(goa groupOrAttrs) *IndentHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}
```

`IndentHandler` 的大部分字段都可以浅层复制，但 `groupOrAttrs` 需要深拷贝，否则克隆和原始副本将指向 相同的基础数组。如果我们使用 `append` 而不是显式 复制，我们会引入那个微妙的混叠错误。

使用 `withGroupOrAttrs`, `With` 方法非常简单:

```go
func (h *IndentHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *IndentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}
```

`Handle` 方法现在可以在内置属性之后和记录中的属性之前处理 `groupOrAttrs` 切片：

```go
func (h *IndentHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)), 0)
	}
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)
	indentLevel := 0
	// Handle state from WithGroup and WithAttrs.
	goas := h.goas
	if r.NumAttrs() == 0 {
		// If the record has no Attrs, remove groups at the end of the list; they are empty.
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}
	for _, goa := range goas {
		if goa.group != "" {
			buf = fmt.Appendf(buf, "%*s%s:\n", indentLevel*4, "", goa.group)
			indentLevel++
		} else {
			for _, a := range goa.attrs {
				buf = h.appendAttr(buf, a, indentLevel)
			}
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a, indentLevel)
		return true
	})
	buf = append(buf, "---\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}
```

您可能已经注意到，因为重复复制，我们记录`WithGroup`和`WithAttrs`信息的算法在调用这些方法的次数上是二次方的。

这在实践中不太重要，但如果它困扰您，您可以使用链表代替，`Handle`必须反转或递归访问该链表，有关实现请参阅https://github.com/jba/slog/tree/main/withsupport

#### 正确使用互斥锁

让我们在看看`Handle`的最后几行：

```go
h.mu.Lock()
defer h.mu.Unlock()
_, err := h.out.Write(buf)
return err
```

这段代码没有任何变化，但是我们可以领会一下为什么`h.mu`是指向`sync.Mutex`的指针。`WithGroup`和`WithAttrs`都要复制`handler`，所有副本指针都指向同一个`mutex`。

假如副本对象和原始对象使用不同的互斥锁对象，并且同时使用，那么输出可能会交错或者丢失一些内容。像这样的代码：

```go
l2 := l1.With("a", 1)
go l1.Info("hello")
l2.Info("goodbye")
```

可能会产生这样的输出：

```shell
hegoollo a=dbye1
```

### 使用预格式化

我们的第二个实现实现了预格式化。此实现比前一个更复杂。额外的复杂性值得吗？这取决于您的情况，但这里有一种情况它可能需要。

假设您希望您的服务器在请求期间发生的每条日志消息中记录有关入站请求的大量信息，典型的处理程序可能如下所示：

```go
func (s *Server) handleWidgets(w http.ResponseWriter, r *http.Request) {
    logger := s.logger.With(
        "url", r.URL,
        "traceID": r.Header.Get("X-Cloud-Trace-Context"),
        // many other attributes
        )
    // ...
}
```

一个`handleWidgets`可能会产线是数百行日志，例如,它可能包含这样的代码：

```go
for _, w := range widgets {
    logger.Info("processing widget", "name", w.Name)
    // ...
}
```

对于每一行，除了日志行本身上的属性之外，我们上面编写的 Handle 方法将格式化使用上面的 With 添加的所有属性。

所有这些额外的工作都不会显著降低您的服务器速度，因为它做了如此多的其他工作，以至于花费在日志记录上的时间只是噪音。但也许您的服务器足够快，以至于所有额外的格式化都出现在 CPU profile 中的顶部。也就是说，预格式化可以产生很大的影响，只需在调用一次 `With` 格式化属性。

为了把预格式化参数给`WithAttrs`，我们必须在`IndentHandler`结构体里面保持跟踪一些额外的状态。

```go
type IndentHandler struct {
	opts           Options
	preformatted   []byte   // data from WithGroup and WithAttrs
	unopenedGroups []string // groups from WithGroup that haven't been opened
	indentLevel    int      // same as number of opened groups so far
	mu             *sync.Mutex
	out            io.Writer
}
```

主要是，我们需要一个缓冲区来保存预格式化的数据。但是我们还需要跟踪哪些组我们已经看到了，但还没有输出。我们将这些组称为“unopened”， 我们还需要跟踪多少个群组已经打开，我们可以使用一个简单的计数器来实现，因为打开组的唯一效果是更改缩进级别。

WitGroup 实现与前一个非常相似：只需要记住初始 unopened 的新组：

```go
func (h *IndentHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	// Add an unopened group to h2 without modifying h.
	h2.unopenedGroups = make([]string, len(h.unopenedGroups)+1)
	copy(h2.unopenedGroups, h.unopenedGroups)
	h2.unopenedGroups[len(h2.unopenedGroups)-1] = name
	return &h2
}
```

`WithAttrs`完成所有预格式化：

```go
func (h *IndentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	// Force an append to copy the underlying array.
	pre := slices.Clip(h.preformatted)
	// Add all groups from WithGroup that haven't already been added.
	h2.preformatted = h2.appendUnopenedGroups(pre, h2.indentLevel)
	// Each of those groups increased the indent level by 1.
	h2.indentLevel += len(h2.unopenedGroups)
	// Now all groups have been opened.
	h2.unopenedGroups = nil
	// Pre-format the attributes.
	for _, a := range attrs {
		h2.preformatted = h2.appendAttr(h2.preformatted, a, h2.indentLevel)
	}
	return &h2
}

func (h *IndentHandler) appendUnopenedGroups(buf []byte, indentLevel int) []byte {
	for _, g := range h.unopenedGroups {
		buf = fmt.Appendf(buf, "%*s%s:\n", indentLevel*4, "", g)
		indentLevel++
	}
	return buf
}
```

它首先打开任何未打开(unopened)的组。这处理调用诸如:

```go
logger.WithGroup("g").WithGroup("h").With("a", 1)
```

此处,WithAttrs 必须在 "a" 之前输出 "g" 和 "h"。由于 ` WithGroup` 建立的组对日志行的其余部分有效,因此 `WithAttrs` 会为打开的每个组增加缩进级别。

最后,`WithAttrs`使用我们上面看到的相同的`appendAttr` 方法格式化其参数属性。

`Handle` 方法的工作是将预先格式化的材料插入到正确的位置，即在内置属性之后和记录中的属性之前：

```go
func (h *IndentHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)), 0)
	}
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)
	// Insert preformatted attributes just after built-in ones.
	buf = append(buf, h.preformatted...)
	if r.NumAttrs() > 0 {
		buf = h.appendUnopenedGroups(buf, h.indentLevel)
		r.Attrs(func(a slog.Attr) bool {
			buf = h.appendAttr(buf, a, h.indentLevel+len(h.unopenedGroups))
			return true
		})
	}
	buf = append(buf, "---\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}
```

它还必须打开尚未打开的任何组。逻辑涵盖的日志行如下那样：

```go
logger.WithGroup("g").Info("msg", "a", 1)
```

`g`在`handle`调用之前未打开，并且必须写入生成正确的输出。

```shell
level: INFO
msg: "msg"
g:
    a: 1
```

在这种情况下检查`r.NumAttrs() > 0`：

```go
logger.WithGroup("g").Info("msg")
```

这里没有记录属性，所以没有要打开的组。

## 测试

`Handler` 指定了对处理程序的多个约束。若要验证处理程序是否遵循这些规则并通常生成正确的输出，请使用 `testing/slogtest` 包。

包`TestHandler`的功能是接受一个handler的实例，和一个返回输出格式化成一个slice map的函数。一下是一个测试函数的例子：

```go
func TestSlogtest(t *testing.T) {
	var buf bytes.Buffer
	err := slogtest.TestHandler(New(&buf, nil), func() []map[string]any {
		return parseLogEntries(t, buf.Bytes())
	})
	if err != nil {
		t.Error(err)
	}
}
```

调用 TestHandler 非常简单。困难的部分是解析`handler`的输出。TestHandler 多次调用`handler`，从而生成一系列日志条目。您的工作是将每个条目解析到 `map[string]any` 中。条目中的组应显示为嵌套map。

如果`handler`输出标准格式，则可以使用现有解析器。例如，如果处理程序每行输出一个 JSON 对象，则可以将输出拆分为几行并调用 `encoding/json.Unmarshal`。其他格式的解析器解析成map就可以开箱即用。我们的示例输出与 YAML 足够像，因此我们可以使用 `gopkg.in/yaml.v3` 包来解析它：

```go
func parseLogEntries(t *testing.T, data []byte) []map[string]any {
	entries := bytes.Split(data, []byte("---\n"))
	entries = entries[:len(entries)-1] // last one is empty
	var ms []map[string]any
	for _, e := range entries {
		var m map[string]any
		if err := yaml.Unmarshal([]byte(e), &m); err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}
	return ms
}
```

如果你必须编写自己的解析器，它可能远非完美。`slogtest`包只使用了几个简单的属性。(它是在测试处理程序的一致性，而不是解析。)你的解析器可以忽略键和值中的空白符和换行符等边缘情况。在切换到YAML解析器之前，我们写了一个65行的自定义解析器。

# 注意事项

## 复制记录

大多数`handler`不需要复制传给`Handle`方法的`slog.Record`。在某些情况下这样做的必须特别小心。

如果一个`handler`不保留原始记录，则可以使用go赋值、channel发送或者函数调用来复制单个`Record`。但是如果操作结果是多个副本，则必须使用`Record.Clone`来创建副本，以便他们不会共享状态。

此 `Handle` 方法将记录传递给单个处理程序，因此它不需要`Clone`：


```go
type Handler1 struct {
    h slog.Handler
    // ...
}

func (h *Handler1) Handle(ctx context.Context, r slog.Record) error {
    return h.h.Handle(ctx, r)
}
```

此`Handle` 方法可能会将记录传递给多个处理程序，因此它应使用 `Clone`：

```go
type Handler2 struct {
    hs []slog.Handler
    // ...
}

func (h *Handler2) Handle(ctx context.Context, r slog.Record) error {
    for _, hh := range h.hs {
        if err := hh.Handle(ctx, r.Clone()); err != nil {
            return err
        }
    }
    return nil
}
```

## 并发安全

当单个 Logger 在多个 goroutine 之间共享时，处理程序必须正常工作。这意味着必须使用锁或其他机制来保护可变状态。在实践中，这并不难实现，因为许多处理程序没有任何可变状态。

- `Enabled` 方法通常仅查询其参数和配置的级别。该级别通常最初设置一次，或者保存在 `LevelVar` 中，这已经是并发安全的。
- 由于上述原因，`WithAttrs` 和 `WithGroup` 方法不应修改接收方(receiver)。
- `Handle` 方法通常仅适用于其参数和存储字段。

对如`io.Writer.Write`等输出方法的调用，除非可以验证不需要锁定，否则应同步。正如我们在示例中看到的，存储指向互斥锁的指针使记录器及其所有克隆能够相互同步。谨防诸如“Unix 写入是原子的”之类的简单声明;情况比这要微妙得多。

一些处理程序有正当理由保持状态。例如，处理程序可能支持 `SetLevel` 方法动态更改其配置的级别。或者，它可能会输出对 `Handle` 的成功调用之间的时间，这需要一个可变字段来保存上次输出时间。同步对此类字段的所有访问，包括读取和写入。

内置处理程序没有直接可变状态。它们仅使用互斥锁对其包含的 `io.Writer` 进行排序调用。

## 健壮性

日志记录通常是最后的调试技术。当难以或不可能检查系统时（例如生产服务器通常的情况），日志提供了了解其行为的最详细方法。因此，处理程序应该对错误的输入具有健壮性。

例如，当函数发现问题（如无效参数）时，通常的建议是恐慌或返回错误。内置处理程序不遵循该建议。没有什么比无法调试导致日志记录失败的问题更令人沮丧的了;它最后是产生一些输出，无论多么不完美，总比什么都不生产要好。这就是为什么像 `Logger.Info` 这样的方法将键值对列表中的编程错误（如缺失值或格式错误的键）转换为包含尽可能多的信息的 `Attrs`。

避免恐慌(panic)的一个地方是处理属性值。想要设置值格式的处理程序会判断值的类型：

```go
switch attr.Value.Kind() {
case KindString: ...
case KindTime: ...
// all other Kinds
default: ...
}
```

在默认情况下，当处理程序遇到它不知道的 Kind 时，会发生什么？内置处理程序尝试使用值的 String 方法的结果进行混淆，就像我们的示例处理程序一样。他们不会惊慌失措或返回错误。此外，您自己的处理程序可能还希望通过生产监视或错误跟踪遥测系统报告问题。对这个问题最可能的解释是，较新版本的 slog 包添加了一个新的 Kind（Go 1 兼容性承诺下的向后兼容更改），并且处理程序未更新。这当然是一个问题，但它不应该剥夺读者看到日志输出的其余部分。

在一种情况下，从 Handler.Handle 返回错误是合适的。如果输出操作本身失败，最佳操作方法是通过返回错误来报告此失败。例如，内置 Handle 方法的最后两行是

```go
_, err := h.w.Write(*state.buf)
return err
```

尽管 `Logger` 的输出方法忽略了该错误，但可以编写一个处理程序来执行某些操作，可能会回退到写入标准错误。

## 性能

大多数程序不需要快速日志记录。在优化处理程序之前，请收集数据（最好是生产数据，而不是基准比较），以证明它需要快速。避免过早优化。

如果您需要快速处理程序，请从预格式化开始。如果对 `Logger.With` 的单个调用之后对生成的记录器进行多次调用，它可能会提供显着的加速。

如果日志输出是瓶颈，请考虑将处理程序设置为异步。在处理程序中执行最少的处理量，然后通过通道发送记录和其他信息。另一个 goroutine 可以收集传入的日志条目，并在后台批量写入它们。您可能希望保留同步日志选项，以便可以看到调试崩溃的所有日志输出。

分配通常是系统缓慢的主要原因。slog 包已经在努力最小化分配。如果处理程序执行自己的分配，并且分析显示它是一个问题，那么看看是否可以最小化它。

您可以进行的一项简单更改是替换对`fmt.Sprintf`或`fmt.Appendf`的调用直接追加到缓冲区。例如，我们的 `IndentHandler` 将字符串属性追加到缓冲区，如下所示：

```go
buf = fmt.Appendf(buf, "%s: %q\n", a.Key, a.Value.String())
```

从 Go 1.21 开始，这会产生两个分配，每个参数传递给 any 参数一个。我们可以通过直接使用 append 将其降至零：

```go
buf = append(buf, a.Key...)
buf = append(buf, ": "...)
buf = strconv.AppendQuote(buf, a.Value.String())
buf = append(buf, '\n')
```

另一个值得注意的更改是使用`sync.Pool`。用于管理大多数处理程序所需的一个内存块的池：保存格式化输出的`[]byte`缓冲区。

我们的示例 Handle 方法从以下行开始：

```go
buf := make([]byte, 0, 1024)
```

如上所述，提供较大的初始容量可以避免在切片增长时重复复制和重新分配切片，从而将分配次数减少到一个。但是，在稳定状态下，我们可以通过保留一个全局缓冲池将其降至零。最初，池将是空的，并将分配新的缓冲区。但最终，假设并发日志调用数达到稳定的最大值，池中将有足够的缓冲区在所有正在进行的处理程序调用之间共享。只要没有日志条目超过缓冲区的容量，从垃圾回收器的角度来看，就不会有分配。

我们将把我们的池隐藏在一对函数后面，`allocBuf` 和 `freeBuf`。在`Handle`顶部获取缓冲区单行会变为两行：

```go
bufp := allocBuf()
defer freeBuf(bufp)
```

其中一个点是创建`sync.Pool`slice建议通过变量名`bufp`：池必须处理指向切片的指针，而不是切片本身。池值必须始终为指针。如果不是，则`sync.Pool`有`any` 参数和返回值的方法本身会导致分配，从而破坏池化的目的。

有两种方法可以处理我们的切片指针：我们可以在整个函数中用 *bufp 替换 buf，或者我们可以取消引用它并记住在释放之前重新分配它：我们可以在整个函数中用 *bufp 替换 buf，或者我们可以取消引用它并记住在释放之前重新分配它：

```go
bufp := allocBuf()
buf := *bufp
defer func() {
	*bufp = buf
	freeBuf(bufp)
}()
```

以下是我们的`pool`及其相关函数：

```go
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

func allocBuf() *[]byte {
	return bufPool.Get().(*[]byte)
}

func freeBuf(b *[]byte) {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}
```
池的 `New` 函数执行与原始代码相同的操作：创建一个长度为 0 且容量充足的字节切片。`allocBuf` 函数只是对池的 `Get` 方法的结果进行类型断言。

`freeBuf` 方法在将缓冲区放回池之前会截断缓冲区，以便 `allocBuf` 始终返回长度为零的切片。它还实现了一项重要的优化：它不会将大型缓冲区返回到池中。要了解为什么这很重要，请考虑如果有一个异常大的日志条目（例如格式化时为兆字节的日志条目）会发生什么。如果将兆字节大小的缓冲区放入池中，它可以无限期地保留在那里，不断被重复使用，但其大部分容量被浪费了。额外的内存可能永远不会被处理程序再次使用，并且由于它位于处理程序的池中，因此它可能永远不会被垃圾回收以在其他地方重用。我们可以通过从池中排除大型缓冲区来避免这种情况。