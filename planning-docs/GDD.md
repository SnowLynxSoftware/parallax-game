# Parallax - Game Design Document

**Genre:** Idle Progression / Resource Management  
**Platform:** Browser-Based (PBBG)  
**Theme:** Parallel Worlds  
**Development Time:** 3 weeks (solo)  
**Tech Stack:** Go + Chi + Postgres + HTML Templates

---

## 1. High-Level Concept

**Elevator Pitch:**  
Manage expedition teams that explore dimensional rifts leading to parallel worlds. Send teams on timed expeditions, collect randomized loot from different dimensions, upgrade your teams, and compete on leaderboards for the rarest cross-dimensional artifacts.

**Core Loop:**  
Send Team → Wait (5min - 1hr) → Receive Random Loot → Upgrade Team → Send to Harder Rift → Repeat

**Theme Integration:**  
Each rift leads to a parallel version of Earth with different physical laws and resources. Players collect dimension-specific loot (Fire World crystals, Ice World shards, Tech World circuits, etc.). The theme is primarily cosmetic but reinforced through consistent world-building.

---

## 2. Core Mechanics

### 2.1 Expedition System

**How It Works:**

1. Player selects a team (unlocked progressively: 1st at start, 2nd after first rift, 3rd after second rift)
2. Player selects a rift destination (5 parallel worlds available)
3. Expedition starts with timestamp recorded in DB
4. Cron job runs every 5 minutes, processes completed expeditions
5. Player returns to see loot rewards and can send team again

**Expedition Durations:**

- **Tutorial Rift:** 5 minutes (guaranteed common loot)
- **Easy Rifts:** 15 minutes (Common/Uncommon)
- **Medium Rifts:** 30 minutes (Uncommon/Rare)
- **Hard Rifts:** 1 hour (Rare/Epic)
- **Legendary Rifts:** 2 hours (Epic/Legendary)

**Database Schema:**

```
expeditions:
  id, user_id, team_id, rift_id, start_time, duration_minutes,
  completed (bool), processed (bool), created_at
```

**Cron Job Logic (runs every 5 min):**

```
1. SELECT * FROM expeditions
   WHERE completed = false
   AND start_time + duration_minutes <= NOW()

2. For each expedition:
   - Generate loot based on rift loot table
   - Insert into user_inventory
   - Mark expedition as completed = true, processed = true
   - Update team status to available

3. (Week 3) Check for echo encounters (10% chance)
```

### 2.2 Teams

**Team Attributes:**

- **Speed Modifier:** Reduces expedition duration (e.g., +10% speed = 54min instead of 60min)
- **Luck Modifier:** Increases rare drop rates (e.g., +5% luck = 6% epic drop instead of 5%)
- **Specialization:** Tank, Scout, or Scientist (affects success rate in specific rifts)

**Progression:**

- **Team 1:** Unlocked at registration (starts with pre-completed tutorial expedition)
- **Team 2:** Unlocked after completing first expedition
- **Team 3:** Unlocked after completing third expedition
- **Team 4 & 5:** (Week 2 stretch) Unlock after 25 and 50 total expeditions

**Upgrades:**
Teams upgrade by consuming loot:

- **Speed Upgrade:** Costs 10 Uncommon + 5 Rare → +5% speed (max 50%)
- **Luck Upgrade:** Costs 15 Uncommon + 3 Rare → +3% luck (max 30%)
- **Specialization:** Costs 1 Epic → Choose Tank/Scout/Scientist (permanent)

### 2.3 Rifts (Parallel Worlds)

**Available Rifts:**

1. **Tutorial Rift** (5 min)

   - Bland, safe dimension
   - 100% Common loot drops
   - Used for onboarding only

2. **Crimson Wastes** (Fire World - 15 min)

   - Volcanic parallel Earth
   - Drops: Ember Shards, Magma Cores, Flame Crystals
   - Difficulty: Easy

3. **Frozen Expanse** (Ice World - 30 min)

   - Eternal winter dimension
   - Drops: Frost Fragments, Glacial Hearts, Permafrost Essence
   - Difficulty: Medium

