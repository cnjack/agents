# 测试策略文档

## 文档信息
- **版本**: 1.0
- **创建日期**: 2026-02-12
- **作者**: QA Engineer
- **项目**: Stardew Community Agent - AI Agent Development Platform

---

## 1. 概述

本文档定义了 Stardew Community Agent 项目的全面测试策略，包括：
- 角色生成质量验证方案
- API 测试用例设计
- 预生成脚本验证方案
- 测试数据准备建议

### 1.1 测试目标

1. **功能正确性**: 确保所有功能按设计文档正常工作
2. **数据完整性**: 验证生成数据符合 JSON Schema 规范
3. **API 可靠性**: 确保 RESTful API 端点稳定可靠
4. **性能指标**: 满足响应时间和并发要求
5. **兼容性**: 支持多种 AI Provider 切换

### 1.2 测试范围

| 模块 | 测试类型 | 优先级 |
|------|----------|--------|
| 角色生成器 | 单元测试、集成测试、E2E测试 | P0 |
| AI Service | 单元测试、Mock测试 | P0 |
| Backend API | API测试、性能测试 | P0 |
| 预生成脚本 | 单元测试、验证测试 | P1 |
| WebSocket | 连接测试、消息测试 | P1 |
| 前端组件 | 组件测试、E2E测试 | P2 |

---

## 2. 角色生成质量验证方案

### 2.1 JSON Schema 验证

#### 2.1.1 验证规则

基于 `CHARACTER_GENERATOR_DESIGN.md` 中定义的 JSON Schema，实现以下验证：

```go
// File: backend/internal/validators/character_validator.go

package validators

import (
    "encoding/json"
    "fmt"
    "regexp"

    "github.com/xeipuuv/gojsonschema"
)

// CharacterValidator 角色数据验证器
type CharacterValidator struct {
    schema *gojsonschema.Schema
}

// NewCharacterValidator 创建验证器
func NewCharacterValidator(schemaPath string) (*CharacterValidator, error) {
    schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
    schema, err := gojsonschema.NewSchema(schemaLoader)
    if err != nil {
        return nil, err
    }
    return &CharacterValidator{schema: schema}, nil
}

// Validate 验证角色数据
func (v *CharacterValidator) Validate(characterJSON []byte) (*ValidationResult, error) {
    docLoader := gojsonschema.NewBytesLoader(characterJSON)
    result, err := v.schema.Validate(docLoader)
    if err != nil {
        return nil, err
    }

    return &ValidationResult{
        Valid:   result.Valid(),
        Errors:  v.formatErrors(result.Errors()),
    }, nil
}

// ValidationResult 验证结果
type ValidationResult struct {
    Valid   bool     `json:"valid"`
    Errors  []string `json:"errors,omitempty"`
}

// 验证规则详解
/*
必需字段验证:
- id: 必需，匹配 ^[a-z0-9_]+$ 模式
- basic_info: 必需，包含 name, gender, birthday
- appearance: 必需，包含 hairstyle, eyes, skin, clothing, body_type
- personality: 必需，至少包含 primary_trait

数据范围验证:
- name: 1-20 字符
- birthday.day: 1-28
- body_type.height: 0.8-1.2
- body_type.build: 0.7-1.3
- secondary_traits: 最多3项
- jewelry: 最多3项

枚举值验证:
- gender: male, female, non_binary
- season: spring, summer, fall, winter
- hair_style: 16种预定义样式
- eye_shape: 8种预定义样式
- 等等...
*/
```

#### 2.1.2 测试用例 - Schema 验证

```go
// File: backend/tests/validators/character_validator_test.go

package validators_test

import (
    "encoding/json"
    "testing"

    "stardew-agent/internal/validators"
)

func TestCharacterValidator_ValidCharacter(t *testing.T) {
    validator, _ := validators.NewCharacterValidator("schemas/character.json")

    validCharacter := map[string]interface{}{
        "id": "test_char_001",
        "basic_info": map[string]interface{}{
            "name":   "Alex",
            "gender": "male",
            "birthday": map[string]interface{}{
                "season": "spring",
                "day":    14,
            },
        },
        "appearance": map[string]interface{}{
            "hairstyle": map[string]interface{}{
                "style": "short_spiky",
                "color": "#3d2314",
            },
            "eyes": map[string]interface{}{
                "shape": "round",
                "color": "#3498db",
            },
            "skin": map[string]interface{}{
                "tone": "#deb887",
            },
            "clothing": map[string]interface{}{
                "top": map[string]interface{}{
                    "style": "flannel",
                    "color": "#c0392b",
                },
                "bottom": map[string]interface{}{
                    "style": "jeans",
                    "color": "#4a6fa5",
                },
                "shoes": "work_boots",
            },
            "body_type": map[string]interface{}{
                "height": 1.0,
                "build":  1.1,
            },
        },
        "personality": map[string]interface{}{
            "primary_trait": "friendly",
        },
    }

    data, _ := json.Marshal(validCharacter)
    result, err := validator.Validate(data)

    if err != nil {
        t.Fatalf("Validation failed with error: %v", err)
    }
    if !result.Valid {
        t.Errorf("Expected valid character, got errors: %v", result.Errors)
    }
}

func TestCharacterValidator_MissingRequiredFields(t *testing.T) {
    testCases := []struct {
        name        string
        character   map[string]interface{}
        expectError string
    }{
        {
            name: "missing name",
            character: map[string]interface{}{
                "id": "test_001",
                "basic_info": map[string]interface{}{
                    "gender": "male",
                },
            },
            expectError: "name is required",
        },
        {
            name: "invalid gender",
            character: map[string]interface{}{
                "id": "test_002",
                "basic_info": map[string]interface{}{
                    "name":   "Test",
                    "gender": "invalid_gender",
                },
            },
            expectError: "gender",
        },
        {
            name: "birthday day out of range",
            character: map[string]interface{}{
                "id": "test_003",
                "basic_info": map[string]interface{}{
                    "name":   "Test",
                    "gender": "male",
                    "birthday": map[string]interface{}{
                        "season": "spring",
                        "day":    30, // invalid: max is 28
                    },
                },
            },
            expectError: "maximum",
        },
    }

    validator, _ := validators.NewCharacterValidator("schemas/character.json")

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            data, _ := json.Marshal(tc.character)
            result, _ := validator.Validate(data)

            if result.Valid {
                t.Errorf("Expected validation error for %s", tc.name)
            }
        })
    }
}
```

