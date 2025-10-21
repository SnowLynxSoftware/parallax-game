-- ############################
-- Parallax Game Core Schema
--
-- https://snowlynxsoftware.net 
--
-- Copyright 2025. Snow Lynx Software, LLC. All Rights Reserved.
-- ############################

-- GAME CORE SCHEMA: Rifts, Teams, Expeditions, Loot, Inventory

-- ############################
-- STEP 1: CREATE ENUM TYPES
-- ############################

-- World types for rifts and loot classification
CREATE TYPE world_type AS ENUM (
    'tutorial',
    'fire',
    'ice',
    'tech',
    'nature',
    'void',
    'light'
);

-- Difficulty levels for rifts
CREATE TYPE difficulty_level AS ENUM (
    'tutorial',
    'easy',
    'medium',
    'hard',
    'legendary'
);

-- Elemental types for relics and rift weaknesses
CREATE TYPE elemental_type AS ENUM (
    'fire',
    'water',
    'wind',
    'earth',
    'void',
    'light',
    'none'
);

-- Item rarity tiers for loot drops
CREATE TYPE item_rarity AS ENUM (
    'common',
    'uncommon',
    'rare',
    'epic',
    'legendary'
);

-- Item types: equipment vs consumable
CREATE TYPE item_type AS ENUM (
    'equipment',
    'consumable'
);

-- Equipment slot types for teams
CREATE TYPE equipment_slot_type AS ENUM (
    'weapon',
    'armor',
    'accessory',
    'artifact',
    'relic'
);

-- ############################
-- STEP 3: CREATE RIFTS TABLE
-- ############################

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

-- Seed 7 rifts with static data
INSERT INTO rifts (name, description, world_type, duration_minutes, difficulty, weak_to_element, unlock_requirement_text, icon) VALUES
    ('Tutorial Rift', 'A safe training ground to learn the basics of dimensional exploration. Perfect for new explorers.', 'tutorial', 5, 'tutorial', 'none', NULL, 'fa-graduation-cap'),
    ('Crimson Wastes', 'A scorched desert world where rivers of lava flow beneath crimson skies. The heat is oppressive, but fire-attuned relics lose their power here.', 'fire', 15, 'easy', 'water', 'Complete Tutorial Rift', 'fa-fire'),
    ('Frozen Expanse', 'An endless tundra locked in eternal winter. Ice storms rage constantly, but water-based powers hold sway.', 'ice', 30, 'medium', 'fire', 'Complete 5 expeditions', 'fa-snowflake'),
    ('Neon Sprawl', 'A cyberpunk megacity frozen in time, where holographic advertisements flicker in abandoned streets. Technology reigns supreme, but vulnerable to natural forces.', 'tech', 30, 'medium', 'wind', 'Complete 5 expeditions', 'fa-microchip'),
    ('Verdant Overgrowth', 'Nature has reclaimed this world completely. Massive trees pierce the clouds, and ancient ruins are wrapped in vines. Earth powers dominate here.', 'nature', 60, 'hard', 'earth', 'Complete 15 expeditions', 'fa-leaf'),
    ('Void Confluence', 'A place where reality itself breaks down. Floating islands drift in an endless dark void. Only light can pierce this darkness.', 'void', 120, 'legendary', 'light', 'Complete 30 expeditions', 'fa-circle-notch'),
    ('Seat of Heaven', 'A radiant realm of pure crystalline light. Celestial architecture defies physics. The void is anathema here.', 'light', 120, 'legendary', 'void', 'Complete 30 expeditions', 'fa-sun');

-- ############################
-- STEP 4: CREATE LOOT_ITEMS TABLE
-- ############################

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

-- Seed loot items (~75 items across all worlds and rarities)
-- TUTORIAL WORLD ITEMS (Starter gear)
INSERT INTO loot_items (name, description, rarity, world_type, item_type, equipment_slot, speed_bonus, luck_bonus, power_bonus, elemental_affinity, power_value, icon) VALUES
    -- Tutorial Weapons
    ('Wooden Sword', 'A simple training blade. Everyone starts somewhere.', 'common', 'tutorial', 'equipment', 'weapon', 1.00, 0.00, 2, 'none', 5, 'fa-sword'),
    ('Practice Bow', 'A basic bow for target practice. Surprisingly effective.', 'common', 'tutorial', 'equipment', 'weapon', 2.00, 1.00, 1, 'none', 5, 'fa-bow-arrow'),
    -- Tutorial Armor
    ('Leather Vest', 'Standard issue protection for new explorers.', 'common', 'tutorial', 'equipment', 'armor', 0.00, 1.00, 3, 'none', 5, 'fa-vest'),
    ('Training Gloves', 'Padded gloves that protect your hands.', 'common', 'tutorial', 'equipment', 'armor', 1.00, 0.00, 2, 'none', 5, 'fa-hand-back-fist'),
    -- Tutorial Accessories
    ('Explorer Badge', 'Proof of your training completion.', 'common', 'tutorial', 'equipment', 'accessory', 2.00, 2.00, 1, 'none', 8, 'fa-medal');

