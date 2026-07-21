# WEEK 2 IMPLEMENTATION PLAN

**Timeline:** 5 working days  
**Phase 1 (Days 1-3):** Polish & Testing  
**Phase 2 (Days 4-5):** Reverse Proxy Implementation

---

## PHASE 1: POLISH & TESTING (Days 1-3)

### Day 1: Code Review & Bug Fixes

**Morning (2 hours):**
- [ ] Review all 37 handlers for consistency
- [ ] Check error handling (are we returning proper status codes?)
- [ ] Verify JSON response format consistency
- [ ] Check for unhandled edge cases

**Checklist:**
```go
// Check all handlers follow this pattern:
func Handler(c *fiber.Ctx) error {
    userID, err := middleware.GetUserID(c)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
    }
    
    // Logic here
    
    return c.Status(200).JSON(fiber.Map{"success": true, "data": result})
}
```

**Afternoon (2 hours):**
- [ ] Fix any inconsistent error responses
- [ ] Add missing error handling in services
- [ ] Validate all database operations have error checks
- [ ] Test all error paths

**Output:**
- All 37 handlers follow consistent patterns
- No panics on invalid input
- Proper HTTP status codes (201, 400, 401, 404, 500)

---

### Day 2: Unit Tests (Core Services)

**Focus:** Test the most critical business logic

**Priority 1 - Auth Service (MUST TEST):**
```go
// tests/auth_service_test.go

func TestRegisterUser(t *testing.T) {
    // Valid registration
    // Duplicate email
    // Weak password
    // Invalid email format
}

func TestLoginUser(t *testing.T) {
    // Correct password
    // Wrong password
    // Non-existent user
    // Risk assessment calculation
}

func TestJWTGeneration(t *testing.T) {
    // Token format valid
    // Token expires in 7 days
    // Token can be verified
}
```

**Priority 2 - Token Generator (MUST TEST):**
```go
func TestTokenGeneration(t *testing.T) {
    for _, tokenType := range []string{"url", "api_key", "document", "dns", "email"} {
        token := GenerateTokenValue(tokenType)
        // Assert format
        // Assert prefix: sk_
        // Assert length: 32+ chars
    }
}
```

**Priority 3 - Classifier (SHOULD TEST):**
```go
func TestTrafficClassification(t *testing.T) {
    // Low risk request (legitimate user)
    // Medium risk (new location)
    // High risk (SQL injection in payload)
    // Critical risk (known bad IP)
}
```

**Output:**
- ✅ All auth tests pass
- ✅ All token generation tests pass
- ✅ All classifier tests pass
- ✅ 80%+ test coverage on core services

---

### Day 3: Integration Tests & API Testing

**API Test Suite (Postman/curl):**

```bash
# Create postman_collection.json with all 37 endpoints

# Test sequence:
1. Register user → Get user_id
2. Login with that user → Get JWT token
3. Use JWT in all subsequent requests
4. Create token → Get token_value
5. List tokens → Verify token appears
6. Create alert config
7. Get dashboard stats
8. Delete token
9. Logout flows
```

**Test Matrix:**

| Endpoint | Status | Input Validation | Auth Check | Response Format |
|----------|--------|------------------|-----------|-----------------|
| POST /api/auth/register | ✅ | TBD | TBD | TBD |
| POST /api/auth/login | ✅ | TBD | TBD | TBD |
| GET /api/user/profile | ✅ | N/A | TBD | TBD |
| POST /api/tokens | ✅ | TBD | TBD | TBD |
| GET /api/tokens | ✅ | N/A | TBD | TBD |
| DELETE /api/tokens/{id} | ✅ | TBD | TBD | TBD |
| ... (37 total) | | | | |

**Output:**
- Test collection file (`tests/postman_collection.json`)
- All endpoints return correct status codes
- All endpoints validate input
- All endpoints check authentication

---

## PHASE 2: REVERSE PROXY IMPLEMENTATION (Days 4-5)

### Day 4: Proxy Server Setup

**Objective:** Create separate HTTP server on port 3001 that intercepts traffic

**Step 1: Create Proxy Server (1 hour)**

