package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"stardew-agent/internal/game"
	aiservice "stardew-agent/internal/services/ai"
	"stardew-agent/internal/websocket"
)

// SetupRouter sets up the API router
func SetupRouter(engine *game.Engine, wsHub *websocket.Hub, aiManager *aiservice.ServiceManager) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Custom middleware
	router.Use(RequestLogger())
	router.Use(Recovery())

	// Create handlers
	handler := NewHandler(engine)

	// Create AI handler with service manager
	var aiHandler *AIHandler
	if aiManager != nil {
		aiHandler = NewAIHandlerWithManager(engine, aiManager)
	} else {
		aiHandler = NewAIHandler(engine)
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Game state endpoints
		v1.GET("/game/state", handler.GetState)
		v1.GET("/game/observation", handler.GetObservation)
		v1.GET("/game/map", handler.GetMap)
		v1.POST("/game/action", handler.ExecuteAction)
		v1.POST("/game/reset", handler.ResetGame)

		// AI behavior decision endpoints
		v1.POST("/ai/decision", handler.GetAIDecision)
		v1.GET("/ai/factors/:npc_id", handler.GetBehaviorFactors)

		// AI dialogue endpoints
		v1.POST("/ai/dialogue", aiHandler.GenerateDialogue)
		v1.POST("/ai/interpret", aiHandler.InterpretNaturalLanguage)
		v1.POST("/ai/execute", aiHandler.ExecuteNLCommand)
		v1.POST("/ai/interact", aiHandler.HandleComplexInteraction)

		// AI provider management
		v1.GET("/ai/providers", aiHandler.ListProviders)
		v1.GET("/ai/providers/:name/health", aiHandler.CheckProviderHealth)

		// NPC endpoints
		v1.GET("/npc/:id/mood", aiHandler.GetNPCMood)
		v1.GET("/npc/mood/:npc_id", handler.GetNPCMood)
		v1.GET("/npc/interactions", handler.GetNPCInteractions)
		v1.POST("/npc/interactions/trigger", handler.TriggerNPCInteraction)

		// Pathfinding endpoints
		v1.POST("/pathfinding/find", handler.FindPath)
		v1.GET("/npc/:id/movement", handler.GetNPCMovement)

		// WebSocket endpoint
		v1.GET("/ws", func(c *gin.Context) {
			websocket.HandleWebSocket(wsHub, engine, c)
		})

	// Connect engine movement channel to hub for broadcasting
	wsHub.SetMovementChannel(engine.MovementUpdateChannel())

	return router
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "stardew-agent",
		})
	})

	return router
}