-- FIRE WORLD ITEMS (Crimson Wastes)
INSERT INTO loot_items (name, description, rarity, world_type, item_type, equipment_slot, speed_bonus, luck_bonus, power_bonus, elemental_affinity, power_value, icon) VALUES
    -- Fire Weapons
    ('Scorched Blade', 'A sword tempered in volcanic heat.', 'common', 'fire', 'equipment', 'weapon', 2.00, 0.00, 5, 'none', 10, 'fa-sword'),
    ('Ember Staff', 'Channels the power of smoldering coals.', 'uncommon', 'fire', 'equipment', 'weapon', 3.00, 2.00, 8, 'none', 20, 'fa-wand-magic-sparkles'),
    ('Flameburst Hammer', 'Each strike releases a burst of flame.', 'rare', 'fire', 'equipment', 'weapon', 4.00, 3.00, 15, 'none', 40, 'fa-hammer'),
    ('Inferno Scythe', 'Wreathed in perpetual flame.', 'epic', 'fire', 'equipment', 'weapon', 6.00, 5.00, 25, 'none', 80, 'fa-sickle'),
    -- Fire Armor
    ('Ash Cloak', 'Protects against heat and embers.', 'common', 'fire', 'equipment', 'armor', 0.00, 1.00, 6, 'none', 10, 'fa-shirt'),
    ('Magma Plate', 'Armor forged from cooled lava.', 'uncommon', 'fire', 'equipment', 'armor', 1.00, 2.00, 10, 'none', 20, 'fa-shield'),
    ('Phoenix Mail', 'Said to grant the resilience of the legendary bird.', 'rare', 'fire', 'equipment', 'armor', 2.00, 4.00, 18, 'none', 40, 'fa-shield-halved'),
    -- Fire Accessories
    ('Cinder Ring', 'Warm to the touch, grants minor protection.', 'common', 'fire', 'equipment', 'accessory', 3.00, 2.00, 3, 'none', 12, 'fa-ring'),
    ('Obsidian Amulet', 'Sharp volcanic glass shaped into jewelry.', 'uncommon', 'fire', 'equipment', 'accessory', 4.00, 4.00, 6, 'none', 22, 'fa-gem'),
    ('Volcanic Heart', 'Pulses with inner heat.', 'rare', 'fire', 'equipment', 'accessory', 6.00, 5.00, 12, 'none', 42, 'fa-heart'),
    -- Fire Artifacts
    ('Eternal Flame', 'A flame that never dies, contained in crystal.', 'epic', 'fire', 'equipment', 'artifact', 10.00, 8.00, 30, 'none', 90, 'fa-fire-flame-curved'),
    -- Fire Relics
    ('Inferno Stone', 'The essence of fire itself, crystallized.', 'legendary', 'fire', 'equipment', 'relic', 15.00, 12.00, 50, 'fire', 150, 'fa-fire'),
    -- Fire Consumables
    ('Fire Elixir', 'Grants permanent power by absorbing fire energy.', 'uncommon', 'fire', 'consumable', NULL, 0.00, 0.00, 5, 'none', 15, 'fa-flask'),
    ('Blaze Essence', 'Pure fire energy in liquid form.', 'rare', 'fire', 'consumable', NULL, 2.00, 2.00, 3, 'none', 30, 'fa-vial');

