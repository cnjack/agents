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
	engine     *game.Engine
	aiManager  *ai.ServiceManager
	moodSystem *ai.MoodSystem
	profiles   map[string]models.NPCDialogueProfile

	// Legacy services (for backward compatibility)
	dialogueService *ai.DialogueService
	nlpService      *ai.NLPService
}

// NewAIHandler creates a new AI handler (legacy, uses mock)
func NewAIHandler(engine *game.Engine) *AIHandler {
	profiles := game.GetAllDialogueProfiles()
	llmClient := ai.NewMockLLMClient(profiles)

	return &AIHandler{
		engine:          engine,
		moodSystem:      ai.NewMoodSystem(),
		profiles:        profiles,
		dialogueService: ai.NewDialogueService(llmClient),
		nlpService:      ai.NewNLPService(llmClient),
	}
}

// NewAIHandlerWithManager creates a new AI handler with service manager
func NewAIHandlerWithManager(engine *game.Engine, manager *ai.ServiceManager) *AIHandler {
	profiles := game.GetAllDialogueProfiles()

	// Still create legacy services for backward compatibility
	llmClient := ai.NewMockLLMClient(profiles)

	return &AIHandler{
		engine:          engine,
		aiManager:       manager,
		moodSystem:      ai.NewMoodSystem(),
		profiles:        profiles,
		dialogueService: ai.NewDialogueService(llmClient),
		nlpService:      ai.NewNLPService(llmClient),
	}
}

// DialogueRequest represents a dialogue generation request
type DialogueRequest struct {
	NPCID           string `json:"npc_id" binding:"required"`
	PlayerInput     string `json:"player_input"`
	NaturalLanguage string `json:"natural_language"`
}

