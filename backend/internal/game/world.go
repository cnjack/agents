package game

import (
	"stardew-agent/internal/models"
)

// World manages the game world map
type World struct {
	config    *models.GameConfig
	tiles     [][]models.Tile
	buildings []models.Building
	areas     []models.MapArea
	trees     []models.Position // Tree positions
}

// NewWorld creates a new world
func NewWorld(config *models.GameConfig) *World {
	return &World{
		config: config,
	}
}

// Initialize sets up the world map
func (w *World) Initialize(config *models.GameConfig) {
	w.config = config
	w.tiles = make([][]models.Tile, config.MapHeight)

	// Initialize all tiles as grass
	for y := 0; y < config.MapHeight; y++ {
		w.tiles[y] = make([]models.Tile, config.MapWidth)
		for x := 0; x < config.MapWidth; x++ {
			w.tiles[y][x] = models.Tile{
				X:        x,
				Y:        y,
				Type:     models.TileTypeGrass,
				Occupied: false,
			}
		}
	}

	// Set up farm area (center of map)
	farmStartX := (config.MapWidth - config.FarmWidth) / 2
	farmStartY := (config.MapHeight - config.FarmHeight) / 2

	for y := farmStartY; y < farmStartY+config.FarmHeight; y++ {
		for x := farmStartX; x < farmStartX+config.FarmWidth; x++ {
			w.tiles[y][x].Type = models.TileTypeDirt
		}
	}

	// Set up buildings
	w.setupBuildings()

	// Set up water areas (pond)
	w.setupWater()

	// Set up trees
	w.setupTrees()

	// Set up areas
	w.setupAreas()
}

// setupWater adds water bodies to the map
func (w *World) setupWater() {
	// Pond at top-right (matching frontend)
	w.addWaterArea(38, 2, 8, 6)
}

// addWaterArea adds a water body
func (w *World) addWaterArea(x, y, width, height int) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px := x + dx
			py := y + dy
			if py < w.config.MapHeight && px < w.config.MapWidth {
				w.tiles[py][px].Type = models.TileTypeWater
				w.tiles[py][px].Occupied = true
			}
		}
	}
}

// setupTrees places trees around the map
func (w *World) setupTrees() {
	// Tree positions matching frontend MAP_FEATURES
	treePositions := []models.Position{
		{X: 1, Y: 1},
		{X: 1, Y: 10},
		{X: 1, Y: 20},
		{X: 1, Y: 28},
		{X: 44, Y: 1},
		{X: 44, Y: 15},
		{X: 44, Y: 25},
		{X: 44, Y: 35},
		{X: 46, Y: 8},
		{X: 46, Y: 20},
		{X: 46, Y: 32},
	}

	w.trees = treePositions

	// Mark tree tiles as occupied (trees block movement)
	for _, pos := range treePositions {
		if pos.Y < w.config.MapHeight && pos.X < w.config.MapWidth {
			w.tiles[pos.Y][pos.X].Occupied = true
		}
	}
}

// setupBuildings places buildings on the map (matching frontend MAP_FEATURES)
func (w *World) setupBuildings() {
	w.buildings = []models.Building{
		{
			ID:          "community_center",
			Name:        "Community Center",
			Type:        "community_center",
			Position:    models.Position{X: 18, Y: 4},
			Width:       12,
			Height:      8,
			Entrance:    models.Position{X: 24, Y: 12},
			Interactive: true,
		},
		{
			ID:          "pierre_shop",
			Name:        "Pierre's General Store",
			Type:        "shop",
			Position:    models.Position{X: 4, Y: 16},
			Width:       8,
			Height:      6,
			Entrance:    models.Position{X: 8, Y: 22},
			Owner:       "pierre",
			Interactive: true,
		},
		{
			ID:          "carpenter_shop",
			Name:        "Robin's Carpenter Shop",
			Type:        "carpenter",
			Position:    models.Position{X: 6, Y: 4},
			Width:       8,
			Height:      6,
			Entrance:    models.Position{X: 10, Y: 10},
			Owner:       "robin",
			Interactive: true,
		},
		{
			ID:          "fish_shop",
			Name:        "Willy's Fish Shop",
			Type:        "fish_shop",
			Position:    models.Position{X: 2, Y: 36},
			Width:       8,
			Height:      6,
			Entrance:    models.Position{X: 6, Y: 42},
			Owner:       "willy",
			Interactive: true,
		},
		{
			ID:          "saloon",
			Name:        "The Stardrop Saloon",
			Type:        "saloon",
			Position:    models.Position{X: 36, Y: 20},
			Width:       8,
			Height:      6,
			Entrance:    models.Position{X: 40, Y: 26},
			Interactive: true,
		},
		{
			ID:          "haley_house",
			Name:        "Haley's House",
			Type:        "house",
			Position:    models.Position{X: 28, Y: 6},
			Width:       6,
			Height:      5,
			Entrance:    models.Position{X: 31, Y: 11},
			Owner:       "haley",
			Interactive: false,
		},
		{
			ID:          "house_2",
			Name:        "House",
			Type:        "house",
			Position:    models.Position{X: 32, Y: 12},
			Width:       6,
			Height:      5,
			Entrance:    models.Position{X: 35, Y: 17},
			Interactive: false,
		},
	}

	// Mark building tiles as occupied
	for _, b := range w.buildings {
		for dy := 0; dy < b.Height; dy++ {
			for dx := 0; dx < b.Width; dx++ {
				x := b.Position.X + dx
				y := b.Position.Y + dy
				if y < w.config.MapHeight && x < w.config.MapWidth {
					w.tiles[y][x].Type = models.TileTypeBuilding
					w.tiles[y][x].Occupied = true
				}
			}
		}
	}
}

