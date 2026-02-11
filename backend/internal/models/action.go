package models

// ActionType represents types of actions an agent can perform
type ActionType string

const (
	ActionMove        ActionType = "move"
	ActionInteract    ActionType = "interact"
	ActionTalk        ActionType = "talk"
	ActionGiveGift    ActionType = "give_gift"
	ActionAcceptQuest ActionType = "accept_quest"
	ActionWait        ActionType = "wait"
	ActionSleep       ActionType = "sleep"
)

// Action represents an agent action
type Action struct {
	Type ActionType `json:"type"`
	Params ActionParams `json:"params"`
}

// ActionParams contains parameters for different action types
type ActionParams struct {
	// Move action
	Direction Direction `json:"direction,omitempty"`

	// Interact action
	Target string `json:"target,omitempty"` // npc id, object id, or "front"

	// Talk action
	NPC string `json:"npc,omitempty"`

	// GiveGift action
	GiftItem string `json:"gift_item,omitempty"`

	// AcceptQuest action
	QuestID string `json:"quest_id,omitempty"`

	// Wait action
	Ticks int `json:"ticks,omitempty"`
}

// ActionResult represents the result of an action
type ActionResult struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	NewState  *GameState `json:"new_state,omitempty"`
	Events    []GameEvent `json:"events,omitempty"`
}

// GameEvent represents something that happened in the game
type GameEvent struct {
	Type      string      `json:"type"` // crop_grown, quest_completed, item_received, npc_dialog
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// MoveActionResult specific result for move action
type MoveActionResult struct {
	Success   bool     `json:"success"`
	OldPos    Position `json:"old_pos"`
	NewPos    Position `json:"new_pos"`
	BlockedBy string   `json:"blocked_by,omitempty"`
}

// SocialActionResult specific result for NPC interactions
type SocialActionResult struct {
	Success        bool   `json:"success"`
	NPCID          string `json:"npc_id"`
	NPCName        string `json:"npc_name"`
	Dialogue       string `json:"dialogue,omitempty"`
	FriendshipGain int    `json:"friendship_gain"`
	NewFriendship  int    `json:"new_friendship"`
}

// QuestActionResult specific result for quest actions
type QuestActionResult struct {
	Success    bool       `json:"success"`
	QuestID    string     `json:"quest_id"`
	QuestName  string     `json:"quest_name"`
	Objectives []Objective `json:"objectives,omitempty"`
}
