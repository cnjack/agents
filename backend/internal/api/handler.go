package api

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"stardew-agent/internal/game"
	"stardew-agent/internal/models"
)

// Handler contains API handlers
type Handler struct {
	engine               *game.Engine
	npcManager           *game.NPCManager
	questSystem          *game.QuestSystem
	behaviorEngine       *game.BehaviorEngine
	moodSystem           *game.MoodSystem
	dynamicSchedule      *game.DynamicScheduleSystem
	npcInteractionSystem *game.NPCInteractionSystem
}

// NewHandler creates a new handler
func NewHandler(engine *game.Engine) *Handler {
	moodSystem := game.NewMoodSystem()
	behaviorEngine := game.NewBehaviorEngine(moodSystem)

	return &Handler{
		engine:               engine,
		npcManager:           game.NewNPCManager(),
		questSystem:          game.NewQuestSystem(),
		behaviorEngine:       behaviorEngine,
		moodSystem:           moodSystem,
		dynamicSchedule:      game.NewDynamicScheduleSystem(),
		npcInteractionSystem: game.NewNPCInteractionSystem(),
	}
}

// GetState returns the full game state
func (h *Handler) GetState(c *gin.Context) {
	state := h.engine.GetState()
	c.JSON(http.StatusOK, state)
}

// GetObservation returns the agent-observable state
func (h *Handler) GetObservation(c *gin.Context) {
	obs := h.engine.GetObservation()
	c.JSON(http.StatusOK, obs)
}

// GetMap returns the map data
func (h *Handler) GetMap(c *gin.Context) {
	state := h.engine.GetState()

	// Return simplified map data
	mapData := gin.H{
		"width":     state.MapWidth,
		"height":    state.MapHeight,
		"buildings": []interface{}{},
		"areas":     []interface{}{},
	}

	c.JSON(http.StatusOK, mapData)
}

// ExecuteAction executes an agent action
func (h *Handler) ExecuteAction(c *gin.Context) {
	var action models.Action
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.processAction(&action)
	c.JSON(http.StatusOK, result)
}

// processAction processes an action and returns the result
func (h *Handler) processAction(action *models.Action) *models.ActionResult {
	switch action.Type {
	case models.ActionMove:
		return h.handleMove(action)

	case models.ActionTalk:
		return h.handleTalk(action)

	case models.ActionGiveGift:
		return h.handleGiveGift(action)

	case models.ActionAcceptQuest:
		return h.handleAcceptQuest(action)

	case models.ActionWait:
		return h.handleWait(action)

	case models.ActionSleep:
		return h.handleSleep()

	default:
		return &models.ActionResult{
			Success: false,
			Message: "Unknown action type",
		}
	}
}

// handleMove handles movement action
func (h *Handler) handleMove(action *models.Action) *models.ActionResult {
	err := h.engine.MovePlayer(action.Params.Direction)
	if err != nil {
		return &models.ActionResult{
			Success: false,
			Message: err.Error(),
		}
	}

	return &models.ActionResult{
		Success:  true,
		Message:  "Moved " + string(action.Params.Direction),
		NewState: h.engine.GetState(),
	}
}

// handleTalk handles talking to NPCs
func (h *Handler) handleTalk(action *models.Action) *models.ActionResult {
	state := h.engine.GetState()

	// Find the NPC
	var targetNPC *models.NPCState
	for i := range state.NPCs {
		if state.NPCs[i].ID == action.Params.NPC {
			targetNPC = &state.NPCs[i]
			break
		}
	}

	if targetNPC == nil {
		return &models.ActionResult{
			Success: false,
			Message: "NPC not found",
		}
	}

	// Check distance (must be within 3 tiles)
	playerPos := state.Player.Position
	npcPos := targetNPC.Position
	dx := playerPos.X - npcPos.X
	dy := playerPos.Y - npcPos.Y
	if dx*dx+dy*dy > 9 {
		return &models.ActionResult{
			Success: false,
			Message: "Too far away to talk. Move closer.",
		}
	}

	// Get dialogue
	dialogue := "..."
	if len(targetNPC.Dialogue) > 0 {
		dialogue = targetNPC.Dialogue[rand.Intn(len(targetNPC.Dialogue))]
	}

	// Update quest progress
	h.questSystem.UpdateObjective("", "talk", action.Params.NPC, 1)

	return &models.ActionResult{
		Success: true,
		Message: dialogue,
		NewState: h.engine.GetState(),
		Events: []models.GameEvent{
			{
				Type: "npc_dialog",
				Data: map[string]string{
					"npc":     targetNPC.Name,
					"dialogue": dialogue,
				},
			},
		},
	}
}

