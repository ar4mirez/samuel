# PRD: API Rate Limiting

> **Generated**: 2025-01-15
> **Status**: Example (demonstrates PRD format)
> **Complexity**: FEATURE mode (6-8 files, Medium complexity)

---

## Introduction

This feature adds rate limiting to all API endpoints to prevent abuse, ensure fair usage, and protect backend resources from DDoS attacks. Currently, the API has no request throttling, making it vulnerable to brute force attacks and resource exhaustion.

**Goal**: Implement configurable rate limiting with IP-based and user-based throttling, graceful degradation, and monitoring.

---

## Goals

1. Prevent API abuse (brute force, scraping, DDoS)
2. Ensure fair resource allocation across users
3. Protect database and backend services from overload
4. Provide clear feedback to clients when rate limited
5. Enable monitoring of rate limit violations
6. Maintain API response time <50ms overhead for rate limit checks

---

## User Stories

**US-001**: As an API consumer, I want to see clear rate limit information in response headers so that I can avoid being throttled.

**US-002**: As a free-tier user, I want fair access to the API (100 req/min) so that I can build and test my application.

**US-003**: As a paid user, I want higher rate limits (1000 req/min) so that I can support production traffic.

**US-004**: As a platform admin, I want to monitor rate limit violations so that I can identify abuse patterns.

**US-005**: As a developer being rate limited, I want a clear error message with retry-after information so that I know when to retry.

---

## Functional Requirements

### Rate Limiting Rules
**FR-001**: The system must implement tiered rate limits based on user authentication status:
- Unauthenticated (IP-based): 60 requests/minute
- Authenticated free tier: 100 requests/minute
- Authenticated paid tier: 1,000 requests/minute
- Admin users: 10,000 requests/minute (monitoring only)

**FR-002**: The system must use sliding window algorithm for accurate rate limiting (not simple counter reset)

**FR-003**: The system must apply rate limits per endpoint category:
- Read operations (GET): Standard limits
- Write operations (POST/PUT/DELETE): 50% of standard limits
- Authentication endpoints: 5 attempts/minute (stricter to prevent brute force)

### Client Communication
**FR-004**: The system must include rate limit headers in all API responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642234567
```

**FR-005**: The system must return HTTP 429 (Too Many Requests) when rate limit exceeded with response body:
```json
{
  "error": "Rate limit exceeded",
  "retryAfter": 45,
  "limit": 100,
  "window": "1 minute"
}
```

**FR-006**: The system must include `Retry-After` header (seconds until reset) when returning 429

### Storage & Performance
**FR-007**: The system must use Redis for rate limit counters (fast, distributed, TTL support)

**FR-008**: The system must add <50ms latency to API response time for rate limit checks

**FR-009**: The system must handle Redis failures gracefully (fail open, log error, don't block requests)

### Monitoring & Admin
**FR-010**: The system must log all rate limit violations with:
- IP address or user ID
- Endpoint accessed
- Current limit and usage
- Timestamp

**FR-011**: The system must provide admin endpoint `/api/admin/rate-limits` showing:
- Current usage per user/IP
- Violation counts
- Top violators

**FR-012**: The system must allow admins to temporarily adjust rate limits for specific users/IPs via API

---

## Non-Goals

- ❌ Geographic-based rate limiting (IP geolocation) - Deferred to v2
- ❌ Dynamic rate limits based on system load - Future enhancement
- ❌ Custom rate limits per API key - Not in scope (use tier system)
- ❌ GraphQL query complexity-based limits - Future consideration
- ❌ Billing integration (charge for overage) - Separate feature

---

## Technical Considerations

### Tech Stack Integration
- **Backend**: Express.js with middleware approach
- **Storage**: Redis (add to infrastructure)
- **Library**: Use `express-rate-limit` + `rate-limit-redis` (mature, well-tested)
- **Deployment**: Add Redis to Docker Compose (development) and production environment

### Architecture
- Implement as Express middleware (applies before route handlers)
- Create separate middleware for different endpoint categories
- Use Redis key structure: `ratelimit:{tier}:{identifier}:{window}`
- Example: `ratelimit:free:user:12345:60` (user 12345, free tier, 60-second window)

### Dependencies
```json
{
  "express-rate-limit": "^7.1.0",
  "rate-limit-redis": "^4.1.0",
  "ioredis": "^5.3.0"
}
```

### File Organization (respecting 300-line limit)
```
src/middleware/
├── rateLimiting/
│   ├── index.ts                 # Main exports (<50 lines)
│   ├── rateLimiter.ts           # Core middleware logic (<200 lines)
│   ├── rateLimiter.test.ts      # Unit tests
│   ├── config.ts                # Rate limit configurations (<100 lines)
│   ├── tiers.ts                 # Tier definitions (<150 lines)
│   └── redis.ts                 # Redis client setup (<100 lines)
src/api/admin/
├── rateLimits.controller.ts     # Admin endpoints (<200 lines)
├── rateLimits.controller.test.ts
```

### Configuration (Environment Variables)
```bash
# Required
REDIS_URL=redis://localhost:6379

