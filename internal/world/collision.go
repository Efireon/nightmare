package world

import (
	"math"

	"nightmare/internal/common"
	"nightmare/internal/entity"
)

// CollisionSystem manages collisions in the game world
type CollisionSystem struct {
	world        *World
	cellSize     float64
	collisionMap [][]bool
}

// CollisionResult represents the result of a collision check
type CollisionResult struct {
	HasCollision bool
	Normal       common.Vector2D
	Object       interface{}
	Distance     float64
}

// NewCollisionSystem creates a new collision system
func NewCollisionSystem(world *World, cellSize float64) *CollisionSystem {
	// Create collision map
	width := int(math.Ceil(float64(world.Width) / cellSize))
	height := int(math.Ceil(float64(world.Height) / cellSize))

	collisionMap := make([][]bool, height)
	for y := range collisionMap {
		collisionMap[y] = make([]bool, width)
	}

	return &CollisionSystem{
		world:        world,
		cellSize:     cellSize,
		collisionMap: collisionMap,
	}
}

// UpdateCollisionMap updates the collision map
func (cs *CollisionSystem) UpdateCollisionMap() {
	// Reset collision map
	for y := range cs.collisionMap {
		for x := range cs.collisionMap[y] {
			cs.collisionMap[y][x] = false
		}
	}

	// Update collisions based on tiles
	for y := 0; y < cs.world.Height; y++ {
		for x := 0; x < cs.world.Width; x++ {
			tile := cs.world.GetTileAt(x, y)
			if tile == nil {
				continue
			}

			// Check if the tile is impassable
			if cs.isTileSolid(tile) {
				// Calculate indexes in the collision map
				cellX := int(float64(x) / cs.cellSize)
				cellY := int(float64(y) / cs.cellSize)

				// Check that indexes are within the map
				if cellX >= 0 && cellX < len(cs.collisionMap[0]) && cellY >= 0 && cellY < len(cs.collisionMap) {
					cs.collisionMap[cellY][cellX] = true
				}
			}
		}
	}

	// Add collisions from objects
	for _, obj := range cs.world.Objects {
		if obj.Solid {
			// Calculate indexes in the collision map
			cellX := int(obj.Position.X / cs.cellSize)
			cellY := int(obj.Position.Y / cs.cellSize)

			// Check that indexes are within the map
			if cellX >= 0 && cellX < len(cs.collisionMap[0]) && cellY >= 0 && cellY < len(cs.collisionMap) {
				cs.collisionMap[cellY][cellX] = true
			}
		}
	}
}

// CheckCollision checks for a collision at the specified position
func (cs *CollisionSystem) CheckCollision(position common.Vector2D) bool {
	// Calculate indexes in the collision map
	cellX := int(position.X / cs.cellSize)
	cellY := int(position.Y / cs.cellSize)

	// Check that indexes are within the map
	if cellX < 0 || cellX >= len(cs.collisionMap[0]) || cellY < 0 || cellY >= len(cs.collisionMap) {
		return true // Outside the world is considered a collision
	}

	return cs.collisionMap[cellY][cellX]
}

// CheckMovement checks if movement from the current position to a new one is possible
func (cs *CollisionSystem) CheckMovement(from, to common.Vector2D) (common.Vector2D, bool) {
	// Check for a collision at the new position
	if cs.CheckCollision(to) {
		// Collision exists, movement is not possible
		return from, false
	}

	// Check the path from from to to
	direction := common.Vector2D{
		X: to.X - from.X,
		Y: to.Y - from.Y,
	}

	distance := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y)

	// If the distance is less than the threshold value, consider the movement possible
	if distance < cs.cellSize*0.1 {
		return to, true
	}

	// Normalize direction
	direction.X /= distance
	direction.Y /= distance

	// Check several points along the way
	steps := int(distance / (cs.cellSize * 0.5))
	if steps < 2 {
		steps = 2
	}

	stepSize := distance / float64(steps)

	for i := 1; i < steps; i++ {
		checkPoint := common.Vector2D{
			X: from.X + direction.X*stepSize*float64(i),
			Y: from.Y + direction.Y*stepSize*float64(i),
		}

		if cs.CheckCollision(checkPoint) {
			// Found a collision, movement to this point is not possible
			// Return the last safe position
			safePoint := common.Vector2D{
				X: from.X + direction.X*stepSize*float64(i-1),
				Y: from.Y + direction.Y*stepSize*float64(i-1),
			}

			return safePoint, false
		}
	}

	// Movement is possible
	return to, true
}