// setupAreas defines named areas
func (w *World) setupAreas() {
	farmStartX := (w.config.MapWidth - w.config.FarmWidth) / 2
	farmStartY := (w.config.MapHeight - w.config.FarmHeight) / 2

	w.areas = []models.MapArea{
		{
			ID:       "farm",
			Name:     "Your Farm",
			Position: models.Position{X: farmStartX, Y: farmStartY},
			Width:    w.config.FarmWidth,
			Height:   w.config.FarmHeight,
			Type:     "farm",
		},
		{
			ID:       "town",
			Name:     "Pelican Town",
			Position: models.Position{X: 14, Y: 18},
			Width:    20,
			Height:   15,
			Type:     "town",
		},
		{
			ID:       "beach",
			Name:     "The Beach",
			Position: models.Position{X: 0, Y: 42},
			Width:    12,
			Height:   6,
			Type:     "wild",
		},
	}
}

// GetTile returns the tile at the given position
func (w *World) GetTile(x, y int) *models.Tile {
	if x < 0 || x >= w.config.MapWidth || y < 0 || y >= w.config.MapHeight {
		return nil
	}
	tile := w.tiles[y][x]
	return &tile
}

// SetTile updates a tile at the given position
func (w *World) SetTile(x, y int, tile models.Tile) {
	if x < 0 || x >= w.config.MapWidth || y < 0 || y >= w.config.MapHeight {
		return
	}
	w.tiles[y][x] = tile
}

// IsWalkable checks if a position is walkable
func (w *World) IsWalkable(x, y int) bool {
	tile := w.GetTile(x, y)
	if tile == nil {
		return false
	}
	return tile.Type != models.TileTypeWater &&
		tile.Type != models.TileTypeBuilding &&
		!tile.Occupied
}

// GetBuildingAt returns the building at the given position
func (w *World) GetBuildingAt(x, y int) *models.Building {
	for i := range w.buildings {
		b := &w.buildings[i]
		if x >= b.Position.X && x < b.Position.X+b.Width &&
			y >= b.Position.Y && y < b.Position.Y+b.Height {
			return b
		}
	}
	return nil
}

// GetBuildingEntrance checks if position is a building entrance
func (w *World) GetBuildingEntrance(x, y int) *models.Building {
	for i := range w.buildings {
		b := &w.buildings[i]
		if b.Entrance.X == x && b.Entrance.Y == y {
			return b
		}
	}
	return nil
}

// GetArea returns the area containing the given position
func (w *World) GetArea(x, y int) *models.MapArea {
	for i := range w.areas {
		a := &w.areas[i]
		if x >= a.Position.X && x < a.Position.X+a.Width &&
			y >= a.Position.Y && y < a.Position.Y+a.Height {
			return a
		}
	}
	return nil
}

// GetAllTiles returns all tiles (for serialization)
func (w *World) GetAllTiles() [][]models.Tile {
	return w.tiles
}

// GetBuildings returns all buildings
func (w *World) GetBuildings() []models.Building {
	return w.buildings
}

// GetAreas returns all areas
func (w *World) GetAreas() []models.MapArea {
	return w.areas
}

// GetTrees returns tree positions
func (w *World) GetTrees() []models.Position {
	return w.trees
}
