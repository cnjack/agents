import type { GameState, TileType, CropStage, Direction, NPCState, RenderNPCState } from '@/types/game'

const TILE_SIZE = 48

// ============== NPC INTERPOLATION STATE ==============
// Client-side interpolated NPC states for smooth animation
const npcRenderStates = new Map<string, RenderNPCState>()

// Get or create render state for an NPC
function getNPCRenderState(npc: NPCState): RenderNPCState {
  const id = npc.id
  if (!npcRenderStates.has(id)) {
    // First time seeing this NPC - initialize at current position
    npcRenderStates.set(id, {
      ...npc,
      renderX: npc.position.x,
      renderY: npc.position.y,
    })
  }

  const renderState = npcRenderStates.get(id)!

  // Check if NPC position changed (server update)
  const posChanged = renderState.position.x !== npc.position.x ||
                     renderState.position.y !== npc.position.y

  if (posChanged) {
    // Store previous position and start moving to new position
    renderState.last_position = { ...renderState.position }
    renderState.position = { ...npc.position }
    renderState.target_position = { ...npc.position }
    renderState.is_moving = true
  }

  return renderState
}

// Update NPC interpolation (call every frame)
function updateNPCInterpolation(deltaTime: number) {
  const lerpFactor = 10 * deltaTime // Adjust for movement speed

  for (const state of npcRenderStates.values()) {
    if (state.is_moving && state.target_position) {
      // Interpolate towards target
      const dx = state.target_position.x - state.renderX
      const dy = state.target_position.y - state.renderY

      state.renderX += dx * Math.min(lerpFactor, 1)
      state.renderY += dy * Math.min(lerpFactor, 1)

      // Check if we've reached the target (with small tolerance)
      const dist = Math.sqrt(dx * dx + dy * dy)
      if (dist < 0.05) {
        state.renderX = state.target_position.x
        state.renderY = state.target_position.y
        state.is_moving = false
      }
    }
  }
}

// Remove NPC render state (when NPC is removed from game)
function removeNPCRenderState(id: string) {
  npcRenderStates.delete(id)
}

// ============== PIXEL ART ASSETS (Base64 encoded) ==============
// These are simple pixel art patterns encoded as data URLs for the tiles

// Create a canvas for each tile type
const tileCache = new Map<string, HTMLCanvasElement>()

function getTileCanvas(type: string, variant = 0): HTMLCanvasElement {
  const key = `${type}_${variant}`
  if (tileCache.has(key)) {
    return tileCache.get(key)!
  }

  const canvas = document.createElement('canvas')
  canvas.width = TILE_SIZE
  canvas.height = TILE_SIZE
  const ctx = canvas.getContext('2d')!
  ctx.imageSmoothingEnabled = false

  // Generate tile based on type
  generateTile(ctx, type, variant)

  tileCache.set(key, canvas)
  return canvas
}

function generateTile(ctx: CanvasRenderingContext2D, type: string, variant: number) {
  switch (type) {
    case 'grass':
      drawGrassTile(ctx, variant)
      break
    case 'dirt':
      drawDirtTile(ctx, variant)
      break
    case 'tilled':
      drawTilledTile(ctx, variant)
      break
    case 'watered':
      drawWateredTile(ctx, variant)
      break
    case 'path':
      drawPathTile(ctx, variant)
      break
    case 'water':
      drawWaterTile(ctx, variant)
      break
    default:
      ctx.fillStyle = '#4a7c23'
      ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)
  }
}

// Grass tile with enhanced pixel art detail
function drawGrassTile(ctx: CanvasRenderingContext2D, variant: number) {
  // Base green with gradient effect
  const baseGreens = ['#4a8c2a', '#3d7a24', '#5a9c3a', '#448c28']
  ctx.fillStyle = baseGreens[variant % baseGreens.length]
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Subtle vertical gradient overlay
  const gradient = ctx.createLinearGradient(0, 0, 0, TILE_SIZE)
  gradient.addColorStop(0, 'rgba(255, 255, 255, 0.08)')
  gradient.addColorStop(0.5, 'rgba(0, 0, 0, 0)')
  gradient.addColorStop(1, 'rgba(0, 0, 0, 0.1)')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Darker texture patches (more detailed pattern)
  ctx.fillStyle = '#2d5a18'
  const pattern = [
    [0,0,1,0,0,1,0,0,1,0,0,1],
    [1,0,0,0,1,0,0,0,1,0,0,0],
    [0,0,1,0,0,0,1,0,0,0,1,0],
    [0,1,0,0,0,0,0,1,0,1,0,0],
    [1,0,0,1,0,1,0,0,0,0,1,0],
    [0,0,0,0,1,0,0,1,0,0,0,1],
  ]
  for (let y = 0; y < 6; y++) {
    for (let x = 0; x < 12; x++) {
      if (pattern[y][x]) {
        ctx.fillRect(x * 4, y * 8, 4, 4)
      }
    }
  }

  // Highlight patches
  ctx.fillStyle = '#7ac04a'
  for (let i = 0; i < 8; i++) {
    const x = ((i * 17 + variant * 7) % (TILE_SIZE - 4))
    const y = ((i * 11 + variant * 5) % (TILE_SIZE - 4))
    ctx.fillRect(x, y, 3, 3)
  }

  // Grass blades with varied heights and colors
  const bladeColors = ['#6aac3a', '#5a9c2a', '#7abc4a', '#4a8c1a']
  const blades = [
    // Top row
    [3,2,8],[9,4,6],[15,3,7],[21,5,5],[27,2,9],[33,4,6],[39,3,7],[45,5,5],
    // Middle row
    [1,18,7],[7,20,5],[13,19,8],[19,21,6],[25,18,7],[31,20,5],[37,19,8],[43,21,6],
    // Bottom row
    [5,34,6],[11,36,4],[17,35,7],[23,33,5],[29,34,6],[35,36,4],[41,35,7],[47,33,5]
  ]
  for (let i = 0; i < blades.length; i++) {
    const [x, y, h] = blades[i]
    if ((variant + i) % 3 !== 0) {
      ctx.fillStyle = bladeColors[i % bladeColors.length]
      // Draw grass blade with slight curve effect
      ctx.fillRect(x, y, 2, h - 2)
      ctx.fillRect(x + 1, y + h - 2, 1, 2)
    }
  }

  // Small clover/flower details
  if (variant % 6 === 0) {
    const flowerColors = ['#ff6b6b', '#ffd93d', '#6bcb77', '#4d96ff', '#ff9ff3', '#ffffff']
    ctx.fillStyle = flowerColors[variant % flowerColors.length]
    ctx.beginPath()
    ctx.arc(24, 24, 5, 0, Math.PI * 2)
    ctx.fill()
    // Flower center
    ctx.fillStyle = '#ffeb3b'
    ctx.beginPath()
    ctx.arc(24, 24, 2, 0, Math.PI * 2)
    ctx.fill()
    // Petal details
    ctx.fillStyle = flowerColors[variant % flowerColors.length]
    for (let p = 0; p < 5; p++) {
      const angle = (p / 5) * Math.PI * 2
      ctx.beginPath()
      ctx.arc(24 + Math.cos(angle) * 6, 24 + Math.sin(angle) * 6, 3, 0, Math.PI * 2)
      ctx.fill()
    }
  }

  // Occasional mushroom
  if (variant % 10 === 1) {
    ctx.fillStyle = '#d2691e'
    ctx.fillRect(38, 38, 4, 6)
    ctx.fillStyle = '#ff6347'
    ctx.beginPath()
    ctx.arc(40, 36, 5, Math.PI, 0)
    ctx.fill()
    ctx.fillStyle = '#fff'
    ctx.fillRect(38, 35, 2, 2)
    ctx.fillRect(42, 34, 2, 1)
  }
}

// Dirt/farm tile with enhanced detail
function drawDirtTile(ctx: CanvasRenderingContext2D, variant: number) {
  const browns = ['#8b6914', '#7a5a10', '#9b7920', '#6a4a0c']
  ctx.fillStyle = browns[variant % browns.length]
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Subtle gradient
  const gradient = ctx.createLinearGradient(0, 0, TILE_SIZE, TILE_SIZE)
  gradient.addColorStop(0, 'rgba(255, 220, 180, 0.1)')
  gradient.addColorStop(1, 'rgba(0, 0, 0, 0.15)')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Dark soil patches
  ctx.fillStyle = '#5a3a00'
  for (let i = 0; i < 18; i++) {
    const x = ((i * 11) + variant * 3) % (TILE_SIZE - 3)
    const y = ((i * 7) + variant * 2) % (TILE_SIZE - 3)
    ctx.fillRect(x, y, 4, 4)
  }

  // Small pebbles
  ctx.fillStyle = '#7a6a50'
  for (let i = 0; i < 6; i++) {
    const x = ((i * 19) + variant * 5) % (TILE_SIZE - 4)
    const y = ((i * 13) + variant * 3) % (TILE_SIZE - 4)
    ctx.beginPath()
    ctx.ellipse(x + 2, y + 2, 3, 2, (i * 0.5), 0, Math.PI * 2)
    ctx.fill()
  }

  // Lighter soil highlights
  ctx.fillStyle = '#b08040'
  for (let i = 0; i < 8; i++) {
    const x = ((i * 23 + variant * 7) % (TILE_SIZE - 2))
    const y = ((i * 17 + variant * 4) % (TILE_SIZE - 2))
    ctx.fillRect(x, y, 2, 2)
  }

  // Occasional small root/twig
  if (variant % 5 === 0) {
    ctx.strokeStyle = '#4a3020'
    ctx.lineWidth = 2
    ctx.beginPath()
    ctx.moveTo(8, 8)
    ctx.lineTo(14, 18)
    ctx.lineTo(20, 14)
    ctx.stroke()
  }
}