### 2.2 业务规则验证

#### 2.2.1 角色属性一致性验证

```go
// File: backend/internal/validators/business_rules.go

package validators

import (
    "errors"
    "regexp"

    "stardew-agent/internal/models"
)

// BusinessRuleValidator 业务规则验证器
type BusinessRuleValidator struct{}

// ValidateAppearanceConsistency 验证外观属性一致性
func (v *BusinessRuleValidator) ValidateAppearanceConsistency(char *models.Character) error {
    // 验证颜色格式
    colorPattern := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

    if !colorPattern.MatchString(char.Appearance.Hairstyle.Color) {
        return errors.New("invalid hair color format")
    }
    if !colorPattern.MatchString(char.Appearance.Eyes.Color) {
        return errors.New("invalid eye color format")
    }
    if !colorPattern.MatchString(char.Appearance.Skin.Tone) {
        return errors.New("invalid skin tone format")
    }

    return nil
}

// ValidatePersonalityConsistency 验证性格一致性
func (v *BusinessRuleValidator) ValidatePersonalityConsistency(char *models.Character) error {
    // 检查性格特质是否与对话风格匹配
    traitStyleMap := map[string][]string{
        "friendly":  {"casual", "warm"},
        "shy":       {"polite", "reserved"},
        "energetic": {"enthusiastic", "casual"},
        "calm":      {"polite", "thoughtful"},
        "serious":   {"formal", "professional"},
        "mysterious":{"cryptic", "reserved"},
    }

    // 检查次要特质数量限制
    if len(char.Personality.SecondaryTraits) > 3 {
        return errors.New("secondary traits cannot exceed 3")
    }

    return nil
}

// ValidateScheduleValidity 验证日程有效性
func (v *BusinessRuleValidator) ValidateScheduleValidity(char *models.Character) error {
    for _, entry := range char.Schedule {
        // 验证时间范围 (6:00 - 26:00 游戏时间)
        if entry.StartHour < 6 || entry.StartHour > 26 {
            return errors.New("schedule start hour must be between 6 and 26")
        }
        if entry.EndHour < entry.StartHour {
            return errors.New("schedule end hour must be after start hour")
        }
    }
    return nil
}
```

### 2.3 AI 生成质量评估

#### 2.3.1 内容质量检查

```go
// File: backend/internal/validators/ai_quality.go

package validators

import (
    "regexp"
    "strings"
)

// AIQualityValidator AI生成内容质量验证器
type AIQualityValidator struct {
    minBackgroundLength int
    maxBackgroundLength int
    forbiddenPatterns   []*regexp.Regexp
}

// NewAIQualityValidator 创建质量验证器
func NewAIQualityValidator() *AIQualityValidator {
    return &AIQualityValidator{
        minBackgroundLength: 50,
        maxBackgroundLength: 2000,
        forbiddenPatterns: []*regexp.Regexp{
            regexp.MustCompile(`(?i)as an ai`),
            regexp.MustCompile(`(?i)i cannot`),
            regexp.MustCompile(`(?i)error:`),
            regexp.MustCompile(`\{.*error.*\}`),
        },
    }
}

// QualityScore 质量评分结果
type QualityScore struct {
    Overall       float64            `json:"overall"`
    Dimensions    map[string]float64 `json:"dimensions"`
    Issues        []string           `json:"issues,omitempty"`
    PassThreshold bool               `json:"pass_threshold"`
}

// ValidateCharacterQuality 验证角色生成质量
func (v *AIQualityValidator) ValidateCharacterQuality(char *models.Character) *QualityScore {
    score := &QualityScore{
        Dimensions: make(map[string]float64),
        Issues:     []string{},
    }

    // 1. 完整性评分 (0-100)
    completeness := v.scoreCompleteness(char)
    score.Dimensions["completeness"] = completeness

    // 2. 一致性评分 (0-100)
    consistency := v.scoreConsistency(char)
    score.Dimensions["consistency"] = consistency

    // 3. 丰富度评分 (0-100)
    richness := v.scoreRichness(char)
    score.Dimensions["richness"] = richness

    // 4. 内容质量评分 (0-100)
    contentQuality := v.scoreContentQuality(char)
    score.Dimensions["content_quality"] = contentQuality

    // 计算总分
    score.Overall = (completeness*0.3 + consistency*0.25 + richness*0.25 + contentQuality*0.2)
    score.PassThreshold = score.Overall >= 70.0

    return score
}

// scoreCompleteness 评分：完整性
func (v *AIQualityValidator) scoreCompleteness(char *models.Character) float64 {
    score := 100.0
    deductions := 0.0

    // 检查必需字段
    requiredFields := []string{
        char.Name,
        char.Role,
        char.Background,
        char.Personality.PrimaryTrait,
    }

    for _, field := range requiredFields {
        if field == "" {
            deductions += 10
        }
    }

    // 检查可选但推荐字段
    if len(char.Traits) == 0 {
        deductions += 5
    }
    if char.Hobby == "" {
        deductions += 5
    }
    if len(char.Schedule) == 0 {
        deductions += 10
    }

    return max(0, score-deductions)
}

// scoreConsistency 评分：一致性
func (v *AIQualityValidator) scoreConsistency(char *models.Character) float64 {
    score := 100.0

    // 检查性格与背景故事是否一致
    // 例如：害羞的角色不应该有过于外向的对话主题

    return score
}

// scoreRichness 评分：丰富度
func (v *AIQualityValidator) scoreRichness(char *models.Character) float64 {
    score := 0.0

    // 背景故事长度
    bgLen := len(char.Background)
    if bgLen >= v.minBackgroundLength {
        score += 25
        if bgLen >= 200 {
            score += 15
        }
    }

    // 特质数量
    if len(char.Traits) >= 2 {
        score += 20
    }

    // 日程条目数量
    if len(char.Schedule) >= 3 {
        score += 20
    }

    // 对话主题数量
    if len(char.DialogueThemes.Themes) >= 2 {
        score += 20
    }

    return min(100, score)
}

// scoreContentQuality 评分：内容质量
func (v *AIQualityValidator) scoreContentQuality(char *models.Character) float64 {
    score := 100.0

    // 检查是否包含 AI 生成痕迹
    texts := []string{
        char.Background,
        char.SpeechStyle,
    }

    for _, text := range texts {
        for _, pattern := range v.forbiddenPatterns {
            if pattern.MatchString(text) {
                score -= 20
                break
            }
        }
    }

    // 检查文本是否有意义（非重复）
    if strings.Contains(char.Background, strings.Repeat(char.Background[:10], 3)) {
        score -= 30
    }

    return max(0, score)
}
```

