CREATE TABLE pack_sizes (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL,
    sizes INTEGER[] NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_pack_sizes_version ON pack_sizes(version DESC);
CREATE INDEX idx_pack_sizes_active ON pack_sizes(is_active) WHERE is_active = true;
