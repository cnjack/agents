package models

import "time"

// Character represents a game character with full AI-driven properties
type Character struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Role             string             `json:"role"`
	Age              string             `json:"age"`
	Gender           string             `json:"gender"`
	Traits           []string           `json:"traits"`
	Values           []string           `json:"values"`
	Dislikes         []string           `json:"dislikes"`
	Fears            []string           `json:"fears,omitempty"`
	Hopes            []string           `json:"hopes,omitempty"`
	Background       string             `json:"background"`
	Secrets          []string           `json:"secrets,omitempty"`
	SpeechStyle      string             `json:"speech_style"`
	Hobby            string             `json:"hobby"`
	GiftPreferences  GiftPreferences    `json:"gift_preferences"`
	DialogueThemes   DialogueThemes     `json:"dialogue_themes"`
	Schedule         []ScheduleEntry    `json:"schedule,omitempty"`
	MaxFriendship    int                `json:"max_friendship"`
	CharacterArc     string             `json:"character_arc,omitempty"`
	CreatedAt        time.Time          `json:"created_at,omitempty"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty"`
}

// NPCDialogueProfile defines the personality profile of an NPC for AI-driven dialogue
// This is separate from the behavior-focused NPCPersonality in behavior.go
type NPCDialogueProfile struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Role               string            `json:"role"`               // mayor, shopkeeper, carpenter, villager, fisherman
	Age                string            `json:"age"`                // young, middle, senior
	Traits             []string          `json:"traits"`             // personality traits
	Background         string            `json:"background"`         // background story
	Values             []string          `json:"values"`             // what they value
	Dislikes           []string          `json:"dislikes"`           // what they dislike
	Fears              []string          `json:"fears,omitempty"`    // what they fear
	Hopes              []string          `json:"hopes,omitempty"`    // what they hope for
	SpeechStyle        string            `json:"speech_style"`       // how they speak
	Hobby              string            `json:"hobby"`              // their hobby
	GiftPreferences    GiftPreferences   `json:"gift_preferences"`
	DialogueThemes     DialogueThemes    `json:"dialogue_themes"`
	Secrets            []string          `json:"secrets"`            // secrets revealed at high friendship
	FriendshipUnlocks  map[string]string `json:"friendship_unlocks"` // heart level -> unlock description
	CharacterArc       string            `json:"character_arc,omitempty"` // character development arc
}

// GiftPreferences defines what items an NPC loves, likes, etc.
type GiftPreferences struct {
	Love    []string `json:"love"`
	Like    []string `json:"like"`
	Neutral []string `json:"neutral"`
	Dislike []string `json:"dislike"`
	Hate    []string `json:"hate"`
}

// DialogueThemes defines dialogue themes for different situations
type DialogueThemes struct {
	Morning               []string `json:"morning,omitempty"`
	Afternoon             []string `json:"afternoon,omitempty"`
	Evening               []string `json:"evening,omitempty"`
	Rainy                 []string `json:"rainy,omitempty"`
	Festival              []string `json:"festival,omitempty"`
	MorningLowFriendship  []string `json:"morning_low_friendship,omitempty"`
	MorningHighFriendship []string `json:"morning_high_friendship,omitempty"`
	AfternoonLowFriendship  []string `json:"afternoon_low_friendship,omitempty"`
	AfternoonHighFriendship []string `json:"afternoon_high_friendship,omitempty"`
	RainyLowFriendship    []string `json:"rainy_low_friendship,omitempty"`
	RainyHighFriendship   []string `json:"rainy_high_friendship,omitempty"`
}

// Interaction represents a past interaction between player and NPC
type Interaction struct {
	Type      string `json:"type"`       // talk, gift, quest
	Timestamp int64  `json:"timestamp"`
	Summary   string `json:"summary"`    // brief summary of interaction
	Result    string `json:"result"`     // outcome of interaction
}

// DialogueContext contains all context needed for AI dialogue generation
type DialogueContext struct {
	NPC                NPCDialogueProfile `json:"npc"`
	Time               TimeState          `json:"time"`
	Weather            Weather            `json:"weather"`
	Friendship         int                `json:"friendship"`
	Hearts             int                `json:"hearts"`
	PlayerName         string             `json:"player_name"`
	RecentInteractions []Interaction      `json:"recent_interactions"`
	RecentEvents       []string           `json:"recent_events"`
	CurrentMood        string             `json:"current_mood"`
	Location           string             `json:"location"`
	PlayerInput        string             `json:"player_input"` // natural language input from player
	DayPhase           string             `json:"day_phase"`    // morning, afternoon, evening, night
}

