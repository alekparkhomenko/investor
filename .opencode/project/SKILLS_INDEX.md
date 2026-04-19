# 🧩 Skills Index

## Overview

This registry contains all reusable engineering skills available to the Backend Agent. Skills MUST be reused - no duplication allowed.

---

## Core Skills

### golang-error-handling
**Category:** Language - Error Management  
**Description:** Patterns for error creation, wrapping, and handling in Go  
**Usage:** All error handling code  
**Path:** `agents/skills/golang-error-handling/SKILL.md`

**Key Patterns:**
- Sentinel errors (var ErrX = errors.New(...))
- Error wrapping with context (fmt.Errorf("%w: ...", err))
- Error checking with errors.Is()
- Error inspection with errors.As()
- Joining multiple errors with errors.Join()

---

### golang-testing
**Category:** Language - Testing  
**Description:** Testing patterns and best practices for Go  
**Usage:** All test files  
**Path:** `agents/skills/golang-testing/SKILL.md`

**Key Patterns:**
- Table-driven tests
- Testify assertions (assert, require)
- Mock generation with mockery
- Parallel test execution
- Coverage requirements (80%+)

---

### golang-concurrency
**Category:** Language - Concurrency  
**Description:** Goroutines, channels, and synchronization patterns  
**Usage:** Concurrent code, background tasks  
**Path:** `agents/skills/golang-concurrency/SKILL.md`

**Key Patterns:**
- Worker pools
- Channel patterns
- sync.WaitGroup
- Context cancellation
- Goroutine leak prevention

---

### golang-database
**Category:** Language - Database  
**Description:** Database access patterns with sql, sqlx, pgx  
**Usage:** All database operations  
**Path:** `agents/skills/golang-database/SKILL.md`

**Key Patterns:**
- Connection pooling
- Transaction handling
- Query building
- Result scanning
- Migration management

---

### golang-dependency-injection
**Category:** Language - Architecture  
**Description:** DI patterns and service wiring  
**Usage:** Service initialization, testing  
**Path:** `agents/skills/golang-dependency-injection/SKILL.md`

**Key Patterns:**
- Constructor injection
- Interface-based design
- Service containers
- Wire code generation
- Testing with mocks

---

### golang-api-design
**Category:** Language - API  
**Description:** REST API design and implementation patterns  
**Usage:** HTTP handlers, middleware  
**Path:** `agents/skills/golang-api-design/SKILL.md`

**Key Patterns:**
- Handler structure
- Middleware chaining
- Request/response models
- Status codes
- Error responses

---

## Domain Skills

### auth-jwt
**Category:** Domain - Authentication  
**Description:** JWT-based authentication implementation  
**Usage:** Login, token management  
**Path:** `agents/skills/auth-jwt/SKILL.md`

**Key Patterns:**
- Token generation
- Token validation
- Refresh tokens
- Middleware
- Secure storage

---

### validation-input
**Category:** Domain - Validation  
**Description:** Input validation patterns  
**Usage:** API input, form validation  
**Path:** `agents/skills/validation-input/SKILL.md`

**Key Patterns:**
- Struct validation with tags
- Custom validators
- Error aggregation
- Sanitization

---

### logging-structured
**Category:** Domain - Observability  
**Description:** Structured logging with zap/slog  
**Usage:** All logging code  
**Path:** `agents/skills/logging-structured/SKILL.md`

**Key Patterns:**
- JSON logging
- Log levels
- Context fields
- Error logging
- Performance

---

## Infrastructure Skills

### docker-containerization
**Category:** Infrastructure - Containerization  
**Description:** Docker patterns for Go applications  
**Usage:** Container builds, deployments  
**Path:** `agents/skills/docker-containerization/SKILL.md`

**Key Patterns:**
- Multi-stage builds
- Scratch images
- Health checks
- Docker Compose

---

### ci-cd-github-actions
**Category:** Infrastructure - CI/CD  
**Description:** GitHub Actions workflows  
**Usage:** Automated testing, deployment  
**Path:** `agents/skills/ci-cd-github-actions/SKILL.md`

**Key Patterns:**
- Test automation
- Lint checking
- Build pipelines
- Deployment

---

## Adding New Skills

**Process:**
1. Identify need during Architecture phase
2. Create SKILL.md following template
3. Add to SKILLS_INDEX.md
4. Reference in EXECUTION_PLAN.md tasks

**Skill Template:**
```markdown
# [Skill Name]

## Overview
Brief description

## When to Use
Applicable scenarios

## Patterns
### Pattern 1
Description and example

### Pattern 2
Description and example

## Anti-Patterns
What NOT to do

## Examples
### Good
```
code example
```

### Bad
```
code example
```

## References
Links to docs, articles
```

---

**Total Skills:** 10  
**Last Updated:** 2026-01-15  
**Enforcement:** Rule 3 - Skills MUST be reused
