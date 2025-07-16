# Go Package Learning Project Specification

## 项目概述

**项目名称**: gopkg  
**项目类型**: Go 语言学习和实验代码库  
**目标**: 提供 Go 语言各种包和功能的代码示例和使用方法，特别专注于完整的设计模式实现集合

## 项目特色

### 🎯 完整的设计模式实现

- **23 种 GoF 设计模式** - 业界最完整的 Go 语言设计模式实现
- **实用场景示例** - 每个模式都包含真实应用场景
- **中文详细注释** - 便于中文开发者学习理解
- **可直接运行** - 所有示例都包含完整的 main 函数

### 📚 系统化学习资源

- **分类清晰** - 按创建型、结构型、行为型模式分类
- **渐进式学习** - 从简单到复杂的学习路径
- **最佳实践** - 遵循 Go 语言惯用法和最佳实践

## 项目结构

### 核心模块

#### 1. 设计模式 (designpattern/) - 项目核心亮点

完整实现 23 种 GoF 设计模式的 Go 语言版本：

##### 创建型模式 (Creational Patterns)

- **抽象工厂模式** (abstractfactory) - 创建相关对象族
- **建造者模式** (builder) - 复杂对象分步构建
- **工厂方法模式** (factory) - 对象创建接口抽象
- **原型模式** (prototype) - 对象克隆复制
- **单例模式** (singleton) - 唯一实例保证

##### 结构型模式 (Structural Patterns)

- **适配器模式** (adapter) - 接口适配和转换
- **桥接模式** (bridge) - 抽象与实现分离
- **组合模式** (composite) - 树形结构统一处理
- **装饰器模式** (decorator) - 动态功能扩展
- **外观模式** (facade) - 简化复杂子系统接口
- **享元模式** (flyweight) - 对象共享内存优化
- **代理模式** (proxy) - 访问控制和延迟加载

##### 行为型模式 (Behavioral Patterns)

- **责任链模式** (chain) - 请求处理链传递
- **命令模式** (command) - 请求封装为对象
- **迭代器模式** (iterator) - 集合元素顺序访问
- **中介者模式** (mediator) - 对象交互中介协调
- **备忘录模式** (memento) - 对象状态保存恢复
- **观察者模式** (observer) - 一对多依赖通知
- **状态模式** (state) - 状态改变行为切换
- **策略模式** (strategy) - 算法族封装互换
- **模板方法模式** (template) - 算法骨架定义
- **访问者模式** (visitor) - 操作与对象结构分离

#### 2. 网络编程 (net/, grpc/)

- HTTP 客户端和服务器实现
- gRPC 服务开发
- 网络协议处理

#### 3. 加密安全 (crypto/)

- **AES 加密** - 对称加密实现
- **RSA 加密** - 非对称加密
- **密码生成** (genpwd) - 安全密码生成工具

#### 4. 数据处理

- **bufio** - 缓冲 I/O 操作
- **bytes** - 字节处理
- **image** - 图像处理

#### 5. 爬虫和自动化 (crawler/, chromedp/)

- **ChromeDP** - 浏览器自动化
- **CSGO 相关爬虫** - 游戏数据抓取

#### 6. 监控和追踪

- **OpenTelemetry** - 分布式追踪
- **expvar** - 运行时变量暴露
- **log** - 日志处理

#### 7. 数据库和 ORM (orm/)

- 数据库操作抽象
- ORM 模式实现

#### 8. 其他工具模块

- **OTP** - 一次性密码
- **Plugin** - 插件系统
- **Reflect** - 反射机制
- **Unsafe** - 不安全操作
- **xsync** - 扩展同步原语
- **iter** - 迭代器模式

## 技术栈

### 核心技术

- **Go 1.23** - 主要编程语言
- **标准库** - 充分利用 Go 标准库功能

### 主要依赖

- **chromedp/chromedp** - 浏览器自动化
- **go.opentelemetry.io** - 可观测性
- **google.golang.org/grpc** - RPC 框架
- **go.uber.org/zap** - 高性能日志
- **github.com/fatih/color** - 终端颜色输出

## 项目目标

### 学习目标

1. **Go 语言特性掌握** - 深入理解 Go 语言各种特性
2. **设计模式实践** - 在 Go 中实现经典设计模式
3. **标准库使用** - 熟练使用 Go 标准库
4. **第三方库集成** - 学习主流 Go 生态库的使用

### 实践目标

1. **代码示例库** - 提供可运行的代码示例
2. **最佳实践** - 展示 Go 语言最佳实践
3. **性能优化** - 学习 Go 性能优化技巧
4. **并发编程** - 掌握 Go 并发编程模式

## 开发规范

### 代码结构

- 每个模块独立目录
- 包含 main.go 作为示例入口
- 清晰的包结构和命名

