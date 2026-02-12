package game

import (
	"fmt"
	"sync"
	"time"

	"stardew-agent/internal/models"
)

// Engine is the main game engine
type Engine struct {
	mu            sync.RWMutex
	state         *models.GameState
	world         *World
	timeSystem    *TimeSystem
	questSystem   *QuestSystem
	friendshipSystem *FriendshipSystem
	pathManager   *PathManager
	npcMovementStates map[string]*models.NPCMovementState

	ticker        *time.Ticker
	running       bool
	eventChan     chan models.GameEvent
	actionChan    chan *models.Action
	movementUpdateChan chan models.NPCMovementUpdate
}

// Config returns default game configuration
func DefaultConfig() *models.GameConfig {
	return &models.GameConfig{
		MapWidth:        48,
		MapHeight:       48,
		FarmWidth:       32,
		FarmHeight:      32,
		TickRate:        100,
		TimeScale:       10, // 10 game minutes per real second
		StartingGold:    500,
		StartingEnergy:  100,
		EnergyPerAction: 2,
	}
}

// NewEngine creates a new game engine
func NewEngine() *Engine {
	config := DefaultConfig()

	engine := &Engine{
		state:      &models.GameState{},
		world:      NewWorld(config),
		timeSystem: NewTimeSystem(),
		questSystem: NewQuestSystem(),
		friendshipSystem: NewFriendshipSystem(),
		pathManager: NewPathManager(nil),
		npcMovementStates: make(map[string]*models.NPCMovementState),
		eventChan:  make(chan models.GameEvent, 100),
		actionChan: make(chan *models.Action, 100),
		movementUpdateChan: make(chan models.NPCMovementUpdate, 100),
	}

	engine.initializeState(config)
	// Initialize path manager after world is created
	engine.pathManager = NewPathManager(engine.world)
	return engine
}

// initializeState sets up initial game state
func (e *Engine) initializeState(config *models.GameConfig) {
	// Initialize player (simplified for community-focused agent)
	e.state.Player = models.PlayerState{
		Position:   models.Position{X: 24, Y: 24},
		Direction:  models.DirectionDown,
		Gold:       config.StartingGold,
		Gifts: []models.InventoryItem{
			{ID: "flower", Name: "Flower", Type: "gift", Quantity: 5, Price: 20},
			{ID: "chocolate", Name: "Chocolate", Type: "gift", Quantity: 3, Price: 50},
		},
		Friendship: make(map[string]int),
	}

	// Initialize time
	e.state.Time = models.TimeState{
		Day:     1,
		Hour:    6,
		Minute:  0,
		Season:  models.SeasonSpring,
		Year:    1,
		Paused:  false,
	}

	// Initialize map dimensions
	e.state.MapWidth = config.MapWidth
	e.state.MapHeight = config.MapHeight

	// Initialize NPCs
	e.state.NPCs = e.createNPCs()

	// Initialize quests
	e.state.Quests = e.questSystem.GetAvailableQuests()

	e.world.Initialize(config)
}

