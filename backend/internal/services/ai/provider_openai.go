package ai

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"stardew-agent/internal/config"
)

// OpenAICompatibleProvider OpenAI 兼容协议提供者
// 支持: 智谱 GLM-4, DeepSeek, Ollama, 以及所有 OpenAI 兼容 API
type OpenAICompatibleProvider struct {
	config     config.ProviderConfig
	httpClient *http.Client
}

// NewOpenAICompatibleProvider 创建 OpenAI 兼容协议提供者
func NewOpenAICompatibleProvider(cfg config.ProviderConfig) *OpenAICompatibleProvider {
	// 创建自定义 HTTP 客户端，配置更长的超时时间
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 120
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	return &OpenAICompatibleProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: transport,
		},
	}
}

// OpenAI API 请求/响应结构

// ChatCompletionRequest Chat Completion 请求
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// Message 对话消息
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// ChatCompletionResponse Chat Completion 响应
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice 选择结果
type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

// Usage Token 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// APIError API 错误
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (type: %s, code: %s)", e.Message, e.Type, e.Code)
}

// GenerateDialogue 生成对话
func (p *OpenAICompatibleProvider) GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error) {
	// 构建消息
	messages := []Message{
		{Role: "system", Content: p.buildDialogueSystemPrompt(req)},
	}

	if req.PlayerInput != "" {
		messages = append(messages, Message{Role: "user", Content: req.PlayerInput})
	} else {
		messages = append(messages, Message{Role: "user", Content: "你好"})
	}

	// 调用 API
	resp, err := p.chatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	// 解析响应
	return p.parseDialogueResponse(resp)
}

// GenerateCharacter 生成角色
func (p *OpenAICompatibleProvider) GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error) {
	messages := []Message{
		{Role: "system", Content: p.buildCharacterSystemPrompt()},
		{Role: "user", Content: p.buildCharacterUserPrompt(req)},
	}

	resp, err := p.chatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	return p.parseCharacterResponse(resp)
}

// Understand 自然语言理解
func (p *OpenAICompatibleProvider) Understand(ctx context.Context, req NLURequest) (*NLUResponse, error) {
	messages := []Message{
		{Role: "system", Content: p.buildNLUSystemPrompt(req)},
		{Role: "user", Content: req.Input},
	}

	resp, err := p.chatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	return p.parseNLUResponse(resp)
}

// DecideBehavior 行为决策
func (p *OpenAICompatibleProvider) DecideBehavior(ctx context.Context, req BehaviorRequest) (*BehaviorResponse, error) {
	messages := []Message{
		{Role: "system", Content: p.buildBehaviorSystemPrompt(req)},
		{Role: "user", Content: p.buildBehaviorUserPrompt(req)},
	}

	resp, err := p.chatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	return p.parseBehaviorResponse(resp)
}

// HealthCheck 健康检查
func (p *OpenAICompatibleProvider) HealthCheck(ctx context.Context) error {
	// 发送一个简单的测试请求
	messages := []Message{
		{Role: "user", Content: "ping"},
	}

	_, err := p.chatCompletion(ctx, messages)
	return err
}

// ProviderInfo 获取提供者信息
func (p *OpenAICompatibleProvider) ProviderInfo() ProviderInfo {
	return ProviderInfo{
		Type:      ProviderTypeOpenAICompatible,
		Model:     p.config.Model,
		Endpoint:  p.config.Endpoint,
		Available: true,
	}
}

// chatCompletion 调用 Chat Completion API
func (p *OpenAICompatibleProvider) chatCompletion(ctx context.Context, messages []Message) (*ChatCompletionResponse, error) {
	// 构建请求
	chatReq := ChatCompletionRequest{
		Model:       p.config.Model,
		Messages:    messages,
		MaxTokens:   p.config.MaxTokens,
		Temperature: p.config.Temperature,
		Stream:      false,
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发送请求（带重试）
	var resp *http.Response
	var lastErr error

	for i := 0; i < p.config.MaxRetries; i++ {
		// 每次重试都创建新的请求（因为 body 是一次性的）
		httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint+"/chat/completions", bytes.NewReader(body))
		if err != nil {
			lastErr = err
			continue
		}

		// 设置请求头
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

		resp, lastErr = p.httpClient.Do(httpReq)
		if lastErr == nil && resp.StatusCode < 500 {
			break
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if i < p.config.MaxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed after %d retries: %w", p.config.MaxRetries, lastErr)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(respBody))
	}

	// 检查错误
	if chatResp.Error != nil {
		return nil, chatResp.Error
	}

	// 检查是否有选择
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response: %s", string(respBody))
	}

	return &chatResp, nil
}

