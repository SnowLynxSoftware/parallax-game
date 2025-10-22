# Next Steps - UX/UI Implementation Plan

## 🎯 Recommended Implementation Priority

Based on the GDD and current implementation state, here's the roadmap for building the game UI.

---

## **Phase 1: Core Game Loop (Week 1 - Priority 1) 🔥**

### **1. Teams Page** (`/teams`)

**Why first?** Foundation of the entire game - players need to see and understand their teams before anything else

**What to show:**

- List of all teams (3 total initially: Team 1 unlocked, Teams 2-3 locked with requirements)
- Team stats: Speed bonus, Luck bonus, Specialization
- Visual indicators: Available, On Expedition, Locked
- Basic team info cards with Font Awesome icons

**Backend Status:** ✅ Already implemented

- `TeamService.GetUserTeams()` - ready to use
- `TeamService.GetTeamById()` - ready to use

**Design Notes:**

- Use card-based layout
- Color-code team status (green=available, blue=on expedition, gray=locked)
- Show unlock requirements for locked teams
- Font Awesome icons for team types (user, user-shield, user-astronaut)

---

### **2. Expeditions Page** (`/expeditions`)

**Why second?** This is THE core gameplay loop - send teams, wait, get rewards

**What to show:**

- Left panel: Available teams (clickable to select)
- Center panel: Rift selection grid (cards with icons, duration, difficulty)
- Right panel: Selected expedition details + "Launch Expedition" button
- Locked rifts shown grayed out with unlock requirements

**Backend Status:** ✅ Already implemented

- `RiftService.GetAllRifts()` - returns all rifts with completion counts
- `RiftService.IsRiftUnlockedForUser()` - checks unlock status
- `ExpeditionService.StartExpedition()` - launches expeditions
- `TeamService.GetUserTeams()` - gets available teams

**Design Notes:**

- Rift cards with Font Awesome icons:
  - Tutorial: `fa-book-open`
  - Crimson Wastes: `fa-fire`
  - Frozen Expanse: `fa-snowflake`
  - Neon Sprawl: `fa-microchip`
  - Verdant Overgrowth: `fa-leaf`
  - Void Confluence: `fa-circle-notch`
- Show duration, difficulty, and reward preview on each card
- Disable "Launch" button if team is busy or rift is locked

---

### **3. Enhanced Dashboard** (update existing `/dashboard`)

**Why third?** Players need a home base to see active expeditions and quick stats

**What to add:**

- **Active Expeditions Panel** with countdown timers (JavaScript)
  - Show team name, rift name, progress bar, time remaining
  - "Claim Rewards" button when complete
- **Recent Loot Preview** (last 5-10 items)
  - Color-coded by rarity
  - Item icon + name + quantity
- **Quick Stats Card:**
  - Total expeditions completed
  - Dimensional Power Score
  - Legendary item count
- **Quick Actions:**
  - "Send Expedition" button → links to `/expeditions`
  - "View Teams" button → links to `/teams`
  - "View Inventory" button → links to `/inventory`

**Backend Status:** ✅ Already implemented

- `ExpeditionService.GetActiveExpeditions()` - gets in-progress expeditions
- `ExpeditionService.ClaimExpeditionRewards()` - claims completed expedition loot
- `InventoryService.GetUserInventory()` - gets recent items (can limit to last 10)

**JavaScript Needed:**

- Countdown timer that updates every second
- Progress bar animation
- Auto-refresh when expedition completes (poll every 30 seconds)

---

## **Phase 2: Progression Systems (Week 2 - Priority 2) 📈**

### **4. Inventory Page** (`/inventory`)

**What to show:**

- All collected loot organized by rarity or world type
- Filter options: All / Common / Uncommon / Rare / Epic / Legendary
- Sort options: By Rarity / By World Type / By Quantity / By Date
- Show which items are equipped on teams (if applicable)
- Tabs for Equipment vs Consumables

**Backend Status:** ✅ Already implemented

- `InventoryService.GetUserInventory()` - gets all items with equipped status
- `InventoryService.GetEquipment()` - filters equipment only
- `InventoryService.GetConsumables()` - filters consumables only

**Design Notes:**

- Card grid layout
- Color-coded borders by rarity
- Show quantity badge
- Equipped indicator (small team icon if equipped)

---

### **5. Leaderboard Page** (`/leaderboards`)

**What to show:**

- Three tabs:
  1. **Legendary Items** - Total legendary items collected
  2. **Dimensional Power Score** - Weighted value calculation (Common=1, Uncommon=5, Rare=25, Epic=125, Legendary=1000)
  3. **Total Expeditions** - Most expeditions completed
- Top 100 players per category
- User's own rank always visible (highlighted row, even if outside top 100)
- Last active timestamp

**Backend Status:** ❌ Not yet implemented

- Need new controller: `LeaderboardController`
- Need new queries for each leaderboard type
- Need caching strategy (update every 5 minutes)

**Design Notes:**

