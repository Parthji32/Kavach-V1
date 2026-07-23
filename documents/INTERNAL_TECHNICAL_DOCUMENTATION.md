# KAVACH V1 - Internal Technical Documentation
**For:** Developers, DevOps, System Engineers  
**Last Updated:** Thursday, July 23, 2026  
**Version:** 1.0.0-prod

---

## Architecture Overview

### High-Level Flow
```
User Request
    ↓
Honeypot Detection Middleware (checks token in params/headers/form)
    ↓
If Token Found:
  - Fingerprint generation
  - Attacker lookup/create
  - Event logging
  - Alert dispatch (async)
    ↓
Otherwise:
  - Route to handler
    ↓
Response sent to client
```

### Deployment Architecture
```
Internet
    ↓
Railway (Load Balancer)
    ↓
Docker Container (Alpine 3.20)
    ├── Go Binary (main.go, 65 handlers)
    ├── SQLite Database (/var/data/kavach.db)
    ├── Templates folder (7 HTML pages)
    └── Static assets (CSS, JS)
    ↓
GitHub (Auto-deploy on push)
```

---

## Database Schema (7 Tables)

### users
```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  full_name TEXT,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  jwt_secret TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_login DATETIME
);
```
**Indexes:** `idx_users_email`

### tokens
```sql
CREATE TABLE tokens (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  name TEXT,
  token_value TEXT UNIQUE NOT NULL,
  token_type TEXT NOT NULL, -- url, api_key, document, dns, email
  description TEXT,
  is_active BOOLEAN DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  triggered_at DATETIME,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```
**Indexes:** `idx_tokens_user_id`, `idx_tokens_value`

### attackers
```sql
CREATE TABLE attackers (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  fingerprint TEXT,  -- MD5(IP + UA + lang + enc)
  ip_address TEXT,
  user_agent TEXT,
  os TEXT,
  browser TEXT,
  device_type TEXT,
  risk_score REAL DEFAULT 0,
  risk_level TEXT,  -- low, medium, high, critical
  detection_count INT DEFAULT 1,
  is_known_user BOOLEAN DEFAULT 0,
  is_blocked BOOLEAN DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```
**Indexes:** `idx_attackers_user_id`, `idx_attackers_fingerprint`

### trigger_events
```sql
CREATE TABLE trigger_events (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  token_id TEXT NOT NULL,
  attacker_id TEXT,
  event_type TEXT DEFAULT 'token_accessed',
  request_method TEXT,
  request_path TEXT,
  request_headers TEXT,  -- JSON
  ip_address TEXT,
  user_agent TEXT,
  classification_reason TEXT,
  risk_score REAL,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(token_id) REFERENCES tokens(id),
  FOREIGN KEY(attacker_id) REFERENCES attackers(id)
);
```
**Indexes:** `idx_events_user_id`, `idx_events_attacker_id`

### alert_configs
```sql
CREATE TABLE alert_configs (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  alert_type TEXT NOT NULL,  -- webhook, slack, email
  destination TEXT NOT NULL,  -- URL or email
  is_enabled BOOLEAN DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
```
**Indexes:** `idx_alert_configs_user_id`

### sent_alerts
```sql
CREATE TABLE sent_alerts (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  config_id TEXT NOT NULL,
  event_id TEXT NOT NULL,
  status TEXT,  -- pending, sent, failed
  response_code INT,
  error_message TEXT,
  sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(config_id) REFERENCES alert_configs(id),
  FOREIGN KEY(event_id) REFERENCES trigger_events(id)
);
```

---

## API Endpoints (28 Total)

### Authentication Endpoints
```
POST   /api/auth/register           (Register new user)
POST   /api/auth/login              (Login, returns JWT)
POST   /api/auth/trust-device       (Mark device as trusted)
GET    /api/user/profile            (Get logged-in user's profile)
```

### Token Management Endpoints
```
POST   /api/tokens                  (Create new token)
GET    /api/tokens                  (List user's tokens, paginated)
GET    /api/tokens/:tokenID         (Get single token details)
DELETE /api/tokens/:tokenID         (Deactivate token)
POST   /api/tokens/bulk             (Create multiple tokens)
```

### Dashboard Endpoints
```
GET    /api/dashboard/stats         (Total tokens, attackers, events)
GET    /api/dashboard/attackers     (List high-risk attackers)
GET    /api/dashboard/events        (List trigger events, 24h)
```

### Alert Configuration Endpoints
```
POST   /api/alerts/config           (Create alert config)
GET    /api/alerts/config           (List user's alert configs)
DELETE /api/alerts/config/:configID (Delete alert config)
```

### Page Routes (Not API)
```
GET    /                            (Homepage)
GET    /login                       (Login form)
GET    /signup                      (Registration form)
GET    /app                         (Dashboard - requires auth)
GET    /tokens                      (Token management page)
GET    /how-it-works                (How It Works page)
GET    /pricing                     (Pricing page)
GET    /use-cases                   (Use Cases page)
GET    /faq                         (FAQ page)
GET    /support                     (Support & Contact page)
GET    /health                      (Health check)
```

