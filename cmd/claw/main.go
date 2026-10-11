// cmd/claw/main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/ybyao07/go-tiny-claw/internal/engine"
	"github.com/ybyao07/go-tiny-claw/internal/feishu"
	"github.com/ybyao07/go-tiny-claw/internal/provider"
	"github.com/ybyao07/go-tiny-claw/internal/tools"
)

func main() {
	// 1. 初始化引擎依赖
	workDir, _ := os.Getwd()

	if os.Getenv("DEEPSEEK_API_KEY") == "" {
		log.Fatal("请先导出 DEEPSEEK_API_KEY 环境变量")
	}
	llmProvider := provider.NewDeepSeekOpenAIProvider("deepseek-v4-flash")

	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))

	// 开启慢思考
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	// 2. 通过 WebSocket 长连接接入飞书（无需公网 HTTP webhook）
	bot := feishu.NewFeishuBot(eng)
	if err := bot.Start(context.Background()); err != nil {
		log.Fatalf("飞书长连接失败: %v", err)
	}
}