// Tilled soil with enhanced furrow detail
function drawTilledTile(ctx: CanvasRenderingContext2D, variant: number) {
  // Darker base for tilled soil
  ctx.fillStyle = '#5a4010'
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Subtle horizontal gradient for depth
  const gradient = ctx.createLinearGradient(0, 0, 0, TILE_SIZE)
  gradient.addColorStop(0, 'rgba(0, 0, 0, 0.1)')
  gradient.addColorStop(0.5, 'rgba(0, 0, 0, 0)')
  gradient.addColorStop(1, 'rgba(0, 0, 0, 0.15)')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Furrows with 3D effect
  for (let i = 0; i < 6; i++) {
    const y = i * 8 + 2
    // Dark shadow
    ctx.fillStyle = '#3a2800'
    ctx.fillRect(0, y, TILE_SIZE, 3)
    // Highlight ridge
    ctx.fillStyle = '#7a6030'
    ctx.fillRect(0, y - 2, TILE_SIZE, 2)
  }

  // Soil clumps with shadow
  ctx.fillStyle = '#6a5020'
  const clumps = [[4,4],[20,10],[12,18],[32,24],[8,34],[28,38],[40,6]]
  for (const [x, y] of clumps) {
    ctx.fillRect(x, y, 5, 5)
    // Highlight on clump
    ctx.fillStyle = '#8a7040'
    ctx.fillRect(x, y, 2, 2)
    ctx.fillStyle = '#6a5020'
  }

  // Small stones in soil
  ctx.fillStyle = '#908060'
  for (let i = 0; i < 4; i++) {
    const x = ((i * 31 + variant * 11) % (TILE_SIZE - 4))
    const y = ((i * 19 + variant * 7) % (TILE_SIZE - 4))
    ctx.fillRect(x, y, 3, 2)
  }
}

// Watered soil with enhanced wet appearance
function drawWateredTile(ctx: CanvasRenderingContext2D, variant: number) {
  // Dark wet base
  ctx.fillStyle = '#3a2a10'
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Wet sheen gradient
  const gradient = ctx.createLinearGradient(0, 0, TILE_SIZE, TILE_SIZE)
  gradient.addColorStop(0, 'rgba(100, 140, 200, 0.15)')
  gradient.addColorStop(0.5, 'rgba(60, 100, 160, 0.08)')
  gradient.addColorStop(1, 'rgba(100, 140, 200, 0.12)')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Wet shine spots
  ctx.fillStyle = 'rgba(140, 180, 220, 0.3)'
  for (let i = 0; i < 5; i++) {
    const x = (i * 19 + variant * 7) % (TILE_SIZE - 10)
    const y = (i * 13 + variant * 5) % (TILE_SIZE - 6)
    ctx.fillRect(x, y, 10, 4)
    // Brighter center
    ctx.fillStyle = 'rgba(180, 210, 240, 0.25)'
    ctx.fillRect(x + 2, y + 1, 6, 2)
    ctx.fillStyle = 'rgba(140, 180, 220, 0.3)'
  }

  // Darker furrows with wet depth
  for (let i = 0; i < 6; i++) {
    const y = i * 8 + 2
    ctx.fillStyle = '#1a0a00'
    ctx.fillRect(0, y, TILE_SIZE, 3)
    // Wet highlight on ridge
    ctx.fillStyle = 'rgba(100, 130, 180, 0.2)'
    ctx.fillRect(0, y - 2, TILE_SIZE, 2)
  }

  // Occasional water droplet reflection
  if (variant % 3 === 0) {
    ctx.fillStyle = 'rgba(200, 230, 255, 0.4)'
    ctx.beginPath()
    ctx.arc(24, 20, 3, 0, Math.PI * 2)
    ctx.fill()
    ctx.fillStyle = 'rgba(255, 255, 255, 0.5)'
    ctx.beginPath()
    ctx.arc(23, 19, 1, 0, Math.PI * 2)
    ctx.fill()
  }
}

// Path tile with enhanced stone and wear detail
function drawPathTile(ctx: CanvasRenderingContext2D, variant: number) {
  // Base path color
  ctx.fillStyle = '#a08060'
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Subtle texture gradient
  const gradient = ctx.createRadialGradient(TILE_SIZE/2, TILE_SIZE/2, 0, TILE_SIZE/2, TILE_SIZE/2, TILE_SIZE)
  gradient.addColorStop(0, 'rgba(255, 255, 255, 0.05)')
  gradient.addColorStop(1, 'rgba(0, 0, 0, 0.1)')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Larger stones with shadow
  ctx.fillStyle = '#806040'
  for (let i = 0; i < 12; i++) {
    const x = ((i * 17) + variant * 3) % (TILE_SIZE - 5)
    const y = ((i * 11) + variant * 2) % (TILE_SIZE - 5)
    // Shadow
    ctx.fillStyle = '#604020'
    ctx.fillRect(x + 1, y + 1, 5, 4)
    // Stone
    ctx.fillStyle = '#9a7a5a'
    ctx.fillRect(x, y, 5, 4)
    // Highlight
    ctx.fillStyle = '#b09070'
    ctx.fillRect(x, y, 2, 2)
  }

  // Worn path center (dirt showing through)
  ctx.fillStyle = '#b09070'
  ctx.fillRect(12, 18, 24, 12)
  // Worn edge texture
  ctx.fillStyle = '#c0a080'
  ctx.fillRect(14, 20, 20, 8)

  // Small pebbles scattered
  ctx.fillStyle = '#706050'
  for (let i = 0; i < 8; i++) {
    const x = ((i * 23 + variant * 5) % (TILE_SIZE - 3))
    const y = ((i * 19 + variant * 3) % (TILE_SIZE - 3))
    ctx.beginPath()
    ctx.arc(x + 1, y + 1, 2, 0, Math.PI * 2)
    ctx.fill()
  }

  // Occasional footprint mark
  if (variant % 4 === 0) {
    ctx.fillStyle = 'rgba(0, 0, 0, 0.15)'
    ctx.beginPath()
    ctx.ellipse(28, 32, 4, 6, 0.2, 0, Math.PI * 2)
    ctx.fill()
  }
}

// Water tile with enhanced animated appearance
function drawWaterTile(ctx: CanvasRenderingContext2D, variant: number) {
  // Base water with gradient depth
  const baseGradient = ctx.createLinearGradient(0, 0, TILE_SIZE, TILE_SIZE)
  baseGradient.addColorStop(0, '#4a90c2')
  baseGradient.addColorStop(0.5, '#3a80b2')
  baseGradient.addColorStop(1, '#2a70a2')
  ctx.fillStyle = baseGradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Depth gradient overlay
  const depthGradient = ctx.createRadialGradient(TILE_SIZE/2, TILE_SIZE/2, 5, TILE_SIZE/2, TILE_SIZE/2, TILE_SIZE)
  depthGradient.addColorStop(0, 'rgba(0, 40, 80, 0.2)')
  depthGradient.addColorStop(1, 'rgba(0, 20, 60, 0.3)')
  ctx.fillStyle = depthGradient
  ctx.fillRect(0, 0, TILE_SIZE, TILE_SIZE)

  // Wave patterns with foam
  const waveOffset = (variant * 3) % 20
  ctx.fillStyle = 'rgba(120, 180, 220, 0.4)'
  // Wave 1
  ctx.beginPath()
  ctx.moveTo(4 + waveOffset, 10)
  ctx.quadraticCurveTo(12 + waveOffset, 6, 20 + waveOffset, 10)
  ctx.quadraticCurveTo(28 + waveOffset, 14, 36 + waveOffset, 10)
  ctx.lineTo(36 + waveOffset, 12)
  ctx.quadraticCurveTo(28 + waveOffset, 16, 20 + waveOffset, 12)
  ctx.quadraticCurveTo(12 + waveOffset, 8, 4 + waveOffset, 12)
  ctx.closePath()
  ctx.fill()

  // Wave 2
  ctx.fillStyle = 'rgba(140, 200, 240, 0.3)'
  ctx.beginPath()
  ctx.moveTo(8 - waveOffset/2, 26)
  ctx.quadraticCurveTo(16 - waveOffset/2, 22, 24 - waveOffset/2, 26)
  ctx.quadraticCurveTo(32 - waveOffset/2, 30, 40 - waveOffset/2, 26)
  ctx.lineTo(40 - waveOffset/2, 28)
  ctx.quadraticCurveTo(32 - waveOffset/2, 32, 24 - waveOffset/2, 28)
  ctx.quadraticCurveTo(16 - waveOffset/2, 24, 8 - waveOffset/2, 28)
  ctx.closePath()
  ctx.fill()

  // Wave 3
  ctx.fillStyle = 'rgba(100, 160, 200, 0.35)'
  ctx.beginPath()
  ctx.moveTo(2 + waveOffset/3, 40)
  ctx.quadraticCurveTo(14 + waveOffset/3, 36, 26 + waveOffset/3, 40)
  ctx.quadraticCurveTo(38 + waveOffset/3, 44, 46 + waveOffset/3, 40)
  ctx.lineTo(46 + waveOffset/3, 42)
  ctx.quadraticCurveTo(38 + waveOffset/3, 46, 26 + waveOffset/3, 42)
  ctx.quadraticCurveTo(14 + waveOffset/3, 38, 2 + waveOffset/3, 42)
  ctx.closePath()
  ctx.fill()

  // Sparkle highlights
  if (variant % 3 === 0) {
    ctx.fillStyle = 'rgba(255, 255, 255, 0.6)'
    ctx.beginPath()
    ctx.arc(32, 16, 3, 0, Math.PI * 2)
    ctx.fill()
    ctx.fillStyle = 'rgba(255, 255, 255, 0.9)'
    ctx.beginPath()
    ctx.arc(31, 15, 1.5, 0, Math.PI * 2)
    ctx.fill()
  }
  if (variant % 5 === 1) {
    ctx.fillStyle = 'rgba(255, 255, 255, 0.5)'
    ctx.beginPath()
    ctx.arc(14, 36, 2, 0, Math.PI * 2)
    ctx.fill()
  }
}

