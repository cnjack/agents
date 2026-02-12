package ai

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"stardew-agent/internal/config"
)

// AIProvider AI 服务提供者接口
type AIProvider interface {
	GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error)
	GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error)
	Understand(ctx context.Context, req NLURequest) (*NLUResponse, error)
	DecideBehavior(ctx context.Context, req BehaviorRequest) (*BehaviorResponse, error)
	HealthCheck(ctx context.Context) error
	ProviderInfo() ProviderInfo
}

// ServiceManager AI 服务管理器
type ServiceManager struct {
	mu             sync.RWMutex
	providers      map[string]AIProvider
	defaultProvider string
	cache          *ResponseCache
	config         *config.AIConfig
}

// ResponseCache 响应缓存
type ResponseCache struct {
	mu       sync.RWMutex
	entries  map[string]cacheEntry
	maxSize  int
	ttl      time.Duration
}

type cacheEntry struct {
	response  interface{}
	timestamp time.Time
}

// NewServiceManager 创建 AI 服务管理器
func NewServiceManager(cfg *config.AIConfig) (*ServiceManager, error) {
	manager := &ServiceManager{
		providers:      make(map[string]AIProvider),
		defaultProvider: cfg.DefaultProvider,
		config:         cfg,
	}

	// 初始化缓存
	if cfg.Cache.Enabled {
		manager.cache = &ResponseCache{
			entries: make(map[string]cacheEntry),
			maxSize: cfg.Cache.MaxEntries,
			ttl:     time.Duration(cfg.Cache.TTLSeconds) * time.Second,
		}
	}

	// 初始化所有提供者
	for name, providerCfg := range cfg.Providers {
		if !providerCfg.Enabled {
			continue
		}

		provider, err := createProvider(name, providerCfg)
		if err != nil {
			log.Printf("Warning: failed to create provider %s: %v", name, err)
			continue
		}

		manager.providers[name] = provider
		log.Printf("Initialized AI provider: %s (model: %s)", name, providerCfg.Model)
	}

	// 验证默认提供者存在
	if _, ok := manager.providers[manager.defaultProvider]; !ok {
		// 尝试使用 mock 作为后备
		if _, ok := manager.providers["mock"]; ok {
			manager.defaultProvider = "mock"
			log.Println("Warning: default provider not available, using mock provider")
		} else if len(manager.providers) > 0 {
			// 使用第一个可用的提供者
			for name := range manager.providers {
				manager.defaultProvider = name
				log.Printf("Warning: using %s as default provider", name)
				break
			}
		} else {
			return nil, fmt.Errorf("no AI providers available")
		}
	}

	return manager, nil
}

// createProvider 创建 AI 提供者
func createProvider(name string, cfg config.ProviderConfig) (AIProvider, error) {
	switch name {
	case "openai_compatible":
		return NewOpenAICompatibleProvider(cfg), nil
	case "mock":
		return NewMockProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", name)
	}
}

// GetProvider 获取指定提供者
func (m *ServiceManager) GetProvider(name string) (AIProvider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if name == "" {
		name = m.defaultProvider
	}

	provider, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

	return provider, nil
}

// SetDefaultProvider 设置默认提供者
func (m *ServiceManager) SetDefaultProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.providers[name]; !ok {
		return fmt.Errorf("provider not found: %s", name)
	}

	m.defaultProvider = name
	return nil
}

// ListProviders 列出所有提供者
func (m *ServiceManager) ListProviders() map[string]ProviderInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]ProviderInfo)
	for name, provider := range m.providers {
		result[name] = provider.ProviderInfo()
	}
	return result
}

// GenerateDialogue 生成对话
func (m *ServiceManager) GenerateDialogue(ctx context.Context, req DialogueRequest) (*DialogueResponse, error) {
	// 检查缓存
	if m.cache != nil {
		cacheKey := m.buildDialogueCacheKey(req)
		if cached, ok := m.cache.Get(cacheKey); ok {
			if resp, ok := cached.(*DialogueResponse); ok {
				return resp, nil
			}
		}
	}

	provider, err := m.GetProvider("")
	if err != nil {
		return nil, err
	}

	resp, err := provider.GenerateDialogue(ctx, req)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	if m.cache != nil {
		cacheKey := m.buildDialogueCacheKey(req)
		m.cache.Set(cacheKey, resp)
	}

	return resp, nil
}

// GenerateCharacter 生成角色
func (m *ServiceManager) GenerateCharacter(ctx context.Context, req CharacterRequest) (*CharacterResponse, error) {
	provider, err := m.GetProvider("")
	if err != nil {
		return nil, err
	}

	return provider.GenerateCharacter(ctx, req)
}

// Understand 自然语言理解
func (m *ServiceManager) Understand(ctx context.Context, req NLURequest) (*NLUResponse, error) {
	provider, err := m.GetProvider("")
	if err != nil {
		return nil, err
	}

	return provider.Understand(ctx, req)
}

// DecideBehavior 行为决策
func (m *ServiceManager) DecideBehavior(ctx context.Context, req BehaviorRequest) (*BehaviorResponse, error) {
	provider, err := m.GetProvider("")
	if err != nil {
		return nil, err
	}

	return provider.DecideBehavior(ctx, req)
}

// HealthCheck 健康检查
func (m *ServiceManager) HealthCheck(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]error)
	for name, provider := range m.providers {
		results[name] = provider.HealthCheck(ctx)
	}
	return results
}

func (m *ServiceManager) buildDialogueCacheKey(req DialogueRequest) string {
	charID := ""
	if req.Character != nil {
		charID = req.Character.ID
	}
	return fmt.Sprintf("dialogue:%s:%s:%s", charID, req.CurrentMood, req.PlayerInput)
}

// Cache methods

func (c *ResponseCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if time.Since(entry.timestamp) > c.ttl {
		return nil, false
	}

	return entry.response, true
}

func (c *ResponseCache) Set(key string, response interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 简单的 LRU：如果满了就删除最旧的
	if len(c.entries) >= c.maxSize {
		var oldestKey string
		var oldestTime time.Time
		for k, v := range c.entries {
			if oldestKey == "" || v.timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.timestamp
			}
		}
		delete(c.entries, oldestKey)
	}

	c.entries[key] = cacheEntry{
		response:  response,
		timestamp: time.Now(),
	}
}

func (c *ResponseCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cacheEntry)
}