- Table layout with alternating row colors
- User's row highlighted in gold/yellow
- Rank badges for top 3 (🥇 🥈 🥉)
- Username truncation if too long

---

### **6. Team Upgrades** (Enhancement to Teams page)

**What to add:**

- Upgrade buttons on each team card
- Modal showing upgrade options:
  - Speed Upgrade (+5% speed) - costs 10 Uncommon + 5 Rare
  - Luck Upgrade (+3% luck) - costs 15 Uncommon + 3 Rare
  - Specialization (Tank/Scout/Scientist) - costs 1 Epic
- Show current vs upgraded stats
- Disable buttons if player lacks resources

**Backend Status:** ✅ Already implemented

- `TeamService.EquipItemToTeam()` - equips items
- `TeamService.UnequipItemFromTeam()` - unequips items
- `TeamService.ConsumeItemOnTeam()` - consumes items for upgrades
- `InventoryService.GetEquipment()` - gets available equipment

**Design Notes:**

- Use modal overlays
- Show resource costs with color-coding (red if insufficient, green if available)
- Confirmation step before consuming items
- Success animation when upgrade completes

---

## 🎨 **Visual Design Guidelines**

### **Color Scheme:**

- Background: Dark theme (`#1a1a2e`, `#16213e`)
- Primary accent: Cyan/Blue (`#0f4c75`, `#3282b8`)
- Rarity colors:
  - Common: `#9e9e9e` (gray)
  - Uncommon: `#4caf50` (green)
  - Rare: `#2196f3` (blue)
  - Epic: `#9c27b0` (purple)
  - Legendary: `#ff9800` (orange)

### **Typography:**

- Headers: System font stack or 'Orbitron' (futuristic feel)
- Body: 'Inter' or system font stack
- Font sizes: 14px base, 16px body, 20px+ headers

### **Layout:**

- Max width: 1200px container
- Card-based design with shadows/borders
- Responsive grid (3 columns desktop, 2 tablet, 1 mobile)
- Consistent spacing: 16px/24px/32px

### **Icons:**

- Font Awesome 6 (free version)
- Use solid style for most icons
- Color icons to match rarity/theme

---

## 🚀 **Quick Win Timeline (7 Days)**

**Day 1:** Build Teams page (read-only, just display)

- Create `teams.html` template
- Create `TeamsController` with one route
- Display team cards with stats

**Day 2:** Build Expeditions page (sending expeditions)

- Create `expeditions.html` template
- Create `ExpeditionsController`
- Implement rift selection + team selection
- Wire up StartExpedition API call

**Day 3:** Enhance Dashboard with active expeditions

- Add active expeditions panel to `dashboard.html`
- Add countdown timers (JavaScript)
- Add "Claim Rewards" functionality

**Day 4:** Add JavaScript polish

- Countdown timers that update every second
- Auto-refresh for completed expeditions
- Progress bar animations
- Modal for claiming rewards

**Day 5:** Build Inventory page

- Create `inventory.html` template
- Create `InventoryController`
- Display all loot with filters

**Day 6:** Build Leaderboard page

- Create `leaderboard.html` template
- Create `LeaderboardController`
- Implement three leaderboard types
- Add caching

**Day 7:** Polish, bug fixes, deploy

- CSS polish across all pages
- Mobile responsive testing
- Bug fixes
- Deploy to production

---

## 📋 **Technical Implementation Notes**

### **Controllers Needed:**

1. ✅ `UIController` - already exists (dashboard, welcome, etc.)
2. ❌ `TeamsController` - needs to be created
3. ❌ `ExpeditionsController` - needs to be created
4. ❌ `InventoryController` - needs to be created
5. ❌ `LeaderboardController` - needs to be created

### **Templates Needed:**

1. ✅ `dashboard.html` - exists, needs enhancement
2. ❌ `teams.html` - needs to be created
3. ❌ `expeditions.html` - needs to be created
4. ❌ `inventory.html` - needs to be created
5. ❌ `leaderboard.html` - needs to be created

### **JavaScript Components:**

1. Countdown timer utility
2. Progress bar animator
3. Modal system (or use Bootstrap modals)
4. Auto-refresh poller for expeditions
5. Filter/sort logic for inventory

### **API Endpoints Needed:**

All service methods already exist! Just need to wire up HTTP handlers:

- `GET /api/teams` - get user teams
- `GET /api/rifts` - get all rifts
- `POST /api/expeditions/start` - start expedition
- `GET /api/expeditions/active` - get active expeditions
- `POST /api/expeditions/:id/claim` - claim rewards
- `GET /api/inventory` - get inventory
- `GET /api/leaderboard/:type` - get leaderboard

---

## 🎮 **Ready to Start?**

The backend is fully tested and ready (185 passing tests!). All the business logic exists in the service layer. Now we just need to build the UI layer on top.

**Recommended starting point: Teams Page**

- Lowest risk
- Establishes visual design patterns
- Simple to implement (just display data)
- Gives players context before they can take actions
