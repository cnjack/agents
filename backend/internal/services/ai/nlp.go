package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"stardew-agent/internal/models"
)

// NLPService handles natural language processing for game interactions
type NLPService struct {
	llmClient      LLMClient
	intentPatterns []IntentPattern
	contextHistory map[string]*ConversationContext // session ID -> context
}

// IntentPattern defines a pattern for intent recognition
type IntentPattern struct {
	Intent     string
	Patterns   []*regexp.Regexp
	Keywords   []string
	Extractors []func(string) map[string]string
}

// ConversationContext maintains conversation state
type ConversationContext struct {
	SessionID       string
	LastIntent      string
	LastNPC         string
	LastItem        string
	ConversationLog []ConversationEntry
	Timestamp       time.Time
}

// ConversationEntry records a single exchange
type ConversationEntry struct {
	PlayerInput string
	Intent      string
	Response    string
	Timestamp   time.Time
}

// ComplexIntent represents a complex interaction intent
type ComplexIntent struct {
	Type          string            `json:"type"`           // bargain, express_emotion, negotiate, chat, trade
	SubIntent     string            `json:"sub_intent"`     // offer, counter, accept, reject
	TargetNPC     string            `json:"target_npc"`
	TargetItem    string            `json:"target_item"`
	Amount        int               `json:"amount"`
	Emotion       string            `json:"emotion"`        // happy, sad, angry, excited, grateful
	Parameters    map[string]string `json:"parameters"`
	Confidence    float64           `json:"confidence"`
	RawExpression string            `json:"raw_expression"`
}

// NLPResult represents the full NLP processing result
type NLPResult struct {
	Understood      bool              `json:"understood"`
	Intent          string            `json:"intent"`
	ComplexIntent   *ComplexIntent    `json:"complex_intent,omitempty"`
	Action          *models.Action    `json:"action,omitempty"`
	TargetNPC       string            `json:"target_npc,omitempty"`
	TargetItem      string            `json:"target_item,omitempty"`
	ResponseText    string            `json:"response_text"`
	Suggestions     []string          `json:"suggestions,omitempty"`
	NeedsClarification bool           `json:"needs_clarification"`
	ClarificationQuestion string      `json:"clarification_question,omitempty"`
	Context         *ConversationContext `json:"context,omitempty"`
}

// NewNLPService creates a new NLP service
func NewNLPService(llmClient LLMClient) *NLPService {
	service := &NLPService{
		llmClient:      llmClient,
		contextHistory: make(map[string]*ConversationContext),
	}

	service.initializePatterns()
	return service
}

// initializePatterns sets up intent recognition patterns
func (n *NLPService) initializePatterns() {
	n.intentPatterns = []IntentPattern{
		// Greeting intents
		{
			Intent:   "greet",
			Keywords: []string{"你好", "嗨", "hello", "hi", "hey", "早上好", "晚上好"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(你好|嗨|hello|hi|hey)`),
			},
		},
		// Movement intents
		{
			Intent:   "move",
			Keywords: []string{"走", "去", "移动", "move", "go", "上", "下", "左", "右"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(向|往)?(上|下|左|右|north|south|east|west)(走|去|移动)?`),
				regexp.MustCompile(`(?i)go\s+(up|down|left|right|north|south|east|west)`),
				regexp.MustCompile(`去(.+)`),
			},
		},
		// Gift intents
		{
			Intent:   "gift",
			Keywords: []string{"送", "给", "礼物", "gift", "give"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(送|给)(.+?)(礼物|东西)?给?(.+)`),
				regexp.MustCompile(`(?i)give\s+(.+?)\s+to\s+(.+)`),
			},
		},
		// Trade/Buy intents
		{
			Intent:   "buy",
			Keywords: []string{"买", "购买", "buy", "purchase", "要"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(买|购买|要)(\d+)?个?(.+?)`),
				regexp.MustCompile(`(?i)buy\s+(\d+)?\s*(.+)`),
			},
		},
		// Trade/Sell intents
		{
			Intent:   "sell",
			Keywords: []string{"卖", "出售", "sell"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(卖|出售)(\d+)?个?(.+?)`),
				regexp.MustCompile(`(?i)sell\s+(\d+)?\s*(.+)`),
			},
		},
		// Bargain intents
		{
			Intent:   "bargain",
			Keywords: []string{"便宜", "折扣", "bargain", "discount", "太贵", "讲价"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(能不能|可以)(便宜|折扣)?(.+)?`),
				regexp.MustCompile(`(?i)too\s+expensive`),
				regexp.MustCompile(`(?i)(便宜点|打个折)`),
			},
		},
		// Express emotion
		{
			Intent:   "express_emotion",
			Keywords: []string{"开心", "难过", "生气", "感谢", "喜欢", "讨厌", "happy", "sad", "angry", "thank", "love", "hate"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(我|很|好)(开心|高兴|快乐|难过|伤心|生气|感谢|喜欢|讨厌)`),
				regexp.MustCompile(`(?i)i\s+(am|feel)\s+(happy|sad|angry|grateful)`),
				regexp.MustCompile(`(?i)(谢谢|感谢)(你|你)`),
			},
		},
		// Quest intents
		{
			Intent:   "quest",
			Keywords: []string{"任务", "quest", "委托", "帮忙"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(接受|接)(任务|委托)?`),
				regexp.MustCompile(`(?i)accept\s+(the\s+)?quest`),
			},
		},
		// Ask about NPC
		{
			Intent:   "ask_about",
			Keywords: []string{"怎么样", "如何", "什么", "how", "what", "why"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(.+)怎么样`),
				regexp.MustCompile(`(?i)how\s+(is|are|about)\s+(.+)`),
			},
		},
		// Farewell
		{
			Intent:   "farewell",
			Keywords: []string{"再见", "拜拜", "bye", "goodbye"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)(再见|拜拜|bye|goodbye)`),
			},
		},
	}
}