```go
// internal/proxy/server.go

package proxy

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "log"
)

type ProxyServer struct {
    targetURL *url.URL
    proxy     *httputil.ReverseProxy
}

func NewProxyServer(targetAddress string) (*ProxyServer, error) {
    target, err := url.Parse(targetAddress)
    if err != nil {
        return nil, err
    }
    
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // Intercept requests before forwarding
    proxy.Director = func(req *http.Request) {
        req.URL.Scheme = target.Scheme
        req.URL.Host = target.Host
        req.Host = target.Host
    }
    
    // Intercept responses
    proxy.ModifyResponse = func(res *http.Response) error {
        // We'll log response here
        return nil
    }
    
    return &ProxyServer{
        targetURL: target,
        proxy:     proxy,
    }, nil
}

func (ps *ProxyServer) Start(listenAddr string) error {
    http.HandleFunc("/", ps.handleRequest)
    log.Printf("Proxy listening on %s, forwarding to %s", listenAddr, ps.targetURL.String())
    return http.ListenAndServe(listenAddr, nil)
}

func (ps *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {
    // Step 2: Fingerprint request
    // Step 3: Detect tokens
    // Step 4: Classify risk
    // Step 5: Forward to real server
    // Step 6: Log everything
    
    ps.proxy.ServeHTTP(w, r)
}
```

**Step 2: Fingerprinting Middleware (1 hour)**

```go
// Add to handleRequest()

func (ps *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {
    // Extract request metadata
    ip := r.RemoteAddr // In production: r.Header.Get("X-Forwarded-For")
    userAgent := r.Header.Get("User-Agent")
    language := r.Header.Get("Accept-Language")
    encoding := r.Header.Get("Accept-Encoding")
    
    // Generate fingerprint
    fingerprint := fp.GenerateFingerprint(ip, userAgent, language, encoding)
    
    // Store in request context for later use
    r.Header.Set("X-KAVACH-Fingerprint", fingerprint)
    r.Header.Set("X-KAVACH-IP", ip)
    
    log.Printf("[PROXY] Fingerprint: %s | IP: %s | Path: %s", fingerprint, ip, r.RequestURI)
    
    ps.proxy.ServeHTTP(w, r)
}
```

**Step 3: Token Detection (1 hour)**

```go
// internal/proxy/token_detector.go

func DetectTokenInRequest(r *http.Request) (string, string) {
    // Check URL parameters
    for _, tokenPrefix := range []string{"sk_", "token="} {
        if token := findTokenInURL(r.RequestURI, tokenPrefix); token != "" {
            return token, "url"
        }
    }
    
    // Check Authorization header
    if auth := r.Header.Get("Authorization"); auth != "" {
        if token, found := extractFromAuthHeader(auth); found {
            return token, "api_key"
        }
    }
    
    // Check common credential parameters
    if r.Method == "POST" {
        r.ParseForm()
        for key, values := range r.Form {
            if isCredentialKey(key) && len(values) > 0 {
                if token, found := validateToken(values[0]); found {
                    return token, "form"
                }
            }
        }
    }
    
    return "", ""
}

func isCredentialKey(key string) bool {
    credentialKeys := []string{"token", "key", "password", "credential", "secret", "api_key"}
    for _, k := range credentialKeys {
        if strings.Contains(strings.ToLower(key), k) {
            return true
        }
    }
    return false
}

func findTokenInURL(url, prefix string) string {
    idx := strings.Index(url, prefix)
    if idx == -1 {
        return ""
    }
    start := idx + len(prefix)
    end := start
    for end < len(url) && (isAlphaNumeric(rune(url[end])) || url[end] == '_') {
        end++
    }
    return url[start:end]
}
```

**Output by end of Day 4:**
- ✅ Proxy server listens on :3001
- ✅ Forwards all requests to target (e.g., localhost:3000)
- ✅ Fingerprints every request
- ✅ Detects tokens in URLs, headers, forms
- ✅ Logs all activity

---

### Day 5: Token Triggering & End-to-End Testing

**Step 1: Token Lookup (1 hour)**

```go
// In handleRequest(), after detecting token:

if detectedToken != "" && detectedTokenType != "" {
    // Look up token in database
    token, err := db.GetTokenByValue(detectedToken)
    if err != nil {
        log.Printf("[HONEYPOT] Token not found: %s", detectedToken)
        // Not a honeypot, let it through
        ps.proxy.ServeHTTP(w, r)
        return
    }
    
    // WE FOUND A HONEYPOT TOKEN!
    log.Printf("[HONEYPOT TRIGGERED] Token: %s | IP: %s", token.Name, ip)
    
    // Mark token as triggered
    db.MarkTokenTriggered(token.ID, fingerprint)
    
    // Step 2: Classify risk
    // Step 3: Create trigger event
    // Step 4: Dispatch alerts
}
```

**Step 2: Risk Classification (30 min)**

```go
// In handleRequest(), after finding token:

// Classify the request
riskScore := classifier.Classify7D(
    ip,
    r.Header.Get("User-Agent"),
    r.Method,
    r.RequestURI,
    r.Header,
)

log.Printf("[CLASSIFIER] Risk Score: %.1f/100", riskScore)
```

