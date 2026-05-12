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

// 定义输入 schema
type RecipeInput struct {
	Ingredient          string `json:"ingredient" jsonschema:"description=主原料或菜系类型"`                  // 主原料或菜系类型
	DietaryRestrictions string `json:"dietaryRestrictions,omitempty" jsonschema:"description=任何饮食限制"` // 任何饮食限制
}

// 定义输出 schema
type Recipe struct {
	Title        string   `json:"title" jsonschema:"description=菜谱标题"`           // 菜谱标题
	Description  string   `json:"description" jsonschema:"description=菜谱描述"`     // 菜谱描述
	PrepTime     string   `json:"prepTime" jsonschema:"description=准备时间"`        // 准备时间
	CookTime     string   `json:"cookTime" jsonschema:"description=烹饪时间"`        // 烹饪时间
	Servings     int      `json:"servings" jsonschema:"description=份数"`          // 份数
	Ingredients  []string `json:"ingredients" jsonschema:"description=食材列表"`     // 食材列表
	Instructions []string `json:"instructions" jsonschema:"description=制作步骤"`    // 制作步骤
	Tips         []string `json:"tips,omitempty" jsonschema:"description=烹饪小贴士"` // 烹饪小贴士
}

func main() {
	ctx := context.Background()

	// 使用 Ollama 插件初始化 Genkit
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
			Type: "generate", // "chat" or "generate"
		},
		&ai.ModelOptions{
			Supports: &ai.ModelSupports{
				Multiturn:  true,
				SystemRole: true,
				Tools:      false,
				Media:      false,
			},
		},
	)

	// 定义一个菜谱生成器 Flow
	recipeGeneratorFlow := genkit.DefineFlow(g, "recipeGeneratorFlow", func(ctx context.Context, input *RecipeInput) (*Recipe, error) {
		// 根据输入构建提示词
		dietaryRestrictions := input.DietaryRestrictions
		if dietaryRestrictions == "" {
			dietaryRestrictions = "none"
		}

		prompt := fmt.Sprintf(`Create a recipe with the following requirements:
            Main ingredient: %s
            Dietary restrictions: %s`, input.Ingredient, dietaryRestrictions)

		// 使用相同 schema 生成结构化菜谱数据
		recipe, _, err := genkit.GenerateData[Recipe](ctx, g,
			ai.WithSystem("使用简体中文回答。"),
			ai.WithPrompt(prompt),
			ai.WithModel(model),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recipe: %w", err)
		}

		return recipe, nil
	})

	// 运行 Flow 进行一次测试
	recipe, err := recipeGeneratorFlow.Run(ctx, &RecipeInput{
		Ingredient:          "avocado",
		DietaryRestrictions: "vegetarian",
	})
	if err != nil {
		log.Fatalf("could not generate recipe: %v", err)
	}

	// 打印结构化菜谱
	recipeJSON, _ := json.MarshalIndent(recipe, "", "  ")
	fmt.Println("Sample recipe generated:")
	fmt.Println(string(recipeJSON))

	// 启动 HTTP 服务以提供 Flow，并保持应用运行以支持 Developer UI
	mux := http.NewServeMux()
	mux.HandleFunc("POST /recipeGeneratorFlow", genkit.Handler(recipeGeneratorFlow))

	log.Println("Starting server on http://localhost:3400")
	log.Println("Flow available at: POST http://localhost:3400/recipeGeneratorFlow")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}
