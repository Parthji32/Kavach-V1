# KAVACH Build & Run Guide

## Prerequisites
- **Go 1.22+** [Download](https://golang.org/dl/)
- **Windows 10/11** (or Linux/macOS)
- **Git** (optional)

---

## Build & Run

### Option 1: Direct Run (Development)

```bash
cd E:\KAVACH_VISION_1

# Download dependencies
go mod download
go mod tidy

# Run server
go run ./cmd/server/main.go
```

**Expected Output:**
```
✓ Database connection established
✓ Database migrations completed

🚀 KAVACH Server Starting

📍 Server: http://localhost:3000
📊 Dashboard: http://localhost:3000/dashboard
📚 API Docs: https://docs.kavach.security
```

### Option 2: Build Binary

```bash
cd E:\KAVACH_VISION_1

# Build Windows binary
go build -o kavach.exe ./cmd/server/main.go

# Run the binary
.\kavach.exe
```

### Option 3: Docker

```bash
cd E:\KAVACH_VISION_1

# Build Docker image
docker build -t kavach:latest .

# Run container
docker-compose up --build
```

Server will be available at: `http://localhost:3000`

---

## Database

SQLite database is automatically created at `./kavach.db` on first run.

**To reset database:**
```bash
# Stop server (Ctrl+C)
rm kavach.db
go run ./cmd/server/main.go
```

---

## Configuration

Edit `.env` file:
```env
PORT=3000                                   # Server port
DATABASE_PATH=./kavach.db                   # Database file location
JWT_SECRET=your_secret_key_here             # Keep this secret!
ENVIRONMENT=development                    # dev or production
LOG_LEVEL=debug                             # debug, info, warn, error
```

---

## Testing the API

### 1. Register User
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123",
    "full_name": "Test User"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123"
  }'
```

**Save the `token` from response** and use in next calls.

### 3. Create Token
```bash
curl -X POST http://localhost:3000/api/tokens \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "token_type": "api_key",
    "description": "Test honeypot"
  }'
```

### 4. List Tokens
```bash
curl -X GET http://localhost:3000/api/tokens \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 5. Dashboard Stats
```bash
curl -X GET http://localhost:3000/api/dashboard/stats \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Using Postman/Insomnia

Import collection from: `API.md`

Or manually create requests with:
- **Base URL:** `http://localhost:3000/api`
- **Auth Header:** `Authorization: Bearer <token>`
- **Content-Type:** `application/json`

---

## Troubleshooting

### Error: "database file is locked"
- Windows: Close any running instances
- SQLite limit: Only one writer at a time
- Solution: Stop server, wait 2 seconds, restart

### Error: "go: not found"
- Install Go from https://golang.org/dl/
- Add Go to PATH (Windows: Settings → Environment Variables)
- Restart terminal

### Error: "port 3000 already in use"
- Change PORT in `.env` to 3001, 3002, etc.
- Or kill process: `lsof -i :3000` (macOS/Linux)

### Slow performance
- Check `LOG_LEVEL=debug` is not enabled in production
- Use `ENVIRONMENT=production` for better performance
- Check disk space (SQLite performance depends on I/O)

---

## Production Deployment

### Minimal Requirements
- Linux server (or Windows Server)
- Go 1.22 installed
- Or: Docker + Docker Compose

### Deploy to Linux VPS

```bash
# SSH into server
ssh user@your-vps.com

# Clone or upload code
git clone https://github.com/jindal-parth/kavach.git
cd kavach

# Build binary
go build -o /usr/local/bin/kavach ./cmd/server

# Create systemd service
sudo tee /etc/systemd/system/kavach.service > /dev/null << EOF
[Unit]
Description=KAVACH Security Platform
After=network.target

[Service]
Type=simple
User=kavach
WorkingDirectory=/opt/kavach
ExecStart=/usr/local/bin/kavach
Restart=on-failure
Environment="PORT=3000"
Environment="DATABASE_PATH=/opt/kavach/kavach.db"
Environment="JWT_SECRET=[REDACTED_PASSWORD]"
Environment="ENVIRONMENT=production"

[Install]
WantedBy=multi-user.target
EOF

# Start service
sudo systemctl enable kavach
sudo systemctl start kavach

# Check status
sudo systemctl status kavach
```

### Deploy with Docker Compose (Recommended)

```bash
# On your server
git clone https://github.com/jindal-parth/kavach.git
cd kavach
docker-compose up -d

# Check logs
docker-compose logs -f
```

---

## Performance Benchmarks

**Tested on:** Windows 11, Intel i7, 16GB RAM

| Operation | Time |
|-----------|------|
| Register | ~50ms |
| Login | ~100ms |
| Create Token | ~20ms |
| List Tokens (100) | ~30ms |
| Dashboard Stats | ~50ms |
| Classify Request | ~5ms |

---

## Next Steps

1. ✅ Run server locally
2. ✅ Test API endpoints
3. ⬜ Build frontend dashboard (templates)
4. ⬜ Implement reverse proxy
5. ⬜ Add rate limiting
6. ⬜ Production deployment

---

**Version:** 1.0.0  
**Last Updated:** 2026-07-18
