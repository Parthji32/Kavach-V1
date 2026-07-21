# Days 3-5: Reverse Proxy Implementation Guide

**Timeline:** July 18-19, 2026  
**Phase:** Week 2, Days 3-5 - Core Feature Implementation  
**Objective:** Build reverse proxy with honeypot detection and alert triggering

---

## 📋 WHAT'S BEEN CREATED

### New Files (Proxy Infrastructure)
1. **cmd/proxy/main.go** (380 lines)
   - Reverse proxy server on port 3001
   - Token detection engine
   - Request fingerprinting
   - Attacker correlation
   - Alert dispatching

2. **internal/database/attacker_methods.go** (86 lines)
   - GetAttackerByFingerprint()
   - UpdateAttacker()
   - CreateAttacker()

3. **internal/services/fingerprint_service.go** (75 lines)
   - Fingerprint generation (MD5 hash)
   - User-Agent parsing
   - Device type detection

4. **internal/alerts/alert_dispatcher_enhanced.go** (195 lines)
   - Webhook alert sending
   - Slack integration
   - Email placeholder
   - Payload construction

---

## 🏗️ ARCHITECTURE

### Request Flow (When Attacker Hits Honeypot)

```
1. ATTACKER → PROXY (port 3001)
   └─ Request with honeypot token

2. FINGERPRINT GENERATION
   └─ Extract: IP, User-Agent, Language, Encoding
   └─ Hash: MD5(IP | UA | Lang | Enc)

3. TOKEN DETECTION
   └─ Check URL parameters
   └─ Check Authorization header
   └─ Check form data
   └─ Find: "sk_*" tokens

4. TOKEN LOOKUP
   └─ Query database for token value
   └─ If found: Continue to attacker analysis
   └─ If not found: Pass through to real server

5. ATTACKER CORRELATION
   └─ Lookup attacker by fingerprint
   └─ If new: Create attacker profile
   └─ If known: Update last_seen, increment count

6. RISK CLASSIFICATION
   └─ 7D classifier analyzes request
   └─ Scores: 0-100 risk
   └─ Actions: Allow, Flag, Challenge, Block

7. EVENT CREATION
   └─ Store trigger event in database
   └─ Record: token, attacker, timestamp, risk

8. ALERT DISPATCH
   └─ Find user's alert configs
   └─ Send webhook/Slack/email
   └─ Async (non-blocking)

9. BLOCK or FORWARD
   └─ If risk > 75: Return 403 Forbidden
   └─ Else: Forward to target server
```

---

## 🚀 HOW TO RUN (Day 3-5)

### Step 1: Prepare Environment

**Update .env:**
```env
# Server (main KAVACH API)
PORT=3000
DATABASE_PATH=./data/kavach.db
JWT_SECRET=[REDACTED_PASSWORD]

# Proxy
PROXY_LISTEN=:3001
PROXY_TARGET=http://localhost:3000
```

### Step 2: Build Proxy Server

```bash
cd E:\KAVACH_VISION_1
go build -o proxy.exe ./cmd/proxy/main.go
```

### Step 3: Run Both Servers (in separate terminals)

**Terminal 1: Main API Server**
```bash
cd E:\KAVACH_VISION_1
docker-compose up --build
# Or: go run ./cmd/server/main.go
```

**Terminal 2: Reverse Proxy**
```bash
cd E:\KAVACH_VISION_1
./proxy.exe
# Or: go run ./cmd/proxy/main.go
```

### Step 4: Test End-to-End

```bash
# Terminal 3: Testing

# 1. Register user (via main server)
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecurePass123!"}'
# Get user_id from response

# 2. Login to get JWT
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecurePass123!"}'
# Save JWT_TOKEN from response

# 3. Create honeypot token
curl -X POST http://localhost:3000/api/tokens \
  -H "Authorization: JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"url","description":"Honeypot Test"}'
# Get token value: sk_xxx...

# 4. **Simulate attacker hitting honeypot via PROXY**
curl "http://localhost:3001/?token=sk_xxx" \
  -v
# Watch proxy logs for token detection

# 5. Check dashboard stats (via main server)
curl http://localhost:3000/api/dashboard/stats \
  -H "Authorization: JWT_TOKEN"
# Should show: 1 attacker detected, 1 event created
```

---

## 📊 EXPECTED OUTPUT

### When Attacker Hits Honeypot (Proxy Logs)

```
[PROXY] Incoming request | IP: 203.0.113.42 | Method: GET | Path: /?token=sk_abc123
[FINGERPRINT] Generated: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
[TOKEN DETECTED] Token: sk_abc123 (in url)
[TOKEN LOOKUP] Found: "Database Backup"
[HONEYPOT TRIGGERED] Token: Database Backup | User: usr_xyz
[CLASSIFIER] Risk Score: 78.5
[NEW ATTACKER] ID: atk_12345 | IP: 203.0.113.42
[EVENT CREATED] Event ID: evt_xyz123
[ALERT SENT] webhook → https://webhook.example.com/alerts
[ALERT SENT] slack → #security
[BLOCKING] Risk score too high: 78.5
```

