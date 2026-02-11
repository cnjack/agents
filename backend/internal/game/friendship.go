package game

import (
	"stardew-agent/internal/models"
)

// FriendshipSystem manages NPC relationships
type FriendshipSystem struct {
	// Max friendship per heart is 250 points
	// Visible hearts = friendship / 250
	maxFriendship int
}

// NewFriendshipSystem creates a new friendship system
func NewFriendshipSystem() *FriendshipSystem {
	return &FriendshipSystem{
		maxFriendship: 2500, // Max 10 hearts
	}
}

// AddFriendship adds friendship points to an NPC
func (fs *FriendshipSystem) AddFriendship(npc *models.NPCState, points int) int {
	npc.Friendship += points
	if npc.Friendship < 0 {
		npc.Friendship = 0
	}
	if npc.Friendship > npc.MaxFriendship {
		npc.Friendship = npc.MaxFriendship
	}
	return npc.Friendship
}

// GetHearts returns the number of hearts for an NPC
func (fs *FriendshipSystem) GetHearts(friendship int) int {
	return friendship / 250
}

// GetHeartPercentage returns the progress towards the next heart
func (fs *FriendshipSystem) GetHeartPercentage(friendship int) int {
	return friendship % 250
}

// GetFriendshipLevel returns the friendship level name
func (fs *FriendshipSystem) GetFriendshipLevel(friendship int) string {
	hearts := fs.GetHearts(friendship)

	switch {
	case hearts >= 10:
		return "Best Friend"
	case hearts >= 8:
		return "Close Friend"
	case hearts >= 6:
		return "Friend"
	case hearts >= 4:
		return "Friendly"
	case hearts >= 2:
		return "Acquaintance"
	case hearts >= 1:
		return "Newcomer"
	default:
		return "Stranger"
	}
}

// DailyDecay applies daily friendship decay
func (fs *FriendshipSystem) DailyDecay(npc *models.NPCState) {
	// Lose 2 points per day if no interaction
	// But not below 0
	if npc.Friendship > 0 {
		npc.Friendship -= 2
		if npc.Friendship < 0 {
			npc.Friendship = 0
		}
	}
}

// GetGiftMultiplier returns the multiplier for gift points based on friendship level
func (fs *FriendshipSystem) GetGiftMultiplier(friendship int) float64 {
	hearts := fs.GetHearts(friendship)

	if hearts >= 8 {
		return 1.5
	} else if hearts >= 6 {
		return 1.3
	} else if hearts >= 4 {
		return 1.2
	} else if hearts >= 2 {
		return 1.1
	}
	return 1.0
}

// CanGiveGift checks if a gift can be given today
// In a full implementation, this would track daily gift limits per NPC
func (fs *FriendshipSystem) CanGiveGift(npcID string) bool {
	// One gift per NPC per week in the full game
	// For simplicity, we allow one gift per day
	return true
}

// GetDialogModifier returns a modifier for NPC dialogue based on friendship
func (fs *FriendshipSystem) GetDialogModifier(friendship int) string {
	hearts := fs.GetHearts(friendship)

	switch {
	case hearts >= 8:
		return "very_friendly"
	case hearts >= 6:
		return "friendly"
	case hearts >= 4:
		return "warm"
	case hearts >= 2:
		return "neutral"
	default:
		return "distant"
	}
}