4. **Neon Sprawl** (Tech World - 30 min)

   - Cyberpunk parallel reality
   - Drops: Circuit Boards, Quantum Chips, Data Crystals
   - Difficulty: Medium

5. **Verdant Overgrowth** (Nature World - 1 hour)

   - Hyper-evolved biosphere
   - Drops: Bio-Catalysts, Spore Clusters, Ancient Seeds
   - Difficulty: Hard

6. **Void Confluence** (Void World - 2 hours)
   - Reality breakdown dimension
   - Drops: Void Essence, Entropy Shards, Paradox Fragments
   - Difficulty: Legendary

**Rift Requirements:**

- Tutorial: Always available
- Crimson/Frozen/Neon: Available after completing Tutorial
- Verdant: Unlocked after 10 total expeditions
- Void: Unlocked after 25 total expeditions + Must have at least 1 Epic item

### 2.4 Loot System

**Rarity Tiers:**

- Common (White)
- Uncommon (Green)
- Rare (Blue)
- Epic (Purple)
- Legendary (Orange)

**Drop Rate Tables:**

**Tutorial Rift:**

- 100% Common

**Easy Rifts (Crimson Wastes):**

- 70% Common
- 25% Uncommon
- 5% Rare

**Medium Rifts (Frozen Expanse, Neon Sprawl):**

- 50% Common
- 30% Uncommon
- 15% Rare
- 5% Epic

**Hard Rifts (Verdant Overgrowth):**

- 30% Common
- 35% Uncommon
- 25% Rare
- 9% Epic
- 1% Legendary

**Legendary Rifts (Void Confluence):**

- 10% Common
- 25% Uncommon
- 35% Rare
- 25% Epic
- 5% Legendary

**Loot Quantity Per Expedition:**

- Tutorial: 3 items
- Easy: 2-4 items
- Medium: 3-5 items
- Hard: 4-6 items
- Legendary: 5-8 items

**Luck Modifier Effect:**
Each point of Luck shifts drop rates toward higher rarities:

- +10% Luck = shift 10% from each tier to the tier above
- Example: Medium rift with +10% luck becomes:
  - 45% Common (was 50%)
  - 28% Uncommon (was 30%)
  - 17% Rare (was 15%)
  - 9% Epic (was 5%)
  - 1% Legendary (was 0.5%, rounded up)

**Database Schema:**

```
loot_items:
  id, name, description, rarity, world_type, icon (font-awesome class)

user_inventory:
  id, user_id, loot_item_id, quantity, acquired_at
```

### 2.5 Progression Systems

**Player Goals:**

1. **Short-term:** Complete expeditions, collect loot, unlock teams
2. **Mid-term:** Upgrade teams, unlock harder rifts, get Epic items
3. **Long-term:** Collect Legendary items, top leaderboard, optimize team setups

**Unlock Progression:**

```
Start → Tutorial Rift → Team 2 → Easy Rifts → Team 3 → Medium Rifts
→ Hard Rifts (10 expeditions) → Legendary Rifts (25 expeditions + 1 Epic)
```

**Power Curve:**

- Hour 1: Tutorial + first few Easy rifts (Teams 1-3 unlocked)
- Hour 2-3: Easy/Medium rifts, first Rare items
- Hour 4-8: Medium/Hard rifts, team upgrades, first Epic
- Hour 10+: Legendary rifts, Legendary items, leaderboard competition

---

## 3. Multiplayer Elements

### 3.1 Leaderboard (Priority 1 - Must Ship)

**Leaderboard Categories:**

1. **Total Legendary Items** - Who has collected the most Legendary loot
2. **Dimensional Power Score** - Weighted value of all inventory (Common=1, Uncommon=5, Rare=25, Epic=125, Legendary=1000)
3. **Total Expeditions Completed** - Volume metric

**Display:**

- Top 100 players
- Shows: Rank, Username, Score, Last Active
- Player's own rank always visible (even if outside top 100)

**Database Query:**

