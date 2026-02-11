package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"stardew-agent/internal/game"
	"stardew-agent/internal/models"
	"stardew-agent/internal/services/ai"
)

// AIHandler contains AI-related API handlers
type AIHandler struct {
	engine          *game.Engine
	dialogueService *ai.DialogueService
	moodSystem      *ai.MoodSystem
	nlpService      *ai.NLPService
	profiles        map[string]models.NPCDialogueProfile
}

// NewAIHandler creates a new AI handler
func NewAIHandler(engine *game.Engine) *AIHandler {
	// Get dialogue profiles from game package
	profiles := game.GetAllDialogueProfiles()

	// Create mock LLM client for now (can be replaced with real Claude client)
	llmClient := ai.NewMockLLMClient(profiles)

	return &AIHandler{
		engine:          engine,
		dialogueService: ai.NewDialogueService(llmClient),
		moodSystem:      ai.NewMoodSystem(),
		nlpService:      ai.NewNLPService(llmClient),
		profiles:        profiles,
	}
}

// DialogueRequest represents a dialogue generation request
type DialogueRequest struct {
	NPCID           string `json:"npc_id" binding:"required"`
	PlayerInput     string `json:"player_input"`
	NaturalLanguage string `json:"natural_language"` // Alternative to structured input
}

// DialogueResponse represents the dialogue API response
type DialogueResponse struct {
	Success bool                      `json:"success"`
	Data    *models.AIDialogueResponse `json:"data,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// GenerateDialogue handles POST /api/v1/ai/dialogue
func (h *AIHandler) GenerateDialogue(c *gin.Context) {
	var req DialogueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, DialogueResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get NPC profile
	profile, exists := h.profiles[req.NPCID]
	if !exists {
		c.JSON(http.StatusNotFound, DialogueResponse{
			Success: false,
			Error:   "NPC not found",
		})
		return
	}

	// Get game state
	gameState := h.engine.GetState()

	// Get current mood
	mood := h.moodSystem.GetMood(req.NPCID)

	// Build dialogue context
	dialogueCtx := ai.BuildDialogueContext(
		req.NPCID,
		profile,
		gameState,
		req.PlayerInput,
		[]models.Interaction{}, // TODO: Track interactions
		mood.CurrentMood,
	)

	// Generate dialogue with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.dialogueService.GenerateDialogue(ctx, dialogueCtx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DialogueResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Update mood if changed
	if response.MoodChange != "" {
		h.moodSystem.SetMood(req.NPCID, response.MoodChange, mood.MoodValue)
	}

	c.JSON(http.StatusOK, DialogueResponse{
		Success: true,
		Data:    response,
	})
}

// MoodResponse represents the mood API response
type MoodResponse struct {
	Success bool             `json:"success"`
	Data    *models.NPCMood  `json:"data,omitempty"`
	Error   string           `json:"error,omitempty"`
}

// GetNPCMood handles GET /api/v1/npc/:id/mood
func (h *AIHandler) GetNPCMood(c *gin.Context) {
	npcID := c.Param("id")

	mood := h.moodSystem.GetMood(npcID)

	c.JSON(http.StatusOK, MoodResponse{
		Success: true,
		Data:    mood,
	})
}

// NLPInterpretRequest represents a natural language interpretation request
type NLPInterpretRequest struct {
	PlayerInput string `json:"player_input" binding:"required"`
	Context     string `json:"context"`
	SessionID   string `json:"session_id"` // For maintaining conversation context
	UseLLM      bool   `json:"use_llm"`    // Whether to use LLM for processing
}

// NLPInterpretResponse represents the NLP interpretation response
type NLPInterpretResponse struct {
	Success               bool                   `json:"success"`
	Data                  *models.NLPResponse    `json:"data,omitempty"`
	EnhancedResult        *ai.NLPResult          `json:"enhanced_result,omitempty"`
	Error                 string                 `json:"error,omitempty"`
}

// InterpretNaturalLanguage handles POST /api/v1/ai/interpret
func (h *AIHandler) InterpretNaturalLanguage(c *gin.Context) {
	var req NLPInterpretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NLPInterpretResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	gameState := h.engine.GetState()
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = "default"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use enhanced NLP service
	var nlpResult *ai.NLPResult
	var err error

	if req.UseLLM {
		nlpResult, err = h.nlpService.ProcessWithLLM(ctx, req.PlayerInput, gameState, sessionID)
	} else {
		nlpResult, err = h.nlpService.Process(ctx, req.PlayerInput, gameState, sessionID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, NLPInterpretResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert to legacy NLPResponse format for backward compatibility
	nlpResponse := &models.NLPResponse{
		Understood:            nlpResult.Understood,
		Intent:                nlpResult.Intent,
		Action:                nlpResult.Action,
		TargetNPC:             nlpResult.TargetNPC,
		TargetItem:            nlpResult.TargetItem,
		ResponseText:          nlpResult.ResponseText,
		NeedsClarification:    nlpResult.NeedsClarification,
		ClarificationQuestion: nlpResult.ClarificationQuestion,
	}

	c.JSON(http.StatusOK, NLPInterpretResponse{
		Success:        true,
		Data:           nlpResponse,
		EnhancedResult: nlpResult,
	})
}

// interpretInput performs simple natural language interpretation
func (h *AIHandler) interpretInput(req models.NLPRequest) *models.NLPResponse {
	input := req.PlayerInput

	// Simple pattern matching for demo
	// In production, use LLM for interpretation

	// Check for greeting + NPC name patterns
	for _, npc := range req.NearbyNPCs {
		if containsAny(input, []string{npc.Name, npc.ID}) {
			if containsAny(input, []string{"你好", "hi", "hello", "嗨", "hey"}) {
				return &models.NLPResponse{
					Understood:   true,
					Intent:       "talk",
					TargetNPC:    npc.ID,
					ResponseText: "你想和 " + npc.Name + " 对话",
					Action: &models.Action{
						Type: models.ActionTalk,
						Params: models.ActionParams{
							NPC: npc.ID,
						},
					},
				}
			}

			if containsAny(input, []string{"送", "给", "gift", "give"}) {
				// Find gift item in inventory
				if len(req.PlayerInventory) > 0 {
					return &models.NLPResponse{
						Understood:   true,
						Intent:       "gift",
						TargetNPC:    npc.ID,
						TargetItem:   req.PlayerInventory[0].ID,
						ResponseText: "你想把 " + req.PlayerInventory[0].Name + " 送给 " + npc.Name,
						Action: &models.Action{
							Type: models.ActionGiveGift,
							Params: models.ActionParams{
								NPC:       npc.ID,
								GiftItem: req.PlayerInventory[0].ID,
							},
						},
					}
				}
			}
		}
	}

	// Check for movement patterns
	if containsAny(input, []string{"去", "走", "move", "go"}) {
		directions := map[string]models.Direction{
			"上":    models.DirectionUp,
			"下":    models.DirectionDown,
			"左":    models.DirectionLeft,
			"右":    models.DirectionRight,
			"up":   models.DirectionUp,
			"down": models.DirectionDown,
			"left": models.DirectionLeft,
			"right": models.DirectionRight,
		}

		for dirWord, dir := range directions {
			if contains(input, dirWord) {
				return &models.NLPResponse{
					Understood:   true,
					Intent:       "move",
					ResponseText: "向 " + dirWord + " 移动",
					Action: &models.Action{
						Type: models.ActionMove,
						Params: models.ActionParams{
							Direction: dir,
						},
					},
				}
			}
		}
	}

	// Check for quest patterns
	if containsAny(input, []string{"任务", "quest", "接受", "accept"}) {
		// Find available quest from nearby NPC
		activeQuests := h.engine.GetState().Quests
		for _, quest := range activeQuests {
			if quest.Status == models.QuestStatusAvailable {
				return &models.NLPResponse{
					Understood:   true,
					Intent:       "accept_quest",
					ResponseText: "接受任务: " + quest.Name,
					Action: &models.Action{
						Type: models.ActionAcceptQuest,
						Params: models.ActionParams{
							QuestID: quest.ID,
						},
					},
				}
			}
		}
	}

	// Default: try to talk to nearest NPC
	if len(req.NearbyNPCs) > 0 {
		return &models.NLPResponse{
			Understood:   true,
			Intent:       "talk",
			TargetNPC:    req.NearbyNPCs[0].ID,
			ResponseText: "你想和附近的 " + req.NearbyNPCs[0].Name + " 对话",
			Action: &models.Action{
				Type: models.ActionTalk,
				Params: models.ActionParams{
					NPC: req.NearbyNPCs[0].ID,
				},
			},
		}
	}

	return &models.NLPResponse{
		Understood:         false,
		Intent:             "unknown",
		ResponseText:       "我不太明白你的意思",
		NeedsClarification: true,
		ClarificationQuestion: "你想做什么？你可以:\n- 和NPC对话\n- 送礼物\n- 移动\n- 接受任务",
	}
}

// Helper functions

func getNearbyNPCs(state *models.GameState) []models.NPCState {
	playerPos := state.Player.Position
	var nearby []models.NPCState

	for _, npc := range state.NPCs {
		dx := npc.Position.X - playerPos.X
		dy := npc.Position.Y - playerPos.Y
		if dx*dx+dy*dy <= 25 { // Within 5 tiles
			nearby = append(nearby, npc)
		}
	}

	return nearby
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

// ExecuteNLCommandRequest represents a request to execute a natural language command
type ExecuteNLCommandRequest struct {
	PlayerInput string `json:"player_input" binding:"required"`
	SessionID   string `json:"session_id"`
	AutoExecute bool   `json:"auto_execute"` // Whether to automatically execute the action
}

// ExecuteNLCommandResponse represents the response from executing a natural language command
type ExecuteNLCommandResponse struct {
	Success        bool                   `json:"success"`
	Interpretation *ai.NLPResult          `json:"interpretation"`
	ActionExecuted bool                   `json:"action_executed"`
	GameResponse   string                 `json:"game_response"`
	GameState      *models.GameState      `json:"game_state,omitempty"`
	Suggestions    []string               `json:"suggestions"`
	Error          string                 `json:"error,omitempty"`
}

// ExecuteNLCommand handles POST /api/v1/ai/execute
// It interprets natural language and optionally executes the action
func (h *AIHandler) ExecuteNLCommand(c *gin.Context) {
	var req ExecuteNLCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecuteNLCommandResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	gameState := h.engine.GetState()
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = "default"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Process natural language
	nlpResult, err := h.nlpService.Process(ctx, req.PlayerInput, gameState, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecuteNLCommandResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := ExecuteNLCommandResponse{
		Success:        true,
		Interpretation: nlpResult,
		Suggestions:    nlpResult.Suggestions,
		GameResponse:   nlpResult.ResponseText,
	}

	// Execute action if requested and action exists
	if req.AutoExecute && nlpResult.Action != nil {
		actionResult := h.executeAction(nlpResult.Action)
		response.ActionExecuted = actionResult.Success
		response.GameResponse = actionResult.Message
		response.GameState = h.engine.GetState()
	}

	c.JSON(http.StatusOK, response)
}

// executeAction executes an action and returns the result
func (h *AIHandler) executeAction(action *models.Action) *models.ActionResult {
	// This is a simplified version - in production, this would use the game engine
	// to execute the action properly

	result := &models.ActionResult{
		Success: true,
		Message: "Action executed",
	}

	switch action.Type {
	case models.ActionMove:
		result.Message = fmt.Sprintf("Moved %s", action.Params.Direction)
	case models.ActionTalk:
		result.Message = fmt.Sprintf("Talked to NPC: %s", action.Params.NPC)
	case models.ActionGiveGift:
		result.Message = fmt.Sprintf("Gave %s to %s", action.Params.GiftItem, action.Params.NPC)
	case models.ActionAcceptQuest:
		result.Message = fmt.Sprintf("Accepted quest: %s", action.Params.QuestID)
	default:
		result.Message = "Unknown action"
	}

	return result
}

// ComplexInteractionRequest represents a complex interaction request
type ComplexInteractionRequest struct {
	Type        string `json:"type" binding:"required"`        // bargain, express_emotion, negotiate
	NPCID       string `json:"npc_id" binding:"required"`
	Content     string `json:"content"`                        // The interaction content
	Amount      int    `json:"amount,omitempty"`               // For bargaining/trading
	Emotion     string `json:"emotion,omitempty"`              // For expressing emotion
	SessionID   string `json:"session_id"`
}

// ComplexInteractionResponse represents the response to a complex interaction
type ComplexInteractionResponse struct {
	Success           bool                      `json:"success"`
	NPCResponse       string                    `json:"npc_response"`
	MoodChange        string                    `json:"mood_change,omitempty"`
	FriendshipDelta   int                       `json:"friendship_delta"`
	UpdatedPrice      int                       `json:"updated_price,omitempty"`    // For bargaining
	NextOptions       []string                  `json:"next_options,omitempty"`
	TriggeredEvent    string                    `json:"triggered_event,omitempty"`
}

// HandleComplexInteraction handles complex interactions like bargaining
func (h *AIHandler) HandleComplexInteraction(c *gin.Context) {
	var req ComplexInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ComplexInteractionResponse{
			Success: false,
		})
		return
	}

	gameState := h.engine.GetState()
	profile, exists := h.profiles[req.NPCID]
	if !exists {
		c.JSON(http.StatusNotFound, ComplexInteractionResponse{
			Success: false,
		})
		return
	}

	response := ComplexInteractionResponse{
		Success: true,
	}

	switch req.Type {
	case "bargain":
		response = h.handleBargain(req, profile, gameState)
	case "express_emotion":
		response = h.handleEmotion(req, profile, gameState)
	case "negotiate":
		response = h.handleNegotiate(req, profile, gameState)
	default:
		response.NPCResponse = "我不太明白你的意思。"
	}

	// Update mood if changed
	if response.MoodChange != "" {
		h.moodSystem.SetMood(req.NPCID, response.MoodChange, 0)
	}

	c.JSON(http.StatusOK, response)
}

// handleBargain handles bargaining interactions
func (h *AIHandler) handleBargain(req ComplexInteractionRequest, profile models.NPCDialogueProfile, gameState *models.GameState) ComplexInteractionResponse {
	response := ComplexInteractionResponse{Success: true}

	// Calculate discount based on friendship
	friendship := 0
	for _, npc := range gameState.NPCs {
		if npc.ID == req.NPCID {
			friendship = npc.Friendship
			break
		}
	}

	discount := friendship / 500 // 1% discount per 500 friendship points, max ~5%

	if discount > 20 {
		discount = 20 // Max 20% discount
	}

	if discount > 0 {
		response.NPCResponse = fmt.Sprintf("好吧，看在我们是朋友的份上，给你打%d%%折扣。", discount)
		response.MoodChange = "friendly"
		response.FriendshipDelta = 5
		response.NextOptions = []string{"接受价格", "继续砍价", "放弃购买"}
	} else {
		response.NPCResponse = "抱歉，这已经是最优惠的价格了。"
		response.FriendshipDelta = -5
		response.NextOptions = []string{"接受价格", "放弃购买"}
	}

	return response
}

// handleEmotion handles emotional expression interactions
func (h *AIHandler) handleEmotion(req ComplexInteractionRequest, profile models.NPCDialogueProfile, gameState *models.GameState) ComplexInteractionResponse {
	response := ComplexInteractionResponse{Success: true}

	// Find matching traits for response style
	emotionResponse := map[string]string{
		"happy":    "你看起来心情不错！今天有什么好事吗？",
		"sad":      "怎么了？看你不太开心，需要聊聊吗？",
		"angry":    "冷静一下，发生什么事了？",
		"grateful": "不客气！能帮到你我也很开心。",
		"love":     "...你、你说什么呢！（脸红）",
		"hate":     "呃...好吧，我知道了。",
	}

	if resp, ok := emotionResponse[req.Emotion]; ok {
		response.NPCResponse = resp
	} else {
		response.NPCResponse = "我理解你的感受。"
	}

	response.MoodChange = "friendly"
	response.FriendshipDelta = 10

	return response
}

// handleNegotiate handles negotiation interactions
func (h *AIHandler) handleNegotiate(req ComplexInteractionRequest, profile models.NPCDialogueProfile, gameState *models.GameState) ComplexInteractionResponse {
	response := ComplexInteractionResponse{Success: true}

	// Simple negotiation logic
	if req.Amount > 0 {
		if req.Amount < 100 {
			response.NPCResponse = fmt.Sprintf("%d金币？嗯...可以接受。", req.Amount)
			response.UpdatedPrice = req.Amount
		} else {
			response.NPCResponse = "这个价格有点高，我们能再商量一下吗？"
			response.NextOptions = []string{"同意降低", "坚持原价", "放弃交易"}
		}
	} else {
		response.NPCResponse = "你想出多少钱？"
	}

	return response
}