### 2.4 自动化质量门禁

```yaml
# File: .github/workflows/character-quality-gate.yml

name: Character Quality Gate

on:
  pull_request:
    paths:
      - 'data/characters/**/*.json'
      - 'backend/internal/models/character*.go'

jobs:
  validate-characters:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run Schema Validation
        run: |
          go test ./tests/validators/... -v -run TestCharacterValidator

      - name: Run Business Rules Validation
        run: |
          go test ./tests/validators/... -v -run TestBusinessRules

      - name: Run Quality Scoring
        run: |
          go test ./tests/validators/... -v -run TestAIQuality

      - name: Check Quality Threshold
        run: |
          # 所有角色质量分必须 >= 70
          go run ./cmd/quality_check.go --threshold=70 ./data/characters/
```

---

## 3. API 测试用例设计

### 3.1 测试框架设计

```
tests/
├── api/
│   ├── fixtures/           # 测试数据
│   │   ├── characters.json
│   │   └── requests.json
│   ├── helpers/            # 测试辅助函数
│   │   ├── client.go
│   │   └── assertions.go
│   ├── game_api_test.go    # 游戏 API 测试
│   ├── character_api_test.go # 角色 API 测试
│   ├── ai_api_test.go      # AI API 测试
│   └── websocket_test.go   # WebSocket 测试
└── integration/
    └── e2e_test.go
```

### 3.2 游戏 API 测试用例

#### 3.2.1 游戏状态 API 测试

```go
// File: tests/api/game_api_test.go

package api_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "stardew-agent/internal/api"
    "stardew-agent/internal/game"
)

func TestGetGameState(t *testing.T) {
    // 设置测试环境
    engine := game.NewEngine()
    handler := api.NewHandler(engine)
    router := setupTestRouter(handler)

    testCases := []struct {
        name           string
        endpoint       string
        expectedStatus int
        validateBody   func(t *testing.T, body map[string]interface{})
    }{
        {
            name:           "获取完整游戏状态",
            endpoint:       "/api/v1/game/state",
            expectedStatus: http.StatusOK,
            validateBody: func(t *testing.T, body map[string]interface{}) {
                assert.Contains(t, body, "player")
                assert.Contains(t, body, "time")
                assert.Contains(t, body, "npcs")

                player := body["player"].(map[string]interface{})
                assert.Contains(t, player, "position")
                assert.Contains(t, player, "energy")
                assert.Contains(t, player, "gold")
            },
        },
        {
            name:           "获取观察状态",
            endpoint:       "/api/v1/game/observation",
            expectedStatus: http.StatusOK,
            validateBody: func(t *testing.T, body map[string]interface{}) {
                assert.Contains(t, body, "player")
                assert.Contains(t, body, "nearby_npcs")
                assert.Contains(t, body, "messages")
            },
        },
        {
            name:           "获取地图数据",
            endpoint:       "/api/v1/game/map",
            expectedStatus: http.StatusOK,
            validateBody: func(t *testing.T, body map[string]interface{}) {
                assert.Contains(t, body, "width")
                assert.Contains(t, body, "height")
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", tc.endpoint, nil)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var body map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &body)
            tc.validateBody(t, body)
        })
    }
}

func TestGameActions(t *testing.T) {
    engine := game.NewEngine()
    handler := api.NewHandler(engine)
    router := setupTestRouter(handler)

    t.Run("移动动作测试", func(t *testing.T) {
        testCases := []struct {
            name         string
            direction    string
            expectSuccess bool
        }{
            {"向上移动", "up", true},
            {"向下移动", "down", true},
            {"向左移动", "left", true},
            {"向右移动", "right", true},
            {"无效方向", "invalid", false},
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                body := map[string]interface{}{
                    "type": "move",
                    "params": map[string]string{
                        "direction": tc.direction,
                    },
                }
                bodyBytes, _ := json.Marshal(body)

                req := httptest.NewRequest("POST", "/api/v1/game/action", bytes.NewReader(bodyBytes))
                req.Header.Set("Content-Type", "application/json")
                w := httptest.NewRecorder()
                router.ServeHTTP(w, req)

                var resp map[string]interface{}
                json.Unmarshal(w.Body.Bytes(), &resp)

                assert.Equal(t, tc.expectSuccess, resp["success"])
            })
        }
    })

    t.Run("NPC交互测试", func(t *testing.T) {
        // 测试与 NPC 对话
        t.Run("与NPC对话", func(t *testing.T) {
            // 先移动到 NPC 附近
            moveToNPC(engine, "lewis")

            body := map[string]interface{}{
                "type": "talk",
                "params": map[string]string{
                    "npc": "lewis",
                },
            }
            bodyBytes, _ := json.Marshal(body)

            req := httptest.NewRequest("POST", "/api/v1/game/action", bytes.NewReader(bodyBytes))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            var resp map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &resp)

            assert.True(t, resp["success"].(bool))
            assert.NotEmpty(t, resp["message"])
        })

        // 测试送礼
        t.Run("赠送礼物", func(t *testing.T) {
            moveToNPC(engine, "haley")

            body := map[string]interface{}{
                "type": "give_gift",
                "params": map[string]string{
                    "npc":       "haley",
                    "gift_item": "flower",
                },
            }
            bodyBytes, _ := json.Marshal(body)

            req := httptest.NewRequest("POST", "/api/v1/game/action", bytes.NewReader(bodyBytes))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            var resp map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &resp)

            // 验证友谊值增加
            if resp["success"].(bool) {
                events := resp["events"].([]interface{})
                assert.NotEmpty(t, events)
            }
        })
    })

    t.Run("边界条件测试", func(t *testing.T) {
        // 测试地图边界移动
        t.Run("尝试移出地图边界", func(t *testing.T) {
            // 重置游戏，玩家在初始位置
            engine.Reset()

            // 尝试向左上方移动多次（尝试移出边界）
            for i := 0; i < 10; i++ {
                body := map[string]interface{}{
                    "type": "move",
                    "params": map[string]string{"direction": "left"},
                }
                bodyBytes, _ := json.Marshal(body)
                req := httptest.NewRequest("POST", "/api/v1/game/action", bytes.NewReader(bodyBytes))
                req.Header.Set("Content-Type", "application/json")
                w := httptest.NewRecorder()
                router.ServeHTTP(w, req)
            }

            state := engine.GetState()
            // 验证玩家位置在有效范围内
            assert.GreaterOrEqual(t, state.Player.Position.X, 0)
            assert.GreaterOrEqual(t, state.Player.Position.Y, 0)
        })
    })
}
```