-- ICE WORLD ITEMS (Frozen Expanse)
INSERT INTO loot_items (name, description, rarity, world_type, item_type, equipment_slot, speed_bonus, luck_bonus, power_bonus, elemental_affinity, power_value, icon) VALUES
    -- Ice Weapons
    ('Frost Dagger', 'Leaves a trail of ice crystals.', 'common', 'ice', 'equipment', 'weapon', 3.00, 1.00, 4, 'none', 10, 'fa-dagger'),
    ('Glacial Spear', 'Frozen solid but surprisingly flexible.', 'uncommon', 'ice', 'equipment', 'weapon', 4.00, 2.00, 8, 'none', 20, 'fa-spear'),
    ('Blizzard Axe', 'Summons ice storms with each swing.', 'rare', 'ice', 'equipment', 'weapon', 5.00, 4.00, 16, 'none', 40, 'fa-axe'),
    ('Winter Claymore', 'A massive blade of eternal ice.', 'epic', 'ice', 'equipment', 'weapon', 7.00, 6.00, 26, 'none', 80, 'fa-sword'),
    -- Ice Armor
    ('Snowdrift Robe', 'Light as snow, cold as ice.', 'common', 'ice', 'equipment', 'armor', 1.00, 0.00, 6, 'none', 10, 'fa-shirt'),
    ('Permafrost Armor', 'Never melts, no matter the heat.', 'uncommon', 'ice', 'equipment', 'armor', 2.00, 2.00, 10, 'none', 20, 'fa-shield'),
    ('Avalanche Guard', 'As unmovable as a mountain of snow.', 'rare', 'ice', 'equipment', 'armor', 3.00, 4.00, 18, 'none', 40, 'fa-shield-halved'),
    -- Ice Accessories
    ('Frozen Tear', 'A pendant of pure ice that never melts.', 'common', 'ice', 'equipment', 'accessory', 4.00, 2.00, 2, 'none', 12, 'fa-droplet'),
    ('Crystal Snowflake', 'Each one truly unique.', 'uncommon', 'ice', 'equipment', 'accessory', 5.00, 4.00, 5, 'none', 22, 'fa-snowflake'),
    ('Winterstorm Crown', 'Worn by ancient ice kings.', 'rare', 'ice', 'equipment', 'accessory', 7.00, 6.00, 11, 'none', 42, 'fa-crown'),
    -- Ice Artifacts
    ('Glacier Core', 'The frozen heart of an ancient glacier.', 'epic', 'ice', 'equipment', 'artifact', 12.00, 9.00, 28, 'none', 90, 'fa-cube-ice'),
    -- Ice Relics
    ('Glacial Heart', 'Water incarnate, frozen in time.', 'legendary', 'ice', 'equipment', 'relic', 18.00, 13.00, 48, 'water', 150, 'fa-snowflake'),
    -- Ice Consumables
    ('Frost Potion', 'Chills you to the bone, but makes you stronger.', 'uncommon', 'ice', 'consumable', NULL, 3.00, 0.00, 2, 'none', 15, 'fa-flask'),
    ('Winter Essence', 'The concentrated cold of eternal winter.', 'rare', 'ice', 'consumable', NULL, 0.00, 3.00, 4, 'none', 30, 'fa-vial');

-- TECH WORLD ITEMS (Neon Sprawl)
INSERT INTO loot_items (name, description, rarity, world_type, item_type, equipment_slot, speed_bonus, luck_bonus, power_bonus, elemental_affinity, power_value, icon) VALUES
    -- Tech Weapons
    ('Plasma Pistol', 'Standard issue energy weapon.', 'common', 'tech', 'equipment', 'weapon', 4.00, 0.00, 4, 'none', 10, 'fa-gun'),
    ('Laser Rifle', 'Precision energy weapon with high output.', 'uncommon', 'tech', 'equipment', 'weapon', 5.00, 2.00, 7, 'none', 20, 'fa-crosshairs'),
    ('Quantum Blade', 'Exists in multiple states simultaneously.', 'rare', 'tech', 'equipment', 'weapon', 6.00, 4.00, 14, 'none', 40, 'fa-sword'),
    ('Singularity Cannon', 'Weaponized gravity well generator.', 'epic', 'tech', 'equipment', 'weapon', 8.00, 6.00, 24, 'none', 80, 'fa-rocket'),
    -- Tech Armor
    ('Flex Weave', 'Lightweight synthetic armor.', 'common', 'tech', 'equipment', 'armor', 2.00, 1.00, 5, 'none', 10, 'fa-shirt'),
    ('Nano Suit', 'Self-repairing armor with millions of nanobots.', 'uncommon', 'tech', 'equipment', 'armor', 3.00, 2.00, 9, 'none', 20, 'fa-shield'),
    ('Exo Frame', 'Powered armor that enhances strength.', 'rare', 'tech', 'equipment', 'armor', 4.00, 4.00, 17, 'none', 40, 'fa-robot'),
    -- Tech Accessories
    ('HUD Visor', 'Heads-up display with tactical information.', 'common', 'tech', 'equipment', 'accessory', 5.00, 3.00, 2, 'none', 12, 'fa-glasses'),
    ('Neural Link', 'Direct brain-computer interface.', 'uncommon', 'tech', 'equipment', 'accessory', 6.00, 5.00, 4, 'none', 22, 'fa-microchip'),
    ('Quantum Processor', 'Computes all possible outcomes simultaneously.', 'rare', 'tech', 'equipment', 'accessory', 8.00, 7.00, 10, 'none', 42, 'fa-cpu'),
    -- Tech Artifacts
    ('Fusion Reactor', 'Portable power source of immense energy.', 'epic', 'tech', 'equipment', 'artifact', 14.00, 10.00, 27, 'none', 90, 'fa-atom'),
    -- Tech Relics
    ('Storm Circuit', 'Living electricity trapped in crystalline circuits.', 'legendary', 'tech', 'equipment', 'relic', 20.00, 14.00, 46, 'wind', 150, 'fa-bolt'),
    -- Tech Consumables
    ('Nano Serum', 'Rewrites your cellular structure permanently.', 'uncommon', 'tech', 'consumable', NULL, 4.00, 1.00, 1, 'none', 15, 'fa-syringe'),
    ('Tech Essence', 'Pure computational power made physical.', 'rare', 'tech', 'consumable', NULL, 5.00, 0.00, 3, 'none', 30, 'fa-microchip');

