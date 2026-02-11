package game

import (
	"fmt"
	"math"
	"math/rand"
	"stardew-agent/internal/models"
	"time"
)

// BehaviorEngine handles AI-driven NPC behavior decisions
type BehaviorEngine struct {
	moodSystem     *MoodSystem
	random         *rand.Rand
	interactionLog map[string][]string // NPC ID -> recent interactions
}

// NewBehaviorEngine creates a new behavior engine
func NewBehaviorEngine(moodSystem *MoodSystem) *BehaviorEngine {
	return &BehaviorEngine{
		moodSystem:     moodSystem,
		random:         rand.New(rand.NewSource(time.Now().UnixNano())),
		interactionLog: make(map[string][]string),
	}
}

// MakeDecision generates a behavior decision for an NPC
func (be *BehaviorEngine) MakeDecision(req *models.AIDecisionRequest, npcs []models.NPCState) *models.AIDecisionResponse {
	// Find the NPC
	var npc *models.NPCState
	for i := range npcs {
		if npcs[i].ID == req.NPCID {
			npc = &npcs[i]
			break
		}
	}
	if npc == nil {
		return nil
	}

	// Build context
	context := be.buildContext(npc, req, npcs)

	// Update mood based on current context
	moodState := be.moodSystem.UpdateMood(npc.ID, context)

	// Generate decision factors
	factors := be.calculateDecisionFactors(context)

	// Determine best behavior
	decision := be.selectBehavior(context, factors, req.OverrideSchedule)

	// Check for NPC interaction trigger
	var interaction *models.NPCInteraction
	if decision.BehaviorType == models.BehaviorInteractNPC || decision.BehaviorType == models.BehaviorInteractPlayer {
		interaction = be.createInteraction(decision, context)
	}

	// Generate mood reason
	moodReason := be.generateMoodReason(moodState, context)

	return &models.AIDecisionResponse{
		Decision:            *decision,
		CurrentMood:         moodState.CurrentMood,
		MoodReason:          moodReason,
		ScheduleModified:    decision.DeviatesFromSchedule,
		InteractionTriggered: interaction,
	}
}

// buildContext builds the behavior context for decision making
func (be *BehaviorEngine) buildContext(npc *models.NPCState, req *models.AIDecisionRequest, npcs []models.NPCState) *models.BehaviorContext {
	// Get nearby NPCs
	var nearbyNPCs []models.NPCState
	for _, other := range npcs {
		if other.ID == npc.ID {
			continue
		}
		dx := other.Position.X - npc.Position.X
		dy := other.Position.Y - npc.Position.Y
		distance := math.Sqrt(float64(dx*dx + dy*dy))
		if distance <= 10 {
			nearbyNPCs = append(nearbyNPCs, other)
		}
	}

	// Get or create personality
	personality := be.moodSystem.GetPersonality(npc.ID)
	if personality == nil {
		personality = be.getDefaultPersonality(npc.ID)
	}

	// Get mood state
	moodState := be.moodSystem.GetMoodState(npc.ID)
	if moodState == nil {
		moodState = &models.NPCMoodState{
			NPCID:         npc.ID,
			CurrentMood:   models.MoodNeutral,
			MoodIntensity: 0.5,
		}
	}

	return &models.BehaviorContext{
		NPC:               *npc,
		PlayerPosition:    req.PlayerPosition,
		PlayerFriendship:  npc.Friendship,
		Weather:           req.Weather,
		Time:              req.CurrentTime,
		NearbyNPCs:        nearbyNPCs,
		RecentInteractions: be.interactionLog[npc.ID],
		DayEvents:         req.RecentEvents,
		Personality:       *personality,
		MoodState:         *moodState,
	}
}