// DialogueAPIResponse represents the dialogue API response
type DialogueAPIResponse struct {
	Success bool                      `json:"success"`
	Data    *models.AIDialogueResponse `json:"data,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// GenerateDialogue handles POST /api/v1/ai/dialogue
func (h *AIHandler) GenerateDialogue(c *gin.Context) {
	var req DialogueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, DialogueAPIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get NPC profile
	profile, exists := h.profiles[req.NPCID]
	if !exists {
		c.JSON(http.StatusNotFound, DialogueAPIResponse{
			Success: false,
			Error:   "NPC not found",
		})
		return
	}

	// Get game state
	gameState := h.engine.GetState()

	// Use new AI manager if available
	if h.aiManager != nil {
		h.generateDialogueWithManager(c, req, profile, gameState)
		return
	}

	// Fallback to legacy implementation
	h.generateDialogueLegacy(c, req, profile, gameState)
}

func (h *AIHandler) generateDialogueWithManager(c *gin.Context, req DialogueRequest, profile models.NPCDialogueProfile, gameState *models.GameState) {
	// Get current mood
	mood := h.moodSystem.GetMood(req.NPCID)

	// Calculate friendship hearts
	hearts := 0
	for _, npc := range gameState.NPCs {
		if npc.ID == req.NPCID {
			hearts = npc.Friendship / 250
			break
		}
	}

	// Convert profile to character for new API
	char := convertProfileToCharacter(req.NPCID, profile)

	// Build dialogue request for new manager
	dialogueReq := ai.DialogueRequest{
		Character:        char,
		PlayerInput:      req.PlayerInput,
		CurrentMood:      mood.CurrentMood,
		FriendshipHearts: hearts,
		GameState:        gameState,
	}

	// Generate dialogue with timeout (use longer timeout for AI calls)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	response, err := h.aiManager.GenerateDialogue(ctx, dialogueReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DialogueAPIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Update mood if changed
	if response.MoodChange != "" {
		h.moodSystem.SetMood(req.NPCID, response.MoodChange, mood.MoodValue)
	}

	// Convert to legacy response format
	aiResp := &models.AIDialogueResponse{
		NPCID:            req.NPCID,
		Dialogue:         response.Dialogue,
		MoodChange:       response.MoodChange,
		FriendshipDelta:  response.FriendshipDelta,
		SuggestedActions: response.SuggestedActions,
		SecretRevealed:   response.SecretRevealed,
	}

	c.JSON(http.StatusOK, DialogueAPIResponse{
		Success: true,
		Data:    aiResp,
	})
}

func (h *AIHandler) generateDialogueLegacy(c *gin.Context, req DialogueRequest, profile models.NPCDialogueProfile, gameState *models.GameState) {
	// Get current mood
	mood := h.moodSystem.GetMood(req.NPCID)

	// Build dialogue context
	dialogueCtx := ai.BuildDialogueContext(
		req.NPCID,
		profile,
		gameState,
		req.PlayerInput,
		[]models.Interaction{},
		mood.CurrentMood,
	)

	// Generate dialogue with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.dialogueService.GenerateDialogue(ctx, dialogueCtx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DialogueAPIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Update mood if changed
	if response.MoodChange != "" {
		h.moodSystem.SetMood(req.NPCID, response.MoodChange, mood.MoodValue)
	}

	c.JSON(http.StatusOK, DialogueAPIResponse{
		Success: true,
		Data:    response,
	})
}

// convertProfileToCharacter converts NPCDialogueProfile to Character
func convertProfileToCharacter(id string, profile models.NPCDialogueProfile) *models.Character {
	return &models.Character{
		ID:              id,
		Name:            profile.Name,
		Role:            profile.Role,
		Age:             profile.Age,
		Traits:          profile.Traits,
		Background:      profile.Background,
		Values:          profile.Values,
		Dislikes:        profile.Dislikes,
		Fears:           profile.Fears,
		Hopes:           profile.Hopes,
		SpeechStyle:     profile.SpeechStyle,
		Hobby:           profile.Hobby,
		GiftPreferences: profile.GiftPreferences,
		DialogueThemes:  profile.DialogueThemes,
		Secrets:         profile.Secrets,
		MaxFriendship:   2500,
	}
}

// MoodResponse represents the mood API response
type MoodResponse struct {
	Success bool            `json:"success"`
	Data    *models.NPCMood `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
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

// ListProvidersResponse represents the list providers response
type ListProvidersResponse struct {
	Success  bool                        `json:"success"`
	Providers map[string]ai.ProviderInfo `json:"providers"`
	Default  string                      `json:"default"`
}

// ListProviders handles GET /api/v1/ai/providers
func (h *AIHandler) ListProviders(c *gin.Context) {
	if h.aiManager == nil {
		c.JSON(http.StatusOK, ListProvidersResponse{
			Success:   true,
			Providers: map[string]ai.ProviderInfo{"mock": {Type: ai.ProviderTypeMock, Model: "mock", Available: true}},
			Default:   "mock",
		})
		return
	}

	c.JSON(http.StatusOK, ListProvidersResponse{
		Success:   true,
		Providers: h.aiManager.ListProviders(),
		Default:   "", // Could store this in handler
	})
}

// CheckProviderHealth handles GET /api/v1/ai/providers/:name/health
func (h *AIHandler) CheckProviderHealth(c *gin.Context) {
	name := c.Param("name")

	if h.aiManager == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"provider": name,
			"healthy": name == "mock",
		})
		return
	}

	provider, err := h.aiManager.GetProvider(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	healthErr := provider.HealthCheck(ctx)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"provider": name,
		"healthy":  healthErr == nil,
		"error":    healthErr,
		"info":     provider.ProviderInfo(),
	})
}

// NLPInterpretRequest represents a natural language interpretation request
type NLPInterpretRequest struct {
	PlayerInput string `json:"player_input" binding:"required"`
	Context     string `json:"context"`
	SessionID   string `json:"session_id"`
	UseLLM      bool   `json:"use_llm"`
}

