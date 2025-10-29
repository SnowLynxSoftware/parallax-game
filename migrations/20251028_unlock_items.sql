-- ############################
-- Migration: Add Unlock Items
-- Date: 2025-10-28
--
-- Purpose: Add Golden Fishing Rod and Explorers Compass as common
-- weapon items that unlock Fishing and Dungeons navigation.
-- 
-- These items are seeded into the loot pool with world_type='tutorial'
-- so they will be dropped from any rift with matching world_type.
-- Tutorial Rift (ID 1) and Crimson Wastes (ID 2) both drop tutorial items.
-- ############################

-- Add the Golden Fishing Rod (unlocks Fishing feature)
INSERT INTO loot_items (
    name, 
    description, 
    rarity, 
    world_type, 
    item_type, 
    equipment_slot, 
    speed_bonus, 
    luck_bonus, 
    power_bonus, 
    elemental_affinity, 
    power_value, 
    icon
) VALUES (
    'Golden Fishing Rod',
    'A shimmering rod of pure gold. Its presence in your inventory unlocks the Fishing feature.',
    'common',
    'tutorial',
    'equipment',
    'weapon',
    1.00,
    0.00,
    1,
    'none',
    5,
    'fa-fishing'
);

-- Add the Explorers Compass (unlocks Dungeons feature)
INSERT INTO loot_items (
    name, 
    description, 
    rarity, 
    world_type, 
    item_type, 
    equipment_slot, 
    speed_bonus, 
    luck_bonus, 
    power_bonus, 
    elemental_affinity, 
    power_value, 
    icon
) VALUES (
    'Explorers Compass',
    'An ornate compass that points to hidden dungeons. Its presence in your inventory unlocks the Dungeons feature.',
    'common',
    'tutorial',
    'equipment',
    'weapon',
    1.00,
    0.00,
    1,
    'none',
    5,
    'fa-compass'
);

-- NOTE: Loot drop tables are already seeded in the initial schema migration.
-- These items will automatically be included in common drops from:
-- - Tutorial Rift (ID 1): world_type='tutorial', common drop rate 100%
-- - Crimson Wastes (ID 2): world_type='fire' (not tutorial), so these items 
--   will NOT drop from Crimson Wastes
--
-- The expedition reward generation selects items by matching world_type.
-- To make these items available from Crimson Wastes as well, they would need
-- world_type='fire' instead. Currently seeded as 'tutorial' for early availability.
