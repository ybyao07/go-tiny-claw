package provider

import (
	"context"
	"os"
	"testing"

	"github.com/ybyao07/go-tiny-claw/internal/schema"
)

func TestDeepSeekReasoningRoundTrip(t *testing.T) {
	if os.Getenv("DEEPSEEK_API_KEY") == "" {
		t.Skip("DEEPSEEK_API_KEY not set")
	}

	p := NewDeepSeekOpenAIProvider("deepseek-v4-flash")
	msgs := []schema.Message{
		{Role: schema.RoleSystem, Content: "You are a helper"},
		{Role: schema.RoleUser, Content: "say hi briefly"},
	}

	think, err := p.Generate(context.Background(), msgs, nil)
	if err != nil {
		t.Fatalf("thinking: %v", err)
	}
	t.Logf("content=%q reasoning_len=%d", think.Content, len(think.ReasoningContent))

	msgs = append(msgs, *think)
	tools := []schema.ToolDefinition{{
		Name:        "bash",
		Description: "run",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command": map[string]interface{}{"type": "string"},
			},
			"required": []string{"command"},
		},
	}}

	if _, err := p.Generate(context.Background(), msgs, tools); err != nil {
		t.Fatalf("action with tools: %v", err)
	}
}
