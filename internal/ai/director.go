package ai

import (
	"math"
	"math/rand"
	"time"

	"nightmare/internal/common"
	"nightmare/internal/entity"
)

// Using ScareEvent from common package
// Note: We're now using the common package's version of ScareEvent which has the Timestamp field

// BehaviorPattern describes a player behavior pattern
type BehaviorPattern struct {
	MovementPreference    float64  // Movement preference (positive - active, negative - passive)
	ExplorationPreference float64  // Exploration preference (high - explorer, low - linear player)
	RiskTolerance         float64  // Risk tolerance (high - brave, low - cautious)
	ReactivityToScares    float64  // Reaction to scares (high - strong reaction, low - weak reaction)
	PreferredInteractions []string // Types of interactions that the player prefers
}

// Director represents the AI director that manages game events
type Director struct {
	player             *entity.Player
	world              interface{} // Using interface to avoid direct import of world
	playerBehavior     BehaviorPattern
	scareHistory       []common.ScareEvent
	scareEffectiveness map[common.ScareEventType]float64 // Effectiveness of different scare types
	lastAnalysisTime   time.Time
	mood               float64 // General "mood" of the director from 0 (calm) to 1 (aggressive)
	tension            float64 // Current tension level from 0 to 1
}

// NewDirector creates a new AI director
func NewDirector(player *entity.Player, world interface{}) *Director {
	return &Director{
		player: player,
		world:  world,
		playerBehavior: BehaviorPattern{
			MovementPreference:    0.5,
			ExplorationPreference: 0.5,
			RiskTolerance:         0.5,
			ReactivityToScares:    0.5,
			PreferredInteractions: []string{},
		},
		scareHistory:       []common.ScareEvent{},
		scareEffectiveness: make(map[common.ScareEventType]float64),
		lastAnalysisTime:   time.Now(),
		mood:               0.3, // Initial mood
		tension:            0.1, // Initial tension
	}
}

// AnalyzePlayerBehavior analyzes player behavior
func (d *Director) AnalyzePlayerBehavior() {
	// If there are no player action logs, do nothing
	if len(d.player.ActionLog) == 0 {
		return
	}

	// Analyze only logs since the last analysis
	recentLogs := []entity.PlayerActionRecord{}
	for _, log := range d.player.ActionLog {
		if log.Timestamp.After(d.lastAnalysisTime) {
			recentLogs = append(recentLogs, log)
		}
	}

	d.lastAnalysisTime = time.Now()

	// If there are no new logs, do nothing
	if len(recentLogs) == 0 {
		return
	}

	// Movement analysis
	moveCount := 0
	for _, log := range recentLogs {
		if log.Action == entity.ActionMove {
			moveCount++
		}
	}

	movementRatio := float64(moveCount) / float64(len(recentLogs))

	// Update movement preferences (with inertia)
	d.playerBehavior.MovementPreference = d.playerBehavior.MovementPreference*0.8 + movementRatio*0.2

	// Analyze exploration (how much the player deviates from the direct path)
	// This is a more complex analysis that we'll simplify for this example

	// Update other aspects of behavior
	// ...

	// Analyze the effectiveness of past attempts to scare
	d.analyzeScareEffectiveness()

	// Update the director's mood and tension level
	d.updateMoodAndTension()
}

// AdjustWorld modifies the world based on analysis of player behavior
func (d *Director) AdjustWorld() {
	// Decide whether to create a scare event
	if d.shouldCreateScareEvent() {
		event := d.createScareEvent()
		d.executeScareEvent(event)
	}

	// Modify the surrounding world
	d.modifyEnvironment()

	// Manage creatures
	d.manageCreatures()
}

// shouldCreateScareEvent determines if a scare event should be created
func (d *Director) shouldCreateScareEvent() bool {
	// Base chance based on tension
	baseChance := d.tension * 0.1

	// Increase chance if player hasn't been scared for a while
	if len(d.scareHistory) > 0 {
		lastScare := d.scareHistory[len(d.scareHistory)-1]
		timeSinceLast := time.Since(lastScare.Timestamp)
		if timeSinceLast > 30*time.Second {
			baseChance += 0.1
		}
		if timeSinceLast > 60*time.Second {
			baseChance += 0.2
		}
	} else {
		// If there hasn't been a scare yet, increase the chance
		baseChance += 0.3
	}

	// Add randomness
	return rand.Float64() < baseChance
}

