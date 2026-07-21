# KAVACH — Day 1 Setup ✅

## What's Been Created

```
E:\KAVACH_VISION_1\
├── cmd\server\main.go              ← Entry point
├── internal\
│   ├── database\db.go              ← SQLite initialization & migrations
│   ├── handlers\                   ← HTTP routes (next)
│   ├── services\                   ← Business logic (next)
│   ├── models\models.go            ← Data structures
│   ├── middleware\                 ← JWT auth (next)
│   ├── alerts\                     ← Notifications (next)
│   ├── fingerprint\                ← Fingerprinting (next)
│   └── classifier\                 ← Attack detection (next)
├── migrations\001_init.sql         ← Database schema
├── templates\                      ← HTML templates (to copy from archive)
├── static\css, static\js           ← Frontend assets
├── go.mod                          ← Dependencies
├── .env                            ← Configuration
├── .gitignore                      ← Git ignore rules
├── Dockerfile                      ← Container image
├── docker-compose.yml              ← Docker orchestration
└── README_SETUP.md                 ← This file
```

## Quick Start (Windows)

### Prerequisites
- Go 1.22+ installed
- Git installed (optional)

### Step 1: Download Dependencies
```bash
cd E:\KAVACH_VISION_1
go mod download
go mod tidy
```

### Step 2: Build & Run
```bash
go run ./cmd/server/main.go
```

**Expected output:**
```
✓ Database connection established
✓ Database migrations completed
🚀 Starting KAVACH server on port 3000
```

### Step 3: Test
Open browser: `http://localhost:3000`

```json
{
  "status": "KAVACH server running",
  "version": "1.0.0"
}
```

### Step 4: With Docker (Optional)
```bash
docker-compose up --build
```

---

## Next: Day 1 Tasks

### Task 2: Authentication System
- [ ] Create user registration handler
- [ ] Create login handler
- [ ] JWT token generation & validation
- [ ] Auth middleware

### Task 3: Token Management
- [ ] Create token generation service
- [ ] Create token API routes (CRUD)
- [ ] Store tokens in database

### Task 4: Dashboard Backend
- [ ] Create dashboard data aggregation service
- [ ] Dashboard stats endpoint
- [ ] Recent events endpoint

---

## Database
- **Type:** SQLite (file-based, zero setup)
- **File:** `./kavach.db` (auto-created on first run)
- **Schema:** Auto-applied from `migrations/001_init.sql`
- **Tables:** users, tokens, attackers, trigger_events, alert_configs, sent_alerts

---

## Environment Variables (.env)
```
PORT=3000                                   # Server port
DATABASE_PATH=./kavach.db                   # SQLite path
JWT_SECRET=[your-secret-key]                # JWT signing key (keep secret!)
ENVIRONMENT=development                    # dev/prod
LOG_LEVEL=debug                             # Verbosity
```

---

## Troubleshooting

**Q: "go mod not found"**
- `go mod init github.com/jindal-parth/kavach` (already done)
- `go mod tidy`

**Q: "sqlite3 not installed"**
- Windows: Should work automatically with CGO_ENABLED=1
- If issues: Install MinGW or use Docker

**Q: Database file not created**
- Check `DATABASE_PATH` in `.env`
- Ensure write permissions in directory

---

## Ready for Day 1? ✅
Commit this setup to git:
```bash
git init
git add .
git commit -m "Day 1: Project setup with Fiber + SQLite"
```

Next session: **Task 2 — Authentication System**
