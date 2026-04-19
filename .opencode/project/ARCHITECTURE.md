# System Architecture Template

**Status:** 🚧 NOT_CREATED  
**Version:** Template v1.0

---

## Overview

This document will contain the system architecture design once the Architecture Agent analyzes the approved Execution Plan.

**Current Status:** Waiting for EXECUTION_PLAN.md approval

---

## Creation Process

### Step 1: Plan Analysis
Architecture Agent reads:
- EXECUTION_PLAN.md (approved)
- MVP_SPEC.md (constraints)
- SKILLS_INDEX.md (available skills)

### Step 2: Component Design
Design:
- Service boundaries
- Communication patterns
- Data flow
- External integrations

### Step 3: Technology Selection
Select and justify:
- Programming languages
- Frameworks
- Databases
- Infrastructure

### Step 4: Documentation
Generate this document with:
- System diagrams
- Component descriptions
- Data models
- API specifications
- Deployment architecture

---

## Expected Structure

Once generated, this document will contain:

```markdown
# System Architecture v1.0

## Overview
One-paragraph description

## Architecture Diagram
[ASCII or image]

## Components
### Component 1: [Name]
- Responsibility
- Technology
- Interfaces

## Data Models
### Model 1: [Name]
```go
struct definition
```

## API Specifications
### Endpoint 1: [Name]
- Method: POST
- Path: /api/v1/...
- Request/Response

## Technology Stack
| Component | Technology | Justification |

## Deployment Architecture
Infrastructure setup

## Required Skills
List of needed skills
```

---

## Architecture Principles

1. **KISS** - Keep it simple
2. **YAGNI** - You ain't gonna need it
3. **Single Responsibility** - One thing well
4. **Loose Coupling** - Interface-based
5. **Fail Fast** - Explicit errors

---

## Next Steps

1. Wait for EXECUTION_PLAN.md approval
2. Architecture Agent creates design
3. Review architecture
4. Begin implementation

---
**Template Version:** 1.0 | **Last Updated:** 2026-01-15