---

## 7-Dimensional ML Classifier

### Dimensions & Weights
| Dimension | Weight | Range | Scoring |
|-----------|--------|-------|---------|
| IP Reputation | 25% | 0-100 | Private IPs (0), Datacenter (20), Known-bad (100) |
| Request Rate | 15% | 0-100 | <1/min (0), >10/sec (100) |
| Payload Analysis | 15% | 0-100 | Clean (0), SQLi/XSS (50), Command inject (100) |
| Header Fingerprint | 12% | 0-100 | Normal headers (0), Bot UAs/missing headers (100) |
| Behavioral Anomaly | 12% | 0-100 | Normal path (0), Admin path (25), Traversal (100) |
| Geolocation | 12% | 0-100 | Known country (0), VPN/proxy (50), Suspicious country (100) |
| Timing Pattern | 9% | 0-100 | Variable (0), Machine-like intervals (100) |

### Risk Actions
- **0-55:** ALLOW (low risk)
- **55-65:** FLAG (monitor, optional MFA)
- **65-75:** CHALLENGE (require MFA)
- **75+:** BLOCK (HTTP 403)
- **Honeypot Token:** 99 (immediate block + alert)

---

## Device Fingerprinting

### Fingerprint Generation
```go
fingerprint := MD5(IP + User-Agent + Accept-Language + Accept-Encoding)
```
- Consistent for same device + network combo
- Allows tracking of repeat attackers
- Used to correlate multiple events

### Fingerprint Fields Captured
- IP Address
- User-Agent (OS, browser, version)
- Accept-Language (region hints)
- Accept-Encoding (compression preferences)
- TLS fingerprint (JA3 – future enhancement)

---

## Honeypot Token Types

### URL Token
```
Format: https://api.internal.company.com/v1/users?token=sk_<random>
Detection: Scanned in request URL query params
Use: Hidden in documentation, runbooks, comments
```

### API Key Token
```
Format: sk_<64-char-random>
Detection: Bearer token in Authorization header
Use: Committed to Git repos, shared in Slack, env files
```

### Document Token
```
Format: PDF/Doc with embedded tracking token
Detection: Access log + watermark validation
Use: Shared with vendors, left on desktops
```

### DNS Token
```
Format: honeypot-<uuid>.internal.company.com
Detection: DNS query (requires DNS server integration)
Use: Embedded in configs, internal wikis
```

### Email Token
```
Format: honeypot-<uuid>@company.com
Detection: Email forwarding rules capture access
Use: Added to mailing lists, contact forms
```

---

## Alert Dispatch System

### Webhook Delivery
```json
{
  "event_type": "token_accessed",
  "timestamp": "2026-07-23T10:45:30Z",
  "user_id": "uuid",
  "token_id": "uuid",
  "token_value": "sk_...",
  "token_type": "api_key",
  "attacker_id": "uuid",
  "attacker_ip": "192.168.1.105",
  "risk_score": 95,
  "risk_level": "critical",
  "severity": "critical",
  "message": "Honeypot token accessed"
}
```

### Retry Logic
- Attempt 1: Immediate
- Attempt 2: Wait 1 second
- Attempt 3: Wait 2 seconds
- Attempt 4: Wait 4 seconds
- Then: Give up + log failure

### Slack Integration
```
Sends formatted blocks with:
- Risk level (color-coded: red=critical, orange=high, yellow=medium, green=low)
- Attacker IP & location
- Token type
- Event timestamp
- Link to dashboard
```

---

## Development Setup

### Local Build & Run
```powershell
cd E:\KAVACH_VISION_1

# Build locally (requires Go 1.22)
go build -o server.exe .

# Or run via Docker (recommended)
docker-compose up --build

# Watch logs
docker-compose logs -f

# Stop
docker-compose down
```

### Environment Variables
```
PORT=3000
DATABASE_PATH=./kavach.db  (or /var/data/kavach.db in Docker)
JWT_SECRET=<random-32-char-string>
ENVIRONMENT=production
LOG_LEVEL=info
```

### Docker Build
```
Multi-stage build (Alpine 3.20):
- Stage 1: golang:1.22-alpine (build binary with CGO)
- Stage 2: alpine:3.20 (runtime, includes ca-certificates for HTTPS)
```

---

## Performance Tuning

### Database Indexes
- `idx_tokens_user_id` — Fast token lookup by user
- `idx_tokens_value` — Fast token validation
- `idx_attackers_fingerprint` — Detect repeat attackers
- `idx_events_user_id` — Dashboard queries
- `idx_alert_configs_user_id` — Alert routing

