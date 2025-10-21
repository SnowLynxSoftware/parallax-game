# Database Schema Plan - Parallax Game

**Version:** 1.0  
**Date:** October 21, 2025  
**Status:** Ready for Implementation

---

## Overview

This document defines the complete PostgreSQL database schema for the Parallax game core loop, including:

- Teams with equipment slots
- Expeditions and loot generation
- Inventory management
- Loot drop tables
- Item stats and elemental affinities

---

## Schema Design Principles

1. **Simplicity First**: Permanent team stats stored directly on teams table (no complex upgrade tracking)
2. **Flexible Equipment**: Items can be equipped to multiple teams simultaneously
3. **Type Safety**: Enums for rarities, world types, item types, slot types, elements
4. **Audit Trail**: Track expedition loot for debugging and analytics
5. **Scalability**: Designed for 50-100 concurrent users, optimized for read-heavy queries

---

## Table Definitions

### 1. **`teams`**

Stores user's expedition teams (up to 5 per user).

```sql
CREATE TYPE team_specialization AS ENUM ('none', 'tank', 'scout', 'scientist');

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_number INT NOT NULL CHECK (team_number >= 1 AND team_number <= 5),

    -- Base stats (permanent upgrades from consumables)
    speed_bonus DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    luck_bonus DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    power_bonus INT NOT NULL DEFAULT 0,

    -- Team configuration
    specialization team_specialization NOT NULL DEFAULT 'none',
    is_unlocked BOOLEAN NOT NULL DEFAULT false,

    -- Equipment slots (FK to user_inventory.id)
    equipped_weapon_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_armor_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_accessory_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_artifact_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_relic_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, team_number)
);

CREATE INDEX idx_teams_user_id ON teams(user_id);
CREATE INDEX idx_teams_unlocked ON teams(user_id, is_unlocked);
```

**Notes:**

- Each user has 5 team slots (unlocked progressively via game logic)
- Base stats are permanent upgrades from consuming items
- Equipment slots reference `user_inventory.id` (not `loot_items.id`)
- Items in inventory can be equipped to multiple teams

---

### 2. **`rifts`** (Static/Seed Data)

Defines the 6 parallel world rifts players can explore.

```sql
CREATE TYPE world_type AS ENUM ('tutorial', 'fire', 'ice', 'tech', 'nature', 'void');
CREATE TYPE difficulty_level AS ENUM ('tutorial', 'easy', 'medium', 'hard', 'legendary');
CREATE TYPE elemental_type AS ENUM ('fire', 'water', 'wind', 'earth', 'void', 'none');

CREATE TABLE rifts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    world_type world_type NOT NULL,
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
    difficulty difficulty_level NOT NULL,

    -- Elemental weakness (for artifact bonuses)
    weak_to_element elemental_type NOT NULL DEFAULT 'none',

    -- Unlock requirements (text description, enforced in Go)
    unlock_requirement_text VARCHAR(255),

    -- UI
    icon VARCHAR(50) NOT NULL, -- Font Awesome class

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rifts_difficulty ON rifts(difficulty);
CREATE INDEX idx_rifts_world_type ON rifts(world_type);
```

**Seed Data (6 Rifts):**

1. Tutorial Rift - 5min, tutorial difficulty, no weakness
2. Crimson Wastes (Fire) - 15min, easy, weak to water
3. Frozen Expanse (Ice) - 30min, medium, weak to fire
4. Neon Sprawl (Tech) - 30min, medium, weak to void
5. Verdant Overgrowth (Nature) - 60min, hard, weak to fire
6. Void Confluence (Void) - 120min, legendary, weak to earth

---

### 3. **`loot_items`** (Static/Seed Data)

Defines all collectible items (~75+ items across 5 worlds).

```sql
CREATE TYPE item_rarity AS ENUM ('common', 'uncommon', 'rare', 'epic', 'legendary');
CREATE TYPE item_type AS ENUM ('equipment', 'consumable');
CREATE TYPE equipment_slot_type AS ENUM ('weapon', 'armor', 'accessory', 'artifact', 'relic');

CREATE TABLE loot_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,

    -- Classification
    rarity item_rarity NOT NULL,
    world_type world_type NOT NULL,
    item_type item_type NOT NULL,

    -- Equipment-specific (NULL for consumables)
    equipment_slot equipment_slot_type,

    -- Stats (all items have all stats, may be 0)
    speed_bonus DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    luck_bonus DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    power_bonus INT NOT NULL DEFAULT 0,

    -- Elemental affinity (for artifacts/relics)
    elemental_affinity elemental_type NOT NULL DEFAULT 'none',

    -- Power value for leaderboard calculations
    power_value INT NOT NULL DEFAULT 1,

    -- UI
    icon VARCHAR(50) NOT NULL, -- Font Awesome class

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_loot_items_rarity ON loot_items(rarity);
CREATE INDEX idx_loot_items_world_type ON loot_items(world_type);
CREATE INDEX idx_loot_items_item_type ON loot_items(item_type);
CREATE INDEX idx_loot_items_slot ON loot_items(equipment_slot) WHERE equipment_slot IS NOT NULL;
```

