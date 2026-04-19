# ⚖️ System Rules

## Overview

These 7 rules are **STRICTLY ENFORCED** by the Reviewer Agent and Orchestrator. Violations result in automatic rejection.

---

## Rule 1: STATE.json is the Only Truth

**Statement:** All system state MUST be maintained in STATE.json. No agent may bypass or ignore it.

**Rationale:**
- Single source of truth prevents confusion
- Enables coordination between agents
- Allows human oversight
- Tracks progress and violations

**Enforcement:**
- Orchestrator writes to STATE.json
- All agents read from STATE.json
- Bypassing state = violation

**Violation:** REJECT

---

## Rule 2: No Scope Expansion Beyond Plan

**Statement:** Agents MUST implement ONLY what is in EXECUTION_PLAN.md. No additional features.

**Rationale:**
- Prevents scope creep
- Keeps project on track
- Maintains focus on MVP
- Enables accurate estimation

**Enforcement:**
- Reviewer checks all deliverables against plan
- Extra features = immediate rejection
- Must remove extra code and resubmit

**Violation:** REJECT

**Example:**
```go
// ❌ WRONG - Adding feature not in plan
// Task: User registration
// Backend also implements: user preferences, settings, notifications
```

---

## Rule 3: Skills MUST Be Reused

**Statement:** Backend Agent MUST use existing skills from SKILLS_INDEX.md. Creating duplicate skills is prohibited.

**Rationale:**
- Prevents knowledge duplication
- Ensures consistency
- Builds reusable library
- Reduces maintenance

**Enforcement:**
- Reviewer checks skill usage
- New skills created = violation
- Must use existing skill or request creation

**Violation:** REJECT

**Example:**
```go
// ❌ WRONG - Creating duplicate skill
// SKILLS_INDEX.md has: golang-error-handling
// Backend creates: internal/errors/custom.go (duplicate!)

// ✅ CORRECT - Using existing skill
// Backend imports and follows golang-error-handling patterns
```

---

## Rule 4: Architecture is Immutable

**Statement:** ARCHITECTURE.md is immutable without explicit human approval. No agent may modify it.

**Rationale:**
- Ensures consistency
- Prevents chaos during implementation
- Requires conscious human decision
- Maintains system integrity

**Enforcement:**
- Reviewer detects any architecture changes
- Changes = immediate rejection
- Human must approve architecture update

**Violation:** REJECT

**Process for Updates:**
1. Identify need for change
2. Request human approval via Orchestrator
3. Human approves
4. Architecture Agent creates new version
5. Backend notified

---

## Rule 5: Backend-Review Loop Max 3 Cycles

**Statement:** Backend Agent has maximum 3 attempts to pass review. After 3 rejections, escalate to human.

**Rationale:**
- Prevents infinite loops
- Detects fundamental issues
- Conserves resources
- Requires human judgment

**Enforcement:**
- Orchestrator tracks loop count
- Reviewer increments counter
- At 3rd rejection = escalate

**Violation:** ESCALATE

**Process:**
```
Cycle 1: Rejected → Fix → Resubmit
Cycle 2: Rejected → Fix → Resubmit
Cycle 3: Rejected → ESCALATE TO HUMAN
```

---

## Rule 6: Product Owner is Fully Isolated

**Statement:** Product Owner Agent has NO access to engineering system. Engineering agents have NO access to Product Owner.

**Rationale:**
- Maintains separation of concerns
- Prevents product direction changes mid-execution
- Ensures clear handoff
- Prevents confusion of roles

**Enforcement:**
- Product Owner doesn't know STATE.json
- Engineering doesn't modify product artifacts
- Any contact = violation

**Violation:** REJECT

---

## Rule 7: KISS Principle (Keep It Simple, Stupid)

**Statement:** Solutions MUST be as simple as possible. No over-engineering.

**Rationale:**
- Easier to understand
- Easier to maintain
- Less prone to bugs
- Faster development

**Enforcement:**
- Reviewer assesses complexity
- Architecture Agent follows KISS
- Backend implements simply

**Violation:** REJECT

**Examples:**
```
❌ Over-engineering:
- Microservices for MVP
- Event sourcing for simple CRUD
- Complex caching before measuring
- Premature optimization

✅ Simple solutions:
- Monolith for MVP
- REST API for simple operations
- Add caching after profiling
- Optimize after identifying bottleneck
```

---

## Violation Severity

| Severity | Action | Examples |
|----------|--------|----------|
| CRITICAL | Immediate reject + log | Architecture modification, security issues |
| HIGH | Reject + require fix | Scope expansion, skill duplication |
| MEDIUM | Reject + guidance | Missing tests, code quality issues |
| LOW | Comment + suggestion | Minor style issues |

## Violation Log

All violations are logged in STATE.json:

```json
{
  "violations": [
    {
      "id": "VIO-001",
      "rule": "RULE-004",
      "severity": "CRITICAL",
      "agent": "backend",
      "description": "Modified data model without approval",
      "detected_at": "2026-01-15T14:30:00Z",
      "resolution": "REJECTED"
    }
  ]
}
```

## Human Override

In exceptional cases, human may override rules:

1. Escalation from Reviewer (Rule 5)
2. Emergency security fix
3. Architecture update approval (Rule 4)
4. Scope change with business justification (Rule 2)

**Process:**
- Orchestrator requests override
- Human reviews and decides
- Decision logged in STATE.json
- System proceeds according to decision

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-15  
**Enforcement:** Automatic via Reviewer Agent
