package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
)

// ============================================================================
// Schema 定义
// ============================================================================

// WeatherInput - 天气工具输入
type WeatherInput struct {
	Location string `json:"location" jsonschema:"description=要查询天气的地点"`
}

// MenuQueryInput - 菜单查询输入
type MenuQueryInput struct {
	Question string `json:"question" jsonschema:"description=关于菜单的问题"`
}

// MenuQueryOutput - 菜单查询输出
type MenuQueryOutput struct {
	Answer string `json:"answer" jsonschema:"description=回答"`
}

// ============================================================================
// main
// ============================================================================

func main() {
	ctx := context.Background()

	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       80,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(ollamaPlugin),
	)
	model := ollamaPlugin.DefineModel(
		g,
		ollama.ModelDefinition{
			Name: "gemma4:e2b",
			Type: "generate",
		},
		&ai.ModelOptions{
			Supports: &ai.ModelSupports{
				Multiturn:  true,
				SystemRole: true,
				Tools:      true,
				Media:      false,
			},
		},
	)

	// ====================================================================
	// 定义工具: getWeather - 获取天气
	// 使用 genkit.DefineTool() 定义工具
	// 参数: g, name, description, handler
	// description 对 LLM 判断何时使用工具至关重要，须详细准确
	// ====================================================================
	getWeatherTool := genkit.DefineTool(
		g,
		"getWeather",
		"获取指定地点的当前天气信息。该工具可以返回任何地点的天气情况，包括温度、天气状况等。",
		func(ctx *ai.ToolContext, input WeatherInput) (string, error) {
			log.Printf("工具被调用: getWeather(%s)", input.Location)
			// 模拟真实 API 调用
			return fmt.Sprintf("%s 当前天气: 22°C, 晴朗", input.Location), nil
		},
	)

	// ====================================================================
	// 定义工具: getDailyMenu - 获取每日菜单
	// ====================================================================
	getDailyMenuTool := genkit.DefineTool(
		g,
		"getDailyMenu",
		"获取餐厅今日菜单。返回今日所有可点的菜品列表，包括名称、描述和价格。",
		func(ctx *ai.ToolContext, _ struct{}) (string, error) {
			menu := []map[string]interface{}{
				{"name": "凯撒沙拉", "price": "¥38", "description": "新鲜罗马生菜配凯撒酱"},
				{"name": "意大利面", "price": "¥58", "description": "番茄肉酱意面"},
				{"name": "提拉米苏", "price": "¥32", "description": "经典意式甜点"},
			}
			b, _ := json.Marshal(menu)
			return string(b), nil
		},
	)

	// ====================================================================
	// 示例 1: 工具调用 - 基础使用
	// 模型自动决定何时使用工具，Genkit 自动处理工具调用循环
	// ====================================================================
	weatherFlow := genkit.DefineFlow(g, "weatherFlow",
		func(ctx context.Context, input *struct{ Location string `json:"location"` }) (string, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("%s 的天气怎么样？", input.Location),
				ai.WithTools(getWeatherTool),  // 注册可用工具
			)
			if err != nil {
				return "", fmt.Errorf("weather flow failed: %w", err)
			}
			return resp.Text(), nil
		})

	// ====================================================================
	// 示例 2: 带工具调用的菜单问答
	// LLM 使用菜单工具获取数据后回答用户问题
	// ====================================================================
	menuFlow := genkit.DefineFlow(g, "menuFlow",
		func(ctx context.Context, input *MenuQueryInput) (*MenuQueryOutput, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithSystem("你是一个餐厅助手。使用工具获取菜单数据来回答客户问题。使用中文。"),
				ai.WithPrompt(input.Question),
				ai.WithTools(getDailyMenuTool),
			)
			if err != nil {
				return nil, fmt.Errorf("menu flow failed: %w", err)
			}
			return &MenuQueryOutput{Answer: resp.Text()}, nil
		})

	// ====================================================================
	// 示例 3: 显式处理工具调用 - ai.WithReturnToolRequests(true)
	// 当需要完全控制工具调用循环时使用
	// ====================================================================
	explicitToolFlow := genkit.DefineFlow(g, "explicitToolFlow",
		func(ctx context.Context, input *struct{ Location string `json:"location"` }) (string, error) {
			// 第一次调用：返回工具请求而非执行
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("%s 的天气怎么样？请使用工具查询。", input.Location),
				ai.WithTools(getWeatherTool),
				ai.WithReturnToolRequests(true),  // 阻止自动执行工具
			)
			if err != nil {
				return "", fmt.Errorf("first generate failed: %w", err)
			}

			// 手动处理每个工具请求
			var parts []*ai.Part
			for _, req := range resp.ToolRequests() {
				tool := genkit.LookupTool(g, req.Name)
				if tool == nil {
					return "", fmt.Errorf("工具 %q 未找到", req.Name)
				}
				output, err := tool.RunRaw(ctx, req.Input)
				if err != nil {
					return "", fmt.Errorf("工具 %q 执行失败: %w", tool.Name(), err)
				}
				parts = append(parts, ai.NewToolResponsePart(&ai.ToolResponse{
					Name:   req.Name,
					Ref:    req.Ref,
					Output: output,
				}))
			}

			// 第二次调用：将工具结果发送回模型
			resp2, err := genkit.Generate(ctx, g,
				ai.WithMessages(append(resp.History(),
					ai.NewMessage(ai.RoleTool, nil, parts...))...),
			)
			if err != nil {
				return "", fmt.Errorf("second generate failed: %w", err)
			}
			return resp2.Text(), nil
		})

	// ====================================================================
	// 本地测试
	// ====================================================================

	fmt.Println("========== 示例 1: 天气查询（自动工具调用）==========")
	result, err := weatherFlow.Run(ctx, &struct{ Location string }{Location: "北京"})
	if err != nil {
		log.Printf("weatherFlow error: %v", err)
	} else {
		fmt.Println(result)
	}

	fmt.Println("\n========== 测试 3: 显式工具调用 ==========")
	result3, err := explicitToolFlow.Run(ctx, &struct{ Location string }{Location: "上海"})
	if err != nil {
		log.Printf("explicitToolFlow error: %v", err)
	} else {
		fmt.Println(result3)
	}

	// ====================================================================
	// HTTP 部署
	// ====================================================================
	mux := http.NewServeMux()
	mux.HandleFunc("POST /weatherFlow", genkit.Handler(weatherFlow))
	mux.HandleFunc("POST /menuFlow", genkit.Handler(menuFlow))
	mux.HandleFunc("POST /explicitToolFlow", genkit.Handler(explicitToolFlow))

	fmt.Println("\n========== 已部署的 Flow 端点 ==========")
	for _, flow := range genkit.ListFlows(g) {
		fmt.Printf("  POST /%s\n", flow.Name())
	}

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}