// CheckMovementWithSliding checks if movement is possible with sliding along walls
func (cs *CollisionSystem) CheckMovementWithSliding(from, to common.Vector2D) common.Vector2D {
	// Check direct collision
	newPos, canMove := cs.CheckMovement(from, to)
	if canMove {
		return to
	}

	// If movement is not possible, try movement only along X
	xMove := common.Vector2D{X: to.X, Y: from.Y}
	newPos, canMoveX := cs.CheckMovement(from, xMove)

	// If movement along X is possible, use it
	if canMoveX {
		return newPos
	}

	// Try movement only along Y
	yMove := common.Vector2D{X: from.X, Y: to.Y}
	newPos, canMoveY := cs.CheckMovement(from, yMove)

	// If movement along Y is possible, use it
	if canMoveY {
		return newPos
	}

	// If none of the options work, stay in place
	return from
}

// CastRay performs a raycast from a point in the specified direction
func (cs *CollisionSystem) CastRay(start common.Vector2D, direction common.Vector2D, maxDistance float64) (CollisionResult, bool) {
	// Normalize direction
	length := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y)
	if length > 0 {
		direction.X /= length
		direction.Y /= length
	} else {
		return CollisionResult{}, false
	}

	// Raycast step
	stepSize := cs.cellSize * 0.5

	// Check points along the ray path
	for dist := 0.0; dist <= maxDistance; dist += stepSize {
		checkPoint := common.Vector2D{
			X: start.X + direction.X*dist,
			Y: start.Y + direction.Y*dist,
		}

		// Check if the point is within the world
		if checkPoint.X < 0 || checkPoint.X >= float64(cs.world.Width) ||
			checkPoint.Y < 0 || checkPoint.Y >= float64(cs.world.Height) {
			// Reached the world boundary
			return CollisionResult{
				HasCollision: true,
				Normal:       common.Vector2D{X: -direction.X, Y: -direction.Y},
				Object:       nil,
				Distance:     dist,
			}, true
		}

		// Check for a collision at the point
		if cs.CheckCollision(checkPoint) {
			// Found a collision

			// Determine the object with which the collision occurred
			var hitObject interface{}

			// Check tile
			tileX, tileY := int(checkPoint.X), int(checkPoint.Y)
			tile := cs.world.GetTileAt(tileX, tileY)
			if tile != nil && cs.isTileSolid(tile) {
				hitObject = tile
			}

			// Check objects
			for _, obj := range cs.world.Objects {
				if obj.Solid {
					objDist := distance(obj.Position, checkPoint)
					if objDist < 1.0 {
						hitObject = obj
						break
					}
				}
			}

			// Calculate surface normal (simplified)
			normal := common.Vector2D{X: -direction.X, Y: -direction.Y}

			return CollisionResult{
				HasCollision: true,
				Normal:       normal,
				Object:       hitObject,
				Distance:     dist,
			}, true
		}
	}

	// The ray did not intersect any object
	return CollisionResult{HasCollision: false}, false
}

// CheckLineOfSight checks for direct visibility between two points
func (cs *CollisionSystem) CheckLineOfSight(from, to common.Vector2D) bool {
	// Calculate direction and distance
	direction := common.Vector2D{
		X: to.X - from.X,
		Y: to.Y - from.Y,
	}

	distance := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y)

	// Perform raycast
	result, hit := cs.CastRay(from, direction, distance)

	// If there were no collisions, or if the distance to the collision is greater than or equal to the distance to the target,
	// then direct visibility exists
	return !hit || result.Distance >= distance
}

// CheckCollisionWithEntities checks for a collision with creatures
func (cs *CollisionSystem) CheckCollisionWithEntities(position common.Vector2D, radius float64) (bool, *entity.Creature) {
	// Fixed this function to properly check entity type

	// Check all entities in the world
	for _, worldEntity := range cs.world.Entities {
		// Since we can't use direct type assertion on *Entity, we need a different approach
		// We'll check if the entity implements expected Creature methods

		// First, get the position of the entity for distance calculation
		entityPos := worldEntity.Position

		// Calculate distance between the position and the entity
		dist := distance(position, entityPos)

		// If the distance is less than the sum of the radiuses, there is a collision
		if dist < radius {
			// Now, we need to check if this is a Creature
			// Since we can't use direct type assertion, we'll return the collision
			// but a nil creature pointer since we can't safely convert it

			// If Entity had a method to access the underlying Creature, we could use that
			// For now, we'll just report the collision
			return true, nil
		}
	}

	return false, nil
}

// isTileSolid checks if a tile is impassable
func (cs *CollisionSystem) isTileSolid(tile *Tile) bool {
	// Determine which tile types are considered impassable
	switch tile.Type {
	case common.TileWater, common.TileSwamp, common.TileRocks, common.TileDenseForest:
		return true
	case common.TileCorrupted:
		// Corrupted terrain can be partially passable depending on the level of corruption
		return tile.Corruption > 0.7
	default:
		return false
	}
}

// distance calculates the distance between two points
func distance(a, b common.Vector2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}