// NLPInterpretResponse represents the NLP interpretation response
type NLPInterpretResponse struct {
	Success        bool                `json:"success"`
	Data           *models.NLPResponse `json:"data,omitempty"`
	EnhancedResult *ai.NLPResult       `json:"enhanced_result,omitempty"`
	Error          string              `json:"error,omitempty"`
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

	// Use new AI manager if available and LLM is requested
	if h.aiManager != nil && req.UseLLM {
		h.interpretWithManager(c, req, gameState)
		return
	}

	// Fallback to legacy implementation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nlpResult, err := h.nlpService.Process(ctx, req.PlayerInput, gameState, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NLPInterpretResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

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

func (h *AIHandler) interpretWithManager(c *gin.Context, req NLPInterpretRequest, gameState *models.GameState) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nluReq := ai.NLURequest{
		Input:      req.PlayerInput,
		GameState:  gameState,
		Context:    req.Context,
	}

	response, err := h.aiManager.Understand(ctx, nluReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NLPInterpretResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert ai.NLUResponse to models.NLPResponse
	nlpResponse := &models.NLPResponse{
		Understood:            response.Understood,
		Intent:                response.Intent,
		Action:                response.Action,
		TargetNPC:             response.TargetNPC,
		TargetItem:            response.TargetItem,
		ResponseText:          response.ResponseText,
		NeedsClarification:    response.NeedsClarification,
		ClarificationQuestion: response.ClarificationQuestion,
	}

	c.JSON(http.StatusOK, NLPInterpretResponse{
		Success: true,
		Data:    nlpResponse,
	})
}

// ExecuteNLCommandRequest represents a request to execute a natural language command
type ExecuteNLCommandRequest struct {
	PlayerInput string `json:"player_input" binding:"required"`
	SessionID   string `json:"session_id"`
	AutoExecute bool   `json:"auto_execute"`
}

// ExecuteNLCommandResponse represents the response from executing a natural language command
type ExecuteNLCommandResponse struct {
	Success        bool              `json:"success"`
	Interpretation *ai.NLPResult     `json:"interpretation"`
	ActionExecuted bool              `json:"action_executed"`
	GameResponse   string            `json:"game_response"`
	GameState      *models.GameState `json:"game_state,omitempty"`
	Suggestions    []string          `json:"suggestions"`
	Error          string            `json:"error,omitempty"`
}

// ExecuteNLCommand handles POST /api/v1/ai/execute
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

	if req.AutoExecute && nlpResult.Action != nil {
		actionResult := h.executeAction(nlpResult.Action)
		response.ActionExecuted = actionResult.Success
		response.GameResponse = actionResult.Message
		response.GameState = h.engine.GetState()
	}

	c.JSON(http.StatusOK, response)
}

func (h *AIHandler) executeAction(action *models.Action) *models.ActionResult {
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
	Type      string `json:"type" binding:"required"`
	NPCID     string `json:"npc_id" binding:"required"`
	Content   string `json:"content"`
	Amount    int    `json:"amount,omitempty"`
	Emotion   string `json:"emotion,omitempty"`
	SessionID string `json:"session_id"`
}

// ComplexInteractionResponse represents the response to a complex interaction
type ComplexInteractionResponse struct {
	Success         bool     `json:"success"`
	NPCResponse     string   `json:"npc_response"`
	MoodChange      string   `json:"mood_change,omitempty"`
	FriendshipDelta int      `json:"friendship_delta"`
	UpdatedPrice    int      `json:"updated_price,omitempty"`
	NextOptions     []string `json:"next_options,omitempty"`
	TriggeredEvent  string   `json:"triggered_event,omitempty"`
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

	if response.MoodChange != "" {
		h.moodSystem.SetMood(req.NPCID, response.MoodChange, 0)
	}

	c.JSON(http.StatusOK, response)
}

func (h *AIHandler) handleBargain(req ComplexInteractionRequest, profile models.NPCDialogueProfile, gameState *models.GameState) ComplexInteractionResponse {
	response := ComplexInteractionResponse{Success: true}

	friendship := 0
	for _, npc := range gameState.NPCs {
		if npc.ID == req.NPCID {
			friendship = npc.Friendship
			break
		}
	}

	discount := friendship / 500
	if discount > 20 {
		discount = 20
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

func (h *AIHandler) handleEmotion(req ComplexInteractionRequest, profile models.NPCDialogueProfile, gameState *models.GameState) ComplexInteractionResponse {
	response := ComplexInteractionResponse{Success: true}

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

func (h *AIHandler) handleNegotiate(req ComplexInteractionRequest, profile models.NPCDialogueProfile, gameState *models.GameState) ComplexInteractionResponse {
	response := ComplexInteractionResponse{Success: true}

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