```sql
-- Legendary Count
SELECT user_id, username, COUNT(*) as legendary_count
FROM user_inventory
JOIN users ON users.id = user_inventory.user_id
JOIN loot_items ON loot_items.id = user_inventory.loot_item_id
WHERE loot_items.rarity = 'legendary'
GROUP BY user_id, username
ORDER BY legendary_count DESC
LIMIT 100
```

### 3.2 Echo Encounters (Priority 3 - Week 3 Only)

**Mechanic:**
When an expedition completes, 10% chance to encounter "echoes" of another player's recent expedition to the same rift.

**Effect:**

- Both players receive +25% loot quantity
- Notification: "Your team encountered echoes of [PlayerName]'s expedition!"
- Tracked stat: "Echo Encounters" (visible on profile)

**Implementation:**

```
On expedition completion:
1. Roll 10% chance
2. If success, query for another expedition to same rift within last 24h
3. If found, award bonus loot to both players
4. Insert notification records
```

**Database:**

```
echo_encounters:
  id, user_id_1, user_id_2, rift_id, occurred_at

notifications:
  id, user_id, message, read (bool), created_at
```

---

## 4. User Interface

### 4.1 Page Structure

**Pages Required:**

1. **Landing/Login** - Registration + Login forms
2. **Dashboard** - Main hub (active expeditions, quick stats)
3. **Expeditions** - Send teams to rifts
4. **Teams** - View team stats, upgrade teams
5. **Inventory** - View collected loot
6. **Leaderboard** - Rankings
7. **Profile** - Player stats, echo encounters

### 4.2 Dashboard Layout

**Top Bar:**

- Username
- Total Expeditions Completed
- Dimensional Power Score
- Legendary Count

**Main Content:**

- **Active Expeditions Panel** (real-time countdown timers via JS)

  ```
  [Team 1] → Crimson Wastes
  [Progress Bar] 8 minutes remaining
  [Team 2] → Available - Send Expedition
  [Team 3] → Locked (Complete 1 more expedition)
  ```

- **Recent Loot Panel** (last 10 items acquired)

  ```
  [Epic Icon] Paradox Fragment (x1)
  [Rare Icon] Magma Core (x3)
  [Common Icon] Ember Shard (x5)
  ```

- **Quick Actions:**
  - [Send Expedition] button (opens modal)
  - [View Leaderboard]
  - [Upgrade Teams]

### 4.3 Expedition Flow

**Modal/Page for Sending Expedition:**

1. Select available team (dropdown or buttons)
2. Display team stats (Speed: +10%, Luck: +5%, Type: Scout)
3. Select rift (show duration, difficulty, rewards preview)
4. Confirm button
5. Success message: "Team 1 departed for Crimson Wastes! Return in 15 minutes."

**Completion Flow:**
Player returns to Dashboard → Sees notification badge → Clicks "Claim Rewards" → Modal shows:

```
Expedition Complete!
Team 1 returned from Crimson Wastes

Loot Acquired:
- [Uncommon] Ember Shard x4
- [Rare] Magma Core x1

[Claim] button → Adds to inventory, clears expedition
```

### 4.4 Visual Design (Minimal Viable)

**Style:**

- Dark theme (easier on eyes, hides lack of art)
- Font Awesome icons for loot/rifts
- Color-coded rarity (CSS classes: .rarity-common, .rarity-legendary, etc.)
- Simple progress bars for expedition timers
- Bootstrap-like grid layout (or custom simple CSS)

**Art Assets:**

- **Rifts:** Font Awesome icons (fire, snowflake, microchip, leaf, circle-notch)
- **Loot:** Font Awesome gems/items colored by rarity
- **Teams:** Simple avatar icons (user, user-shield, user-astronaut for specializations)

**Week 3 Polish:**

- Generate AI art for each rift (background images)
- Better loot icons (Gemini-generated simple icons)
- Animated expedition "departure" effect

---

## 5. Technical Implementation

### 5.1 Database Schema

