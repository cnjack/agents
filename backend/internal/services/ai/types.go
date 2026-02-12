package ai

import (
	"stardew-agent/internal/models"
)

// ProviderType AI 提供者类型
type ProviderType string

const (
	ProviderTypeOpenAICompatible ProviderType = "openai_compatible"
	ProviderTypeMock             ProviderType = "mock"
	ProviderTypeOllama           ProviderType = "ollama"
)

// ProviderInfo 提供者信息
type ProviderInfo struct {
	Type      ProviderType `json:"type"`
	Model     string       `json:"model"`
	Endpoint  string       `json:"endpoint"`
	Available bool         `json:"available"`
	Latency   int64        `json:"latency_ms"`
}

// DialogueRequest 对话生成请求
type DialogueRequest struct {
	Character       *models.Character `json:"character"`
	PlayerInput     string            `json:"player_input"`
	Context         string            `json:"context"`
	CurrentMood     string            `json:"current_mood"`
	FriendshipHearts int              `json:"friendship_hearts"`
	GameState       *models.GameState `json:"game_state,omitempty"`
	RecentInteractions []models.Interaction `json:"recent_interactions,omitempty"`
}

// DialogueResponse 对话生成响应
type DialogueResponse struct {
	Dialogue        string   `json:"dialogue"`
	MoodChange      string   `json:"mood_change"`
	FriendshipDelta int      `json:"friendship_delta"`
	SuggestedActions []string `json:"suggested_actions"`
	SecretRevealed  string   `json:"secret_revealed,omitempty"`
}

// CharacterRequest 角色生成请求
type CharacterRequest struct {
	NameHint         string   `json:"name_hint"`
	Role             string   `json:"role"`
	Age              string   `json:"age"`
	Gender           string   `json:"gender"`
	PersonalityHint  string   `json:"personality_hint"`
	TraitsHint       []string `json:"traits_hint"`
	GenerateSchedule bool     `json:"generate_schedule"`
	GenerateDialogue bool     `json:"generate_dialogue"`
	GenerateBehaviors bool    `json:"generate_behaviors"`
	Creativity       float64  `json:"creativity"`
}

// CharacterResponse 角色生成响应
type CharacterResponse struct {
	Name            string                  `json:"name"`
	Role            string                  `json:"role"`
	Age             string                  `json:"age"`
	Gender          string                  `json:"gender"`
	Traits          []string                `json:"traits"`
	Values          []string                `json:"values"`
	Dislikes        []string                `json:"dislikes"`
	Fears           []string                `json:"fears,omitempty"`
	Hopes           []string                `json:"hopes,omitempty"`
	Background      string                  `json:"background"`
	Secrets         []string                `json:"secrets"`
	SpeechStyle     string                  `json:"speech_style"`
	Hobby           string                  `json:"hobby"`
	GiftPreferences models.GiftPreferences  `json:"gift_preferences"`
	DialogueThemes  models.DialogueThemes   `json:"dialogue_themes"`
	CharacterArc    string                  `json:"character_arc"`
	Schedule        []models.ScheduleEntry  `json:"schedule,omitempty"`
}

// NLURequest 自然语言理解请求
type NLURequest struct {
	Input          string            `json:"input"`
	GameState      *models.GameState `json:"game_state"`
	NearbyNPCs     []models.NPCState `json:"nearby_npcs"`
	PlayerInventory []models.InventoryItem `json:"player_inventory"`
	Context        string            `json:"context"`
}

// NLUResponse 自然语言理解响应
type NLUResponse struct {
	Understood            bool            `json:"understood"`
	Intent                string          `json:"intent"`
	Action                *models.Action  `json:"action,omitempty"`
	TargetNPC             string          `json:"target_npc,omitempty"`
	TargetItem            string          `json:"target_item,omitempty"`
	ResponseText          string          `json:"response_text"`
	NeedsClarification    bool            `json:"needs_clarification"`
	ClarificationQuestion string          `json:"clarification_question,omitempty"`
}

// BehaviorRequest 行为决策请求
type BehaviorRequest struct {
	NPCID           string   `json:"npc_id"`
	CurrentMood     string   `json:"current_mood"`
	CurrentGoal     string   `json:"current_goal"`
	Energy          int      `json:"energy"`
	SocialNeed      int      `json:"social_need"`
	Season          string   `json:"season"`
	Hour            int      `json:"hour"`
	Minute          int      `json:"minute"`
	Weather         string   `json:"weather"`
	PlayerNearby    bool     `json:"player_nearby"`
	PossibleActions []string `json:"possible_actions"`
	Context         string   `json:"context"`
}

// BehaviorResponse 行为决策响应
type BehaviorResponse struct {
	Action          string          `json:"action"`
	TargetLocation  string          `json:"target_location"`
	TargetPosition  models.Position `json:"target_position"`
	Reason          string          `json:"reason"`
	MoodAfterAction string          `json:"mood_after_action"`
	CanInterrupt    bool            `json:"can_interrupt"`
}
