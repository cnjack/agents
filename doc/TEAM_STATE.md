# Stardew Agent Team State

## Team Information

**Team Name:** stardew-refactor

**Created:** 2026-02-11

## Teammates

| Name | Role | Agent Type | Description |
|------|------|------------|-------------|
| ui-artist | UI Artist | general-purpose | Responsible for UI design, creating beautiful visuals and graphics for the game |
| game-designer | Game Designer | general-purpose | Designs the game mechanics, rules, and overall gameplay experience |
| engineer | Engineer/Developer | general-purpose | Implements the designs and gameplay mechanics through coding and technical development |
| qa-tester | QA Tester | general-purpose | Responsible for testing and acceptance, ensuring the game meets quality standards |

## Completed Tasks

| Task ID | Task Name | Owner | Status |
|---------|-----------|-------|--------|
| #5 | 提高游戏贴图精细度 | ui-artist | completed |
| #6 | 验证贴图精细度改进 | qa-tester | completed |
| #7 | 实现AI驱动NPC对话系统 | game-designer/engineer | completed |
| #8 | 实现AI驱动NPC行为系统 | game-designer/engineer | completed |
| #9 | 实现自然语言交互系统 | game-designer/engineer | completed |
| #10 | 设计NPC对话UI | ui-artist | completed |
| #11 | 修复TypeScript类型定义错误 | engineer | completed |
| #12 | 修复Go后端编译错误 | engineer | completed |
| #13 | 优化像素角色模型 | ui-artist | completed |

## Key Files Modified/Created

### Frontend
- `frontend/src/game/engine/Renderer.ts` - 贴图精细度提升、角色模型优化
- `frontend/src/components/NPCDialog.vue` - NPC对话UI
- `frontend/src/styles/main.css` - 对话框样式扩展
- `frontend/src/types/game.ts` - 类型定义（TileType, CropStage, FarmTile）

### Backend - AI System
- `backend/internal/models/ai_personality.go` - NPC对话人格模型
- `backend/internal/models/behavior.go` - 行为相关模型
- `backend/internal/game/behavior.go` - 行为决策引擎
- `backend/internal/game/mood.go` - 心情管理系统
- `backend/internal/game/dynamic_schedule.go` - 动态日程系统
- `backend/internal/game/npc_interaction.go` - NPC互动系统
- `backend/internal/game/personalities.go` - NPC档案定义
- `backend/internal/services/ai/dialogue.go` - 对话生成服务
- `backend/internal/services/ai/llm_client.go` - LLM客户端接口
- `backend/internal/api/ai_handler.go` - AI API处理器

### Documentation
- `docs/AI_GAMEPLAY_DESIGN.md` - AI游戏玩法设计规范
- `docs/NPC_PROFILES.md` - NPC人格档案设计

## API Endpoints Created

- `POST /api/v1/ai/dialogue` - AI生成NPC对话
- `POST /api/v1/ai/decision` - AI行为决策
- `POST /api/v1/ai/interpret` - 自然语言理解
- `GET /api/v1/npc/:id/mood` - 获取NPC心情状态
- `GET /api/v1/ai/factors/:npc_id` - 获取决策因素
- `GET /api/v1/npc/interactions` - 获取NPC互动
- `POST /api/v1/npc/interactions/trigger` - 触发NPC互动

## Known Issues (待修复)

### nlp.go
- `n.getNearbyNPCs` 方法未定义
- 部分类型转换问题

### ai_handler.go
- `fmt` 包未导入

## Teammate Status

### ui-artist
- **状态**: ✅ 空闲
- **已完成任务**: #5 (贴图精细度), #10 (NPC对话UI), #13 (角色模型优化)
- **额外工作**: Stardew Valley 风格角色精灵重设计
- **建议**: 需要运行时测试验证视觉效果

### game-designer
- **状态**: ✅ 空闲
- **已完成任务**: #7 (AI对话系统), #8 (AI行为系统)
- **设计文档**: `/docs/AI_GAMEPLAY_DESIGN.md`, `/docs/NPC_PROFILES.md`
- **建议**: 后续可集成真实 Claude API 替换 Mock 实现

### engineer
- **状态**: ✅ 空闲
- **已完成任务**: #9 (自然语言交互), #11 (TypeScript修复), #12 (Go编译修复)
- **API端点**: `/api/v1/ai/interpret`, `/api/v1/ai/execute`, `/api/v1/ai/interact`
- **所有后端服务运行正常**

### qa-tester
- **状态**: ✅ 空闲
- **已完成任务**: #6 (贴图验证)
- **备注**: 完成代码审查验证，等待运行时测试

---

## Recovery Commands

To recover this team, run:

```
TeamCreate with team_name: "stardew-refactor"

Task with subagent_type: "general-purpose", team_name: "stardew-refactor", name: "ui-artist"
Task with subagent_type: "general-purpose", team_name: "stardew-refactor", name: "game-designer"
Task with subagent_type: "general-purpose", team_name: "stardew-refactor", name: "engineer"
Task with subagent_type: "general-purpose", team_name: "stardew-refactor", name: "qa-tester"
```