### Middleware Caching
- JWT validation cached in `c.Locals()` per request
- Fingerprint generation cached (immutable per request)
- Alert configs cached in memory for 1 hour

### Async Operations
- Alert dispatch: Goroutine (non-blocking)
- Event logging: Synchronous (critical path)
- Attacker profiling: Async (background enrichment — future)

---

## Security Considerations

### Password Hashing
- bcrypt with cost 12 (slow enough for security, fast enough for UX)
- Never store plaintext passwords
- Validate strength: 8+ chars, upper+lower+digit+special

### JWT Security
- HS256 (HMAC-SHA256) with 32-byte secret
- 7-day expiration
- Token stored in HTTP-only cookies (no JS access)
- Rotation possible via device-trust system

### HTTPS
- All deployments use HTTPS (Railway enforces)
- TLS 1.2+ only
- Certificate auto-renewal (Railway handles)

### Input Validation
- Email: RFC 5322 regex
- URLs: Must be https://, 10-2000 chars
- Slack webhooks: Must contain `hooks.slack.com`
- Passwords: 8+ chars, mixed case, special chars
- Token types: Whitelist (url, api_key, document, dns, email)

---

## Known Limitations & TODOs

### Phase 2 (Future Enhancements)
- [ ] Autonomous response system (auto-generate decoys)
- [ ] Threat intel feed (aggregate + anonymize + sell)
- [ ] Industry vertical templates (Finance, Healthcare, E-commerce)
- [ ] Web UI for token placement guidance
- [ ] Advanced analytics dashboard (charts, heatmaps)
- [ ] Machine learning model training (on attack patterns)
- [ ] Managed Detection & Response (MDR) service
- [ ] SIEM/Splunk integration
- [ ] Compliance reporting automation (GDPR, HIPAA, PCI-DSS)

### Known Issues
- Email alerts not implemented (placeholder only)
- DNS token detection requires external DNS server (not bundled)
- No rate limiting on API endpoints (add for production)
- No request logging/audit trail (add for compliance)
- No two-factor authentication (MFA placeholder only)

---

## Deployment Checklist

### Pre-Production
- [ ] Set strong JWT_SECRET in .env
- [ ] Enable HTTPS (Railway does this)
- [ ] Configure database backup (Railway does this)
- [ ] Set up monitoring/alerts (Railway console)
- [ ] Test honeypot detection (webhook.site)
- [ ] Load test with 100+ concurrent users (future)

### Post-Deployment
- [ ] Verify all 7 pages loading
- [ ] Test auth flow (signup → login → dashboard)
- [ ] Create sample token
- [ ] Trigger honeypot (simulate attack)
- [ ] Check alert dispatch
- [ ] Monitor logs for errors
- [ ] Check database size growth

---

## Troubleshooting

### Docker Build Fails
**Error:** `go.sum file not found`
**Fix:** Add `GOSUMDB=off go mod tidy` to Dockerfile

### Database Lock
**Error:** `database is locked`
**Fix:** Restart container (SQLite single-writer limitation)

### JWT Token Expired
**Error:** `401 Unauthorized`
**Fix:** Login again (7-day expiration is by design)

### Honeypot Not Detecting
**Check:**
- Is middleware registered BEFORE routes?
- Is token value in database (not empty string)?
- Is request actually sending the token?

---

## Monitoring & Logging

### Key Metrics to Track
- API response time (target: <100ms)
- Database query time (target: <50ms)
- Alert dispatch latency (target: <5s)
- Honeypot detection accuracy (target: 100%)
- False positive rate (target: 0%)

### Log Levels
- `ERROR` — Critical issues (unrecoverable)
- `WARN` — Potential issues (recoverable)
- `INFO` — Key events (auth, alerts, deployments)
- `DEBUG` — Detailed flow (SQL queries, middleware)

### Key Log Markers
- `[TOKEN-DETECTION]` — Token found in request
- `[HONEYPOT-DETECTED]` — Valid honeypot triggered
- `[ATTACKER-CREATED]` — New attacker profile
- `[EVENT-CREATED]` — Trigger event logged
- `[ALERT-DISPATCH]` — Alert sent to config
- `[WEBHOOK-SUCCESS]` — Alert delivered

---

## Cost Analysis (Railway Free Tier)

| Resource | Limit | Used | Status |
|----------|-------|------|--------|
| Compute | 500 hours/month | ~30/month (dev only) | ✅ Unlimited |
| Memory | 8GB | ~256MB | ✅ OK |
| Disk | 1GB | ~50MB (DB + code) | ✅ OK |
| Bandwidth | Unlimited | ~10MB/day (testing) | ✅ OK |
| Cost | Free | $0 | ✅ Free |

**Expected Cost at Scale:**
- 100 active users: ~$50-100/month (Railway Pro tier)
- 1000 active users: ~$500-1000/month
- 10000 active users: Requires infrastructure upgrade

---

**Next Update:** After first customer deployment or major feature addition.
