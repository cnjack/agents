package game

import (
	"math"
	"math/rand"
	"stardew-agent/internal/models"
	"sync"
	"time"
)

// NPCInteractionSystem manages interactions between NPCs
type NPCInteractionSystem struct {
	activeInteractions map[string]*models.NPCInteraction // initiator ID -> interaction
	interactionHistory map[string][]models.NPCInteraction
	mu                 sync.RWMutex
	random             *rand.Rand
}

// NewNPCInteractionSystem creates a new NPC interaction system
func NewNPCInteractionSystem() *NPCInteractionSystem {
	return &NPCInteractionSystem{
		activeInteractions: make(map[string]*models.NPCInteraction),
		interactionHistory: make(map[string][]models.NPCInteraction),
		random:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CheckForInteraction checks if two NPCs should interact
func (nis *NPCInteractionSystem) CheckForInteraction(npc1, npc2 *models.NPCState, mood1, mood2 *models.NPCMoodState) *models.NPCInteraction {
	// Calculate distance
	dx := npc1.Position.X - npc2.Position.X
	dy := npc1.Position.Y - npc2.Position.Y
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	// Only interact if close enough
	if distance > 3 {
		return nil
	}

	// Check if either is already in an interaction
	nis.mu.RLock()
	_, busy1 := nis.activeInteractions[npc1.ID]
	_, busy2 := nis.activeInteractions[npc2.ID]
	nis.mu.RUnlock()

	if busy1 || busy2 {
		return nil
	}

	// Calculate interaction probability
	prob := nis.calculateInteractionProbability(npc1, npc2, mood1, mood2)

	if nis.random.Float64() < prob {
		return nis.createInteraction(npc1, npc2, mood1)
	}

	return nil
}

// calculateInteractionProbability calculates the probability of interaction
func (nis *NPCInteractionSystem) calculateInteractionProbability(npc1, npc2 *models.NPCState, mood1, mood2 *models.NPCMoodState) float64 {
	baseProb := 0.1 // 10% base chance per check

	// Mood modifiers
	moodMod := 0.0
	if mood1 != nil {
		switch mood1.CurrentMood {
		case models.MoodHappy, models.MoodExcited:
			moodMod += 0.1
		case models.MoodSad:
			moodMod -= 0.05 // Less likely to initiate when sad
		case models.MoodAngry:
			moodMod -= 0.15
		}
	}

	// Friendship modifier
	friendshipMod := 0.0
	if npc1.Friendship > 500 {
		friendshipMod = 0.15
	} else if npc1.Friendship > 250 {
		friendshipMod = 0.1
	} else if npc1.Friendship > 100 {
		friendshipMod = 0.05
	}

	// Time of day modifier
	timeMod := 0.0
	// Social hours boost
	// This would need access to game time, using placeholder

	return baseProb + moodMod + friendshipMod + timeMod
}

// createInteraction creates a new interaction between NPCs
func (nis *NPCInteractionSystem) createInteraction(initiator, target *models.NPCState, initiatorMood *models.NPCMoodState) *models.NPCInteraction {
	interaction := &models.NPCInteraction{
		InitiatorID:    initiator.ID,
		TargetID:       target.ID,
		InteractionType: nis.selectInteractionType(initiatorMood),
		Duration:       15 + nis.random.Intn(15), // 15-30 game minutes
	}

	// Generate dialogue based on interaction type
	interaction.Dialogue = nis.generateInteractionDialogue(interaction.InteractionType, initiatorMood)

	// Store as active
	nis.mu.Lock()
	nis.activeInteractions[initiator.ID] = interaction
	nis.mu.Unlock()

	return interaction
}

// selectInteractionType selects an interaction type based on mood
func (nis *NPCInteractionSystem) selectInteractionType(mood *models.NPCMoodState) string {
	types := []string{"talk", "greet", "joke", "share_news", "ask_about_day", "gossip"}

	if mood == nil {
		return types[nis.random.Intn(len(types))]
	}

	// Mood-based interaction preferences
	switch mood.CurrentMood {
	case models.MoodHappy, models.MoodExcited:
		happyTypes := []string{"joke", "share_news", "gossip"}
		return happyTypes[nis.random.Intn(len(happyTypes))]
	case models.MoodSad:
		sadTypes := []string{"talk", "ask_about_day"}
		return sadTypes[nis.random.Intn(len(sadTypes))]
	case models.MoodRomantic:
		return "flirt"
	default:
		return types[nis.random.Intn(len(types))]
	}
}

// generateInteractionDialogue generates dialogue for an interaction
func (nis *NPCInteractionSystem) generateInteractionDialogue(interactionType string, mood *models.NPCMoodState) string {
	dialogues := map[string][]string{
		"talk": {
			"Hey, how's your day going?",
			"Got a moment to chat?",
			"I've been meaning to talk to you.",
		},
		"greet": {
			"Hello there!",
			"Good to see you!",
			"Hey! How are things?",
		},
		"joke": {
			"Did you hear the one about the chicken?",
			"I've got a joke for you!",
			"You're going to love this one!",
		},
		"share_news": {
			"Did you hear what happened?",
			"I've got some news to share!",
			"You won't believe what I just learned!",
		},
		"ask_about_day": {
			"How has your day been?",
			"What have you been up to?",
			"Anything interesting happen today?",
		},
		"gossip": {
			"Have you heard the latest gossip?",
			"There's something I heard...",
			"Between you and me...",
		},
		"flirt": {
			"You look great today!",
			"I was hoping to run into you.",
			"Your smile brightens my day.",
		},
	}

	if options, ok := dialogues[interactionType]; ok {
		return options[nis.random.Intn(len(options))]
	}

	return "Hello!"
}

// CompleteInteraction marks an interaction as complete
func (nis *NPCInteractionSystem) CompleteInteraction(initiatorID string, result string) {
	nis.mu.Lock()
	defer nis.mu.Unlock()

	interaction, exists := nis.activeInteractions[initiatorID]
	if !exists {
		return
	}

	// Set result
	interaction.Result = result

	// Move to history
	nis.interactionHistory[initiatorID] = append(nis.interactionHistory[initiatorID], *interaction)
	if len(nis.interactionHistory[initiatorID]) > 20 {
		nis.interactionHistory[initiatorID] = nis.interactionHistory[initiatorID][len(nis.interactionHistory[initiatorID])-20:]
	}

	// Also add to target's history
	nis.interactionHistory[interaction.TargetID] = append(nis.interactionHistory[interaction.TargetID], *interaction)
	if len(nis.interactionHistory[interaction.TargetID]) > 20 {
		nis.interactionHistory[interaction.TargetID] = nis.interactionHistory[interaction.TargetID][len(nis.interactionHistory[interaction.TargetID])-20:]
	}

	// Remove from active
	delete(nis.activeInteractions, initiatorID)
}

// GetActiveInteraction returns the active interaction for an NPC
func (nis *NPCInteractionSystem) GetActiveInteraction(npcID string) *models.NPCInteraction {
	nis.mu.RLock()
	defer nis.mu.RUnlock()
	return nis.activeInteractions[npcID]
}

// IsInInteraction checks if an NPC is currently in an interaction
func (nis *NPCInteractionSystem) IsInInteraction(npcID string) bool {
	nis.mu.RLock()
	defer nis.mu.RUnlock()

	_, isInitiator := nis.activeInteractions[npcID]
	if isInitiator {
		return true
	}

	// Check if target of any active interaction
	for _, interaction := range nis.activeInteractions {
		if interaction.TargetID == npcID {
			return true
		}
	}

	return false
}

// GetRecentInteractions returns recent interactions for an NPC
func (nis *NPCInteractionSystem) GetRecentInteractions(npcID string, count int) []models.NPCInteraction {
	nis.mu.RLock()
	defer nis.mu.RUnlock()

	history := nis.interactionHistory[npcID]
	if len(history) == 0 {
		return []models.NPCInteraction{}
	}

	if len(history) <= count {
		return history
	}

	return history[len(history)-count:]
}

// TriggerInteractionWithPlayer triggers an NPC interaction with the player
func (nis *NPCInteractionSystem) TriggerInteractionWithPlayer(npc *models.NPCState, mood *models.NPCMoodState, playerPos models.Position) *models.NPCInteraction {
	// Check distance
	dx := npc.Position.X - playerPos.X
	dy := npc.Position.Y - playerPos.Y
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	if distance > 2 {
		return nil
	}

	// Check if already in interaction
	if nis.IsInInteraction(npc.ID) {
		return nil
	}

	interaction := &models.NPCInteraction{
		InitiatorID:     npc.ID,
		TargetID:        "player",
		InteractionType: "greet",
		Duration:        10 + nis.random.Intn(10),
	}

	interaction.Dialogue = nis.generatePlayerGreeting(npc, mood)

	nis.mu.Lock()
	nis.activeInteractions[npc.ID] = interaction
	nis.mu.Unlock()

	return interaction
}

// generatePlayerGreeting generates a greeting for the player
func (nis *NPCInteractionSystem) generatePlayerGreeting(npc *models.NPCState, mood *models.NPCMoodState) string {
	friendship := npc.Friendship

	var greetings []string

	if friendship > 500 {
		// Close friend
		greetings = []string{
			"Hey there, friend! Great to see you!",
			"Oh, it's you! I was hoping we'd run into each other.",
			"There's my favorite farmer!",
		}
	} else if friendship > 250 {
		// Friend
		greetings = []string{
			"Hello! Good to see you around.",
			"Hey! How's the farm coming along?",
			"Oh, hi there! Nice to see you.",
		}
	} else if friendship > 100 {
		// Acquaintance
		greetings = []string{
			"Hello there.",
			"Oh, hi. How's it going?",
			"Hey.",
		}
	} else {
		// Stranger
		greetings = []string{
			"Oh, hello. You're the new farmer, right?",
			"Hi there.",
			"...Hello.",
		}
	}

	// Mood modifier
	if mood != nil && mood.CurrentMood == models.MoodHappy {
		greetings = append(greetings,
			"What a great day! And it just got better!",
			"Oh wonderful, you're here!",
		)
	} else if mood != nil && mood.CurrentMood == models.MoodSad {
		greetings = []string{
			"Oh... hi.",
			"Hello... sorry, I'm not great company right now.",
		}
	}

	return greetings[nis.random.Intn(len(greetings))]
}

// ProcessNPCInteractions processes all potential NPC interactions
func (nis *NPCInteractionSystem) ProcessNPCInteractions(npcs []models.NPCState, moodStates map[string]*models.NPCMoodState) []*models.NPCInteraction {
	var newInteractions []*models.NPCInteraction

	for i := 0; i < len(npcs); i++ {
		for j := i + 1; j < len(npcs); j++ {
			npc1 := &npcs[i]
			npc2 := &npcs[j]

			// Skip if either is already in interaction
			if nis.IsInInteraction(npc1.ID) || nis.IsInInteraction(npc2.ID) {
				continue
			}

			mood1 := moodStates[npc1.ID]
			mood2 := moodStates[npc2.ID]

			if interaction := nis.CheckForInteraction(npc1, npc2, mood1, mood2); interaction != nil {
				newInteractions = append(newInteractions, interaction)
			}
		}
	}

	return newInteractions
}
