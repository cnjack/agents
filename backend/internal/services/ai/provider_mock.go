package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"stardew-agent/internal/config"
	"stardew-agent/internal/models"
)

// MockProvider Mock AI 提供者（用于开发测试）
type MockProvider struct {
	config     config.ProviderConfig
	delay      time.Duration
	dialogues  map[string]map[string][]string
}

// NewMockProvider 创建 Mock 提供者
func NewMockProvider(cfg config.ProviderConfig) *MockProvider {
	delay := time.Duration(cfg.DelayMs) * time.Millisecond
	if delay == 0 {
		delay = 100 * time.Millisecond
	}

	return &MockProvider{
		config: cfg,
		delay:  delay,
		dialogues: getMockDialogues(),
	}
}

func getMockDialogues() map[string]map[string][]string {
	return map[string]map[string][]string{
		"Lewis": {
			"morning": {
				"早安!美好的一天开始了,我们小镇需要像你这样勤劳的农民。",
				"今天天气不错,适合在社区中心附近走走。",
				"亲爱的朋友,有什么需要我帮忙的吗?",
			},
			"afternoon": {
				"下午好!社区中心的重建工作进行得怎么样了?",
				"我们小镇的发展离不开大家的努力啊。",
				"最近有什么新鲜事吗?",
			},
			"evening": {
				"晚上好,辛苦了一天,好好休息吧。",
				"夜晚的小镇特别宁静,不是吗?",
			},
		},
		"Pierre": {
			"morning": {
				"早安!今天有什么需要买的吗?新进的货物很棒!",
				"欢迎光临!早起的鸟儿有虫吃,早来的顾客有优惠!",
			},
			"afternoon": {
				"来看看今天的特价商品吧!",
				"刚从城里进了新货,品质一流!",
				"最近生意还不错,谢谢大家支持。",
			},
			"evening": {
				"快打烊了,还有什么需要的吗?",
				"辛苦了,明天见!",
			},
		},
		"Robin": {
			"morning": {
				"早!今天有什么建筑项目需要帮忙吗?",
				"早上空气好,适合在工坊干活。",
			},
			"afternoon": {
				"需要建造什么?我可以帮你设计!",
				"来看看我最近的作品,怎么样?",
				"建筑是一门艺术,需要用心对待。",
			},
			"evening": {
				"辛苦了一天,该休息了。",
				"晚上好好休息,明天继续努力!",
			},
		},
		"Haley": {
			"morning": {
				"...你想干嘛?",
				"嗯...早。",
			},
			"afternoon": {
				"我在拍照,别打扰我。",
				"这个角度的光线很完美...",
			},
			"evening": {
				"晚上好...还是没什么好聊的。",
			},
		},
		"Willy": {
			"morning": {
				"早安!今天是个钓鱼的好日子!",
				"海风很舒服,适合出海。",
			},
			"afternoon": {
				"钓到什么好鱼了吗?",
				"来聊聊钓鱼心得吧!",
				"海洋总是充满惊喜。",
			},
			"evening": {
				"晚上好!要去酒吧喝一杯吗?",
				"一天的收获怎么样?",
			},
		},
	}
}

// GenerateDialogue 生成对话（Mock 实现）
func (p *MockProvider) GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error) {
	// 模拟延迟
	time.Sleep(p.delay)

	charName := "NPC"
	dayPhase := "afternoon"
	hearts := req.FriendshipHearts

	if req.Character != nil {
		charName = req.Character.Name
	}

	if req.GameState != nil {
		hour := req.GameState.Time.Hour
		if hour >= 6 && hour < 12 {
			dayPhase = "morning"
		} else if hour >= 12 && hour < 18 {
			dayPhase = "afternoon"
		} else {
			dayPhase = "evening"
		}
	}

	dialogue := p.getMockDialogue(charName, dayPhase, hearts)

	return &DialogueResponse{
		Dialogue:        dialogue,
		MoodChange:      "neutral",
		FriendshipDelta: 0,
		SuggestedActions: []string{"talk", "give_gift"},
	}, nil
}

func (p *MockProvider) getMockDialogue(npcName, dayPhase string, hearts int) string {
	rand.Seed(time.Now().UnixNano())

	if npcDialogues, ok := p.dialogues[npcName]; ok {
		if phaseDialogues, ok := npcDialogues[dayPhase]; ok && len(phaseDialogues) > 0 {
			if hearts >= 6 && npcName == "Haley" {
				return fmt.Sprintf("其实...看到你挺好的。%s", phaseDialogues[rand.Intn(len(phaseDialogues))])
			}
			return phaseDialogues[rand.Intn(len(phaseDialogues))]
		}
	}

	return "你好,有什么事吗?"
}