**Notes:**

- **Equipment**: Has `equipment_slot` defined (weapon/armor/accessory/artifact/relic)
- **Consumables**: `equipment_slot` is NULL, consumed to permanently boost team stats
- **All items have stats**: Even if 0, always populated for consistency
- **Elemental affinity**: Mainly for artifacts/relics, provides bonus when rift has weakness
- **Power value**: Used for leaderboard "Dimensional Power Score" calculation

**Example Items:**

```
Weapon (Common): "Rusty Blade" - +3 speed, +1 power, +0 luck, fire affinity
Armor (Rare): "Glacial Plate" - +5 power, +2 luck, +1 speed, water affinity
Artifact (Epic): "Eternal Flame" - +10 power, +8 luck, +5 speed, fire affinity
Consumable (Uncommon): "Speed Elixir" - +2 speed, +0 luck, +1 power, none
Relic (Legendary): "Void Infinity Stone" - +20 power, +15 luck, +10 speed, void
```

---

### 4. **`loot_drop_tables`**

Defines drop rates for each rift and rarity tier.

```sql
CREATE TABLE loot_drop_tables (
    id SERIAL PRIMARY KEY,
    rift_id INT NOT NULL REFERENCES rifts(id) ON DELETE CASCADE,
    rarity item_rarity NOT NULL,

    -- Drop rates
    drop_rate_percent DECIMAL(5,2) NOT NULL CHECK (drop_rate_percent >= 0 AND drop_rate_percent <= 100),

    -- Quantity per expedition
    min_quantity INT NOT NULL DEFAULT 1 CHECK (min_quantity >= 1),
    max_quantity INT NOT NULL DEFAULT 1 CHECK (max_quantity >= min_quantity),

    UNIQUE(rift_id, rarity)
);

CREATE INDEX idx_drop_tables_rift ON loot_drop_tables(rift_id);
```

**Seed Data Examples:**

**Tutorial Rift:**

- Common: 100%, 3-3 items

**Easy Rifts (Crimson Wastes):**

- Common: 70%, 2-4 items
- Uncommon: 25%, 2-4 items
- Rare: 5%, 2-4 items

**Medium Rifts (Frozen/Neon):**

- Common: 50%, 3-5 items
- Uncommon: 30%, 3-5 items
- Rare: 15%, 3-5 items
- Epic: 5%, 3-5 items

**Hard Rifts (Verdant):**

- Common: 30%, 4-6 items
- Uncommon: 35%, 4-6 items
- Rare: 25%, 4-6 items
- Epic: 9%, 4-6 items
- Legendary: 1%, 4-6 items

**Legendary Rifts (Void):**

- Common: 10%, 5-8 items
- Uncommon: 25%, 5-8 items
- Rare: 35%, 5-8 items
- Epic: 25%, 5-8 items
- Legendary: 5%, 5-8 items

---

### 5. **`user_inventory`**

Stores player's collected loot (stacking by loot_item_id).

```sql
CREATE TABLE user_inventory (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loot_item_id INT NOT NULL REFERENCES loot_items(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, loot_item_id)
);

CREATE INDEX idx_inventory_user ON user_inventory(user_id);
CREATE INDEX idx_inventory_item ON user_inventory(loot_item_id);
CREATE INDEX idx_inventory_user_item ON user_inventory(user_id, loot_item_id);
```

**Notes:**

- One row per unique item per user (quantity stacks)
- `id` is used as FK in `teams` equipment slots
- When quantity reaches 0 (consumed), row can be deleted or kept for history

---

### 6. **`expeditions`**

Tracks all expedition instances (active and completed).

```sql
CREATE TABLE expeditions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    rift_id INT NOT NULL REFERENCES rifts(id) ON DELETE CASCADE,

    -- Timing
    start_time TIMESTAMP NOT NULL DEFAULT NOW(),
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),

    -- Status flags
    completed BOOLEAN NOT NULL DEFAULT false,  -- true when start_time + duration <= NOW()
    processed BOOLEAN NOT NULL DEFAULT false,  -- true when cron generated loot
    claimed BOOLEAN NOT NULL DEFAULT false,    -- true when user claimed rewards UI

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expeditions_user ON expeditions(user_id);
CREATE INDEX idx_expeditions_team ON expeditions(team_id);
CREATE INDEX idx_expeditions_status ON expeditions(completed, processed, claimed);
CREATE INDEX idx_expeditions_completion_check ON expeditions(start_time, duration_minutes)
    WHERE completed = false AND processed = false;
```

**Notes:**

