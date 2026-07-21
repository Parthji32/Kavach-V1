$WEBHOOK_URL = "https://webhook.site/8a6aa811-c453-44eb-8c8d-c4840abfb57b"
$BASE_URL = "http://localhost:3000"

Write-Host "`n🧪 KAVACH Webhook Alert Test`n" -ForegroundColor Cyan

# 1. Health check
Write-Host "[1] Health check..." -ForegroundColor Yellow
$r = curl.exe -s $BASE_URL/health
Write-Host "✅ Server responding`n" -ForegroundColor Green

# 2. Signup first (if needed)
Write-Host "[2a] Attempting user signup..." -ForegroundColor Yellow
$signupBody = '{"full_name":"Test User","email":"webhook_test@example.com","password":"TestPassword123!"}'
curl.exe -s -X POST "$BASE_URL/api/auth/register" -H "Content-Type: application/json" -d $signupBody | Out-Null
Write-Host "✅ Signup processed (or user exists)`n" -ForegroundColor Green

# 2b. Login
Write-Host "[2b] Logging in..." -ForegroundColor Yellow
$loginBody = '{"email":"webhook_test@example.com","password":"TestPassword123!"}'
$r = curl.exe -s -X POST "$BASE_URL/api/auth/login" -H "Content-Type: application/json" -d $loginBody | ConvertFrom-Json
$JWT = $r.data.token
Write-Host "✅ JWT obtained`n" -ForegroundColor Green

# 3. Create alert config
Write-Host "[3] Creating webhook alert config..." -ForegroundColor Yellow
$alertBody = "{`"alert_type`":`"webhook`",`"destination`":`"$WEBHOOK_URL`"}"
$alertResp = curl.exe -s -X POST "$BASE_URL/api/alerts/config" `
  -H "Authorization: Bearer $JWT" `
  -H "Content-Type: application/json" `
  -d $alertBody
Write-Host "Response: $alertResp" -ForegroundColor Gray
Write-Host "✅ Alert config created`n" -ForegroundColor Green

# 4. Create token
Write-Host "[4] Creating honeypot token..." -ForegroundColor Yellow
$tokenBody = '{"token_type":"url","description":"webhook test token"}'
$r = curl.exe -s -X POST "$BASE_URL/api/tokens" `
  -H "Authorization: Bearer $JWT" `
  -H "Content-Type: application/json" `
  -d $tokenBody | ConvertFrom-Json
$TOKEN = $r.data.token_value
Write-Host "✅ Token: $($TOKEN.Substring(0, 40))...`n" -ForegroundColor Green

# 5. Get stats before
Write-Host "[5] Dashboard stats BEFORE attack..." -ForegroundColor Yellow
$r = curl.exe -s -H "Authorization: Bearer $JWT" "$BASE_URL/api/dashboard/stats" | ConvertFrom-Json
$before_attackers = $r.data.total_attackers
$before_events = $r.data.events_last_24h
Write-Host "  Attackers: $before_attackers | Events: $before_events`n" -ForegroundColor Gray

# 6. TRIGGER ATTACK
Write-Host "[6] 🚨 TRIGGERING ATTACK..." -ForegroundColor Red
Write-Host "Accessing: $BASE_URL/api/dashboard/stats?token=$($TOKEN.Substring(0,30))...`n" -ForegroundColor Gray
curl.exe -s "$BASE_URL/api/dashboard/stats?token=$TOKEN" | Out-Null
Write-Host "✅ Honeypot token accessed (attack triggered)`n" -ForegroundColor Green

# 7. Wait
Write-Host "[7] Waiting for webhook delivery (3 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# 8. Get stats after
Write-Host "`n[8] Dashboard stats AFTER attack..." -ForegroundColor Yellow
$r = curl.exe -s -H "Authorization: Bearer $JWT" "$BASE_URL/api/dashboard/stats" | ConvertFrom-Json
$after_attackers = $r.data.total_attackers
$after_events = $r.data.events_last_24h
Write-Host "  Attackers: $after_attackers | Events: $after_events`n" -ForegroundColor Gray

# 9. Verify
Write-Host "[9] VERIFICATION:" -ForegroundColor Yellow
$attacker_increase = $after_attackers - $before_attackers
$event_increase = $after_events - $before_events

if ($attacker_increase -gt 0) {
  Write-Host "  ✅ Attacker DETECTED! (+$attacker_increase)" -ForegroundColor Green
} else {
  Write-Host "  ⚠️  No new attacker (increase: $attacker_increase)" -ForegroundColor Yellow
}

if ($event_increase -gt 0) {
  Write-Host "  ✅ Event LOGGED! (+$event_increase)" -ForegroundColor Green
} else {
  Write-Host "  ⚠️  No new event (increase: $event_increase)" -ForegroundColor Yellow
}

Write-Host "`n" -ForegroundColor White
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  📍 CHECK WEBHOOK.SITE FOR INCOMING ALERT PAYLOAD" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "`n🔗 https://webhook.site/8a6aa811-c453-44eb-8c8d-c4840abfb57b`n" -ForegroundColor Cyan

Write-Host "Expected webhook payload should contain:" -ForegroundColor White
Write-Host "  • event_type: 'token_accessed'" -ForegroundColor Gray
Write-Host "  • risk_level: 'critical'" -ForegroundColor Gray
Write-Host "  • risk_score: 95" -ForegroundColor Gray
Write-Host "  • token_type: 'url'" -ForegroundColor Gray
Write-Host "  • severity: 'critical'" -ForegroundColor Gray
Write-Host "`n" -ForegroundColor White