**Step 3: Create Trigger Event (30 min)**

```go
// In handleRequest(), log the event:

event := &models.TriggerEvent{
    ID:              uuid.New().String(),
    UserID:          token.UserID,
    TokenID:         token.ID,
    AttackerID:      attackerFingerprint, // or create new
    EventType:       "token_access",
    HTTPMethod:      r.Method,
    RequestPath:     r.RequestURI,
    RequestHeaders:  headersJSON,
    IPAddress:       ip,
    UserAgent:       userAgent,
    RiskScore:       riskScore,
    CreatedAt:       time.Now(),
}

db.CreateTriggerEvent(event)
```

**Step 4: Dispatch Alerts (30 min)**

```go
// In handleRequest(), after creating event:

configs, _ := db.GetEnabledAlertConfigs(token.UserID)
for _, config := range configs {
    if riskScore >= config.MinRisk {
        dispatcher.SendAlert(config, event, token)
    }
}
```

**End-to-End Test (30 min):**

```bash
# Terminal 1: Start main server
cd E:\KAVACH_VISION_1
docker-compose up --build

# Terminal 2: Start proxy (once backend is ready)
go run ./cmd/proxy/main.go

# Terminal 3: Register user and create honeypot token
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecurePass123"}'

# Get JWT and create token
curl -X POST http://localhost:3000/api/tokens \
  -H "Authorization: [JWT_TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{"type":"url","name":"Honeypot Test"}'

# Terminal 4: Simulate attacker accessing honeypot via proxy
curl "http://localhost:3001/?token=sk_honeypot_token_value"

# Check:
# 1. Server logs show token detected
# 2. Slack/email alert received
# 3. Dashboard shows event
```

**Output by end of Day 5:**
- ✅ Proxy detects tokens
- ✅ Creates trigger events
- ✅ Classifies risk
- ✅ Dispatches alerts in real-time
- ✅ Full flow tested end-to-end

---

## SUCCESS CRITERIA

### After Phase 1 (Polish & Testing):
- ✅ All 37 endpoints tested
- ✅ No unhandled errors
- ✅ Consistent response format
- ✅ Core services have 80%+ test coverage
- ✅ Postman collection created

### After Phase 2 (Reverse Proxy):
- ✅ Proxy runs on :3001
- ✅ Fingerprints requests
- ✅ Detects honeypot tokens
- ✅ Classifies risk (0-100)
- ✅ Creates events in database
- ✅ Dispatches alerts to Slack
- ✅ Full end-to-end test passes

---

## DELIVERABLES

By end of Week 2:

1. **Code:**
   - `internal/proxy/server.go` (200 lines)
   - `internal/proxy/token_detector.go` (150 lines)
   - `tests/auth_service_test.go` (200 lines)
   - `tests/classifier_test.go` (100 lines)

2. **Tests:**
   - `tests/postman_collection.json` (all 37 endpoints)
   - Integration test results (pass/fail matrix)

3. **Documentation:**
   - `WEEK_2_COMPLETION.md` (what was built, what works, what's next)

---

## DAILY CHECKLIST

### Day 1: Code Review
- [ ] Review all handlers
- [ ] Fix error handling
- [ ] Verify response formats
- [ ] Test 10 random endpoints manually

### Day 2: Unit Tests
- [ ] Write auth service tests
- [ ] Write token generator tests
- [ ] Write classifier tests
- [ ] Run all tests (pass > 90%)

### Day 3: Integration Tests
- [ ] Create Postman collection
- [ ] Test all 37 endpoints
- [ ] Document failures
- [ ] Fix any issues

### Day 4: Proxy Setup
- [ ] Create proxy server on :3001
- [ ] Add fingerprinting
- [ ] Add token detection
- [ ] Test manual flow

### Day 5: End-to-End
- [ ] Implement token lookup
- [ ] Add risk classification
- [ ] Create events
- [ ] Dispatch alerts
- [ ] Full test cycle

---

## TIME BREAKDOWN

| Task | Time | Day |
|------|------|-----|
| Code review & bug fixes | 4h | Day 1 |
| Unit tests (auth, token) | 6h | Day 2 |
| Integration tests | 6h | Day 3 |
| Proxy server setup | 3h | Day 4 |
| Token detection | 2h | Day 4 |
| Token triggering | 2h | Day 5 |
| Risk classification | 1h | Day 5 |
| Events & alerts | 2h | Day 5 |
| End-to-end testing | 2h | Day 5 |
| **Total** | **28h** | **Week 2** |

---

**Ready to start Day 1? Let's do this! 🚀**