// GenerateCharacter 生成角色（Mock 实现）
func (p *MockProvider) GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error) {
	time.Sleep(p.delay)

	// 生成一个简单的 Mock 角色
	name := req.NameHint
	if name == "" {
		names := []string{"艾米", "杰克", "莉莉", "汤姆", "苏珊"}
		name = names[rand.Intn(len(names))]
	}

	role := req.Role
	if role == "" {
		roles := []string{"villager", "farmer", "merchant", "fisherman"}
		role = roles[rand.Intn(len(roles))]
	}

	return &CharacterResponse{
		Name:       name,
		Role:       role,
		Age:        "young",
		Gender:     "unknown",
		Traits:     []string{"友好", "勤劳"},
		Values:     []string{"家庭", "友谊"},
		Dislikes:   []string{"懒惰", "欺骗"},
		Background: fmt.Sprintf("%s是一个普通的村民，过着平静的生活。", name),
		Secrets:    []string{"有一个不为人知的爱好"},
		SpeechStyle: "友好热情",
		Hobby:      "钓鱼",
		GiftPreferences: models.GiftPreferences{
			Love:    []string{"花", "手工制品"},
			Like:    []string{"水果", "蔬菜"},
			Neutral: []string{"石头", "木头"},
			Dislike: []string{"垃圾"},
			Hate:    []string{},
		},
		DialogueThemes: models.DialogueThemes{
			Morning:   []string{"早上的空气真好"},
			Afternoon: []string{"今天过得怎么样"},
			Evening:   []string{"晚上了，该休息了"},
		},
		CharacterArc: "从陌生人到好朋友",
	}, nil
}

// Understand 自然语言理解（Mock 实现）
func (p *MockProvider) Understand(ctx context.Context, req NLURequest) (*NLUResponse, error) {
	time.Sleep(p.delay)

	input := strings.ToLower(req.Input)

	// 简单的关键词匹配
	if strings.Contains(input, "你好") || strings.Contains(input, "hello") || strings.Contains(input, "hi") {
		return &NLUResponse{
			Understood:   true,
			Intent:       "talk",
			ResponseText: "你想和谁对话？",
			Action: &models.Action{
				Type: models.ActionTalk,
			},
		}, nil
	}

	if strings.Contains(input, "送礼") || strings.Contains(input, "gift") {
		return &NLUResponse{
			Understood:   true,
			Intent:       "gift",
			ResponseText: "你想送什么礼物给谁？",
			NeedsClarification: true,
			ClarificationQuestion: "请指定收礼人和礼物",
		}, nil
	}

	if strings.Contains(input, "移动") || strings.Contains(input, "move") || strings.Contains(input, "去") {
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
			if strings.Contains(input, dirWord) {
				return &NLUResponse{
					Understood:   true,
					Intent:       "move",
					ResponseText: fmt.Sprintf("向%s移动", dirWord),
					Action: &models.Action{
						Type: models.ActionMove,
						Params: models.ActionParams{
							Direction: dir,
						},
					},
				}, nil
			}
		}
	}

	// 默认返回
	return &NLUResponse{
		Understood:         false,
		Intent:             "unknown",
		ResponseText:       "我不太明白你的意思",
		NeedsClarification: true,
		ClarificationQuestion: "你想做什么？你可以：对话、送礼、移动",
	}, nil
}

// DecideBehavior 行为决策（Mock 实现）
func (p *MockProvider) DecideBehavior(ctx context.Context, req BehaviorRequest) (*BehaviorResponse, error) {
	time.Sleep(p.delay)

	// 简单的行为决策
	actions := req.PossibleActions
	if len(actions) == 0 {
		actions = []string{"idle", "wander", "greet"}
	}

	action := actions[rand.Intn(len(actions))]

	return &BehaviorResponse{
		Action:          action,
		TargetLocation:  "town",
		Reason:          "Mock decision based on random selection",
		MoodAfterAction: "neutral",
		CanInterrupt:    true,
	}, nil
}

// HealthCheck 健康检查
func (p *MockProvider) HealthCheck(ctx context.Context) error {
	return nil // Mock 提供者总是健康
}

// ProviderInfo 获取提供者信息
func (p *MockProvider) ProviderInfo() ProviderInfo {
	return ProviderInfo{
		Type:      ProviderTypeMock,
		Model:     "mock",
		Available: true,
	}
}

// ToJSON 辅助方法：将响应转为 JSON
func (r *DialogueResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *CharacterResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *NLUResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *BehaviorResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}
