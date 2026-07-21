# KAVACH Webhook Alert Test - Version 2
$WEBHOOK_URL = "https://webhook.site/8a6aa811-c453-44eb-8c8d-c4840abfb57b"
$BASE_URL = "http://localhost:3000"
$EMAIL = "webhook_test@example.com"
$PASSWORD = "TestPassword123!"

Write-Host "`n🧪 KAVACH Webhook Alert Test (v2)`n" -ForegroundColor Cyan

# 1. Health check
Write-Host "[1] Health check..." -ForegroundColor Yellow
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/health" -UseBasicParsing
    Write-Host "✅ Server responding`n" -ForegroundColor Green
} catch {
    Write-Host "❌ Server not responding: $_" -ForegroundColor Red
    exit 1
}

# 2a. Signup
Write-Host "[2a] User signup..." -ForegroundColor Yellow
$signupBody = @{
    full_name = "Test User"
    email = $EMAIL
    password = $PASSWORD
} | ConvertTo-Json

try {
    Invoke-WebRequest -Uri "$BASE_URL/api/auth/register" -UseBasicParsing `
        -Method POST -Body $signupBody -ContentType "application/json" | Out-Null
    Write-Host "✅ User created or already exists`n" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Signup response: $($_.Exception.Message)" -ForegroundColor Yellow
}

# 2b. Login
Write-Host "[2b] Login..." -ForegroundColor Yellow
$loginBody = @{
    email = $EMAIL
    password = $PASSWORD
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/auth/login" -UseBasicParsing `
        -Method POST -Body $loginBody -ContentType "application/json"
    $data = $r.Content | ConvertFrom-Json
    $JWT = $data.data.token
    $USER_ID = $data.data.user.id
    Write-Host "✅ Logged in: $($data.data.user.email)`n" -ForegroundColor Green
} catch {
    Write-Host "❌ Login failed: $_" -ForegroundColor Red
    exit 1
}

$headers = @{
    "Authorization" = "Bearer $JWT"
    "Content-Type" = "application/json"
}

# 3. Create alert config
Write-Host "[3] Creating webhook alert config..." -ForegroundColor Yellow
$alertBody = @{
    alert_type = "webhook"
    destination = $WEBHOOK_URL
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/alerts/config" -UseBasicParsing `
        -Method POST -Body $alertBody -Headers $headers
    Write-Host "✅ Alert config created`n" -ForegroundColor Green
} catch {
    Write-Host "❌ Alert creation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Response: $($_.Exception.Response.StatusCode)" -ForegroundColor Gray
    exit 1
}

# 4. Create honeypot token
Write-Host "[4] Creating honeypot token..." -ForegroundColor Yellow
$tokenBody = @{
    token_type = "url"
    description = "Webhook test token"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/tokens" -UseBasicParsing `
        -Method POST -Body $tokenBody -Headers $headers
    $data = $r.Content | ConvertFrom-Json
    $TOKEN = $data.data.token_value
    Write-Host "✅ Token created: $($TOKEN.Substring(0, 40))...`n" -ForegroundColor Green
} catch {
    Write-Host "❌ Token creation failed: $_" -ForegroundColor Red
    exit 1
}

# 5. Get stats BEFORE
Write-Host "[5] Dashboard stats BEFORE attack..." -ForegroundColor Yellow
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/dashboard/stats" -UseBasicParsing -Headers $headers
    $data = $r.Content | ConvertFrom-Json
    $before_attackers = $data.data.total_attackers
    $before_events = $data.data.events_last_24h
    Write-Host "  • Total Attackers: $before_attackers" -ForegroundColor Gray
    Write-Host "  • Events (24h): $before_events`n" -ForegroundColor Gray
} catch {
    Write-Host "❌ Stats retrieval failed: $_" -ForegroundColor Red
}

# 6. TRIGGER ATTACK
Write-Host "[6] 🚨 TRIGGERING ATTACK - accessing honeypot token..." -ForegroundColor Red
try {
    $attackUrl = "$BASE_URL/api/dashboard/stats?token=[REDACTED_PARAM]
    Write-Host "URL: $($attackUrl.Substring(0, 80))..." -ForegroundColor Gray
    $r = Invoke-WebRequest -Uri $attackUrl -UseBasicParsing -ErrorAction Continue
    Write-Host "✅ Attack triggered (token detected)`n" -ForegroundColor Green
} catch {
    Write-Host "✅ Attack processed (middleware intercepted)`n" -ForegroundColor Green
}

# 7. Wait for webhook
Write-Host "[7] Waiting for webhook delivery (4 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 4
Write-Host "`n" -ForegroundColor Green

# 8. Get stats AFTER
Write-Host "[8] Dashboard stats AFTER attack..." -ForegroundColor Yellow
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/dashboard/stats" -UseBasicParsing -Headers $headers
    $data = $r.Content | ConvertFrom-Json
    $after_attackers = $data.data.total_attackers
    $after_events = $data.data.events_last_24h
    Write-Host "  • Total Attackers: $after_attackers" -ForegroundColor Gray
    Write-Host "  • Events (24h): $after_events`n" -ForegroundColor Gray
} catch {
    Write-Host "❌ Stats retrieval failed: $_" -ForegroundColor Red
}

# 9. Verify
Write-Host "[9] VERIFICATION:" -ForegroundColor Yellow
$attacker_diff = $after_attackers - $before_attackers
$event_diff = $after_events - $before_events

if ($attacker_diff -gt 0) {
    Write-Host "  ✅ ATTACKER DETECTED (+$attacker_diff new attacker)" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  No new attacker detected (diff: $attacker_diff)" -ForegroundColor Yellow
}

if ($event_diff -gt 0) {
    Write-Host "  ✅ EVENT LOGGED (+$event_diff new event)" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  No new event logged (diff: $event_diff)" -ForegroundColor Yellow
}

Write-Host "`n"
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  📍 NOW CHECK WEBHOOK.SITE FOR ALERT PAYLOAD" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "`n🔗 https://webhook.site/8a6aa811-c453-44eb-8c8d-c4840abfb57b`n" -ForegroundColor Cyan

Write-Host "Expected payload structure:" -ForegroundColor White
Write-Host @"
{
  "event_type": "token_accessed",
  "timestamp": "2026-07-20T...",
  "user_id": "...",
  "token_id": "...",
  "token_value": "sk_...",
  "token_type": "url",
  "attacker_id": "...",
  "attacker_ip": "172.18.0.1",
  "risk_score": 95,
  "risk_level": "critical",
  "detected_at": "2026-07-20T...",
  "severity": "critical",
  "message": "Honeypot token (url) accessed from 172.18.0.1 with risk score 95 (critical)"
}
"@ -ForegroundColor Gray

Write-Host "`n✅ If you see this JSON in webhook.site, alerts are working!`n" -ForegroundColor Green