### 3.3 角色 API 测试用例

```go
// File: tests/api/character_api_test.go

package api_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCharacterCRUD(t *testing.T) {
    router := setupTestRouter(nil)

    t.Run("创建角色", func(t *testing.T) {
        char := map[string]interface{}{
            "name":   "TestCharacter",
            "role":   "villager",
            "gender": "male",
            "personality": map[string]interface{}{
                "primary_trait": "friendly",
            },
        }
        bodyBytes, _ := json.Marshal(char)

        req := httptest.NewRequest("POST", "/api/v1/characters", bytes.NewReader(bodyBytes))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusCreated, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.True(t, resp["success"].(bool))
        assert.NotEmpty(t, resp["character"].(map[string]interface{})["id"])
    })

    t.Run("获取角色列表", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/characters", nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.Contains(t, resp, "characters")
        assert.Contains(t, resp, "pagination")
    })

    t.Run("获取单个角色", func(t *testing.T) {
        // 先创建一个角色
        charID := createTestCharacter(t, router)

        req := httptest.NewRequest("GET", "/api/v1/characters/"+charID, nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.Equal(t, charID, resp["character"].(map[string]interface{})["id"])
    })

    t.Run("更新角色", func(t *testing.T) {
        charID := createTestCharacter(t, router)

        updates := map[string]interface{}{
            "name": "UpdatedName",
        }
        bodyBytes, _ := json.Marshal(updates)

        req := httptest.NewRequest("PUT", "/api/v1/characters/"+charID, bytes.NewReader(bodyBytes))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("删除角色", func(t *testing.T) {
        charID := createTestCharacter(t, router)

        req := httptest.NewRequest("DELETE", "/api/v1/characters/"+charID, nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        // 验证已删除
        req = httptest.NewRequest("GET", "/api/v1/characters/"+charID, nil)
        w = httptest.NewRecorder()
        router.ServeHTTP(w, req)
        assert.Equal(t, http.StatusNotFound, w.Code)
    })
}

func TestCharacterGeneration(t *testing.T) {
    router := setupTestRouter(nil)

    t.Run("AI生成角色-异步", func(t *testing.T) {
        req := map[string]interface{}{
            "role":          "villager",
            "gender":        "female",
            "generate_schedule": true,
            "generate_dialogue": true,
        }
        bodyBytes, _ := json.Marshal(req)

        httpReq := httptest.NewRequest("POST", "/api/v1/characters/generate", bytes.NewReader(bodyBytes))
        httpReq.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, httpReq)

        assert.Equal(t, http.StatusAccepted, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NotEmpty(t, resp["job_id"])

        // 轮询任务状态
        jobID := resp["job_id"].(string)
        eventually := func() bool {
            req := httptest.NewRequest("GET", "/api/v1/characters/generate/"+jobID, nil)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            var status map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &status)
            return status["status"] == "completed"
        }

        assert.Eventually(t, eventually, 30*time.Second, 500*time.Millisecond)
    })

    t.Run("生成角色对话", func(t *testing.T) {
        charID := createTestCharacter(t, router)

        req := map[string]interface{}{
            "context":      "first_meeting",
            "player_input": "Hello!",
        }
        bodyBytes, _ := json.Marshal(req)

        httpReq := httptest.NewRequest("POST", "/api/v1/characters/"+charID+"/dialogue", bytes.NewReader(bodyBytes))
        httpReq.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, httpReq)

        assert.Equal(t, http.StatusOK, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NotEmpty(t, resp["dialogue"])
    })
}

// API 测试用例矩阵

/*
+-------------------+-------------------+-------------------+-------------------+
| 端点              | 正常场景          | 边界场景          | 错误场景          |
+-------------------+-------------------+-------------------+-------------------+
| POST /characters  | 创建完整角色      | 最小必需字段      | 缺少必需字段      |
|                   | 带所有可选字段    | 字段最大长度      | 无效字段值        |
|                   |                   |                   | 重复ID            |
+-------------------+-------------------+-------------------+-------------------+
| GET /characters   | 无过滤获取        | 大量数据分页      | 无效页码          |
|                   | 按角色类型过滤    |                   | 无效过滤参数      |
|                   | 按名称搜索        |                   |                   |
+-------------------+-------------------+-------------------+-------------------+
| POST /generate    | 全参数生成        | 仅必需参数        | 无效AI Provider   |
|                   | 指定Provider      | 低创造性参数      | AI服务不可用      |
|                   |                   |                   | 超时              |
+-------------------+-------------------+-------------------+-------------------+
*/
```