```sql
-- Users (from your existing template)
users:
  id, username, email, password_hash, created_at, last_login

-- Teams
teams:
  id, user_id, team_number (1-5), speed_bonus (decimal),
  luck_bonus (decimal), specialization (enum: none/tank/scout/scientist),
  is_unlocked (bool), created_at

-- Rifts (static data, could be hardcoded or DB)
rifts:
  id, name, description, world_type, duration_minutes,
  difficulty (enum), unlock_requirement_text, icon

-- Loot Items (static data)
loot_items:
  id, name, description, rarity, world_type, icon, power_value

-- User Inventory
user_inventory:
  id, user_id, loot_item_id, quantity, acquired_at

-- Expeditions
expeditions:
  id, user_id, team_id, rift_id, start_time, duration_minutes,
  completed (bool), processed (bool), claimed (bool), created_at

-- Loot Drops (records what was found)
expedition_loot:
  id, expedition_id, loot_item_id, quantity

-- Echo Encounters (Week 3)
echo_encounters:
  id, expedition_id_1, expedition_id_2, user_id_1, user_id_2,
  rift_id, occurred_at

-- Notifications (Week 3)
notifications:
  id, user_id, message, type, read (bool), created_at

-- Leaderboard Cache (optional optimization)
leaderboard_cache:
  user_id, legendary_count, power_score, total_expeditions,
  updated_at
```

### 5.2 Cron Job (Expedition Processor)

**Schedule:** Every 5 minutes

**Pseudocode:**

```go
func ProcessExpeditions() {
    // Find completed but unprocessed expeditions
    expeditions := db.Query(`
        SELECT id, user_id, team_id, rift_id, start_time, duration_minutes
        FROM expeditions
        WHERE completed = false
        AND processed = false
        AND start_time + INTERVAL duration_minutes MINUTE <= NOW()
    `)

    for _, exp := range expeditions {
        // Generate loot
        rift := GetRift(exp.rift_id)
        team := GetTeam(exp.team_id)
        loot := GenerateLoot(rift, team.luck_bonus)

        // Save loot to inventory
        for _, item := range loot {
            db.Exec(`
                INSERT INTO user_inventory (user_id, loot_item_id, quantity)
                VALUES (?, ?, ?)
                ON DUPLICATE KEY UPDATE quantity = quantity + ?
            `, exp.user_id, item.id, item.quantity, item.quantity)

            db.Exec(`
                INSERT INTO expedition_loot (expedition_id, loot_item_id, quantity)
                VALUES (?, ?, ?)
            `, exp.id, item.id, item.quantity)
        }

        // Mark expedition complete
        db.Exec(`
            UPDATE expeditions
            SET completed = true, processed = true
            WHERE id = ?
        `, exp.id)

        // Week 3: Check for echo encounters
        if rand.Float64() < 0.10 {
            CheckForEchoEncounter(exp)
        }

        // Update leaderboard cache
        UpdateLeaderboardCache(exp.user_id)
    }
}

func GenerateLoot(rift Rift, luckBonus float64) []LootDrop {
    lootTable := GetLootTable(rift.difficulty)
    quantity := RandomQuantity(rift.difficulty) // 2-4 for easy, 5-8 for legendary

    var drops []LootDrop
    for i := 0; i < quantity; i++ {
        rarity := RollRarity(lootTable, luckBonus)
        item := RandomItemFromWorld(rift.world_type, rarity)
        drops = append(drops, LootDrop{item: item, quantity: 1})
    }

    return drops
}

func RollRarity(table LootTable, luckBonus float64) string {
    // Apply luck bonus (shift probabilities up)
    adjusted := AdjustDropRates(table, luckBonus)

    roll := rand.Float64()
    cumulative := 0.0

    for _, tier := range []string{"legendary", "epic", "rare", "uncommon", "common"} {
        cumulative += adjusted[tier]
        if roll <= cumulative {
            return tier
        }
    }

    return "common" // fallback
}
```

### 5.3 Key API Endpoints

