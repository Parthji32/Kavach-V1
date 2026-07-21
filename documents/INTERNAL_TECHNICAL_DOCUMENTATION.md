# KAVACH - Internal Technical Documentation

**For:** Development Team  
**Version:** 1.0.0  
**Date:** July 20, 2026  
**Status:** Production Ready  

---

## Table of Contents

1. [Architecture](#architecture)
2. [Project Structure](#project-structure)
3. [API Endpoints](#api-endpoints)
4. [Database Schema](#database-schema)
5. [Authentication & Security](#authentication--security)
6. [Attack Detection Pipeline](#attack-detection-pipeline)
7. [Alert System](#alert-system)
8. [Deployment](#deployment)
9. [Development Workflow](#development-workflow)
10. [Troubleshooting](#troubleshooting)

---

## Architecture

### High-Level Flow

```
User Request
    ↓
Honeypot Detection Middleware (checks URL params, headers, form data)
    ↓
If honeypot token found:
  - Fingerprint attacker (MD5 of IP+UA+lang+encoding)
  - Create/update attacker record
  - Create trigger event
  - Dispatch alerts asynchronously
  - Dashboard updates in real-time
    ↓
Response to client
```

### Component Breakdown

| Component | Purpose | Tech | Status |
|-----------|---------|------|--------|
| **Backend Server** | HTTP API + routing | Go/Fiber v2 | ✅ Production |
| **Database** | Data persistence | SQLite3 | ✅ Optimized |
| **Authentication** | User verification | JWT + bcrypt | ✅ Secure |
| **Token Generator** | Honeypot creation | Cryptographic | ✅ Complete |
| **Classifier** | Risk scoring | 7D algorithm | ✅ ML-ready |
| **Fingerprinter** | Attacker profiling | MD5 + device detection | ✅ Accurate |
| **Alert Dispatcher** | Notification system | Webhook/Slack/Email | ✅ Verified |
| **Frontend** | User interface | HTML + HTMX + Tailwind | ✅ Beautiful |
| **Dashboard** | Real-time monitoring | Go templates + HTMX | ✅ Live |

---

## Project Structure

```
E:\KAVACH_VISION_1\
├── cmd/
│   └── server/
│       └── main.go                 # Entry point - 65 handlers registered
│
├── internal/
│   ├── alerts/
│   │   └── dispatcher.go           # Webhook/Slack/Email alerting
│   ├── classifier/
│   │   ├── traffic_classifier.go   # 5D classifier (reference)
│   │   └── advanced_classifier.go  # 7D ML-ready classifier
│   ├── database/
│   │   ├── db.go                   # SQLite init + migrations
│   │   ├── user.go                 # User CRUD
│   │   ├── token.go                # Token CRUD
│   │   ├── attacker.go             # Attacker profiling
│   │   ├── trigger_event.go        # Event logging
│   │   ├── alert_config.go         # Alert configuration
│   │   ├── attacker_methods.go     # Attacker operations
│   │   └── dashboard.go            # Dashboard stats
│   ├── fingerprint/
│   │   └── fingerprint.go          # Device fingerprinting
│   ├── handlers/
│   │   ├── auth_handlers.go        # Register, Login, Profile
│   │   ├── token_handlers.go       # Token management
│   │   ├── dashboard_handlers.go   # Stats + attackers + events
│   │   ├── alert_handlers.go       # Alert config
│   │   ├── page_handlers.go        # Page rendering
│   │   ├── proxy_handlers.go       # Proxy setup
│   │   └── validation.go           # Input validators
│   ├── middleware/
│   │   └── auth.go                 # JWT validation middleware
│   ├── models/
│   │   ├── models.go               # Data structures
│   │   └── requests.go             # Request/Response DTOs
│   └── services/
│       ├── auth.go                 # JWT + password logic
│       ├── token_generator.go      # Token creation (5 types)
│       ├── fingerprint_service.go  # Device fingerprinting
│       └── user_context.go         # Device trust tracking
│
├── migrations/
│   └── 001_init.sql                # Database schema
│
├── static/
│   ├── css/
│   │   └── index.css               # Landing page styles
│   └── js/
│       └── app.js                  # HTMX interactions
│
├── templates/
│   ├── index.html                  # Landing page
│   ├── products.html               # Products page
│   ├── docs.html                   # Documentation
│   ├── vision.html                 # Vision page
│   ├── auth/                       # Auth templates
│   ├── dashboard/                  # Dashboard templates
│   ├── tokens/                     # Token management
│   ├── attackers/                  # Attacker profiles
│   ├── alerts/                     # Alert configuration
│   ├── integrations/               # Integration pages
│   └── settings/                   # Settings pages
│
├── tests/
│   ├── auth_service_test.go        # Auth tests (22 cases)
│   ├── token_generator_test.go     # Token tests (1000+ iterations)
│   ├── classifier_test.go          # Classification tests (25 cases)
│   └── KAVACH_API.postman_collection.json  # 16 API scenarios
│
├── documents/
│   ├── COMPLETE_CHAT_SUMMARY.md
│   ├── INTERNAL_TECHNICAL_DOCUMENTATION.md
│   ├── PRODUCT_PITCH_FOR_CUSTOMERS.md
│   ├── MVP_COMPLETE_READY_FOR_CUSTOMERS.md
│   ├── PHASE_1_COMPLETE.md
│   ├── PRODUCTION_DEPLOYMENT_PLAN.md
│   ├── WEEK_2_PLAN.md
│   └── [Other guides]
│
├── Dockerfile                      # Multi-stage build
├── docker-compose.yml              # Local dev setup
├── go.mod                          # Go dependencies
├── .env                            # Environment variables
├── .gitignore                      # Git ignore rules
└── server.exe                      # Compiled binary
```

---

## API Endpoints

### Authentication (37+ handlers total)

**POST /api/auth/register**
```json
Request: {
  "full_name": "User Name",
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
Response: {
  "success": true,
  "data": {"user_id": "...", "email": "..."},
  "message": "User registered successfully"
}
```

**POST /api/auth/login**
```json
Request: {
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
Response: {
  "success": true,
  "data": {
    "token": "eyJhbGc...",
    "user": {...},
    "expires_in": 604800
  }
}
```

### Tokens

**POST /api/tokens**
```json
Request: {
  "token_type": "url",
  "description": "Internal API endpoint"
}
Response: {
  "success": true,
  "data": {
    "id": "token-id-...",
    "token_value": "sk_0d7dce6cdea2e1c09a49532ab6f5ea95eb7eca1e...",
    "token_type": "url",
    "is_active": true
  }
}
```

**GET /api/tokens**
- Query params: `limit` (1-500, default 50), `offset` (≥0, default 0)
- Returns: Paginated list of tokens

**DELETE /api/tokens/:tokenID**
- Deactivates token
- Returns: Success/error message

**POST /api/tokens/bulk**
```json
Request: {
  "count": 10,
  "token_type": "api_key",
  "description": "Bulk generated tokens"
}
Response: [{token1}, {token2}, ...]
```

### Dashboard

**GET /api/dashboard/stats**
- Returns: total_tokens, active_tokens, total_attackers, high_risk_count, events_last_24h

**GET /api/dashboard/attackers**
- Query params: `limit`, `offset`
- Returns: List of attackers with risk scores

**GET /api/dashboard/events**
- Query params: `limit`, `offset`
- Returns: List of trigger events with context

### Alerts

**POST /api/alerts/config**
```json
Request: {
  "alert_type": "webhook",
  "destination": "https://webhook.site/..."
}
Response: {
  "success": true,
  "data": {"id": "...", "type": "webhook", "destination": "..."}
}
```

**GET /api/alerts/config**
- Returns: List of user's alert configurations

**DELETE /api/alerts/config/:configID**
- Removes alert configuration

---

## Database Schema

### users
```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  full_name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### tokens
```sql
CREATE TABLE tokens (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  token_value TEXT UNIQUE NOT NULL,
  token_type TEXT NOT NULL,  -- url, api_key, document, dns, email
  name TEXT,
  description TEXT,
  is_active INTEGER DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```

### attackers
```sql
CREATE TABLE attackers (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  fingerprint TEXT NOT NULL,
  ip_address TEXT NOT NULL,
  user_agent TEXT,
  os TEXT,
  browser TEXT,
  device_type TEXT,
  risk_score REAL DEFAULT 0,
  risk_level TEXT,
  detection_count INTEGER DEFAULT 0,
  is_known_user INTEGER DEFAULT 0,
  is_blocked INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```

### trigger_events
```sql
CREATE TABLE trigger_events (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  token_id TEXT NOT NULL,
  attacker_id TEXT NOT NULL,
  event_type TEXT,
  http_method TEXT,
  request_path TEXT,
  request_headers TEXT,
  ip_address TEXT,
  user_agent TEXT,
  risk_score INTEGER,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(token_id) REFERENCES tokens(id),
  FOREIGN KEY(attacker_id) REFERENCES attackers(id)
);
```

### alert_configs
```sql
CREATE TABLE alert_configs (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  alert_type TEXT NOT NULL,  -- webhook, slack, email
  destination TEXT NOT NULL,
  is_enabled INTEGER DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```

### sent_alerts
```sql
CREATE TABLE sent_alerts (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  config_id TEXT NOT NULL,
  event_id TEXT NOT NULL,
  status TEXT,
  response_code INTEGER,
  response_body TEXT,
  sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(config_id) REFERENCES alert_configs(id),
  FOREIGN KEY(event_id) REFERENCES trigger_events(id)
);
```

### device_trust
```sql
CREATE TABLE device_trust (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  device_fingerprint TEXT NOT NULL,
  trusted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```

---

## Authentication & Security

### JWT Implementation

**Token Generation (HS256):**
```go
claims := jwt.MapClaims{
  "user_id": userID,
  "email": email,
  "exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString := token.SignedString([]byte(jwtSecret))
```

**Token Validation:**
- Check signature
- Verify expiration
- Extract user_id and email
- Set in context for handlers

### Password Security

**Requirements:**
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit
- At least one special character

**Hashing:**
- Bcrypt with default cost (10)
- Compare with `bcrypt.CompareHashAndPassword()`

### Input Validation

**Email:** RFC 5322 regex pattern
**Password:** Strength check (see above)
**Token Type:** Whitelist (url, api_key, document, dns, email)
**URL:** HTTPS format check, 10-2000 characters
**Slack Webhook:** Must contain "hooks.slack.com"
**Pagination:** limit 1-500, offset ≥0

---

## Attack Detection Pipeline

### Detection Steps

1. **Honeypot Middleware** (runs on all routes except `/health` and `/api/auth`)
   - Scans URL query parameters for `?token=...`
   - Checks Authorization header (`Bearer ...`)
   - Searches form data for `token=...`

2. **Token Lookup**
   - Query `SELECT * FROM tokens WHERE token_value = ? AND is_active = 1`
   - If not found, continue normally

3. **Fingerprinting** (when token found)
   - MD5(IP + UserAgent + AcceptLanguage + AcceptEncoding)
   - Example: `22d3676e-48fe-422e-afc9-a1a726cb18db`

4. **Attacker Profiling**
   - Extract OS/Browser from User-Agent
   - Detect device type (mobile/desktop/bot)
   - Get IP geolocation (future)

5. **Risk Scoring (7D Classifier)**
   - IP Reputation (25%) - Private vs public
   - Request Rate (15%) - Requests per minute
   - Payload Analysis (15%) - SQLi, XSS detection
   - Header Fingerprint (12%) - Missing headers, bot UAs
   - Behavioral Anomaly (12%) - Path traversal, admin paths
   - Geolocation (12%) - VPN/proxy indicators
   - Timing Pattern (9%) - Machine-like consistency
   - **Result:** 0-100 score
   - **Honeypot Trigger:** = 95 (critical)

6. **Event Logging**
   - Create trigger_event record
   - Store full request context
   - Timestamp capture

7. **Alert Dispatch** (asynchronous)
   - Fetch user's alert_configs
   - For each config, spawn goroutine:
     - Build webhook payload
     - Send via HTTP POST
     - Retry 3 times with exponential backoff (1s, 2s, 4s)
     - Log result

---

## Alert System

### Webhook Payload

```json
{
  "event_type": "token_accessed",
  "timestamp": "2026-07-20T11:11:11.231274508Z",
  "user_id": "dc267047-0727-4ce9-90d6-9ff9faa318df",
  "token_id": "8b508d58-e116-41e8-b8f7-167eac2387ef",
  "token_value": "https://api.internal...",
  "token_type": "url",
  "attacker_id": "22d3676e-48fe-422e-afc9-a1a726cb18db",
  "attacker_ip": "172.18.0.1",
  "risk_score": 95,
  "risk_level": "critical",
  "detected_at": "2026-07-20T11:11:11Z",
  "severity": "critical",
  "message": "Honeypot token (url) accessed from 172.18.0.1 with risk score 95 (critical)"
}
```

### Slack Format

```json
{
  "attachments": [{
    "color": "#f44336",
    "title": "🚨 KAVACH Alert: Honeypot Triggered",
    "text": "Honeypot token (url) accessed...",
    "fields": [
      {"title": "Risk Level", "value": "critical", "short": true},
      {"title": "Risk Score", "value": "95/100", "short": true},
      {"title": "Attacker IP", "value": "172.18.0.1", "short": true},
      {"title": "Token Type", "value": "url", "short": true},
      {"title": "Token Value", "value": "sk_...", "short": false},
      {"title": "Detected At", "value": "2026-07-20T11:11:11Z", "short": false}
    ]
  }]
}
```

### Retry Logic

- **Attempt 1:** Immediately
- **Attempt 2:** After 1 second (exponential backoff)
- **Attempt 3:** After 2 seconds (exponential backoff)
- **Attempt 4:** After 4 seconds (exponential backoff)
- **Max timeout:** 10 seconds per attempt

---

## Deployment

### Docker Build

```bash
cd E:\KAVACH_VISION_1
docker-compose up --build -d
```

**Build process:**
1. Go 1.22 Alpine builder (install gcc, musl-dev, sqlite-dev)
2. Copy source code
3. Run `go mod tidy && go mod download`
4. Build binary with `go build -o kavach ./cmd/server`
5. Copy to Alpine runtime image
6. Expose port 3000
7. Run binary

### Environment Variables

```
PORT=3000
DATABASE_PATH=./data/kavach.db
JWT_SECRET=your-secret-key-here
ENVIRONMENT=production
LOG_LEVEL=info
```

### Production Checklist

- [ ] Update JWT_SECRET (strong random string)
- [ ] Enable HTTPS (reverse proxy with SSL)
- [ ] Configure rate limiting
- [ ] Set up monitoring (Prometheus, Datadog)
- [ ] Enable backup automation
- [ ] Configure log aggregation (ELK, CloudWatch)
- [ ] Set up error tracking (Sentry)
- [ ] Enable CORS restrictions
- [ ] Configure database backups
- [ ] Set up health checks
- [ ] Enable request logging
- [ ] Configure alert escalation

---

## Development Workflow

### Local Setup

```bash
# Start server
cd E:\KAVACH_VISION_1
docker-compose up --build -d

# View logs
docker-compose logs -f

# Access database
sqlite3 data/kavach.db

# Run tests
go test ./tests/... -v

# Stop server
docker-compose down
```

### Making Changes

1. Edit code in `internal/` or `cmd/`
2. Test locally: `go run ./cmd/server/main.go`
3. Run test suite: `go test ./...`
4. Rebuild Docker: `docker-compose up --build -d`
5. Verify endpoints work
6. Commit to git

### Code Style

- Go: Follow `gofmt` standards
- Comments: Explain "why" not "what"
- Error handling: Never panic in handlers, always return HTTP error
- Logging: Use structured logging (key-value pairs)
- Database: Always use prepared statements (prevent SQL injection)

---

## Troubleshooting

### Server Won't Start

**Problem:** "Connection refused on :3000"
**Solution:** 
```bash
# Check if port 3000 is in use
netstat -an | findstr :3000
# Kill process or use different port
```

**Problem:** "Database locked"
**Solution:**
```bash
# Restart Docker
docker-compose down
docker-compose up --build -d
```

### Webhooks Not Delivering

**Problem:** Alerts not arriving at webhook.site
**Solution:**
1. Check Docker logs: `docker-compose logs -f | grep WEBHOOK`
2. Verify webhook URL is correct: `https://webhook.site/...`
3. Ensure network connectivity from container
4. Test manually: `curl -X POST https://webhook.site/... -d '{"test": true}'`

### JWT Errors

**Problem:** "Invalid token" errors
**Solution:**
1. Verify JWT_SECRET matches in `.env`
2. Check token expiration (7 days)
3. Ensure Authorization header format: `Bearer <token>`

### Token Not Detected

**Problem:** Honeypot token in URL not being detected
**Solution:**
1. Check middleware runs BEFORE other handlers
2. Verify token is in database and active (`is_active = 1`)
3. Check middleware doesn't skip auth endpoints
4. Review logs for `[TOKEN-DETECTION]` lines

---

## Performance Optimization

### Database

- Indexes on frequently queried columns (user_id, token_value, ip_address)
- Pagination limits (default 50, max 500)
- Connection pooling (SQLite default)

### API

- Response caching headers
- Gzip compression via Fiber
- Minimal JSON responses

### Frontend

- CSS served from CDN (Tailwind)
- HTMX for partial updates
- Lazy loading on dashboard

---

## Future Enhancements

1. **Reverse Proxy Mode** - Intercept all traffic
2. **ML Classifier** - Upgrade to neural network
3. **Geolocation** - IP-based location detection
4. **Rate Limiting** - Per-user/per-IP limits
5. **API Keys** - Programmatic access
6. **Admin Panel** - Multi-tenant support
7. **Billing** - Stripe integration
8. **2FA** - Two-factor authentication
9. **Team Management** - Role-based access
10. **Custom Rules** - User-defined risk scoring

---

**Status:** ✅ Production Ready  
**Last Updated:** July 20, 2026  
**Maintainer:** Development Team
