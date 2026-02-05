# Plan: Review & Improve Contestant Module

## Goal
Review the "Contestant" (Thí sinh) module to enhance security functions (rate limiting, input validation) and optimize code quality (variable usage, function clarity).

## Analysis
After reviewing `contestant.resolvers.go`, `contestant_svc.go`, and `validation/contestant.go`:

### 1. Code Optimization (Variables & Functions)
- **Issue**: The resolver `CreateContestantProfile` and `UpdateContestantProfile` contain verbose manual mapping logic (`input.X != nil` checks) which clutters the resolver with temporary variables.
- **Improvement**: Extract this mapping logic into a dedicated `mapper` package or helper functions. This cleans up the resolver and separates concerns (DTO -> Domain translation).

### 2. Security Functions
- **Rate Limiting**: Currently implemented directly inside the resolver (`checkVoteRateLimit`).
- **Improvement**: Isolate this logic. While keeping it in the resolver is acceptable for now, using a dedicated service method or middleware pattern makes it more testable and reusable.
- **Input Validation**: The `ContainsHTMLTags` function uses a blacklist approach (`<script`, etc.). This is brittle.
- **Improvement**: Consider using a robust library or stricter `html.EscapeString` at the boundary if rich text is not required.
- **Authorization**: Role checks are present (`RequirePermission`). Audit logs are present in `contestant_svc.go`.

## Proposed Changes

### Task 1: Refactor Mappers (Optimization)
- [ ] Create `internal/adapter/graph/mapper/contestant_input.go`.
- [ ] Move DTO-to-Domain mapping logic from `resolvers` to this new mapper.
- [ ] Simplify resolver functions to just call `mapper.ToDomain(input)`.

### Task 2: Harden Security
- [ ] **Rate Limiting**: Review `checkVoteRateLimit`. Ensure Atomic increments are handled correctly (currently looks okay but "fail open" on cache error might be a risk).
- [ ] **Validation**: Review `validation/contestant.go`.
    -   Enhance `ValidatePhone` to be more robust (if needed).
    -   Review `ContainsHTMLTags` for bypass vectors.

### Task 3: Variable Optimization
- [ ] reduce usage of pointer dereferencing in the main flow by using helper accessors if possible, or just the mapper refactor will cover this.

## Agent Assignment
- `backend-specialist`: For Go refactoring and security hardening.

## Questions for User
- Do you have specific performance bottlenecks with the current variables?
- Are there specific security incidents that prompted this review?
