# Backend Service 架构设计文档

## 目录
- [概述](#概述)
- [AI Service 接口设计](#ai-service-接口设计)
- [角色生成 API 设计](#角色生成-api-设计)
- [数据库选型建议](#数据库选型建议)
- [服务架构图](#服务架构图)
- [实现路线图](#实现路线图)

---

## 概述

本文档描述 Stardew Community Agent 项目的 Backend Service 架构设计，重点包括：
- AI Service 接口（支持 SD API 和云端 API 两种模式）
- 角色生成 API（RESTful 设计）
- 数据库选型（角色数据存储）
- 服务架构图（backend <-> AI service 交互流程）

### 设计原则

1. **可扩展性**: AI Service 支持多种后端切换（本地 SD API / 云端 API）
2. **解耦设计**: Backend 与 AI Service 通过接口抽象，便于独立扩展
3. **高可用**: 支持降级策略，AI 服务不可用时使用 Mock 数据
4. **性能优先**: 异步处理、缓存策略、批量操作

---

## AI Service 接口设计

### 1. 接口抽象层

```go
// File: internal/services/ai/provider.go

package ai

import (
    "context"
)

// AIProvider 定义 AI 服务提供者接口
type AIProvider interface {
    // 生成对话
    GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error)

    // 生成角色
    GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error)

    // 自然语言理解
    Understand(ctx context.Context, req NLURequest) (*NLUResponse, error)

    // 行为决策
    DecideBehavior(ctx context.Context, req BehaviorRequest) (*BehaviorResponse, error)

    // 健康检查
    HealthCheck(ctx context.Context) error

    // 获取提供者信息
    ProviderInfo() ProviderInfo
}

// ProviderConfig AI 提供者配置
type ProviderConfig struct {
    Type        ProviderType `json:"type"`         // sd_api, cloud_api, mock
    Endpoint    string       `json:"endpoint"`     // API 端点
    APIKey      string       `json:"api_key"`      // API 密钥
    Model       string       `json:"model"`        // 模型名称
    MaxTokens   int          `json:"max_tokens"`   // 最大 token 数
    Temperature float64      `json:"temperature"`  // 温度参数
    Timeout     int          `json:"timeout"`      // 超时时间(秒)
    MaxRetries  int          `json:"max_retries"`  // 最大重试次数
}

type ProviderType string

const (
    ProviderTypeSDAPI   ProviderType = "sd_api"    // Stable Diffusion API (本地)
    ProviderTypeCloud   ProviderType = "cloud_api" // 云端 API (Claude, OpenAI 等)
    ProviderTypeMock    ProviderType = "mock"      // Mock 实现
    ProviderTypeOllama  ProviderType = "ollama"    // Ollama 本地模型
)

type ProviderInfo struct {
    Type      ProviderType `json:"type"`
    Model     string       `json:"model"`
    Available bool         `json:"available"`
    Latency   int64        `json:"latency_ms"`
}
```

### 2. SD API (Stable Diffusion API) 适配器

```go
// File: internal/services/ai/provider_sdapi.go

package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// SDAPIProvider Stable Diffusion API 提供者
type SDAPIProvider struct {
    config     ProviderConfig
    httpClient *http.Client
}

func NewSDAPIProvider(config ProviderConfig) *SDAPIProvider {
    return &SDAPIProvider{
        config: config,
        httpClient: &http.Client{
            Timeout: time.Duration(config.Timeout) * time.Second,
        },
    }
}

func (p *SDAPIProvider) GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error) {
    // SD API 使用 LLM 扩展进行对话生成
    payload := map[string]interface{}{
        "prompt":      p.buildDialoguePrompt(req),
        "max_length":  p.config.MaxTokens,
        "temperature": p.config.Temperature,
    }

    resp, err := p.post(ctx, "/api/v1/generate", payload)
    if err != nil {
        return nil, err
    }

    return p.parseDialogueResponse(resp)
}

func (p *SDAPIProvider) GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error) {
    // 使用 SD API 生成角色描述
    payload := map[string]interface{}{
        "prompt":      p.buildCharacterPrompt(req),
        "max_length":  2048,
        "temperature": p.config.Temperature,
    }

    resp, err := p.post(ctx, "/api/v1/generate", payload)
    if err != nil {
        return nil, err
    }

    return p.parseCharacterResponse(resp)
}

func (p *SDAPIProvider) post(ctx context.Context, path string, payload interface{}) (*APIResponse, error) {
    body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint+path, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := p.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("SD API error: %s", string(respBody))
    }

    var apiResp APIResponse
    json.Unmarshal(respBody, &apiResp)
    return &apiResp, nil
}

func (p *SDAPIProvider) HealthCheck(ctx context.Context) error {
    resp, err := p.httpClient.Get(p.config.Endpoint + "/health")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("SD API unhealthy")
    }
    return nil
}

func (p *SDAPIProvider) ProviderInfo() ProviderInfo {
    return ProviderInfo{
        Type:  ProviderTypeSDAPI,
        Model: p.config.Model,
    }
}

// APIResponse SD API 响应结构
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   string      `json:"error,omitempty"`
}
```

### 3. 云端 API 适配器 (支持 Claude, OpenAI 等)

```go
// File: internal/services/ai/provider_cloud.go

package ai

import (
    "context"
    "encoding/json"
    "fmt"
)

// CloudAPIProvider 云端 API 提供者 (统一接口)
type CloudAPIProvider struct {
    config   ProviderConfig
    client   CloudClient
    fallback AIProvider // 降级提供者
}

// CloudClient 云端 API 客户端接口
type CloudClient interface {
    Chat(ctx context.Context, messages []Message) (*ChatResponse, error)
    Embed(ctx context.Context, text string) ([]float32, error)
}

type Message struct {
    Role    string `json:"role"`    // system, user, assistant
    Content string `json:"content"`
}

type ChatResponse struct {
    Content      string `json:"content"`
    TokensUsed   int    `json:"tokens_used"`
    FinishReason string `json:"finish_reason"`
}

func NewCloudAPIProvider(config ProviderConfig) *CloudAPIProvider {
    var client CloudClient

    // 根据配置创建不同的客户端
    switch config.Type {
    case "claude":
        client = NewClaudeClient(config)
    case "openai":
        client = NewOpenAIClient(config)
    case "deepseek":
        client = NewDeepSeekClient(config)
    default:
        client = NewClaudeClient(config) // 默认使用 Claude
    }

    return &CloudAPIProvider{
        config: config,
        client: client,
    }
}

func (p *CloudAPIProvider) GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error) {
    messages := []Message{
        {Role: "system", Content: p.buildSystemPrompt(req)},
        {Role: "user", Content: p.buildUserPrompt(req)},
    }

    resp, err := p.client.Chat(ctx, messages)
    if err != nil {
        // 尝试降级
        if p.fallback != nil {
            return p.fallback.GenerateDialogue(ctx, req)
        }
        return nil, fmt.Errorf("cloud API failed: %w", err)
    }

    return p.parseDialogueResponse(resp.Content)
}

func (p *CloudAPIProvider) GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error) {
    messages := []Message{
        {Role: "system", Content: p.buildCharacterSystemPrompt()},
        {Role: "user", Content: p.buildCharacterUserPrompt(req)},
    }

    resp, err := p.client.Chat(ctx, messages)
    if err != nil {
        if p.fallback != nil {
            return p.fallback.GenerateCharacter(ctx, req)
        }
        return nil, err
    }

    return p.parseCharacterResponse(resp.Content)
}

func (p *CloudAPIProvider) SetFallback(provider AIProvider) {
    p.fallback = provider
}

func (p *CloudAPIProvider) HealthCheck(ctx context.Context) error {
    // 简单的健康检查
    _, err := p.client.Chat(ctx, []Message{
        {Role: "user", Content: "ping"},
    })
    return err
}

func (p *CloudAPIProvider) ProviderInfo() ProviderInfo {
    return ProviderInfo{
        Type:  ProviderTypeCloud,
        Model: p.config.Model,
    }
}

// Claude 客户端实现
type ClaudeClient struct {
    apiKey   string
    model    string
    baseURL  string
    maxTokens int
}

func NewClaudeClient(config ProviderConfig) *ClaudeClient {
    baseURL := config.Endpoint
    if baseURL == "" {
        baseURL = "https://api.anthropic.com"
    }
    return &ClaudeClient{
        apiKey:    config.APIKey,
        model:     config.Model,
        baseURL:   baseURL,
        maxTokens: config.MaxTokens,
    }
}

func (c *ClaudeClient) Chat(ctx context.Context, messages []Message) (*ChatResponse, error) {
    // 实现 Claude API 调用
    // POST /v1/messages
    return &ChatResponse{}, nil
}

func (c *ClaudeClient) Embed(ctx context.Context, text string) ([]float32, error) {
    // Claude 暂不支持 embedding，返回空
    return nil, fmt.Errorf("claude does not support embedding")
}
```

### 4. AI Service 管理器

```go
// File: internal/services/ai/manager.go

package ai

import (
    "context"
    "sync"
)

// AIServiceManager AI 服务管理器
type AIServiceManager struct {
    mu          sync.RWMutex
    providers   map[ProviderType]AIProvider
    defaultType ProviderType
    config      *ServiceConfig
}

type ServiceConfig struct {
    DefaultProvider ProviderType     `json:"default_provider"`
    Providers       []ProviderConfig `json:"providers"`
    EnableCache     bool             `json:"enable_cache"`
    CacheTTL        int              `json:"cache_ttl_seconds"`
}

func NewAIServiceManager(config *ServiceConfig) (*AIServiceManager, error) {
    manager := &AIServiceManager{
        providers:   make(map[ProviderType]AIProvider),
        defaultType: config.DefaultProvider,
        config:      config,
    }

    // 初始化所有配置的提供者
    for _, pc := range config.Providers {
        provider, err := CreateProvider(pc)
        if err != nil {
            return nil, err
        }
        manager.providers[pc.Type] = provider
    }

    // 设置降级链
    manager.setupFallbackChain()

    return manager, nil
}

func CreateProvider(config ProviderConfig) (AIProvider, error) {
    switch config.Type {
    case ProviderTypeSDAPI:
        return NewSDAPIProvider(config), nil
    case ProviderTypeCloud:
        return NewCloudAPIProvider(config), nil
    case ProviderTypeOllama:
        return NewOllamaProvider(config), nil
    case ProviderTypeMock:
        return NewMockProvider(config), nil
    default:
        return nil, fmt.Errorf("unknown provider type: %s", config.Type)
    }
}

func (m *AIServiceManager) setupFallbackChain() {
    // 设置降级链: cloud -> sd_api -> mock
    if cloud, ok := m.providers[ProviderTypeCloud]; ok {
        if sdapi, ok := m.providers[ProviderTypeSDAPI]; ok {
            cloud.(*CloudAPIProvider).SetFallback(sdapi)
        }
    }
}

func (m *AIServiceManager) GetProvider(t ProviderType) (AIProvider, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    if t == "" {
        t = m.defaultType
    }

    provider, ok := m.providers[t]
    if !ok {
        return nil, fmt.Errorf("provider not found: %s", t)
    }

    return provider, nil
}

func (m *AIServiceManager) GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error) {
    provider, err := m.GetProvider("")
    if err != nil {
        return nil, err
    }
    return provider.GenerateDialogue(ctx, req)
}

func (m *AIServiceManager) GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error) {
    provider, err := m.GetProvider("")
    if err != nil {
        return nil, err
    }
    return provider.GenerateCharacter(ctx, req)
}

// SwitchProvider 切换默认提供者
func (m *AIServiceManager) SwitchProvider(t ProviderType) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, ok := m.providers[t]; !ok {
        return fmt.Errorf("provider not found: %s", t)
    }

    m.defaultType = t
    return nil
}
```

---

## 角色生成 API 设计

### 1. RESTful API 端点

```
基础路径: /api/v1/characters

端点列表:
├── POST   /characters                    # 创建角色 (AI 生成)
├── GET    /characters                    # 获取角色列表
├── GET    /characters/:id                # 获取单个角色
├── PUT    /characters/:id                # 更新角色
├── DELETE /characters/:id                # 删除角色
├── POST   /characters/generate           # AI 生成角色 (异步)
├── GET    /characters/generate/:job_id   # 查询生成任务状态
├── POST   /characters/:id/dialogue       # 生成角色对话
├── POST   /characters/:id/behavior       # 生成角色行为
└── GET    /characters/templates          # 获取角色模板列表
```

### 2. 数据模型

```go
// File: internal/models/character.go

package models

import "time"

// Character 角色数据模型
type Character struct {
    ID              string            `json:"id" bson:"_id"`
    Name            string            `json:"name" bson:"name"`
    Role            string            `json:"role" bson:"role"`           // mayor, shopkeeper, villager, etc.
    Age             string            `json:"age" bson:"age"`             // young, middle, senior
    Gender          string            `json:"gender" bson:"gender"`

    // 性格系统
    Personality     Personality       `json:"personality" bson:"personality"`
    Traits          []string          `json:"traits" bson:"traits"`
    Values          []string          `json:"values" bson:"values"`
    Dislikes        []string          `json:"dislikes" bson:"dislikes"`
    Fears           []string          `json:"fears,omitempty" bson:"fears,omitempty"`
    Hopes           []string          `json:"hopes,omitempty" bson:"hopes,omitempty"`

    // 背景故事
    Background      string            `json:"background" bson:"background"`
    Secrets         []string          `json:"secrets,omitempty" bson:"secrets,omitempty"`
    CharacterArc    string            `json:"character_arc,omitempty" bson:"character_arc,omitempty"`

    // 互动系统
    GiftPreferences GiftPreferences   `json:"gift_preferences" bson:"gift_preferences"`
    DialogueThemes  DialogueThemes    `json:"dialogue_themes" bson:"dialogue_themes"`
    SpeechStyle     string            `json:"speech_style" bson:"speech_style"`
    Hobby           string            `json:"hobby" bson:"hobby"`

    // 行为系统
    Schedule        []ScheduleEntry   `json:"schedule,omitempty" bson:"schedule,omitempty"`
    Behaviors       []BehaviorConfig  `json:"behaviors,omitempty" bson:"behaviors,omitempty"`

    // 友谊系统
    MaxFriendship   int               `json:"max_friendship" bson:"max_friendship"`
    FriendshipUnlocks map[string]string `json:"friendship_unlocks" bson:"friendship_unlocks"`

    // 元数据
    Version         int               `json:"version" bson:"version"`
    CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
    UpdatedAt       time.Time         `json:"updated_at" bson:"updated_at"`
    CreatedBy       string            `json:"created_by" bson:"created_by"` // "ai", "manual", "template"
    IsActive        bool              `json:"is_active" bson:"is_active"`
}

// Personality 性格配置
type Personality struct {
    Openness     int `json:"openness"`      // 开放性 (0-100)
    Conscientiousness int `json:"conscientiousness"` // 尽责性
    Extraversion int `json:"extraversion"`  // 外向性
    Agreeableness int `json:"agreeableness"` // 宜人性
    Neuroticism  int `json:"neuroticism"`   // 神经质
}

// BehaviorConfig 行为配置
type BehaviorConfig struct {
    Trigger   string   `json:"trigger"`     // 触发条件
    Action    string   `json:"action"`      // 行为类型
    Priority  int      `json:"priority"`    // 优先级
    Cooldown  int      `json:"cooldown"`    // 冷却时间(游戏分钟)
    Conditions []string `json:"conditions"` // 附加条件
}

// CharacterGenerationRequest 角色生成请求
type CharacterGenerationRequest struct {
    // 基础信息
    NameHint     string `json:"name_hint,omitempty"`     // 名字提示
    Role         string `json:"role,omitempty"`          // 角色类型
    Age          string `json:"age,omitempty"`           // 年龄段
    Gender       string `json:"gender,omitempty"`        // 性别

    // 性格倾向
    PersonalityHint string `json:"personality_hint,omitempty"` // 性格提示
    TraitsHint      []string `json:"traits_hint,omitempty"`    // 特质提示

    // 生成选项
    GenerateSchedule    bool `json:"generate_schedule"`     // 是否生成日程
    GenerateDialogue    bool `json:"generate_dialogue"`     // 是否生成对话主题
    GenerateBehaviors   bool `json:"generate_behaviors"`    // 是否生成行为配置
    DetailedBackground  bool `json:"detailed_background"`   // 是否生成详细背景

    // AI 提供者选项
    Provider string `json:"provider,omitempty"`    // 指定 AI 提供者
    Creativity float64 `json:"creativity,omitempty"` // 创造性 (temperature)
}

// CharacterGenerationResponse 角色生成响应
type CharacterGenerationResponse struct {
    JobID       string    `json:"job_id"`
    Status      string    `json:"status"`      // pending, processing, completed, failed
    Character   *Character `json:"character,omitempty"`
    Error       string    `json:"error,omitempty"`
    GeneratedAt time.Time `json:"generated_at,omitempty"`
}

// CharacterFilter 角色查询过滤器
type CharacterFilter struct {
    Role      string `form:"role"`
    Age       string `form:"age"`
    Active    *bool  `form:"active"`
    Name      string `form:"name"`
    Page      int    `form:"page"`
    PageSize  int    `form:"page_size"`
}

// CharacterTemplate 角色模板
type CharacterTemplate struct {
    ID          string       `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    BaseConfig  Character    `json:"base_config"`
    Category    string       `json:"category"` // villager, merchant, special
}
```

### 3. API Handler 实现

```go
// File: internal/api/character_handler.go

package api

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "stardew-agent/internal/models"
    "stardew-agent/internal/services/ai"
    "stardew-agent/internal/services/storage"
)

// CharacterHandler 角色相关 API 处理器
type CharacterHandler struct {
    aiManager   *ai.AIServiceManager
    charStore   storage.CharacterStore
    jobManager  *GenerationJobManager
}

func NewCharacterHandler(aiManager *ai.AIServiceManager, store storage.CharacterStore) *CharacterHandler {
    return &CharacterHandler{
        aiManager:  aiManager,
        charStore:  store,
        jobManager: NewGenerationJobManager(),
    }
}

// CreateCharacter 手动创建角色
// POST /api/v1/characters
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
    var char models.Character
    if err := c.ShouldBindJSON(&char); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 设置默认值
    char.ID = uuid.New().String()
    char.CreatedAt = time.Now()
    char.UpdatedAt = time.Now()
    char.CreatedBy = "manual"
    char.Version = 1

    if char.MaxFriendship == 0 {
        char.MaxFriendship = 2500
    }

    // 保存到数据库
    if err := h.charStore.Create(c.Request.Context(), &char); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success":  true,
        "character": char,
    })
}

// GenerateCharacter AI 生成角色
// POST /api/v1/characters/generate
func (h *CharacterHandler) GenerateCharacter(c *gin.Context) {
    var req models.CharacterGenerationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 创建生成任务
    job := h.jobManager.CreateJob(req)

    // 异步执行生成
    go h.executeGeneration(job.ID, req)

    c.JSON(http.StatusAccepted, gin.H{
        "success": true,
        "job_id":  job.ID,
        "status":  "pending",
        "message": "Character generation started",
    })
}

// executeGeneration 执行角色生成
func (h *CharacterHandler) executeGeneration(jobID string, req models.CharacterGenerationRequest) {
    ctx := context.Background()

    h.jobManager.UpdateStatus(jobID, "processing")

    // 调用 AI 服务生成角色
    char, err := h.aiManager.GenerateCharacter(ctx, ai.CharacterRequest{
        NameHint:        req.NameHint,
        Role:           req.Role,
        Age:            req.Age,
        Gender:         req.Gender,
        PersonalityHint: req.PersonalityHint,
        TraitsHint:      req.TraitsHint,
        GenerateSchedule: req.GenerateSchedule,
        GenerateDialogue: req.GenerateDialogue,
        Creativity:      req.Creativity,
    })

    if err != nil {
        h.jobManager.UpdateError(jobID, err.Error())
        return
    }

    // 设置元数据
    char.ID = uuid.New().String()
    char.CreatedAt = time.Now()
    char.UpdatedAt = time.Now()
    char.CreatedBy = "ai"
    char.Version = 1
    char.IsActive = true

    // 保存到数据库
    if err := h.charStore.Create(ctx, char); err != nil {
        h.jobManager.UpdateError(jobID, err.Error())
        return
    }

    h.jobManager.Complete(jobID, char)
}

// GetGenerationJob 查询生成任务状态
// GET /api/v1/characters/generate/:job_id
func (h *CharacterHandler) GetGenerationJob(c *gin.Context) {
    jobID := c.Param("job_id")

    job, err := h.jobManager.Get(jobID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
        return
    }

    c.JSON(http.StatusOK, job)
}

// ListCharacters 获取角色列表
// GET /api/v1/characters
func (h *CharacterHandler) ListCharacters(c *gin.Context) {
    var filter models.CharacterFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.PageSize <= 0 {
        filter.PageSize = 20
    }

    characters, total, err := h.charStore.List(c.Request.Context(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":    true,
        "characters": characters,
        "pagination": gin.H{
            "page":       filter.Page,
            "page_size":  filter.PageSize,
            "total":      total,
            "total_pages": (total + filter.PageSize - 1) / filter.PageSize,
        },
    })
}

// GetCharacter 获取单个角色
// GET /api/v1/characters/:id
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
    id := c.Param("id")

    char, err := h.charStore.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":   true,
        "character": char,
    })
}

