package ai

import (
	"sync"
	"time"

	"stardew-agent/internal/models"
)

// MoodSystem manages NPC moods
type MoodSystem struct {
	moods map[string]*models.NPCMood
	mu    sync.RWMutex
}

// NewMoodSystem creates a new mood system
func NewMoodSystem() *MoodSystem {
	return &MoodSystem{
		moods: make(map[string]*models.NPCMood),
	}
}

// GetMood returns the current mood for an NPC
func (ms *MoodSystem) GetMood(npcID string) *models.NPCMood {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if mood, exists := ms.moods[npcID]; exists {
		return mood
	}

	// Return default neutral mood
	return &models.NPCMood{
		NPCID:       npcID,
		CurrentMood: string(models.MoodNeutral),
		MoodValue:   0,
		MoodFactors: []models.MoodFactor{},
	}
}

// SetMood sets the mood for an NPC
func (ms *MoodSystem) SetMood(npcID string, mood string, value int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.moods[npcID] = &models.NPCMood{
		NPCID:          npcID,
		CurrentMood:    mood,
		MoodValue:      value,
		MoodFactors:    []models.MoodFactor{},
		LastMoodChange: time.Now().Unix(),
	}
}

// AddMoodFactor adds a mood factor to an NPC
func (ms *MoodSystem) AddMoodFactor(npcID string, factor models.MoodFactor) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	mood, exists := ms.moods[npcID]
	if !exists {
		mood = &models.NPCMood{
			NPCID:       npcID,
			CurrentMood: string(models.MoodNeutral),
			MoodValue:   0,
			MoodFactors: []models.MoodFactor{},
		}
		ms.moods[npcID] = mood
	}

	// Add factor
	factor.StartTime = time.Now().Unix()
	mood.MoodFactors = append(mood.MoodFactors, factor)

	// Apply immediate impact
	mood.MoodValue += factor.Impact
	mood.MoodValue = clamp(mood.MoodValue, -100, 100)
	mood.CurrentMood = ms.valueToMood(mood.MoodValue)
	mood.LastMoodChange = time.Now().Unix()
}

// UpdateMoods updates all NPC moods (call this on game tick)
func (ms *MoodSystem) UpdateMoods(gameHour int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now().Unix()

	for npcID, mood := range ms.moods {
		// Process mood factors
		var activeFactors []models.MoodFactor
		totalImpact := 0

		for _, factor := range mood.MoodFactors {
			// Check if factor is still active
			elapsedHours := int((now - factor.StartTime) / 3600)
			if elapsedHours < factor.Duration {
				activeFactors = append(activeFactors, factor)
				totalImpact += factor.Impact
			}
		}

		mood.MoodFactors = activeFactors

		// Decay mood towards neutral
		if len(activeFactors) == 0 {
			if mood.MoodValue > 0 {
				mood.MoodValue -= 1
				if mood.MoodValue < 0 {
					mood.MoodValue = 0
				}
			} else if mood.MoodValue < 0 {
				mood.MoodValue += 1
				if mood.MoodValue > 0 {
					mood.MoodValue = 0
				}
			}
		} else {
			mood.MoodValue = clamp(totalImpact, -100, 100)
		}

		mood.CurrentMood = ms.valueToMood(mood.MoodValue)

		// Update map
		ms.moods[npcID] = mood
	}
}

// ApplyGiftMood applies mood change from receiving a gift
func (ms *MoodSystem) ApplyGiftMood(npcID string, reaction string) {
	var impact int
	var duration int

	switch reaction {
	case "love":
		impact = 40
		duration = 24
	case "like":
		impact = 25
		duration = 24
	case "neutral":
		impact = 5
		duration = 12
	case "dislike":
		impact = -30
		duration = 24
	case "hate":
		impact = -50
		duration = 48
	default:
		impact = 0
		duration = 0
	}

	if impact != 0 {
		ms.AddMoodFactor(npcID, models.MoodFactor{
			Factor:      "gift",
			Impact:      impact,
			Duration:    duration,
			Description: "收到礼物: " + reaction,
		})
	}
}

// ApplyWeatherMood applies mood change from weather
func (ms *MoodSystem) ApplyWeatherMood(npcID string, weather models.Weather, profile models.NPCDialogueProfile) {
	// Some NPCs dislike rain
	for _, dislike := range profile.Dislikes {
		if weather == models.WeatherRainy && (dislike == "rain" || dislike == "bad_weather") {
			ms.AddMoodFactor(npcID, models.MoodFactor{
				Factor:      "weather",
				Impact:      -15,
				Duration:    12,
				Description: "下雨天心情不好",
			})
			return
		}
	}

	// Willy likes rainy weather for fishing
	if npcID == "willy" && weather == models.WeatherRainy {
		ms.AddMoodFactor(npcID, models.MoodFactor{
			Factor:      "weather",
			Impact:      10,
			Duration:    12,
			Description: "雨天钓鱼的好时机",
		})
	}
}

// ApplyInteractionMood applies mood change from player interaction
func (ms *MoodSystem) ApplyInteractionMood(npcID string, interactionType string, positive bool) {
	var impact int
	var description string

	switch interactionType {
	case "talk":
		if positive {
			impact = 5
			description = "愉快的对话"
		} else {
			impact = -5
			description = "不愉快的对话"
		}
	case "gift":
		// Gift mood is handled separately
		return
	case "quest_complete":
		impact = 20
		description = "帮助完成了任务"
	case "quest_fail":
		impact = -15
		description = "任务失败"
	default:
		return
	}

	ms.AddMoodFactor(npcID, models.MoodFactor{
		Factor:      "player_interaction",
		Impact:      impact,
		Duration:    8,
		Description: description,
	})
}

// valueToMood converts a mood value to a mood string
func (ms *MoodSystem) valueToMood(value int) string {
	switch {
	case value >= 50:
		return string(models.MoodExcited)
	case value >= 20:
		return string(models.MoodHappy)
	case value >= -20:
		return string(models.MoodNeutral)
	case value >= -50:
		return string(models.MoodSad)
	default:
		return string(models.MoodAngry)
	}
}

// clamp clamps a value between min and max
func clamp(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}