```
POST   /expedition/start        - Start new expedition
GET    /expedition/status       - Get all user's expeditions
POST   /expedition/claim/:id    - Claim completed expedition loot
GET    /teams                   - Get user's teams
POST   /teams/:id/upgrade       - Upgrade team (speed/luck/spec)
GET    /inventory               - Get user inventory
GET    /leaderboard/:type       - Get leaderboard (legendary/power/expeditions)
GET    /rifts                   - Get available rifts for user
GET    /dashboard               - Main dashboard data (expeditions, stats, recent loot)
```

### 5.4 Frontend JavaScript

**Expedition Timers (dashboard):**

```javascript
// Poll every 10 seconds for expedition updates
setInterval(async () => {
  const response = await fetch("/expedition/status");
  const expeditions = await response.json();

  expeditions.forEach((exp) => {
    const remaining = calculateRemaining(exp.start_time, exp.duration_minutes);
    updateProgressBar(exp.id, remaining);

    if (remaining <= 0 && !exp.claimed) {
      showClaimButton(exp.id);
    }
  });
}, 10000);
```

**No complex frontend framework needed** - vanilla JS or jQuery for simple interactivity.

---

## 6. Development Roadmap

### Week 1: Core Loop (Target: 40 hours)

**Day 1-2 (10h): Database + Expedition System**

- [ ] Create DB schema (teams, rifts, expeditions, loot_items, user_inventory)
- [ ] Seed rifts and loot_items (static data)
- [ ] Implement expedition start endpoint
- [ ] Basic expedition display on dashboard

**Day 3-4 (12h): Cron Job + Loot Generation**

- [ ] Build cron job (runs every 5 min)
- [ ] Implement RNG loot generation with drop tables
- [ ] Save loot to inventory
- [ ] Test: Send expedition → wait 5 min → verify loot appears

**Day 5-6 (10h): Teams + Upgrades**

- [ ] Team unlock system (1st at start, 2nd/3rd after expeditions)
- [ ] Team upgrade endpoints (speed, luck)
- [ ] Teams page UI
- [ ] Test: Upgrade team → verify faster expeditions

**Day 7 (8h): UI Polish + First Playthrough**

- [ ] Dashboard shows active expeditions with timers
- [ ] Inventory page displays loot with rarity colors
- [ ] Tutorial rift flow (pre-completed expedition for new users)
- [ ] **Milestone: Full loop playable**

### Week 2: Content + Progression (Target: 30 hours)

**Day 8-9 (8h): Expand Rifts**

- [ ] Add 3 more rifts (Verdant, Void, + 1 more if time)
- [ ] Implement unlock requirements (expedition count, items)
- [ ] Different loot tables per rift
- [ ] Test: Progression from Tutorial → Legendary rift

**Day 10-11 (8h): Loot Variety + Balancing**

- [ ] Create 10+ loot items per world (50+ total items)
- [ ] Implement luck modifier effect on drop rates
- [ ] Inventory sorting/filtering
- [ ] Test: Playthrough for 2-3 hours, tune drop rates

**Day 12-13 (8h): Team Specializations**

- [ ] Add specialization system (Tank/Scout/Scientist)
- [ ] Specializations affect success rate in rifts
- [ ] UI for choosing specialization
- [ ] Test: Verify specialization impact

**Day 14 (6h): Progression Tuning**

- [ ] Balance expedition durations (too fast/slow?)
- [ ] Balance upgrade costs
- [ ] Add more team upgrade tiers
- [ ] **Milestone: 3-4 hours of engaging content**

### Week 3: Polish + Multiplayer (Target: 10-15 hours)

**Day 15-16 (6h): Leaderboard**

- [ ] Implement leaderboard queries (legendary count, power score)
- [ ] Leaderboard page UI
- [ ] Show player's rank even if outside top 100
- [ ] Test: Verify correct rankings

**Day 17-18 (6h): UI/UX Polish**

- [ ] Add Font Awesome icons for all rifts/loot
- [ ] Improve expedition timer display (countdown, progress bars)
- [ ] Tutorial/onboarding text (explain first expedition)
- [ ] Color-coded rarity highlighting
- [ ] Responsive design (mobile-friendly)

**Day 19 (4h): Echo Encounters (IF TIME)**