// Process processes natural language input and returns structured result
func (n *NLPService) Process(ctx context.Context, input string, gameState *models.GameState, sessionID string) (*NLPResult, error) {
	// Get or create conversation context
	convCtx := n.getOrCreateContext(sessionID)

	// Normalize input
	normalizedInput := strings.TrimSpace(strings.ToLower(input))

	// Step 1: Pattern-based intent recognition
	intent := n.recognizeIntent(normalizedInput)

	// Step 2: Extract entities
	entities := n.extractEntities(normalizedInput, gameState, convCtx)

	// Step 3: Build complex intent
	complexIntent := n.buildComplexIntent(intent, entities, normalizedInput, gameState)

	// Step 4: Generate action
	action := n.intentToAction(intent, entities, gameState)

	// Step 5: Generate response text
	responseText := n.generateResponse(intent, entities, gameState, convCtx)

	// Step 6: Update context
	n.updateContext(convCtx, input, intent, responseText)

	// Step 7: Generate suggestions
	suggestions := n.generateSuggestions(intent, entities, gameState)

	return &NLPResult{
		Understood:            intent != "unknown",
		Intent:                intent,
		ComplexIntent:         complexIntent,
		Action:                action,
		TargetNPC:             entities["npc"],
		TargetItem:            entities["item"],
		ResponseText:          responseText,
		Suggestions:           suggestions,
		NeedsClarification:    intent == "unknown" || intent == "clarify",
		ClarificationQuestion: n.getClarificationQuestion(intent, entities),
		Context:               convCtx,
	}, nil
}

// recognizeIntent identifies the intent from input
func (n *NLPService) recognizeIntent(input string) string {
	bestMatch := "unknown"
	bestScore := 0

	for _, pattern := range n.intentPatterns {
		score := 0

		// Check patterns
		for _, re := range pattern.Patterns {
			if re.MatchString(input) {
				score += 2
			}
		}

		// Check keywords
		for _, keyword := range pattern.Keywords {
			if strings.Contains(input, strings.ToLower(keyword)) {
				score++
			}
		}

		if score > bestScore {
			bestScore = score
			bestMatch = pattern.Intent
		}
	}

	return bestMatch
}

