package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"stardew-agent/internal/models"
)

// DialogueService handles AI-driven dialogue generation
type DialogueService struct {
	llmClient LLMClient
	cache     *DialogueCache
}

// LLMClient defines the interface for LLM API calls
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
	GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// DialogueCache caches recent dialogues to reduce API calls
type DialogueCache struct {
	entries map[string]CachedDialogue
	maxSize int
}

// CachedDialogue represents a cached dialogue response
type CachedDialogue struct {
	Context   string
	Response  models.AIDialogueResponse
	Timestamp int64
}

// NewDialogueService creates a new dialogue service
func NewDialogueService(llmClient LLMClient) *DialogueService {
	return &DialogueService{
		llmClient: llmClient,
		cache: &DialogueCache{
			entries: make(map[string]CachedDialogue),
			maxSize: 100,
		},
	}
}

// BuildDialogueContext creates a dialogue context from game state
func BuildDialogueContext(
	npcID string,
	npcProfile models.NPCDialogueProfile,
	gameState *models.GameState,
	playerInput string,
	recentInteractions []models.Interaction,
	currentMood string,
) models.DialogueContext {
	// Calculate friendship hearts
	friendship := 0
	for _, npc := range gameState.NPCs {
		if npc.ID == npcID {
			friendship = npc.Friendship
			break
		}
	}
	hearts := friendship / 250

	// Determine day phase
	dayPhase := "morning"
	hour := gameState.Time.Hour
	if hour >= 6 && hour < 12 {
		dayPhase = "morning"
	} else if hour >= 12 && hour < 18 {
		dayPhase = "afternoon"
	} else if hour >= 18 && hour < 24 {
		dayPhase = "evening"
	} else {
		dayPhase = "night"
	}

	// Get weather (default to sunny if not set)
	weather := models.WeatherSunny

	// Get NPC location
	location := "unknown"
	for _, npc := range gameState.NPCs {
		if npc.ID == npcID {
			location = npc.Location
			break
		}
	}

	return models.DialogueContext{
		NPC:                npcProfile,
		Time:               gameState.Time,
		Weather:            weather,
		Friendship:         friendship,
		Hearts:             hearts,
		PlayerName:         "Player", // TODO: Get from player state
		RecentInteractions: recentInteractions,
		RecentEvents:       []string{}, // TODO: Track recent events
		CurrentMood:        currentMood,
		Location:           location,
		PlayerInput:        playerInput,
		DayPhase:           dayPhase,
	}
}