# Optional (with defaults)
RATE_LIMIT_WINDOW=60              # seconds (default: 60)
RATE_LIMIT_FREE=100               # requests per window
RATE_LIMIT_PAID=1000
RATE_LIMIT_ADMIN=10000
RATE_LIMIT_AUTH_ENDPOINTS=5       # stricter for auth
```

---

## Design Considerations

### API Response Headers
Every API response includes rate limit information:
```
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642234567
Content-Type: application/json
```

### Error Response (429 Too Many Requests)
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please try again in 45 seconds.",
    "details": {
      "limit": 100,
      "window": "1 minute",
      "retryAfter": 45
    }
  }
}
```

### Admin Dashboard (Future)
- Show real-time rate limit usage
- Alert on unusual patterns (sudden spikes)
- Allow temporary limit adjustments
- Export violation logs

---

## Guardrails Affected

### Security (CRITICAL)
- ✓ All API boundaries have input validation (validate admin endpoints)
- ✓ All environment variables have secure defaults (REDIS_URL)
- ✓ All async operations have timeout/cancellation (Redis operations)
- ✓ Dependencies checked for vulnerabilities (express-rate-limit, ioredis)

**New guardrail**: All API endpoints must include rate limiting middleware

### Testing (CRITICAL)
- ✓ Coverage targets: >80% for rate limiting logic
- ✓ All public APIs have unit tests (middleware functions)
- ✓ Edge cases tested (Redis failure, concurrent requests, boundary conditions)
- ✓ Integration tests for rate limit enforcement

### Performance (CRITICAL)
- ✓ API responses <200ms (rate limit check must be <50ms)
- ✓ No N+1 queries (single Redis call per request)
- ✓ Expensive computations cached (rate limit configs cached in memory)

### Code Quality
- ✓ No file exceeds 300 lines (already organized into small modules)
- ✓ Cyclomatic complexity ≤10 per function
- ✓ All exported functions have type signatures and JSDoc

---

## Success Metrics

### Technical Metrics
- Rate limit check latency <50ms (p95)
- Test coverage >80% for rate limiting module
- Zero rate limit bypasses (security audit)
- Redis availability >99.9% (monitoring)

### Business Metrics
- <1% of legitimate users hit rate limits (good threshold selection)
- >90% reduction in API abuse incidents
- API response time impact <5% (minimal overhead)

### Security Metrics
- 100% of brute force attempts blocked at rate limit layer
- All rate limit violations logged
- Admin monitoring dashboard shows violations in real-time

---

## Implementation Estimate

### Complexity Analysis
- **Middleware Implementation**: ~6,000 tokens (Medium)
- **Redis Integration**: ~4,000 tokens (Medium)
- **Admin Endpoints**: ~5,000 tokens (Medium)
- **Tests**: ~8,000 tokens (Medium)
- **Documentation**: ~2,000 tokens (Simple)

**Total**: ~25,000 tokens (FEATURE mode appropriate)

### Recommended Approach
1. Use FEATURE mode (not complex enough for full COMPLEX with PRD in real project)
2. This is an example PRD to demonstrate format
3. Implement in 2 phases:
   - Phase 1: Core rate limiting (IP-based, basic tiers)
   - Phase 2: Admin endpoints and monitoring
4. Frequent checkpoints after each subtask

---

## Open Questions

1. **Redis Infrastructure**: Use existing Redis instance or deploy new one? Shared or dedicated?
2. **Rate Limit Bypass**: Should admins be able to whitelist specific IPs?
3. **Distributed Systems**: How to handle rate limits across multiple API servers? (Redis solves this, but confirm)
4. **Testing in CI/CD**: How to test Redis-dependent code in GitHub Actions? (Use Redis container?)
5. **Monitoring Tool**: Integrate with existing monitoring (Datadog/Sentry) or custom dashboard?

**Action**: Clarify these before generating task list.

---

## Next Steps

1. ✅ PRD reviewed and approved
2. ⏭️ Use `@.agent/skills/generate-tasks/SKILL.md` to create task breakdown
3. ⏭️ Implement tasks step-by-step
4. ⏭️ Update `.agent/project.md` with rate limiting architecture
5. ⏭️ Add to `.agent/patterns.md` (rate limiting pattern for future features)

---

**This is an example PRD**. In a real project, AI would generate this after asking clarifying questions using `@.agent/skills/create-prd/SKILL.md`.
