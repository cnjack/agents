package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"stardew-agent/internal/game"
	"stardew-agent/internal/websocket"
)

// SetupRouter sets up the API router
func SetupRouter(engine *game.Engine, wsHub *websocket.Hub) *gin.Engine {
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
	aiHandler := NewAIHandler(engine)

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

		// NPC endpoints
		v1.GET("/npc/:id/mood", aiHandler.GetNPCMood)
		v1.GET("/npc/mood/:npc_id", handler.GetNPCMood)
		v1.GET("/npc/interactions", handler.GetNPCInteractions)
		v1.POST("/npc/interactions/trigger", handler.TriggerNPCInteraction)

		// WebSocket endpoint
		v1.GET("/ws", func(c *gin.Context) {
			websocket.HandleWebSocket(wsHub, engine, c)
		})
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