### 3.4 AI API 测试用例

```go
// File: tests/api/ai_api_test.go

package api_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAIProviderSwitching(t *testing.T) {
    router := setupTestRouter(nil)

    /*
    AI Provider 测试场景:

    1. Cloud API (Claude/OpenAI)
       - 正常响应
       - 速率限制
       - 认证失败
       - 超时

    2. SD API (本地)
       - 正常响应
       - 服务不可用
       - 模型加载失败

    3. Mock Provider
       - 固定响应
       - 可配置延迟

    4. 降级链测试
       - Cloud -> SD -> Mock
       - 验证降级触发条件
    */

    t.Run("测试 Mock Provider", func(t *testing.T) {
        req := map[string]interface{}{
            "provider": "mock",
            "prompt":   "Generate a character",
        }
        bodyBytes, _ := json.Marshal(req)

        httpReq := httptest.NewRequest("POST", "/api/v1/ai/dialogue", bytes.NewReader(bodyBytes))
        httpReq.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, httpReq)

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("测试降级机制", func(t *testing.T) {
        // 配置: cloud_api -> sd_api -> mock
        // 当 cloud_api 不可用时，应该自动降级

        // 模拟 Cloud API 故障
        // 验证降级到下一个 Provider
    })
}

func TestAIEndpoints(t *testing.T) {
    router := setupTestRouter(nil)

    t.Run("AI对话生成", func(t *testing.T) {
        req := map[string]interface{}{
            "npc_id":       "lewis",
            "context":      "greeting",
            "player_input": "Hello Mayor!",
            "mood":         "neutral",
        }
        bodyBytes, _ := json.Marshal(req)

        httpReq := httptest.NewRequest("POST", "/api/v1/ai/dialogue", bytes.NewReader(bodyBytes))
        httpReq.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, httpReq)

        assert.Equal(t, http.StatusOK, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NotEmpty(t, resp["dialogue"])
    })

    t.Run("AI行为决策", func(t *testing.T) {
        req := map[string]interface{}{
            "npc_id":          "haley",
            "weather":         "sunny",
            "current_time":    map[string]int{"hour": 10, "minute": 0},
            "player_position": map[string]int{"x": 20, "y": 20},
        }
        bodyBytes, _ := json.Marshal(req)

        httpReq := httptest.NewRequest("POST", "/api/v1/ai/decision", bytes.NewReader(bodyBytes))
        httpReq.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, httpReq)

        assert.Equal(t, http.StatusOK, w.Code)

        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        assert.Contains(t, resp, "behavior")
        assert.Contains(t, resp, "confidence")
    })

    t.Run("自然语言指令解析", func(t *testing.T) {
        testCases := []struct {
            input    string
            expected map[string]interface{}
        }{
            {
                input: "go to the shop",
                expected: map[string]interface{}{
                    "type": "move",
                    "destination": "shop",
                },
            },
            {
                input: "talk to Lewis",
                expected: map[string]interface{}{
                    "type": "talk",
                    "npc":  "lewis",
                },
            },
            {
                input: "plant some potatoes",
                expected: map[string]interface{}{
                    "type": "plant",
                    "seed": "potato",
                },
            },
        }

        for _, tc := range testCases {
            req := map[string]interface{}{
                "command": tc.input,
            }
            bodyBytes, _ := json.Marshal(req)

            httpReq := httptest.NewRequest("POST", "/api/v1/ai/interpret", bytes.NewReader(bodyBytes))
            httpReq.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, httpReq)

            var resp map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &resp)

            // 验证解析结果符合预期
            if action, ok := resp["action"].(map[string]interface{}); ok {
                assert.Equal(t, tc.expected["type"], action["type"])
            }
        }
    })
}
```

### 3.5 性能测试

```go
// File: tests/api/performance_test.go

package api_test

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestAPIPerformance(t *testing.T) {
    router := setupTestRouter(nil)

    /*
    性能基准:

    | 端点                | P50 响应时间 | P95 响应时间 | P99 响应时间 |
    |---------------------|-------------|-------------|-------------|
    | GET /game/state     | < 10ms      | < 30ms      | < 50ms      |
    | POST /game/action   | < 20ms      | < 50ms      | < 100ms     |
    | POST /characters    | < 30ms      | < 80ms      | < 150ms     |
    | POST /ai/dialogue   | < 500ms     | < 2000ms    | < 5000ms    |
    */

    t.Run("响应时间基准测试", func(t *testing.T) {
        endpoints := []struct {
            name     string
            method   string
            path     string
            maxP95   time.Duration
        }{
            {"GetGameState", "GET", "/api/v1/game/state", 30 * time.Millisecond},
            {"ExecuteAction", "POST", "/api/v1/game/action", 50 * time.Millisecond},
        }

        for _, ep := range endpoints {
            t.Run(ep.name, func(t *testing.T) {
                times := make([]time.Duration, 100)

                for i := 0; i < 100; i++ {
                    start := time.Now()

                    var req *http.Request
                    if ep.method == "GET" {
                        req = httptest.NewRequest("GET", ep.path, nil)
                    } else {
                        body := `{"type":"move","params":{"direction":"up"}}`
                        req = httptest.NewRequest("POST", ep.path, strings.NewReader(body))
                        req.Header.Set("Content-Type", "application/json")
                    }

                    w := httptest.NewRecorder()
                    router.ServeHTTP(w, req)

                    times[i] = time.Since(start)
                }

                // 计算 P95
                sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
                p95 := times[94] // 95th percentile

                assert.Less(t, p95, ep.maxP95,
                    fmt.Sprintf("P95 response time %v exceeds threshold %v", p95, ep.maxP95))
            })
        }
    })

    t.Run("并发测试", func(t *testing.T) {
        concurrency := 50
        requestsPerWorker := 20

        var wg sync.WaitGroup
        errors := make(chan error, concurrency)

        for i := 0; i < concurrency; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for j := 0; j < requestsPerWorker; j++ {
                    req := httptest.NewRequest("GET", "/api/v1/game/state", nil)
                    w := httptest.NewRecorder()
                    router.ServeHTTP(w, req)

                    if w.Code != http.StatusOK {
                        errors <- fmt.Errorf("unexpected status: %d", w.Code)
                        return
                    }
                }
            }()
        }

        wg.Wait()
        close(errors)

        var errorList []error
        for err := range errors {
            errorList = append(errorList, err)
        }

        assert.Empty(t, errorList, "Concurrent requests should not fail")
    })
}
```

