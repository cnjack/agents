package game

import (
	"container/heap"
	"fmt"

	"stardew-agent/internal/models"
)

// Path represents a sequence of positions for NPC movement
type Path struct {
	Positions []models.Position `json:"positions"`
	CurrentIndex int           `json:"current_index"`
	Complete   bool           `json:"complete"`
}

// PathNode represents a node in the pathfinding graph
type PathNode struct {
	Position models.Position
	G       float64 // Cost from start
	H       float64 // Heuristic to goal
	F       float64 // G + H
	Parent  *PathNode
	Index   int // For heap interface
}

// PathPriorityQueue implements heap.Interface for pathfinding
type PathPriorityQueue []*PathNode

func (pq PathPriorityQueue) Len() int { return len(pq) }

func (pq PathPriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PathPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PathPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*PathNode)
	node.Index = n
	*pq = append(*pq, node)
}

func (pq *PathPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // avoid memory leak
	node.Index = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

// PathfindingService handles A* pathfinding
type PathfindingService struct {
	world *World
}

// NewPathfindingService creates a new pathfinding service
func NewPathfindingService(world *World) *PathfindingService {
	return &PathfindingService{
		world: world,
	}
}

// FindPath finds a path from start to goal using A* algorithm
func (ps *PathfindingService) FindPath(start, goal models.Position) *Path {
	// Check if start or goal is invalid
	if !ps.isValidPosition(start) {
		return &Path{Positions: []models.Position{}, Complete: true}
	}
	if !ps.isValidPosition(goal) {
		return &Path{Positions: []models.Position{}, Complete: true}
	}

	// If start equals goal, return empty path
	if start.X == goal.X && start.Y == goal.Y {
		return &Path{Positions: []models.Position{}, Complete: true}
	}

	// Check if goal is reachable (walkable)
	if !ps.world.IsWalkable(goal.X, goal.Y) {
		// Try to find nearest walkable tile to goal
		adjustedGoal := ps.findNearestWalkable(goal)
		if adjustedGoal == nil {
			return &Path{Positions: []models.Position{}, Complete: true}
		}
		goal = *adjustedGoal
	}

	// Initialize open and closed sets
	openSet := make(PathPriorityQueue, 0)
	heap.Init(&openSet)

	closedSet := make(map[string]bool)
	cameFrom := make(map[string]models.Position)
	gScore := make(map[string]float64)

	// Add start node
	startNode := &PathNode{
		Position: start,
		G:       0,
		H:       ps.heuristic(start, goal),
		F:       ps.heuristic(start, goal),
	}
	heap.Push(&openSet, startNode)
	gScore[posToKey(start)] = 0

	// Direction vectors for 8-way movement (including diagonals)
	directions := []models.Position{
		{X: 0, Y: -1},  // N
		{X: 1, Y: -1},  // NE
		{X: 1, Y: 0},   // E
		{X: 1, Y: 1},   // SE
		{X: 0, Y: 1},   // S
		{X: -1, Y: 1},  // SW
		{X: -1, Y: 0},  // W
		{X: -1, Y: -1}, // NW
	}

	maxIterations := 5000 // Prevent infinite loops
	iterations := 0

	for openSet.Len() > 0 && iterations < maxIterations {
		iterations++

		// Get node with lowest F score
		current := heap.Pop(&openSet).(*PathNode)
		currentKey := posToKey(current.Position)

		// Check if we reached the goal
		if current.Position.X == goal.X && current.Position.Y == goal.Y {
			return ps.reconstructPath(cameFrom, current.Position, start)
		}

		closedSet[currentKey] = true

		// Explore neighbors
		for _, dir := range directions {
			neighbor := models.Position{
				X: current.Position.X + dir.X,
				Y: current.Position.Y + dir.Y,
			}

			neighborKey := posToKey(neighbor)

			// Skip if invalid or already closed
			if !ps.isValidPosition(neighbor) {
				continue
			}
			if !ps.world.IsWalkable(neighbor.X, neighbor.Y) {
				continue
			}
			if closedSet[neighborKey] {
				continue
			}

			// Calculate movement cost (diagonal costs more)
			moveCost := 1.0
			if dir.X != 0 && dir.Y != 0 {
				moveCost = 1.414 // sqrt(2)
			}

			tentativeG := gScore[currentKey] + moveCost

			// Check if this is a better path to neighbor
			if existingG, exists := gScore[neighborKey]; !exists || tentativeG < existingG {
				cameFrom[neighborKey] = current.Position
				gScore[neighborKey] = tentativeG

				h := ps.heuristic(neighbor, goal)
				f := tentativeG + h

				// Check if neighbor is in open set
				found := false
				for _, node := range openSet {
					if node.Position.X == neighbor.X && node.Position.Y == neighbor.Y {
						node.G = tentativeG
						node.H = h
						node.F = f
						node.Parent = current
						heap.Init(&openSet) // Re-heap after update
						found = true
						break
					}
				}

				if !found {
					neighborNode := &PathNode{
						Position: neighbor,
						G:       tentativeG,
						H:       h,
						F:       f,
						Parent:  current,
					}
					heap.Push(&openSet, neighborNode)
				}
			}
		}
	}

	// No path found
	return &Path{Positions: []models.Position{}, Complete: true}
}

// heuristic calculates the estimated cost from a to b using Manhattan distance
// for 4-directional or Euclidean for 8-directional movement
func (ps *PathfindingService) heuristic(a, b models.Position) float64 {
	dx := abs(a.X - b.X)
	dy := abs(a.Y - b.Y)
	// Octile distance for 8-way movement
	if dx > dy {
		return float64(dx-dy) + 1.414*float64(dy)
	}
	return float64(dy-dx) + 1.414*float64(dx)
}

// isValidPosition checks if a position is within map bounds
func (ps *PathfindingService) isValidPosition(pos models.Position) bool {
	if pos.X < 0 || pos.X >= ps.world.config.MapWidth {
		return false
	}
	if pos.Y < 0 || pos.Y >= ps.world.config.MapHeight {
		return false
	}
	return true
}

// findNearestWalkable finds the nearest walkable tile to the given position
func (ps *PathfindingService) findNearestWalkable(pos models.Position) *models.Position {
	searchRadius := 10
	for radius := 0; radius <= searchRadius; radius++ {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				if dx == 0 && dy == 0 && radius == 0 {
					continue
				}
				// Only check perimeter at each radius
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}
				checkX := pos.X + dx
				checkY := pos.Y + dy
				if ps.isValidPosition(models.Position{X: checkX, Y: checkY}) {
					if ps.world.IsWalkable(checkX, checkY) {
						return &models.Position{X: checkX, Y: checkY}
					}
				}
			}
		}
	}
	return nil
}