// createScareEvent creates a scare event based on player behavior
func (d *Director) createScareEvent() common.ScareEvent {
	// Choose event type based on effectiveness
	eventType := d.chooseEventType()

	// Determine intensity based on mood and player analysis
	intensity := d.mood * (0.7 + rand.Float64()*0.3)

	// If the player reacts weakly to scares, increase intensity
	if d.playerBehavior.ReactivityToScares < 0.3 {
		intensity *= 1.5
	}

	// Limit intensity
	if intensity > 1.0 {
		intensity = 1.0
	}

	// Create the event
	event := common.ScareEvent{
		Type:      eventType,
		Intensity: intensity,
		Position:  common.VectorFromEntity(d.player.Position), // By default near the player
		Duration:  time.Duration(2+rand.Intn(5)) * time.Second,
		Timestamp: time.Now(), // Add timestamp to fix the missing field error
	}

	// For some event types, additional configuration is needed
	if eventType == common.EventCreatureAppearance {
		// Choose creature type
		event.CreatureType = d.chooseCreatureType()

		// Set creature spawn position
		angle := rand.Float64() * 2 * math.Pi
		distance := 10.0 + rand.Float64()*20.0
		event.Position = common.Vector2D{
			X: d.player.Position.X + math.Cos(angle)*distance,
			Y: d.player.Position.Y + math.Sin(angle)*distance,
		}
	}

	return event
}

// chooseEventType chooses an event type based on effectiveness
func (d *Director) chooseEventType() common.ScareEventType {
	// If we don't have data on effectiveness, choose a random type
	if len(d.scareEffectiveness) == 0 {
		return common.ScareEventType(rand.Intn(6))
	}

	// Choose more effective types with higher probability
	// ...

	// Simplified version - random choice
	return common.ScareEventType(rand.Intn(6))
}

// chooseCreatureType chooses a creature type
func (d *Director) chooseCreatureType() string {
	creatureTypes := []string{
		"shadow", "spider", "phantom", "doppelganger", "wendigo", "faceless",
	}

	return creatureTypes[rand.Intn(len(creatureTypes))]
}

// executeScareEvent executes a scare event
func (d *Director) executeScareEvent(event common.ScareEvent) {
	// Add the event to history
	d.scareHistory = append(d.scareHistory, event)

	// Perform actions depending on the event type
	switch event.Type {
	case common.EventAmbientSound:
		// Play sound
		// ...

	case common.EventSuddenNoise:
		// Sudden loud noise
		// ...

	case common.EventCreatureAppearance:
		// Create creature - this requires interacting with the world interface
		// We'll add a type check to ensure we're using the world correctly
		if worldObj, ok := d.world.(interface {
			SpawnCreature(string, common.Vector2D) interface{}
		}); ok {
			worldObj.SpawnCreature(event.CreatureType, event.Position)
		}

	case common.EventEnvironmentChange:
		// Change environment
		if worldObj, ok := d.world.(interface {
			ModifyEnvironment(common.Vector2D, float64)
		}); ok {
			worldObj.ModifyEnvironment(event.Position, event.Intensity)
		}

	case common.EventHallucination:
		// Create hallucination
		// ...

	case common.EventWhisper:
		// Whisper
		// ...
	}

	// Reduce player's sanity based on event intensity
	d.player.ReduceSanity(event.Intensity * 5)
}

// analyzeScareEffectiveness analyzes the effectiveness of past attempts to scare
func (d *Director) analyzeScareEffectiveness() {
	// This method will analyze how much the sanity level changed
	// after each scare event
	// ...
}

// updateMoodAndTension updates the director's mood and tension level
func (d *Director) updateMoodAndTension() {
	// Mood depends on how effectively we're scaring the player
	// ...

	// Tension increases over time and drops after a successful scare
	d.tension += 0.01

	// Limit tension
	if d.tension > 1.0 {
		d.tension = 1.0
	}
}

// modifyEnvironment modifies the surrounding world
func (d *Director) modifyEnvironment() {
	// Logic for modifying the surrounding world will go here
	// ...
}

// manageCreatures manages creatures in the world
func (d *Director) manageCreatures() {
	// Logic for managing creatures will go here
	// ...
}
