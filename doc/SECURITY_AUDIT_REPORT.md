# Security Vulnerability Audit Report

**Date**: 2026-02-12
**Project**: Stardew Community Agent
**Auditor**: R&D Team

---

## Executive Summary

This security audit identifies vulnerabilities and security concerns in the codebase. The audit covers:
- Backend Go code
- Frontend TypeScript/Vue code
- Dependencies
- Configuration management
- API security

### Risk Summary

| Severity | Count |
|----------|-------|
| Critical | 0 |
| High | 2 |
| Medium | 5 |
| Low | 3 |

---

## Critical Issues

**None identified.**

---

## High Severity Issues

### 1. API Key Exposed in Configuration File

**Location**: `backend/config/local.yaml:15`

**Description**: The API key for the AI service is stored in plaintext in the configuration file.

**Risk**: If this file is accidentally committed to version control, the API key would be exposed, allowing unauthorized access to the AI service.

**Current Mitigation**:
- File is added to `.gitignore`

**Recommendation**:
- Use environment variables for sensitive data
- Implement secrets management (e.g., HashiCorp Vault, AWS Secrets Manager)
- Add pre-commit hooks to prevent accidental commits

**Status**: Partially mitigated

---

### 2. No Authentication/Authorization

**Location**: `backend/internal/api/middleware.go:56-63`

**Description**: The API endpoints have no authentication mechanism. The `AuthMiddleware` is a placeholder that passes all requests through.

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // For now, no authentication required
        // In production, validate API key or JWT token
        c.Next()
    }
}
```

**Risk**: Anyone with network access can call the API endpoints, including potentially sensitive operations like resetting the game state.

**Recommendation**:
- Implement JWT-based authentication
- Add API key validation for external access
- Implement role-based access control

**Status**: Not implemented

---

## Medium Severity Issues

### 3. Outdated Dependencies with Known Vulnerabilities (Frontend)

**Location**: `frontend/package.json`

**Vulnerabilities Found**:

| Package | Severity | Issue |
|---------|----------|-------|
| esbuild <= 0.24.2 | Moderate | Development server allows any website to read responses |
| vite 0.11.0 - 6.1.6 | Moderate | Depends on vulnerable esbuild |
| vue-template-compiler | Moderate | Client-side XSS vulnerability |
| @vue/language-core <= 2.0.28 | Moderate | Depends on vulnerable vue-template-compiler |
| vue-tsc 1.7.0 - 2.0.28 | Moderate | Depends on vulnerable @vue/language-core |

**Recommendation**:
```bash
cd frontend && npm audit fix --force
```

**Status**: Not fixed (breaking changes required)

---

### 4. Outdated Go Dependencies

**Location**: `backend/go.mod`

**Vulnerable Packages**:

| Package | Current Version | Issue |
|---------|----------------|-------|
| golang.org/x/crypto | v0.14.0 | Update recommended |
| golang.org/x/net | v0.17.0 | Update recommended |

**Recommendation**:
```bash
cd backend && go get -u ./... && go mod tidy
```

**Status**: Not updated

---

### 5. CORS Configuration Too Permissive

**Location**: `backend/internal/api/routes.go:16-23`

**Description**: CORS is configured to allow all origins (`*`).

```go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
```

**Risk**: Allows any website to make requests to the API, enabling potential CSRF attacks.

**Recommendation**:
- Restrict `AllowOrigins` to known domains
- Disable `AllowCredentials` when using wildcard origins (currently violates CORS spec)

**Status**: Not fixed

---

### 6. Rate Limiting Not Implemented

**Location**: `backend/internal/api/middleware.go:47-54`

**Description**: The rate limiter middleware is a placeholder.

```go
func RateLimiter() gin.HandlerFunc {
    // In production, use a proper rate limiter like redis-based
    return func(c *gin.Context) {
        // For now, just pass through
        c.Next()
    }
}
```

**Risk**: API is vulnerable to DoS attacks and abuse.

**Recommendation**:
- Implement token bucket or sliding window rate limiting
- Use Redis for distributed rate limiting
- Add per-IP and per-user limits

**Status**: Not implemented

---

### 7. Sensitive Data in Logs

**Location**: Multiple files

**Description**: Request logging captures full request paths which might contain sensitive parameters.

**Risk**: Sensitive data might be logged and persisted.

**Recommendation**:
- Implement log sanitization
- Avoid logging request bodies
- Use structured logging with sensitive field masking

**Status**: Not addressed

---

## Low Severity Issues

### 8. Missing Security Headers

**Location**: HTTP responses

**Description**: The following security headers are not set:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy`

**Recommendation**: Add security headers middleware.

**Status**: Not implemented

---

### 9. Error Messages May Leak Information

**Location**: Various API handlers

**Description**: Some error messages include internal details that could help attackers.

```go
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
```

**Recommendation**:
- Return generic error messages to clients
- Log detailed errors server-side only

**Status**: Not addressed

---

### 10. Insecure Random Number Generation

**Location**: `backend/internal/api/handler.go:4,160`

**Description**: Uses `math/rand` instead of `crypto/rand` for game logic.

```go
import "math/rand"
// ...
dialogue = targetNPC.Dialogue[rand.Intn(len(targetNPC.Dialogue))]
```

**Risk**: Predictable random numbers could be exploited in game mechanics.

**Recommendation**: For any security-sensitive randomness, use `crypto/rand`.

**Status**: Low priority (game logic, not security-critical)

---

## Dependency Audit

### Backend (Go)

| Package | Version | Status |
|---------|---------|--------|
| gin-gonic/gin | v1.9.1 | OK |
| gorilla/websocket | v1.5.1 | OK |
| golang.org/x/crypto | v0.14.0 | Update recommended |
| golang.org/x/net | v0.17.0 | Update recommended |

### Frontend (npm)

| Package | Version | Status |
|---------|---------|--------|
| vue | ^3.4.0 | OK |
| vite | ^5.0.0 | Vulnerable - update needed |
| pinia | ^2.1.7 | OK |

---

## Security Best Practices Checklist

- [x] API keys not committed to version control (`.gitignore`)
- [ ] Input validation on all API endpoints
- [ ] Authentication implemented
- [ ] Authorization/role-based access control
- [ ] Rate limiting
- [ ] CORS properly configured
- [ ] Security headers set
- [ ] HTTPS enforced in production
- [ ] Dependency scanning in CI/CD
- [ ] Secrets management solution
- [ ] Error handling doesn't leak information
- [ ] Logging sanitized
- [ ] Regular security audits scheduled

---

## Recommended Actions

### Immediate (P0)

1. **Update frontend dependencies**:
   ```bash
   cd frontend && npm audit fix --force
   ```

2. **Update Go dependencies**:
   ```bash
   cd backend && go get -u ./... && go mod tidy
   ```

### Short-term (P1)

3. **Implement authentication middleware**
4. **Fix CORS configuration**
5. **Add security headers**

### Long-term (P2)

6. **Implement rate limiting**
7. **Set up secrets management**
8. **Add dependency scanning to CI/CD**
9. **Implement comprehensive logging**

---

## Conclusion

The application has several security concerns that should be addressed before production deployment. The most critical issues are:
1. Lack of authentication
2. Outdated dependencies with known vulnerabilities
3. Permissive CORS configuration

These issues should be prioritized for remediation.

---

*Report generated by Security Audit Tool*
*Version: 1.0*