- [ ] Implement 10% echo encounter check
- [ ] Notifications system
- [ ] Test: Verify echoes trigger and both players get bonus

**Day 20 (4h): AI Art + Final Polish (IF TIME)**

- [ ] Generate rift background images (Gemini)
- [ ] Generate loot icons
- [ ] Add achievement system (stretch)
- [ ] Bug fixes

**Day 21: Buffer/Testing**

- Final playtesting
- Balance adjustments
- Deployment preparation
- **Submission**

---

## 7. Success Criteria (Judging Alignment)

### Theme Adherence (15%)

**Target: 3-4/5**

- ✅ Parallel worlds are central to the game (rifts to different dimensions)
- ✅ Each world has distinct identity (Fire, Ice, Tech, Nature, Void)
- ⚠️ Theme is somewhat cosmetic (could be "dungeons" instead)
- **Strategy:** Strong world-building through descriptions, lore snippets

### Gameplay & Engagement (30%)

**Target: 4/5**

- ✅ Clear progression loop (expeditions → loot → upgrades → harder rifts)
- ✅ RNG creates excitement (rare drop dopamine)
- ✅ Multiple goals (unlock teams, collect legendaries, leaderboard)
- ✅ Idle mechanics allow check-ins (respects player time)
- **Strategy:** Tune drop rates for frequent rewards, balance grind

### Accessibility & Usability (30%)

**Target: 4-5/5**

- ✅ Simple, intuitive UI (send team, claim loot)
- ✅ No complex mechanics (no crafting, no confusing systems)
- ✅ Tutorial expedition teaches loop immediately
- ✅ Progress bars and timers set clear expectations
- **Strategy:** Focus on clarity over complexity, test with non-gamers

### Artistic Presentation (15%)

**Target: 3/5**

- ⚠️ Limited art assets (Font Awesome icons, minimal graphics)
- ✅ Clean, consistent UI design
- ✅ Color-coded rarity creates visual hierarchy
- **Strategy:** Week 3 AI art for rifts, polished dark theme

### Performance and Polish (10%)

**Target: 4/5**

- ✅ Go backend is fast and reliable
- ✅ Simple DB queries (minimal optimization needed)
- ✅ Cron job handles load easily
- ⚠️ Potential bugs in RNG, edge cases
- **Strategy:** Week 3 focus on testing, bug fixes

**Realistic Total: 3.6-4.0/5 average (72-80%)**

---

## 8. Risk Mitigation

### Risk 1: Loot System Feels Unrewarding

**Probability:** High  
**Impact:** Critical (kills engagement)

**Mitigation:**

- Playtest extensively in Week 2
- Ensure players get rare+ items within first hour
- Add "pity timer" (guaranteed epic after 20 expeditions without one)
- Transparent drop rates (show percentages in UI)

### Risk 2: Not Enough Content

**Probability:** Medium  
**Impact:** High (judges get bored)

**Mitigation:**

- Focus on 5 rifts minimum (achievable)
- 50+ loot items (use AI to generate names/descriptions quickly)
- Progression should last 4-6 hours for judges to see full game

### Risk 3: Timer System Bugs

