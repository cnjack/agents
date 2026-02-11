# Stardew Community Agent

> An AI Agent Development Platform Built as a Stardew Valley-Style Pixel Art Community Simulation

A pixel art community simulation game designed as an AI Agent development and testing environment. AI agents can control the game character through a RESTful API to interact with NPCs, complete quests, and build friendships in a dynamic community world.

## Features

- **Pixel Art Community Game**: Vue 3 + Canvas 2D rendering with Stardew Valley-inspired visuals
- **AI Agent API**: RESTful endpoints for programmatic game control
- **Real-time Updates**: WebSocket support for live state synchronization
- **Rich Community Systems**:
  - NPC interactions (dialogue, friendship, gifts)
  - Dynamic NPC schedules based on time and relationships
  - Quest system with objectives and rewards
  - Day/night cycle with seasons
  - AI-powered NPC personalities and dialogue generation

## Tech Stack

### Frontend
- **Vue 3** + TypeScript
- **Pinia** - State Management
- **Vite** - Build Tool
- **Canvas 2D** - Game Rendering
- **TailwindCSS** - Styling
- **WebSocket** - Real-time updates

### Backend
- **Go 1.21+**
- **Gin** - Web Framework
- **Gorilla WebSocket** - Real-time communication
- **AI Integration** - Claude API for NPC dialogue and behavior

### Deployment
- Docker + Docker Compose
- Nginx (Reverse Proxy)

## Quick Start

### Prerequisites
- Node.js 20+
- Go 1.21+
- Docker & Docker Compose (optional)

### Development Setup

1. Start the backend:
```bash
cd backend
go mod download
go run cmd/server/main.go
```

2. In a new terminal, start the frontend:
```bash
cd frontend
npm install
npm run dev
```

3. Open http://localhost:5173 in your browser

### Docker Deployment

```bash
docker-compose up --build
```

Access the game at http://localhost

## API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/game/state` | Get complete game state |
| GET | `/game/observation` | Get agent-observable state |
| GET | `/game/map` | Get map data |
| POST | `/game/action` | Execute an action |
| POST | `/game/reset` | Reset game to initial state |

### AI Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/ai/dialogue` | Generate NPC dialogue using AI |
| POST | `/ai/decision` | AI behavior decision making |
| POST | `/ai/interpret` | Natural language command interpretation |
| GET | `/npc/:id/mood` | Get NPC mood state |

## Action Types

### Movement
```json
{ "type": "move", "params": { "direction": "up" } }
```
Directions: `up`, `down`, `left`, `right`

### Interaction
```json
{ "type": "talk", "params": { "npc": "lewis" } }
{ "type": "give_gift", "params": { "npc": "haley", "gift_item": "flower" } }
```

### Quest Management
```json
{ "type": "accept_quest", "params": { "quest_id": "q001" } }
```

### Time Control
```json
{ "type": "wait", "params": { "ticks": 60 } }
{ "type": "sleep", "params": {} }
```

## Game Controls (Human Play)

| Key | Action |
|-----|--------|
| `W` / `↑` | Move up |
| `S` / `↓` | Move down |
| `A` / `←` | Move left |
| `D` / `→` | Move right |
| `E` | Talk to nearby NPC |
| `G` | Give gift to nearby NPC |
| `Space` | Interact |

## Project Structure

```
stardew-agent/
├── frontend/                  # Vue 3 frontend
│   ├── src/
│   │   ├── components/        # Vue components
│   │   ├── game/
│   │   │   ├── engine/        # Core engine (Renderer, GameLoop, Input)
│   │   │   └── entities/      # Game entities (Player, NPC)
│   │   ├── stores/            # Pinia stores
│   │   ├── api/               # API/WebSocket client
│   │   └── types/             # TypeScript types
│   └── ...
├── backend/                   # Go backend
│   ├── cmd/server/            # Entry point
│   └── internal/
│       ├── game/              # Game logic
│       ├── api/               # HTTP handlers
│       ├── models/            # Data models
│       ├── services/ai/       # AI integration services
│       └── websocket/         # WebSocket hub
├── doc/                       # Documentation
│   ├── API.md                 # API documentation
│   ├── ARCHITECTURE.md        # System architecture
│   ├── BACKEND.md             # Backend implementation details
│   ├── FRONTEND.md            # Frontend implementation details
│   ├── GAME_DESIGN.md         # Game design document
│   ├── AI_GAMEPLAY_DESIGN.md  # AI gameplay mechanics
│   ├── NPC_PROFILES.md        # NPC personality profiles
│   └── TEAM_STATE.md          # Development team state
├── docker-compose.yml
├── Dockerfile.frontend
├── Dockerfile.backend
└── nginx.conf
```

## NPC Characters

| NPC | Role | Location | Schedule |
|-----|------|----------|----------|
| Lewis | Mayor / Quest Giver | Community Center | Home → Community Center → Town Square |
| Pierre | Shop Owner | General Store | Home → Shop → Home |
| Robin | Carpenter | Carpenter Shop | Home → Shop → Home |
| Haley | Neighbor | Various | Home → Beach → Saloon |
| Willy | Fisherman | Fish Shop | Shop → Dock → Saloon |

## Time System

- **Day Cycle**: 6:00 AM - 2:00 AM (20 hours)
- **Seasons**: Spring, Summer, Fall, Winter (28 days each)
- **Year**: 4 seasons (112 days)
- **NPC Schedules**: NPCs move dynamically based on time and relationships

## Development

This project was developed using **Claude Teams** - a collaborative AI agent system where specialized AI teammates work together:

- **UI Artist**: Responsible for visual design, pixel art, and UI components
- **Game Designer**: Designs game mechanics, NPC behaviors, and gameplay systems
- **Engineer**: Implements features, fixes bugs, and manages technical architecture
- **QA Tester**: Validates implementations and ensures quality standards

### Development Roadmap

- [x] Core game engine
- [x] NPC system with AI-powered dialogue
- [x] Quest system
- [x] Agent API
- [x] Community-focused gameplay
- [x] Dynamic NPC schedules
- [ ] Enhanced NPC AI behaviors
- [ ] Event system
- [ ] More NPCs and locations
- [ ] Multiplayer support

## Documentation

Comprehensive documentation is available in the `doc/` directory:

- [API Documentation](doc/API.md) - Complete API reference
- [Architecture](doc/ARCHITECTURE.md) - System design and components
- [Backend Guide](doc/BACKEND.md) - Go backend implementation
- [Frontend Guide](doc/FRONTEND.md) - Vue 3 frontend implementation
- [Game Design](doc/GAME_DESIGN.md) - Core game mechanics
- [AI Gameplay Design](doc/AI_GAMEPLAY_DESIGN.md) - AI-driven NPC behaviors
- [NPC Profiles](doc/NPC_PROFILES.md) - Character personalities and schedules

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

---

*Built with Claude Teams - Collaborative AI Development*
