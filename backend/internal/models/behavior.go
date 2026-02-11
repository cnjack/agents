package models

// Mood represents NPC emotional state
type Mood string

const (
	MoodHappy      Mood = "happy"
	MoodNeutral    Mood = "neutral"
	MoodSad        Mood = "sad"
	MoodExcited    Mood = "excited"
	MoodAnxious    Mood = "anxious"
	MoodAngry      Mood = "angry"
	MoodTired      Mood = "tired"
	MoodRomantic   Mood = "romantic"
)

// BehaviorType represents types of NPC behaviors
type BehaviorType string

const (
	BehaviorIdle          BehaviorType = "idle"
	BehaviorWander        BehaviorType = "wander"
	BehaviorGoToLocation  BehaviorType = "go_to_location"
	BehaviorInteractPlayer BehaviorType = "interact_player"
	BehaviorInteractNPC   BehaviorType = "interact_npc"
	BehaviorWork          BehaviorType = "work"
	BehaviorSocialize     BehaviorType = "socialize"
	BehaviorRest          BehaviorType = "rest"
	BehaviorShop          BehaviorType = "shop"
	BehaviorFish          BehaviorType = "fish"
	BehaviorFarm          BehaviorType = "farm"
)

// DecisionFactor represents a factor in AI decision making
type DecisionFactor struct {
	Name       string  `json:"name"`
	Weight     float64 `json:"weight"`
	Value      float64 `json:"value"`
	RawScore   float64 `json:"raw_score"`
}

// BehaviorDecision represents an AI behavior decision
type BehaviorDecision struct {
	NPCID           string            `json:"npc_id"`
	BehaviorType    BehaviorType      `json:"behavior_type"`
	TargetLocation  *Position         `json:"target_location,omitempty"`
	TargetNPCID     string            `json:"target_npc_id,omitempty"`
	TargetPlayer    bool              `json:"target_player,omitempty"`
	Priority        int               `json:"priority"`
	Duration        int               `json:"duration"` // in game minutes
	Reason          string            `json:"reason"`
	Factors         []DecisionFactor  `json:"factors"`
	DeviatesFromSchedule bool         `json:"deviates_from_schedule"`
}

// AIDecisionRequest is the request for AI decision
type AIDecisionRequest struct {
	NPCID           string   `json:"npc_id"`
	PlayerPosition  Position `json:"player_position"`
	Weather         Weather  `json:"weather"`
	CurrentTime     TimeState `json:"current_time"`
	NearbyNPCs      []string `json:"nearby_npcs,omitempty"`
	RecentEvents    []string `json:"recent_events,omitempty"`
	OverrideSchedule bool    `json:"override_schedule,omitempty"`
}

// AIDecisionResponse is the response from AI decision
type AIDecisionResponse struct {
	Decision        BehaviorDecision `json:"decision"`
	CurrentMood     Mood             `json:"current_mood"`
	MoodReason      string           `json:"mood_reason"`
	ScheduleModified bool            `json:"schedule_modified"`
	InteractionTriggered *NPCInteraction `json:"interaction_triggered,omitempty"`
}

// NPCInteraction represents an interaction between NPCs or NPC and player
type NPCInteraction struct {
	InitiatorID    string        `json:"initiator_id"`
	TargetID       string        `json:"target_id"`
	InteractionType string       `json:"interaction_type"` // talk, greet, gift, argue, joke
	Dialogue       string        `json:"dialogue"`
	Duration       int           `json:"duration"` // game minutes
	Result         string        `json:"result,omitempty"`
}

// NPCMoodState tracks NPC mood over time
type NPCMoodState struct {
	NPCID           string        `json:"npc_id"`
	CurrentMood     Mood          `json:"current_mood"`
	MoodIntensity   float64       `json:"mood_intensity"` // 0.0 - 1.0
	MoodDuration    int           `json:"mood_duration"`  // game minutes remaining
	MoodHistory     []MoodEvent   `json:"mood_history"`
	MoodInfluencers []string      `json:"mood_influencers"`
}

// MoodEvent represents a mood-changing event
type MoodEvent struct {
	Timestamp   TimeState `json:"timestamp"`
	OldMood     Mood      `json:"old_mood"`
	NewMood     Mood      `json:"new_mood"`
	Trigger     string    `json:"trigger"`
	Intensity   float64   `json:"intensity"`
}

// NPCMood represents the current mood state of an NPC (used by services/ai)
type NPCMood struct {
	NPCID          string       `json:"npc_id"`
	CurrentMood    string       `json:"current_mood"`    // happy, neutral, sad, angry, excited
	MoodValue      int          `json:"mood_value"`      // -100 to 100
	MoodFactors    []MoodFactor `json:"mood_factors"`
	LastMoodChange int64        `json:"last_mood_change"`
}

// MoodFactor represents something affecting NPC mood (used by services/ai)
type MoodFactor struct {
	Factor      string `json:"factor"`      // weather, player_interaction, event, gift
	Impact      int    `json:"impact"`      // -100 to 100
	Duration    int    `json:"duration"`    // duration in game hours
	Description string `json:"description"`
	StartTime   int64  `json:"start_time"`
}

// DynamicScheduleEntry extends schedule with AI modifications
type DynamicScheduleEntry struct {
	BaseSchedule   ScheduleEntry `json:"base_schedule"`
	Modified       bool          `json:"modified"`
	ModificationReason string    `json:"modification_reason,omitempty"`
	OverridePosition *Position   `json:"override_position,omitempty"`
	OverrideLocation string      `json:"override_location,omitempty"`
	Priority       int           `json:"priority"`
}

// NPCPersonality defines NPC personality traits
type NPCPersonality struct {
	NPCID              string  `json:"npc_id"`
	Sociability        float64 `json:"sociability"`        // 0.0 - 1.0, how much they seek social interaction
	Spontaneity        float64 `json:"spontaneity"`        // 0.0 - 1.0, how likely to deviate from schedule
	WorkEthic          float64 `json:"work_ethic"`         // 0.0 - 1.0, how dedicated to work tasks
	OutdoorsPreference float64 `json:"outdoors_preference"` // 0.0 - 1.0
	MorningPerson      float64 `json:"morning_person"`     // 0.0 - 1.0, affects mood at different times
	GiftAppreciation   float64 `json:"gift_appreciation"`  // 0.0 - 1.0, mood boost from gifts
	WeatherSensitivity float64 `json:"weather_sensitivity"` // 0.0 - 1.0, how weather affects mood
}

// BehaviorContext provides context for decision making
type BehaviorContext struct {
	NPC             NPCState        `json:"npc"`
	PlayerPosition  Position        `json:"player_position"`
	PlayerFriendship int            `json:"player_friendship"`
	Weather         Weather         `json:"weather"`
	Time            TimeState       `json:"time"`
	NearbyNPCs      []NPCState      `json:"nearby_npcs"`
	RecentInteractions []string     `json:"recent_interactions"`
	DayEvents       []string        `json:"day_events"`
	Personality     NPCPersonality  `json:"personality"`
	MoodState       NPCMoodState    `json:"mood_state"`
}
