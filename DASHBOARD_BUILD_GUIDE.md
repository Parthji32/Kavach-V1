# KAVACH Dashboard - Build & Test Guide

## Quick Start (After Templates Integration)

### Build & Run

```powershell
cd E:\KAVACH_VISION_1
docker-compose down          # Clean previous build (optional)
docker-compose up --build    # Build and run
```

Wait for: `Server: http://localhost:3000`

### Test Flow

**1. Sign Up**
```
URL: http://localhost:3000/signup
Email: test@example.com
Password: Test@12345 (must have uppercase, lowercase, digit, special char)
Full Name: Test User
```

**2. Login**
```
URL: http://localhost:3000/login
Email: test@example.com
Password: Test@12345
```

**3. Dashboard**
```
After login, redirects to: http://localhost:3000/app
You should see:
- Total Tokens count
- Active Tokens count
- Total Attackers count
- High-Risk Attackers count
- Events (last 24 hours)
```

**4. Create Token**
```
Click: + New Token (or go to /tokens/new)
Token Type: api_key
Description: Test API Key
Click: Create Token
Token value will be displayed (sk_xxxx...)
```

**5. Create Alert Config**
```
Go to: http://localhost:3000/alerts
Click: + New Alert
Webhook URL: https://webhook.site/[your-id] (or any webhook endpoint)
Event Type: token_triggered
Click: Create
```

**6. Simulate Attack**
```powershell
# Use the token you created
$token = "sk_xxxxx"
curl "http://localhost:3000/app?token=$token"
```

Expected result:
- Dashboard shows new attacker
- Event logged in database
- (Alert sends if webhook configured)

---

## Architecture Overview

### New Files Added
- `internal/handlers/page_handlers.go` — Page rendering (10 endpoints)
- `internal/database/dashboard.go` — Dashboard stats queries
- `static/js/app.js` — Frontend interactivity
- `templates/` — All HTML templates (your purple theme)

### Updated Files
- `cmd/server/main.go` — Added template routes
- `go.mod` — Added fiber/template dependency
- `internal/database/token.go` — Added pagination support

### Routes Available

**Unauthenticated:**
- `GET /login` — Login page
- `GET /signup` — Signup page

**Authenticated (require JWT):**
- `GET /app` — Dashboard
- `GET /tokens` — Tokens list
- `GET /tokens/new` — Create token form
- `GET /attackers` — Attackers list
- `GET /alerts` — Alerts configuration
- `GET /settings` — Settings page
- `GET /integrations` — Integrations page

---

## Frontend Design System

Your purple dark theme is preserved:
- **Background:** #0A0A14 (dark)
- **Surface:** #0D0B1A (darker)
- **Primary (Purple):** #7C3AED
- **Accent (Cyan):** #06B6D4
- **Border:** #1E1A30

All templates use **Tailwind CSS (CDN)** + **HTMX** for real-time updates.

---

## Troubleshooting

### "Port 3000 already in use"
```powershell
docker-compose down
docker system prune
docker-compose up --build
```

### "Cannot connect to database"
```powershell
# Clean data folder and restart
rm E:\KAVACH_VISION_1\data\kavach.db
docker-compose down
docker-compose up --build
```

### "Template not found"
```
Make sure templates/ folder exists with all 8 subfolders:
✓ auth/
✓ dashboard/
✓ tokens/
✓ attackers/
✓ alerts/
✓ settings/
✓ integrations/
✓ layouts/
```

---

## Next Phase (After Dashboard Works)

1. **Wire more pages** — Token detail, attacker profile, etc.
2. **Add real-time updates** — HTMX for live dashboard refresh
3. **Connect alert webhooks** — Test Slack, email, custom endpoints
4. **Create deployment guide** — For beta launch

---

**Status:** Dashboard ready to build! 🚀