// ============== CHARACTER SPRITES ==============
const characterCache = new Map<string, HTMLCanvasElement>()

function getCharacterSprite(type: string, id: string, direction: Direction = 'down'): HTMLCanvasElement {
  const key = `${type}_${id}_${direction}`
  if (characterCache.has(key)) {
    return characterCache.get(key)!
  }

  const canvas = document.createElement('canvas')
  canvas.width = TILE_SIZE * 2
  canvas.height = TILE_SIZE * 2
  const ctx = canvas.getContext('2d')!
  ctx.imageSmoothingEnabled = false

  if (type === 'player') {
    drawPlayerSprite(ctx, direction)
  } else {
    drawNPCSprite(ctx, id, direction)
  }

  characterCache.set(key, canvas)
  return canvas
}

// Player sprite - Stardew Valley inspired style
// Key: 16x32 equivalent proportions, simple colors, minimal outlines
function drawPlayerSprite(ctx: CanvasRenderingContext2D, direction: Direction) {
  const centerX = TILE_SIZE
  const centerY = TILE_SIZE
  const s = TILE_SIZE / 32 // Scale factor

  // Color palette (Stardew-style: limited, clear colors)
  const COLORS = {
    skin: '#f5cba7',
    skinShadow: '#d4a574',
    hair: '#5d4037',
    hairShadow: '#3e2723',
    shirt: '#c0392b',
    shirtShadow: '#922b21',
    pants: '#2980b9',
    pantsShadow: '#1a5276',
    boots: '#5d4037',
    bootsShadow: '#3e2723',
    outline: '#1a1a1a'
  }

  // ===== SHADOW (simple oval) =====
  ctx.fillStyle = 'rgba(0, 0, 0, 0.2)'
  ctx.beginPath()
  ctx.ellipse(centerX, centerY + 28 * s, 10 * s, 5 * s, 0, 0, Math.PI * 2)
  ctx.fill()

  // ===== BOOTS =====
  ctx.fillStyle = COLORS.boots
  ctx.fillRect(centerX - 6 * s, centerY + 18 * s, 5 * s, 8 * s)
  ctx.fillRect(centerX + 1 * s, centerY + 18 * s, 5 * s, 8 * s)
  // Boot shadow
  ctx.fillStyle = COLORS.bootsShadow
  ctx.fillRect(centerX - 5 * s, centerY + 24 * s, 4 * s, 2 * s)
  ctx.fillRect(centerX + 2 * s, centerY + 24 * s, 4 * s, 2 * s)

  // ===== PANTS/LEGS =====
  ctx.fillStyle = COLORS.pants
  ctx.fillRect(centerX - 5 * s, centerY + 8 * s, 4 * s, 12 * s)
  ctx.fillRect(centerX + 1 * s, centerY + 8 * s, 4 * s, 12 * s)
  // Pants shadow
  ctx.fillStyle = COLORS.pantsShadow
  ctx.fillRect(centerX + 2 * s, centerY + 8 * s, 3 * s, 12 * s)

  // ===== BODY/TORSO =====
  ctx.fillStyle = COLORS.shirt
  ctx.fillRect(centerX - 7 * s, centerY - 4 * s, 14 * s, 14 * s)
  // Body shadow
  ctx.fillStyle = COLORS.shirtShadow
  ctx.fillRect(centerX + 3 * s, centerY - 4 * s, 5 * s, 14 * s)

  // ===== ARMS =====
  ctx.fillStyle = COLORS.shirt
  if (direction === 'up') {
    // Arms spread when walking away
    ctx.fillRect(centerX - 11 * s, centerY - 2 * s, 4 * s, 10 * s)
    ctx.fillRect(centerX + 7 * s, centerY - 2 * s, 4 * s, 10 * s)
  } else if (direction === 'down' || direction === 'left' || direction === 'right') {
    // Arms at side
    ctx.fillRect(centerX - 10 * s, centerY - 2 * s, 4 * s, 10 * s)
    ctx.fillRect(centerX + 6 * s, centerY - 2 * s, 4 * s, 10 * s)
  }
  // Arm shadow
  ctx.fillStyle = COLORS.shirtShadow
  ctx.fillRect(centerX + 7 * s, centerY - 2 * s, 3 * s, 10 * s)

  // ===== HANDS =====
  ctx.fillStyle = COLORS.skin
  if (direction === 'up') {
    ctx.fillRect(centerX - 10 * s, centerY + 6 * s, 3 * s, 3 * s)
    ctx.fillRect(centerX + 7 * s, centerY + 6 * s, 3 * s, 3 * s)
  } else {
    ctx.fillRect(centerX - 9 * s, centerY + 6 * s, 3 * s, 3 * s)
    ctx.fillRect(centerX + 6 * s, centerY + 6 * s, 3 * s, 3 * s)
  }

  // ===== HEAD (proportional, not oversized) =====
  // Head base - oval shape
  ctx.fillStyle = COLORS.skin
  ctx.beginPath()
  ctx.ellipse(centerX, centerY - 14 * s, 7 * s, 8 * s, 0, 0, Math.PI * 2)
  ctx.fill()
  // Head shadow
  ctx.fillStyle = COLORS.skinShadow
  ctx.beginPath()
  ctx.ellipse(centerX + 2 * s, centerY - 12 * s, 5 * s, 6 * s, 0.3, Math.PI * 0.5, Math.PI * 1.5)
  ctx.fill()

  // ===== HAIR =====
  ctx.fillStyle = COLORS.hair
  // Hair top
  ctx.beginPath()
  ctx.ellipse(centerX, centerY - 20 * s, 7 * s, 5 * s, 0, Math.PI, 0)
  ctx.fill()
  // Hair sides
  ctx.fillRect(centerX - 8 * s, centerY - 20 * s, 3 * s, 8 * s)
  ctx.fillRect(centerX + 5 * s, centerY - 20 * s, 3 * s, 8 * s)
  // Hair shadow
  ctx.fillStyle = COLORS.hairShadow
  ctx.fillRect(centerX + 5 * s, centerY - 18 * s, 2 * s, 6 * s)

  // ===== EYES (Stardew style: small, 2-3 pixel dots) =====
  if (direction === 'up') {
    // Back of head - no eyes
  } else if (direction === 'down') {
    // Small black dots for eyes (Stardew style)
    ctx.fillStyle = '#1a1a1a'
    ctx.fillRect(centerX - 4 * s, centerY - 16 * s, 2 * s, 2 * s)
    ctx.fillRect(centerX + 2 * s, centerY - 16 * s, 2 * s, 2 * s)
  } else if (direction === 'left') {
    ctx.fillStyle = '#1a1a1a'
    ctx.fillRect(centerX - 6 * s, centerY - 16 * s, 2 * s, 2 * s)
  } else {
    ctx.fillStyle = '#1a1a1a'
    ctx.fillRect(centerX + 4 * s, centerY - 16 * s, 2 * s, 2 * s)
  }

  // ===== OUTLINE (thin, only where needed for contrast) =====
  ctx.fillStyle = COLORS.outline
  // Head outline - bottom
  ctx.fillRect(centerX - 7 * s, centerY - 8 * s, 14 * s, 1 * s)
  // Body outline - sides
  ctx.fillRect(centerX - 8 * s, centerY - 4 * s, 1 * s, 18 * s)
  ctx.fillRect(centerX + 7 * s, centerY - 4 * s, 1 * s, 18 * s)
}

