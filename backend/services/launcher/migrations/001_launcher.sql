CREATE TABLE IF NOT EXISTS builds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    mod_loader TEXT NOT NULL DEFAULT 'fabric',
    mc_version TEXT NOT NULL,
    server_id TEXT NOT NULL,
    file_hash TEXT NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    changelog TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_builds_server_id ON builds(server_id);
