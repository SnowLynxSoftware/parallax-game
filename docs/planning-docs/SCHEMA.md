# Database Schema Plan - Parallax Game

**Version:** 1.0  
**Date:** October 21, 2025  
**Status:** Ready for Implementation

---

## Implementation Checklist

- [x] 1. Create migration file: `migrations/20251021_game_core_schema.sql`
- [x] 2. Create ENUM types (world_type, difficulty_level, elemental_type, item_rarity, item_type, equipment_slot_type)
- [x] 3. Create `rifts` table and seed 7 rifts with static data
- [x] 4. Create `loot_items` table and seed ~75 loot items (relics only have non-'none' elemental_affinity)
- [x] 5. Create `loot_drop_tables` table and seed drop rates for all rift/rarity combinations
- [x] 6. Create `teams` table with all base fields
- [x] 7. Create `user_inventory` table with all base fields
- [x] 8. Create `expeditions` table with all base fields
- [x] 9. Create `expedition_loot` table with all base fields
- [x] 10. Add all indexes to tables for performance
- [x] 11. Seed user teams (Team 1 unlocked by default for existing users)
- [x] 12. Run migration: `make migrate`
- [x] 13. Verify schema with test queries (check all tables, enums, seed data)
- [x] 14. Begin repository/service implementation

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
2. **Exclusive Equipment**: Each item can only be equipped to one team at a time. Equipping to Team A unequips it from Team B (managed in Go logic).
3. **Type Safety**: Enums for rarities, world types, item types, slot types, elements
4. **Audit Trail**: Track expedition loot for debugging and analytics
5. **Scalability**: Designed for 50-100 concurrent users, optimized for read-heavy queries

---

## Table Definitions

### 1. **`teams`**

Stores user's expedition teams (up to 5 per user).

```sql
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_number INT NOT NULL CHECK (team_number >= 1 AND team_number <= 5),

    -- Base stats (permanent upgrades from consumables)
    speed_bonus DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    luck_bonus DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    power_bonus INT NOT NULL DEFAULT 0,

    -- Team configuration
    is_unlocked BOOLEAN NOT NULL DEFAULT false,

    -- Equipment slots (FK to user_inventory.id)
    equipped_weapon_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_armor_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_accessory_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_artifact_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,
    equipped_relic_slot INT REFERENCES user_inventory(id) ON DELETE SET NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false,

    UNIQUE(user_id, team_number)
);

CREATE INDEX idx_teams_user_id ON teams(user_id);
CREATE INDEX idx_teams_unlocked ON teams(user_id, is_unlocked);
```

**Notes:**