-- NATURE WORLD ITEMS (Verdant Overgrowth)
INSERT INTO loot_items (name, description, rarity, world_type, item_type, equipment_slot, speed_bonus, luck_bonus, power_bonus, elemental_affinity, power_value, icon) VALUES
    -- Nature Weapons
    ('Thorn Whip', 'Barbed vines that ensnare enemies.', 'uncommon', 'nature', 'equipment', 'weapon', 4.00, 3.00, 9, 'none', 20, 'fa-whip'),
    ('Ancient Bow', 'Carved from wood older than civilization.', 'rare', 'nature', 'equipment', 'weapon', 6.00, 5.00, 17, 'none', 40, 'fa-bow-arrow'),
    ('Treant Greatclub', 'Shaped from the arm of a living tree.', 'epic', 'nature', 'equipment', 'weapon', 8.00, 7.00, 28, 'none', 80, 'fa-staff'),
    ('Worldroot Staff', 'Connected to the root network of the entire forest.', 'legendary', 'nature', 'equipment', 'weapon', 10.00, 10.00, 40, 'none', 120, 'fa-wand-sparkles'),
    -- Nature Armor
    ('Bark Plate', 'Natural armor as strong as steel.', 'uncommon', 'nature', 'equipment', 'armor', 2.00, 3.00, 11, 'none', 20, 'fa-shield'),
    ('Vine Mail', 'Living armor that regenerates.', 'rare', 'nature', 'equipment', 'armor', 3.00, 5.00, 19, 'none', 40, 'fa-leaf'),
    ('Grove Guardian', 'Blessed by ancient forest spirits.', 'epic', 'nature', 'equipment', 'armor', 5.00, 7.00, 30, 'none', 80, 'fa-tree'),
    -- Nature Accessories
    ('Acorn Charm', 'From the first tree, holds great potential.', 'uncommon', 'nature', 'equipment', 'accessory', 6.00, 6.00, 5, 'none', 22, 'fa-seedling'),
    ('Moonflower Pendant', 'Blooms only in moonlight.', 'rare', 'nature', 'equipment', 'accessory', 8.00, 8.00, 11, 'none', 42, 'fa-flower'),
    ('Forest Crown', 'Woven from living branches that never die.', 'epic', 'nature', 'equipment', 'accessory', 10.00, 10.00, 18, 'none', 70, 'fa-crown'),
    -- Nature Artifacts
    ('Life Seed', 'Contains the potential for infinite growth.', 'legendary', 'nature', 'equipment', 'artifact', 15.00, 12.00, 35, 'none', 100, 'fa-spa'),
    -- Nature Relics
    ('Earthheart Stone', 'The beating heart of the planet itself.', 'legendary', 'nature', 'equipment', 'relic', 22.00, 15.00, 52, 'earth', 150, 'fa-mountain'),
    -- Nature Consumables
    ('Growth Tonic', 'Accelerates natural development.', 'rare', 'nature', 'consumable', NULL, 3.00, 3.00, 4, 'none', 30, 'fa-flask-vial'),
    ('Nature Essence', 'The concentrated life force of the forest.', 'epic', 'nature', 'consumable', NULL, 2.00, 4.00, 8, 'none', 50, 'fa-vial');