### Alert Example (Webhook Payload)

```json
{
  "event_id": "evt_xyz123",
  "timestamp": 1689707400,
  "severity": "high",
  "token": {
    "id": "tok_xyz",
    "name": "Database Backup",
    "type": "url",
    "created_at": 1689600000
  },
  "attacker": {
    "id": "atk_12345",
    "fingerprint": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "ip_address": "203.0.113.42",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "risk_score": 78.5,
    "detection_count": 1,
    "first_seen": 1689707400,
    "last_seen": 1689707400
  },
  "request": {
    "method": "GET",
    "path": "/?token=sk_abc123",
    "timestamp": 1689707400
  }
}
```

---

## 🧪 TESTING CHECKLIST

### Day 3: Proxy Setup
- [ ] Proxy server builds without errors
- [ ] Proxy listens on :3001
- [ ] Proxy forwards requests to :3000
- [ ] Fingerprinting works (logs show fingerprint)
- [ ] Token detection works (logs show token)

### Day 4: Token Triggering
- [ ] Token lookup works (finds token in DB)
- [ ] Attacker profile created (new attacker)
- [ ] Attacker profile updated (known attacker)
- [ ] Risk classification works (score 0-100)
- [ ] Events saved to database

### Day 5: Alerts & End-to-End
- [ ] Webhook alerts sent successfully
- [ ] Slack alerts sent successfully
- [ ] Dashboard shows new events
- [ ] Dashboard shows new attackers
- [ ] Full flow works: token → alert → dashboard

---

## 📝 DEBUGGING GUIDE

### Issue: "Token not detected"
**Check:**
- Is token in correct format? (should start with "sk_")
- Is token in URL? (e.g., `?token=sk_xxx`)
- Check proxy logs for "TOKEN DETECTED" message

**Solution:**
- Verify token value from DB
- Manually test: `curl "http://localhost:3001/?token=sk_abc123"`

### Issue: "Honeypot not triggered"
**Check:**
- Did proxy find the token? (logs: "TOKEN DETECTED")
- Is token in database? (logs: "TOKEN LOOKUP")
- Is token active? (check DB: `is_active = 1`)

**Solution:**
- Create new token if necessary
- Verify token is in database

### Issue: "Alerts not sent"
**Check:**
- Did event get created? (logs: "EVENT CREATED")
- Are alert configs set up? (check DB: alert_configs)
- Is webhook URL valid? (logs: "ALERT SENT" or "ALERT ERROR")

**Solution:**
- Create alert config in database
- Test webhook URL manually
- Check proxy error logs

### Issue: "Port already in use"
**Problem:** Another process using port 3001
**Solution:**
```bash
# Find process using port
netstat -ano | findstr :3001

# Kill process (replace PID)
taskkill /PID 12345 /F

# Or use different port
set PROXY_LISTEN=:3002
```

---

## ✅ SUCCESS CRITERIA (When Complete)

### Day 3
- ✅ Proxy server running
- ✅ Fingerprinting working
- ✅ Token detection working
- ✅ Request forwarding working

### Day 4
- ✅ Token lookup working
- ✅ Attacker correlation working
- ✅ Risk classification working
- ✅ Events created in DB

### Day 5
- ✅ Alerts dispatched via webhook
- ✅ Alerts dispatched via Slack
- ✅ Dashboard updated in real-time
- ✅ Full end-to-end test passing
- ✅ Production ready

---

## 🎯 NEXT AFTER DAY 5

Once proxy is complete, you'll have:
- ✅ Running honeypot platform
- ✅ Real-time attack detection
- ✅ Automatic alert triggering
- ✅ Complete event logging
- ✅ Attacker profiling

**Remaining (for Week 3):**
- Dashboard UI (copy from archive)
- Real-time updates via HTMX
- Admin panel
- User management

---

## 📞 REFERENCE

**Port Map:**
- :3000 = Main API (Kavach backend)
- :3001 = Proxy (Attacker entry point)
- Real server = Behind proxy

**Key Files:**
- `cmd/proxy/main.go` - Proxy server
- `internal/database/attacker_methods.go` - Attacker ops
- `internal/alerts/alert_dispatcher_enhanced.go` - Alert sending

**Database Tables:**
- `tokens` - Honeypot tokens
- `attackers` - Attacker profiles
- `trigger_events` - When tokens accessed
- `alert_configs` - User alert settings

---

*Implementation Guide: Days 3-5 Proxy Development*  
*Status: Ready to Build*  
*Next: Execute steps above*