---

## 4. 预生成脚本验证方案

### 4.1 预生成脚本概述

预生成脚本用于在游戏启动前批量生成角色数据，包括：
- NPC 基础属性
- 日程安排
- 对话模板
- 礼物偏好

### 4.2 脚本验证测试

```go
// File: tests/pregen/pregen_test.go

package pregen_test

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "stardew-agent/internal/validators"
    "stardew-agent/scripts/pregen"
)

func TestPregeneratedCharacters(t *testing.T) {
    validator, err := validators.NewCharacterValidator("schemas/character.json")
    require.NoError(t, err)
    qualityValidator := validators.NewAIQualityValidator()

    // 遍历所有预生成角色
    charDir := "data/characters"
    files, err := filepath.Glob(filepath.Join(charDir, "*.json"))
    require.NoError(t, err)

    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            data, err := os.ReadFile(file)
            require.NoError(t, err)

            // 1. Schema 验证
            result, err := validator.Validate(data)
            require.NoError(t, err)
            assert.True(t, result.Valid, "Schema validation failed: %v", result.Errors)

            // 2. 解析角色
            var char models.Character
            err = json.Unmarshal(data, &char)
            require.NoError(t, err)

            // 3. 业务规则验证
            businessValidator := &validators.BusinessRuleValidator{}
            err = businessValidator.ValidateAppearanceConsistency(&char)
            assert.NoError(t, err)

            err = businessValidator.ValidateScheduleValidity(&char)
            assert.NoError(t, err)

            // 4. 质量评分
            qualityScore := qualityValidator.ValidateCharacterQuality(&char)
            assert.True(t, qualityScore.PassThreshold,
                "Quality score %.2f below threshold 70. Issues: %v",
                qualityScore.Overall, qualityScore.Issues)
        })
    }
}

func TestPregenerationScript(t *testing.T) {
    // 测试预生成脚本功能

    t.Run("生成NPC配置", func(t *testing.T) {
        generator := pregen.NewCharacterGenerator(&pregen.Config{
            OutputDir:    t.TempDir(),
            ProviderType: "mock",
            Count:        5,
        })

        results := generator.GenerateBatch("villager", 5)

        assert.Len(t, results, 5)
        for _, result := range results {
            assert.NoError(t, result.Error)
            assert.NotNil(t, result.Character)
            assert.NotEmpty(t, result.Character.ID)
        }
    })

    t.Run("生成日程", func(t *testing.T) {
        scheduler := pregen.NewScheduleGenerator()

        schedule := scheduler.Generate(&models.Character{
            Role: "shopkeeper",
            Personality: models.Personality{
                PrimaryTrait: "serious",
            },
        })

        // 验证日程完整性
        assert.NotEmpty(t, schedule)

        // 验证时间覆盖 (6:00 - 26:00)
        hours := make(map[int]bool)
        for _, entry := range schedule {
            for h := entry.StartHour; h < entry.EndHour; h++ {
                hours[h] = true
            }
        }

        // 检查主要时间段覆盖
        for h := 9; h <= 18; h++ {
            assert.True(t, hours[h], "Hour %d should be covered in schedule", h)
        }
    })

    t.Run("生成对话模板", func(t *testing.T) {
        dialogueGen := pregen.NewDialogueGenerator()

        dialogues := dialogueGen.Generate(&models.Character{
            Name: "TestNPC",
            Personality: models.Personality{
                PrimaryTrait: "friendly",
            },
            SpeechStyle: "casual",
        })

        // 验证对话类型
        assert.NotEmpty(t, dialogues.Greetings)
        assert.NotEmpty(t, dialogues.Farewells)

        // 验证对话不包含占位符
        for _, d := range dialogues.Greetings {
            assert.NotContains(t, d, "{")
            assert.NotContains(t, d, "}")
        }
    })
}

func TestPregenerationIdempotency(t *testing.T) {
    // 测试预生成幂等性

    config := &pregen.Config{
        Seed:        12345, // 固定种子
        OutputDir:   t.TempDir(),
        ProviderType: "mock",
    }

    generator1 := pregen.NewCharacterGenerator(config)
    generator2 := pregen.NewCharacterGenerator(config)

    char1 := generator1.GenerateOne("villager")
    char2 := generator2.GenerateOne("villager")

    // 相同种子应该产生相同结果
    assert.Equal(t, char1.Name, char2.Name)
    assert.Equal(t, char1.Personality.PrimaryTrait, char2.Personality.PrimaryTrait)
}
```

### 4.3 预生成数据验证流程

