package models

// Entity is the base interface for game entities
type Entity interface {
	GetID() string
	GetPosition() Position
	SetPosition(pos Position)
}

// BaseEntity contains common entity fields
type BaseEntity struct {
	ID       string   `json:"id"`
	Position Position `json:"position"`
}

// ToolType represents available tools
type ToolType string

const (
	ToolHoe          ToolType = "hoe"
	ToolWateringCan  ToolType = "watering_can"
	ToolAxe          ToolType = "axe"
	ToolPickaxe      ToolType = "pickaxe"
	ToolScythe       ToolType = "scythe"
	ToolFishingRod   ToolType = "fishing_rod"
)

// ToolInfo contains tool configuration
type ToolInfo struct {
	Type        ToolType
	Name        string
	EnergyCost  int
	Description string
}

// ItemType represents item categories
type ItemType string

const (
	ItemTypeSeed   ItemType = "seed"
	ItemTypeCrop   ItemType = "crop"
	ItemTypeTool   ItemType = "tool"
	ItemTypeGift   ItemType = "gift"
	ItemTypeFood   ItemType = "food"
	ItemTypeResource ItemType = "resource"
)

// ItemInfo contains item configuration
type ItemInfo struct {
	ID          string
	Name        string
	Type        ItemType
	Price       int
	SellPrice   int
	Stackable   bool
	Description string
}

// ShopItem represents an item in a shop
type ShopItem struct {
	Item     InventoryItem `json:"item"`
	Price    int           `json:"price"`
	Stock    int           `json:"stock"` // -1 = infinite
	Season   []Season      `json:"available_seasons"`
}

// Shop represents a shop with items
type Shop struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Owner    string     `json:"owner"` // NPC ID
	Items    []ShopItem `json:"items"`
	Location string     `json:"location"`
}

// Building represents a building on the map
type Building struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // house, shop, community_center
	Position    Position `json:"position"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Entrance    Position `json:"entrance"`
	Owner       string   `json:"owner,omitempty"` // NPC ID if applicable
	Interactive bool     `json:"interactive"`
}

// MapArea represents a named area on the map
type MapArea struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Position Position   `json:"position"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Type     string     `json:"type"` // farm, shop, community_center, house, wild
}

// Weather represents weather state
type Weather string

const (
	WeatherSunny  Weather = "sunny"
	WeatherRainy  Weather = "rainy"
	WeatherWindy  Weather = "windy"
	WeatherStormy Weather = "stormy"
)

// DailyInfo contains information about the current day
type DailyInfo struct {
	Weather       Weather `json:"weather"`
	LuckyDay      bool    `json:"lucky_day"`
	SpecialEvent  string  `json:"special_event,omitempty"`
	BirthdayNPC   string  `json:"birthday_npc,omitempty"`
}

// DialogueLine represents a dialogue option
type DialogueLine struct {
	Text       string   `json:"text"`
	Responses  []string `json:"responses,omitempty"`
	Condition  string   `json:"condition,omitempty"` // condition to show this dialogue
	Action     string   `json:"action,omitempty"`    // action to trigger
}

// NPCDialogue contains dialogue for an NPC
type NPCDialogue struct {
	NPCID      string         `json:"npc_id"`
	Greetings  []DialogueLine `json:"greetings"`
	Quests     []DialogueLine `json:"quest_dialogues"`
	Shop       []DialogueLine `json:"shop_dialogues"`
	Friendly   []DialogueLine `json:"friendly_dialogues"`
}

// GiftReaction represents NPC reaction to gifts
type GiftReaction struct {
	Item     string `json:"item"`
	Reaction string `json:"reaction"` // love, like, neutral, dislike, hate
	Points   int    `json:"points"`
}

// NPCConfig contains NPC configuration
type NPCConfig struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Role         string          `json:"role"`
	HomeLocation string          `json:"home_location"`
	DefaultPos   Position        `json:"default_position"`
	Schedule     []ScheduleEntry `json:"schedule"`
	Dialogue     NPCDialogue     `json:"dialogue"`
	GiftReactions []GiftReaction `json:"gift_reactions"`
	ShopID       string          `json:"shop_id,omitempty"`
}

// GameConfig contains overall game configuration
type GameConfig struct {
	MapWidth      int   `json:"map_width"`
	MapHeight     int   `json:"map_height"`
	FarmWidth     int   `json:"farm_width"`
	FarmHeight    int   `json:"farm_height"`
	TickRate      int   `json:"tick_rate"` // ms per tick
	TimeScale     int   `json:"time_scale"` // game minutes per real second
	StartingGold  int   `json:"starting_gold"`
	StartingEnergy int  `json:"starting_energy"`
	EnergyPerAction int `json:"energy_per_action"`
}
