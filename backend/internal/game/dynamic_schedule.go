package game

import (
	"math/rand"
	"stardew-agent/internal/models"
	"time"
)

// DynamicScheduleSystem manages AI-modifiable NPC schedules
type DynamicScheduleSystem struct {
	modifiedSchedules map[string][]models.DynamicScheduleEntry
	random            *rand.Rand
}

// NewDynamicScheduleSystem creates a new dynamic schedule system
func NewDynamicScheduleSystem() *DynamicScheduleSystem {
	return &DynamicScheduleSystem{
		modifiedSchedules: make(map[string][]models.DynamicScheduleEntry),
		random:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// InitializeNPC sets up dynamic schedule tracking for an NPC
func (dss *DynamicScheduleSystem) InitializeNPC(npc *models.NPCState) {
	var entries []models.DynamicScheduleEntry
	for _, entry := range npc.Schedule {
		entries = append(entries, models.DynamicScheduleEntry{
			BaseSchedule: entry,
			Modified:     false,
			Priority:     0,
		})
	}
	dss.modifiedSchedules[npc.ID] = entries
}

// GetCurrentEntry returns the current schedule entry for an NPC
func (dss *DynamicScheduleSystem) GetCurrentEntry(npcID string, currentTime models.TimeState) *models.DynamicScheduleEntry {
	entries := dss.modifiedSchedules[npcID]
	if len(entries) == 0 {
		return nil
	}

	// Find the most recent entry that has passed
	var current *models.DynamicScheduleEntry
	for i := range entries {
		entry := &entries[i]
		entryMinutes := entry.BaseSchedule.Hour*60 + entry.BaseSchedule.Minute
		currentMinutes := currentTime.Hour*60 + currentTime.Minute

		if currentMinutes >= entryMinutes {
			current = entry
		}
	}

	return current
}

// GetNextEntry returns the next schedule entry for an NPC
func (dss *DynamicScheduleSystem) GetNextEntry(npcID string, currentTime models.TimeState) *models.DynamicScheduleEntry {
	entries := dss.modifiedSchedules[npcID]
	if len(entries) == 0 {
		return nil
	}

	currentMinutes := currentTime.Hour*60 + currentTime.Minute

	for i := range entries {
		entry := &entries[i]
		entryMinutes := entry.BaseSchedule.Hour*60 + entry.BaseSchedule.Minute
		if entryMinutes > currentMinutes {
			return entry
		}
	}

	// No more entries today - return first entry for tomorrow
	return &entries[0]
}

// ModifySchedule applies a temporary modification to an NPC's schedule
func (dss *DynamicScheduleSystem) ModifySchedule(npcID string, decision *models.BehaviorDecision, currentTime models.TimeState) {
	entries := dss.modifiedSchedules[npcID]
	if len(entries) == 0 {
		return
	}

	currentMinutes := currentTime.Hour*60 + currentTime.Minute

	// Find current or next entry to modify
	for i := range entries {
		entry := &entries[i]
		entryMinutes := entry.BaseSchedule.Hour*60 + entry.BaseSchedule.Minute

		if entryMinutes >= currentMinutes && entryMinutes <= currentMinutes+decision.Duration {
			// This entry overlaps with our modification
			if decision.TargetLocation != nil {
				entry.Modified = true
				entry.ModificationReason = decision.Reason
				entry.OverridePosition = decision.TargetLocation
				entry.Priority = decision.Priority
			}
			break
		}
	}
}

// ApplyDecision applies a behavior decision to modify schedule
func (dss *DynamicScheduleSystem) ApplyDecision(npcID string, decision *models.BehaviorDecision, currentTime models.TimeState) *models.Position {
	if !decision.DeviatesFromSchedule {
		return nil
	}

	// Get current position from schedule
	current := dss.GetCurrentEntry(npcID, currentTime)
	if current == nil {
		return nil
	}

	// Check if decision should override schedule
	if decision.Priority > current.Priority || dss.random.Float64() < 0.7 {
		dss.ModifySchedule(npcID, decision, currentTime)
		if decision.TargetLocation != nil {
			return decision.TargetLocation
		}
	}

	return nil
}

// ResetDaily resets schedule modifications for a new day
func (dss *DynamicScheduleSystem) ResetDaily(npcID string, baseSchedule []models.ScheduleEntry) {
	var entries []models.DynamicScheduleEntry
	for _, entry := range baseSchedule {
		entries = append(entries, models.DynamicScheduleEntry{
			BaseSchedule: entry,
			Modified:     false,
			Priority:     0,
		})
	}
	dss.modifiedSchedules[npcID] = entries
}

// ShouldDeviatFromSchedule determines if NPC should deviate from schedule
func (dss *DynamicScheduleSystem) ShouldDeviateFromSchedule(npcID string, personality *models.NPCPersonality, mood *models.NPCMoodState) bool {
	// High spontaneity + good mood = more likely to deviate
	baseChance := personality.Spontaneity * 0.3

	// Mood modifier
	switch mood.CurrentMood {
	case models.MoodExcited:
		baseChance += 0.2
	case models.MoodHappy:
		baseChance += 0.1
	case models.MoodSad, models.MoodAnxious:
		baseChance -= 0.1
	}

	return dss.random.Float64() < baseChance
}

// GetTimeUntilNextEvent returns minutes until next schedule event
func (dss *DynamicScheduleSystem) GetTimeUntilNextEvent(npcID string, currentTime models.TimeState) int {
	next := dss.GetNextEntry(npcID, currentTime)
	if next == nil {
		return -1
	}

	nextMinutes := next.BaseSchedule.Hour*60 + next.BaseSchedule.Minute
	currentMinutes := currentTime.Hour*60 + currentTime.Minute

	if nextMinutes > currentMinutes {
		return nextMinutes - currentMinutes
	}

	// Next event is tomorrow
	return (24*60 - currentMinutes) + nextMinutes
}

// InsertTemporaryEntry inserts a temporary schedule entry
func (dss *DynamicScheduleSystem) InsertTemporaryEntry(npcID string, position models.Position, location string, duration int, currentTime models.TimeState) {
	entries := dss.modifiedSchedules[npcID]

	// Create temporary entry
	tempEntry := models.DynamicScheduleEntry{
		BaseSchedule: models.ScheduleEntry{
			Hour:     currentTime.Hour,
			Minute:   currentTime.Minute,
			Location: location,
			Position: position,
		},
		Modified:           true,
		ModificationReason: "AI decision override",
		Priority:           10, // Higher than normal
	}

	// Find position to insert
	insertPos := 0
	currentMinutes := currentTime.Hour*60 + currentTime.Minute

	for i, entry := range entries {
		entryMinutes := entry.BaseSchedule.Hour*60 + entry.BaseSchedule.Minute
		if entryMinutes > currentMinutes {
			insertPos = i
			break
		}
		insertPos = i + 1
	}

	// Insert entry
	if insertPos >= len(entries) {
		entries = append(entries, tempEntry)
	} else {
		entries = append(entries[:insertPos], append([]models.DynamicScheduleEntry{tempEntry}, entries[insertPos:]...)...)
	}

	dss.modifiedSchedules[npcID] = entries
}

// GetNPCPosition returns the current position an NPC should be at
func (dss *DynamicScheduleSystem) GetNPCPosition(npcID string, currentTime models.TimeState) *models.Position {
	current := dss.GetCurrentEntry(npcID, currentTime)
	if current == nil {
		return nil
	}

	if current.Modified && current.OverridePosition != nil {
		return current.OverridePosition
	}

	return &current.BaseSchedule.Position
}

// GetNPCLocation returns the current location an NPC should be at
func (dss *DynamicScheduleSystem) GetNPCLocation(npcID string, currentTime models.TimeState) string {
	current := dss.GetCurrentEntry(npcID, currentTime)
	if current == nil {
		return ""
	}

	if current.Modified && current.OverrideLocation != "" {
		return current.OverrideLocation
	}

	return current.BaseSchedule.Location
}