// Prompt 构建方法

func (p *OpenAICompatibleProvider) buildDialogueSystemPrompt(req DialogueRequest) string {
	char := req.Character
	if char == nil {
		return "你是一个友好的NPC角色，请用自然的语气与玩家对话。"
	}

	return fmt.Sprintf(`你是一个游戏中的NPC角色，请严格按照以下设定进行角色扮演：

【角色信息】
- 名字: %s
- 职业: %s
- 年龄段: %s
- 性别: %s

【性格特质】
%s

【背景故事】
%s

【说话风格】
%s

【喜好】
- 喜欢: %v
- 讨厌: %v

【对话规则】
1. 严格遵守角色设定，保持性格一致性
2. 用符合角色身份的语气说话
3. 根据友谊度(%d心)调整亲密度
4. 回复要自然、有趣、符合游戏氛围
5. 控制回复长度在50-100字以内

【回复格式】请用JSON格式回复:
{
  "dialogue": "你的对话内容",
  "mood_change": "happy/neutral/sad/angry/excited",
  "friendship_delta": 0,
  "suggested_actions": ["建议的后续动作"]
}`,
		char.Name,
		char.Role,
		char.Age,
		char.Gender,
		char.Traits,
		char.Background,
		char.SpeechStyle,
		char.GiftPreferences.Love,
		char.GiftPreferences.Hate,
		req.FriendshipHearts,
	)
}

func (p *OpenAICompatibleProvider) buildCharacterSystemPrompt() string {
	return `你是一个游戏角色设计师。你的任务是根据用户的描述生成一个完整的游戏NPC角色。

请严格按照以下JSON格式输出角色信息：

{
  "name": "角色名称",
  "role": "角色职业类型(mayor/shopkeeper/villager/fisherman/farmer/etc)",
  "age": "年龄段(young/middle/senior)",
  "gender": "性别",
  "traits": ["性格特质1", "性格特质2"],
  "values": ["重视的事物1", "重视的事物2"],
  "dislikes": ["讨厌的事物1", "讨厌的事物2"],
  "background": "详细的背景故事，200-300字",
  "secrets": ["秘密1（高友谊度时透露）"],
  "speech_style": "说话风格描述",
  "hobby": "兴趣爱好",
  "gift_preferences": {
    "love": ["最喜欢的礼物"],
    "like": ["喜欢的礼物"],
    "neutral": ["普通礼物"],
    "dislike": ["不喜欢的礼物"],
    "hate": ["讨厌的礼物"]
  },
  "dialogue_themes": {
    "morning": ["早晨对话主题"],
    "afternoon": ["下午对话主题"],
    "evening": ["晚上对话主题"]
  },
  "character_arc": "角色成长弧线描述"
}

设计要求：
1. 角色要有鲜明的个性和独特的魅力
2. 背景故事要有趣，能引发玩家的探索欲
3. 喜好要与角色性格相符
4. 对话主题要体现角色特点
5. 秘密要能增加角色的深度`
}

func (p *OpenAICompatibleProvider) buildCharacterUserPrompt(req CharacterRequest) string {
	return fmt.Sprintf(`请根据以下描述生成一个NPC角色：

- 名字提示: %s
- 角色类型: %s
- 年龄段: %s
- 性别: %s
- 性格提示: %s
- 特质提示: %v
- 需要生成日程: %v
- 需要生成对话主题: %v

请生成一个完整的角色设定。`, req.NameHint, req.Role, req.Age, req.Gender, req.PersonalityHint, req.TraitsHint, req.GenerateSchedule, req.GenerateDialogue)
}