-- VOID WORLD ITEMS (Void Confluence & Seat of Heaven)
INSERT INTO loot_items (name, description, rarity, world_type, item_type, equipment_slot, speed_bonus, luck_bonus, power_bonus, elemental_affinity, power_value, icon) VALUES
    -- Void Weapons
    ('Shadow Blade', 'Forged from solidified darkness.', 'rare', 'void', 'equipment', 'weapon', 7.00, 6.00, 20, 'none', 40, 'fa-knife'),
    ('Void Reaper', 'Harvests the essence of reality itself.', 'epic', 'void', 'equipment', 'weapon', 9.00, 8.00, 32, 'none', 80, 'fa-scythe'),
    ('Oblivion Edge', 'A sword that cuts through space and time.', 'legendary', 'void', 'equipment', 'weapon', 12.00, 12.00, 45, 'none', 120, 'fa-sword'),
    ('Radiant Lance', 'Pure crystallized light given form.', 'epic', 'void', 'equipment', 'weapon', 9.00, 8.00, 32, 'none', 80, 'fa-staff'),
    -- Void Armor
    ('Twilight Shroud', 'Exists between light and shadow.', 'rare', 'void', 'equipment', 'armor', 4.00, 6.00, 22, 'none', 40, 'fa-shirt'),
    ('Stellar Plate', 'Forged from collapsed starlight.', 'epic', 'void', 'equipment', 'armor', 6.00, 8.00, 34, 'none', 80, 'fa-shield'),
    ('Reality Weave', 'Armor that bends physics itself.', 'legendary', 'void', 'equipment', 'armor', 8.00, 10.00, 48, 'none', 120, 'fa-shield-halved'),
    -- Void Accessories
    ('Cosmos Ring', 'Contains a miniature universe.', 'rare', 'void', 'equipment', 'accessory', 9.00, 9.00, 14, 'none', 42, 'fa-ring'),
    ('Paradox Amulet', 'Exists and does not exist simultaneously.', 'epic', 'void', 'equipment', 'accessory', 11.00, 11.00, 22, 'none', 70, 'fa-infinity'),
    ('Dimensional Prism', 'Refracts reality into infinite possibilities.', 'legendary', 'void', 'equipment', 'accessory', 14.00, 14.00, 30, 'none', 100, 'fa-gem'),
    -- Void Artifacts
    ('Entropy Orb', 'The end of all things, contained.', 'legendary', 'void', 'equipment', 'artifact', 18.00, 15.00, 40, 'none', 120, 'fa-circle'),
    ('Genesis Sphere', 'The beginning of everything, crystallized.', 'legendary', 'void', 'equipment', 'artifact', 18.00, 15.00, 40, 'none', 120, 'fa-sun'),
    -- Void Relics
    ('Void Infinity Stone', 'The absence of all, the presence of nothing.', 'legendary', 'void', 'equipment', 'relic', 25.00, 18.00, 60, 'void', 200, 'fa-circle-notch'),
    ('Light Infinity Stone', 'The sum of all creation, infinite radiance.', 'legendary', 'void', 'equipment', 'relic', 25.00, 18.00, 60, 'light', 200, 'fa-star'),
    -- Void Consumables
    ('Void Essence', 'The taste of nothingness.', 'epic', 'void', 'consumable', NULL, 4.00, 4.00, 8, 'none', 50, 'fa-flask'),
    ('Light Essence', 'Bottled starlight and hope.', 'epic', 'void', 'consumable', NULL, 4.00, 4.00, 8, 'none', 50, 'fa-sun');

-- ############################
-- STEP 5: CREATE LOOT_DROP_TABLES TABLE
-- ############################

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

-- Seed drop rates for all rift/rarity combinations
-- Tutorial Rift (ID 1)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (1, 'common', 100.00, 1, 2);

-- Crimson Wastes - Easy (ID 2)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (2, 'common', 70.00, 2, 4),
    (2, 'uncommon', 25.00, 2, 4),
    (2, 'rare', 5.00, 1, 1);

-- Frozen Expanse - Medium (ID 3)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (3, 'common', 50.00, 3, 5),
    (3, 'uncommon', 30.00, 3, 5),
    (3, 'rare', 15.00, 3, 5),
    (3, 'epic', 5.00, 1, 1);