// UpdateCharacter 更新角色
// PUT /api/v1/characters/:id
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
    id := c.Param("id")

    var updates models.Character
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updates.ID = id
    updates.UpdatedAt = time.Now()
    updates.Version++ // 版本递增

    if err := h.charStore.Update(c.Request.Context(), &updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":   true,
        "character": updates,
    })
}

// DeleteCharacter 删除角色
// DELETE /api/v1/characters/:id
func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
    id := c.Param("id")

    if err := h.charStore.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Character deleted",
    })
}

// GenerateDialogue 为角色生成对话
// POST /api/v1/characters/:id/dialogue
func (h *CharacterHandler) GenerateDialogue(c *gin.Context) {
    id := c.Param("id")

    var req struct {
        Context     string `json:"context"`
        PlayerInput string `json:"player_input"`
        Mood        string `json:"mood"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 获取角色信息
    char, err := h.charStore.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
        return
    }

    // 调用 AI 生成对话
    dialogue, err := h.aiManager.GenerateDialogue(c.Request.Context(), ai.DialogueRequest{
        Character:    char,
        Context:      req.Context,
        PlayerInput:  req.PlayerInput,
        CurrentMood:  req.Mood,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":  true,
        "dialogue": dialogue,
    })
}

// ListTemplates 获取角色模板列表
// GET /api/v1/characters/templates
func (h *CharacterHandler) ListTemplates(c *gin.Context) {
    templates := h.getDefaultTemplates()

    c.JSON(http.StatusOK, gin.H{
        "success":   true,
        "templates": templates,
    })
}

func (h *CharacterHandler) getDefaultTemplates() []models.CharacterTemplate {
    return []models.CharacterTemplate{
        {
            ID:          "template_villager",
            Name:        "村民模板",
            Description: "普通村民角色模板",
            Category:    "villager",
            BaseConfig: models.Character{
                Role:          "villager",
                MaxFriendship: 2500,
                GiftPreferences: models.GiftPreferences{
                    Like:    []string{"flower", "cooked_dish"},
                    Neutral: []string{"mineral", "forage"},
                },
            },
        },
        {
            ID:          "template_merchant",
            Name:        "商人模板",
            Description: "商店老板角色模板",
            Category:    "merchant",
            BaseConfig: models.Character{
                Role:          "shopkeeper",
                MaxFriendship: 2500,
                GiftPreferences: models.GiftPreferences{
                    Love:    []string{"rare_item", "artisan_good"},
                    Like:    []string{"crop", "mineral"},
                },
            },
        },
    }
}
```

### 4. 路由配置

```go
// File: internal/api/character_routes.go

package api

import (
    "github.com/gin-gonic/gin"
)

func SetupCharacterRoutes(router *gin.RouterGroup, handler *CharacterHandler) {
    characters := router.Group("/characters")
    {
        // CRUD 操作
        characters.POST("", handler.CreateCharacter)
        characters.GET("", handler.ListCharacters)
        characters.GET("/:id", handler.GetCharacter)
        characters.PUT("/:id", handler.UpdateCharacter)
        characters.DELETE("/:id", handler.DeleteCharacter)

        // AI 生成
        characters.POST("/generate", handler.GenerateCharacter)
        characters.GET("/generate/:job_id", handler.GetGenerationJob)

        // 互动功能
        characters.POST("/:id/dialogue", handler.GenerateDialogue)

        // 模板
        characters.GET("/templates", handler.ListTemplates)
    }
}
```

---

## 数据库选型建议

### 1. 选型对比分析

| 数据库 | 优势 | 劣势 | 适用场景 |
|--------|------|------|----------|
| **MongoDB** | 文档结构灵活、嵌套数据支持好、查询丰富 | 内存占用高、无事务ACID | 角色数据、动态属性 |
| **PostgreSQL** | ACID事务、JSON支持、成熟稳定 | 需要预定义schema | 关系数据、交易记录 |
| **Redis** | 极快速度、缓存友好 | 内存限制、持久化弱 | 缓存、会话、实时状态 |
| **SQLite** | 轻量、无服务器、易部署 | 并发限制、无网络访问 | 开发测试、小型部署 |

### 2. 推荐方案：MongoDB + Redis 组合

```
┌─────────────────────────────────────────────────────────────┐
│                      数据存储架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌───────────────┐       ┌───────────────┐                │
│   │   MongoDB     │       │    Redis      │                │
│   │   主数据库     │       │   缓存层       │                │
│   └───────┬───────┘       └───────┬───────┘                │
│           │                       │                         │
│   ┌───────▼───────────────────────▼───────┐                │
│   │              数据访问层                 │                │
│   │         (Repository Pattern)          │                │
│   └───────────────────────────────────────┘                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3. 数据模型设计

```go
// File: internal/services/storage/mongo_store.go

package storage

import (
    "context"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "stardew-agent/internal/models"
)

// CharacterStore 角色存储接口
type CharacterStore interface {
    Create(ctx context.Context, char *models.Character) error
    Get(ctx context.Context, id string) (*models.Character, error)
    Update(ctx context.Context, char *models.Character) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter models.CharacterFilter) ([]models.Character, int, error)
    FindByRole(ctx context.Context, role string) ([]models.Character, error)
}

// MongoCharacterStore MongoDB 实现
type MongoCharacterStore struct {
    collection *mongo.Collection
}

func NewMongoCharacterStore(db *mongo.Database) *MongoCharacterStore {
    collection := db.Collection("characters")

    // 创建索引
    indexes := []mongo.IndexModel{
        {Keys: bson.D{{Key: "name", Value: 1}}},
        {Keys: bson.D{{Key: "role", Value: 1}}},
        {Keys: bson.D{{Key: "is_active", Value: 1}}},
        {Keys: bson.D{{Key: "created_at", Value: -1}}},
    }
    collection.Indexes().CreateMany(context.Background(), indexes)

    return &MongoCharacterStore{collection: collection}
}

func (s *MongoCharacterStore) Create(ctx context.Context, char *models.Character) error {
    _, err := s.collection.InsertOne(ctx, char)
    return err
}

func (s *MongoCharacterStore) Get(ctx context.Context, id string) (*models.Character, error) {
    var char models.Character
    err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&char)
    if err != nil {
        return nil, err
    }
    return &char, nil
}

