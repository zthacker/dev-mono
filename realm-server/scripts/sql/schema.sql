-- =============================================================================
-- REALM SERVER DATABASE SCHEMA
-- =============================================================================
--
-- This schema follows patterns from WoW private server emulators.
-- Designed for MySQL 8.0+
--
-- Tables are organized into:
-- 1. Account/Auth tables (might be separate DB in production)
-- 2. Character tables (player data)
-- 3. World tables (templates, spawns - often read-only)

-- =============================================================================
-- ACCOUNTS (could be separate auth database)
-- =============================================================================

CREATE TABLE IF NOT EXISTS accounts (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(32) NOT NULL,
    -- Password hash (use bcrypt or argon2 in production, NOT plain SHA)
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255),

    -- Permissions
    gm_level TINYINT UNSIGNED NOT NULL DEFAULT 0,

    -- Status
    banned TINYINT UNSIGNED NOT NULL DEFAULT 0,
    ban_reason VARCHAR(255),
    ban_expires DATETIME,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY idx_username (username),
    KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =============================================================================
-- CHARACTERS
-- =============================================================================

CREATE TABLE IF NOT EXISTS characters (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    account_id INT UNSIGNED NOT NULL,
    name VARCHAR(12) NOT NULL,

    -- Appearance
    race TINYINT UNSIGNED NOT NULL,
    class TINYINT UNSIGNED NOT NULL,
    gender TINYINT UNSIGNED NOT NULL,
    skin TINYINT UNSIGNED NOT NULL DEFAULT 0,
    face TINYINT UNSIGNED NOT NULL DEFAULT 0,
    hair_style TINYINT UNSIGNED NOT NULL DEFAULT 0,
    hair_color TINYINT UNSIGNED NOT NULL DEFAULT 0,

    -- Position
    map_id INT UNSIGNED NOT NULL DEFAULT 0,
    zone_id INT UNSIGNED NOT NULL DEFAULT 0,
    position_x FLOAT NOT NULL DEFAULT 0,
    position_y FLOAT NOT NULL DEFAULT 0,
    position_z FLOAT NOT NULL DEFAULT 0,
    orientation FLOAT NOT NULL DEFAULT 0,

    -- Progression
    level TINYINT UNSIGNED NOT NULL DEFAULT 1,
    xp INT UNSIGNED NOT NULL DEFAULT 0,
    money INT UNSIGNED NOT NULL DEFAULT 0,

    -- Stats (current values - max values derived from level/gear)
    health INT NOT NULL DEFAULT 100,
    max_health INT NOT NULL DEFAULT 100,
    mana INT NOT NULL DEFAULT 0,
    max_mana INT NOT NULL DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    total_time INT UNSIGNED NOT NULL DEFAULT 0,

    -- Flags
    at_login INT UNSIGNED NOT NULL DEFAULT 0,
    is_deleted TINYINT UNSIGNED NOT NULL DEFAULT 0,

    PRIMARY KEY (id),
    UNIQUE KEY idx_name (name),
    KEY idx_account (account_id),
    KEY idx_zone (zone_id),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =============================================================================
-- CHARACTER INVENTORY
-- =============================================================================

CREATE TABLE IF NOT EXISTS character_inventory (
    character_id BIGINT UNSIGNED NOT NULL,
    bag TINYINT UNSIGNED NOT NULL,      -- 0 = equipped, 1 = backpack, 2-5 = bags
    slot TINYINT UNSIGNED NOT NULL,
    item_id INT UNSIGNED NOT NULL,       -- Reference to item_template
    stack_count SMALLINT UNSIGNED NOT NULL DEFAULT 1,
    durability INT UNSIGNED NOT NULL DEFAULT 0,

    -- Item modifications
    enchant_id INT UNSIGNED DEFAULT 0,
    -- gems, reforging, transmogrification, etc. would be additional columns

    PRIMARY KEY (character_id, bag, slot),
    KEY idx_item (item_id),
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =============================================================================
-- CHARACTER SPELLS
-- =============================================================================

CREATE TABLE IF NOT EXISTS character_spells (
    character_id BIGINT UNSIGNED NOT NULL,
    spell_id INT UNSIGNED NOT NULL,
    active TINYINT UNSIGNED NOT NULL DEFAULT 1,

    PRIMARY KEY (character_id, spell_id),
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =============================================================================
-- CHARACTER QUESTS
-- =============================================================================

CREATE TABLE IF NOT EXISTS character_quests (
    character_id BIGINT UNSIGNED NOT NULL,
    quest_id INT UNSIGNED NOT NULL,
    status TINYINT UNSIGNED NOT NULL DEFAULT 0, -- 0=in progress, 1=complete, 2=failed

    -- Objective tracking (could be JSON or separate table for complex quests)
    objective_1 INT UNSIGNED NOT NULL DEFAULT 0,
    objective_2 INT UNSIGNED NOT NULL DEFAULT 0,
    objective_3 INT UNSIGNED NOT NULL DEFAULT 0,
    objective_4 INT UNSIGNED NOT NULL DEFAULT 0,

    accepted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (character_id, quest_id),
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =============================================================================
-- WORLD TEMPLATES (read-only reference data)
-- =============================================================================

-- Item definitions
CREATE TABLE IF NOT EXISTS item_template (
    id INT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    class TINYINT UNSIGNED NOT NULL,     -- Weapon, Armor, Consumable, etc.
    subclass TINYINT UNSIGNED NOT NULL,  -- Sword, Axe, Cloth, Plate, etc.
    quality TINYINT UNSIGNED NOT NULL DEFAULT 1,
    level TINYINT UNSIGNED NOT NULL DEFAULT 1,
    required_level TINYINT UNSIGNED NOT NULL DEFAULT 0,

    -- Stats
    armor INT UNSIGNED DEFAULT 0,
    damage_min INT UNSIGNED DEFAULT 0,
    damage_max INT UNSIGNED DEFAULT 0,
    attack_speed INT UNSIGNED DEFAULT 2000, -- milliseconds

    -- TODO: Add stat bonuses, on-use effects, etc.

    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- NPC definitions
CREATE TABLE IF NOT EXISTS npc_template (
    id INT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    subname VARCHAR(100),
    level TINYINT UNSIGNED NOT NULL DEFAULT 1,
    health INT NOT NULL DEFAULT 100,
    mana INT NOT NULL DEFAULT 0,
    faction INT UNSIGNED NOT NULL DEFAULT 0,
    npc_flags INT UNSIGNED NOT NULL DEFAULT 0,

    -- TODO: Add AI script, loot table reference, etc.

    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Spawn points
CREATE TABLE IF NOT EXISTS npc_spawns (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    template_id INT UNSIGNED NOT NULL,
    map_id INT UNSIGNED NOT NULL,
    zone_id INT UNSIGNED NOT NULL,
    position_x FLOAT NOT NULL,
    position_y FLOAT NOT NULL,
    position_z FLOAT NOT NULL,
    orientation FLOAT NOT NULL DEFAULT 0,
    respawn_time INT UNSIGNED NOT NULL DEFAULT 300, -- seconds
    wander_radius FLOAT NOT NULL DEFAULT 0,

    PRIMARY KEY (id),
    KEY idx_zone (zone_id),
    FOREIGN KEY (template_id) REFERENCES npc_template(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =============================================================================
-- INDEXES FOR COMMON QUERIES
-- =============================================================================

-- TODO: Add indexes based on actual query patterns:
-- - Character lookup by name (for /who, mail, etc.)
-- - Characters by zone (for zone population)
-- - Item searches (auction house)