// handleGiveGift handles giving gifts to NPCs
func (h *Handler) handleGiveGift(action *models.Action) *models.ActionResult {
	state := h.engine.GetState()
	player := &state.Player

	// Check if player has the gift
	hasGift := false
	giftIndex := -1
	for i, gift := range player.Gifts {
		if gift.ID == action.Params.GiftItem && gift.Quantity > 0 {
			hasGift = true
			giftIndex = i
			break
		}
	}

	if !hasGift {
		return &models.ActionResult{
			Success: false,
			Message: "You don't have this gift",
		}
	}

	// Find the NPC
	var targetNPC *models.NPCState
	for i := range state.NPCs {
		if state.NPCs[i].ID == action.Params.NPC {
			targetNPC = &state.NPCs[i]
			break
		}
	}

	if targetNPC == nil {
		return &models.ActionResult{
			Success: false,
			Message: "NPC not found",
		}
	}

	// Check distance
	playerPos := player.Position
	npcPos := targetNPC.Position
	dx := playerPos.X - npcPos.X
	dy := playerPos.Y - npcPos.Y
	if dx*dx+dy*dy > 9 {
		return &models.ActionResult{
			Success: false,
			Message: "Too far away. Move closer.",
		}
	}

	// Remove gift from player
	player.Gifts[giftIndex].Quantity--
	if player.Gifts[giftIndex].Quantity <= 0 {
		player.Gifts = append(player.Gifts[:giftIndex], player.Gifts[giftIndex+1:]...)
	}

	// Calculate friendship gain (10-30 points based on gift)
	friendshipGain := 10 + rand.Intn(21)

	// Update friendship
	if player.Friendship == nil {
		player.Friendship = make(map[string]int)
	}
	player.Friendship[action.Params.NPC] += friendshipGain

	// Get reaction
	reactions := []string{
		"Thank you! This is wonderful!",
		"Oh, for me? How thoughtful!",
		"This is exactly what I wanted!",
		"How did you know? I love this!",
	}
	reaction := reactions[rand.Intn(len(reactions))]

	// Update quest progress
	h.questSystem.UpdateObjective("", "give_gift", action.Params.NPC, 1)

	return &models.ActionResult{
		Success: true,
		Message: reaction + " (+" + string(rune('0'+friendshipGain/10)) + " friendship)",
		NewState: h.engine.GetState(),
		Events: []models.GameEvent{
			{
				Type: "gift_given",
				Data: map[string]interface{}{
					"npc":            targetNPC.Name,
					"gift":           action.Params.GiftItem,
					"friendship_gain": friendshipGain,
				},
			},
		},
	}
}

// handleAcceptQuest handles accepting a quest
func (h *Handler) handleAcceptQuest(action *models.Action) *models.ActionResult {
	state := h.engine.GetState()

	// Find quest
	var targetQuest *models.QuestData
	for i := range state.Quests {
		if state.Quests[i].ID == action.Params.QuestID {
			targetQuest = &state.Quests[i]
			break
		}
	}

	if targetQuest == nil {
		return &models.ActionResult{
			Success: false,
			Message: "Quest not found",
		}
	}

	if targetQuest.Status != models.QuestStatusAvailable {
		return &models.ActionResult{
			Success: false,
			Message: "Quest is not available",
		}
	}

	// Accept quest
	targetQuest.Status = models.QuestStatusActive

	return &models.ActionResult{
		Success:  true,
		Message:  "Accepted quest: " + targetQuest.Name,
		NewState: h.engine.GetState(),
	}
}

// handleWait handles waiting action
func (h *Handler) handleWait(action *models.Action) *models.ActionResult {
	state := h.engine.GetState()

	// Wait for specified ticks (each tick = 1 game minute)
	for i := 0; i < action.Params.Ticks; i++ {
		state.Time.Minute++
		if state.Time.Minute >= 60 {
			state.Time.Minute = 0
			state.Time.Hour++
		}
	}

	return &models.ActionResult{
		Success:  true,
		Message:  "Waited",
		NewState: h.engine.GetState(),
	}
}

// handleSleep handles sleeping action
func (h *Handler) handleSleep() *models.ActionResult {
	h.engine.EndDay()

	return &models.ActionResult{
		Success:  true,
		Message:  "Slept and woke up on a new day",
		NewState: h.engine.GetState(),
	}
}

// ResetGame resets the game state
func (h *Handler) ResetGame(c *gin.Context) {
	state := h.engine.Reset()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Game reset",
		"state":   state,
	})
}

// GetAIDecision handles AI behavior decision requests
func (h *Handler) GetAIDecision(c *gin.Context) {
	var req models.AIDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current game state for context
	state := h.engine.GetState()

	// Ensure weather is set (default to sunny if not provided)
	if req.Weather == "" {
		req.Weather = models.WeatherSunny
	}

	// Ensure current time is set
	if req.CurrentTime.Day == 0 {
		req.CurrentTime = state.Time
	}

	// Ensure player position is set
	if req.PlayerPosition.X == 0 && req.PlayerPosition.Y == 0 {
		req.PlayerPosition = state.Player.Position
	}

	// Initialize NPC mood state if needed
	for _, npc := range state.NPCs {
		if h.moodSystem.GetMoodState(npc.ID) == nil {
			h.moodSystem.InitializeNPC(&npc, nil)
		}
		if h.dynamicSchedule.GetCurrentEntry(npc.ID, state.Time) == nil {
			h.dynamicSchedule.InitializeNPC(&npc)
		}
	}

	// Make AI decision
	response := h.behaviorEngine.MakeDecision(&req, state.NPCs)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetNPCMood returns the mood state for an NPC
func (h *Handler) GetNPCMood(c *gin.Context) {
	npcID := c.Param("npc_id")

	state := h.moodSystem.GetMoodState(npcID)
	if state == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NPC mood state not found"})
		return
	}

	personality := h.moodSystem.GetPersonality(npcID)

	c.JSON(http.StatusOK, gin.H{
		"mood_state":  state,
		"personality": personality,
	})
}

