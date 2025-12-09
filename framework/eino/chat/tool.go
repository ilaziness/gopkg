package main

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// 自定义工具
// 官方预定义的工具列表：https://github.com/cloudwego/eino-ext/tree/main/components/tool

// GetCurrentTimeInput 获取当前时间的输入参数
// 标签说明：https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/#%E6%96%B9%E5%BC%8F-2---openapi3schema
type GetCurrentTimeInput struct {
	Timezone string `json:"timezone" jsonschema:"required,description=获取时间的时区"`
}

// GetCurrentTimeOutput 获取当前时间的输出参数
type GetCurrentTimeOutput struct {
	Time string `json:"time" jsonschema:"description=time"`
}

// GetCurrentTime 获取当前时间的工具，返回参数没有复杂结构也可以简单地定义为string这种简单类型
func GetCurrentTime(_ context.Context, input GetCurrentTimeInput) (GetCurrentTimeOutput, error) {
	t := time.Now().In(time.FixedZone(input.Timezone, 8*60*60)).Format(time.RFC3339)
	return GetCurrentTimeOutput{Time: t}, nil
}

// GetCurrentTimeTool 创建工具
func GetCurrentTimeTool() tool.InvokableTool {
	t, err := utils.InferTool("get_current_time", "获取指定时区的当前时间", GetCurrentTime)
	if err != nil {
		panic(err)
	}
	return t
}
