CREATE TABLE IF NOT EXISTS authenticated_sessions (
    session_id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_jti VARCHAR(36) NOT NULL,
    start_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    exp TIMESTAMPTZ NOT NULL,
    last_ip VARCHAR(45) NOT NULL,
    device_type TEXT DEFAULT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS anonymous_sessions (
    session_id UUID PRIMARY KEY,
    client_id VARCHAR(36) NOT NULL,
    token_jti VARCHAR(36) NOT NULL,
    start_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    exp TIMESTAMPTZ NOT NULL,
    last_ip VARCHAR(45) NOT NULL,
    device_type TEXT DEFAULT NULL,
    is_active BOOLEAN DEFAULT TRUE
);