// GetNPCInteractions returns active and recent NPC interactions
func (h *Handler) GetNPCInteractions(c *gin.Context) {
	npcID := c.Query("npc_id")

	if npcID != "" {
		active := h.npcInteractionSystem.GetActiveInteraction(npcID)
		recent := h.npcInteractionSystem.GetRecentInteractions(npcID, 10)

		c.JSON(http.StatusOK, gin.H{
			"npc_id":            npcID,
			"active_interaction": active,
			"recent_interactions": recent,
		})
		return
	}

	// Return all active interactions
	c.JSON(http.StatusOK, gin.H{
		"message": "Use ?npc_id=<id> to get specific NPC interactions",
	})
}

// TriggerNPCInteraction manually triggers an NPC interaction
func (h *Handler) TriggerNPCInteraction(c *gin.Context) {
	var req struct {
		InitiatorID string `json:"initiator_id"`
		TargetID    string `json:"target_id"`
		Type        string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	state := h.engine.GetState()

	// Find initiator NPC
	var initiator *models.NPCState
	for i := range state.NPCs {
		if state.NPCs[i].ID == req.InitiatorID {
			initiator = &state.NPCs[i]
			break
		}
	}

	if initiator == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Initiator NPC not found"})
		return
	}

	moodState := h.moodSystem.GetMoodState(req.InitiatorID)

	var interaction *models.NPCInteraction
	if req.TargetID == "player" {
		interaction = h.npcInteractionSystem.TriggerInteractionWithPlayer(initiator, moodState, state.Player.Position)
	} else {
		// Find target NPC
		var target *models.NPCState
		for i := range state.NPCs {
			if state.NPCs[i].ID == req.TargetID {
				target = &state.NPCs[i]
				break
			}
		}

		if target == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target NPC not found"})
			return
		}

		targetMood := h.moodSystem.GetMoodState(req.TargetID)
		interaction = h.npcInteractionSystem.CheckForInteraction(initiator, target, moodState, targetMood)
	}

	if interaction == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "No interaction triggered (too far, already busy, or probability check failed)",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"interaction": interaction,
	})
}

// GetBehaviorFactors returns the decision factors for an NPC
func (h *Handler) GetBehaviorFactors(c *gin.Context) {
	npcID := c.Param("npc_id")

	state := h.engine.GetState()

	// Find NPC
	var npc *models.NPCState
	for i := range state.NPCs {
		if state.NPCs[i].ID == npcID {
			npc = &state.NPCs[i]
			break
		}
	}

	if npc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
		return
	}

	// Build context
	req := &models.AIDecisionRequest{
		NPCID:          npcID,
		PlayerPosition: state.Player.Position,
		Weather:        models.WeatherSunny,
		CurrentTime:    state.Time,
	}

	context := h.behaviorEngine.BuildContext(npc, req, state.NPCs)
	factors := h.behaviorEngine.CalculateDecisionFactors(context)

	c.JSON(http.StatusOK, gin.H{
		"npc_id":   npcID,
		"factors":  factors,
		"context":  context,
	})
}

// FindPathRequest represents a pathfinding request
type FindPathRequest struct {
	NPCID  string          `json:"npc_id"`
	Start  models.Position `json:"start"`
	Goal   models.Position `json:"goal"`
}

// FindPath handles pathfinding requests
func (h *Handler) FindPath(c *gin.Context) {
	var req FindPathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := h.engine.FindPath(req.NPCID, req.Start, req.Goal)

	if path == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "No path found",
			"path":    []models.Position{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"npc_id":   req.NPCID,
		"path":     path.Positions,
		"length":   len(path.Positions),
	})
}

// GetNPCMovement returns the current movement state for an NPC
func (h *Handler) GetNPCMovement(c *gin.Context) {
	npcID := c.Param("id")

	movState := h.engine.GetNPCMovementState(npcID)

	if movState == nil {
		c.JSON(http.StatusOK, gin.H{
			"npc_id":  npcID,
			"is_moving": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"npc_id":  npcID,
		"is_moving": movState.IsMoving,
		"target_pos": movState.TargetPos,
		"path":      movState.CurrentPath,
		"speed":     movState.MoveSpeed,
	})
}