// NPC sprites - Stardew Valley inspired style
function drawNPCSprite(ctx: CanvasRenderingContext2D, id: string, direction: Direction) {
  const centerX = TILE_SIZE
  const centerY = TILE_SIZE
  const s = TILE_SIZE / 32

  // NPC color palettes (Stardew-style: 2 main colors + shadow)
  const NPC_PALETTES: Record<string, {
    outfit: string; outfitShadow: string;
    hair: string; hairShadow: string;
    skin: string; skinShadow: string;
    accessory: string
  }> = {
    lewis: {
      outfit: '#2e7d32', outfitShadow: '#1b5e20',
      hair: '#616161', hairShadow: '#424242',
      skin: '#f5cba7', skinShadow: '#d4a574',
      accessory: 'sash'
    },
    pierre: {
      outfit: '#1565c0', outfitShadow: '#0d47a1',
      hair: '#5d4037', hairShadow: '#3e2723',
      skin: '#f5cba7', skinShadow: '#d4a574',
      accessory: 'apron'
    },
    robin: {
      outfit: '#e65100', outfitShadow: '#bf360c',
      hair: '#c62828', hairShadow: '#8e0000',
      skin: '#f5cba7', skinShadow: '#d4a574',
      accessory: 'toolbelt'
    },
    haley: {
      outfit: '#ec407a', outfitShadow: '#c2185b',
      hair: '#fdd835', hairShadow: '#f9a825',
      skin: '#ffe0b2', skinShadow: '#ffcc80',
      accessory: 'bow'
    },
    willy: {
      outfit: '#1565c0', outfitShadow: '#0d47a1',
      hair: '#37474f', hairShadow: '#263238',
      skin: '#ffcc80', skinShadow: '#ffb74d',
      accessory: 'hat'
    }
  }

  const p = NPC_PALETTES[id] || {
    outfit: '#7b1fa2', outfitShadow: '#4a148c',
    hair: '#424242', hairShadow: '#212121',
    skin: '#f5cba7', skinShadow: '#d4a574',
    accessory: 'none'
  }

  // ===== SHADOW =====
  ctx.fillStyle = 'rgba(0, 0, 0, 0.2)'
  ctx.beginPath()
  ctx.ellipse(centerX, centerY + 28 * s, 10 * s, 5 * s, 0, 0, Math.PI * 2)
  ctx.fill()

  // ===== BOOTS =====
  ctx.fillStyle = '#4a3728'
  ctx.fillRect(centerX - 6 * s, centerY + 18 * s, 5 * s, 8 * s)
  ctx.fillRect(centerX + 1 * s, centerY + 18 * s, 5 * s, 8 * s)
  ctx.fillStyle = '#3e2723'
  ctx.fillRect(centerX - 5 * s, centerY + 24 * s, 4 * s, 2 * s)
  ctx.fillRect(centerX + 2 * s, centerY + 24 * s, 4 * s, 2 * s)

  // ===== PANTS/LEGS =====
  ctx.fillStyle = p.outfitShadow
  ctx.fillRect(centerX - 5 * s, centerY + 8 * s, 4 * s, 12 * s)
  ctx.fillRect(centerX + 1 * s, centerY + 8 * s, 4 * s, 12 * s)

  // ===== BODY =====
  ctx.fillStyle = p.outfit
  ctx.fillRect(centerX - 7 * s, centerY - 4 * s, 14 * s, 14 * s)
  ctx.fillStyle = p.outfitShadow
  ctx.fillRect(centerX + 3 * s, centerY - 4 * s, 5 * s, 14 * s)

  // ===== ARMS =====
  ctx.fillStyle = p.outfit
  ctx.fillRect(centerX - 10 * s, centerY - 2 * s, 4 * s, 10 * s)
  ctx.fillRect(centerX + 6 * s, centerY - 2 * s, 4 * s, 10 * s)
  ctx.fillStyle = p.outfitShadow
  ctx.fillRect(centerX + 7 * s, centerY - 2 * s, 3 * s, 10 * s)

  // ===== HANDS =====
  ctx.fillStyle = p.skin
  ctx.fillRect(centerX - 9 * s, centerY + 6 * s, 3 * s, 3 * s)
  ctx.fillRect(centerX + 6 * s, centerY + 6 * s, 3 * s, 3 * s)

  // ===== HEAD =====
  ctx.fillStyle = p.skin
  ctx.beginPath()
  ctx.ellipse(centerX, centerY - 14 * s, 7 * s, 8 * s, 0, 0, Math.PI * 2)
  ctx.fill()
  ctx.fillStyle = p.skinShadow
  ctx.beginPath()
  ctx.ellipse(centerX + 2 * s, centerY - 12 * s, 5 * s, 6 * s, 0.3, Math.PI * 0.5, Math.PI * 1.5)
  ctx.fill()

  // ===== HAIR =====
  ctx.fillStyle = p.hair
  ctx.beginPath()
  ctx.ellipse(centerX, centerY - 20 * s, 7 * s, 5 * s, 0, Math.PI, 0)
  ctx.fill()
  ctx.fillRect(centerX - 8 * s, centerY - 20 * s, 3 * s, 8 * s)
  ctx.fillRect(centerX + 5 * s, centerY - 20 * s, 3 * s, 8 * s)
  ctx.fillStyle = p.hairShadow
  ctx.fillRect(centerX + 5 * s, centerY - 18 * s, 2 * s, 6 * s)

  // ===== EYES (Stardew style: small dots) =====
  if (direction !== 'up') {
    ctx.fillStyle = '#1a1a1a'
    if (direction === 'down') {
      ctx.fillRect(centerX - 4 * s, centerY - 16 * s, 2 * s, 2 * s)
      ctx.fillRect(centerX + 2 * s, centerY - 16 * s, 2 * s, 2 * s)
    } else if (direction === 'left') {
      ctx.fillRect(centerX - 6 * s, centerY - 16 * s, 2 * s, 2 * s)
    } else {
      ctx.fillRect(centerX + 4 * s, centerY - 16 * s, 2 * s, 2 * s)
    }
  }

  // ===== ACCESSORIES =====
  if (p.accessory === 'sash') {
    // Mayor sash
    ctx.fillStyle = '#ffd700'
    ctx.fillRect(centerX + 4 * s, centerY - 2 * s, 3 * s, 8 * s)
  } else if (p.accessory === 'apron') {
    // Shop apron
    ctx.fillStyle = '#f5f5f5'
    ctx.fillRect(centerX - 5 * s, centerY - 2 * s, 10 * s, 8 * s)
  } else if (p.accessory === 'toolbelt') {
    // Tool belt
    ctx.fillStyle = '#5d4037'
    ctx.fillRect(centerX - 6 * s, centerY + 2 * s, 12 * s, 4 * s)
  } else if (p.accessory === 'bow') {
    // Hair bow
    ctx.fillStyle = '#ff69b4'
    ctx.beginPath()
    ctx.arc(centerX - 10 * s, centerY - 22 * s, 3 * s, 0, Math.PI * 2)
    ctx.fill()
    ctx.beginPath()
    ctx.arc(centerX - 5 * s, centerY - 22 * s, 3 * s, 0, Math.PI * 2)
    ctx.fill()
  } else if (p.accessory === 'hat') {
    // Fisherman hat
    ctx.fillStyle = '#37474f'
    ctx.beginPath()
    ctx.ellipse(centerX, centerY - 24 * s, 10 * s, 4 * s, 0, 0, Math.PI * 2)
    ctx.fill()
    ctx.fillStyle = '#ffd700'
    ctx.fillRect(centerX - 8 * s, centerY - 26 * s, 16 * s, 2 * s)
  }

  // ===== SUBTLE OUTLINE =====
  ctx.fillStyle = '#1a1a1a'
  ctx.fillRect(centerX - 7 * s, centerY - 4 * s, 14 * s, 1 * s)  // Body top
  ctx.fillRect(centerX - 8 * s, centerY - 4 * s, 1 * s, 18 * s)  // Left side
  ctx.fillRect(centerX + 7 * s, centerY - 4 * s, 1 * s, 18 * s)  // Right side
}

// ============== BUILDING SPRITES ==============
const buildingCache = new Map<string, HTMLCanvasElement>()

function getBuildingSprite(type: string, width: number, height: number): HTMLCanvasElement {
  const key = `building_${type}_${width}_${height}`
  if (buildingCache.has(key)) {
    return buildingCache.get(key)!
  }

  const canvas = document.createElement('canvas')
  canvas.width = width * TILE_SIZE
  canvas.height = height * TILE_SIZE
  const ctx = canvas.getContext('2d')!
  ctx.imageSmoothingEnabled = false

  drawBuilding(ctx, type, width, height)

  buildingCache.set(key, canvas)
  return canvas
}

