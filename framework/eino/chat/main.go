package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func main() {
	chat("我的代码一直报错，感觉好沮丧，该怎么办？")
	chat("现在纽约时间几点？使用工具获取时间。")
}

func chat(question string) {
	// 0. 定义模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统提示
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),
		// 历史消息， 消息占位符：支持插入一组消息（如对话历史）
		schema.MessagesPlaceholder("chat_history", true),
		// 用户输入问题
		schema.UserMessage("问题：{question}"),
	)
	//-----------------------------------------------------------------------------------------------//

	// 1. 使用模板生成消息
	message, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": question,
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------------------------------------------------------------------------------//

	// 2. 创建聊天模型，这里指定连接哪个大语言模型
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:  "gpt-4o",                    // 使用的模型版本
		APIKey: os.Getenv("OPENAI_API_KEY"), // OpenAI API 密钥
	})
	if err != nil {
		log.Fatal(err)
	}
	// 添加工具，不需要工具调用则直接使用上面的chatModel对象
	tools := []tool.BaseTool{GetCurrentTimeTool()}
	toolInfos := make([]*schema.ToolInfo, 0)
	for _, tl := range tools {
		info, err := tl.Info(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		toolInfos = append(toolInfos, info)
	}
	toolChatModel, err := chatModel.WithTools(toolInfos)
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------------------------------------------------------------------------------//

	// 3. 请求模型接口输出结果
	// 3.1 一次性输出
	result, err := toolChatModel.Generate(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}
	// 结果输出
	if len(result.ToolCalls) > 0 {
		// 需要调用工具
		log.Println("调用工具：", result.ToolCalls)
		// 手动执行工具并收集结果
		var toolResponses []*schema.Message
		for _, tc := range result.ToolCalls {
			if tc.Function.Name == "get_current_time" {
				var args GetCurrentTimeInput
				err = json.Unmarshal([]byte(tc.Function.Arguments), &args)
				if err != nil {
					log.Println("参数解析错误：", tc.Function.Arguments)
					log.Fatal(err)
				}

				output, err := GetCurrentTime(context.Background(), args)
				if err != nil {
					output.Time = fmt.Sprintf("执行出错: %v", err)
				}
				toolResponses = append(toolResponses, schema.ToolMessage(output.Time, tc.ID))
			}
		}

		// 第二种调用工具的方法，使用工具节点调用
		/*conf := &compose.ToolsNodeConfig{
			Tools: tools,
		}
		toolsNode, err := compose.NewToolNode(context.Background(), conf)
		if err != nil {
			log.Fatal(err)
		}
		toolMessages, err := toolsNode.Invoke(context.Background(), result)
		if err != nil {
			log.Fatal(err)
		}
		finalMessages := append(message, toolMessages...)
		*/

		// 将工具结果追加，再次调用模型
		finalMessages := append(message, toolResponses...)
		finalResult, err := chatModel.Generate(context.Background(), finalMessages)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("✅ 最终回答：", finalResult.Content)
	} else {
		// 不需要调用工具
		log.Println("不需要调用工具")
		log.Println("✅ 最终回答:：", result.Content)
	}

	// 3.2 流式输出
	// 流式输出通常用于最终回答
	/*
		streamResult, err := toolChatModel.Stream(context.Background(), message)
		if err != nil {
			log.Fatal(err)
		}
		defer streamResult.Close()
		i := 0
		for {
			msg, err := streamResult.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("recv failed: %v", err)
			}
			log.Printf("message[%d]: %+v\n", i, msg)
			i++
		}
	*/
	//-----------------------------------------------------------------------------------------------//
}
