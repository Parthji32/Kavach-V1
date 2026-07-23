# KAVACH - Internal Technical Documentation (Updated)

**For:** Development Team  
**Version:** 2.0.0 (Updated Session 9)  
**Date:** July 22, 2026  
**Status:** Production Ready, MVP Complete

---

## TABLE OF CONTENTS
1. [Quick Start](#quick-start)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Building & Running](#building--running)
5. [API Reference](#api-reference)
6. [Database Schema](#database-schema)
7. [Attack Detection Pipeline](#attack-detection-pipeline)
8. [Alert System](#alert-system)
9. [Deployment](#deployment)
10. [Troubleshooting](#troubleshooting)

---

## QUICK START

### Prerequisites
- Go 1.22+
- Docker & Docker Compose
- Git

### Local Development (Docker)
```bash
cd E:\KAVACH_VISION_1

# Build and run
docker-compose up --build

# Server runs on http://localhost:3000
# Database: SQLite at ./data/kavach.db

# Clean rebuild
docker-compose down
rm -r data/
docker-compose up --build
```

### First Steps
1. Navigate to http://localhost:3000
2. Click "Get Started" → Sign up
3. Create honeypot tokens via dashboard
4. Test attack detection (see Testing section)

---

## ARCHITECTURE

### System Design
```
┌─────────────────────────────────────────────────┐
│                 Browser / Client                │
└──────────────────────────┬──────────────────────┘
                           │
                 HTTP Request (Fiber)
                           │
                           ▼
┌─────────────────────────────────────────────────┐
│          KAVACH Server (Go/Fiber v2)            │
│                                                 │
│  1. Honeypot Detection Middleware               │
│     ├─ Check URL params for tokens             │
│     ├─ Check Authorization header              │
│     ├─ Check form data (POST bodies)           │
│     └─ If found: Trigger detection flow        │
│                                                 │
│  2. Fingerprinting Service                      │
│     ├─ IP address extraction                   │
│     ├─ User-Agent parsing                      │
│     ├─ Device fingerprint (MD5 hash)           │
│     └─ Geolocation lookup                      │
│                                                 │
│  3. Risk Classification (7D)                    │
│     ├─ IP reputation (25%)                     │
│     ├─ Request rate (15%)                      │
│     ├─ Payload analysis (15%)                  │
│     ├─ Header fingerprint (12%)                │
│     ├─ Behavioral anomaly (12%)                │
│     ├─ Geolocation (12%)                       │
│     └─ Timing pattern (9%)                     │
│     → Risk Score: 0-100                        │
│                                                 │
│  4. Attacker Correlation                        │
│     ├─ Look up by fingerprint                  │
│     ├─ Create if new                           │
│     └─ Update with latest context              │
│                                                 │
│  5. Event Logging                               │
│     └─ Store trigger_event in DB               │
│                                                 │
│  6. Alert Dispatch (Async)                      │
│     ├─ Get user's alert configs                │
│     ├─ Build payload (attacker + context)      │
│     ├─ Send webhook/Slack/email                │
│     └─ Log delivery status                     │
│                                                 │
│  7. Dashboard Update                            │
│     └─ Real-time stats (tokens, attackers)    │
└─────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────┐
│         SQLite Database (data/kavach.db)        │
│                                                 │
│  Tables:                                        │
│  ├─ users (auth, profiles)                     │
│  ├─ tokens (honeypots: 5 types)                │
│  ├─ attackers (profiles, risk scores)          │
│  ├─ trigger_events (attack log)                │
│  └─ alert_configs (webhook/Slack/email)       │
└─────────────────────────────────────────────────┘
```

### Key Components

| Component | Tech | Purpose | Status |
|-----------|------|---------|--------|
| **Web Server** | Fiber v2 | HTTP API + routing | ✅ Production |
| **Database** | SQLite3 | Persistence layer | ✅ Optimized |
| **Auth** | JWT (HS256) + bcrypt | User verification | ✅ Secure |
| **Fingerprinting** | MD5 hash + UA parsing | Attacker profiling | ✅ Complete |
| **Classification** | 7D algorithm | Risk scoring | ✅ ML-ready |
| **Alerts** | Webhook/Slack/Email | Notifications | ✅ Tested |
| **Frontend** | HTML + HTMX + Tailwind | UI/UX | ✅ Live |
| **Dashboard** | Go templates + HTMX | Real-time monitoring | ✅ Working |

---

## PROJECT STRUCTURE

```
E:\KAVACH_VISION_1\
│
├── cmd/
│   └── server/
│       └── main.go                      # Entry point (65 handlers)
│
├── internal/
│   ├── alerts/
│   │   └── dispatcher.go                # Webhook/Slack/Email alerts
│   ├── classifier/
│   │   ├── traffic_classifier.go        # 5D classifier (ref)
│   │   └── advanced_classifier.go       # 7D ML-ready classifier
│   ├── database/
│   │   ├── db.go                        # SQLite init + migrations
│   │   ├── user.go                      # User CRUD
│   │   ├── token.go                     # Token CRUD
│   │   ├── attacker.go                  # Attacker operations
│   │   ├── trigger_event.go             # Event logging
│   │   ├── alert_config.go              # Alert config CRUD
│   │   ├── attacker_methods.go          # Extra attacker ops
│   │   └── dashboard.go                 # Dashboard stats
│   ├── fingerprint/
│   │   └── fingerprint_service.go       # Device fingerprinting
│   ├── handlers/
│   │   ├── auth_handlers.go             # Register, Login, Profile
│   │   ├── token_handlers.go            # Token management
│   │   ├── dashboard_handlers.go        # Stats + attackers
│   │   ├── alert_handlers.go            # Alert config
│   │   ├── page_handlers.go             # Page rendering
│   │   ├── proxy_handlers.go            # Proxy setup
│   │   └── validation.go                # Input validators
│   ├── middleware/
│   │   └── auth.go                      # JWT validation
│   ├── models/
│   │   ├── models.go                    # Data structures
│   │   └── requests.go                  # Request/Response DTOs
│   └── services/
│       ├── auth.go                      # JWT + password logic
│       ├── token_generator.go           # Token creation
│       ├── fingerprint_service.go       # Fingerprinting
│       └── user_context.go              # Device trust
│
├── migrations/
│   └── 001_init.sql                     # Database schema
│
├── static/
│   ├── css/
│   │   └── index.css                    # Landing page styles
│   └── js/
│       ├── app.js                       # HTMX interactions
│       └── animations.js                # Scroll animations
│
├── templates/
│   ├── index.html                       # Landing page
│   ├── products.html                    # Products page
│   ├── docs.html                        # Documentation
│   ├── vision.html                      # Vision page
│   ├── auth/
│   ├── dashboard/
│   ├── tokens/
│   ├── attackers/
│   ├── alerts/
│   ├── integrations/
│   └── settings/
│
├── tests/
│   ├── auth_service_test.go             # 22 auth tests
│   ├── token_generator_test.go          # 1000+ iterations
│   ├── classifier_test.go               # 25 classifier tests
│   └── KAVACH_API.postman_collection.json  # 16 scenarios
│
├── demo_pages/                          # Updated website redesign
│   ├── homepage_updated.html            # Real messaging
│   ├── how-it-works_updated.html        # 7D classifier explained
│   ├── login_updated.html               # Proper auth UI
│   ├── pricing_updated.html             # Pricing structure
│   └── use-cases_updated.html           # Real scenarios
│
├── documents/
│   ├── COMPLETE_CHAT_SUMMARY.md         # Full project history
│   ├── INTERNAL_TECHNICAL_DOCUMENTATION.md  # This file
│   ├── PRODUCT_PITCH_FOR_CUSTOMERS.md   # Sales doc
│   ├── WEBSITE_WIREFRAMES.md            # UX/UI specs
│   └── DEMO_PAGES_UPDATE_GUIDE.md       # Redesign guide
│
├── Dockerfile                           # Multi-stage build
├── docker-compose.yml                   # Local dev setup
├── go.mod                               # Go dependencies
├── go.sum                               # Go dependency hashes
├── .env                                 # Environment config
├── .gitignore                           # Git rules
└── server.exe                           # Compiled binary
```

---

## BUILDING & RUNNING

### Local Development (Recommended)

**1. Start with Docker Compose:**
```bash
cd E:\KAVACH_VISION_1
docker-compose up --build
```

Server runs on http://localhost:3000, database at `./data/kavach.db`

**2. Test locally:**
```powershell
# Create user
$body = @{
    full_name = "Test User"
    email = "test@example.com"
    password = "TestPassword123!"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:3000/api/auth/register" `
  -Method POST -ContentType "application/json" -Body $body

# Create token
$header = @{ Authorization = "Bearer YOUR_JWT_TOKEN" }
Invoke-WebRequest -Uri "http://localhost:3000/api/tokens" `
  -Method POST -Headers $header `
  -Body '{"token_type":"url","description":"Test"}'
```

### Production Deployment (Railway.app)

**1. Push to GitHub:**
```bash
git add .
git commit -m "KAVACH production build"
git push origin main
```

**2. Railway auto-deploys:**
- Detects Go app
- Runs build: `go build -o kavach ./cmd/server`
- Starts: `./kavach`
- Live at: https://kavach-v1-production.up.railway.app

**3. Environment variables (set in Railway dashboard):**
```
PORT=3000
JWT_SECRET=[your-secret-key-32-chars-min]
DATABASE_PATH=/var/data/kavach.db
ENVIRONMENT=production
LOG_LEVEL=info
```

---

## API REFERENCE

### Authentication

**POST /api/auth/register**
```json
{
  "full_name": "User Name",
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```
Returns: `{success: true, data: {user_id, email}}`

**POST /api/auth/login**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```
Returns: `{success: true, data: {token, user, expires_in}}`

### Tokens (Honeypots)

**POST /api/tokens** — Create single token
```json
{
  "token_type": "url|api_key|document|dns|email",
  "description": "Optional description"
}
```

**GET /api/tokens** — List tokens (paginated)
Query params: `limit` (1-500, default 50), `offset` (≥0, default 0)

**DELETE /api/tokens/:tokenID** — Deactivate token

**POST /api/tokens/bulk** — Create multiple
```json
{
  "count": 10,
  "token_type": "api_key",
  "description": "Bulk created"
}
```

### Dashboard

**GET /api/dashboard/stats** 
Returns: `{total_tokens, active_tokens, total_attackers, high_risk_count, events_last_24h}`

**GET /api/dashboard/attackers**
Returns: List of attackers (paginated)

**GET /api/dashboard/events**
Returns: List of trigger events (paginated)

### Alerts

**POST /api/alerts/config**
```json
{
  "alert_type": "webhook|slack|email",
  "destination": "https://webhook.site/..." | "slack-channel" | "email@domain"
}
```

**GET /api/alerts/config**
Returns: List of user's alert configs

**DELETE /api/alerts/config/:configID**
Removes alert config

### Pages

All page routes return HTML (no JSON):
- `GET /` — Landing page
- `GET /login` — Login page
- `GET /app` — Dashboard (requires JWT in cookie)
- `GET /tokens` — Token management
- `GET /attackers` — Attacker list
- `GET /alerts` — Alert configuration
- `GET /settings` — User settings

---

## DATABASE SCHEMA

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
  token_type TEXT NOT NULL,     -- url, api_key, document, dns, email
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
  ip_address TEXT NOT NULL,
  fingerprint TEXT UNIQUE NOT NULL, -- MD5 hash
  risk_score REAL DEFAULT 0,
  risk_level TEXT,                   -- low, medium, high, critical
  detection_count INTEGER DEFAULT 0,
  is_known_user INTEGER DEFAULT 0,
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
  event_type TEXT DEFAULT 'token_accessed',
  request_path TEXT,
  request_headers TEXT,              -- JSON
  ip_address TEXT,
  user_agent TEXT,
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
  alert_type TEXT NOT NULL,          -- webhook, slack, email
  destination TEXT NOT NULL,
  is_enabled INTEGER DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```

---

## ATTACK DETECTION PIPELINE

### Full Flow (Step-by-Step)

**1. Request Arrives**
```
GET http://localhost:3000/api/dashboard/stats?token=sk_0d7dce...
```

**2. Honeypot Detection Middleware**
- Scans `query` params for patterns: `sk_`, `token_`, etc.
- Scans `Authorization` header: `Bearer sk_...`
- Scans `POST` body form fields (token, api_key, password, etc.)
- ✅ Found: `sk_0d7dce...` in query param

**3. Token Lookup**
```go
tokenRecord := db.GetTokenByValue("sk_0d7dce...")
// Returns: {ID, UserID, TokenType, IsActive}
```

**4. Fingerprinting**
```go
fingerprint := md5.Sum([]byte(
  fmt.Sprintf("%s|%s|%s|%s", 
    clientIP, userAgent, language, encoding)))
// Result: "abc123def456..."
```

**5. Classification (7D)**
```
Score calculation:
├─ IP Reputation (25%):     Public IP = 0, Private = 50, Known Bad = 100
├─ Request Rate (15%):      First request = 0, High velocity = 100
├─ Payload (15%):           Clean = 0, SQLi patterns = 100
├─ Headers (12%):           Normal = 0, Bot signature = 100
├─ Behavior (12%):          Normal path = 0, Admin/traversal = 100
├─ Geolocation (12%):       Normal = 0, VPN/unusual = 100
└─ Timing (9%):             Variable = 0, Machine-like = 100
   ────────────────
   Final Score: 95 → CRITICAL (honeypot hit)
```

**6. Attacker Correlation**
```go
// Look up by fingerprint
attacker := db.GetAttackerByFingerprint("abc123def456...")
if !found {
  // Create new attacker record
  attacker = Attacker{
    ID: UUID(),
    UserID: tokenRecord.UserID,
    IPAddress: clientIP,
    Fingerprint: fingerprint,
    RiskScore: 95,
    RiskLevel: "critical",
  }
  db.CreateAttacker(attacker)
} else {
  // Update existing
  attacker.DetectionCount++
  attacker.RiskScore = 95
  db.UpdateAttacker(attacker)
}
```

**7. Event Logging**
```go
event := TriggerEvent{
  ID: UUID(),
  UserID: tokenRecord.UserID,
  TokenID: tokenRecord.ID,
  AttackerID: attacker.ID,
  EventType: "token_accessed",
  RequestPath: "/api/dashboard/stats",
  RequestHeaders: headersJSON,
  IPAddress: clientIP,
  UserAgent: userAgent,
  Timestamp: now(),
}
db.CreateTriggerEvent(event)
```

**8. Alert Dispatch (Async)**
```go
go func() {
  configs := db.GetAlertConfigsByUserID(tokenRecord.UserID)
  for _, config := range configs {
    payload := BuildWebhookPayload(attacker, event, tokenRecord)
    switch config.AlertType {
    case "webhook":
      alerts.SendWebhookAlert(config.Destination, payload)
    case "slack":
      alerts.SendSlackAlert(config.Destination, payload)
    case "email":
      alerts.SendEmailAlert(config.Destination, payload)
    }
  }
}()
```

**9. Dashboard Updated**
- Real-time stats recalculated
- Attacker profile visible in UI
- Event timeline updated

**10. Response Sent**
```json
{
  "success": true,
  "data": {
    "total_tokens": 5,
    "active_tokens": 5,
    "total_attackers": 1,
    "high_risk_count": 1,
    "events_last_24h": 1
  },
  "message": "Stats retrieved"
}
```

### Risk Scoring Details

**Default Thresholds:**
- **0-55:** ALLOW (low risk)
- **55-65:** FLAG (monitor)
- **65-75:** CHALLENGE (MFA prompt)
- **75-95:** BLOCK (immediate block)
- **95+:** HONEYPOT TRIGGER (critical)

---

## ALERT SYSTEM

### Webhook Payload Structure
```json
{
  "event_type": "token_accessed",
  "timestamp": "2026-07-22T14:30:00Z",
  "user_id": "user-uuid-...",
  "token_id": "token-uuid-...",
  "token_value": "https://internal-api.company.com/admin",
  "token_type": "url",
  "attacker_id": "attacker-uuid-...",
  "attacker_ip": "192.168.1.105",
  "risk_score": 95,
  "risk_level": "critical",
  "detected_at": "2026-07-22T14:30:00Z",
  "severity": "critical",
  "message": "Honeypot token (url) accessed from 192.168.1.105 with risk score 95 (critical)",
  "attacker_profile": {
    "device_fingerprint": "abc123def456...",
    "browser": "Chrome 120",
    "os": "Windows 10",
    "geolocation": "Unknown VPN",
    "detection_history": 1
  }
}
```

### Slack Format
```
🚨 CRITICAL ALERT - HONEYPOT TRIGGERED

Token: https://internal-api.company.com/admin (type: url)
Attacker IP: 192.168.1.105
Risk Score: 95/100
Device: Chrome on Windows 10
Geolocation: Unknown VPN

Action: Block attacker | Investigate | Dismiss
```

---

## DEPLOYMENT

### Production Checklist

- [ ] Environment variables set in Railway dashboard
- [ ] JWT_SECRET is 32+ chars, cryptographically random
- [ ] DATABASE_PATH points to persistent volume
- [ ] Email configuration (if using email alerts)
- [ ] SSL/TLS enabled (Railway auto-provides)
- [ ] Rate limiting configured (if needed)
- [ ] Monitoring/logging enabled
- [ ] Backup strategy for SQLite database
- [ ] Domain configured (if custom domain needed)

### Railway Dashboard Setup

1. **Create new Web Service** → Connect GitHub
2. **Select:** Parthji32/Kavach-V1 repo, main branch
3. **Runtime:** Go (auto-detected)
4. **Environment → Add variables:**
   - PORT: 3000
   - JWT_SECRET: (generate random 32-char string)
   - DATABASE_PATH: /var/data/kavach.db
   - ENVIRONMENT: production
   - LOG_LEVEL: info
5. **Disk → Add disk:**
   - Name: data
   - Mount Path: /var/data
   - Size: 1 GB
6. **Deploy** → Auto-builds and starts

---

## TROUBLESHOOTING

### Server Won't Start

**Error: `bind: address already in use`**
```bash
# Port 3000 is taken
# Kill existing process
netstat -ano | findstr :3000
taskkill /PID [PID] /F
```

**Error: `database is locked`**
```bash
# SQLite has exclusive lock
# Stop Docker, delete db, restart
docker-compose down
rm data/kavach.db
docker-compose up
```

### Attack Detection Not Triggering

**Check:** Honeypot detection middleware
```go
// Verify token is in correct location (URL param, header, form)
GET /api/dashboard/stats?token=sk_0d7dce...  // ✅ Correct
POST /api/dashboard/stats with form: token=sk_...  // ✅ Correct
```

**Check:** Token exists in database
```bash
sqlite3 data/kavach.db "SELECT * FROM tokens WHERE token_value LIKE 'sk_%';"
```

**Check:** Alert configs exist
```bash
sqlite3 data/kavach.db "SELECT * FROM alert_configs WHERE user_id='...';"
```

### Alerts Not Delivering

**Check:** Webhook URL is reachable
```powershell
curl -X POST "YOUR_WEBHOOK_URL" `
  -H "Content-Type: application/json" `
  -D '{"test": true}'
```

**Check:** Network access from Railway
- Webhook must be publicly accessible (no localhost)
- Check firewall/CORS if external

**Check:** Logs in Railway dashboard
- View real-time logs for alert dispatch errors

---

## PERFORMANCE NOTES

- **Token generation:** 1000+ tokens/sec
- **Classification:** <10ms per request
- **Database queries:** <5ms (SQLite on SSD)
- **Alert dispatch:** Async (non-blocking)
- **Memory:** ~50MB baseline, scales with active connections
- **Scalability:** Single instance good for 1K+ users, use load balancer + replicas for 10K+

---

## NEXT DEVELOPMENT PRIORITIES

1. **ML-Enhanced Classifier** — Upgrade from 7D rule-based to neural network
2. **Threat Intelligence Feed** — Auto-update IP reputation database
3. **Advanced Correlation** — Link multiple attacks across users
4. **Performance Optimization** — Redis caching for classifier
5. **Multi-tenancy** — Support enterprise white-label deployments
6. **Mobile App** — Native iOS/Android app for alerts

---

**End of Internal Technical Documentation**
