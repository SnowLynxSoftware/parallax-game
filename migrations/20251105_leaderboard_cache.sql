-- ############################
-- Parallax Leaderboard Cache Schema
--
-- https://snowlynxsoftware.net 
--
-- Copyright 2025. Snow Lynx Software, LLC. All Rights Reserved.
-- ############################

-- LEADERBOARD CACHE: Three leaderboard categories with async refresh

-- ############################
-- STEP 1: CREATE LEADERBOARD CACHE METADATA TABLE
-- ############################

-- Stores metadata about each leaderboard type
CREATE TABLE leaderboard_cache (
    id SERIAL PRIMARY KEY,
    leaderboard_type VARCHAR(50) NOT NULL UNIQUE, -- 'legendary', 'power', 'expeditions'
    last_synced TIMESTAMP NOT NULL DEFAULT NOW(),
    is_syncing BOOLEAN NOT NULL DEFAULT FALSE, -- Prevents concurrent rebuilds
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leaderboard_cache_type ON leaderboard_cache(leaderboard_type);

-- ############################
-- STEP 2: CREATE LEADERBOARD CACHE ITEMS TABLE
-- ############################

-- Stores actual rankings for each leaderboard
CREATE TABLE leaderboard_cache_items (
    id SERIAL PRIMARY KEY,
    leaderboard_type VARCHAR(50) NOT NULL, -- 'legendary', 'power', 'expeditions'
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL, -- Denormalized for performance
    score BIGINT NOT NULL, -- Actual value (item count, power score, expedition count)
    rank INTEGER NOT NULL, -- Calculated rank (1, 2, 3, etc.)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_leaderboard_user UNIQUE (leaderboard_type, user_id)
);

CREATE INDEX idx_leaderboard_cache_items_type_rank ON leaderboard_cache_items(leaderboard_type, rank);
CREATE INDEX idx_leaderboard_cache_items_user ON leaderboard_cache_items(user_id);

-- ############################
-- STEP 3: SEED INITIAL LEADERBOARD METADATA
-- ############################

-- Seed the three leaderboard types with old last_synced to force initial sync
INSERT INTO leaderboard_cache (leaderboard_type, last_synced) VALUES
    ('legendary', '2000-01-01 00:00:00'),  -- Force initial sync
    ('power', '2000-01-01 00:00:00'),
    ('expeditions', '2000-01-01 00:00:00')
ON CONFLICT (leaderboard_type) DO NOTHING;