// getDefaultPersonality returns default personality for an NPC
func (be *BehaviorEngine) getDefaultPersonality(npcID string) *models.NPCPersonality {
	// NPC-specific personalities
	personalities := map[string]*models.NPCPersonality{
		"lewis": {
			NPCID:              "lewis",
			Sociability:        0.7,
			Spontaneity:        0.2,
			WorkEthic:          0.8,
			OutdoorsPreference: 0.4,
			MorningPerson:      0.9,
			GiftAppreciation:   0.6,
			WeatherSensitivity: 0.3,
		},
		"pierre": {
			NPCID:              "pierre",
			Sociability:        0.5,
			Spontaneity:        0.3,
			WorkEthic:          0.9,
			OutdoorsPreference: 0.2,
			MorningPerson:      0.7,
			GiftAppreciation:   0.8,
			WeatherSensitivity: 0.2,
		},
		"robin": {
			NPCID:              "robin",
			Sociability:        0.6,
			Spontaneity:        0.4,
			WorkEthic:          0.85,
			OutdoorsPreference: 0.6,
			MorningPerson:      0.8,
			GiftAppreciation:   0.5,
			WeatherSensitivity: 0.3,
		},
		"haley": {
			NPCID:              "haley",
			Sociability:        0.6,
			Spontaneity:        0.7,
			WorkEthic:          0.3,
			OutdoorsPreference: 0.7,
			MorningPerson:      0.2,
			GiftAppreciation:   0.9,
			WeatherSensitivity: 0.6,
		},
		"willy": {
			NPCID:              "willy",
			Sociability:        0.5,
			Spontaneity:        0.4,
			WorkEthic:          0.7,
			OutdoorsPreference: 0.9,
			MorningPerson:      0.85,
			GiftAppreciation:   0.6,
			WeatherSensitivity: 0.5,
		},
	}

	if p, ok := personalities[npcID]; ok {
		return p
	}

	return &models.NPCPersonality{
		NPCID:              npcID,
		Sociability:        0.5,
		Spontaneity:        0.3,
		WorkEthic:          0.6,
		OutdoorsPreference: 0.5,
		MorningPerson:      0.5,
		GiftAppreciation:   0.7,
		WeatherSensitivity: 0.4,
	}
}