// extractEntities extracts entities from input
func (n *NLPService) extractEntities(input string, gameState *models.GameState, ctx *ConversationContext) map[string]string {
	entities := make(map[string]string)

	// Extract NPC names
	for _, npc := range gameState.NPCs {
		if strings.Contains(strings.ToLower(input), strings.ToLower(npc.Name)) ||
			strings.Contains(strings.ToLower(input), strings.ToLower(npc.ID)) {
			entities["npc"] = npc.ID
			break
		}
	}

	// Use context if NPC not found
	if entities["npc"] == "" && ctx.LastNPC != "" {
		if strings.Contains(input, "他") || strings.Contains(input, "她") ||
			strings.Contains(input, "him") || strings.Contains(input, "her") {
			entities["npc"] = ctx.LastNPC
		}
	}

	// Extract item names from inventory
	for _, item := range gameState.Player.Gifts {
		if strings.Contains(strings.ToLower(input), strings.ToLower(item.Name)) ||
			strings.Contains(strings.ToLower(input), strings.ToLower(item.ID)) {
			entities["item"] = item.ID
			entities["item_name"] = item.Name
			break
		}
	}

	// Extract direction
	directions := map[string]string{
		"上": "up", "北": "up", "north": "up",
		"下": "down", "南": "down", "south": "down",
		"左": "left", "西": "left", "west": "left",
		"右": "right", "东": "right", "east": "right",
	}
	for cn, en := range directions {
		if strings.Contains(input, cn) || strings.Contains(input, en) {
			entities["direction"] = en
			break
		}
	}

	// Extract numbers
	re := regexp.MustCompile(`\d+`)
	if matches := re.FindStringSubmatch(input); len(matches) > 0 {
		entities["amount"] = matches[0]
	}

	// Extract emotion
	emotions := []string{"开心", "难过", "生气", "感谢", "喜欢", "讨厌", "happy", "sad", "angry", "grateful", "love", "hate"}
	for _, emotion := range emotions {
		if strings.Contains(input, emotion) {
			entities["emotion"] = emotion
			break
		}
	}

	// Extract price/bargain amount
	priceRe := regexp.MustCompile(`(\d+)(块|元|金|gold|g)`)
	if matches := priceRe.FindStringSubmatch(input); len(matches) > 1 {
		entities["price"] = matches[1]
	}

	return entities
}

// buildComplexIntent creates a complex intent structure
func (n *NLPService) buildComplexIntent(intent string, entities map[string]string, input string, gameState *models.GameState) *ComplexIntent {
	complex := &ComplexIntent{
		Type:          intent,
		TargetNPC:     entities["npc"],
		TargetItem:    entities["item"],
		Parameters:    entities,
		RawExpression: input,
		Confidence:    0.8,
	}

	// Determine sub-intent
	switch intent {
	case "bargain":
		if strings.Contains(input, "太贵") || strings.Contains(input, "expensive") {
			complex.SubIntent = "complain_price"
		} else if strings.Contains(input, "便宜") || strings.Contains(input, "cheap") || strings.Contains(input, "折扣") {
			complex.SubIntent = "request_discount"
		}
	case "express_emotion":
		complex.Emotion = entities["emotion"]
		if strings.Contains(input, "谢谢") || strings.Contains(input, "thank") {
			complex.SubIntent = "gratitude"
		} else if strings.Contains(input, "喜欢") || strings.Contains(input, "love") {
			complex.SubIntent = "affection"
		} else if strings.Contains(input, "讨厌") || strings.Contains(input, "hate") {
			complex.SubIntent = "dislike"
		}
	case "buy", "sell":
		if amountStr, ok := entities["amount"]; ok {
			if amount, err := strconv.Atoi(amountStr); err == nil {
				complex.Amount = amount
			} else {
				complex.Amount = 1
			}
		} else {
			complex.Amount = 1
		}
	}

	return complex
}

// intentToAction converts intent to game action
func (n *NLPService) intentToAction(intent string, entities map[string]string, gameState *models.GameState) *models.Action {
	switch intent {
	case "greet", "talk":
		npcID := entities["npc"]
		if npcID == "" && len(gameState.NPCs) > 0 {
			// Find nearest NPC
			npcID = n.findNearestNPC(gameState)
		}
		if npcID != "" {
			return &models.Action{
				Type: models.ActionTalk,
				Params: models.ActionParams{
					NPC: npcID,
				},
			}
		}

	case "move":
		if dir, ok := entities["direction"]; ok {
			return &models.Action{
				Type: models.ActionMove,
				Params: models.ActionParams{
					Direction: models.Direction(dir),
				},
			}
		}

	case "gift":
		npcID := entities["npc"]
		itemID := entities["item"]
		if npcID != "" && itemID != "" {
			return &models.Action{
				Type: models.ActionGiveGift,
				Params: models.ActionParams{
					NPC:       npcID,
					GiftItem:  itemID,
				},
			}
		}

	case "quest":
		// Find available quest
		for _, quest := range gameState.Quests {
			if quest.Status == models.QuestStatusAvailable {
				return &models.Action{
					Type: models.ActionAcceptQuest,
					Params: models.ActionParams{
						QuestID: quest.ID,
					},
				}
			}
		}

	case "bargain":
		// Bargaining is handled by dialogue system
		return nil
	}

	return nil
}