- **completed**: Set by cron when time elapsed
- **processed**: Set by cron when loot generated and added to inventory
- **claimed**: Set by frontend when user views rewards (optional UX flag)
- Cron job queries for `completed=false AND start_time + duration <= NOW()`

---

### 7. **`expedition_loot`**

Records specific loot drops from each expedition (for debugging/analytics).

```sql
CREATE TABLE expedition_loot (
    id SERIAL PRIMARY KEY,
    expedition_id INT NOT NULL REFERENCES expeditions(id) ON DELETE CASCADE,
    loot_item_id INT NOT NULL REFERENCES loot_items(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expedition_loot_expedition ON expedition_loot(expedition_id);
CREATE INDEX idx_expedition_loot_item ON expedition_loot(loot_item_id);
```

**Notes:**

- Pure audit trail - not used in core game logic
- Helps debug: "Why did I get 3 Epic drops in one expedition?"
- Can power analytics: "Which rift gives most Legendary items?"

---

## Computed Values & Game Logic

### Team Total Stats Calculation

When calculating a team's effective stats for an expedition:

```
total_speed = team.speed_bonus
    + SUM(equipped_items.speed_bonus)
    + elemental_affinity_bonus (if artifact matches rift weakness)

total_luck = team.luck_bonus
    + SUM(equipped_items.luck_bonus)
    + elemental_affinity_bonus

total_power = team.power_bonus
    + SUM(equipped_items.power_bonus)
    + elemental_affinity_bonus
```

### Elemental Affinity Bonus

When a team has an equipped artifact/relic with elemental affinity matching the rift's weakness:

```
If equipped_artifact.elemental_affinity == rift.weak_to_element:
    total_power += (total_power * 0.20)  // 20% power boost
```

Example:

- Rift: Crimson Wastes (weak_to_element = 'water')
- Team has equipped artifact: "Glacial Heart" (elemental_affinity = 'water')
- Team gets +20% total power for this expedition

### Expedition Duration with Speed Bonus

```
actual_duration = rift.duration_minutes * (1 - (total_speed / 100))

Example:
- Rift duration: 60 minutes
- Team speed: +10% (0.10)
- Actual duration: 60 * (1 - 0.10) = 54 minutes
```

### Luck Effect on Drop Rates

Luck shifts drop rates toward higher rarities:

```
For each rarity tier (starting from lowest):
    adjusted_rate = base_rate - (base_rate * luck_modifier)
    shifted_amount += (base_rate * luck_modifier)

Next tier up:
    adjusted_rate += shifted_amount
```

Example with +10% luck on Medium Rift:

- Common: 50% → 45% (5% shifted up)
- Uncommon: 30% → 28% (3% shifted up, +5% from common) = 33%
- Rare: 15% → 16.5% (+3% from uncommon)
- Epic: 5% → 6% (+remainder)

### Leaderboard Power Score

```sql
SELECT
    user_id,
    SUM(loot_items.power_value * user_inventory.quantity) as power_score
FROM user_inventory
JOIN loot_items ON loot_items.id = user_inventory.loot_item_id
GROUP BY user_id
ORDER BY power_score DESC;
```

### Total Expeditions Count (for unlock checks)

```sql
SELECT COUNT(*)
FROM expeditions
WHERE user_id = ? AND completed = true;
```

---

## Migration Implementation Plan

### **Migration File: `20251021_game_core_schema.sql`**

**Step-by-step execution order:**

1. **Create ENUMs** (must come first, referenced by tables)

   - team_specialization
   - world_type
   - difficulty_level
   - elemental_type
   - item_rarity
   - item_type
   - equipment_slot_type

2. **Create `rifts` table** (no FK dependencies)

   - Seed 6 rifts with static data

3. **Create `loot_items` table** (references world_type enum)

   - Seed initial loot items (~75 items, can expand later)

4. **Create `loot_drop_tables` table** (references rifts, item_rarity)

   - Seed drop rates for all rift/rarity combinations

5. **Create `teams` table** (references users, user_inventory)

   - Note: user_inventory FK will be deferred until after inventory table created
   - Use `ON DELETE SET NULL` for equipment slots

6. **Create `user_inventory` table** (references users, loot_items)

   - Simple many-to-many with quantity

7. **Create `expeditions` table** (references users, teams, rifts)

   - Core expedition tracking

8. **Create `expedition_loot` table** (references expeditions, loot_items)

   - Audit trail for generated loot

9. **Add indexes** to all tables for performance

10. **Seed user teams** for existing users
    - Each user gets Team 1 unlocked by default

---

## Rollback Plan

If migration fails or needs to be reverted:

```sql
DROP TABLE IF EXISTS expedition_loot CASCADE;
DROP TABLE IF EXISTS expeditions CASCADE;
DROP TABLE IF EXISTS user_inventory CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS loot_drop_tables CASCADE;
DROP TABLE IF EXISTS loot_items CASCADE;
DROP TABLE IF EXISTS rifts CASCADE;

DROP TYPE IF EXISTS equipment_slot_type;
DROP TYPE IF EXISTS item_type;
DROP TYPE IF EXISTS item_rarity;
DROP TYPE IF EXISTS elemental_type;
DROP TYPE IF EXISTS difficulty_level;
DROP TYPE IF EXISTS world_type;
DROP TYPE IF EXISTS team_specialization;
```

**Order matters**: Drop tables first (in reverse dependency order), then drop types.

---

## Seed Data Requirements

### **Rifts** (6 rifts)

- Tutorial Rift
- Crimson Wastes (Fire)
- Frozen Expanse (Ice)
- Neon Sprawl (Tech)
- Verdant Overgrowth (Nature)
- Void Confluence (Void)

### **Loot Items** (~75+ items minimum)

- 15 items per world × 5 worlds
- Distribute across:
  - 5 rarities (common to legendary)
  - 2 item types (equipment vs consumable)
  - 5 equipment slots (weapon, armor, accessory, artifact, relic)
- Set appropriate stats based on rarity
- Assign elemental affinities (mainly for artifacts/relics)

### **Loot Drop Tables** (30 rows)

- 6 rifts × 5 rarities = 30 combinations
- Set drop rates per GDD specifications
- Set quantity ranges per difficulty

---

## Post-Migration Verification

**Run these queries to verify schema:**

```sql
-- Check all tables created
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;

-- Check all enums created
SELECT typname
FROM pg_type
WHERE typtype = 'e'
ORDER BY typname;

-- Verify rifts seeded
SELECT id, name, world_type, difficulty, duration_minutes
FROM rifts
ORDER BY difficulty, duration_minutes;

-- Verify loot items seeded (should be 75+)
SELECT rarity, item_type, equipment_slot, COUNT(*)
FROM loot_items
GROUP BY rarity, item_type, equipment_slot
ORDER BY rarity, item_type;

-- Verify drop tables seeded (should be 30 rows)
SELECT r.name, ldt.rarity, ldt.drop_rate_percent, ldt.min_quantity, ldt.max_quantity
FROM loot_drop_tables ldt
JOIN rifts r ON r.id = ldt.rift_id
ORDER BY r.difficulty, ldt.rarity;

-- Verify existing users have Team 1 unlocked
SELECT u.username, t.team_number, t.is_unlocked
FROM users u
LEFT JOIN teams t ON t.user_id = u.id AND t.team_number = 1;
```

---

## Performance Considerations

### **Expected Query Patterns:**

**High Frequency:**

- Check active expeditions for user (indexed on user_id)
- Get user inventory (indexed on user_id)
- Get team details with equipped items (indexed on user_id)
- Cron job: Find completed expeditions (indexed on completed, start_time)

**Medium Frequency:**

- Leaderboard queries (can add materialized view later)
- Loot generation (join drop tables with items)

**Low Frequency:**

- Create new expedition
- Equip/unequip items
- Consume items for upgrades

### **Optimization Notes:**

- All foreign keys have indexes
- Composite indexes for common query patterns
- `user_inventory` uses UNIQUE constraint for upsert operations
- Expedition completion check has dedicated partial index
- Consider materialized view for leaderboard if it becomes slow (Week 2/3)

---

## Future Enhancements (Post-MVP)

**Potential Schema Extensions:**

1. **Achievements System**

   - `achievements` table (static)
   - `user_achievements` table (progress tracking)

2. **Echo Encounters** (Week 3 feature)

   - `echo_encounters` table
   - `notifications` table

3. **Leaderboard Cache**

   - Materialized view or dedicated cache table
   - Refreshed by cron job

4. **Team Presets**

   - Save/load equipment configurations
   - `team_presets` table

5. **Item Crafting** (if scope expands)
   - `crafting_recipes` table
   - Track combined items

---

## Summary

**Total Tables:** 7 core tables + `users` (existing)  
**Total ENUMs:** 7 types  
**Estimated Seed Data:**

- 6 rifts
- 75+ loot items
- 30 drop table entries

**Key Features Enabled:**
✅ Teams with 5 equipment slots (typed)  
✅ Expeditions with timed completion  
✅ Loot generation from drop tables  
✅ Inventory management with stacking  
✅ Equipment with stat bonuses  
✅ Consumables for permanent upgrades  
✅ Elemental affinity system  
✅ Leaderboard power calculations  
✅ Full expedition audit trail

**Ready for Implementation:** ✅

---

**Next Steps:**

1. Create migration file: `migrations/20251021_game_core_schema.sql`
2. Write seed data for rifts, loot_items, and loot_drop_tables
3. Run migration: `make migrate`
4. Verify with test queries
5. Begin repository/service implementation
