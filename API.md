# KAVACH API Reference v1.0

## Base URL
```
http://localhost:3000/api
```

---

## Authentication Routes (No JWT Required)

### 1. User Registration
**POST** `/auth/register`

Register a new KAVACH account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123",
  "full_name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid-here",
    "email": "user@example.com",
    "full_name": "John Doe",
    "created_at": "2026-07-18T10:30:00Z",
    "updated_at": "2026-07-18T10:30:00Z"
  }
}
```

---

### 2. User Login
**POST** `/auth/login`

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-here",
    "email": "user@example.com",
    "full_name": "John Doe",
    "created_at": "2026-07-18T10:30:00Z",
    "updated_at": "2026-07-18T10:30:00Z"
  },
  "expires_at": 1721314200
}
```

**Note:** Use the `token` in the `Authorization` header for all subsequent requests:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## Protected Routes (JWT Required)

### 3. Get User Profile
**GET** `/user/profile`

Retrieve current user's profile information.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK):**
```json
{
  "id": "uuid-here",
  "email": "user@example.com",
  "full_name": "John Doe",
  "created_at": "2026-07-18T10:30:00Z",
  "updated_at": "2026-07-18T10:30:00Z"
}
```

---

## Token Routes

### 4. Create Single Token
**POST** `/tokens`

Create a new honeypot token.

**Request Body:**
```json
{
  "token_type": "api_key",
  "description": "Test API key"
}
```

**Supported Token Types:**
- `url` - URL-based token
- `api_key` - API key format (sk_...)
- `document` - Document filename
- `dns` - DNS subdomain
- `email` - Email address

**Response (201 Created):**
```json
{
  "message": "Token created successfully",
  "token": {
    "id": "token-uuid",
    "user_id": "user-uuid",
    "token_type": "api_key",
    "token_value": "sk_ab1234cd5678ef90...",
    "description": "Test API key",
    "is_active": true,
    "created_at": "2026-07-18T10:30:00Z",
    "triggered_at": null
  }
}
```

---

### 5. Create Multiple Tokens (Bulk)
**POST** `/tokens/bulk`

Create multiple tokens at once.

**Request Body:**
```json
{
  "count": 10
}
```

**Response (201 Created):**
```json
{
  "message": "Tokens created successfully",
  "tokens": [
    {
      "id": "token-uuid-1",
      "token_type": "url",
      "token_value": "https://api.internal.com/webhook?token=...",
      "is_active": true,
      "created_at": "2026-07-18T10:30:00Z"
    },
    // ... more tokens
  ],
  "count": 10
}
```

---

### 6. List All Tokens
**GET** `/tokens`

Retrieve all tokens for current user.

**Response (200 OK):**
```json
{
  "tokens": [
    {
      "id": "token-uuid",
      "token_type": "api_key",
      "token_value": "sk_...",
      "description": "Test API key",
      "is_active": true,
      "created_at": "2026-07-18T10:30:00Z",
      "triggered_at": null
    }
  ],
  "count": 5
}
```

---

### 7. Deactivate Token
**DELETE** `/tokens/{tokenID}`

Deactivate/delete a token.

**Response (200 OK):**
```json
{
  "message": "Token deleted successfully"
}
```

---

## Dashboard Routes

### 8. Get Dashboard Statistics
**GET** `/dashboard/stats`

Get real-time dashboard statistics.

**Response (200 OK):**
```json
{
  "total_tokens": 15,
  "active_tokens": 12,
  "total_attackers": 8,
  "high_risk_count": 3,
  "events_last_24h": 42,
  "recent_attackers": [
    {
      "id": "attacker-uuid",
      "ip_address": "192.168.1.100",
      "risk_score": 85,
      "last_seen": "2026-07-18T15:30:00Z",
      "is_blocked": false
    }
  ],
  "recent_events": [
    {
      "id": "event-uuid",
      "token_type": "api_key",
      "event_type": "token_accessed",
      "risk_score": 90,
      "timestamp": "2026-07-18T15:25:00Z"
    }
  ]
}
```

---

### 9. List Attackers
**GET** `/dashboard/attackers?limit=50`

Get list of detected attackers.

**Query Parameters:**
- `limit` (optional, default: 50, max: 1000) - Number of records to return