// calculateDecisionFactors calculates all decision factors
func (be *BehaviorEngine) calculateDecisionFactors(context *models.BehaviorContext) []models.DecisionFactor {
	var factors []models.DecisionFactor

	// Player proximity factor
	playerDist := math.Sqrt(float64(
		(context.NPC.Position.X-context.PlayerPosition.X)*(context.NPC.Position.X-context.PlayerPosition.X) +
			(context.NPC.Position.Y-context.PlayerPosition.Y)*(context.NPC.Position.Y-context.PlayerPosition.Y),
	))
	playerProximityScore := 1.0 / (1.0 + playerDist/10.0)
	factors = append(factors, models.DecisionFactor{
		Name:     "player_proximity",
		Weight:   0.2 + float64(context.PlayerFriendship)/2500.0, // Higher friendship = more weight
		Value:    playerProximityScore,
		RawScore: playerProximityScore,
	})

	// Social need factor
	socialNeed := 0.0
	if context.Personality.Sociability > 0.5 {
		if len(context.NearbyNPCs) == 0 {
			socialNeed = context.Personality.Sociability
		} else {
			socialNeed = 0.3 * context.Personality.Sociability
		}
	}
	factors = append(factors, models.DecisionFactor{
		Name:     "social_need",
		Weight:   0.15,
		Value:    socialNeed,
		RawScore: socialNeed,
	})

	// Weather preference
	weatherScore := 0.0
	if context.Weather == models.WeatherSunny && context.Personality.OutdoorsPreference > 0.5 {
		weatherScore = context.Personality.OutdoorsPreference
	} else if context.Weather == models.WeatherRainy && context.Personality.OutdoorsPreference < 0.3 {
		weatherScore = 0.3 // Enjoy indoor activities
	} else if context.Weather == models.WeatherStormy {
		weatherScore = -0.3 // Generally negative
	}
	factors = append(factors, models.DecisionFactor{
		Name:     "weather_preference",
		Weight:   0.1,
		Value:    weatherScore,
		RawScore: weatherScore,
	})

	// Mood influence
	moodScore := 0.0
	switch context.MoodState.CurrentMood {
	case models.MoodHappy, models.MoodExcited:
		moodScore = 0.4 * context.MoodState.MoodIntensity
	case models.MoodSad, models.MoodAnxious:
		moodScore = -0.3 * context.MoodState.MoodIntensity
	case models.MoodAngry:
		moodScore = -0.2
	case models.MoodTired:
		moodScore = -0.4
	}
	factors = append(factors, models.DecisionFactor{
		Name:     "mood_influence",
		Weight:   0.2,
		Value:    moodScore,
		RawScore: moodScore,
	})

	// Time-based work ethic
	timeScore := 0.0
	if context.Time.Hour >= 9 && context.Time.Hour < 17 {
		if context.Personality.WorkEthic > 0.6 {
			timeScore = context.Personality.WorkEthic
		}
	} else if context.Time.Hour >= 17 && context.Time.Hour < 22 {
		timeScore = 0.3 // Social hours
	} else {
		timeScore = -0.2 // Rest hours
	}
	factors = append(factors, models.DecisionFactor{
		Name:     "time_of_day",
		Weight:   0.15,
		Value:    timeScore,
		RawScore: timeScore,
	})

	// Spontaneity factor
	spontaneityScore := context.Personality.Spontaneity * be.random.Float64()
	factors = append(factors, models.DecisionFactor{
		Name:     "spontaneity",
		Weight:   0.1,
		Value:    spontaneityScore,
		RawScore: spontaneityScore,
	})

	// Nearby NPCs factor
	nearbyScore := float64(len(context.NearbyNPCs)) * 0.2 * context.Personality.Sociability
	factors = append(factors, models.DecisionFactor{
		Name:     "nearby_npcs",
		Weight:   0.1,
		Value:    nearbyScore,
		RawScore: nearbyScore,
	})

	return factors
}

