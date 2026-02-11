package game

import (
	"errors"
	"stardew-agent/internal/models"
)

// QuestSystem manages quests
type QuestSystem struct {
	quests map[string]*models.QuestData
}

// NewQuestSystem creates a new quest system
func NewQuestSystem() *QuestSystem {
	qs := &QuestSystem{
		quests: make(map[string]*models.QuestData),
	}

	qs.initQuests()
	return qs
}

// initQuests initializes the default quests
func (qs *QuestSystem) initQuests() {
	qs.quests = map[string]*models.QuestData{
		"q001": {
			ID:          "q001",
			Name:        "First Harvest",
			Description: "Harvest your first crop. This will help you learn the basics of farming.",
			Objectives: []models.Objective{
				{
					ID:          "q001_obj1",
					Description: "Harvest any crop",
					Type:        "harvest",
					Target:      "any",
					Count:       1,
					Progress:    0,
					Completed:   false,
				},
			},
			Rewards: models.Reward{
				Gold: 100,
				Items: []models.InventoryItem{
					{ID: "seed_potato", Name: "Potato Seeds", Type: "seed", Quantity: 10, Price: 50},
				},
			},
			Status:    models.QuestStatusAvailable,
			Giver:     "lewis",
			Deadline:  0,
		},
		"q002": {
			ID:          "q002",
			Name:        "Potato Farmer",
			Description: "Lewis wants you to grow and harvest 5 potatoes.",
			Objectives: []models.Objective{
				{
					ID:          "q002_obj1",
					Description: "Harvest potatoes",
					Type:        "harvest",
					Target:      "potato",
					Count:       5,
					Progress:    0,
					Completed:   false,
				},
			},
			Rewards: models.Reward{
				Gold: 250,
				Items: []models.InventoryItem{
					{ID: "seed_tomato", Name: "Tomato Seeds", Type: "seed", Quantity: 15, Price: 50},
				},
			},
			Status:    models.QuestStatusAvailable,
			Giver:     "lewis",
			Deadline:  14,
		},
		"q003": {
			ID:          "q003",
			Name:        "Making Friends",
			Description: "Introduce yourself to the villagers. Go talk to Pierre at the general store.",
			Objectives: []models.Objective{
				{
					ID:          "q003_obj1",
					Description: "Talk to Pierre",
					Type:        "talk",
					Target:      "pierre",
					Count:       1,
					Progress:    0,
					Completed:   false,
				},
			},
			Rewards: models.Reward{
				Gold: 50,
			},
			Status:    models.QuestStatusAvailable,
			Giver:     "lewis",
			Deadline:  7,
		},
		"q004": {
			ID:          "q004",
			Name:        "Robin's Request",
			Description: "Robin needs some wood for her workshop. Gather 10 pieces of wood.",
			Objectives: []models.Objective{
				{
					ID:          "q004_obj1",
					Description: "Collect wood",
					Type:        "collect",
					Target:      "wood",
					Count:       10,
					Progress:    0,
					Completed:   false,
				},
			},
			Rewards: models.Reward{
				Gold: 150,
				Items: []models.InventoryItem{
					{ID: "chest", Name: "Chest", Type: "furniture", Quantity: 1, Price: 50},
				},
			},
			Status:    models.QuestStatusAvailable,
			Giver:     "robin",
			Deadline:  10,
		},
		"q005": {
			ID:          "q005",
			Name:        "A Gift for Haley",
			Description: "Haley has been looking down lately. Maybe a gift would cheer her up?",
			Objectives: []models.Objective{
				{
					ID:          "q005_obj1",
					Description: "Give Haley a gift she likes",
					Type:        "give_gift",
					Target:      "haley",
					Count:       1,
					Progress:    0,
					Completed:   false,
				},
			},
			Rewards: models.Reward{
				Gold: 100,
				Items: []models.InventoryItem{
					{ID: "flower", Name: "Tulip", Type: "gift", Quantity: 3, Price: 30},
				},
			},
			Status:    models.QuestStatusAvailable,
			Giver:     "lewis",
			Deadline:  0,
		},
		"q006": {
			ID:          "q006",
			Name:        "Corn for the Community",
			Description: "Grow and harvest 10 ears of corn for the community potluck.",
			Objectives: []models.Objective{
				{
					ID:          "q006_obj1",
					Description: "Harvest corn",
					Type:        "harvest",
					Target:      "corn",
					Count:       10,
					Progress:    0,
					Completed:   false,
				},
			},
			Rewards: models.Reward{
				Gold: 500,
				Items: []models.InventoryItem{
					{ID: "seed_pumpkin", Name: "Pumpkin Seeds", Type: "seed", Quantity: 20, Price: 100},
				},
			},
			Status:    models.QuestStatusAvailable,
			Giver:     "lewis",
			Deadline:  28,
		},
	}
}

