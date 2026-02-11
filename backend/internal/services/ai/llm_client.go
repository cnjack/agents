package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"stardew-agent/internal/models"
)

// MockLLMClient is a mock implementation for testing without actual LLM
type MockLLMClient struct {
	profiles map[string]models.NPCDialogueProfile
}

// NewMockLLMClient creates a new mock LLM client
func NewMockLLMClient(profiles map[string]models.NPCDialogueProfile) *MockLLMClient {
	return &MockLLMClient{
		profiles: profiles,
	}
}

// Generate generates a simple response (mock implementation)
func (m *MockLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	// Simulate some processing delay
	time.Sleep(100 * time.Millisecond)

	// Simple mock response
	return `{"dialogue": "你好!今天天气不错。", "mood_change": "neutral", "friendship_delta": 0}`, nil
}

// GenerateWithSystem generates a response with system context (mock implementation)
func (m *MockLLMClient) GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// Simulate some processing delay
	time.Sleep(100 * time.Millisecond)

	// Extract NPC info from system prompt
	npcName := "NPC"
	for name := range m.profiles {
		if strings.Contains(systemPrompt, name) {
			npcName = name
			break
		}
	}

	// Extract context info
	dayPhase := "afternoon"
	if strings.Contains(userPrompt, "morning") {
		dayPhase = "morning"
	} else if strings.Contains(userPrompt, "evening") {
		dayPhase = "evening"
	}

	hearts := 0
	if strings.Contains(userPrompt, "友谊度") {
		// Try to parse hearts
		fmt.Sscanf(userPrompt, "友谊度: %d心", &hearts)
	}

	// Generate context-aware mock response
	dialogue := m.generateMockDialogue(npcName, dayPhase, hearts)

	response := models.AIDialogueResponse{
		NPCID:           strings.ToLower(npcName),
		Dialogue:        dialogue,
		MoodChange:       "happy",
		FriendshipDelta:  0,
		SuggestedActions: []string{"talk", "give_gift"},
	}

	jsonBytes, _ := json.Marshal(response)
	return string(jsonBytes), nil
}

// generateMockDialogue creates contextual mock dialogue
func (m *MockLLMClient) generateMockDialogue(npcName, dayPhase string, hearts int) string {
	rand.Seed(time.Now().UnixNano())

	dialogues := map[string]map[string][]string{
		"Lewis": {
			"morning": {
				"早安!美好的一天开始了,我们小镇需要像你这样勤劳的农民。",
				"今天天气不错,适合在社区中心附近走走。",
				"亲爱的朋友,有什么需要我帮忙的吗?",
			},
			"afternoon": {
				"下午好!社区中心的重建工作进行得怎么样了?",
				"我们小镇的发展离不开大家的努力啊。",
				"最近有什么新鲜事吗?",
			},
			"evening": {
				"晚上好,辛苦了一天,好好休息吧。",
				"夜晚的小镇特别宁静,不是吗?",
			},
		},
		"Pierre": {
			"morning": {
				"早安!今天有什么需要买的吗?新进的货物很棒!",
				"欢迎光临!早起的鸟儿有虫吃,早来的顾客有优惠!",
			},
			"afternoon": {
				"来看看今天的特价商品吧!",
				"刚从城里进了新货,品质一流!",
				"最近生意还不错,谢谢大家支持。",
			},
			"evening": {
				"快打烊了,还有什么需要的吗?",
				"辛苦了,明天见!",
			},
		},
		"Robin": {
			"morning": {
				"早!今天有什么建筑项目需要帮忙吗?",
				"早上空气好,适合在工坊干活。",
			},
			"afternoon": {
				"需要建造什么?我可以帮你设计!",
				"来看看我最近的作品,怎么样?",
				"建筑是一门艺术,需要用心对待。",
			},
			"evening": {
				"辛苦了一天,该休息了。",
				"晚上好好休息,明天继续努力!",
			},
		},
		"Haley": {
			"morning": {
				"...你想干嘛?",
				"嗯...早。",
			},
			"afternoon": {
				"我在拍照,别打扰我。",
				"这个角度的光线很完美...",
			},
			"evening": {
				"晚上好...还是没什么好聊的。",
			},
		},
		"Willy": {
			"morning": {
				"早安!今天是个钓鱼的好日子!",
				"海风很舒服,适合出海。",
			},
			"afternoon": {
				"钓到什么好鱼了吗?",
				"来聊聊钓鱼心得吧!",
				"海洋总是充满惊喜。",
			},
			"evening": {
				"晚上好!要去酒吧喝一杯吗?",
				"一天的收获怎么样?",
			},
		},
	}

	if npcDialogues, ok := dialogues[npcName]; ok {
		if phaseDialogues, ok := npcDialogues[dayPhase]; ok && len(phaseDialogues) > 0 {
			// Add friendship-based variations
			if hearts >= 6 && npcName == "Haley" {
				return fmt.Sprintf("其实...看到你挺好的。%s", phaseDialogues[rand.Intn(len(phaseDialogues))])
			}
			return phaseDialogues[rand.Intn(len(phaseDialogues))]
		}
	}

	return "你好,有什么事吗?"
}

// ClaudeLLMClient implements LLMClient using Anthropic Claude API
type ClaudeLLMClient struct {
	apiKey     string
	model      string
	maxTokens  int
	baseURL    string
}

// ClaudeConfig contains configuration for Claude API
type ClaudeConfig struct {
	APIKey    string
	Model     string // e.g., "claude-3-sonnet-20240229"
	MaxTokens int
	BaseURL   string // optional, defaults to Anthropic API
}

// NewClaudeLLMClient creates a new Claude LLM client
func NewClaudeLLMClient(config ClaudeConfig) *ClaudeLLMClient {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1024
	}

	return &ClaudeLLMClient{
		apiKey:    config.APIKey,
		model:     config.Model,
		maxTokens: maxTokens,
		baseURL:   baseURL,
	}
}

// Generate generates a response using Claude API
func (c *ClaudeLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	return c.GenerateWithSystem(ctx, "你是一个游戏中的NPC角色,请用自然的语气与玩家对话。", prompt)
}

// GenerateWithSystem generates a response with system context using Claude API
func (c *ClaudeLLMClient) GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// Note: Actual implementation would use the Anthropic SDK
	// This is a placeholder that should be implemented with proper HTTP client
	/*
		Example implementation:

		client := anthropic.NewClient(c.apiKey)
		resp, err := client.CreateMessages(ctx, &anthropic.MessagesRequest{
			Model: c.model,
			Messages: []anthropic.Message{
				{
					Role:    anthropic.RoleUser,
					Content: userPrompt,
				},
			},
			System:    systemPrompt,
			MaxTokens: c.maxTokens,
		})
		if err != nil {
			return "", err
		}
		return resp.Content[0].Text, nil
	*/

	// For now, return mock response
	return fmt.Sprintf(`{"dialogue": "这是一条来自AI的回复。系统提示: %s", "mood_change": "neutral", "friendship_delta": 0}`,
		systemPrompt[:min(50, len(systemPrompt))]), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
