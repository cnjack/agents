package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"stardew-agent/internal/game"
	"stardew-agent/internal/models"
)

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Client represents a WebSocket client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	movChan     chan models.NPCMovementUpdate // NPC movement updates from engine
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		movChan:     make(chan models.NPCMovementUpdate, 100),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS update
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()

		case movUpdate := <-h.movChan:
			// Broadcast NPC movement update to all clients
			h.broadcastNPCMovement(movUpdate)

		case <-ticker.C:
			// Periodic state broadcasts are handled by BroadcastState
		}
	}
}

// BroadcastState broadcasts the current game state to all clients
func (h *Hub) BroadcastState(state *models.GameState) {
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error marshaling state: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			// Client buffer full, skip
		}
	}
}

// BroadcastMessage broadcasts a custom message to all clients
func (h *Hub) BroadcastMessage(msgType string, data interface{}) {
	message := map[string]interface{}{
		"type": msgType,
		"data": data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- jsonData
}

// broadcastNPCMovement broadcasts NPC movement update to all clients
func (h *Hub) broadcastNPCMovement(mov models.NPCMovementUpdate) {
	message := map[string]interface{}{
		"type": "npc_movement",
		"data": mov,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling NPC movement: %v", err)
		return
	}

	h.mu.Lock()
	h.broadcast <- jsonData
	h.mu.Unlock()
}

// SetMovementChannel sets the movement update channel from engine
func (h *Hub) SetMovementChannel(movChan chan models.NPCMovementUpdate) {
	go func() {
		for movUpdate := range movChan {
			h.broadcastNPCMovement(movUpdate)
		}
	}()
}

// HandleWebSocket handles WebSocket connections

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub, engine *game.Engine, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	// Send initial state
	state := engine.GetState()
	initialState, _ := json.Marshal(map[string]interface{}{
		"type": "state",
		"data": state,
	})
	client.send <- initialState

	// Start read and write goroutines
	go client.writePump()
	go client.readPump(engine)
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump(engine *game.Engine) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse action from message
		var action models.Action
		if err := json.Unmarshal(message, &action); err == nil {
			// Handle action
			c.handleAction(engine, &action)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Batch queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleAction handles an incoming action from WebSocket
func (c *Client) handleAction(engine *game.Engine, action *models.Action) {
	// Process action and send result
	state := engine.GetState()

	// Simple action processing
	result := map[string]interface{}{
		"type":   "action_result",
		"action": action.Type,
		"state":  state,
	}

	data, _ := json.Marshal(result)
	c.send <- data

	// Broadcast updated state to all clients
	c.hub.BroadcastState(engine.GetState())
}
