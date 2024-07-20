package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"sync"

	"github.com/fatih/color"
)

// https://colobu.com/2024/03/10/slog-the-ultimate-guide/

// 颜色(github.com/fatih/color)，错误堆栈(github.com/MDobak/go-xerrors)，隐藏敏感字段

const timeFormat = "2006-01-02 15:04:05.000"

type MyHandler struct {
	opts Options
	mu   *sync.Mutex
	out  io.Writer
	file *os.File
	// json已格式化属性
	jsonPreformated []byte
	// 控制台已格式化属性
	consolePreFormated []byte
	// 已格式化group
	groupPreFormated []byte
	// group列表
	groups []string
}

type Options struct {
	slog.HandlerOptions
	Trace         bool
	Color         bool
	OutputFile    string
	OutputConsole bool
}

func NewMyHandler(opts *Options) *MyHandler {
	h := &MyHandler{mu: &sync.Mutex{}, out: os.Stdout}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.OutputFile != "" {
		file, err := os.OpenFile(h.opts.OutputFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			panic(err)
		}
		h.file = file
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelDebug
	}
	return h
}

// Enabled 判断是否需要记录日志
func (h *MyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle 处理每一条日志
// 比如记录到文件，转发到其他地方
func (h *MyHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.file != nil {
		// 输出到文件
		if err := h.handleJson(ctx, r); err != nil {
			log.Println(err)
		}
	}
	if !h.opts.OutputConsole {
		return nil
	}
	// 输出到控制台
	buf := make([]byte, 0, 1024)

	// 输出内置属性
	if !r.Time.IsZero() {
		buf = fmt.Appendf(buf, "%s%2s", r.Time.Format(timeFormat), "")
	}
	// %*s语法：*的值=7-len(r.Level.String())的计算结果，必须是整数
	buf = fmt.Appendf(buf, "%s%*s", h.getLevelColor(r.Level), 7-len(r.Level.String()), "")
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = fmt.Appendf(buf, "%s%2s", fmt.Sprintf("%s:%d", f.File, f.Line), "")
	}
	// 日志消息
	buf = fmt.Appendf(buf, "%s%2s", r.Message, "")
	buf = append(buf, h.consolePreFormated...)
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttrConsole(buf, a)
		return true
	})
	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

// WithAttrs 使用现有处理程序创建一个新处理程序，并向其中添加指定的属性 (attrs)
// 添加公共属性
func (h *MyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	ch := *h
	for _, a := range attrs {
		if h.opts.OutputConsole {
			ch.consolePreFormated = ch.appendAttrConsole(ch.consolePreFormated, a)
		}
		if h.opts.OutputFile != "" {
			for _, g := range h.groups {
				ch.jsonPreformated = fmt.Appendf(ch.jsonPreformated, `"%s":{`, g)
			}
			ch.jsonPreformated = ch.appendAttrJson(ch.jsonPreformated, a)
		}
	}
	return &ch
}

// WithGroup 使用现有处理程序创建一个新处理程序，并向其中添加指定的组名 (group)，该名称限定后续的属性
// 添加分组
func (h *MyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	ch := *h
	ch.groups = append(ch.groups, name)
	ch.groupPreFormated = fmt.Appendf(ch.groupPreFormated, "%s.", name)
	return &ch
}

// handleJson 输出json格式
func (h *MyHandler) handleJson(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	buf = append(buf, '{')
	if !r.Time.IsZero() {
		buf = h.appendAttrJson(buf, slog.Time(slog.TimeKey, r.Time))
	}
	buf = h.appendAttrJson(buf, slog.Any(slog.LevelKey, r.Level))
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttrJson(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)))
	}
	buf = h.appendAttrJson(buf, slog.String(slog.MessageKey, r.Message))
	buf = append(buf, h.jsonPreformated...)
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttrJson(buf, a)
		return true
	})
	buf = buf[0 : len(buf)-1]
	if len(h.groups) > 0 {
		buf = fmt.Appendf(buf, "%*s", len(h.groups), "}")
	}
	buf = fmt.Appendf(buf, "}\n")
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.file.Write(buf)
	return err
}

// appendAttrJson 格式化日志项的单个属性
func (h *MyHandler) appendAttrJson(buf []byte, a slog.Attr) []byte {
	// 执行任何操作前先解析值
	a.Value = a.Value.Resolve()
	// 忽略空属性
	if a.Equal(slog.Attr{}) {
		return buf
	}
	// 用类型断言确定属性要使用的格式
	switch a.Value.Kind() {
	case slog.KindTime:
		buf = fmt.Appendf(buf, `"%s":"%s"`, a.Key, a.Value.Time().Format(timeFormat))
	case slog.KindGroup:
		// todo 处理 { 和 }
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			buf = fmt.Appendf(buf, `"%s":`, a.Key)
		}
		for _, ga := range attrs {
			buf = h.appendAttrJson(buf, ga)
		}
	default:
		buf = fmt.Appendf(buf, `"%s":"%s"`, a.Key, a.Value)
	}
	buf = append(buf, ',')
	return buf
}

// appendAttrConsole 格式化输出到控制台的属性
func (h *MyHandler) appendAttrConsole(buf []byte, a slog.Attr) []byte {
	// 执行任何操作前先解析值
	a.Value = a.Value.Resolve()
	// 忽略空属性
	if a.Equal(slog.Attr{}) {
		return buf
	}
	// 所有属性之前先加上group
	if a.Value.Kind() != slog.KindGroup {
		buf = append(buf, h.groupPreFormated...)
	}
	// 用类型断言确定属性要使用的格式
	switch a.Value.Kind() {
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s: %s", a.Key, a.Value.Time().Format(timeFormat))
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			buf = append(buf, h.groupPreFormated...)
			buf = fmt.Appendf(buf, "%s.", a.Key)
		}
		for _, ga := range attrs {
			buf = fmt.Appendf(buf, "%s: %s", ga.Key, ga.Value)
		}
	default:
		buf = fmt.Appendf(buf, "%s: %s", a.Key, a.Value)
	}
	buf = fmt.Appendf(buf, "%2s", "")
	return buf
}

func (h *MyHandler) getLevelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return color.MagentaString(level.String())
	case slog.LevelInfo:
		return color.GreenString(level.String())
	case slog.LevelWarn:
		return color.YellowString(level.String())
	case slog.LevelError:
		return color.RedString(level.String())
	default:
		return level.String()
	}
}