// GetAvailableQuests returns all available quests
func (qs *QuestSystem) GetAvailableQuests() []models.QuestData {
	var quests []models.QuestData
	for _, q := range qs.quests {
		quests = append(quests, *q)
	}
	return quests
}

// GetQuest returns a quest by ID
func (qs *QuestSystem) GetQuest(id string) (*models.QuestData, error) {
	quest, ok := qs.quests[id]
	if !ok {
		return nil, errors.New("quest not found")
	}
	return quest, nil
}

// AcceptQuest marks a quest as active
func (qs *QuestSystem) AcceptQuest(id string) (*models.QuestActionResult, error) {
	quest, ok := qs.quests[id]
	if !ok {
		return nil, errors.New("quest not found")
	}

	if quest.Status != models.QuestStatusAvailable {
		return nil, errors.New("quest not available")
	}

	quest.Status = models.QuestStatusActive

	return &models.QuestActionResult{
		Success:   true,
		QuestID:   id,
		QuestName: quest.Name,
	}, nil
}

// UpdateObjective updates progress on a quest objective
func (qs *QuestSystem) UpdateObjective(questID, objectiveType, target string, amount int) bool {
	quest, ok := qs.quests[questID]
	if !ok || quest.Status != models.QuestStatusActive {
		return false
	}

	updated := false
	for i := range quest.Objectives {
		obj := &quest.Objectives[i]
		if obj.Type == objectiveType && (obj.Target == target || obj.Target == "any") {
			obj.Progress += amount
			if obj.Progress >= obj.Count {
				obj.Progress = obj.Count
				obj.Completed = true
			}
			updated = true
		}
	}

	// Check if quest is complete
	if updated {
		qs.checkQuestCompletion(questID)
	}

	return updated
}

// checkQuestCompletion checks if all objectives are completed
func (qs *QuestSystem) checkQuestCompletion(questID string) {
	quest, ok := qs.quests[questID]
	if !ok {
		return
	}

	allComplete := true
	for _, obj := range quest.Objectives {
		if !obj.Completed {
			allComplete = false
			break
		}
	}

	if allComplete {
		quest.Status = models.QuestStatusCompleted
	}
}

// CompleteQuest finalizes a completed quest and gives rewards
func (qs *QuestSystem) CompleteQuest(id string) (*models.Reward, error) {
	quest, ok := qs.quests[id]
	if !ok {
		return nil, errors.New("quest not found")
	}

	if quest.Status != models.QuestStatusCompleted {
		return nil, errors.New("quest objectives not completed")
	}

	// Return rewards
	rewards := quest.Rewards

	return &rewards, nil
}

// GetActiveQuests returns all active quests
func (qs *QuestSystem) GetActiveQuests() []models.QuestData {
	var quests []models.QuestData
	for _, q := range qs.quests {
		if q.Status == models.QuestStatusActive {
			quests = append(quests, *q)
		}
	}
	return quests
}

// GetCompletedQuests returns all completed quests
func (qs *QuestSystem) GetCompletedQuests() []models.QuestData {
	var quests []models.QuestData
	for _, q := range qs.quests {
		if q.Status == models.QuestStatusCompleted {
			quests = append(quests, *q)
		}
	}
	return quests
}

// GetQuestsFromNPC returns quests given by a specific NPC
func (qs *QuestSystem) GetQuestsFromNPC(npcID string) []models.QuestData {
	var quests []models.QuestData
	for _, q := range qs.quests {
		if q.Giver == npcID {
			quests = append(quests, *q)
		}
	}
	return quests
}

// HasActiveQuest checks if player has an active quest with given type/target
func (qs *QuestSystem) HasActiveQuest(objectiveType, target string) bool {
	for _, q := range qs.quests {
		if q.Status != models.QuestStatusActive {
			continue
		}
		for _, obj := range q.Objectives {
			if obj.Type == objectiveType && (obj.Target == target || obj.Target == "any") && !obj.Completed {
				return true
			}
		}
	}
	return false
}

// GetActiveQuestForObjective returns active quest that matches objective
func (qs *QuestSystem) GetActiveQuestForObjective(objectiveType, target string) *models.QuestData {
	for _, q := range qs.quests {
		if q.Status != models.QuestStatusActive {
			continue
		}
		for _, obj := range q.Objectives {
			if obj.Type == objectiveType && (obj.Target == target || obj.Target == "any") && !obj.Completed {
				return q
			}
		}
	}
	return nil
}
