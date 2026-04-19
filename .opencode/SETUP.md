# Integration Guide

## Setup Instructions

This guide helps you integrate the AI Multi-Agent System into your existing project.

---

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/alekparkhomenko/ai-multi-agent-system.git
cd ai-multi-agent-system
```

### 2. Copy to Your Project

```bash
# Navigate to your project
cd /path/to/your/project

# Create .opencode directory
mkdir -p .opencode

# Copy agents
cp -r /path/to/ai-multi-agent-system/agents .opencode/

# Create product directory
mkdir -p .opencode/product

# Create project directory
mkdir -p .opencode/project
```

### 3. Copy Project Files

```bash
# Copy state schema
cp /path/to/ai-multi-agent-system/project/STATE.json .opencode/project/

# Copy rules
cp /path/to/ai-multi-agent-system/project/RULES.md .opencode/project/

# Copy skills index
cp /path/to/ai-multi-agent-system/project/SKILLS_INDEX.md .opencode/project/

# Copy decisions template
cp /path/to/ai-multi-agent-system/project/DECISIONS.md .opencode/project/

# Copy architecture template
cp /path/to/ai-multi-agent-system/project/ARCHITECTURE.md .opencode/project/

# Copy execution plan template
cp /path/to/ai-multi-agent-system/project/EXECUTION_PLAN.md .opencode/project/
```

### 4. Initialize Product Artifacts

Copy example product files (optional):

```bash
# Copy product examples (or create your own)
cp /path/to/ai-multi-agent-system/product/PRODUCT_VISION.md .opencode/product/
cp /path/to/ai-multi-agent-system/product/ROADMAP.md .opencode/product/
cp /path/to/ai-multi-agent-system/product/MVP_SPEC.md .opencode/product/
cp /path/to/ai-multi-agent-system/product/BACKLOG.md .opencode/product/
```

---

## Directory Structure After Setup

```
your-project/
├── .opencode/
│   ├── agents/
│   │   ├── product-owner.md
│   │   ├── planner.md
│   │   ├── orchestrator.md
│   │   ├── architecture.md
│   │   ├── backend.md
│   │   └── reviewer.md
│   ├── product/
│   │   ├── PRODUCT_VISION.md
│   │   ├── ROADMAP.md
│   │   ├── MVP_SPEC.md
│   │   └── BACKLOG.md
│   └── project/
│       ├── STATE.json
│       ├── EXECUTION_PLAN.md
│       ├── ARCHITECTURE.md
│       ├── DECISIONS.md
│       ├── RULES.md
│       └── SKILLS_INDEX.md
├── your-existing-code/
└── README.md
```

---

## First Run

### 1. Customize STATE.json

Edit `.opencode/project/STATE.json`:

```json
{
  "system": {
    "name": "your-project-name",
    "description": "Your project description",
    "project_type": "software",
    "language": "go"
  }
  ...
}
```

### 2. Create Product Artifacts

Start with Product Owner agent to generate:
- PRODUCT_VISION.md
- ROADMAP.md
- MVP_SPEC.md
- BACKLOG.md

### 3. Run Planner

Planner will convert BACKLOG.md → EXECUTION_PLAN.md

### 4. Human Approval

**STOP!** Review EXECUTION_PLAN.md and approve.

### 5. Continue Workflow

- Architecture design
- Backend implementation
- Review loops
- Completion

---

## Configuration

### STATE.json Config Section

```json
{
  "config": {
    "review_max_loops": 3,
    "min_test_coverage": 80,
    "auto_escalate_after_max_loops": true,
    "require_human_approval": true,
    "architecture_immutable": true,
    "enforce_skill_reuse": true
  }
}
```

### Customize for Your Project

- Adjust `min_test_coverage` (default: 80)
- Change `review_max_loops` (default: 3)
- Modify rules as needed

---

## Adding Skills

Create skill files in `agents/skills/`:

```
opencode/agents/skills/
├── golang-error-handling/
│   └── SKILL.md
├── golang-testing/
│   └── SKILL.md
└── your-custom-skill/
    └── SKILL.md
```

Update `SKILLS_INDEX.md` to register new skills.

---

## Troubleshooting

### Issue: Agents not working

**Solution:** Check STATE.json format is valid JSON.

### Issue: Approval gate not stopping

**Solution:** Verify `config.require_human_approval` is `true`.

### Issue: Skills not being reused

**Solution:** Ensure skills are registered in SKILLS_INDEX.md.

---

## Next Steps

1. Read agent definitions in `agents/`
2. Understand RULES.md enforcement
3. Review example product artifacts
4. Start with Product Owner workflow
5. Iterate and improve

---

**Document Version:** 1.0 | **Last Updated:** 2026-01-15
