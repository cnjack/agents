package game

import (
	"errors"
	"math/rand"
	"stardew-agent/internal/models"
)

// NPCManager handles NPC interactions
type NPCManager struct {
	npcs map[string]*models.NPCState
}

// NewNPCManager creates a new NPC manager
func NewNPCManager() *NPCManager {
	return &NPCManager{
		npcs: make(map[string]*models.NPCState),
	}
}

// Register adds an NPC to the manager
func (nm *NPCManager) Register(npc *models.NPCState) {
	nm.npcs[npc.ID] = npc
}

// Get retrieves an NPC by ID
func (nm *NPCManager) Get(id string) (*models.NPCState, error) {
	npc, ok := nm.npcs[id]
	if !ok {
		return nil, errors.New("NPC not found")
	}
	return npc, nil
}

// GetNearby returns NPCs near a position
func (nm *NPCManager) GetNearby(pos models.Position, radius int) []*models.NPCState {
	var result []*models.NPCState
	for _, npc := range nm.npcs {
		dx := npc.Position.X - pos.X
		dy := npc.Position.Y - pos.Y
		distance := dx*dx + dy*dy
		if distance <= radius*radius {
			result = append(result, npc)
		}
	}
	return result
}

// GetFacing returns the NPC the player is facing
func (nm *NPCManager) GetFacing(playerPos models.Position, direction models.Direction, npcs []models.NPCState) *models.NPCState {
	facingPos := models.Position{X: playerPos.X, Y: playerPos.Y}
	switch direction {
	case models.DirectionUp:
		facingPos.Y--
	case models.DirectionDown:
		facingPos.Y++
	case models.DirectionLeft:
		facingPos.X--
	case models.DirectionRight:
		facingPos.X++
	}

	for i := range npcs {
		if npcs[i].Position.X == facingPos.X && npcs[i].Position.Y == facingPos.Y {
			return &npcs[i]
		}
	}
	return nil
}

// Talk handles talking to an NPC
func (nm *NPCManager) Talk(npc *models.NPCState) *models.SocialActionResult {
	if len(npc.Dialogue) == 0 {
		return &models.SocialActionResult{
			Success:  true,
			NPCID:    npc.ID,
			NPCName:  npc.Name,
			Dialogue: "...",
		}
	}

	// Pick a random dialogue
	dialogue := npc.Dialogue[rand.Intn(len(npc.Dialogue))]

	return &models.SocialActionResult{
		Success:  true,
		NPCID:    npc.ID,
		NPCName:  npc.Name,
		Dialogue: dialogue,
	}
}

// GiveGift handles giving a gift to an NPC
func (nm *NPCManager) GiveGift(npc *models.NPCState, item string) *models.SocialActionResult {
	reaction := nm.getGiftReaction(npc.ID, item)
	points := nm.getGiftPoints(npc.ID, item, reaction)

	// Update friendship (max 250)
	npc.Friendship += points
	if npc.Friendship > npc.MaxFriendship {
		npc.Friendship = npc.MaxFriendship
	}

	dialogue := nm.getGiftDialogue(reaction)

	return &models.SocialActionResult{
		Success:        true,
		NPCID:          npc.ID,
		NPCName:        npc.Name,
		Dialogue:       dialogue,
		FriendshipGain: points,
		NewFriendship:  npc.Friendship,
	}
}

// getGiftReaction returns how an NPC reacts to a gift
func (nm *NPCManager) getGiftReaction(npcID, item string) string {
	// Define gift preferences per NPC
	preferences := map[string]map[string]string{
		"lewis": {
			"wine":        "love",
			"hot_pepper":  "love",
			"potato":      "like",
			"tomato":      "like",
			"corn":        "neutral",
			"pumpkin":     "like",
			"flower":      "like",
		},
		"pierre": {
			"wine":        "like",
			"coffee":      "love",
			"potato":      "neutral",
			"tomato":      "like",
			"flower":      "like",
		},
		"robin": {
			"wood":        "love",
			"stone":       "like",
			"iron_ore":    "like",
			"potato":      "neutral",
			"pumpkin":     "like",
		},
		"haley": {
			"flower":      "love",
			"coconut":     "love",
			"strawberry":  "love",
			"potato":      "dislike",
			"corn":        "neutral",
		},
		"willy": {
			"fish":        "love",
			"wine":        "like",
			"potato":      "neutral",
			"pumpkin":     "like",
		},
	}

	if npcPrefs, ok := preferences[npcID]; ok {
		if reaction, ok := npcPrefs[item]; ok {
			return reaction
		}
	}
	return "neutral"
}

// getGiftPoints returns friendship points for a gift
func (nm *NPCManager) getGiftPoints(npcID, item, reaction string) int {
	switch reaction {
	case "love":
		return 80
	case "like":
		return 45
	case "neutral":
		return 20
	case "dislike":
		return -20
	case "hate":
		return -40
	default:
		return 20
	}
}

// getGiftDialogue returns dialogue for gift reaction
func (nm *NPCManager) getGiftDialogue(reaction string) string {
	switch reaction {
	case "love":
		dialogues := []string{
			"This is amazing! Thank you so much!",
			"Oh, I love this! You're so thoughtful!",
			"Wow! This is exactly what I wanted!",
		}
		return dialogues[rand.Intn(len(dialogues))]
	case "like":
		dialogues := []string{
			"Hey, this is nice! Thanks!",
			"Oh, for me? How kind of you!",
			"I appreciate this, thank you.",
		}
		return dialogues[rand.Intn(len(dialogues))]
	case "neutral":
		dialogues := []string{
			"Oh, thanks.",
			"That's... nice of you.",
			"I guess I can use this.",
		}
		return dialogues[rand.Intn(len(dialogues))]
	case "dislike":
		dialogues := []string{
			"Um... thanks, I guess?",
			"This isn't really my thing...",
			"Oh. Well, it's the thought that counts...",
		}
		return dialogues[rand.Intn(len(dialogues))]
	case "hate":
		dialogues := []string{
			"...Why would you give me this?",
			"I really don't like this at all.",
			"Please don't give me things like this.",
		}
		return dialogues[rand.Intn(len(dialogues))]
	default:
		return "Thanks."
	}
}

// GetShopForNPC returns the shop ID for an NPC
func (nm *NPCManager) GetShopForNPC(npcID string) string {
	shops := map[string]string{
		"pierre": "general_store",
		"robin":  "carpenter_shop",
		"willy":  "fish_shop",
	}
	return shops[npcID]
}

// IsShopkeeper checks if an NPC is a shopkeeper
func (nm *NPCManager) IsShopkeeper(npcID string) bool {
	shopkeepers := map[string]bool{
		"pierre": true,
		"robin":  true,
		"willy":  true,
	}
	return shopkeepers[npcID]
}