-- Neon Sprawl - Medium (ID 4)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (4, 'common', 50.00, 3, 5),
    (4, 'uncommon', 30.00, 3, 5),
    (4, 'rare', 15.00, 3, 5),
    (4, 'epic', 5.00, 1, 1);

-- Verdant Overgrowth - Hard (ID 5)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (5, 'common', 30.00, 4, 6),
    (5, 'uncommon', 35.00, 4, 6),
    (5, 'rare', 25.00, 4, 6),
    (5, 'epic', 9.00, 4, 6),
    (5, 'legendary', 1.00, 1, 1);

-- Void Confluence - Legendary (ID 6)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (6, 'common', 10.00, 5, 8),
    (6, 'uncommon', 25.00, 5, 8),
    (6, 'rare', 35.00, 5, 8),
    (6, 'epic', 25.00, 5, 8),
    (6, 'legendary', 5.00, 1, 1);

-- Seat of Heaven - Legendary (ID 7)
INSERT INTO loot_drop_tables (rift_id, rarity, drop_rate_percent, min_quantity, max_quantity) VALUES
    (7, 'common', 10.00, 5, 8),
    (7, 'uncommon', 25.00, 5, 8),
    (7, 'rare', 35.00, 5, 8),
    (7, 'epic', 25.00, 5, 8),
    (7, 'legendary', 5.00, 1, 1);

-- ############################
-- STEP 6: CREATE TEAMS TABLE
-- ############################

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
    equipped_weapon_slot INT,
    equipped_armor_slot INT,
    equipped_accessory_slot INT,
    equipped_artifact_slot INT,
    equipped_relic_slot INT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT false,

    UNIQUE(user_id, team_number)
);

CREATE INDEX idx_teams_user_id ON teams(user_id);
CREATE INDEX idx_teams_unlocked ON teams(user_id, is_unlocked);

-- ############################
-- STEP 7: CREATE USER_INVENTORY TABLE
-- ############################

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

-- Now add the foreign key constraints for teams equipment slots
ALTER TABLE teams
    ADD CONSTRAINT fk_teams_weapon FOREIGN KEY (equipped_weapon_slot) REFERENCES user_inventory(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_teams_armor FOREIGN KEY (equipped_armor_slot) REFERENCES user_inventory(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_teams_accessory FOREIGN KEY (equipped_accessory_slot) REFERENCES user_inventory(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_teams_artifact FOREIGN KEY (equipped_artifact_slot) REFERENCES user_inventory(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_teams_relic FOREIGN KEY (equipped_relic_slot) REFERENCES user_inventory(id) ON DELETE SET NULL;

-- ############################
-- STEP 8: CREATE EXPEDITIONS TABLE
-- ############################

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

-- ############################
-- STEP 9: CREATE EXPEDITION_LOOT TABLE
-- ############################

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

-- ############################
-- STEP 10: VERIFY ALL INDEXES
-- ############################

-- All required indexes were created inline with their tables in steps 3-9:
-- - rifts: idx_rifts_difficulty, idx_rifts_world_type
-- - loot_items: idx_loot_items_rarity, idx_loot_items_world_type, idx_loot_items_item_type, idx_loot_items_slot
-- - loot_drop_tables: idx_drop_tables_rift
-- - teams: idx_teams_user_id, idx_teams_unlocked
-- - user_inventory: idx_inventory_user, idx_inventory_item, idx_inventory_user_item
-- - expeditions: idx_expeditions_user, idx_expeditions_team, idx_expeditions_status, idx_expeditions_completion_check
-- - expedition_loot: idx_expedition_loot_expedition, idx_expedition_loot_item

-- ############################
-- STEP 11: SEED USER TEAMS
-- ############################

-- Create Team 1 (unlocked by default) for all existing users
-- Teams 2-5 are created but locked (unlocked through gameplay)
INSERT INTO teams (user_id, team_number, is_unlocked)
SELECT id, 1, true
FROM users
WHERE NOT EXISTS (
    SELECT 1 FROM teams WHERE teams.user_id = users.id AND teams.team_number = 1
);

-- Create Teams 2-5 (locked) for all existing users
INSERT INTO teams (user_id, team_number, is_unlocked)
SELECT id, team_num, false
FROM users
CROSS JOIN (SELECT 2 AS team_num UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5) AS team_numbers
WHERE NOT EXISTS (
    SELECT 1 FROM teams WHERE teams.user_id = users.id AND teams.team_number = team_num
);

-- Note: New users will get their teams created by the UserService.CreateUser() method
-- This seed data only handles existing users who were created before the game schema was added
