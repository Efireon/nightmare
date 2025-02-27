package ai

import (
	"fmt"
	"math"
	"sort"
	"time"

	"nightmare/internal/common"
	"nightmare/internal/entity"
)

// PlayerPattern represents a player behavior pattern
type PlayerPattern struct {
	Name        string
	Description string
	Weight      float64
}

// MovementAnalysis contains the results of player movement analysis
type MovementAnalysis struct {
	AverageSpeed     float64
	DirectionChanges int
	ExplorationArea  float64
	PathRepetition   float64
	PreferredAreas   []entity.Vector2D
}

// InteractionAnalysis contains the results of interaction analysis
type InteractionAnalysis struct {
	InteractionRate       float64
	PreferredInteractions map[string]int
	ResponseToScareEvents map[common.ScareEventType]float64 // Changed to use common.ScareEventType
	HealthLossRate        float64
	SanityLossRate        float64
}

// Analyzer analyzes player behavior
type Analyzer struct {
	player              *entity.Player
	detectedPatterns    []PlayerPattern
	movementAnalysis    MovementAnalysis
	interactionAnalysis InteractionAnalysis
	scareHistory        []common.ScareEvent // Changed to use common.ScareEvent
	lastAnalysisTime    time.Time

	positionHistory []entity.Vector2D
	areaVisits      map[string]int  // key: "x,y" for sector, value: number of visits
	sectorsExplored map[string]bool // key: "x,y" for sector, value: whether explored

	sectorSize float64     // size of one sector for analysis
	heatmap    [][]float64 // visit heatmap

	scareResponses map[common.ScareEventType][]float64 // Changed to use common.ScareEventType
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(player *entity.Player) *Analyzer {
	return &Analyzer{
		player:           player,
		detectedPatterns: []PlayerPattern{},
		movementAnalysis: MovementAnalysis{
			AverageSpeed:     0,
			DirectionChanges: 0,
			ExplorationArea:  0,
			PathRepetition:   0,
			PreferredAreas:   []entity.Vector2D{},
		},
		interactionAnalysis: InteractionAnalysis{
			InteractionRate:       0,
			PreferredInteractions: make(map[string]int),
			ResponseToScareEvents: make(map[common.ScareEventType]float64), // Changed to use common.ScareEventType
			HealthLossRate:        0,
			SanityLossRate:        0,
		},
		scareHistory:     []common.ScareEvent{}, // Changed to use common.ScareEvent
		lastAnalysisTime: time.Now(),
		positionHistory:  []entity.Vector2D{},
		areaVisits:       make(map[string]int),
		sectorsExplored:  make(map[string]bool),
		sectorSize:       5.0,                                       // World unit sector size
		heatmap:          make([][]float64, 50),                     // 50x50 heatmap
		scareResponses:   make(map[common.ScareEventType][]float64), // Changed to use common.ScareEventType
	}
}

// AnalyzePlayer performs comprehensive analysis of player behavior
func (a *Analyzer) AnalyzePlayer() {
	// Record current player position
	a.recordPlayerPosition()

	// Analyze movements
	a.analyzeMovement()

	// Analyze interactions
	a.analyzeInteractions()

	// Update heatmap
	a.updateHeatmap()

	// Detect behavior patterns
	a.detectPatterns()

	// Update last analysis time
	a.lastAnalysisTime = time.Now()
}

// recordPlayerPosition records the current player position
func (a *Analyzer) recordPlayerPosition() {
	a.positionHistory = append(a.positionHistory, a.player.Position)

	// Limit history size
	if len(a.positionHistory) > 1000 {
		a.positionHistory = a.positionHistory[len(a.positionHistory)-1000:]
	}

	// Update sector visits
	sectorX := int(a.player.Position.X / a.sectorSize)
	sectorY := int(a.player.Position.Y / a.sectorSize)
	sectorKey := makeKey(sectorX, sectorY)

	a.areaVisits[sectorKey]++
	a.sectorsExplored[sectorKey] = true
}

// makeKey creates a string key from coordinates
func makeKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// analyzeMovement analyzes player movement
func (a *Analyzer) analyzeMovement() {
	if len(a.positionHistory) < 2 {
		return
	}

	// Calculate average speed
	totalDistance := 0.0
	directionChanges := 0
	for i := 1; i < len(a.positionHistory); i++ {
		// Distance between consecutive points
		dist := distance(a.positionHistory[i-1], a.positionHistory[i])
		totalDistance += dist

		// Direction changes
		if i > 1 {
			prev := a.positionHistory[i-2]
			curr := a.positionHistory[i-1]
			next := a.positionHistory[i]

			dir1 := math.Atan2(curr.Y-prev.Y, curr.X-prev.X)
			dir2 := math.Atan2(next.Y-curr.Y, next.X-curr.X)

			// Normalize angle difference
			diff := math.Abs(dir2 - dir1)
			if diff > math.Pi {
				diff = 2*math.Pi - diff
			}

			// Count direction change if angle is greater than threshold
			if diff > math.Pi/4 {
				directionChanges++
			}

		}
	}

	// Update analysis results
	a.movementAnalysis.AverageSpeed = totalDistance / float64(len(a.positionHistory)-1)
	a.movementAnalysis.DirectionChanges = directionChanges

	// Calculate explored area
	a.movementAnalysis.ExplorationArea = float64(len(a.sectorsExplored)) * (a.sectorSize * a.sectorSize)

	// Find preferred areas
	a.findPreferredAreas()

	// Calculate path repetition
	a.calculatePathRepetition()
}

// analyzeInteractions analyzes player interactions
func (a *Analyzer) analyzeInteractions() {
	// Analyze action records
	totalActions := len(a.player.ActionLog)
	if totalActions == 0 {
		return
	}

	interactions := 0
	for _, action := range a.player.ActionLog {
		if action.Action == entity.ActionInteract {
			interactions++

			// Add interaction type if present
			if action.InteractionType != "" {
				a.interactionAnalysis.PreferredInteractions[action.InteractionType]++
			}
		}
	}

	// Calculate interaction frequency
	a.interactionAnalysis.InteractionRate = float64(interactions) / float64(totalActions)
}

// RecordScareResponse records player's response to a scare event
func (a *Analyzer) RecordScareResponse(event common.ScareEvent, response float64) { // Changed to use common.ScareEvent
	// Record event in history
	a.scareHistory = append(a.scareHistory, event)

	// Limit history size
	if len(a.scareHistory) > 50 {
		a.scareHistory = a.scareHistory[len(a.scareHistory)-50:]
	}

	// Record response
	if _, ok := a.scareResponses[event.Type]; !ok {
		a.scareResponses[event.Type] = []float64{}
	}
	a.scareResponses[event.Type] = append(a.scareResponses[event.Type], response)

	// Calculate average response
	total := 0.0
	for _, resp := range a.scareResponses[event.Type] {
		total += resp
	}
	avgResponse := total / float64(len(a.scareResponses[event.Type]))

	// Update analysis
	a.interactionAnalysis.ResponseToScareEvents[event.Type] = avgResponse
}

// updateHeatmap updates the visit heatmap
func (a *Analyzer) updateHeatmap() {
	// Initialize heatmap if needed
	if len(a.heatmap) == 0 || len(a.heatmap[0]) == 0 {
		for i := range a.heatmap {
			a.heatmap[i] = make([]float64, 50)
		}
	}

	// Update heatmap values
	for key, visits := range a.areaVisits {
		var x, y int
		fmt.Sscanf(key, "%d,%d", &x, &y)

		// Convert world coordinates to heatmap coordinates
		heatmapX := x * 50 / 256 // Assuming 256x256 world
		heatmapY := y * 50 / 256

		// Check that coordinates are within the map
		if heatmapX >= 0 && heatmapX < 50 && heatmapY >= 0 && heatmapY < 50 {
			// Increase value based on number of visits
			a.heatmap[heatmapY][heatmapX] = math.Min(1.0, float64(visits)/10.0)
		}
	}
}

// findPreferredAreas finds preferred areas
func (a *Analyzer) findPreferredAreas() {
	// Clear previous results
	a.movementAnalysis.PreferredAreas = []entity.Vector2D{}

	// Find most visited sectors
	type sectorVisit struct {
		key    string
		visits int
	}

	// Convert map to slice for sorting
	visits := []sectorVisit{}
	for key, count := range a.areaVisits {
		visits = append(visits, sectorVisit{key: key, visits: count})
	}

	// Sort by descending number of visits
	sort.Slice(visits, func(i, j int) bool {
		return visits[i].visits > visits[j].visits
	})

	// Take top-5 or fewer
	count := min(5, len(visits))
	for i := 0; i < count; i++ {
		var x, y int
		fmt.Sscanf(visits[i].key, "%d,%d", &x, &y)

		// Convert sector coordinates to world coordinates
		worldX := float64(x)*a.sectorSize + a.sectorSize/2
		worldY := float64(y)*a.sectorSize + a.sectorSize/2

		a.movementAnalysis.PreferredAreas = append(a.movementAnalysis.PreferredAreas,
			entity.Vector2D{X: worldX, Y: worldY})
	}
}

// calculatePathRepetition calculates path repetition
func (a *Analyzer) calculatePathRepetition() {
	totalSectors := len(a.sectorsExplored)
	if totalSectors == 0 {
		a.movementAnalysis.PathRepetition = 0
		return
	}

	// Count visits for each sector
	totalVisits := 0
	for _, count := range a.areaVisits {
		totalVisits += count
	}

	// Calculate average visits per sector
	avgVisits := float64(totalVisits) / float64(totalSectors)

	// Repetition - ratio of average visits to expected (1)
	a.movementAnalysis.PathRepetition = math.Max(0, avgVisits-1)
}

// detectPatterns detects player behavior patterns
func (a *Analyzer) detectPatterns() {
	// Clear previous results
	a.detectedPatterns = []PlayerPattern{}

	// Analyze based on movement
	a.detectMovementPatterns()

	// Analyze based on interactions
	a.detectInteractionPatterns()

	// Analyze based on fear responses
	a.detectFearResponsePatterns()

	// Sort patterns by weight
	sort.Slice(a.detectedPatterns, func(i, j int) bool {
		return a.detectedPatterns[i].Weight > a.detectedPatterns[j].Weight
	})
}

// detectMovementPatterns detects movement patterns
func (a *Analyzer) detectMovementPatterns() {
	// Explorer - explores a lot, repeats little
	if a.movementAnalysis.ExplorationArea > 500 && a.movementAnalysis.PathRepetition < 2 {
		a.addPattern(PlayerPattern{
			Name:        "explorer",
			Description: "Player actively explores the world without lingering in one place",
			Weight:      0.8 - (a.movementAnalysis.PathRepetition / 10),
		})
	}

	// Cautious - moves slowly, changes direction often
	if a.movementAnalysis.AverageSpeed < 1.5 && a.movementAnalysis.DirectionChanges > 30 {
		a.addPattern(PlayerPattern{
			Name:        "cautious",
			Description: "Player is cautious, moves slowly and changes direction often",
			Weight:      0.9 - (a.movementAnalysis.AverageSpeed / 3),
		})
	}

	// Determined - moves quickly and directly
	if a.movementAnalysis.AverageSpeed > 2.0 && a.movementAnalysis.DirectionChanges < 15 {
		a.addPattern(PlayerPattern{
			Name:        "determined",
			Description: "Player moves quickly in the chosen direction",
			Weight:      0.7 + (a.movementAnalysis.AverageSpeed / 5),
		})
	}

	// Indecisive - stays in the same place a lot
	if a.movementAnalysis.ExplorationArea < 200 && a.movementAnalysis.PathRepetition > 3 {
		a.addPattern(PlayerPattern{
			Name:        "indecisive",
			Description: "Player is indecisive, often returns to the same places",
			Weight:      0.6 + (a.movementAnalysis.PathRepetition / 5),
		})
	}
}

// detectInteractionPatterns detects interaction patterns
func (a *Analyzer) detectInteractionPatterns() {
	// Interactive - interacts with the world often
	if a.interactionAnalysis.InteractionRate > 0.3 {
		a.addPattern(PlayerPattern{
			Name:        "interactive",
			Description: "Player actively interacts with the environment",
			Weight:      0.7 + a.interactionAnalysis.InteractionRate,
		})
	}

	// Passive - rarely interacts with the world
	if a.interactionAnalysis.InteractionRate < 0.1 {
		a.addPattern(PlayerPattern{
			Name:        "passive",
			Description: "Player rarely interacts with the environment",
			Weight:      0.6 + (0.1 - a.interactionAnalysis.InteractionRate),
		})
	}
}

// detectFearResponsePatterns detects fear response patterns
func (a *Analyzer) detectFearResponsePatterns() {
	// General response to fear
	avgResponse := 0.0
	count := 0

	for _, response := range a.interactionAnalysis.ResponseToScareEvents {
		avgResponse += response
		count++
	}

	if count > 0 {
		avgResponse /= float64(count)

		// Fear resistant
		if avgResponse < 0.3 {
			a.addPattern(PlayerPattern{
				Name:        "fearless",
				Description: "Player reacts weakly to frightening events",
				Weight:      0.8 - avgResponse,
			})
		}

		// Easily frightened
		if avgResponse > 0.7 {
			a.addPattern(PlayerPattern{
				Name:        "easily_scared",
				Description: "Player reacts strongly to frightening events",
				Weight:      0.7 + avgResponse,
			})
		}
	}

	// Reaction to specific types of fear
	if response, ok := a.interactionAnalysis.ResponseToScareEvents[common.EventSuddenNoise]; ok && response > 0.8 { // Changed to use common.EventSuddenNoise
		a.addPattern(PlayerPattern{
			Name:        "startles_easily",
			Description: "Player is particularly sensitive to sudden sounds",
			Weight:      0.7 + response,
		})
	}

	if response, ok := a.interactionAnalysis.ResponseToScareEvents[common.EventCreatureAppearance]; ok && response > 0.8 { // Changed to use common.EventCreatureAppearance
		a.addPattern(PlayerPattern{
			Name:        "monster_phobia",
			Description: "Player is particularly afraid of encountering creatures",
			Weight:      0.7 + response,
		})
	}
}

// addPattern adds a behavior pattern
func (a *Analyzer) addPattern(pattern PlayerPattern) {
	// Check if pattern already exists
	for i, p := range a.detectedPatterns {
		if p.Name == pattern.Name {
			// Update weight of existing pattern
			a.detectedPatterns[i].Weight = (a.detectedPatterns[i].Weight + pattern.Weight) / 2
			return
		}
	}

	// Add new pattern
	a.detectedPatterns = append(a.detectedPatterns, pattern)
}

// GetTopPatterns returns the most prominent behavior patterns
func (a *Analyzer) GetTopPatterns(count int) []PlayerPattern {
	if count > len(a.detectedPatterns) {
		count = len(a.detectedPatterns)
	}

	return a.detectedPatterns[:count]
}

// GetHeatmap returns the visit heatmap
func (a *Analyzer) GetHeatmap() [][]float64 {
	return a.heatmap
}

// GetMovementAnalysis returns the movement analysis results
func (a *Analyzer) GetMovementAnalysis() MovementAnalysis {
	return a.movementAnalysis
}

// GetInteractionAnalysis returns the interaction analysis results
func (a *Analyzer) GetInteractionAnalysis() InteractionAnalysis {
	return a.interactionAnalysis
}

// distance calculates the distance between two points
func distance(a, b entity.Vector2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// min returns the minimum of two numbers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