```
┌─────────────────────────────────────────────────────────────────────┐
│                      预生成数据验证流程                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│   │ JSON 语法   │───>│ Schema 验证 │───>│ 业务规则   │            │
│   │ 解析检查    │    │             │    │ 验证        │            │
│   └─────────────┘    └─────────────┘    └─────────────┘            │
│          │                  │                  │                    │
│          v                  v                  v                    │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│   │ 语法错误    │    │ 类型错误    │    │ 逻辑错误    │            │
│   │ 报告        │    │ 报告        │    │ 报告        │            │
│   └─────────────┘    └─────────────┘    └─────────────┘            │
│                                                                     │
│   ┌─────────────────────────────────────────────────────────┐      │
│   │                    质量评估                              │      │
│   │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │      │
│   │  │ 完整性  │  │ 一致性  │  │ 丰富度  │  │ 内容    │   │      │
│   │  │ 30%     │  │ 25%     │  │ 25%     │  │ 质量20% │   │      │
│   │  └─────────┘  └─────────┘  └─────────┘  └─────────┘   │      │
│   └─────────────────────────────────────────────────────────┘      │
│                          │                                         │
│                          v                                         │
│   ┌─────────────────────────────────────────────────────────┐      │
│   │                 质量分 >= 70 ?                           │      │
│   │                      │                                  │      │
│   │            ┌─────────┴─────────┐                        │      │
│   │            v                   v                        │      │
│   │     ┌───────────┐       ┌───────────┐                  │      │
│   │     │   通过    │       │   拒绝    │                  │      │
│   │     │  部署     │       │  重新生成 │                  │      │
│   │     └───────────┘       └───────────┘                  │      │
│   └─────────────────────────────────────────────────────────┘      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 5. 测试数据准备建议

### 5.1 测试数据分类

```
tests/fixtures/
├── characters/
│   ├── valid/
│   │   ├── character_minimal.json      # 最小必需字段
│   │   ├── character_full.json         # 完整字段
│   │   └── character_edge.json         # 边界值
│   ├── invalid/
│   │   ├── missing_name.json           # 缺少必需字段
│   │   ├── invalid_gender.json         # 无效枚举值
│   │   └── out_of_range.json           # 超出范围值
│   └── ai_generated/
│       ├── high_quality.json           # 质量分 >= 90
│       ├── medium_quality.json         # 质量分 70-90
│       └── low_quality.json            # 质量分 < 70
├── game_states/
│   ├── initial_state.json              # 初始状态
│   ├── mid_game.json                   # 游戏中期
│   └── end_of_day.json                 # 一天结束
├── api_requests/
│   ├── move_actions.json               # 移动动作集合
│   ├── npc_interactions.json           # NPC交互集合
│   └── quest_operations.json           # 任务操作集合
└── mock_responses/
    ├── ai_dialogue_success.json        # AI对话成功响应
    ├── ai_dialogue_timeout.json        # AI对话超时
    └── ai_character_generated.json     # AI角色生成响应
```

### 5.2 测试数据生成器

```go
// File: tests/fixtures/generator.go

package fixtures

import (
    "math/rand"
    "time"

    "stardew-agent/internal/models"
)

// CharacterBuilder 角色构建器
type CharacterBuilder struct {
    character *models.Character
}