// createNPCs creates the initial NPC roster
func (e *Engine) createNPCs() []models.NPCState {
	return []models.NPCState{
		{
			ID:           "lewis",
			Name:         "Lewis",
			Position:     models.Position{X: 40, Y: 10},
			Direction:    models.DirectionDown,
			Friendship:   0,
			MaxFriendship: 250,
			Location:     "community_center",
			Dialogue: []string{
				"Welcome to Stardew Valley! I'm Lewis, the mayor.",
				"There's plenty to do around here. Check the community center for tasks.",
				"The farm has been in disrepair for years. It's good to see someone working it.",
			},
			Schedule: []models.ScheduleEntry{
				{Hour: 6, Minute: 0, Location: "home", Position: models.Position{X: 40, Y: 10}},
				{Hour: 9, Minute: 0, Location: "community_center", Position: models.Position{X: 35, Y: 15}},
				{Hour: 17, Minute: 0, Location: "town_square", Position: models.Position{X: 30, Y: 20}},
				{Hour: 21, Minute: 0, Location: "home", Position: models.Position{X: 40, Y: 10}},
			},
		},
		{
			ID:           "pierre",
			Name:         "Pierre",
			Position:     models.Position{X: 10, Y: 30},
			Direction:    models.DirectionDown,
			Friendship:   0,
			MaxFriendship: 250,
			Location:     "shop",
			Dialogue: []string{
				"Welcome to Pierre's General Store!",
				"I've got seeds, tools, and all sorts of supplies.",
				"Bring me your crops and I'll give you a fair price.",
			},
			Schedule: []models.ScheduleEntry{
				{Hour: 6, Minute: 0, Location: "home", Position: models.Position{X: 12, Y: 32}},
				{Hour: 9, Minute: 0, Location: "shop", Position: models.Position{X: 10, Y: 30}},
				{Hour: 18, Minute: 0, Location: "home", Position: models.Position{X: 12, Y: 32}},
			},
		},
		{
			ID:           "robin",
			Name:         "Robin",
			Position:     models.Position{X: 15, Y: 35},
			Direction:    models.DirectionDown,
			Friendship:   0,
			MaxFriendship: 250,
			Location:     "carpenter_shop",
			Dialogue: []string{
				"Hey there! Need any building work done?",
				"I can upgrade your farmhouse or build new structures.",
				"Quality craftsmanship is my specialty!",
			},
			Schedule: []models.ScheduleEntry{
				{Hour: 6, Minute: 0, Location: "home", Position: models.Position{X: 15, Y: 35}},
				{Hour: 9, Minute: 0, Location: "carpenter_shop", Position: models.Position{X: 16, Y: 36}},
				{Hour: 17, Minute: 0, Location: "home", Position: models.Position{X: 15, Y: 35}},
			},
		},
		{
			ID:           "haley",
			Name:         "Haley",
			Position:     models.Position{X: 25, Y: 15},
			Direction:    models.DirectionDown,
			Friendship:   0,
			MaxFriendship: 250,
			Location:     "home",
			Dialogue: []string{
				"Oh, a new farmer? How... rustic.",
				"I prefer the finer things in life. Photography, fashion...",
				"Maybe there's more to this place than I thought.",
			},
			Schedule: []models.ScheduleEntry{
				{Hour: 6, Minute: 0, Location: "home", Position: models.Position{X: 25, Y: 15}},
				{Hour: 12, Minute: 0, Location: "beach", Position: models.Position{X: 5, Y: 45}},
				{Hour: 18, Minute: 0, Location: "saloon", Position: models.Position{X: 32, Y: 25}},
				{Hour: 22, Minute: 0, Location: "home", Position: models.Position{X: 25, Y: 15}},
			},
		},
		{
			ID:           "willy",
			Name:         "Willy",
			Position:     models.Position{X: 5, Y: 45},
			Direction:    models.DirectionUp,
			Friendship:   0,
			MaxFriendship: 250,
			Location:     "fish_shop",
			Dialogue: []string{
				"Ahoy there! Fine day for fishing!",
				"The ocean's full of surprises if you've got the patience.",
				"I've got rods, bait, and all sorts of fishing gear.",
			},
			Schedule: []models.ScheduleEntry{
				{Hour: 6, Minute: 0, Location: "fish_shop", Position: models.Position{X: 5, Y: 45}},
				{Hour: 10, Minute: 0, Location: "dock", Position: models.Position{X: 3, Y: 47}},
				{Hour: 17, Minute: 0, Location: "saloon", Position: models.Position{X: 32, Y: 25}},
				{Hour: 22, Minute: 0, Location: "fish_shop", Position: models.Position{X: 5, Y: 45}},
			},
		},
	}
}

// Start begins the game loop
func (e *Engine) Start() {
	e.mu.Lock()
	e.running = true
	e.ticker = time.NewTicker(time.Duration(DefaultConfig().TickRate) * time.Millisecond)
	e.mu.Unlock()

	go e.gameLoop()
}

// Stop stops the game loop
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.running = false
	if e.ticker != nil {
		e.ticker.Stop()
	}
}

// gameLoop runs the main game loop
func (e *Engine) gameLoop() {
	for {
		e.mu.RLock()
		running := e.running
		e.mu.RUnlock()

		if !running {
			return
		}

		select {
		case <-e.ticker.C:
			e.tick()
		case action := <-e.actionChan:
			e.processAction(action)
		}
	}
}

// tick processes one game tick
func (e *Engine) tick() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state.Time.Paused {
		return
	}

	// Update time
	e.timeSystem.AdvanceTime(&e.state.Time)

	// Update NPCs based on schedule
	e.updateNPCs()

	// Check for day end
	if e.state.Time.Hour >= 2 && e.state.Time.Hour < 6 {
		e.state.DayEnded = true
	}
}

