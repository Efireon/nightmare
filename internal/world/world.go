package world

import (
	"math"
	"math/rand"

	"nightmare/internal/common"

	"github.com/ojrac/opensimplex-go"
)

// Using common.TileType instead
type TileType = common.TileType

// Constants are already defined in common.TileType

// Tile represents a tile in the world
type Tile struct {
	Type       TileType
	Position   common.Vector2D // Added Position field to fix the error
	Elevation  float64
	Moisture   float64
	Corruption float64 // "Nightmareness" level from 0 to 1
	Objects    []common.WorldObject
}

// We'll use common.WorldObject instead

// Entity represents an entity in the world
type Entity struct {
	ID        int
	Type      string
	Position  common.Vector2D // Using common.Vector2D to fix type compatibility
	Direction float64
	Model     *EntityModel
	Behavior  EntityBehavior
}

// EntityModel represents an entity model
type EntityModel struct {
	Parts      []EntityPart
	Animations map[string][]int
}

// EntityPart represents a part of an entity model
type EntityPart struct {
	Type     string
	Texture  int
	Offset   common.Vector2D // Using common.Vector2D
	Scale    float64
	Rotation float64
}

// EntityBehavior defines entity behavior
type EntityBehavior interface {
	Update(world *World, entity *Entity)
}

// Using common.Vector2D instead of defining our own

// World represents the game world
type World struct {
	Width    int
	Height   int
	Tiles    [][]Tile
	Entities []*Entity
	Objects  []common.WorldObject // Using common.WorldObject
	nextID   int
	noise    opensimplex.Noise // Noise generator for procedural generation
}

// NewWorld creates a new world
func NewWorld(width, height int) (*World, error) {
	// Create world
	world := &World{
		Width:    width,
		Height:   height,
		Tiles:    make([][]Tile, height),
		Entities: []*Entity{},
		Objects:  []common.WorldObject{},
		nextID:   1,
		noise:    opensimplex.New(rand.Int63()),
	}

	// Initialize tiles
	for y := 0; y < height; y++ {
		world.Tiles[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			world.Tiles[y][x] = Tile{
				Type:       common.TileGrass,
				Position:   common.Vector2D{X: float64(x), Y: float64(y)}, // Initialize Position field
				Elevation:  0,
				Moisture:   0,
				Corruption: 0,
				Objects:    []common.WorldObject{},
			}
		}
	}

	// Generate landscape
	world.generateTerrain()

	// Place objects
	world.placeObjects()

	return world, nil
}

// generateTerrain generates landscape using noise algorithms
func (w *World) generateTerrain() {
	const elevationScale = 0.05
	const moistureScale = 0.07

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			// Generate height using Simplex noise
			elevation := (w.noise.Eval2(float64(x)*elevationScale, float64(y)*elevationScale) + 1) / 2

			// Generate moisture with a different noise octave
			moisture := (w.noise.Eval2(float64(x)*moistureScale+100, float64(y)*moistureScale+100) + 1) / 2

			// Save base values
			w.Tiles[y][x].Elevation = elevation
			w.Tiles[y][x].Moisture = moisture

			// Determine tile type based on height and moisture
			w.Tiles[y][x].Type = w.determineTileType(elevation, moisture)
		}
	}
}

// determineTileType determines tile type based on height and moisture
func (w *World) determineTileType(elevation, moisture float64) TileType {
	if elevation < 0.3 {
		if moisture > 0.7 {
			return common.TileSwamp
		} else if moisture > 0.4 {
			return common.TileWater
		} else {
			return common.TileGrass
		}
	} else if elevation < 0.6 {
		if moisture > 0.6 {
			return common.TileForest
		} else if moisture > 0.3 {
			return common.TileGrass
		} else {
			return common.TilePath
		}
	} else {
		if moisture > 0.7 {
			return common.TileDenseForest
		} else if moisture > 0.4 {
			return common.TileForest
		} else {
			return common.TileRocks
		}
	}
}

