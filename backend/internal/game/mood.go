package game

import (
	"math/rand"
	"stardew-agent/internal/models"
	"time"
)

// MoodSystem manages NPC mood states
type MoodSystem struct {
	moodStates map[string]*models.NPCMoodState
	personalities map[string]*models.NPCPersonality
	random      *rand.Rand
}

// NewMoodSystem creates a new mood system
func NewMoodSystem() *MoodSystem {
	return &MoodSystem{
		moodStates:   make(map[string]*models.NPCMoodState),
		personalities: make(map[string]*models.NPCPersonality),
		random:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// InitializeNPC sets up mood state for an NPC
func (ms *MoodSystem) InitializeNPC(npc *models.NPCState, personality *models.NPCPersonality) {
	ms.moodStates[npc.ID] = &models.NPCMoodState{
		NPCID:           npc.ID,
		CurrentMood:     models.MoodNeutral,
		MoodIntensity:   0.5,
		MoodDuration:    0,
		MoodHistory:     []models.MoodEvent{},
		MoodInfluencers: []string{},
	}

	if personality != nil {
		ms.personalities[npc.ID] = personality
	} else {
		// Default personality
		ms.personalities[npc.ID] = &models.NPCPersonality{
			NPCID:              npc.ID,
			Sociability:        0.5,
			Spontaneity:        0.3,
			WorkEthic:          0.6,
			OutdoorsPreference: 0.5,
			MorningPerson:      0.5,
			GiftAppreciation:   0.7,
			WeatherSensitivity: 0.4,
		}
	}
}

// GetMoodState returns the mood state for an NPC
func (ms *MoodSystem) GetMoodState(npcID string) *models.NPCMoodState {
	return ms.moodStates[npcID]
}

// GetPersonality returns the personality for an NPC
func (ms *MoodSystem) GetPersonality(npcID string) *models.NPCPersonality {
	return ms.personalities[npcID]
}

// UpdateMood updates NPC mood based on various factors
func (ms *MoodSystem) UpdateMood(npcID string, context *models.BehaviorContext) *models.NPCMoodState {
	state := ms.moodStates[npcID]
	if state == nil {
		return nil
	}

	personality := ms.personalities[npcID]

	// Calculate mood influences
	var influences []string
	var totalMoodShift float64
	newMood := state.CurrentMood

	// Weather influence
	if personality.WeatherSensitivity > 0.3 {
		switch context.Weather {
		case models.WeatherSunny:
			totalMoodShift += 0.1 * personality.WeatherSensitivity
			influences = append(influences, "sunny_weather")
		case models.WeatherRainy:
			totalMoodShift -= 0.15 * personality.WeatherSensitivity
			influences = append(influences, "rainy_weather")
		case models.WeatherStormy:
			totalMoodShift -= 0.25 * personality.WeatherSensitivity
			influences = append(influences, "stormy_weather")
		}
	}

	// Time of day influence
	if personality.MorningPerson > 0.5 {
		if context.Time.Hour >= 6 && context.Time.Hour < 12 {
			totalMoodShift += 0.1
			influences = append(influences, "morning_energy")
		} else if context.Time.Hour >= 22 {
			totalMoodShift -= 0.1
			influences = append(influences, "late_hour_tired")
		}
	} else {
		if context.Time.Hour >= 18 && context.Time.Hour < 24 {
			totalMoodShift += 0.1
			influences = append(influences, "evening_energy")
		}
	}

	// Friendship influence
	if context.PlayerFriendship > 500 {
		totalMoodShift += 0.15
		influences = append(influences, "close_friendship")
	} else if context.PlayerFriendship > 250 {
		totalMoodShift += 0.05
		influences = append(influences, "good_friendship")
	}

	// Social interaction need
	nearbyCount := len(context.NearbyNPCs)
	if personality.Sociability > 0.6 && nearbyCount > 0 {
		totalMoodShift += 0.1 * float64(minInt(nearbyCount, 3))
		influences = append(influences, "social_proximity")
	} else if personality.Sociability > 0.7 && nearbyCount == 0 {
		totalMoodShift -= 0.1
		influences = append(influences, "lonely")
	}

	// Determine new mood based on total shift
	intensity := state.MoodIntensity + totalMoodShift
	intensity = maxFloat(0, minFloat(1, intensity))

	// Map intensity and current factors to mood
	newMood = ms.determineMood(state.CurrentMood, intensity, influences, context)

	// Record mood change
	if newMood != state.CurrentMood {
		event := models.MoodEvent{
			Timestamp: context.Time,
			OldMood:   state.CurrentMood,
			NewMood:   newMood,
			Trigger:   influences[len(influences)-1],
			Intensity: intensity,
		}
		state.MoodHistory = append(state.MoodHistory, event)
		// Keep only last 20 mood events
		if len(state.MoodHistory) > 20 {
			state.MoodHistory = state.MoodHistory[len(state.MoodHistory)-20:]
		}
	}

	state.CurrentMood = newMood
	state.MoodIntensity = intensity
	state.MoodInfluencers = influences

	return state
}

// determineMood determines the mood based on factors
func (ms *MoodSystem) determineMood(currentMood models.Mood, intensity float64, influences []string, context *models.BehaviorContext) models.Mood {
	// Check for specific conditions
	for _, inf := range influences {
		switch inf {
		case "sunny_weather":
			if intensity > 0.7 {
				return models.MoodExcited
			}
			return models.MoodHappy
		case "stormy_weather":
			if intensity < 0.3 {
				return models.MoodAnxious
			}
			return models.MoodSad
		case "close_friendship":
			if intensity > 0.6 {
				return models.MoodHappy
			}
		case "lonely":
			return models.MoodSad
		}
	}

	// Time-based defaults
	if context.Time.Hour >= 22 || context.Time.Hour < 6 {
		return models.MoodTired
	}

	// Default based on intensity
	if intensity > 0.7 {
		return models.MoodHappy
	} else if intensity > 0.4 {
		return models.MoodNeutral
	} else if intensity > 0.2 {
		return models.MoodSad
	}

	return currentMood
}

// ApplyMoodEvent applies a specific mood-changing event
func (ms *MoodSystem) ApplyMoodEvent(npcID string, trigger string, intensityChange float64, time models.TimeState) {
	state := ms.moodStates[npcID]
	if state == nil {
		return
	}

	personality := ms.personalities[npcID]
	multiplier := 1.0

	// Personality-based multiplier
	switch trigger {
	case "gift_received":
		multiplier = personality.GiftAppreciation
	case "social_interaction":
		multiplier = personality.Sociability
	case "work_completed":
		multiplier = personality.WorkEthic
	}

	newIntensity := state.MoodIntensity + (intensityChange * multiplier)
	newIntensity = maxFloat(0, minFloat(1, newIntensity))

	// Determine new mood
	var newMood models.Mood
	if intensityChange > 0 {
		if newIntensity > 0.8 {
			newMood = models.MoodExcited
		} else if newIntensity > 0.5 {
			newMood = models.MoodHappy
		} else {
			newMood = models.MoodNeutral
		}
	} else {
		if newIntensity < 0.2 {
			newMood = models.MoodSad
		} else if newIntensity < 0.4 {
			newMood = models.MoodNeutral
		} else {
			newMood = state.CurrentMood
		}
	}

	// Record event
	event := models.MoodEvent{
		Timestamp: time,
		OldMood:   state.CurrentMood,
		NewMood:   newMood,
		Trigger:   trigger,
		Intensity: newIntensity,
	}
	state.MoodHistory = append(state.MoodHistory, event)
	if len(state.MoodHistory) > 20 {
		state.MoodHistory = state.MoodHistory[len(state.MoodHistory)-20:]
	}

	state.CurrentMood = newMood
	state.MoodIntensity = newIntensity
	state.MoodInfluencers = append(state.MoodInfluencers, trigger)
	if len(state.MoodInfluencers) > 10 {
		state.MoodInfluencers = state.MoodInfluencers[len(state.MoodInfluencers)-10:]
	}
}

// GetMoodDialogueModifier returns a modifier for dialogue based on mood
func (ms *MoodSystem) GetMoodDialogueModifier(mood models.Mood) string {
	switch mood {
	case models.MoodHappy, models.MoodExcited:
		return "cheerful"
	case models.MoodSad:
		return "melancholy"
	case models.MoodAngry:
		return "irritated"
	case models.MoodTired:
		return "weary"
	case models.MoodAnxious:
		return "nervous"
	case models.MoodRomantic:
		return "affectionate"
	default:
		return "neutral"
	}
}

// ShouldSeekInteraction determines if NPC should seek social interaction
func (ms *MoodSystem) ShouldSeekInteraction(npcID string) bool {
	state := ms.moodStates[npcID]
	personality := ms.personalities[npcID]

	if state == nil || personality == nil {
		return false
	}

	// Lonely NPCs with high sociability seek interaction
	if state.CurrentMood == models.MoodSad && personality.Sociability > 0.6 {
		return true
	}

	// Happy, social NPCs want to share their joy
	if state.CurrentMood == models.MoodHappy && personality.Sociability > 0.7 {
		return ms.random.Float64() < 0.3
	}

	return false
}

// Helper functions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