// updateNPCs updates NPC positions based on schedule with pathfinding
func (e *Engine) updateNPCs() {
	for i := range e.state.NPCs {
		npc := &e.state.NPCs[i]
		if len(npc.Schedule) == 0 {
			continue
		}

		// Check if NPC needs to move to a new scheduled location
		e.updateNPCScheduleMovement(npc)
	}

	// Update NPC movement along paths
	e.updateNPCMovement()
}

// updateNPCScheduleMovement checks if NPC needs to start moving to scheduled location
func (e *Engine) updateNPCScheduleMovement(npc *models.NPCState) {
	// Find current schedule entry
	var targetEntry *models.ScheduleEntry
	for _, entry := range npc.Schedule {
		entryMinutes := entry.Hour*60 + entry.Minute
		currentMinutes := e.state.Time.Hour*60 + e.state.Time.Minute

		// Find the most recent schedule entry (last one before or at current time)
		if entryMinutes <= currentMinutes {
			targetEntry = &entry
		}
	}

	if targetEntry == nil {
		return
	}

	// Check if we already have a path for this NPC
	movState := e.npcMovementStates[npc.ID]
	if movState != nil && movState.IsMoving {
		return // Already moving
	}

	// Check if NPC is not at target position
	if npc.Position.X != targetEntry.Position.X || npc.Position.Y != targetEntry.Position.Y {
		// Need to move - find path
		path := e.pathManager.FindAndSetPath(npc.ID, npc.Position, targetEntry.Position)

		if path != nil && len(path.Positions) > 0 {
			// Initialize movement state
			e.npcMovementStates[npc.ID] = &models.NPCMovementState{
				NPCID:       npc.ID,
				IsMoving:     true,
				CurrentPath:  path.Positions,
				TargetPos:    targetEntry.Position,
				MoveSpeed:    0.1, // Move speed (tiles per tick)
				StartedAt:    e.state.Time,
			}
		}
	}

	// Update location name
	npc.Location = targetEntry.Location
}

// updateNPCMovement updates NPC positions along their paths
func (e *Engine) updateNPCMovement() {
	for npcID, movState := range e.npcMovementStates {
		if !movState.IsMoving {
			continue
		}

		// Find NPC in state
		var npc *models.NPCState
		for i := range e.state.NPCs {
			if e.state.NPCs[i].ID == npcID {
				npc = &e.state.NPCs[i]
				break
			}
		}
		if npc == nil {
			delete(e.npcMovementStates, npcID)
			continue
		}

		// Get next position from path
		nextPos := e.pathManager.GetNextPosition(npcID, npc.Position)
		if nextPos == nil {
			// Path complete or no path
			movState.IsMoving = false
			e.npcMovementStates[npcID].IsMoving = false

			// Send update notification
			select {
			case e.movementUpdateChan <- models.NPCMovementUpdate{
				NPCID:    npcID,
				OldPos:    npc.Position,
				NewPos:    npc.Position,
				Direction:  npc.Direction,
				IsMoving:  false,
				PathEnded: true,
			}:
			default:
			}
			continue
		}

		// Calculate direction
		direction := e.calculateDirection(npc.Position, *nextPos)

		// Update NPC position
		oldPos := npc.Position
		npc.Position = *nextPos
		npc.Direction = direction

		// Send movement update
		select {
		case e.movementUpdateChan <- models.NPCMovementUpdate{
			NPCID:   npcID,
			OldPos:   oldPos,
			NewPos:   *nextPos,
			Direction: direction,
			IsMoving: true,
		}:
		default:
		}
	}
}

// calculateDirection determines the direction from pos1 to pos2
func (e *Engine) calculateDirection(pos1, pos2 models.Position) models.Direction {
	dx := pos2.X - pos1.X
	dy := pos2.Y - pos1.Y

	if dx == 0 && dy < 0 {
		return models.DirectionUp
	}
	if dx == 0 && dy > 0 {
		return models.DirectionDown
	}
	if dy == 0 && dx < 0 {
		return models.DirectionLeft
	}
	if dy == 0 && dx > 0 {
		return models.DirectionRight
	}

	// Diagonal - prefer primary direction
	if abs(dx) > abs(dy) {
		if dx > 0 {
			return models.DirectionRight
		}
		return models.DirectionLeft
	}
	if dy > 0 {
		return models.DirectionDown
	}
	return models.DirectionUp
}

// processAction processes an action from the action channel
func (e *Engine) processAction(action *models.Action) {
	// Actions are processed through ExecuteAction
}