// placeObjects places objects in the world
func (w *World) placeObjects() {
	// Place trees in forest areas
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			tile := &w.Tiles[y][x]

			if tile.Type == common.TileForest {
				if rand.Float64() < 0.2 {
					tree := common.WorldObject{
						ID:          w.nextID,
						Type:        "tree",
						Position:    common.Vector2D{X: float64(x), Y: float64(y)},
						Solid:       true,
						Interactive: false,
					}
					tile.Objects = append(tile.Objects, tree)
					w.Objects = append(w.Objects, tree)
					w.nextID++
				}
			} else if tile.Type == common.TileDenseForest {
				if rand.Float64() < 0.5 {
					tree := common.WorldObject{
						ID:          w.nextID,
						Type:        "dense_tree",
						Position:    common.Vector2D{X: float64(x), Y: float64(y)},
						Solid:       true,
						Interactive: false,
					}
					tile.Objects = append(tile.Objects, tree)
					w.Objects = append(w.Objects, tree)
					w.nextID++
				}
			} else if tile.Type == common.TileRocks {
				if rand.Float64() < 0.1 {
					rock := common.WorldObject{
						ID:          w.nextID,
						Type:        "rock",
						Position:    common.Vector2D{X: float64(x), Y: float64(y)},
						Solid:       true,
						Interactive: false,
					}
					tile.Objects = append(tile.Objects, rock)
					w.Objects = append(w.Objects, rock)
					w.nextID++
				}
			}
		}
	}
}

// Update updates the world state
func (w *World) Update() {
	// Update all entities
	for _, entity := range w.Entities {
		if entity.Behavior != nil {
			entity.Behavior.Update(w, entity)
		}
	}
}

// GetTileAt returns the tile at the specified position
func (w *World) GetTileAt(x, y int) *Tile {
	if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return nil
	}
	return &w.Tiles[y][x]
}

// SpawnCreature creates a creature of the specified type at the specified position
func (w *World) SpawnCreature(creatureType string, position common.Vector2D) *Entity {
	// Create creature model
	model := generateCreatureModel(creatureType)

	// Create entity
	entity := &Entity{
		ID:        w.nextID,
		Type:      creatureType,
		Position:  position,
		Direction: rand.Float64() * 2 * math.Pi,
		Model:     model,
		Behavior:  NewBasicCreatureBehavior(),
	}

	w.Entities = append(w.Entities, entity)
	w.nextID++

	return entity
}

// ModifyEnvironment changes the environment around the specified position
func (w *World) ModifyEnvironment(position common.Vector2D, intensity float64) {
	// Influence radius
	radius := 10.0 + intensity*20.0
	radiusSq := radius * radius

	// Change tiles around the position
	for y := int(position.Y - radius); y <= int(position.Y+radius); y++ {
		for x := int(position.X - radius); x <= int(position.X+radius); x++ {
			// Check that coordinates are within the world
			if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
				continue
			}

			// Calculate distance to center
			dx := float64(x) - position.X
			dy := float64(y) - position.Y
			distSq := dx*dx + dy*dy

			// If the point is within the radius
			if distSq <= radiusSq {
				// Calculate influence strength (closer to center - stronger)
				strength := (1.0 - distSq/radiusSq) * intensity

				// Increase "nightmareness" level
				tile := &w.Tiles[y][x]
				tile.Corruption += strength * 0.3

				// Limit value
				if tile.Corruption > 1.0 {
					tile.Corruption = 1.0
				}

				// With high corruption level, change tile type
				if tile.Corruption > 0.7 {
					tile.Type = common.TileCorrupted
				}
			}
		}
	}
}