// selectBehavior selects the best behavior based on factors
func (be *BehaviorEngine) selectBehavior(context *models.BehaviorContext, factors []models.DecisionFactor, overrideSchedule bool) *models.BehaviorDecision {
	// Calculate weighted scores for each behavior type
	scores := make(map[models.BehaviorType]float64)

	for _, f := range factors {
		score := f.Value * f.Weight
		// Apply factor to relevant behaviors
		switch f.Name {
		case "player_proximity":
			scores[models.BehaviorInteractPlayer] += score * 2
			scores[models.BehaviorWander] += score * 0.5
		case "social_need":
			scores[models.BehaviorInteractNPC] += score * 2
			scores[models.BehaviorSocialize] += score * 1.5
		case "weather_preference":
			if score > 0 {
				scores[models.BehaviorWander] += score
				scores[models.BehaviorFish] += score * 0.5
			} else {
				scores[models.BehaviorRest] += -score
				scores[models.BehaviorShop] += -score * 0.5
			}
		case "mood_influence":
			if score > 0 {
				scores[models.BehaviorSocialize] += score
				scores[models.BehaviorInteractPlayer] += score * 0.5
			} else {
				scores[models.BehaviorRest] += -score
				scores[models.BehaviorIdle] += -score * 0.5
			}
		case "time_of_day":
			if score > 0.5 {
				scores[models.BehaviorWork] += score
			} else if score > 0 {
				scores[models.BehaviorSocialize] += score
			} else {
				scores[models.BehaviorRest] += -score
			}
		case "spontaneity":
			scores[models.BehaviorWander] += score
			scores[models.BehaviorInteractNPC] += score * 0.5
		case "nearby_npcs":
			scores[models.BehaviorInteractNPC] += score
			scores[models.BehaviorSocialize] += score * 0.5
		}
	}

	// Add base scores
	scores[models.BehaviorIdle] += 0.1 // Default low-priority

	// Check schedule unless overriding
	deviatesFromSchedule := false
	if !overrideSchedule && len(context.NPC.Schedule) > 0 {
		// Check if we should follow schedule
		scheduleScore := 1.0 - context.Personality.Spontaneity
		scores[models.BehaviorGoToLocation] += scheduleScore * 0.5

		// Determine if deviating
		currentEntry := be.getCurrentScheduleEntry(context)
		if currentEntry != nil {
			// Reduce spontaneity-based scores if on schedule
			for behavior := range scores {
				if behavior != models.BehaviorGoToLocation && behavior != models.BehaviorWork {
					scores[behavior] *= (1.0 - scheduleScore*0.3)
				}
			}
		}
	}

	// Select highest scoring behavior
	var bestBehavior models.BehaviorType
	bestScore := -math.MaxFloat64
	for behavior, score := range scores {
		if score > bestScore {
			bestScore = score
			bestBehavior = behavior
		}
	}

	// Build decision
	decision := &models.BehaviorDecision{
		NPCID:                context.NPC.ID,
		BehaviorType:         bestBehavior,
		Priority:             int(bestScore * 10),
		Duration:             be.calculateDuration(bestBehavior, context),
		Reason:               be.generateReason(bestBehavior, factors, context),
		Factors:              factors,
		DeviatesFromSchedule: deviatesFromSchedule,
	}

	// Set target based on behavior
	switch bestBehavior {
	case models.BehaviorInteractPlayer:
		decision.TargetPlayer = true
		decision.TargetLocation = &context.PlayerPosition
		deviatesFromSchedule = true
	case models.BehaviorInteractNPC:
		if len(context.NearbyNPCs) > 0 {
			// Pick nearest NPC
			nearest := context.NearbyNPCs[0]
			minDist := math.MaxFloat64
			for _, npc := range context.NearbyNPCs {
				dist := math.Sqrt(float64(
					(npc.Position.X-context.NPC.Position.X)*(npc.Position.X-context.NPC.Position.X) +
						(npc.Position.Y-context.NPC.Position.Y)*(npc.Position.Y-context.NPC.Position.Y),
				))
				if dist < minDist {
					minDist = dist
					nearest = npc
				}
			}
			decision.TargetNPCID = nearest.ID
			decision.TargetLocation = &nearest.Position
			deviatesFromSchedule = true
		}
	case models.BehaviorGoToLocation:
		entry := be.getCurrentScheduleEntry(context)
		if entry != nil {
			decision.TargetLocation = &entry.Position
		}
	case models.BehaviorWander:
		// Random nearby location
		target := be.getWanderTarget(context)
		decision.TargetLocation = &target
		deviatesFromSchedule = true
	}

	decision.DeviatesFromSchedule = deviatesFromSchedule

	return decision
}

// getCurrentScheduleEntry returns the current schedule entry
func (be *BehaviorEngine) getCurrentScheduleEntry(context *models.BehaviorContext) *models.ScheduleEntry {
	for i := len(context.NPC.Schedule) - 1; i >= 0; i-- {
		entry := context.NPC.Schedule[i]
		entryMinutes := entry.Hour*60 + entry.Minute
		currentMinutes := context.Time.Hour*60 + context.Time.Minute
		if currentMinutes >= entryMinutes {
			return &entry
		}
	}
	return nil
}

// calculateDuration calculates duration for a behavior
func (be *BehaviorEngine) calculateDuration(behavior models.BehaviorType, context *models.BehaviorContext) int {
	baseDuration := 30 // 30 game minutes default

	switch behavior {
	case models.BehaviorInteractPlayer:
		return 15 + be.random.Intn(15)
	case models.BehaviorInteractNPC:
		return 20 + be.random.Intn(20)
	case models.BehaviorWork:
		return 60 + be.random.Intn(60)
	case models.BehaviorSocialize:
		return 30 + be.random.Intn(30)
	case models.BehaviorRest:
		return 60
	case models.BehaviorWander:
		return 10 + be.random.Intn(20)
	case models.BehaviorFish:
		return 60 + be.random.Intn(60)
	default:
		return baseDuration
	}
}