func (s *MongoCharacterStore) Update(ctx context.Context, char *models.Character) error {
    filter := bson.M{"_id": char.ID}
    update := bson.M{"$set": char}

    _, err := s.collection.UpdateOne(ctx, filter, update)
    return err
}

func (s *MongoCharacterStore) Delete(ctx context.Context, id string) error {
    _, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (s *MongoCharacterStore) List(ctx context.Context, filter models.CharacterFilter) ([]models.Character, int, error) {
    // 构建查询条件
    query := bson.M{}
    if filter.Role != "" {
        query["role"] = filter.Role
    }
    if filter.Age != "" {
        query["age"] = filter.Age
    }
    if filter.Active != nil {
        query["is_active"] = *filter.Active
    }
    if filter.Name != "" {
        query["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
    }

    // 计算总数
    total, err := s.collection.CountDocuments(ctx, query)
    if err != nil {
        return nil, 0, err
    }

    // 分页查询
    opts := options.Find().
        SetSkip(int64((filter.Page - 1) * filter.PageSize)).
        SetLimit(int64(filter.PageSize)).
        SetSort(bson.D{{Key: "created_at", Value: -1}})

    cursor, err := s.collection.Find(ctx, query, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    var characters []models.Character
    if err := cursor.All(ctx, &characters); err != nil {
        return nil, 0, err
    }

    return characters, int(total), nil
}
```

### 4. Redis 缓存层

```go
// File: internal/services/storage/redis_cache.go

package storage

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"

    "stardew-agent/internal/models"
)

// CharacterCache Redis 缓存实现
type CharacterCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewCharacterCache(client *redis.Client) *CharacterCache {
    return &CharacterCache{
        client: client,
        ttl:    time.Hour * 2, // 2小时缓存
    }
}

func (c *CharacterCache) Get(ctx context.Context, id string) (*models.Character, error) {
    data, err := c.client.Get(ctx, c.key(id)).Bytes()
    if err != nil {
        return nil, err
    }

    var char models.Character
    if err := json.Unmarshal(data, &char); err != nil {
        return nil, err
    }

    return &char, nil
}

func (c *CharacterCache) Set(ctx context.Context, char *models.Character) error {
    data, err := json.Marshal(char)
    if err != nil {
        return err
    }

    return c.client.Set(ctx, c.key(char.ID), data, c.ttl).Err()
}

func (c *CharacterCache) Delete(ctx context.Context, id string) error {
    return c.client.Del(ctx, c.key(id)).Err()
}

func (c *CharacterCache) key(id string) string {
    return "character:" + id
}

// CachedCharacterStore 带缓存的角色存储
type CachedCharacterStore struct {
    store CharacterStore
    cache *CharacterCache
}

func NewCachedCharacterStore(store CharacterStore, cache *CharacterCache) *CachedCharacterStore {
    return &CachedCharacterStore{
        store: store,
        cache: cache,
    }
}

func (s *CachedCharacterStore) Get(ctx context.Context, id string) (*models.Character, error) {
    // 先查缓存
    if char, err := s.cache.Get(ctx, id); err == nil {
        return char, nil
    }

    // 查数据库
    char, err := s.store.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    go s.cache.Set(context.Background(), char)

    return char, nil
}
```

### 5. 配置文件

```yaml
# config/database.yaml

database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "stardew_agent"
    options:
      maxPoolSize: 100
      minPoolSize: 10
      maxIdleTimeMS: 60000

  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    poolSize: 50

ai_service:
  default_provider: "cloud_api"
  providers:
    - type: "cloud_api"
      endpoint: "https://api.anthropic.com"
      model: "claude-sonnet-4-5-20250929"
      max_tokens: 4096
      timeout: 30

    - type: "sd_api"
      endpoint: "http://localhost:8080"
      model: "llama-3"
      max_tokens: 2048
      timeout: 60

    - type: "mock"
      enabled: true
```

---

## 服务架构图

### 1. 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层                                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │   Web Frontend  │  │   AI Agent      │  │   Admin Panel   │             │
│  │   (Vue.js)      │  │   (Python/JS)   │  │   (React)       │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
└───────────┼─────────────────────┼─────────────────────┼─────────────────────┘
            │                     │                     │
            └─────────────────────┼─────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API 网关层                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Nginx Reverse Proxy                          │   │
│  │   - SSL 终止                                                         │   │
│  │   - 负载均衡                                                         │   │
│  │   - 请求路由                                                         │   │
│  │   - WebSocket 代理                                                   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Backend Service                                  │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                      Gin Web Server (Go)                              │  │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐        │  │
│  │  │   Game     │ │ Character  │ │    AI      │ │   Admin    │        │  │
│  │  │   Handler  │ │  Handler   │ │  Handler   │ │  Handler   │        │  │
│  │  └─────┬──────┘ └─────┬──────┘ └─────┬──────┘ └────────────┘        │  │
│  │        │              │              │                               │  │
│  │        └──────────────┼──────────────┘                               │  │
│  │                       │                                              │  │
│  │  ┌────────────────────▼────────────────────────────────────────┐    │  │
│  │  │                    Service Layer                             │    │  │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │    │  │
│  │  │  │  Game    │ │Character │ │   AI     │ │   NPC    │       │    │  │
│  │  │  │  Engine  │ │ Service  │ │ Service  │ │ Service  │       │    │  │
│  │  │  └──────────┘ └──────────┘ └────┬─────┘ └──────────┘       │    │  │
│  │  └─────────────────────────────────┼──────────────────────────┘    │  │
│  └──────────────────────────────────────┼──────────────────────────────┘  │
└─────────────────────────────────────────┼─────────────────────────────────┘
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    │                     │                     │
                    ▼                     ▼                     ▼
┌───────────────────────────┐ ┌───────────────────┐ ┌───────────────────────┐
│      AI Service Layer     │ │   Data Layer      │ │   External Services   │
│  ┌─────────────────────┐  │ │ ┌───────────────┐ │ │ ┌───────────────────┐ │
│  │  AI Service Manager │  │ │ │   MongoDB     │ │ │ │  Claude API       │ │
│  │  ┌───────────────┐  │  │ │ │  (Characters) │ │ │ │  (Cloud LLM)      │ │
│  │  │ Cloud API     │  │  │ │ └───────────────┘ │ │ └───────────────────┘ │
│  │  │ Provider      │  │  │ │ ┌───────────────┐ │ │ ┌───────────────────┐ │
│  │  ├───────────────┤  │  │ │ │    Redis      │ │ │ │  SD API           │ │
│  │  │ SD API        │  │  │ │ │   (Cache)     │ │ │ │  (Local LLM)      │ │
│  │  │ Provider      │  │  │ │ └───────────────┘ │ │ └───────────────────┘ │
│  │  ├───────────────┤  │  │ │ ┌───────────────┐ │ │ ┌───────────────────┐ │
│  │  │ Ollama        │  │  │ │ │   SQLite      │ │ │ │  Ollama           │ │
│  │  │ Provider      │  │  │ │ │  (Dev Mode)   │ │ │ │  (Local Models)   │ │
│  │  ├───────────────┤  │  │ │ └───────────────┘ │ │ └───────────────────┘ │
│  │  │ Mock          │  │  │ └───────────────────┘ └───────────────────────┘
│  │  │ Provider      │  │  │
│  │  └───────────────┘  │  │
│  └─────────────────────┘  │
└───────────────────────────┘
```

### 2. Backend <-> AI Service 交互流程

```
┌──────────────────────────────────────────────────────────────────────────┐
│                    角色生成流程 (Character Generation)                     │
└──────────────────────────────────────────────────────────────────────────┘

Client                Backend               AI Service Manager        External AI
  │                     │                         │                       │
  │  POST /characters   │                         │                       │
  │  /generate          │                         │                       │
  │────────────────────>│                         │                       │
  │                     │                         │                       │
  │                     │  1. 创建生成任务         │                       │
  │                     │  (存入 JobManager)       │                       │
  │                     │                         │                       │
  │  202 Accepted       │                         │                       │
  │  {job_id: "xxx"}    │                         │                       │
  │<────────────────────│                         │                       │
  │                     │                         │                       │
  │                     │  [异步执行]              │                       │
  │                     │                         │                       │
  │                     │  GenerateCharacter()    │                       │
  │                     │────────────────────────>│                       │
  │                     │                         │                       │
  │                     │                         │ GetProvider()         │
  │                     │                         │ (选择: cloud/sd/mock)  │
  │                     │                         │                       │
  │                     │                         │  POST /v1/messages    │
  │                     │                         │──────────────────────>│
  │                     │                         │                       │
  │                     │                         │  AI Response          │
  │                     │                         │<──────────────────────│
  │                     │                         │                       │
  │                     │  CharacterResponse      │                       │
  │                     │<────────────────────────│                       │
  │                     │                         │                       │
  │                     │  2. 保存到 MongoDB       │                       │
  │                     │  3. 更新任务状态         │                       │
  │                     │                         │                       │
  │  GET /characters    │                         │                       │
  │  /generate/:job_id  │                         │                       │
  │────────────────────>│                         │                       │
  │                     │                         │                       │
  │  {status: completed,│                         │                       │
  │   character: {...}} │                         │                       │
  │<────────────────────│                         │                       │
  │                     │                         │                       │
```

### 3. 对话生成流程

```
┌──────────────────────────────────────────────────────────────────────────┐
│                    对话生成流程 (Dialogue Generation)                      │
└──────────────────────────────────────────────────────────────────────────┘

Client                Backend               AI Service                 AI Provider
  │                     │                    Manager                       │
  │                     │                      │                          │
  │  POST /characters   │                      │                          │
  │  /:id/dialogue      │                      │                          │
  │────────────────────>│                      │                          │
  │                     │                      │                          │
  │                     │ 1. 获取角色          │                          │
  │                     │    (MongoDB/Cache)   │                          │
  │                     │                      │                          │
  │                     │ GenerateDialogue()   │                          │
  │                     │─────────────────────>│                          │
  │                     │                      │                          │
  │                     │                      │ BuildContext()           │
  │                     │                      │ - 角色 Profile           │
  │                     │                      │ - 游戏 State             │
  │                     │                      │ - 历史 Interaction       │
  │                     │                      │ - 当前 Mood              │
  │                     │                      │                          │
  │                     │                      │ CheckCache()             │
  │                     │                      │ (相同上下文复用)          │
  │                     │                      │                          │
  │                     │                      │ Chat Completion          │
  │                     │                      │─────────────────────────>│
  │                     │                      │                          │
  │                     │                      │ Response                 │
  │                     │                      │<─────────────────────────│
  │                     │                      │                          │
  │                     │ DialogueResponse     │                          │
  │                     │<─────────────────────│                          │
  │                     │                      │                          │
  │                     │ 2. 更新缓存          │                          │
  │                     │ 3. 更新角色心情      │                          │
  │                     │                      │                          │
  │  {dialogue: "...",  │                      │                          │
  │   mood_change: "...",│                     │                          │
  │   friendship_delta}  │                     │                          │
  │<────────────────────│                      │                          │
  │                     │                      │                          │
```

### 4. 服务部署架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Docker Compose 架构                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              Docker Network                                  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          nginx (反向代理)                            │   │
│  │                         Port: 80, 443                               │   │
│  └───────────────────────────────┬─────────────────────────────────────┘   │
│                                  │                                          │
│          ┌───────────────────────┼───────────────────────┐                 │
│          │                       │                       │                 │
│          ▼                       ▼                       ▼                 │
│  ┌───────────────┐       ┌───────────────┐       ┌───────────────┐        │
│  │   frontend    │       │   backend     │       │   ollama      │        │
│  │   (Vue.js)    │       │   (Go)        │       │   (Optional)  │        │
│  │   Port: 3000  │       │   Port: 8080  │       │   Port: 11434 │        │
│  └───────────────┘       └───────┬───────┘       └───────────────┘        │
│                                  │                                          │
│                  ┌───────────────┼───────────────┐                         │
│                  │               │               │                         │
│                  ▼               ▼               ▼                         │
│          ┌───────────────┐ ┌───────────────┐ ┌───────────────┐            │
│          │   mongodb     │ │    redis      │ │  sd-api       │            │
│          │   Port: 27017 │ │   Port: 6379  │ │  (Optional)   │            │
│          │               │ │               │ │  Port: 8080   │            │
│          └───────────────┘ └───────────────┘ └───────────────┘            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 实现路线图

### Phase 1: 基础架构 (Week 1-2)

1. **AI Service 接口层**
   - [ ] 定义 `AIProvider` 接口
   - [ ] 实现 `MockProvider` (用于开发测试)
   - [ ] 实现 `AIServiceManager`

2. **数据存储层**
   - [ ] MongoDB 连接和基础配置
   - [ ] 实现 `CharacterStore` 接口
   - [ ] Redis 缓存层实现

### Phase 2: 角色生成 API (Week 3-4)

1. **API 端点**
   - [ ] CRUD 操作 (Create, Read, Update, Delete)
   - [ ] 异步生成任务管理
   - [ ] 角色模板系统

2. **AI 集成**
   - [ ] `CloudAPIProvider` 实现 (Claude/OpenAI)
   - [ ] 角色生成 Prompt 工程
   - [ ] 降级策略

### Phase 3: 高级功能 (Week 5-6)

1. **对话系统增强**
   - [ ] 上下文管理
   - [ ] 对话缓存
   - [ ] 情绪系统

2. **本地 AI 支持**
   - [ ] `SDAPIProvider` 实现
   - [ ] `OllamaProvider` 实现
   - [ ] 模型切换机制

### Phase 4: 优化与部署 (Week 7-8)

1. **性能优化**
   - [ ] 批量操作
   - [ ] 连接池优化
   - [ ] 缓存策略优化

2. **部署配置**
   - [ ] Docker Compose 完善
   - [ ] 环境变量配置
   - [ ] 健康检查和监控

---

*文档版本: 1.0*
*创建日期: 2026-02-12*
*作者: R&D Team*