**Probability:** Medium  
**Impact:** High (expeditions don't complete)

**Mitigation:**

- Test cron job thoroughly in Week 1
- Add logging to track expedition processing
- Manual "force complete" admin button for testing

### Risk 4: Multiplayer (Echoes) Doesn't Work

**Probability:** Low  
**Impact:** Low (it's optional Week 3 feature)

**Mitigation:**

- Build core game first
- Echoes are nice-to-have, not required
- Can ship without multiplayer and still score well

### Risk 5: Scope Creep

**Probability:** High  
**Impact:** Critical (nothing ships)

**Mitigation:**

- **Stick to the roadmap**
- No crafting system
- No additional features in Week 1-2
- Week 3 is polish only, not new systems

---

## 9. Post-Jam Potential (Optional)

If the game is successful and you want to continue:

**Possible Expansions:**

- Prestige system (reset for permanent bonuses)
- PvP expeditions (race other players to same rift)
- Guilds/Clans (shared rift discoveries)
- Seasonal events (limited-time rifts)
- Crafting system (combine loot for unique items)
- More rifts (expand to 15+ parallel worlds)

**Monetization (if desired):**

- Cosmetic team skins
- Additional team slots (4th, 5th team as IAP)
- "Expedite" consumables (speed up one expedition)

---

## 10. Final Notes

**What Makes This Viable:**

1. ✅ Proven idle game mechanics (Cookie Clicker, Torn inspiration)
2. ✅ Scoped to 3 weeks (no crafting, simple UI)
3. ✅ Tech stack you know (Go + Postgres)
4. ✅ Clear Week 1 milestone (playable loop)
5. ✅ Multiplayer is optional (leaderboard is easy, echoes are stretch)

**What Could Go Wrong:**

1. ⚠️ Drop rates feel bad (requires tuning)
2. ⚠️ Not enough content (need 50+ items)
3. ⚠️ Theme feels generic (need strong world-building)

**Your Job:**

- Build Week 1 core loop (40h budget)
- Playtest constantly (feel the progression)
- Cut features ruthlessly if behind schedule
- **Ship a working game, not a perfect game**

**This is buildable. Now build it.**

---

## Appendix A: Loot Item Examples

### Crimson Wastes (Fire World)

- Common: Ash Pile, Ember Shard, Charcoal Fragment
- Uncommon: Magma Core, Flame Crystal, Igneous Rock
- Rare: Inferno Heart, Volcanic Diamond, Phoenix Feather
- Epic: Eternal Flame, Magma Titan Core, Solar Rune
- Legendary: Crimson Infinity Stone, Primordial Fire Essence

### Frozen Expanse (Ice World)

- Common: Ice Chip, Frost Dust, Snow Crystal
- Uncommon: Glacial Shard, Frozen Heart, Permafrost Chunk
- Rare: Aurora Fragment, Blizzard Orb, Cryo-Crystal
- Epic: Eternal Ice, Glacial Titan Core, Lunar Rune
- Legendary: Frozen Infinity Stone, Absolute Zero Essence

### Neon Sprawl (Tech World)

- Common: Scrap Wire, Broken Circuit, Rusty Chip
- Uncommon: Circuit Board, Data Crystal, Quantum Bit
- Rare: Neural Processor, Holographic Emitter, Plasma Cell
- Epic: AI Core, Singularity Chip, Digital Rune
- Legendary: Tech Infinity Stone, Pure Information Essence

### Verdant Overgrowth (Nature World)

- Common: Leaf Fragment, Spore Dust, Bark Chip
- Uncommon: Bio-Catalyst, Ancient Seed, Vine Tendril
- Rare: Treant Heart, Bloom Crystal, Fungal Core
- Epic: World Tree Sapling, Gaia's Tear, Nature Rune
- Legendary: Verdant Infinity Stone, Primordial Life Essence

### Void Confluence (Void World)

- Common: Void Dust, Entropy Fragment, Shadow Wisp
- Uncommon: Paradox Shard, Reality Crack, Null Stone
- Rare: Chaos Orb, Dimensional Tear, Void Crystal
- Epic: Entropy Core, Reality Anchor, Void Rune
- Legendary: Void Infinity Stone, Pure Nothingness Essence

_Total: 75 unique items (15 per world × 5 worlds)_

---

## Appendix B: Specialization Effects

**Tank:**

- +15% success rate in Crimson Wastes (Fire World)
- +10% success rate in Verdant Overgrowth (Nature World)
- Reduces expedition failure chance (stretch feature)

**Scout:**

- +15% success rate in Frozen Expanse (Ice World)
- +10% success rate in Void Confluence (Void World)
- +5% base luck (stacks with upgrades)

**Scientist:**

- +15% success rate in Neon Sprawl (Tech World)
- +10% success rate in Void Confluence (Void World)
- +5% expedition speed (stacks with upgrades)

_Specializations cost 1 Epic item to unlock (encourages progression)_

---

**END OF GDD**

_Last Updated: October 19th, 2025_  
_Version: 1.0_  
_Status: Ready for Development_