// GetState returns the current game state
func (e *Engine) GetState() *models.GameState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Return a copy
	state := *e.state
	return &state
}

// Reset resets the game to initial state
func (e *Engine) Reset() *models.GameState {
	e.mu.Lock()
	defer e.mu.Unlock()

	config := DefaultConfig()
	e.initializeState(config)
	return e.state
}

// GetObservation returns what an agent can observe
func (e *Engine) GetObservation() *models.Observation {
	e.mu.RLock()
	defer e.mu.RUnlock()

	obs := &models.Observation{
		Player:       e.state.Player,
		Time:         e.state.Time,
		NearbyTiles:  e.getNearbyTiles(5),
		NearbyNPCs:   e.getNearbyNPCs(10),
		ActiveQuests: e.getActiveQuests(),
		Messages:     []string{},
	}

	return obs
}

// getNearbyTiles returns tiles near the player
func (e *Engine) getNearbyTiles(radius int) []models.Tile {
	playerPos := e.state.Player.Position
	var tiles []models.Tile

	for y := playerPos.Y - radius; y <= playerPos.Y+radius; y++ {
		for x := playerPos.X - radius; x <= playerPos.X+radius; x++ {
			if tile := e.world.GetTile(x, y); tile != nil {
				tiles = append(tiles, *tile)
			}
		}
	}

	return tiles
}

// getNearbyNPCs returns NPCs near the player
func (e *Engine) getNearbyNPCs(radius int) []models.NPCState {
	playerPos := e.state.Player.Position
	var npcs []models.NPCState

	for _, npc := range e.state.NPCs {
		dx := npc.Position.X - playerPos.X
		dy := npc.Position.Y - playerPos.Y
		if dx*dx+dy*dy <= radius*radius {
			npcs = append(npcs, npc)
		}
	}

	return npcs
}

// getActiveQuests returns active quests
func (e *Engine) getActiveQuests() []models.QuestData {
	var quests []models.QuestData
	for _, q := range e.state.Quests {
		if q.Status == models.QuestStatusActive {
			quests = append(quests, q)
		}
	}
	return quests
}

// Pause pauses the game
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.state.Time.Paused = true
}

// Resume resumes the game
func (e *Engine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.state.Time.Paused = false
}

// World returns the world manager
func (e *Engine) World() *World {
	return e.world
}

// MovePlayer moves the player in a direction (directly modifies engine state)
func (e *Engine) MovePlayer(direction models.Direction) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	player := &e.state.Player
	player.Direction = direction

	newPos := player.Position
	switch direction {
	case models.DirectionUp:
		newPos.Y--
	case models.DirectionDown:
		newPos.Y++
	case models.DirectionLeft:
		newPos.X--
	case models.DirectionRight:
		newPos.X++
	}

	if !e.world.IsWalkable(newPos.X, newPos.Y) {
		return fmt.Errorf("cannot move there")
	}

	player.Position = newPos
	return nil
}

// EndDay ends the current day and starts a new one
func (e *Engine) EndDay() {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Reset daily states
	e.state.DayEnded = false

	// Advance to next day
	e.timeSystem.AdvanceDay(&e.state.Time)
}

// MovementUpdateChannel returns the channel for NPC movement updates
func (e *Engine) MovementUpdateChannel() chan models.NPCMovementUpdate {
	return e.movementUpdateChan
}

// FindPath finds a path from start to end for an NPC
func (e *Engine) FindPath(npcID string, start, goal models.Position) *Path {
	e.mu.Lock()
	defer e.mu.Unlock()

	path := e.pathManager.FindAndSetPath(npcID, start, goal)
	return path
}

// GetNPCMovementState returns the movement state for an NPC
func (e *Engine) GetNPCMovementState(npcID string) *models.NPCMovementState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.npcMovementStates[npcID]
}

// SetNPCPath manually sets a path for an NPC
func (e *Engine) SetNPCPath(npcID string, path *Path) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.pathManager.SetPath(npcID, path)

	if path != nil && len(path.Positions) > 0 {
		e.npcMovementStates[npcID] = &models.NPCMovementState{
			NPCID:       npcID,
			IsMoving:     true,
			CurrentPath:  path.Positions,
			TargetPos:    path.Positions[len(path.Positions)-1],
			MoveSpeed:    0.1,
			StartedAt:    e.state.Time,
		}
	}
}

// PathManager returns the path manager
func (e *Engine) PathManager() *PathManager {
	return e.pathManager
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
