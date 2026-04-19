| Слой системы     | Роль агента      | Primary модель   | Fallback     | Cheap mode   | Задача                           |
| ---------------- | ---------------- | ---------------- | ------------ | ------------ | -------------------------------- |
| 🟡 Product       | Product Owner    | **Kimi K2.5**    | Qwen3.6 Plus | Qwen3.5 Plus | продукт, roadmap, MVP → prod     |
| 🧠 Planning      | Planner          | **Qwen3.6 Plus** | Kimi K2.5    | GLM-5        | backlog → execution plan         |
| ⚙️ Orchestration | Orchestrator     | **Qwen3.6 Plus** | Qwen3.5 Plus | GLM-5        | STATE.json, контроль pipeline    |
| 🏗️ Architecture | System Architect | **Kimi K2.5**    | Qwen3.6 Plus | GLM-5        | system design, структура backend |
| ⚙️ Backend       | Execution Engine | **GLM-5.1**      | Qwen3.6 Plus | Qwen3.5 Plus | код, API, implementation         |
| 🔍 Review        | Code Reviewer    | **Qwen3.6 Plus** | Kimi K2.5    | GLM-5        | контроль качества, violations    |
| 🧩 Skills        | Module Generator | **GLM-5 / 5.1**  | Qwen3.5 Plus | MiniMax M2.7 | auth/db/logging skills           |


Kimi (Product Owner)
      ↓
Qwen (Planner)
      ↓
Qwen (Orchestrator)
      ↓
Kimi (Architecture)
      ↓
GLM (Backend Execution)
      ↓
Qwen (Reviewer Loop)
      ↓
DONE