// generateResponse creates a response text
func (n *NLPService) generateResponse(intent string, entities map[string]string, gameState *models.GameState, ctx *ConversationContext) string {
	npcName := entities["npc"]
	if npcName == "" {
		npcName = "NPC"
	} else {
		// Find NPC name from ID
		for _, npc := range gameState.NPCs {
			if npc.ID == npcName {
				npcName = npc.Name
				break
			}
		}
	}

	itemName := entities["item_name"]
	if itemName == "" {
		itemName = entities["item"]
	}

	switch intent {
	case "greet", "talk":
		return fmt.Sprintf("你想和 %s 对话", npcName)

	case "move":
		if dir, ok := entities["direction"]; ok {
			directions := map[string]string{
				"up":    "上", "down": "下",
				"left":  "左", "right": "右",
			}
			return fmt.Sprintf("向%s移动", directions[dir])
		}
		return "你想去哪里？"

	case "gift":
		if itemName != "" && npcName != "" {
			return fmt.Sprintf("你想把 %s 送给 %s", itemName, npcName)
		}
		return "你想送什么礼物给谁？"

	case "buy":
		amount := 1
		if aStr, ok := entities["amount"]; ok {
			if a, err := strconv.Atoi(aStr); err == nil {
				amount = a
			}
		}
		if itemName != "" {
			return fmt.Sprintf("你想购买 %d 个 %s", amount, itemName)
		}
		return "你想买什么？"

	case "sell":
		amount := 1
		if aStr, ok := entities["amount"]; ok {
			if a, err := strconv.Atoi(aStr); err == nil {
				amount = a
			}
		}
		if itemName != "" {
			return fmt.Sprintf("你想出售 %d 个 %s", amount, itemName)
		}
		return "你想卖什么？"

	case "bargain":
		return fmt.Sprintf("你想和 %s 讨价还价", npcName)

	case "express_emotion":
		emotion := entities["emotion"]
		return fmt.Sprintf("你向 %s 表达了 %s 的情感", npcName, emotion)

	case "quest":
		return "你接受了任务"

	case "farewell":
		return "再见！期待下次见面"

	case "ask_about":
		return fmt.Sprintf("你想了解关于 %s 的信息", npcName)

	default:
		return "我不太明白你的意思，你能换个方式说吗？"
	}
}

// generateSuggestions creates action suggestions
func (n *NLPService) generateSuggestions(intent string, entities map[string]string, gameState *models.GameState) []string {
	suggestions := []string{}

	nearbyNPCs := n.getNearbyNPCs(gameState)

	switch intent {
	case "greet":
		if len(nearbyNPCs) > 0 {
			suggestions = append(suggestions, fmt.Sprintf("和 %s 聊聊", nearbyNPCs[0].Name))
		}
		suggestions = append(suggestions, "问问最近怎么样")
		suggestions = append(suggestions, "送个礼物")

	case "gift":
		for _, item := range gameState.Player.Gifts {
			suggestions = append(suggestions, fmt.Sprintf("送 %s", item.Name))
		}

	case "unknown":
		if len(nearbyNPCs) > 0 {
			suggestions = append(suggestions, fmt.Sprintf("和 %s 对话", nearbyNPCs[0].Name))
		}
		suggestions = append(suggestions, "四处走走")
		suggestions = append(suggestions, "查看任务")

	case "bargain":
		suggestions = append(suggestions, "接受价格")
		suggestions = append(suggestions, "继续砍价")
		suggestions = append(suggestions, "放弃购买")
	}

	return suggestions
}

// getClarificationQuestion returns a clarification question if needed
func (n *NLPService) getClarificationQuestion(intent string, entities map[string]string) string {
	if intent == "unknown" {
		return "你想做什么？你可以尝试：\n- 和NPC对话\n- 送礼物\n- 移动探索\n- 接受任务\n- 购买物品"
	}

	if intent == "gift" && (entities["npc"] == "" || entities["item"] == "") {
		if entities["npc"] == "" {
			return "你想把礼物送给谁？"
		}
		return "你想送什么礼物？"
	}

	if intent == "move" && entities["direction"] == "" {
		return "你想往哪个方向走？（上/下/左/右）"
	}

	return ""
}

// Helper functions

func (n *NLPService) getOrCreateContext(sessionID string) *ConversationContext {
	if ctx, exists := n.contextHistory[sessionID]; exists {
		return ctx
	}

	ctx := &ConversationContext{
		SessionID:       sessionID,
		ConversationLog: []ConversationEntry{},
		Timestamp:       time.Now(),
	}
	n.contextHistory[sessionID] = ctx
	return ctx
}

func (n *NLPService) updateContext(ctx *ConversationContext, input, intent, response string) {
	ctx.LastIntent = intent
	ctx.ConversationLog = append(ctx.ConversationLog, ConversationEntry{
		PlayerInput: input,
		Intent:      intent,
		Response:    response,
		Timestamp:   time.Now(),
	})

	// Keep only last 20 entries
	if len(ctx.ConversationLog) > 20 {
		ctx.ConversationLog = ctx.ConversationLog[len(ctx.ConversationLog)-20:]
	}
}