func (p *OpenAICompatibleProvider) buildNLUSystemPrompt(req NLURequest) string {
	return `你是一个游戏指令解析助手。你的任务是理解玩家的自然语言输入，并将其转换为游戏动作。

请根据玩家输入，分析其意图并返回JSON格式的结果：

{
  "understood": true/false,
  "intent": "move/talk/gift/buy/sell/interact/unknown",
  "action": {
    "type": "动作类型",
    "params": {}
  },
  "target_npc": "目标NPC的ID（如果有）",
  "target_item": "目标物品ID（如果有）",
  "response_text": "给玩家的回复",
  "needs_clarification": false,
  "clarification_question": "如果需要澄清，这里提问"
}

可选的动作类型：
- move: 移动 (params: {"direction": "up/down/left/right"})
- talk: 对话 (params: {"npc": "npc_id"})
- give_gift: 送礼 (params: {"npc": "npc_id", "gift_item": "item_id"})
- buy: 购买 (params: {"item": "item_id", "quantity": 数量})
- sell: 出售 (params: {"item": "item_id", "quantity": 数量})
- interact: 交互 (params: {"target": "target_id"})`
}

func (p *OpenAICompatibleProvider) buildBehaviorSystemPrompt(req BehaviorRequest) string {
	return fmt.Sprintf(`你是一个游戏AI行为决策系统。根据NPC的状态和环境，决定NPC下一步的行为。

NPC信息:
- ID: %s
- 当前心情: %s
- 当前目标: %s
- 体力: %d
- 社交需求: %d

游戏时间: %s, %d:%02d
天气: %s
玩家在附近: %v

可选动作: %v

请返回JSON格式的行为决策：
{
  "action": "动作名称",
  "target_location": "目标地点",
  "target_position": {"x": 0, "y": 0},
  "reason": "决策理由",
  "mood_after_action": "执行后的心情",
  "can_interrupt": true/false
}`,
		req.NPCID, req.CurrentMood, req.CurrentGoal, req.Energy, req.SocialNeed,
		req.Season, req.Hour, req.Minute, req.Weather, req.PlayerNearby, req.PossibleActions)
}

func (p *OpenAICompatibleProvider) buildBehaviorUserPrompt(req BehaviorRequest) string {
	return fmt.Sprintf("请根据以上信息决定NPC的行为。当前情境：%s", req.Context)
}

// 响应解析方法

func (p *OpenAICompatibleProvider) parseDialogueResponse(resp *ChatCompletionResponse) (*DialogueResponse, error) {
	content := resp.Choices[0].Message.Content

	// 尝试解析JSON
	result := &DialogueResponse{
		Dialogue: content,
	}

	// 尝试提取JSON
	jsonStart := indexOf(content, "{")
	jsonEnd := lastIndexOf(content, "}")

	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := content[jsonStart : jsonEnd+1]

		var parsed struct {
			Dialogue        string   `json:"dialogue"`
			MoodChange      string   `json:"mood_change"`
			FriendshipDelta int      `json:"friendship_delta"`
			SuggestedActions []string `json:"suggested_actions"`
		}

		if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
			result.Dialogue = parsed.Dialogue
			result.MoodChange = parsed.MoodChange
			result.FriendshipDelta = parsed.FriendshipDelta
			result.SuggestedActions = parsed.SuggestedActions
		}
	}

	return result, nil
}

func (p *OpenAICompatibleProvider) parseCharacterResponse(resp *ChatCompletionResponse) (*CharacterResponse, error) {
	content := resp.Choices[0].Message.Content

	// 提取JSON
	jsonStart := indexOf(content, "{")
	jsonEnd := lastIndexOf(content, "}")

	if jsonStart < 0 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := content[jsonStart : jsonEnd+1]

	var char CharacterResponse
	if err := json.Unmarshal([]byte(jsonStr), &char); err != nil {
		return nil, fmt.Errorf("failed to parse character: %w", err)
	}

	return &char, nil
}

func (p *OpenAICompatibleProvider) parseNLUResponse(resp *ChatCompletionResponse) (*NLUResponse, error) {
	content := resp.Choices[0].Message.Content

	jsonStart := indexOf(content, "{")
	jsonEnd := lastIndexOf(content, "}")

	if jsonStart < 0 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := content[jsonStart : jsonEnd+1]

	var result NLUResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse NLU response: %w", err)
	}

	return &result, nil
}

func (p *OpenAICompatibleProvider) parseBehaviorResponse(resp *ChatCompletionResponse) (*BehaviorResponse, error) {
	content := resp.Choices[0].Message.Content

	jsonStart := indexOf(content, "{")
	jsonEnd := lastIndexOf(content, "}")

	if jsonStart < 0 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := content[jsonStart : jsonEnd+1]

	var result BehaviorResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse behavior response: %w", err)
	}

	return &result, nil
}

// 辅助函数

func indexOf(s string, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func lastIndexOf(s string, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