// reconstructPath builds the path from the cameFrom map
func (ps *PathfindingService) reconstructPath(cameFrom map[string]models.Position, current, start models.Position) *Path {
	path := make([]models.Position, 0)
	path = append(path, current)

	currentKey := posToKey(current)
	for {
		if prev, exists := cameFrom[currentKey]; exists {
			path = append([]models.Position{prev}, path...)
			currentKey = posToKey(prev)
			if prev.X == start.X && prev.Y == start.Y {
				break
			}
		} else {
			break
		}
	}

	// Remove start position from path (we're already there)
	if len(path) > 0 && path[0].X == start.X && path[0].Y == start.Y {
		path = path[1:]
	}

	return &Path{
		Positions:    path,
		CurrentIndex: 0,
		Complete:     false,
	}
}

// GetNextPosition returns the next position in the path
func (p *Path) GetNextPosition() *models.Position {
	if p.CurrentIndex >= len(p.Positions) {
		p.Complete = true
		return nil
	}

	pos := p.Positions[p.CurrentIndex]
	p.CurrentIndex++
	return &pos
}

// IsComplete returns true if the path is complete
func (p *Path) IsComplete() bool {
	if p.Complete {
		return true
	}
	if p.CurrentIndex >= len(p.Positions) {
		p.Complete = true
		return true
	}
	return false
}

// RemainingSteps returns the number of steps remaining in the path
func (p *Path) RemainingSteps() int {
	if len(p.Positions) == 0 {
		return 0
	}
	return len(p.Positions) - p.CurrentIndex
}

// SimplifyPath removes redundant waypoints from the path
// (keeps only direction changes)
func (p *Path) SimplifyPath() *Path {
	if len(p.Positions) <= 2 {
		return p
	}

	simplified := []models.Position{p.Positions[0]}

	for i := 1; i < len(p.Positions)-1; i++ {
		prev := simplified[len(simplified)-1]
		curr := p.Positions[i]
		next := p.Positions[i+1]

		// Check if direction changes
		dir1 := getDirection(prev, curr)
		dir2 := getDirection(curr, next)

		if dir1 != dir2 {
			simplified = append(simplified, curr)
		}
	}

	simplified = append(simplified, p.Positions[len(p.Positions)-1])

	p.Positions = simplified
	return p
}

// getDirection returns the direction from a to b
func getDirection(a, b models.Position) string {
	dx := b.X - a.X
	dy := b.Y - a.Y

	if dx > 0 && dy < 0 {
		return "ne"
	}
	if dx > 0 && dy == 0 {
		return "e"
	}
	if dx > 0 && dy > 0 {
		return "se"
	}
	if dx == 0 && dy > 0 {
		return "s"
	}
	if dx < 0 && dy > 0 {
		return "sw"
	}
	if dx < 0 && dy == 0 {
		return "w"
	}
	if dx < 0 && dy < 0 {
		return "nw"
	}
	return "n"
}

// posToKey converts a position to a string key for maps
func posToKey(pos models.Position) string {
	return fmt.Sprintf("%d,%d", pos.X, pos.Y)
}

// PathManager manages active paths for all NPCs
type PathManager struct {
	pathfinding *PathfindingService
	paths      map[string]*Path // NPC ID -> Path
}

// NewPathManager creates a new path manager
func NewPathManager(world *World) *PathManager {
	return &PathManager{
		pathfinding: NewPathfindingService(world),
		paths:      make(map[string]*Path),
	}
}

// SetPath sets a path for an NPC
func (pm *PathManager) SetPath(npcID string, path *Path) {
	pm.paths[npcID] = path
}

// GetPath gets the current path for an NPC
func (pm *PathManager) GetPath(npcID string) *Path {
	return pm.paths[npcID]
}

// FindAndSetPath finds a path and sets it for the NPC
func (pm *PathManager) FindAndSetPath(npcID string, start, goal models.Position) *Path {
	path := pm.pathfinding.FindPath(start, goal)
	path.SimplifyPath()
	pm.paths[npcID] = path
	return path
}

// ClearPath removes the path for an NPC
func (pm *PathManager) ClearPath(npcID string) {
	delete(pm.paths, npcID)
}

// GetNextPosition returns the next position for an NPC along their path
func (pm *PathManager) GetNextPosition(npcID string, current models.Position) *models.Position {
	path := pm.paths[npcID]
	if path == nil || path.IsComplete() {
		pm.ClearPath(npcID)
		return nil
	}

	return path.GetNextPosition()
}

// UpdateWorld updates the world reference for pathfinding
func (pm *PathManager) UpdateWorld(world *World) {
	pm.pathfinding = NewPathfindingService(world)
}