**Response (200 OK):**
```json
{
  "attackers": [
    {
      "id": "attacker-uuid",
      "user_id": "user-uuid",
      "ip_address": "203.0.113.50",
      "user_agent": "Mozilla/5.0...",
      "os": "Windows",
      "browser": "Chrome",
      "device_type": "desktop",
      "fingerprint": "a1b2c3d4e5f6...",
      "risk_score": 75,
      "is_blocked": false,
      "first_seen": "2026-07-17T10:00:00Z",
      "last_seen": "2026-07-18T15:30:00Z"
    }
  ],
  "count": 8
}
```

---

### 10. List Trigger Events
**GET** `/dashboard/events?limit=100`

Get trigger events (when tokens are accessed).

**Query Parameters:**
- `limit` (optional, default: 100, max: 1000) - Number of records

**Response (200 OK):**
```json
{
  "events": [
    {
      "id": "event-uuid",
      "user_id": "user-uuid",
      "token_id": "token-uuid",
      "attacker_id": "attacker-uuid",
      "event_type": "token_accessed",
      "http_method": "GET",
      "endpoint": "/admin?token=...",
      "request_payload": "{...}",
      "response_status": 200,
      "timestamp": "2026-07-18T15:25:00Z"
    }
  ],
  "count": 42
}
```

---

## Alert Configuration Routes

### 11. Create Alert Configuration
**POST** `/alerts/config`

Configure an alert destination (webhook, email, Slack, etc.).

**Request Body:**
```json
{
  "alert_type": "webhook",
  "destination": "https://your-webhook-url.com/alerts"
}
```

**Supported Alert Types:**
- `webhook` - HTTP webhook
- `email` - Email address
- `slack` - Slack webhook URL

**Response (201 Created):**
```json
{
  "message": "Alert config created successfully",
  "config": {
    "id": "config-uuid",
    "user_id": "user-uuid",
    "alert_type": "webhook",
    "destination": "https://your-webhook-url.com/alerts",
    "is_enabled": true,
    "created_at": "2026-07-18T10:30:00Z"
  }
}
```

---

### 12. List Alert Configurations
**GET** `/alerts/config`

Retrieve all alert configurations.

**Response (200 OK):**
```json
{
  "configs": [
    {
      "id": "config-uuid",
      "alert_type": "webhook",
      "destination": "https://...",
      "is_enabled": true,
      "created_at": "2026-07-18T10:30:00Z"
    },
    {
      "id": "config-uuid-2",
      "alert_type": "slack",
      "destination": "https://hooks.slack.com/...",
      "is_enabled": true,
      "created_at": "2026-07-18T11:00:00Z"
    }
  ],
  "count": 2
}
```

---

### 13. Delete Alert Configuration
**DELETE** `/alerts/config/{configID}`

Remove an alert configuration.

**Response (200 OK):**
```json
{
  "message": "Alert config deleted successfully"
}
```

---

## Error Responses

All endpoints return consistent error format:

**Response (400/401/500):**
```json
{
  "error": "Error Type",
  "message": "Detailed error message"
}
```

**Common Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `409` - Conflict (e.g., email already exists)
- `500` - Internal Server Error

---

## Example Workflow

### 1. Register
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123",
    "full_name": "John Doe"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

### 3. Create Token
```bash
curl -X POST http://localhost:3000/api/tokens \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "token_type": "api_key",
    "description": "Production API key"
  }'
```

### 4. Setup Alert Webhook
```bash
curl -X POST http://localhost:3000/api/alerts/config \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "alert_type": "webhook",
    "destination": "https://your-domain.com/webhooks/kavach"
  }'
```

### 5. Check Dashboard
```bash
curl -X GET http://localhost:3000/api/dashboard/stats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Rate Limiting
Currently not enforced. Will be added in production.

## Webhooks
When a honeypot token is triggered, webhook alerts receive:
```json
{
  "event_id": "event-uuid",
  "token_id": "token-uuid",
  "attacker_id": "attacker-uuid",
  "event_type": "token_accessed",
  "http_method": "GET",
  "endpoint": "/admin?token=...",
  "response_status": 200,
  "timestamp": "2026-07-18T15:25:00Z",
  "request_payload": "..."
}
```

---

**Last Updated:** 2026-07-18  
**Version:** 1.0.0