function drawBuilding(ctx: CanvasRenderingContext2D, type: string, width: number, height: number) {
  const w = width * TILE_SIZE
  const h = height * TILE_SIZE
  const scale = TILE_SIZE / 32

  // Building color schemes with more detail colors
  const colors: Record<string, { wall: string; wallDark: string; wallLight: string; roof: string; roofDark: string; door: string; trim: string; accent: string }> = {
    community_center: { wall: '#d4a574', wallDark: '#b8956a', wallLight: '#e8c0a0', roof: '#8b4513', roofDark: '#6b3510', door: '#654321', trim: '#ffd700', accent: '#fff8dc' },
    shop: { wall: '#b8a090', wallDark: '#9a8a7a', wallLight: '#d0c0b0', roof: '#2f4f4f', roofDark: '#1f3f3f', door: '#4a3728', trim: '#cd853f', accent: '#90ee90' },
    carpenter: { wall: '#d2691e', wallDark: '#b25a18', wallLight: '#e89050', roof: '#556b2f', roofDark: '#3d4f1f', door: '#4a3728', trim: '#deb887', accent: '#8b4513' },
    fish_shop: { wall: '#87ceeb', wallDark: '#6ab8d8', wallLight: '#a8e0f8', roof: '#191970', roofDark: '#101050', door: '#4169e1', trim: '#4682b4', accent: '#ffd700' },
    saloon: { wall: '#8b4513', wallDark: '#6b3510', wallLight: '#a86030', roof: '#2f4f4f', roofDark: '#1f2f2f', door: '#654321', trim: '#ffd700', accent: '#ff6347' },
    house: { wall: '#deb887', wallDark: '#c4a070', wallLight: '#f0d0a0', roof: '#8b0000', roofDark: '#6b0000', door: '#654321', trim: '#8b4513', accent: '#87ceeb' },
  }

  const c = colors[type] || colors.house

  // Enhanced shadow
  ctx.fillStyle = 'rgba(0, 0, 0, 0.2)'
  ctx.fillRect(12, h - 6, w, 10)
  ctx.fillStyle = 'rgba(0, 0, 0, 0.1)'
  ctx.fillRect(6, h - 4, w + 12, 6)

  // Wall base with gradient
  const wallGradient = ctx.createLinearGradient(0, 0, w, h)
  wallGradient.addColorStop(0, c.wallLight)
  wallGradient.addColorStop(0.5, c.wall)
  wallGradient.addColorStop(1, c.wallDark)
  ctx.fillStyle = wallGradient
  ctx.fillRect(0, 0, w, h)

  // Wall texture - brick/wood pattern
  ctx.fillStyle = 'rgba(0, 0, 0, 0.08)'
  for (let y = 0; y < h; y += 16 * scale) {
    const offset = (Math.floor(y / (16 * scale)) % 2) * 24 * scale
    for (let x = -offset; x < w; x += 48 * scale) {
      ctx.fillRect(x, y, 46 * scale, 14 * scale)
      ctx.strokeStyle = 'rgba(0, 0, 0, 0.05)'
      ctx.lineWidth = 1
      ctx.strokeRect(x, y, 46 * scale, 14 * scale)
    }
  }

  // Wall highlight strip
  ctx.fillStyle = 'rgba(255, 255, 255, 0.1)'
  ctx.fillRect(0, 0, w, 8 * scale)

  // Foundation
  ctx.fillStyle = c.wallDark
  ctx.fillRect(0, h - 12 * scale, w, 12 * scale)
  ctx.fillStyle = 'rgba(0, 0, 0, 0.15)'
  ctx.fillRect(0, h - 12 * scale, w, 3 * scale)

  // Enhanced roof with 3D effect
  const roofHeight = h / 3

  // Roof shadow
  ctx.fillStyle = 'rgba(0, 0, 0, 0.15)'
  ctx.beginPath()
  ctx.moveTo(-12, 4 * scale)
  ctx.lineTo(w / 2, -roofHeight + 8 * scale)
  ctx.lineTo(w + 12, 4 * scale)
  ctx.closePath()
  ctx.fill()

  // Main roof
  const roofGradient = ctx.createLinearGradient(0, -roofHeight, w, 0)
  roofGradient.addColorStop(0, c.roof)
  roofGradient.addColorStop(0.5, c.roofDark)
  roofGradient.addColorStop(1, c.roof)
  ctx.fillStyle = roofGradient
  ctx.beginPath()
  ctx.moveTo(-10, 0)
  ctx.lineTo(w / 2, -roofHeight)
  ctx.lineTo(w + 10, 0)
  ctx.closePath()
  ctx.fill()

  // Roof ridge line
  ctx.strokeStyle = c.roofDark
  ctx.lineWidth = 4 * scale
  ctx.beginPath()
  ctx.moveTo(w / 2, -roofHeight)
  ctx.lineTo(w / 2, 0)
  ctx.stroke()

  // Roof tiles pattern
  ctx.fillStyle = 'rgba(0, 0, 0, 0.05)'
  for (let row = 0; row < 6; row++) {
    const y = -roofHeight + row * (roofHeight / 6) + roofHeight / 12
    const rowWidth = w * (1 - row / 8)
    const startX = (w - rowWidth) / 2
    ctx.fillRect(startX, y, rowWidth, 4 * scale)
  }

  // Roof trim/fascia
  ctx.fillStyle = c.trim
  ctx.fillRect(-6, -2 * scale, w + 12, 8 * scale)

  // Decorative trim under roof
  ctx.fillStyle = c.accent
  ctx.fillRect(0, 6 * scale, w, 3 * scale)

  // Enhanced door with frame
  const doorWidth = Math.min(28 * scale, w - 32 * scale)
  const doorHeight = Math.min(40 * scale, h - 24 * scale)
  const doorX = w / 2 - doorWidth / 2
  const doorY = h - doorHeight - 8 * scale

  // Door frame
  ctx.fillStyle = c.trim
  ctx.fillRect(doorX - 4 * scale, doorY - 4 * scale, doorWidth + 8 * scale, doorHeight + 8 * scale)

  // Door with wood grain effect
  const doorGradient = ctx.createLinearGradient(doorX, doorY, doorX + doorWidth, doorY)
  doorGradient.addColorStop(0, c.door)
  doorGradient.addColorStop(0.3, '#5a4232')
  doorGradient.addColorStop(0.7, c.door)
  doorGradient.addColorStop(1, '#4a3222')
  ctx.fillStyle = doorGradient
  ctx.fillRect(doorX, doorY, doorWidth, doorHeight)

  // Door panels
  ctx.fillStyle = 'rgba(0, 0, 0, 0.1)'
  ctx.fillRect(doorX + 4 * scale, doorY + 6 * scale, doorWidth - 8 * scale, doorHeight / 2 - 10 * scale)
  ctx.fillRect(doorX + 4 * scale, doorY + doorHeight / 2 + 4 * scale, doorWidth - 8 * scale, doorHeight / 2 - 10 * scale)

  // Door knob with shine
  ctx.fillStyle = c.trim
  ctx.beginPath()
  ctx.arc(doorX + doorWidth - 8 * scale, doorY + doorHeight / 2, 4 * scale, 0, Math.PI * 2)
  ctx.fill()
  ctx.fillStyle = 'rgba(255, 255, 255, 0.5)'
  ctx.beginPath()
  ctx.arc(doorX + doorWidth - 9 * scale, doorY + doorHeight / 2 - 1 * scale, 2 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Enhanced windows
  const windowWidth = 16 * scale
  const windowHeight = 20 * scale
  const windowY = h / 2 - 8 * scale

  const drawWindow = (wx: number, wy: number) => {
    // Window frame
    ctx.fillStyle = c.trim
    ctx.fillRect(wx - 3 * scale, wy - 3 * scale, windowWidth + 6 * scale, windowHeight + 6 * scale)

    // Window glass with gradient
    const glassGradient = ctx.createLinearGradient(wx, wy, wx + windowWidth, wy + windowHeight)
    glassGradient.addColorStop(0, '#a8d8ea')
    glassGradient.addColorStop(0.5, '#87ceeb')
    glassGradient.addColorStop(1, '#6bb8d8')
    ctx.fillStyle = glassGradient
    ctx.fillRect(wx, wy, windowWidth, windowHeight)

    // Window reflection
    ctx.fillStyle = 'rgba(255, 255, 255, 0.3)'
    ctx.fillRect(wx + 2 * scale, wy + 2 * scale, windowWidth / 3, windowHeight / 3)

    // Window cross bars
    ctx.fillStyle = c.trim
    ctx.fillRect(wx + windowWidth / 2 - 1 * scale, wy, 2 * scale, windowHeight)
    ctx.fillRect(wx, wy + windowHeight / 2 - 1 * scale, windowWidth, 2 * scale)

    // Window sill
    ctx.fillStyle = c.wallLight
    ctx.fillRect(wx - 4 * scale, wy + windowHeight, windowWidth + 8 * scale, 4 * scale)
  }

  if (w > 80 * scale) {
    // Two windows
    drawWindow(16 * scale, windowY)
    drawWindow(w - 16 * scale - windowWidth, windowY)
  } else {
    // One centered window
    drawWindow(w / 2 - windowWidth / 2, windowY)
  }

  // Building-specific decorations
  if (type === 'community_center') {
    // Banner/sign
    ctx.fillStyle = '#ffd700'
    ctx.fillRect(w / 2 - 30 * scale, 16 * scale, 60 * scale, 16 * scale)
    ctx.fillStyle = '#8b4513'
    ctx.font = `${10 * scale}px Arial`
    ctx.textAlign = 'center'
    ctx.fillText('COMMUNITY', w / 2, 26 * scale)
  } else if (type === 'shop') {
    // Awning
    ctx.fillStyle = '#ff6347'
    ctx.fillRect(8 * scale, 14 * scale, w - 16 * scale, 12 * scale)
    // Awning stripes
    ctx.fillStyle = '#fff'
    for (let i = 0; i < 4; i++) {
      ctx.fillRect(12 * scale + i * (w - 24 * scale) / 4, 14 * scale, (w - 24 * scale) / 8, 12 * scale)
    }
  } else if (type === 'fish_shop') {
    // Fish sign
    ctx.fillStyle = '#4682b4'
    ctx.beginPath()
    ctx.ellipse(w / 2, 20 * scale, 16 * scale, 8 * scale, 0, 0, Math.PI * 2)
    ctx.fill()
    ctx.fillStyle = '#ffd700'
    ctx.beginPath()
    ctx.moveTo(w / 2 + 12 * scale, 20 * scale)
    ctx.lineTo(w / 2 + 24 * scale, 12 * scale)
    ctx.lineTo(w / 2 + 24 * scale, 28 * scale)
    ctx.closePath()
    ctx.fill()
  }
}

// ============== TREE SPRITE ==============
const treeCache = new Map<string, HTMLCanvasElement>()

function getTreeSprite(variant: number): HTMLCanvasElement {
  const key = `tree_${variant}`
  if (treeCache.has(key)) {
    return treeCache.get(key)!
  }

  const canvas = document.createElement('canvas')
  canvas.width = TILE_SIZE * 2
  canvas.height = TILE_SIZE * 3
  const ctx = canvas.getContext('2d')!
  ctx.imageSmoothingEnabled = false

  drawTree(ctx, variant)

  treeCache.set(key, canvas)
  return canvas
}

function drawTree(ctx: CanvasRenderingContext2D, variant: number) {
  const centerX = TILE_SIZE
  const baseY = TILE_SIZE * 3
  const scale = TILE_SIZE / 32

  // Enhanced shadow
  ctx.fillStyle = 'rgba(0, 0, 0, 0.2)'
  ctx.beginPath()
  ctx.ellipse(centerX + 8, baseY - 4, 22 * scale, 10 * scale, 0, 0, Math.PI * 2)
  ctx.fill()

  // Roots
  ctx.fillStyle = '#4a2a10'
  ctx.fillRect(centerX - 16 * scale, baseY - 12 * scale, 8 * scale, 12 * scale)
  ctx.fillRect(centerX + 8 * scale, baseY - 10 * scale, 8 * scale, 10 * scale)

  // Main trunk with gradient
  const trunkGradient = ctx.createLinearGradient(centerX - 12 * scale, baseY - 48 * scale, centerX + 12 * scale, baseY)
  trunkGradient.addColorStop(0, '#6a4a2a')
  trunkGradient.addColorStop(0.3, '#5a3a1a')
  trunkGradient.addColorStop(0.7, '#4a2a10')
  trunkGradient.addColorStop(1, '#3a1a00')
  ctx.fillStyle = trunkGradient
  ctx.fillRect(centerX - 12 * scale, baseY - 48 * scale, 24 * scale, 48 * scale)

  // Trunk bark texture
  ctx.fillStyle = '#3a1a00'
  for (let i = 0; i < 8; i++) {
    const x = centerX - 10 * scale + (i * 6) % 20 * scale
    const y = baseY - 44 * scale + i * 5 * scale
    ctx.fillRect(x, y, 4 * scale, 6 * scale)
  }

  // Trunk highlights
  ctx.fillStyle = 'rgba(255, 255, 255, 0.1)'
  ctx.fillRect(centerX - 10 * scale, baseY - 46 * scale, 6 * scale, 42 * scale)

  // Branch stubs
  ctx.fillStyle = '#5a3a1a'
  ctx.fillRect(centerX - 18 * scale, baseY - 36 * scale, 8 * scale, 6 * scale)
  ctx.fillRect(centerX + 10 * scale, baseY - 42 * scale, 10 * scale, 6 * scale)

  // Foliage layers with enhanced shading
  const greens = ['#1a5a1a', '#1a6b1a', '#228b22', '#32cd32', '#90ee90']
  const greenLights = ['#2a7a2a', '#2a8b2a', '#3a9b32', '#4adb42', '#b0ffb0']

  // Bottom foliage layer (darkest)
  ctx.fillStyle = greens[0]
  ctx.beginPath()
  ctx.arc(centerX, baseY - 56 * scale, 38 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Bottom layer highlights
  ctx.fillStyle = greenLights[0]
  ctx.beginPath()
  ctx.arc(centerX - 16 * scale, baseY - 52 * scale, 20 * scale, 0, Math.PI * 2)
  ctx.fill()
  ctx.beginPath()
  ctx.arc(centerX + 18 * scale, baseY - 50 * scale, 18 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Second layer
  ctx.fillStyle = greens[1]
  ctx.beginPath()
  ctx.arc(centerX, baseY - 68 * scale, 32 * scale, 0, Math.PI * 2)
  ctx.fill()

  ctx.fillStyle = greenLights[1]
  ctx.beginPath()
  ctx.arc(centerX - 14 * scale, baseY - 64 * scale, 16 * scale, 0, Math.PI * 2)
  ctx.fill()
  ctx.beginPath()
  ctx.arc(centerX + 12 * scale, baseY - 66 * scale, 14 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Third layer
  ctx.fillStyle = greens[2]
  ctx.beginPath()
  ctx.arc(centerX, baseY - 80 * scale, 26 * scale, 0, Math.PI * 2)
  ctx.fill()

  ctx.fillStyle = greenLights[2]
  ctx.beginPath()
  ctx.arc(centerX - 10 * scale, baseY - 78 * scale, 14 * scale, 0, Math.PI * 2)
  ctx.fill()
  ctx.beginPath()
  ctx.arc(centerX + 8 * scale, baseY - 82 * scale, 12 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Fourth layer
  ctx.fillStyle = greens[3]
  ctx.beginPath()
  ctx.arc(centerX, baseY - 92 * scale, 20 * scale, 0, Math.PI * 2)
  ctx.fill()

  ctx.fillStyle = greenLights[3]
  ctx.beginPath()
  ctx.arc(centerX - 6 * scale, baseY - 92 * scale, 10 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Top layer (brightest)
  ctx.fillStyle = greens[4]
  ctx.beginPath()
  ctx.arc(centerX, baseY - 102 * scale, 14 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Specular highlight
  ctx.fillStyle = 'rgba(255, 255, 255, 0.3)'
  ctx.beginPath()
  ctx.arc(centerX - 4 * scale, baseY - 104 * scale, 5 * scale, 0, Math.PI * 2)
  ctx.fill()

  // Add some leaf clusters for depth
  ctx.fillStyle = 'rgba(0, 50, 0, 0.3)'
  for (let i = 0; i < 6; i++) {
    const angle = (i / 6) * Math.PI * 2 + variant * 0.3
    const dist = 30 * scale + variant * 2
    const x = centerX + Math.cos(angle) * dist
    const y = baseY - 70 * scale + Math.sin(angle) * 20 * scale
    ctx.beginPath()
    ctx.arc(x, y, 8 * scale, 0, Math.PI * 2)
    ctx.fill()
  }

  // Occasional fruit/apples
  if (variant % 4 === 0) {
    ctx.fillStyle = '#ff4444'
    ctx.beginPath()
    ctx.arc(centerX - 20 * scale, baseY - 60 * scale, 5 * scale, 0, Math.PI * 2)
    ctx.fill()
    ctx.beginPath()
    ctx.arc(centerX + 24 * scale, baseY - 74 * scale, 5 * scale, 0, Math.PI * 2)
    ctx.fill()
    // Apple shine
    ctx.fillStyle = 'rgba(255, 255, 255, 0.5)'
    ctx.beginPath()
    ctx.arc(centerX - 21 * scale, baseY - 62 * scale, 2 * scale, 0, Math.PI * 2)
    ctx.fill()
    ctx.beginPath()
    ctx.arc(centerX + 23 * scale, baseY - 76 * scale, 2 * scale, 0, Math.PI * 2)
    ctx.fill()
  }
}

// ============== MAIN RENDERER ==============
export class Renderer {
  private canvas: HTMLCanvasElement
  private ctx: CanvasRenderingContext2D
  private cameraX = 0
  private cameraY = 0
  private targetCameraX = 0
  private targetCameraY = 0
  private frame = 0
  private hoveredAgentId: string | null = null
  private lastTime = performance.now()

  // Map data
  // Building positions must match backend world.go setupBuildings()
  private readonly MAP_DATA = {
    buildings: [
      { x: 18, y: 4, w: 12, h: 8, type: 'community_center' },
      { x: 4, y: 16, w: 8, h: 6, type: 'shop' },
      { x: 6, y: 4, w: 8, h: 6, type: 'carpenter' },
      { x: 2, y: 36, w: 8, h: 6, type: 'fish_shop' },
      { x: 36, y: 20, w: 8, h: 6, type: 'saloon' },
      // Fixed: Moved houses to avoid collision with community_center
      // haley_house: (28, 6) -> (2, 28)
      // house_2: (32, 12) -> (40, 6)
      { x: 2, y: 28, w: 6, h: 5, type: 'house' },
      { x: 40, y: 6, w: 6, h: 5, type: 'house' },
    ],
    trees: [
      { x: 1, y: 1 }, { x: 1, y: 10 }, { x: 1, y: 20 }, { x: 1, y: 28 },
      { x: 44, y: 1 }, { x: 44, y: 15 }, { x: 44, y: 25 }, { x: 44, y: 35 },
      { x: 46, y: 8 }, { x: 46, y: 20 }, { x: 46, y: 32 },
    ],
    paths: [
      { x: 14, y: 24, w: 20, h: 1 },
      { x: 24, y: 12, w: 1, h: 12 },
      { x: 8, y: 22, w: 1, h: 14 },
      { x: 10, y: 22, w: 4, h: 1 },
    ],
    water: [
      { x: 38, y: 2, w: 8, h: 6 },
    ],
  }

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas
    this.ctx = canvas.getContext('2d')!
    this.ctx.imageSmoothingEnabled = false

    // Pre-generate tile variations
    for (let v = 0; v < 16; v++) {
      getTileCanvas('grass', v)
      getTileCanvas('dirt', v)
      getTileCanvas('path', v)
    }

    // Pre-generate sprites
    for (const dir of ['up', 'down', 'left', 'right'] as Direction[]) {
      getCharacterSprite('player', 'main', dir)
    }
  }

  resize(width: number, height: number) {
    this.canvas.width = width
    this.canvas.height = height
  }

  render(state: GameState) {
    const { ctx, canvas } = this
    this.frame++

    // Calculate delta time for smooth interpolation
    const currentTime = performance.now()
    const deltaTime = Math.min((currentTime - this.lastTime) / 1000, 0.1) // Cap at 100ms
    this.lastTime = currentTime

    // Update NPC interpolation
    updateNPCInterpolation(deltaTime)

    // Camera follows player
    this.targetCameraX = state.player.position.x * TILE_SIZE - canvas.width / 2
    this.targetCameraY = state.player.position.y * TILE_SIZE - canvas.height / 2

    this.cameraX += (this.targetCameraX - this.cameraX) * 0.15
    this.cameraY += (this.targetCameraY - this.cameraY) * 0.15

    // Clamp camera
    const maxX = state.map_width * TILE_SIZE - canvas.width
    const maxY = state.map_height * TILE_SIZE - canvas.height
    this.cameraX = Math.max(0, Math.min(maxX, this.cameraX))
    this.cameraY = Math.max(0, Math.min(maxY, this.cameraY))

    ctx.save()
    ctx.translate(-Math.floor(this.cameraX), -Math.floor(this.cameraY))

    // Sky background
    ctx.fillStyle = this.getSkyColor(state.time.hour)
    ctx.fillRect(Math.floor(this.cameraX), Math.floor(this.cameraY), canvas.width, canvas.height)

    // Calculate visible range
    const startTileX = Math.max(0, Math.floor(this.cameraX / TILE_SIZE) - 2)
    const startTileY = Math.max(0, Math.floor(this.cameraY / TILE_SIZE) - 2)
    const endTileX = Math.min(state.map_width, startTileX + Math.ceil(canvas.width / TILE_SIZE) + 4)
    const endTileY = Math.min(state.map_height, startTileY + Math.ceil(canvas.height / TILE_SIZE) + 4)

    // Draw base tiles
    this.drawBaseTiles(ctx, state, startTileX, startTileY, endTileX, endTileY)

    // Draw paths
    this.drawPaths(ctx)

    // Draw water
    this.drawWater(ctx)

    // Draw farm tiles and crops
    this.drawFarmArea(ctx, state, startTileX, startTileY, endTileX, endTileY)

    // Draw buildings
    this.drawBuildings(ctx)

    // Draw trees
    this.drawTrees(ctx)

    // Draw NPCs
    for (const npc of state.npcs) {
      const renderState = getNPCRenderState(npc)
      // Use interpolated position for smooth animation
      this.drawCharacter(renderState.renderX, renderState.renderY, 'npc', npc.id, npc.direction)
      // Draw highlight if hovered
      if (this.hoveredAgentId === npc.id) {
        this.drawHighlight(ctx, renderState.renderX, renderState.renderY)
      }
    }

    // Draw player
    this.drawCharacter(
      state.player.position.x,
      state.player.position.y,
      'player',
      'main',
      state.player.direction
    )

    ctx.restore()

    // Draw UI overlay
    this.drawUI(ctx, state)
  }

  private drawBaseTiles(
    ctx: CanvasRenderingContext2D,
    state: GameState,
    startX: number,
    startY: number,
    endX: number,
    endY: number
  ) {
    const farmOffsetX = Math.floor((state.map_width - (state.farm?.width || 0)) / 2)
    const farmOffsetY = Math.floor((state.map_height - (state.farm?.height || 0)) / 2)

    for (let y = startY; y < endY; y++) {
      for (let x = startX; x < endX; x++) {
        const px = x * TILE_SIZE
        const py = y * TILE_SIZE

        // Check if in farm area
        const farmX = x - farmOffsetX
        const farmY = y - farmOffsetY
        const isFarm = state.farm && farmX >= 0 && farmX < state.farm.width && farmY >= 0 && farmY < state.farm.height

        if (!isFarm) {
          // Draw grass with variation
          const variant = (x * 13 + y * 7) % 16
          const tileCanvas = getTileCanvas('grass', variant)
          ctx.drawImage(tileCanvas, px, py)
        }
      }
    }
  }

  private drawFarmArea(
    ctx: CanvasRenderingContext2D,
    state: GameState,
    startX: number,
    startY: number,
    endX: number,
    endY: number
  ) {
    if (!state.farm) return

    const farmOffsetX = Math.floor((state.map_width - state.farm.width) / 2)
    const farmOffsetY = Math.floor((state.map_height - state.farm.height) / 2)

    for (let y = startY; y < endY; y++) {
      for (let x = startX; x < endX; x++) {
        const farmX = x - farmOffsetX
        const farmY = y - farmOffsetY

        if (farmX >= 0 && farmX < state.farm.width && farmY >= 0 && farmY < state.farm.height) {
          const tile = state.farm.tiles[farmY]?.[farmX]
          if (tile) {
            const px = x * TILE_SIZE
            const py = y * TILE_SIZE
            const variant = (x * 7 + y * 11) % 8

            const tileCanvas = getTileCanvas(tile.type, variant)
            ctx.drawImage(tileCanvas, px, py)

            // Draw crop if present
            if (tile.crop) {
              this.drawCrop(ctx, px, py, tile.crop.stage)
            }
          }
        }
      }
    }
  }

  private drawCrop(ctx: CanvasRenderingContext2D, px: number, py: number, stage: CropStage) {
    const cx = px + TILE_SIZE / 2
    const cy = py + TILE_SIZE / 2
    const scale = TILE_SIZE / 32

    ctx.save()

    switch (stage) {
      case 'seed':
        // Small mound of dirt with seed visible
        ctx.fillStyle = '#6b4423'
        ctx.beginPath()
        ctx.ellipse(cx, cy + 12 * scale, 10 * scale, 5 * scale, 0, 0, Math.PI * 2)
        ctx.fill()
        // Darker shadow side
        ctx.fillStyle = '#5a3820'
        ctx.beginPath()
        ctx.ellipse(cx + 2 * scale, cy + 13 * scale, 8 * scale, 4 * scale, 0.2, 0, Math.PI * 2)
        ctx.fill()
        // Seed visible on top
        ctx.fillStyle = '#8b6914'
        ctx.beginPath()
        ctx.ellipse(cx, cy + 8 * scale, 4 * scale, 3 * scale, 0, 0, Math.PI * 2)
        ctx.fill()
        // Seed highlight
        ctx.fillStyle = '#a08030'
        ctx.fillRect(cx - 1 * scale, cy + 6 * scale, 2 * scale, 2 * scale)
        break

      case 'sprout':
        // Small green stem with first leaves
        // Stem shadow
        ctx.fillStyle = 'rgba(0, 50, 0, 0.2)'
        ctx.fillRect(cx - 1 * scale, cy + 2 * scale, 3 * scale, 12 * scale)
        // Main stem
        ctx.strokeStyle = '#60c060'
        ctx.lineWidth = 4 * scale
        ctx.lineCap = 'round'
        ctx.beginPath()
        ctx.moveTo(cx, cy + 12 * scale)
        ctx.lineTo(cx, cy - 2 * scale)
        ctx.stroke()
        // Stem highlight
        ctx.strokeStyle = '#80e080'
        ctx.lineWidth = 2 * scale
        ctx.beginPath()
        ctx.moveTo(cx - 1 * scale, cy + 10 * scale)
        ctx.lineTo(cx - 1 * scale, cy)
        ctx.stroke()
        // First leaves (cotyledons)
        ctx.fillStyle = '#70d070'
        ctx.beginPath()
        ctx.ellipse(cx - 7 * scale, cy + 2 * scale, 5 * scale, 3 * scale, -0.5, 0, Math.PI * 2)
        ctx.fill()
        ctx.beginPath()
        ctx.ellipse(cx + 7 * scale, cy + 2 * scale, 5 * scale, 3 * scale, 0.5, 0, Math.PI * 2)
        ctx.fill()
        // Leaf highlights
        ctx.fillStyle = '#90f090'
        ctx.beginPath()
        ctx.ellipse(cx - 8 * scale, cy, 3 * scale, 2 * scale, -0.5, 0, Math.PI * 2)
        ctx.fill()
        ctx.beginPath()
        ctx.ellipse(cx + 6 * scale, cy, 3 * scale, 2 * scale, 0.5, 0, Math.PI * 2)
        ctx.fill()
        break

      case 'growing':
        // Taller plant with more leaves
        // Stem shadow
        ctx.fillStyle = 'rgba(0, 50, 0, 0.15)'
        ctx.fillRect(cx + 2 * scale, cy, 3 * scale, 18 * scale)
        // Main stem
        ctx.strokeStyle = '#228b22'
        ctx.lineWidth = 5 * scale
        ctx.lineCap = 'round'
        ctx.beginPath()
        ctx.moveTo(cx, cy + 14 * scale)
        ctx.lineTo(cx, cy - 8 * scale)
        ctx.stroke()
        // Stem highlight
        ctx.strokeStyle = '#32a032'
        ctx.lineWidth = 2 * scale
        ctx.beginPath()
        ctx.moveTo(cx - 1 * scale, cy + 12 * scale)
        ctx.lineTo(cx - 1 * scale, cy - 6 * scale)
        ctx.stroke()
        // Large leaves
        ctx.fillStyle = '#32cd32'
        ctx.beginPath()
        ctx.ellipse(cx - 10 * scale, cy - 2 * scale, 7 * scale, 10 * scale, -0.3, 0, Math.PI * 2)
        ctx.fill()
        ctx.beginPath()
        ctx.ellipse(cx + 10 * scale, cy - 2 * scale, 7 * scale, 10 * scale, 0.3, 0, Math.PI * 2)
        ctx.fill()
        // Leaf midribs
        ctx.strokeStyle = '#28a428'
        ctx.lineWidth = 1 * scale
        ctx.beginPath()
        ctx.moveTo(cx - 10 * scale, cy + 4 * scale)
        ctx.lineTo(cx - 14 * scale, cy - 6 * scale)
        ctx.moveTo(cx + 10 * scale, cy + 4 * scale)
        ctx.lineTo(cx + 14 * scale, cy - 6 * scale)
        ctx.stroke()
        // Leaf highlights
        ctx.fillStyle = '#50e050'
        ctx.beginPath()
        ctx.ellipse(cx - 12 * scale, cy - 4 * scale, 4 * scale, 5 * scale, -0.3, 0, Math.PI * 2)
        ctx.fill()
        ctx.beginPath()
        ctx.ellipse(cx + 8 * scale, cy - 4 * scale, 4 * scale, 5 * scale, 0.3, 0, Math.PI * 2)
        ctx.fill()
        break

      case 'mature':
        // Full grown crop with produce
        // Stem shadow
        ctx.fillStyle = 'rgba(0, 50, 0, 0.15)'
        ctx.fillRect(cx + 2 * scale, cy - 2 * scale, 3 * scale, 20 * scale)
        // Main stem
        ctx.strokeStyle = '#1a6b1a'
        ctx.lineWidth = 6 * scale
        ctx.lineCap = 'round'
        ctx.beginPath()
        ctx.moveTo(cx, cy + 14 * scale)
        ctx.lineTo(cx, cy - 10 * scale)
        ctx.stroke()
        // Stem highlight
        ctx.strokeStyle = '#2a8b2a'
        ctx.lineWidth = 2 * scale
        ctx.beginPath()
        ctx.moveTo(cx - 1 * scale, cy + 12 * scale)
        ctx.lineTo(cx - 1 * scale, cy - 8 * scale)
        ctx.stroke()
        // Large mature leaves
        ctx.fillStyle = '#228b22'
        ctx.beginPath()
        ctx.ellipse(cx - 12 * scale, cy - 4 * scale, 9 * scale, 12 * scale, -0.4, 0, Math.PI * 2)
        ctx.fill()
        ctx.beginPath()
        ctx.ellipse(cx + 12 * scale, cy - 4 * scale, 9 * scale, 12 * scale, 0.4, 0, Math.PI * 2)
        ctx.fill()
        // Leaf texture
        ctx.strokeStyle = '#1a5a1a'
        ctx.lineWidth = 1 * scale
        ctx.beginPath()
        ctx.moveTo(cx - 12 * scale, cy + 4 * scale)
        ctx.lineTo(cx - 16 * scale, cy - 10 * scale)
        ctx.moveTo(cx + 12 * scale, cy + 4 * scale)
        ctx.lineTo(cx + 16 * scale, cy - 10 * scale)
        ctx.stroke()
        // Leaf highlights
        ctx.fillStyle = '#32a032'
        ctx.beginPath()
        ctx.ellipse(cx - 14 * scale, cy - 6 * scale, 5 * scale, 6 * scale, -0.4, 0, Math.PI * 2)
        ctx.fill()
        ctx.beginPath()
        ctx.ellipse(cx + 10 * scale, cy - 6 * scale, 5 * scale, 6 * scale, 0.4, 0, Math.PI * 2)
        ctx.fill()
        // Fruit/produce with gradient
        const fruitGradient = ctx.createRadialGradient(cx - 2 * scale, cy - 8 * scale, 0, cx, cy - 4 * scale, 10 * scale)
        fruitGradient.addColorStop(0, '#ffe033')
        fruitGradient.addColorStop(0.5, '#ffd700')
        fruitGradient.addColorStop(1, '#cc9900')
        ctx.fillStyle = fruitGradient
        ctx.beginPath()
        ctx.arc(cx, cy - 4 * scale, 10 * scale, 0, Math.PI * 2)
        ctx.fill()
        // Fruit shadow
        ctx.fillStyle = 'rgba(0, 0, 0, 0.15)'
        ctx.beginPath()
        ctx.ellipse(cx + 2 * scale, cy + 2 * scale, 8 * scale, 4 * scale, 0, 0, Math.PI * 2)
        ctx.fill()
        // Fruit shine
        ctx.fillStyle = 'rgba(255, 255, 255, 0.6)'
        ctx.beginPath()
        ctx.arc(cx - 4 * scale, cy - 8 * scale, 4 * scale, 0, Math.PI * 2)
        ctx.fill()
        ctx.fillStyle = 'rgba(255, 255, 255, 0.4)'
        ctx.beginPath()
        ctx.arc(cx - 3 * scale, cy - 10 * scale, 2 * scale, 0, Math.PI * 2)
        ctx.fill()
        // Small stem on fruit
        ctx.fillStyle = '#4a3020'
        ctx.fillRect(cx - 2 * scale, cy - 14 * scale, 4 * scale, 4 * scale)
        break
    }

    ctx.restore()
  }

  private drawPaths(ctx: CanvasRenderingContext2D) {
    for (const path of this.MAP_DATA.paths) {
      const px = path.x * TILE_SIZE
      const py = path.y * TILE_SIZE
      const variant = (path.x + path.y) % 8

      if (path.h > 1) {
        for (let i = 0; i < path.h; i++) {
          const tileCanvas = getTileCanvas('path', (variant + i) % 8)
          ctx.drawImage(tileCanvas, px, py + i * TILE_SIZE)
        }
      } else {
        for (let i = 0; i < path.w; i++) {
          const tileCanvas = getTileCanvas('path', (variant + i) % 8)
          ctx.drawImage(tileCanvas, px + i * TILE_SIZE, py)
        }
      }
    }
  }

  private drawWater(ctx: CanvasRenderingContext2D) {
    for (const water of this.MAP_DATA.water) {
      for (let dy = 0; dy < water.h; dy++) {
        for (let dx = 0; dx < water.w; dx++) {
          const px = (water.x + dx) * TILE_SIZE
          const py = (water.y + dy) * TILE_SIZE
          const tileCanvas = getTileCanvas('water', (this.frame + dx + dy) % 8)
          ctx.drawImage(tileCanvas, px, py)
        }
      }
    }
  }

  private drawBuildings(ctx: CanvasRenderingContext2D) {
    for (const b of this.MAP_DATA.buildings) {
      const sprite = getBuildingSprite(b.type, b.w, b.h)
      ctx.drawImage(sprite, b.x * TILE_SIZE, b.y * TILE_SIZE)
    }
  }

  private drawTrees(ctx: CanvasRenderingContext2D) {
    for (const tree of this.MAP_DATA.trees) {
      const variant = (tree.x + tree.y) % 4
      const sprite = getTreeSprite(variant)
      ctx.drawImage(sprite, tree.x * TILE_SIZE - TILE_SIZE / 2, tree.y * TILE_SIZE - TILE_SIZE)
    }
  }

  private drawCharacter(x: number, y: number, type: string, id: string, direction: Direction) {
    const sprite = getCharacterSprite(type, id, direction)
    this.ctx.drawImage(sprite, x * TILE_SIZE - TILE_SIZE / 2, y * TILE_SIZE - TILE_SIZE / 2)
  }

  private getSkyColor(hour: number): string {
    if (hour >= 6 && hour < 8) return '#ffe4c4'
    if (hour >= 8 && hour < 12) return '#87ceeb'
    if (hour >= 12 && hour < 17) return '#6bb3d9'
    if (hour >= 17 && hour < 20) return '#ff9966'
    if (hour >= 20 && hour < 22) return '#4a4a6a'
    return '#1a1a2e'
  }

  private drawUI(ctx: CanvasRenderingContext2D, state: GameState) {
    // Mini-map
    const minimapSize = 120
    const padding = 10
    const mx = padding
    const my = ctx.canvas.height - minimapSize - padding

    // Background
    ctx.fillStyle = 'rgba(0, 0, 0, 0.7)'
    ctx.fillRect(mx, my, minimapSize, minimapSize)
    ctx.strokeStyle = '#ffd700'
    ctx.lineWidth = 2
    ctx.strokeRect(mx, my, minimapSize, minimapSize)

    const scaleX = minimapSize / state.map_width
    const scaleY = minimapSize / state.map_height

    // Draw buildings
    ctx.fillStyle = '#8b4513'
    for (const b of this.MAP_DATA.buildings) {
      ctx.fillRect(mx + b.x * scaleX, my + b.y * scaleY, b.w * scaleX, b.h * scaleY)
    }

    // Draw player
    ctx.fillStyle = '#ff0000'
    ctx.beginPath()
    ctx.arc(mx + state.player.position.x * scaleX, my + state.player.position.y * scaleY, 4, 0, Math.PI * 2)
    ctx.fill()

    // Draw NPCs (using interpolated positions)
    ctx.fillStyle = '#9370db'
    for (const npc of state.npcs) {
      const renderState = getNPCRenderState(npc)
      ctx.beginPath()
      ctx.arc(mx + renderState.renderX * scaleX, my + renderState.renderY * scaleY, 3, 0, Math.PI * 2)
      ctx.fill()
    }
  }

  // Public methods for mouse interaction
  getCameraPosition(): { x: number; y: number } {
    return { x: this.cameraX, y: this.cameraY }
  }

  setHoveredAgent(agentId: string | null) {
    this.hoveredAgentId = agentId
  }

  // Clear NPC interpolation cache (call when resetting/changing maps)
  clearNPCStates() {
    npcRenderStates.clear()
  }

  // Remove a specific NPC from interpolation cache
  removeNPCState(id: string) {
    removeNPCRenderState(id)
  }

  // Draw highlight effect for hovered agent
  private drawHighlight(ctx: CanvasRenderingContext2D, tileX: number, tileY: number) {
    const px = tileX * TILE_SIZE
    const py = tileY * TILE_SIZE

    // Pulsing glow effect
    const pulse = Math.sin(this.frame * 0.1) * 0.3 + 0.7
    const alpha = 0.3 * pulse

    // Outer glow
    ctx.fillStyle = `rgba(255, 215, 0, ${alpha * 0.5})`
    ctx.beginPath()
    ctx.ellipse(px + TILE_SIZE / 2, py + TILE_SIZE, TILE_SIZE * 0.8, TILE_SIZE * 0.2, 0, 0, Math.PI * 2)
    ctx.fill()

    // Inner highlight ring
    ctx.strokeStyle = `rgba(255, 215, 0, ${alpha})`
    ctx.lineWidth = 2
    ctx.beginPath()
    ctx.arc(px + TILE_SIZE / 2, py + TILE_SIZE / 2, TILE_SIZE * 0.6, 0, Math.PI * 2)
    ctx.stroke()

    // Small indicator above character
    ctx.fillStyle = `rgba(255, 215, 0, ${alpha + 0.2})`
    ctx.beginPath()
    ctx.moveTo(px + TILE_SIZE / 2, py - 8)
    ctx.lineTo(px + TILE_SIZE / 2 - 6, py - 16)
    ctx.lineTo(px + TILE_SIZE / 2 + 6, py - 16)
    ctx.closePath()
    ctx.fill()
  }
}