// generateReason generates a human-readable reason for the behavior
func (be *BehaviorEngine) generateReason(behavior models.BehaviorType, factors []models.DecisionFactor, context *models.BehaviorContext) string {
	var topFactor string
	topWeight := 0.0
	for _, f := range factors {
		weightedValue := math.Abs(f.Value * f.Weight)
		if weightedValue > topWeight {
			topWeight = weightedValue
			topFactor = f.Name
		}
	}

	reasons := map[models.BehaviorType]map[string]string{
		models.BehaviorInteractPlayer: {
			"player_proximity": "Player is nearby and seems friendly",
			"social_need":      "Feeling social and wants to talk",
			"mood_influence":   "In a good mood and wants company",
			"default":          "Noticed the player and wants to interact",
		},
		models.BehaviorInteractNPC: {
			"social_need":    "Wants to catch up with neighbors",
			"nearby_npcs":    "Friends are nearby",
			"spontaneity":    "Feeling spontaneous",
			"default":        "Wants some company",
		},
		models.BehaviorWork: {
			"time_of_day":    "Work hours - time to be productive",
			"work_ethic":     "Dedicated to getting things done",
			"default":        "Has work to do",
		},
		models.BehaviorSocialize: {
			"mood_influence": "In a good mood and feeling social",
			"social_need":    "Craves social interaction",
			"default":        "Time to be social",
		},
		models.BehaviorRest: {
			"time_of_day":    "Late - time to wind down",
			"mood_influence": "Feeling tired",
			"weather_preference": "Bad weather - staying in",
			"default":        "Taking a break",
		},
		models.BehaviorWander: {
			"weather_preference": "Nice weather - perfect for a walk",
			"spontaneity":        "Feeling spontaneous",
			"default":            "Exploring the area",
		},
		models.BehaviorGoToLocation: {
			"default": "Following usual routine",
		},
	}

	if behaviorReasons, ok := reasons[behavior]; ok {
		if reason, ok := behaviorReasons[topFactor]; ok {
			return reason
		}
		if reason, ok := behaviorReasons["default"]; ok {
			return reason
		}
	}

	return fmt.Sprintf("Decided to %s", behavior)
}

// getWanderTarget returns a random wander target
func (be *BehaviorEngine) getWanderTarget(context *models.BehaviorContext) models.Position {
	// Random offset within 5-15 tiles
	angle := be.random.Float64() * 2 * math.Pi
	distance := 5 + be.random.Float64()*10

	target := models.Position{
		X: context.NPC.Position.X + int(math.Cos(angle)*distance),
		Y: context.NPC.Position.Y + int(math.Sin(angle)*distance),
	}

	// Keep within map bounds
	if target.X < 0 {
		target.X = 0
	}
	if target.Y < 0 {
		target.Y = 0
	}
	if target.X > 47 {
		target.X = 47
	}
	if target.Y > 47 {
		target.Y = 47
	}

	return target
}

// createInteraction creates an NPC interaction
func (be *BehaviorEngine) createInteraction(decision *models.BehaviorDecision, context *models.BehaviorContext) *models.NPCInteraction {
	interaction := &models.NPCInteraction{
		InitiatorID: decision.NPCID,
		Duration:    decision.Duration,
	}

	if decision.TargetPlayer {
		interaction.TargetID = "player"
		interaction.InteractionType = "greet"
		interaction.Dialogue = be.generateGreeting(context)
	} else if decision.TargetNPCID != "" {
		interaction.TargetID = decision.TargetNPCID
		interaction.InteractionType = be.randomInteractionType()
		interaction.Dialogue = be.generateNPCDialogue(context)
	}

	// Log interaction
	be.interactionLog[decision.NPCID] = append(be.interactionLog[decision.NPCID], interaction.InteractionType)
	if len(be.interactionLog[decision.NPCID]) > 10 {
		be.interactionLog[decision.NPCID] = be.interactionLog[decision.NPCID][len(be.interactionLog[decision.NPCID])-10:]
	}

	return interaction
}

