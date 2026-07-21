param()

# KAVACH Webhook Alert Test - Version 3 (ASCII safe)
$WEBHOOK_URL = "https://webhook.site/8a6aa811-c453-44eb-8c8d-c4840abfb57b"
$BASE_URL = "http://localhost:3000"
$EMAIL = "webhook_test@example.com"
$PASSWORD = "TestPassword123!"

Write-Host "KAVACH Webhook Alert Test`n" -ForegroundColor Cyan

# 1. Health check
Write-Host "[1] Health check..." -ForegroundColor Yellow
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/health" -UseBasicParsing
    Write-Host "OK - Server responding`n" -ForegroundColor Green
} catch {
    Write-Host "FAIL - Server not responding" -ForegroundColor Red
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
    Write-Host "OK - User created or already exists`n" -ForegroundColor Green
} catch {
    Write-Host "OK - User likely exists`n" -ForegroundColor Yellow
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
    Write-Host "OK - Logged in as $($data.data.user.email)`n" -ForegroundColor Green
} catch {
    Write-Host "FAIL - Login error: $_`n" -ForegroundColor Red
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
    Write-Host "OK - Alert config created`n" -ForegroundColor Green
} catch {
    Write-Host "FAIL - Alert creation: $($_.Exception.Message)`n" -ForegroundColor Red
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
    Write-Host "OK - Token: $($TOKEN.Substring(0, 40))...`n" -ForegroundColor Green
} catch {
    Write-Host "FAIL - Token creation: $_`n" -ForegroundColor Red
    exit 1
}

# 5. Get stats BEFORE
Write-Host "[5] Dashboard stats BEFORE attack..." -ForegroundColor Yellow
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/dashboard/stats" -UseBasicParsing -Headers $headers
    $data = $r.Content | ConvertFrom-Json
    $before_attackers = $data.data.total_attackers
    $before_events = $data.data.events_last_24h
    Write-Host "  Attackers: $before_attackers | Events: $before_events`n" -ForegroundColor Gray
} catch {
    Write-Host "FAIL - Stats retrieval: $_`n" -ForegroundColor Red
}

# 6. TRIGGER ATTACK
Write-Host "[6] TRIGGERING ATTACK - accessing honeypot token..." -ForegroundColor Red
try {
    $attackUrl = "$BASE_URL/api/dashboard/stats?token=$TOKEN"
    Write-Host "Accessing: $($attackUrl.Substring(0, 80))..." -ForegroundColor Gray
    $r = Invoke-WebRequest -Uri $attackUrl -UseBasicParsing -ErrorAction Continue
    Write-Host "OK - Attack triggered`n" -ForegroundColor Green
} catch {
    Write-Host "OK - Attack processed by middleware`n" -ForegroundColor Green
}

# 7. Wait for webhook
Write-Host "[7] Waiting 4 seconds for webhook delivery..." -ForegroundColor Yellow
Start-Sleep -Seconds 4
Write-Host "`n"

# 8. Get stats AFTER
Write-Host "[8] Dashboard stats AFTER attack..." -ForegroundColor Yellow
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/api/dashboard/stats" -UseBasicParsing -Headers $headers
    $data = $r.Content | ConvertFrom-Json
    $after_attackers = $data.data.total_attackers
    $after_events = $data.data.events_last_24h
    Write-Host "  Attackers: $after_attackers | Events: $after_events`n" -ForegroundColor Gray
} catch {
    Write-Host "FAIL - Stats retrieval: $_`n" -ForegroundColor Red
}

# 9. Verify
Write-Host "[9] VERIFICATION:" -ForegroundColor Yellow
$attacker_diff = $after_attackers - $before_attackers
$event_diff = $after_events - $before_events

if ($attacker_diff -gt 0) {
    Write-Host "  OK - ATTACKER DETECTED (plus $attacker_diff)" -ForegroundColor Green
} else {
    Write-Host "  WARN - No new attacker (diff: $attacker_diff)" -ForegroundColor Yellow
}

if ($event_diff -gt 0) {
    Write-Host "  OK - EVENT LOGGED (plus $event_diff)" -ForegroundColor Green
} else {
    Write-Host "  WARN - No new event (diff: $event_diff)" -ForegroundColor Yellow
}

Write-Host "`n"
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "  NEXT: Check webhook.site for alert payload" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "`nURL: https://webhook.site/8a6aa811-c453-44eb-8c8d-c4840abfb57b`n" -ForegroundColor Cyan

Write-Host "Expected webhook payload should have:" -ForegroundColor White
Write-Host "  - event_type: token_accessed" -ForegroundColor Gray
Write-Host "  - risk_level: critical" -ForegroundColor Gray
Write-Host "  - risk_score: 95" -ForegroundColor Gray
Write-Host "  - token_type: url" -ForegroundColor Gray
Write-Host "  - severity: critical" -ForegroundColor Gray
Write-Host "`n"