// GenerateDialogue generates AI dialogue for an NPC
func (ds *DialogueService) GenerateDialogue(ctx context.Context, dialogueCtx models.DialogueContext) (*models.AIDialogueResponse, error) {
	// Check cache first
	cacheKey := ds.buildCacheKey(dialogueCtx)
	if cached, ok := ds.cache.Get(cacheKey); ok {
		return &cached, nil
	}

	// Build system prompt
	systemPrompt := ds.buildSystemPrompt(dialogueCtx)

	// Build user prompt
	userPrompt := ds.buildUserPrompt(dialogueCtx)

	// Call LLM
	response, err := ds.llmClient.GenerateWithSystem(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Parse response
	aiResponse, err := ds.parseResponse(dialogueCtx.NPC.ID, response)
	if err != nil {
		log.Printf("Failed to parse LLM response: %v, using raw response", err)
		aiResponse = &models.AIDialogueResponse{
			NPCID:    dialogueCtx.NPC.ID,
			Dialogue: response,
		}
	}

	// Cache the response
	ds.cache.Set(cacheKey, *aiResponse)

	return aiResponse, nil
}

// buildSystemPrompt creates the system prompt for the LLM
func (ds *DialogueService) buildSystemPrompt(ctx models.DialogueContext) string {
	npc := ctx.NPC

	// Build personality description
	personalityDesc := fmt.Sprintf(
		"你是%s,一个%s角色。你的职业是%s。\n\n"+
			"【性格特质】\n%s\n\n"+
			"【背景故事】\n%s\n\n"+
			"【重视的事物】\n%s\n\n"+
			"【讨厌的事物】\n%s\n\n"+
			"【说话风格】\n%s\n\n"+
			"【爱好】\n%s",
		npc.Name,
		npc.Age,
		npc.Role,
		strings.Join(npc.Traits, ", "),
		npc.Background,
		strings.Join(npc.Values, ", "),
		strings.Join(npc.Dislikes, ", "),
		npc.SpeechStyle,
		npc.Hobby,
	)

	// Add secrets hint based on friendship level
	secretsHint := ""
	if ctx.Hearts >= 8 && len(npc.Secrets) > 0 {
		secretsHint = fmt.Sprintf("\n\n【秘密】(你信任对方,可以适当透露)\n%s", strings.Join(npc.Secrets, "\n"))
	}

	// Add dialogue theme hints
	themeHint := ds.getThemeHint(ctx)

	return personalityDesc + secretsHint + "\n\n" + themeHint + `

【对话规则】
1. 保持角色一致性,符合性格特质和说话风格
2. 根据友谊度调整亲密度:当前` + fmt.Sprintf("%d", ctx.Hearts) + `心(共10心)
3. 根据心情调整语气:当前心情是` + ctx.CurrentMood + `
4. 回复要自然,像真实对话一样
5. 可以主动提出话题或建议
6. 如果玩家的问题与你的职业/爱好相关,可以深入聊

【回复格式】请用JSON格式回复:
{
  "dialogue": "你的对话内容",
  "mood_change": "happy/neutral/sad/angry/excited",
  "friendship_delta": 0,
  "suggested_actions": ["建议的后续动作"],
  "secret_revealed": "如果透露了秘密,写在这里"
}`
}

// buildUserPrompt creates the user prompt with context
func (ds *DialogueService) buildUserPrompt(ctx models.DialogueContext) string {
	var prompt strings.Builder

	// Game context
	prompt.WriteString(fmt.Sprintf("【当前情境】\n"))
	prompt.WriteString(fmt.Sprintf("时间: %s %s, %d年\n", ctx.Time.Season, formatDay(ctx.Time.Day), ctx.Time.Year))
	prompt.WriteString(fmt.Sprintf("具体时间: %d:%02d (%s)\n", ctx.Time.Hour, ctx.Time.Minute, ctx.DayPhase))
	prompt.WriteString(fmt.Sprintf("天气: %s\n", ctx.Weather))
	prompt.WriteString(fmt.Sprintf("地点: %s\n", ctx.Location))
	prompt.WriteString(fmt.Sprintf("友谊度: %d心 (%d/2500点)\n", ctx.Hearts, ctx.Friendship))

	// Recent interactions
	if len(ctx.RecentInteractions) > 0 {
		prompt.WriteString("\n【最近的互动】\n")
		for _, interaction := range ctx.RecentInteractions {
			prompt.WriteString(fmt.Sprintf("- %s: %s\n", interaction.Type, interaction.Summary))
		}
	}

	// Player input
	prompt.WriteString(fmt.Sprintf("\n【玩家说】\n%s\n", ctx.PlayerInput))

	return prompt.String()
}

// getThemeHint returns dialogue theme hints based on context
func (ds *DialogueService) getThemeHint(ctx models.DialogueContext) string {
	npc := ctx.NPC
	themes := npc.DialogueThemes

	var hints []string

	// Check friendship-based themes first (for Haley)
	if ctx.Hearts >= 6 {
		switch ctx.DayPhase {
		case "morning":
			if len(themes.MorningHighFriendship) > 0 {
				hints = themes.MorningHighFriendship
			} else if len(themes.Morning) > 0 {
				hints = themes.Morning
			}
		case "afternoon":
			if len(themes.AfternoonHighFriendship) > 0 {
				hints = themes.AfternoonHighFriendship
			} else if len(themes.Afternoon) > 0 {
				hints = themes.Afternoon
			}
		}
	} else {
		switch ctx.DayPhase {
		case "morning":
			if len(themes.MorningLowFriendship) > 0 {
				hints = themes.MorningLowFriendship
			} else if len(themes.Morning) > 0 {
				hints = themes.Morning
			}
		case "afternoon":
			if len(themes.AfternoonLowFriendship) > 0 {
				hints = themes.AfternoonLowFriendship
			} else if len(themes.Afternoon) > 0 {
				hints = themes.Afternoon
			}
		case "evening":
			hints = themes.Evening
		}
	}

	// Weather override
	if ctx.Weather == models.WeatherRainy {
		if ctx.Hearts >= 6 && len(themes.RainyHighFriendship) > 0 {
			hints = themes.RainyHighFriendship
		} else if ctx.Hearts < 6 && len(themes.RainyLowFriendship) > 0 {
			hints = themes.RainyLowFriendship
		} else if len(themes.Rainy) > 0 {
			hints = themes.Rainy
		}
	}

	if len(hints) > 0 {
		return fmt.Sprintf("【对话主题参考】\n%s", strings.Join(hints, "\n"))
	}
	return ""
}

// parseResponse parses the LLM response into structured data
func (ds *DialogueService) parseResponse(npcID, response string) (*models.AIDialogueResponse, error) {
	// Try to extract JSON from response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		// No valid JSON found, use raw response as dialogue
		return &models.AIDialogueResponse{
			NPCID:    npcID,
			Dialogue: strings.TrimSpace(response),
		}, nil
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var aiResp models.AIDialogueResponse
	if err := json.Unmarshal([]byte(jsonStr), &aiResp); err != nil {
		return nil, err
	}

	aiResp.NPCID = npcID
	return &aiResp, nil
}

// buildCacheKey creates a cache key from dialogue context
func (ds *DialogueService) buildCacheKey(ctx models.DialogueContext) string {
	return fmt.Sprintf("%s_%s_%d_%s_%s",
		ctx.NPC.ID,
		ctx.DayPhase,
		ctx.Hearts,
		ctx.CurrentMood,
		ctx.PlayerInput,
	)
}

// formatDay formats day number with suffix
func formatDay(day int) string {
	return fmt.Sprintf("第%d天", day)
}

// Cache methods

// Get retrieves a cached dialogue
func (dc *DialogueCache) Get(key string) (models.AIDialogueResponse, bool) {
	entry, exists := dc.entries[key]
	if !exists {
		return models.AIDialogueResponse{}, false
	}
	return entry.Response, true
}

// Set stores a dialogue in cache
func (dc *DialogueCache) Set(key string, response models.AIDialogueResponse) {
	// Simple cache eviction if full
	if len(dc.entries) >= dc.maxSize {
		// Remove first entry (simple FIFO)
		for k := range dc.entries {
			delete(dc.entries, k)
			break
		}
	}

	dc.entries[key] = CachedDialogue{
		Context:  key,
		Response: response,
	}
}