// randomInteractionType returns a random interaction type
func (be *BehaviorEngine) randomInteractionType() string {
	types := []string{"talk", "greet", "joke", "share_news", "ask_about_day"}
	return types[be.random.Intn(len(types))]
}

// generateGreeting generates a greeting for the player
func (be *BehaviorEngine) generateGreeting(context *models.BehaviorContext) string {
	greetings := map[models.Mood][]string{
		models.MoodHappy: {
			"Hey there! Great to see you!",
			"Oh, hi! What a lovely day!",
			"Hello, friend! How are you?",
		},
		models.MoodExcited: {
			"Oh! Perfect timing!",
			"Hey! I was just thinking about you!",
			"Hello there! I have something to share!",
		},
		models.MoodNeutral: {
			"Oh, hello.",
			"Hi there.",
			"Good to see you.",
		},
		models.MoodSad: {
			"Oh... hi.",
			"Hello...",
			"Hey...",
		},
		models.MoodTired: {
			"Oh, hello... *yawn*",
			"Hi there... excuse my tiredness.",
		},
	}

	mood := context.MoodState.CurrentMood
	if moodGreeting, ok := greetings[mood]; ok {
		return moodGreeting[be.random.Intn(len(moodGreeting))]
	}
	return "Hello!"
}

// generateNPCDialogue generates dialogue for NPC-NPC interaction
func (be *BehaviorEngine) generateNPCDialogue(context *models.BehaviorContext) string {
	dialogues := []string{
		"Hey, how's your day going?",
		"Did you hear the news?",
		"Nice weather we're having!",
		"I was just thinking about you!",
		"Want to catch up later?",
	}
	return dialogues[be.random.Intn(len(dialogues))]
}

// generateMoodReason generates a reason for the current mood
func (be *BehaviorEngine) generateMoodReason(state *models.NPCMoodState, context *models.BehaviorContext) string {
	if len(state.MoodInfluencers) == 0 {
		return "Just a normal day"
	}

	reasons := map[string]string{
		"sunny_weather":   "The sunshine is uplifting",
		"rainy_weather":   "The rain is a bit gloomy",
		"stormy_weather":  "The storm is unsettling",
		"close_friendship": "Having good friends nearby",
		"good_friendship":  "Feeling connected to others",
		"social_proximity": "Being around people helps",
		"lonely":          "Feeling a bit isolated",
		"morning_energy":  "Starting the day right",
		"evening_energy":  "Enjoying the evening",
		"late_hour_tired": "Getting late",
	}

	inf := state.MoodInfluencers[len(state.MoodInfluencers)-1]
	if reason, ok := reasons[inf]; ok {
		return reason
	}
	return inf
}

// LogInteraction logs an interaction for future reference
func (be *BehaviorEngine) LogInteraction(npcID string, interactionType string) {
	be.interactionLog[npcID] = append(be.interactionLog[npcID], interactionType)
	if len(be.interactionLog[npcID]) > 10 {
		be.interactionLog[npcID] = be.interactionLog[npcID][len(be.interactionLog[npcID])-10:]
	}
}

// BuildContext is a public wrapper for buildContext
func (be *BehaviorEngine) BuildContext(npc *models.NPCState, req *models.AIDecisionRequest, npcs []models.NPCState) *models.BehaviorContext {
	return be.buildContext(npc, req, npcs)
}

// CalculateDecisionFactors is a public wrapper for calculateDecisionFactors
func (be *BehaviorEngine) CalculateDecisionFactors(context *models.BehaviorContext) []models.DecisionFactor {
	return be.calculateDecisionFactors(context)
}