// generateCreatureModel generates a model for a creature of the specified type
func generateCreatureModel(creatureType string) *EntityModel {
	model := &EntityModel{
		Parts:      []EntityPart{},
		Animations: make(map[string][]int),
	}

	switch creatureType {
	case "shadow":
		// Body
		model.Parts = append(model.Parts, EntityPart{
			Type:     "body",
			Texture:  1, // Texture ID
			Offset:   common.Vector2D{X: 0, Y: 0},
			Scale:    1.0,
			Rotation: 0,
		})

		// Arms
		model.Parts = append(model.Parts, EntityPart{
			Type:     "arm",
			Texture:  2,
			Offset:   common.Vector2D{X: -0.5, Y: 0},
			Scale:    0.7,
			Rotation: 0,
		})

		model.Parts = append(model.Parts, EntityPart{
			Type:     "arm",
			Texture:  2,
			Offset:   common.Vector2D{X: 0.5, Y: 0},
			Scale:    0.7,
			Rotation: 0,
		})

	case "spider":
		// Body
		model.Parts = append(model.Parts, EntityPart{
			Type:     "body",
			Texture:  10,
			Offset:   common.Vector2D{X: 0, Y: 0},
			Scale:    1.0,
			Rotation: 0,
		})

		// Legs (8)
		for i := 0; i < 8; i++ {
			angle := float64(i) * math.Pi / 4
			model.Parts = append(model.Parts, EntityPart{
				Type:     "leg",
				Texture:  11,
				Offset:   common.Vector2D{X: math.Cos(angle) * 0.5, Y: math.Sin(angle) * 0.5},
				Scale:    0.6,
				Rotation: angle,
			})
		}

	// Other creature types...
	default:
		// Base model for unknown type
		model.Parts = append(model.Parts, EntityPart{
			Type:     "body",
			Texture:  0,
			Offset:   common.Vector2D{X: 0, Y: 0},
			Scale:    1.0,
			Rotation: 0,
		})
	}

	// Add base animations
	model.Animations["idle"] = []int{0, 1, 2, 1}
	model.Animations["walk"] = []int{3, 4, 5, 6}
	model.Animations["attack"] = []int{7, 8, 9}

	return model
}

// BasicCreatureBehavior represents basic creature behavior
type BasicCreatureBehavior struct {
	state          string
	targetPoint    common.Vector2D
	idleTime       int
	attackRange    float64
	detectionRange float64
}

// NewBasicCreatureBehavior creates new basic creature behavior
func NewBasicCreatureBehavior() *BasicCreatureBehavior {
	return &BasicCreatureBehavior{
		state:          "wander",
		targetPoint:    common.Vector2D{X: 0, Y: 0},
		idleTime:       0,
		attackRange:    5.0,
		detectionRange: 20.0,
	}
}

// Update updates creature behavior
func (b *BasicCreatureBehavior) Update(world *World, entity *Entity) {
	// Creature behavior logic will go here
	// ...
}

// Add these functions to resolve undefined references in other packages
// These are the behavior types that were referenced but not defined

// PatrolBehavior represents patrolling behavior
type PatrolBehavior struct {
	*BasicCreatureBehavior
	patrolPoints []common.Vector2D
	currentPoint int
}

// NewPatrolBehavior creates a new patrol behavior
func NewPatrolBehavior() *PatrolBehavior {
	return &PatrolBehavior{
		BasicCreatureBehavior: NewBasicCreatureBehavior(),
		patrolPoints:          []common.Vector2D{},
		currentPoint:          0,
	}
}

// Update updates patrol behavior
func (b *PatrolBehavior) Update(world *World, entity *Entity) {
	// Patrol behavior logic will go here
	// ...
}

// AggressiveBehavior represents aggressive behavior
type AggressiveBehavior struct {
	*BasicCreatureBehavior
	aggressionLevel float64
}

// NewAggressiveBehavior creates a new aggressive behavior
func NewAggressiveBehavior() *AggressiveBehavior {
	return &AggressiveBehavior{
		BasicCreatureBehavior: NewBasicCreatureBehavior(),
		aggressionLevel:       0.8,
	}
}

// Update updates aggressive behavior
func (b *AggressiveBehavior) Update(world *World, entity *Entity) {
	// Aggressive behavior logic will go here
	// ...
}

// StalkerBehavior represents stalker behavior
type StalkerBehavior struct {
	*BasicCreatureBehavior
	stalkDistance float64
	revealChance  float64
}

// NewStalkerBehavior creates a new stalker behavior
func NewStalkerBehavior() *StalkerBehavior {
	return &StalkerBehavior{
		BasicCreatureBehavior: NewBasicCreatureBehavior(),
		stalkDistance:         15.0,
		revealChance:          0.2,
	}
}

// Update updates stalker behavior
func (b *StalkerBehavior) Update(world *World, entity *Entity) {
	// Stalker behavior logic will go here
	// ...
}
