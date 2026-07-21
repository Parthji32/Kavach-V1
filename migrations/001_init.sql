-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tokens table (honeypot tokens)
CREATE TABLE IF NOT EXISTS tokens (
    id TEXT PRIMARY KEY,
    name TEXT,
    user_id TEXT NOT NULL,
    token_type TEXT NOT NULL, -- 'url', 'api_key', 'document', 'dns', 'email'
    token_value TEXT NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    triggered_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Attackers table
CREATE TABLE IF NOT EXISTS attackers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    os TEXT,
    browser TEXT,
    device_type TEXT, -- 'desktop', 'mobile', 'bot'
    fingerprint TEXT UNIQUE,
    risk_level TEXT DEFAULT 'low', -- 'low', 'medium', 'high', 'critical'
    risk_score INTEGER DEFAULT 0, -- 0-100
    detection_count INTEGER DEFAULT 1,
    is_known_user BOOLEAN DEFAULT 0,
    is_blocked BOOLEAN DEFAULT 0,
    first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger events (when a token is triggered)
CREATE TABLE IF NOT EXISTS trigger_events (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token_id TEXT NOT NULL,
    attacker_id TEXT,
    event_type TEXT NOT NULL, -- 'token_accessed', 'token_used'
    http_method TEXT,
    request_path TEXT,
    request_headers TEXT,
    ip_address TEXT,
    user_agent TEXT,
    endpoint TEXT,
    request_payload TEXT,
    response_status INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (token_id) REFERENCES tokens(id) ON DELETE CASCADE,
    FOREIGN KEY (attacker_id) REFERENCES attackers(id) ON DELETE SET NULL
);

-- Alert configurations
CREATE TABLE IF NOT EXISTS alert_configs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    alert_type TEXT NOT NULL, -- 'webhook', 'email', 'slack'
    destination TEXT NOT NULL, -- URL, email, or slack channel
    is_enabled BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Sent alerts (log of alerts sent)
CREATE TABLE IF NOT EXISTS sent_alerts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    trigger_event_id TEXT NOT NULL,
    alert_config_id TEXT NOT NULL,
    status TEXT DEFAULT 'pending', -- 'pending', 'sent', 'failed'
    error_message TEXT,
    sent_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (trigger_event_id) REFERENCES trigger_events(id) ON DELETE CASCADE,
    FOREIGN KEY (alert_config_id) REFERENCES alert_configs(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_tokens_active ON tokens(is_active);
CREATE INDEX IF NOT EXISTS idx_attackers_user_id ON attackers(user_id);
CREATE INDEX IF NOT EXISTS idx_attackers_fingerprint ON attackers(fingerprint);
CREATE INDEX IF NOT EXISTS idx_trigger_events_user_id ON trigger_events(user_id);
CREATE INDEX IF NOT EXISTS idx_trigger_events_token_id ON trigger_events(token_id);
CREATE INDEX IF NOT EXISTS idx_trigger_events_timestamp ON trigger_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_alert_configs_user_id ON alert_configs(user_id);
CREATE INDEX IF NOT EXISTS idx_sent_alerts_user_id ON sent_alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_sent_alerts_trigger_event_id ON sent_alerts(trigger_event_id);