func NewCharacterBuilder() *CharacterBuilder {
    return &CharacterBuilder{
        character: &models.Character{
            ID:        generateID(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }
}

func (b *CharacterBuilder) WithName(name string) *CharacterBuilder {
    b.character.Name = name
    return b
}

func (b *CharacterBuilder) WithRole(role string) *CharacterBuilder {
    b.character.Role = role
    return b
}

func (b *CharacterBuilder) WithGender(gender string) *CharacterBuilder {
    b.character.Gender = gender
    return b
}

func (b *CharacterBuilder) WithPersonality(trait string) *CharacterBuilder {
    b.character.Personality = models.Personality{
        PrimaryTrait: trait,
    }
    return b
}

func (b *CharacterBuilder) WithFullAppearance() *CharacterBuilder {
    b.character.Appearance = models.Appearance{
        Hairstyle: models.Hairstyle{
            Style: randomChoice(hairstyles),
            Color: randomColor(),
        },
        Eyes: models.Eyes{
            Shape: randomChoice(eyeShapes),
            Color: randomColor(),
        },
        Skin: models.Skin{
            Tone: randomColor(),
        },
        BodyType: models.BodyType{
            Height: 0.8 + rand.Float64()*0.4,
            Build:  0.7 + rand.Float64()*0.6,
        },
    }
    return b
}

func (b *CharacterBuilder) Build() *models.Character {
    return b.character
}

// 预定义测试数据生成函数

func ValidMinimalCharacter() *models.Character {
    return NewCharacterBuilder().
        WithName("Test").
        WithRole("villager").
        WithGender("male").
        WithPersonality("friendly").
        Build()
}

func ValidFullCharacter() *models.Character {
    return NewCharacterBuilder().
        WithName("CompleteCharacter").
        WithRole("shopkeeper").
        WithGender("female").
        WithPersonality("serious").
        WithFullAppearance().
        Build()
}

// GameStateBuilder 游戏状态构建器
type GameStateBuilder struct {
    state *models.GameState
}

func NewGameStateBuilder() *GameStateBuilder {
    return &GameStateBuilder{
        state: &models.GameState{
            Player: models.PlayerState{
                Position:     models.Position{X: 24, Y: 24},
                Energy:       100,
                MaxEnergy:    100,
                Gold:         500,
                Inventory:    []models.InventoryItem{},
                Tools:        []string{"hoe", "watering_can", "axe", "pickaxe", "scythe"},
                ActiveTool:   "hoe",
                Friendship:   make(map[string]int),
            },
            Time: models.GameTime{
                Day:     1,
                Hour:    6,
                Minute:  0,
                Season:  "spring",
                Year:    1,
                Paused:  false,
            },
            NPCs:      []models.NPCState{},
            Quests:    []models.QuestData{},
        },
    }
}

func (b *GameStateBuilder) WithPlayerPosition(x, y int) *GameStateBuilder {
    b.state.Player.Position = models.Position{X: x, Y: y}
    return b
}

func (b *GameStateBuilder) WithTime(hour, minute int) *GameStateBuilder {
    b.state.Time.Hour = hour
    b.state.Time.Minute = minute
    return b
}

func (b *GameStateBuilder) WithNPC(npc models.NPCState) *GameStateBuilder {
    b.state.NPCs = append(b.state.NPCs, npc)
    return b
}

func (b *GameStateBuilder) Build() *models.GameState {
    return b.state
}
```

### 5.3 Mock 数据配置

```yaml
# File: tests/fixtures/mock_config.yaml

mock_providers:
  ai_service:
    # 模拟延迟
    latency_ms: 100

    # 预设响应
    responses:
      dialogue:
        - text: "Hello, traveler!"
          mood_change: "+1"
        - text: "Nice weather today."
          mood_change: "0"

      character:
        - name: "Generated NPC"
          role: "villager"
          personality:
            primary_trait: "friendly"

  # 错误模拟
  error_simulation:
    - type: timeout
      probability: 0.1
      delay_ms: 5000

    - type: rate_limit
      probability: 0.05
      status_code: 429

    - type: server_error
      probability: 0.02
      status_code: 500
```

### 5.4 数据库测试数据

```javascript
// File: tests/fixtures/database/seed.js

// MongoDB 测试数据种子

db.characters.insertMany([
  {
    _id: "test_npc_001",
    name: "Test Lewis",
    role: "mayor",
    gender: "male",
    personality: {
      primary_trait: "friendly",
      secondary_traits: ["hardworking", "honest"]
    },
    is_active: true,
    created_by: "test_seed"
  },
  {
    _id: "test_npc_002",
    name: "Test Haley",
    role: "villager",
    gender: "female",
    personality: {
      primary_trait: "cheerful"
    },
    is_active: true,
    created_by: "test_seed"
  }
]);

db.quests.insertMany([
  {
    _id: "test_quest_001",
    name: "Test Quest",
    description: "A test quest for validation",
    status: "available",
    objectives: [
      { type: "talk", target: "test_npc_001", count: 1 }
    ]
  }
]);
```

---

## 6. 测试执行策略

### 6.1 测试层次金字塔

```
                    ┌─────────┐
                    │   E2E   │  <- 少量，关键流程
                    │  Tests  │
                   ─┴─────────┴─
                  ┌─────────────┐
                  │ Integration │  <- API集成，服务交互
                  │    Tests    │
                 ─┴─────────────┴─
                ┌─────────────────┐
                │   Unit Tests    │  <- 业务逻辑，验证器
                │                 │
               ─┴─────────────────┴─
```

### 6.2 CI/CD 集成

```yaml
# File: .github/workflows/test.yml

name: Test Suite

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run Unit Tests
        run: |
          go test ./... -v -coverprofile=coverage.out
          go tool cover -func=coverage.out

      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out

  api-tests:
    runs-on: ubuntu-latest
    services:
      mongodb:
        image: mongo:7
        ports:
          - 27017:27017
      redis:
        image: redis:7
        ports:
          - 6379:6379

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run API Tests
        run: |
          go test ./tests/api/... -v
        env:
          MONGODB_URI: mongodb://localhost:27017
          REDIS_ADDR: localhost:6379

  quality-gate:
    runs-on: ubuntu-latest
    needs: [unit-tests, api-tests]
    steps:
      - name: Check Quality Gate
        run: |
          echo "All tests passed. Quality gate satisfied."
```

### 6.3 测试覆盖率目标

| 模块 | 行覆盖率 | 分支覆盖率 |
|------|---------|-----------|
| Validators | >= 90% | >= 85% |
| API Handlers | >= 80% | >= 75% |
| AI Service | >= 70% | >= 65% |
| Game Engine | >= 75% | >= 70% |
| **Overall** | **>= 80%** | **>= 75%** |

---

## 7. 缺陷管理流程

### 7.1 Bug 报告模板

```markdown
## Bug 报告

### 基本信息
- **发现者**:
- **日期**:
- **版本**:
- **环境**: [开发/测试/生产]

### 描述
[简洁描述问题]

### 复现步骤
1.
2.
3.

### 预期结果
[应该发生什么]

### 实际结果
[实际发生了什么]

### 影响范围
- [ ] 功能阻塞
- [ ] 功能降级
- [ ] UI问题
- [ ] 性能问题

### 附件
- 截图
- 日志
- 相关数据文件
```

### 7.2 严重程度分级

| 级别 | 描述 | 响应时间 | 解决时限 |
|------|------|---------|---------|
| P0 - 致命 | 服务不可用，数据丢失 | 1小时 | 4小时 |
| P1 - 严重 | 核心功能不可用 | 4小时 | 24小时 |
| P2 - 一般 | 功能受限，有替代方案 | 24小时 | 3天 |
| P3 - 轻微 | UI问题，优化建议 | 3天 | 1周 |

---

## 8. 附录

### 8.1 测试工具清单

| 工具 | 用途 | 版本 |
|------|------|------|
| Go testing | 单元测试框架 | 1.21+ |
| testify | 断言库 | latest |
| gojsonschema | JSON Schema 验证 | latest |
| httptest | HTTP 测试 | stdlib |
| MongoDB Testcontainers | 数据库测试 | latest |

### 8.2 相关文档链接

- [API Documentation](./API.md)
- [Architecture](./ARCHITECTURE.md)
- [Character Generator Design](./CHARACTER_GENERATOR_DESIGN.md)
- [Backend Service Architecture](./BACKEND_SERVICE_ARCHITECTURE.md)

---

*文档版本: 1.0*
*创建日期: 2026-02-12*
*作者: QA Engineer*
