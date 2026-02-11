package models

// Position represents a 2D coordinate
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Direction represents facing direction
type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

// Season represents game season
type Season string

const (
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonFall   Season = "fall"
	SeasonWinter Season = "winter"
)

// TileType represents the type of a map tile
type TileType string

const (
	TileTypeGrass   TileType = "grass"
	TileTypeDirt    TileType = "dirt"
	TileTypeTilled  TileType = "tilled"
	TileTypeWatered TileType = "watered"
	TileTypePath    TileType = "path"
	TileTypeWater   TileType = "water"
	TileTypeBuilding TileType = "building"
)

// CropStage represents growth stage of a crop
type CropStage string

const (
	CropStageSeed   CropStage = "seed"
	CropStageSprout CropStage = "sprout"
	CropStageGrowing CropStage = "growing"
	CropStageMature CropStage = "mature"
)

// CropType represents different crop varieties
type CropType string

const (
	CropPotato  CropType = "potato"
	CropTomato  CropType = "tomato"
	CropCorn    CropType = "corn"
	CropPumpkin CropType = "pumpkin"
	CropStrawberry CropType = "strawberry"
	CropCauliflower CropType = "cauliflower"
)

// CropInfo contains crop configuration
type CropInfo struct {
	Type          CropType
	Name          string
	GrowthDays    int
	Seasons       []Season
	BasePrice     int
	SeedPrice     int
	Stages        []int // days per stage
}

// CropData represents a planted crop
type CropData struct {
	Type         CropType  `json:"type"`
	Stage        CropStage `json:"stage"`
	DaysGrown    int       `json:"days_grown"`
	WateredToday bool      `json:"watered_today"`
}

// Tile represents a single map tile
type Tile struct {
	X       int        `json:"x"`
	Y       int        `json:"y"`
	Type    TileType   `json:"type"`
	Crop    *CropData  `json:"crop,omitempty"`
	Occupied bool      `json:"occupied"`
}

// InventoryItem represents an item in inventory
type InventoryItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"` // seed, crop, tool, gift
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

// PlayerState represents player state
type PlayerState struct {
	Position   Position        `json:"position"`
	Direction  Direction       `json:"direction"`
	Gold       int             `json:"gold"`
	Gifts      []InventoryItem `json:"gifts"` // 简化为礼物系统
	Friendship map[string]int  `json:"friendship"` // NPC ID -> friendship points
}

// NPCState represents NPC state
type NPCState struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Position   Position `json:"position"`
	Direction  Direction `json:"direction"`
	Friendship int      `json:"friendship"`
	MaxFriendship int   `json:"max_friendship"`
	Location   string   `json:"location"` // current building/area
	Dialogue   []string `json:"dialogue"`
	Schedule   []ScheduleEntry `json:"schedule,omitempty"`
}

// ScheduleEntry represents NPC schedule
type ScheduleEntry struct {
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	Location string `json:"location"`
	Position Position `json:"position"`
}

// QuestStatus represents quest state
type QuestStatus string

const (
	QuestStatusAvailable QuestStatus = "available"
	QuestStatusActive    QuestStatus = "active"
	QuestStatusCompleted QuestStatus = "completed"
	QuestStatusFailed    QuestStatus = "failed"
)

// QuestData represents a quest
type QuestData struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Objectives  []Objective `json:"objectives"`
	Rewards     Reward      `json:"rewards"`
	Status      QuestStatus `json:"status"`
	Giver       string      `json:"giver"` // NPC ID
	Deadline    int         `json:"deadline"` // days to complete, 0 = no deadline
}

// Objective represents a quest objective
type Objective struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type"` // harvest, deliver, talk, collect
	Target      string `json:"target"` // item ID, NPC ID, etc.
	Count       int    `json:"count"`
	Progress    int    `json:"progress"`
	Completed   bool   `json:"completed"`
}

// Reward represents quest reward
type Reward struct {
	Gold  int              `json:"gold"`
	Items []InventoryItem  `json:"items"`
}

// TimeState represents game time
type TimeState struct {
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minute  int    `json:"minute"`
	Season  Season `json:"season"`
	Year    int    `json:"year"`
	Paused  bool   `json:"paused"`
}

// FarmState represents the farm area
type FarmState struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Tiles  [][]Tile `json:"tiles"`
}

// GameState represents the complete game state
type GameState struct {
	Player    PlayerState   `json:"player"`
	Time      TimeState     `json:"time"`
	NPCs      []NPCState    `json:"npcs"`
	Quests    []QuestData   `json:"quests"`
	MapWidth  int           `json:"map_width"`
	MapHeight int           `json:"map_height"`
	GameOver  bool          `json:"game_over"`
	DayEnded  bool          `json:"day_ended"`
}

// Observation is what the Agent can see
type Observation struct {
	Player       PlayerState   `json:"player"`
	Time         TimeState     `json:"time"`
	NearbyTiles  []Tile        `json:"nearby_tiles"`
	NearbyNPCs   []NPCState    `json:"nearby_npcs"`
	ActiveQuests []QuestData   `json:"active_quests"`
	Messages     []string      `json:"messages"`
}