// AIDialogueResponse is the response from AI dialogue generation
type AIDialogueResponse struct {
	NPCID            string   `json:"npc_id"`
	Dialogue         string   `json:"dialogue"`          // generated dialogue text
	MoodChange       string   `json:"mood_change"`       // how mood changed
	FriendshipDelta  int      `json:"friendship_delta"`  // friendship change from this interaction
	SuggestedActions []string `json:"suggested_actions"` // suggested follow-up actions
	TriggerEvent     string   `json:"trigger_event"`     // optional event to trigger
	SecretRevealed   string   `json:"secret_revealed,omitempty"` // secret revealed if friendship high enough
}

// NPCMoodData represents the current mood state of an NPC
type NPCMoodData struct {
	NPCID          string        `json:"npc_id"`
	CurrentMood    string        `json:"current_mood"`    // happy, neutral, sad, angry, excited
	MoodValue      int           `json:"mood_value"`      // -100 to 100
	MoodFactors    []MoodFactorData `json:"mood_factors"`
	LastMoodChange int64         `json:"last_mood_change"`
}

// MoodFactorData represents something affecting NPC mood
type MoodFactorData struct {
	Factor      string `json:"factor"`      // weather, player_interaction, event, gift
	Impact      int    `json:"impact"`      // -100 to 100
	Duration    int    `json:"duration"`    // duration in game hours
	Description string `json:"description"`
	StartTime   int64  `json:"start_time"`
}

// NLPRequest is for natural language understanding
type NLPRequest struct {
	PlayerInput     string           `json:"player_input"`
	GameState       *GameState       `json:"game_state"`
	NearbyNPCs      []NPCState       `json:"nearby_npcs"`
	PlayerInventory []InventoryItem  `json:"player_inventory"`
	Context         string           `json:"context"`
}

// NLPResponse is the result of natural language processing
type NLPResponse struct {
	Understood            bool    `json:"understood"`
	Intent                string  `json:"intent"`               // talk, gift, buy, sell, move, etc.
	Action                *Action `json:"action,omitempty"`     // converted game action
	TargetNPC             string  `json:"target_npc,omitempty"` // target NPC if any
	TargetItem            string  `json:"target_item,omitempty"`
	ResponseText          string  `json:"response_text"`
	NeedsClarification    bool    `json:"needs_clarification"`
	ClarificationQuestion string  `json:"clarification_question,omitempty"`
}

// AIBehaviorRequest is for AI-driven NPC behavior decisions
type AIBehaviorRequest struct {
	NPCID           string            `json:"npc_id"`
	Profile         NPCDialogueProfile `json:"profile"`
	BehaviorState   NPCBehaviorStateData `json:"behavior_state"`
	GameState       *GameState        `json:"game_state"`
	PlayerNearby    bool              `json:"player_nearby"`
	Weather         Weather           `json:"weather"`
	Time            TimeState         `json:"time"`
	PossibleActions []string          `json:"possible_actions"`
}

// NPCBehaviorStateData tracks NPC's current behavioral state
type NPCBehaviorStateData struct {
	NPCID          string `json:"npc_id"`
	CurrentGoal    string `json:"current_goal"`
	CurrentMood    string `json:"current_mood"`
	Energy         int    `json:"energy"`
	SocialNeed     int    `json:"social_need"`
	LastDecision   string `json:"last_decision"`
	DecisionReason string `json:"decision_reason"`
}

// AIBehaviorResponseData is the AI's decision for NPC behavior
type AIBehaviorResponseData struct {
	NPCID           string   `json:"npc_id"`
	Action          string   `json:"action"`            // move, stay, interact, talk
	TargetLocation  string   `json:"target_location"`
	TargetPosition  Position `json:"target_position"`
	Reason          string   `json:"reason"`
	MoodAfterAction string   `json:"mood_after_action"`
	CanInterrupt    bool     `json:"can_interrupt"`
}

// DialogueMood constants for dialogue system
const (
	DialogueMoodHappy   string = "happy"
	DialogueMoodNeutral string = "neutral"
	DialogueMoodSad     string = "sad"
	DialogueMoodAngry   string = "angry"
	DialogueMoodExcited string = "excited"
)

// MoodDialogueModifiers defines how mood affects dialogue
var MoodDialogueModifiers = map[string]string{
	DialogueMoodHappy:   "语气愉快,愿意分享更多信息",
	DialogueMoodNeutral: "正常语气",
	DialogueMoodSad:     "语气低落,回复简短",
	DialogueMoodAngry:   "语气生硬,可能拒绝互动",
	DialogueMoodExcited: "语气兴奋,可能主动分享秘密",
}

// GetMoodModifier returns the dialogue modifier for a mood
func GetMoodModifier(mood string) string {
	if modifier, ok := MoodDialogueModifiers[mood]; ok {
		return modifier
	}
	return MoodDialogueModifiers[DialogueMoodNeutral]
}