### 文档要求

- 每个模块包含 README 说明
- 代码注释完整
- 使用示例清晰

### 测试要求

- 关键功能包含单元测试
- 性能基准测试
- 示例代码可执行

## 扩展计划

### 短期目标

1. 完善现有模块的文档和示例
2. 添加更多设计模式实现
3. 增加性能测试和基准

### 长期目标

1. 添加微服务相关示例
2. 云原生技术集成
3. 更多实际应用场景

## 使用指南

### 快速开始

```bash
# 克隆项目
git clone <repository-url>

# 进入项目目录
cd gopkg

# 运行设计模式示例
go run designpattern/singleton/main.go
go run designpattern/observer/main.go
go run designpattern/strategy/main.go

# 运行其他模块示例
go run otp/main.go
go run crypto/aes/main.go
```

### 设计模式学习路径

建议按以下顺序学习设计模式，从简单到复杂：

#### 🚀 入门级模式

1. **单例模式** (singleton) - 理解实例控制
2. **工厂方法模式** (factory) - 学习对象创建
3. **适配器模式** (adapter) - 掌握接口转换

#### 🎯 进阶模式

4. **观察者模式** (observer) - 理解发布订阅
5. **策略模式** (strategy) - 学习算法封装
6. **装饰器模式** (decorator) - 掌握功能扩展
7. **代理模式** (proxy) - 理解访问控制

#### 🏗️ 结构型模式

8. **建造者模式** (builder) - 复杂对象构建
9. **组合模式** (composite) - 树形结构处理
10. **外观模式** (facade) - 接口简化
11. **桥接模式** (bridge) - 抽象实现分离

#### 🔄 行为型模式

12. **命令模式** (command) - 请求封装
13. **状态模式** (state) - 状态行为管理
14. **责任链模式** (chain) - 请求链式处理
15. **模板方法模式** (template) - 算法骨架
16. **迭代器模式** (iterator) - 集合遍历

#### 🎨 高级模式

17. **访问者模式** (visitor) - 操作与结构分离
18. **中介者模式** (mediator) - 对象交互协调
19. **备忘录模式** (memento) - 状态保存恢复
20. **抽象工厂模式** (abstractfactory) - 对象族创建
21. **原型模式** (prototype) - 对象克隆
22. **享元模式** (flyweight) - 内存优化

### 模块探索指南

#### 设计模式模块 (推荐优先学习)

```bash
# 查看所有设计模式
ls designpattern/

# 运行具体模式示例
go run designpattern/observer/main.go
go run designpattern/strategy/main.go
```

#### 网络编程模块

```bash
# HTTP 服务示例
go run net/http/main.go

# gRPC 服务示例
go run grpc/server/server.go
go run grpc/client/client.go
```

#### 加密安全模块

```bash
# AES 加密示例
go run crypto/aes/main.go

# RSA 加密示例
go run crypto/rsa/main.go

# 密码生成工具
go run crypto/genpwd/main.go
```

#### 其他实用模块

```bash
# 一次性密码
go run otp/main.go

# 图像处理
go run image/main.go

# 反射机制
go run reflect/main.go
```

## 项目亮点

### 🌟 设计模式完整性

- **业界最全** - 23 种 GoF 设计模式完整实现
- **Go 语言特色** - 充分利用 Go 的接口、组合、并发特性
- **实用导向** - 每个模式都有真实应用场景示例

### 📖 学习友好性

- **中文注释** - 详细的中文代码注释和说明
- **渐进学习** - 从简单到复杂的学习路径
- **即时运行** - 所有示例都可以直接运行查看效果

### 🔧 代码质量

- **最佳实践** - 遵循 Go 语言编码规范和最佳实践
- **错误处理** - 完善的错误处理机制
- **性能考虑** - 注重内存管理和性能优化

### 🚀 实用价值

- **生产就绪** - 代码质量达到生产环境标准
- **可扩展性** - 易于扩展和定制
- **文档完善** - 详细的使用说明和 API 文档

## 贡献指南

### 代码贡献

1. Fork 项目到个人仓库
2. 创建功能分支 (`git checkout -b feature/new-pattern`)
3. 提交更改 (`git commit -am 'Add new design pattern'`)
4. 推送到分支 (`git push origin feature/new-pattern`)
5. 创建 Pull Request

### 文档贡献

- 改进现有文档
- 添加使用示例
- 翻译多语言版本

### 问题反馈

- 通过 Issues 报告 bug
- 提出功能建议
- 分享使用经验

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

---

_此规格文档描述了 gopkg 项目的整体架构和学习路径，为 Go 语言学习者提供系统性的代码示例和实践指导。项目特别专注于提供业界最完整的 Go 语言设计模式实现集合。_