func (n *NLPService) findNearestNPC(gameState *models.GameState) string {
	playerPos := gameState.Player.Position
	minDist := int(^uint(0) >> 1) // Max int
	nearestID := ""

	for _, npc := range gameState.NPCs {
		dx := npc.Position.X - playerPos.X
		dy := npc.Position.Y - playerPos.Y
		dist := dx*dx + dy*dy
		if dist < minDist {
			minDist = dist
			nearestID = npc.ID
		}
	}

	return nearestID
}

// ProcessWithLLM processes input using LLM for more complex understanding
func (n *NLPService) ProcessWithLLM(ctx context.Context, input string, gameState *models.GameState, sessionID string) (*NLPResult, error) {
	// Build prompt for LLM
	systemPrompt := `你是一个游戏中的自然语言理解系统。分析玩家的输入，识别意图并提取实体。

可能的意图类型:
- greet: 打招呼/对话
- move: 移动
- gift: 送礼物
- buy: 购买物品
- sell: 出售物品
- bargain: 讨价还价
- express_emotion: 表达情感
- quest: 任务相关
- ask_about: 询问信息
- farewell: 告别

请用JSON格式回复:
{
  "intent": "意图类型",
  "target_npc": "目标NPC的ID",
  "target_item": "目标物品的ID",
  "parameters": {"其他参数": "值"},
  "confidence": 0.9,
  "response_text": "理解后的描述"
}`

	userPrompt := fmt.Sprintf("玩家输入: %s\n\n当前游戏状态: 玩家位置(%d,%d), 附近NPC: %v, 背包物品: %v",
		input,
		gameState.Player.Position.X,
		gameState.Player.Position.Y,
		getNPCIDs(gameState.NPCs),
		getItemIDs(gameState.Player.Gifts),
	)

	response, err := n.llmClient.GenerateWithSystem(ctx, systemPrompt, userPrompt)
	if err != nil {
		// Fallback to rule-based processing
		return n.Process(ctx, input, gameState, sessionID)
	}

	// Parse LLM response
	var llmResult struct {
		Intent       string            `json:"intent"`
		TargetNPC    string            `json:"target_npc"`
		TargetItem   string            `json:"target_item"`
		Parameters   map[string]string `json:"parameters"`
		Confidence   float64           `json:"confidence"`
		ResponseText string            `json:"response_text"`
	}

	if err := json.Unmarshal([]byte(response), &llmResult); err != nil {
		return n.Process(ctx, input, gameState, sessionID)
	}

	// Convert to NLPResult
	result := &NLPResult{
		Understood:   llmResult.Intent != "",
		Intent:       llmResult.Intent,
		TargetNPC:    llmResult.TargetNPC,
		TargetItem:   llmResult.TargetItem,
		ResponseText: llmResult.ResponseText,
		Context:      n.getOrCreateContext(sessionID),
	}

	// Build action from LLM result
	result.Action = n.intentToAction(llmResult.Intent, map[string]string{
		"npc":  llmResult.TargetNPC,
		"item": llmResult.TargetItem,
	}, gameState)

	return result, nil
}

// ClearContext clears the conversation context for a session
func (n *NLPService) ClearContext(sessionID string) {
	delete(n.contextHistory, sessionID)
}

func getNPCIDs(npcs []models.NPCState) []string {
	ids := make([]string, len(npcs))
	for i, npc := range npcs {
		ids[i] = npc.ID
	}
	return ids
}

func getItemIDs(items []models.InventoryItem) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}

// getNearbyNPCs returns NPCs near the player
func (n *NLPService) getNearbyNPCs(gameState *models.GameState) []models.NPCState {
	playerPos := gameState.Player.Position
	var nearby []models.NPCState

	for _, npc := range gameState.NPCs {
		dx := npc.Position.X - playerPos.X
		dy := npc.Position.Y - playerPos.Y
		if dx*dx+dy*dy <= 25 { // Within 5 tiles
			nearby = append(nearby, npc)
		}
	}

	return nearby
}

func getNearbyNPCs(gameState *models.GameState) []models.NPCState {
	playerPos := gameState.Player.Position
	var nearby []models.NPCState

	for _, npc := range gameState.NPCs {
		dx := npc.Position.X - playerPos.X
		dy := npc.Position.Y - playerPos.Y
		if dx*dx+dy*dy <= 25 { // Within 5 tiles
			nearby = append(nearby, npc)
		}
	}

	return nearby
}