- Each user has 5 team slots (unlocked progressively via game logic)
- Base stats are permanent upgrades from consuming items
- Equipment slots reference `user_inventory.id` (not `loot_items.id`)
- Each `user_inventory.id` can only be equipped to ONE team at a time (one row per team's slot)
- When a user equips item X to Team A's weapon slot, it must be unequipped from any other team first (handled in Go layer)
- All tables follow standard base fields: `id`, `created_at`, `modified_at`, `is_archived`

**Important Constraint (Enforced in Go):**

- Each `user_inventory.id` (especially equipment) can appear in **at most ONE** equipment slot across all teams.
- When equipping an item to a new team, the Go layer must first check if that item is already equipped elsewhere and unequip it.
- The database does NOT enforce this uniqueness (would require a complex trigger), so **it is a business logic contract** that the application must maintain.
- If this constraint is violated, queries will return inconsistent team stats.

**Equipment Inventory Semantics:**

- Equipment loot_items can have multiple `user_inventory` rows per user (no UNIQUE constraint).
- Each row is a separate "copy" of that item in the player's inventory.
- When a "Rusty Blade" drops twice, the player gets two rows with `loot_item_id=X` and `quantity=1` each.
- When equipping, a specific `user_inventory.id` is chosen (not just the loot_item_id).

**Recommendation for future:** If bugs arise from this contract, add a CHECK trigger or a separate `equipment_assignments` tracking table.

---

### 2. **`rifts`** (Static/Seed Data)

Defines the parallel world rifts players can explore.

```sql
CREATE TYPE world_type AS ENUM ('tutorial', 'fire', 'ice', 'tech', 'nature', 'void');
CREATE TYPE difficulty_level AS ENUM ('tutorial', 'easy', 'medium', 'hard', 'legendary');
CREATE TYPE elemental_type AS ENUM ('fire', 'water', 'wind', 'earth', 'void', 'light', 'none');

CREATE TABLE rifts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    world_type world_type NOT NULL,
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
    difficulty difficulty_level NOT NULL,

    -- Elemental weakness (for relic bonuses)
    weak_to_element elemental_type NOT NULL DEFAULT 'none',

    -- Unlock requirements (text description, enforced in Go)
    unlock_requirement_text VARCHAR(255),

    -- UI
    icon VARCHAR(50) NOT NULL, -- Font Awesome class

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_rifts_difficulty ON rifts(difficulty);
CREATE INDEX idx_rifts_world_type ON rifts(world_type);
```

**Seed Data (7 Rifts):**

1. Tutorial Rift - 5min, tutorial difficulty, no weakness
2. Crimson Wastes (Fire) - 15min, easy, weak to water
3. Frozen Expanse (Ice) - 30min, medium, weak to fire
4. Neon Sprawl (Tech) - 30min, medium, weak to wind
5. Verdant Overgrowth (Nature) - 60min, hard, weak to earth
6. Void Confluence (Void) - 120min, legendary, weak to light
7. Seat of Heaven (Light) - 120min, legendary, weak to void

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

    -- Elemental affinity (ONLY for relics; all other items must be 'none')
    elemental_affinity elemental_type NOT NULL DEFAULT 'none',

    -- Power value for leaderboard calculations
    power_value INT NOT NULL DEFAULT 1,

    -- UI
    icon VARCHAR(50) NOT NULL, -- Font Awesome class

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false
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
- **Elemental affinity**: ONLY relics can have an elemental affinity other than 'none'. All weapons, armor, accessories, artifacts, and consumables MUST have `elemental_affinity = 'none'`. Your relic determines which element you are attuned to.
- **Power value**: Used for leaderboard "Dimensional Power Score" calculation

**Example Items:**

```
Weapon (Common): "Rusty Blade" - +3 speed, +1 power, +0 luck, none
Armor (Rare): "Glacial Plate" - +5 power, +2 luck, +1 speed, none
Artifact (Epic): "Eternal Flame" - +10 power, +8 luck, +5 speed, none
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

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false,

    UNIQUE(rift_id, rarity)
);

CREATE INDEX idx_drop_tables_rift ON loot_drop_tables(rift_id);
```

**Seed Data Examples:**

**Tutorial Rift:**

- Common: 100%, 1-2 items (only weapon or armor types)

**Easy Rifts (Crimson Wastes):**

- Common: 70%, 2-4 items
- Uncommon: 25%, 2-4 items
- Rare: 5%, 1 items

**Medium Rifts (Frozen/Neon):**

- Common: 50%, 3-5 items
- Uncommon: 30%, 3-5 items
- Rare: 15%, 3-5 items
- Epic: 5%, 1 items

**Hard Rifts (Verdant):**

- Common: 30%, 4-6 items
- Uncommon: 35%, 4-6 items
- Rare: 25%, 4-6 items
- Epic: 9%, 4-6 items
- Legendary: 1%, 1 items

**Legendary Rifts (Void):**

- Common: 10%, 5-8 items
- Uncommon: 25%, 5-8 items
- Rare: 35%, 5-8 items
- Epic: 25%, 5-8 items
- Legendary: 5%, 1 items (Light Relics)

**Legendary Rifts (Light):**

- Common: 10%, 5-8 items
- Uncommon: 25%, 5-8 items
- Rare: 35%, 5-8 items
- Epic: 25%, 5-8 items
- Legendary: 5%, 1 items (Void Relics)

---

### 5. **`user_inventory`**

Stores player's collected loot. Equipment items are never stacked; consumables stack via `quantity`.

```sql
CREATE TABLE user_inventory (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loot_item_id INT NOT NULL REFERENCES loot_items(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity >= 1),
    acquired_at TIMESTAMP NOT NULL DEFAULT NOW(),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_inventory_user ON user_inventory(user_id);
CREATE INDEX idx_inventory_item ON user_inventory(loot_item_id);
CREATE INDEX idx_inventory_user_item ON user_inventory(user_id, loot_item_id);
```

**Notes:**

- **Equipment items:** Each drop is a separate row (no UNIQUE constraint). A player can have multiple copies of the same weapon. `quantity` is always 1 for equipment.
- **Consumable items:** Stacked into a single row per loot_item per user. `quantity` increments when duplicates are acquired, decrements when used.
- `id` is used as FK in `teams` equipment slots (only equipment rows are referenced).
- **Equipment semantics:** When an item is equipped to a team, its `user_inventory.id` is stored in that team's slot. Only ONE team can reference that `id` at a time (managed in Go logic via application-level checks).
- **Consumable semantics:** When a consumable is used, `quantity` is decremented. If `quantity` reaches 0, the row is deleted (via DELETE statement, not soft-delete).
- `acquired_at` tracks when the item first entered inventory (useful for UI sorting)
- **Go Validation:** On inventory upsert, if loot_item is equipment type, enforce `quantity = 1`. If consumable, allow stack increments.

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

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false
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
- **Duration is stored at expedition creation** (no need to recalculate at completion)
- Cron job queries: `SELECT * FROM expeditions WHERE completed = false AND start_time + INTERVAL '1 minute' * duration_minutes <= NOW()`
- On expedition completion: Cron calculates loot, adds to `user_inventory`, creates `expedition_loot` audit records, sets `completed=true` then `processed=true`

---

### 7. **`expedition_loot`**

Records specific loot drops from each expedition (for debugging/analytics).

```sql
CREATE TABLE expedition_loot (
    id SERIAL PRIMARY KEY,
    expedition_id INT NOT NULL REFERENCES expeditions(id) ON DELETE CASCADE,
    loot_item_id INT NOT NULL REFERENCES loot_items(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false
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

When a team has an equipped relic with elemental affinity matching the rift's weakness:

```
If equipped_relic.elemental_affinity == rift.weak_to_element:
    total_power += (total_power * 0.20)  // 20% power boost to total AFTER all other equipped items added.
```

Example:

- Rift: Crimson Wastes (weak_to_element = 'water')
- Team has equipped relic: "Glacial Heart" (elemental_affinity = 'water')
- Team gets +20% total power for this expedition

### Expedition Duration with Speed Bonus

```
actual_duration_minutes = MAX(5, rift.duration_minutes * (1 - (total_speed / 100)))

Example:
- Rift duration: 60 minutes
- Team speed: +10% (0.10)
- Actual duration: MAX(5, 60 * (1 - 0.10)) = MAX(5, 54) = 54 minutes

Example with high speed:
- Rift duration: 60 minutes
- Team speed: +95% (0.95)
- Actual duration: MAX(5, 60 * (1 - 0.95)) = MAX(5, 3) = 5 minutes (floor enforced)
```

**Important:** The 5-minute floor prevents negative or zero durations.

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

### Loot Generation Algorithm

Each expedition triggers a **single loot roll** at start time that determines all items the player will receive:

1. Read `loot_drop_tables` for the rift
2. For each rarity tier (common → legendary):
   - Apply luck modifier to shift percentages
   - Roll random value 0-100
   - If roll < adjusted_rate: This rarity "drops"
3. For each dropping rarity:
   - Roll quantity between `min_quantity` and `max_quantity`
   - Roll that many item IDs from loot_items matching the rarity and rift's world_type
4. Store all rolled items in `expedition_loot` and sum into `user_inventory`

**Key Design:** A single expedition generates 1-N items based on the drop table. Rarity probabilities determine _which tiers_ drop, quantity ranges determine _how many_ from each tier.

### Leaderboard Power Score

**Only equipped items count:**

```sql
SELECT
    u.id as user_id,
    SUM(li.power_value) as power_score
FROM users u
LEFT JOIN teams t ON t.user_id = u.id AND t.is_unlocked = true
LEFT JOIN user_inventory wi ON wi.id = t.equipped_weapon_slot
LEFT JOIN user_inventory ai ON ai.id = t.equipped_armor_slot
LEFT JOIN user_inventory aci ON aci.id = t.equipped_accessory_slot
LEFT JOIN user_inventory afi ON afi.id = t.equipped_artifact_slot
LEFT JOIN user_inventory ri ON ri.id = t.equipped_relic_slot
LEFT JOIN loot_items li ON li.id IN (wi.loot_item_id, ai.loot_item_id, aci.loot_item_id, afi.loot_item_id, ri.loot_item_id)
GROUP BY u.id
ORDER BY power_score DESC;
```

**Notes:**

- Only unlocked teams count
- Only equipped items (exactly 5 per team, but some may be NULL)
- Sum of `power_value` across all equipped items across all teams
- Consumables in inventory do NOT count (they don't have equipment_slot defined)

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

   - world_type
   - difficulty_level
   - elemental_type
   - item_rarity
   - item_type
   - equipment_slot_type

2. **Create `rifts` table** (no FK dependencies)

   - Seed 7 rifts with static data

3. **Create `loot_items` table** (references world_type enum)

   - Seed initial loot items (~75 items, can expand later)
   - **IMPORTANT**: Only relics can have elemental_affinity != 'none'. All other items must have `elemental_affinity = 'none'`

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
```

**Order matters**: Drop tables first (in reverse dependency order), then drop types.

---

## Seed Data Requirements

### **Rifts** (2 rifts)

- Tutorial Rift
- Crimson Wastes (Fire)
- Frozen Expanse (Ice)
- Neon Sprawl (Tech)
- Verdant Overgrowth (Nature)
- Void Confluence (Void)
- Seat of Heaven (Light)

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
**Total ENUMs:** 6 types  
**Estimated Seed Data:**

---

## Critical Implementation Notes for Go Services

### Equipment Equip/Unequip Logic

**When equipping an item to a team's slot:**

```
1. Validate: user_id owns the inventory item
2. Validate: user_id owns the team
3. Validate: loot_item type matches equipment slot
4. Query: SELECT * FROM teams WHERE user_id = ? (all 5 teams)
5. Check: Does any team currently have this inventory.id in ANY slot?
6. If yes: SET that slot to NULL (unequip from other team)
7. SET target team's slot to this inventory.id
8. Return: New team stats (for UI update)
```

**Why step 4-6 matters:** Without this multi-team check, the same item can end up equipped to multiple teams, breaking the constraint.

### Loot Generation Service (When Expedition Completes)

**Input:** `expedition_id`, `completed_expedition` row with team_id, rift_id, duration_minutes

**Output:** Loot added to `user_inventory`, audit trail in `expedition_loot`

```
1. Load expedition, team, rift
2. Calculate team stats at expedition START time:
   - Get equipped items from team's 5 slots
   - Sum stats from all equipped items
   - Apply elemental bonus if relic matches rift weakness
3. Roll loot:
   a. Load drop tables for rift
   b. For each rarity tier:
      - Apply luck modifier
      - Roll 0-100
      - If roll < adjusted_rate: rarity drops
   c. For each dropping rarity:
      - Roll quantity between min and max
      - For each unit: roll a random loot_item matching rarity + world_type
4. For each rolled item:
   - Upsert into user_inventory (increment quantity if exists)
   - INSERT into expedition_loot (for audit trail)
5. Set expedition.processed = true
6. Return: Summary of loot received (for notification/UI)
```

**Key gotcha:** If a consumable is used between expedition start and completion, the team stats _should not recalculate_—they were locked at start. Store them in the expedition record if you need to audit "what stats did this expedition use?"

### Consumable Use (Upgrade) Logic

**When player consumes an item:**

```
1. Validate: user_id owns the item (user_inventory row)
2. Load loot_item to determine type (must be consumable)
3. Load team
4. Calculate new team bonus:
   - team.speed_bonus += loot_item.speed_bonus
   - team.luck_bonus += loot_item.luck_bonus
   - team.power_bonus += loot_item.power_bonus
5. Decrement user_inventory.quantity by 1
6. If quantity reaches 0: DELETE the row
7. Return: Updated team stats
```

**Important:** Consumables can be used on any team (the bonus is team-specific, stored in `teams.speed_bonus` etc.). So a Speed Elixir used on Team 1 doesn't affect Team 2.

### Equipment Inventory Logic

**When equipment is dropped (during loot generation):**

```
1. Roll loot for a rarity (e.g., Rusty Blade, common weapon)
2. Create NEW user_inventory row:
   - user_id = expediting player
   - loot_item_id = item ID
   - quantity = 1 (always 1 for equipment)
   - acquired_at = NOW()
3. NO UNIQUE constraint check—allow duplicates
4. Return: "You found a Rusty Blade!"
```

**When equipment is already in inventory (duplicate drop):**

```
1. Check: Is this loot_item an equipment type?
2. If yes: Create a separate row (don't increment quantity)
3. If consumable: Upsert (increment quantity on existing row, or create new)
```

**Go Validation on Upsert:**

- If `loot_item.item_type = 'equipment'`: Always INSERT new row with `quantity = 1`
- If `loot_item.item_type = 'consumable'`: INSERT or UPDATE to increment quantity
- Reject any attempt to insert equipment with `quantity != 1`

---

## Revised Data Integrity Risks & Mitigations

| Risk                                         | Mitigation                                                      |
| -------------------------------------------- | --------------------------------------------------------------- |
| Same equipment instance equipped to 2+ teams | Go layer maintains constraint via select-before-equip check     |
| Equipment quantity > 1 in inventory          | Go validation: equipment always quantity = 1, consumables stack |
| Negative/zero expedition duration            | Database-enforced 5-min floor in Go (MAX(5, calculated))        |
| Relic with non-matching elemental affinity   | Seed data validated; optionally add CHECK trigger               |
| Quantity 0 in inventory                      | Immediate DELETE on consume, not soft-delete                    |
| Consumed item affects all teams              | Bonus stored per-team in `teams.speed_bonus`, so isolated       |
| Cron processes same expedition twice         | Use `processed` flag: only process if `processed = false`       |

---

## Revised Summary

**Total Tables:** 7 core tables + `users` (existing)  
**Total ENUMs:** 6 types  
**Estimated Seed Data:**

- 7 rifts
- 75+ loot items (only relics have non-'none' elemental_affinity)
- 35 drop table entries
- 1 team per existing user (Team 1 unlocked)

**Key Features Enabled:**
✅ Teams with 5 exclusive equipment slots (typed)  
✅ Expeditions with timed completion and locked-in stats  
✅ Single-roll loot generation from drop tables  
✅ Equipment as individual inventory rows (duplicates allowed, quantity=1)  
✅ Consumables with stacking (quantity increments)  
✅ Equipment shared only via unequip/re-equip (not simultaneous)  
✅ Consumables for permanent team upgrades  
✅ Elemental affinity system (relics only, 20% power bonus)  
✅ 5-minute minimum expedition duration  
✅ Leaderboard based on equipped items only  
✅ Full expedition audit trail  
✅ Consistent base fields across all tables

**Design Trade-off:** Equipment exclusivity is enforced in Go, not the database. This is simpler but requires careful code review and tests to avoid bugs.

**Design Decision:** No UNIQUE constraint on `(user_id, loot_item_id)`. Equipment duplicates create separate rows; consumables upsert on collision.

**Ready for Implementation:** ✅

---

---

## Go Implementation Guidelines

### Repository Methods Required

**UserInventoryRepository:**

```go
// AddLoot handles both equipment and consumable loot
// Equipment: Always INSERT new row with quantity=1
// Consumable: Upsert (INSERT or UPDATE increment quantity)
AddLoot(ctx, userID, lootItemID) (*UserInventory, error)

// ConsumeLoot decrements quantity, deletes if 0
ConsumeLoot(ctx, inventoryID) error

// GetByUserAndItem for consumable lookups
GetByUserAndItem(ctx, userID, lootItemID) (*UserInventory, error)

// GetEquippedItems returns only items currently in team slots
GetEquippedItems(ctx, userID) ([]UserInventory, error)
```

**TeamsRepository:**

```go
// EquipItem must check all other teams first, unequip if already equipped elsewhere
EquipItem(ctx, teamID, slot, inventoryID) error

// UnequipItem sets slot to NULL
UnequipItem(ctx, teamID, slot) error
```

**ExpeditionsRepository:**

```go
// FindReadyToProcess for cron job
FindReadyToProcess(ctx) ([]Expedition, error)

// MarkCompleted and MarkProcessed for cron workflow
MarkCompleted(ctx, expeditionID) error
MarkProcessed(ctx, expeditionID) error
```

### Critical Validation Rules

**On every AddLoot operation:**

```go
lootItem := getLootItem(lootItemID)

if lootItem.ItemType == "equipment" {
    // Always INSERT new row, never increment quantity
    inventory := UserInventory{
        UserID:     userID,
        LootItemID: lootItemID,
        Quantity:   1,  // MUST be 1
        AcquiredAt: now,
    }
    db.Insert(inventory)
} else if lootItem.ItemType == "consumable" {
    // Upsert: INSERT or UPDATE increment
    existing := GetByUserAndItem(userID, lootItemID)
    if existing != nil {
        IncrementQuantity(existing.ID)
    } else {
        Insert(UserInventory{Quantity: 1})
    }
}
```

**On every EquipItem operation:**

```go
// Check if inventory_id is already equipped anywhere
result := db.Query(`
    SELECT team_id,
           CASE
               WHEN equipped_weapon_slot = ? THEN 'weapon'
               WHEN equipped_armor_slot = ? THEN 'armor'
               WHEN equipped_accessory_slot = ? THEN 'accessory'
               WHEN equipped_artifact_slot = ? THEN 'artifact'
               WHEN equipped_relic_slot = ? THEN 'relic'
           END as slot
    FROM teams
    WHERE user_id = ? AND (
        equipped_weapon_slot = ? OR
        equipped_armor_slot = ? OR
        equipped_accessory_slot = ? OR
        equipped_artifact_slot = ? OR
        equipped_relic_slot = ?
    )
`, inventoryID, inventoryID, inventoryID, inventoryID, inventoryID, userID,
   inventoryID, inventoryID, inventoryID, inventoryID, inventoryID)

if result != nil {
    // Unequip from old team first
    db.Update("UPDATE teams SET ?_slot = NULL WHERE team_id = ?", result.Slot, result.TeamID)
}

// Now equip to new team
db.Update("UPDATE teams SET ?_slot = ? WHERE team_id = ?", slot, inventoryID, teamID)
```

### Query Patterns

**Get all equipment for user:**

```sql
SELECT ui.*, li.name, li.rarity, li.equipment_slot, li.power_value
FROM user_inventory ui
JOIN loot_items li ON li.id = ui.loot_item_id
WHERE ui.user_id = ? AND li.item_type = 'equipment'
ORDER BY ui.acquired_at DESC;
```

**Get consumables with quantities:**

```sql
SELECT ui.*, li.name, li.speed_bonus, li.luck_bonus, li.power_bonus
FROM user_inventory ui
JOIN loot_items li ON li.id = ui.loot_item_id
WHERE ui.user_id = ? AND li.item_type = 'consumable'
ORDER BY li.name;
```

**Check if equipment is equipped anywhere:**

```sql
SELECT COUNT(*) FROM teams
WHERE user_id = ? AND (
    equipped_weapon_slot = ? OR
    equipped_armor_slot = ? OR
    equipped_accessory_slot = ? OR
    equipped_artifact_slot = ? OR
    equipped_relic_slot = ?
);
```

**Leaderboard (equipped items only):**

```sql
SELECT
    u.id as user_id,
    u.username,
    SUM(li.power_value) as power_score
FROM users u
LEFT JOIN teams t ON t.user_id = u.id AND t.is_unlocked = true
LEFT JOIN user_inventory wi ON wi.id = t.equipped_weapon_slot
LEFT JOIN user_inventory ai ON ai.id = t.equipped_armor_slot
LEFT JOIN user_inventory aci ON aci.id = t.equipped_accessory_slot
LEFT JOIN user_inventory afi ON afi.id = t.equipped_artifact_slot
LEFT JOIN user_inventory ri ON ri.id = t.equipped_relic_slot
LEFT JOIN loot_items li ON li.id IN (
    wi.loot_item_id, ai.loot_item_id, aci.loot_item_id,
    afi.loot_item_id, ri.loot_item_id
)
GROUP BY u.id, u.username
ORDER BY power_score DESC;
```

---

## Next Steps

1. Create migration file: `migrations/20251021_game_core_schema.sql`
2. Write seed data for rifts, loot_items, and loot_drop_tables
3. Run migration: `make migrate`
4. Verify with test queries
5. Begin repository/service implementation
   - Expeditions repository (create, complete, loot generation)
   - Teams repository (get, equip/unequip with validation)
   - Inventory repository (get, consume, upsert with type-aware logic)
   - Game services (stats calculation, loot generation, leaderboard)

```

```
