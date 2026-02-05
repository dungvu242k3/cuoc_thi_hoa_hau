# Plan: Backend Code Review (Vietnamese)

## Goal
Comprehensive code review and function documentation for the Backend (`BE`) codebase.
The output will be a review document explaining what each key function does in Vietnamese, along with code review notes.

## Scope
- Directory: `BE`
- Key Areas: `cmd`, `config`, `internal`
- Focus: Graph Resolvers, Middleware, Business Logic.

## Task Breakdown

### Phase 1: Preparation & Config Review
- [ ] Review `cmd/api/main.go` and `cmd/seeder/main.go` (System Entry points)
- [ ] Review `config/` (Configuration & Environment variables)
- [ ] Review `internal/adapter/cache` (Infrastructure adapters like Redis)

### Phase 2: Core Resolvers (API Layer)
Review files in `internal/adapter/graph/resolver/`:
- [ ] `auth.resolvers.go` (Authentication logic)
- [ ] `contestant.resolvers.go` (Contestant management)
- [ ] `score.resolvers.go` (Scoring System)
- [ ] `schedule.resolvers.go` (Scheduling System)
- [ ] `feedback.resolvers.go`, `user.resolvers.go`, etc.

### Phase 3: Middleware & Shared Logic
- [ ] `internal/adapter/graph/middleware/` (Authentication, Rate limiting)
- [ ] `internal/adapter/graph/resolver/helper.go` (Helper functions)
- [ ] `internal/adapter/graph/resolver/mapper.go` (Data mappers)

## Output Format
The results will be compiled into `docs/BE-CODE-REVIEW.md` using the following format:

### [Module Name]

| File | Function | Description (VN) | Review Notes (VN) |
|------|----------|------------------|-------------------|
| `auth.resolvers.go` | `Login` | Xử lý đăng nhập user, trả về JWT | Nên, thêm rate_limit strict hơn |
| `score.resolvers.go` | `SubmitScore` | Chấm điểm cho thí sinh | Logic tính điểm cần double check |

## Verification Plan
- **Manual Review**: Self-check the generated descriptions against the code to ensure accuracy.
- **User Verification**: Submit the document for user feedback.
