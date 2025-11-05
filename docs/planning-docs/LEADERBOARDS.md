# Leaderboard Implementation Plan

**Status:** Ready for Implementation  
**Target:** Top 20 per category, 10-20 concurrent players  
**Caching Strategy:** 1-hour TTL with async refresh

---

## 🎯 Overview

### Three Leaderboard Categories:

1. **Legendary Items** - Total legendary items in inventory
2. **Dimensional Power Score** - Weighted inventory value (Common=1, Uncommon=5, Rare=25, Epic=125, Legendary=1000)
3. **Total Expeditions** - Total completed expeditions

### Technical Approach:

- Cache results in `leaderboard_cache` (metadata) and `leaderboard_cache_items` (rankings)
- Check cache age on every request (1-hour TTL)
- If stale, return current cache BUT spawn async goroutine to rebuild
- Show top 20 + current user's rank (if outside top 20)

---

## ⚠️ Known Trade-offs & Limitations

### 1. **Race Condition Risk**

Multiple concurrent requests could spawn multiple resync goroutines. **Mitigation:** Use sync.Mutex or atomic flag to ensure only one resync runs at a time.

### 2. **Data Staleness**

Players may see up to 1-hour-old rankings. Two players might both see themselves at the same rank. **Acceptable for game jam scope.**

### 3. **Random Tie-Breaking**

Tied players will appear in random order on each cache rebuild. This is unpredictable but meets requirements. **Consider:** Use `user_id` for deterministic ordering instead.

### 4. **No Pagination**

Only top 20 shown. With 10-20 players, everyone will likely be visible. **Future:** Add pagination if player base grows.

### 5. **Query Performance**

Power Score calculation requires joining inventory → loot_items and summing weighted values. With small player base, this is fine. **Watch:** Query time as inventory grows.

---

## 📊 Database Schema Design

### New Tables

```sql
-- Stores metadata about each leaderboard type
CREATE TABLE leaderboard_cache (
    id SERIAL PRIMARY KEY,
    leaderboard_type VARCHAR(50) NOT NULL UNIQUE, -- 'legendary', 'power', 'expeditions'
    last_synced TIMESTAMP NOT NULL DEFAULT NOW(),
    is_syncing BOOLEAN NOT NULL DEFAULT FALSE, -- Prevents concurrent rebuilds
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

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

-- Indices for fast lookups
CREATE INDEX idx_leaderboard_cache_type ON leaderboard_cache(leaderboard_type);
CREATE INDEX idx_leaderboard_cache_items_type_rank ON leaderboard_cache_items(leaderboard_type, rank);
CREATE INDEX idx_leaderboard_cache_items_user ON leaderboard_cache_items(user_id);
```

### Seeding Initial Data

```sql
-- Seed the three leaderboard types
INSERT INTO leaderboard_cache (leaderboard_type, last_synced) VALUES
    ('legendary', '2000-01-01 00:00:00'), -- Force initial sync
    ('power', '2000-01-01 00:00:00'),
    ('expeditions', '2000-01-01 00:00:00');
```

---

## 🏗️ Implementation Phases

---

## **Phase 1: Database Setup** ⏱️ 30 min

### Tasks:

- [ ] Create migration file: `migrations/YYYYMMDD_leaderboard_cache.sql`
- [ ] Add schema from above
- [ ] Run migration: `make migrate`
- [ ] Verify tables exist in Postgres

### Validation:

```sql
-- Check tables created
\dt leaderboard_*

-- Verify seed data
SELECT * FROM leaderboard_cache;
```

### Deliverable:

Migration file committed, tables exist in local DB.

---

## **Phase 2: Repository Layer** ⏱️ 2 hours

### File: `server/database/repositories/leaderboard.repository.go`

### Entities:

```go
type LeaderboardCacheEntity struct {
    ID              int       `db:"id"`
    LeaderboardType string    `db:"leaderboard_type"`
    LastSynced      time.Time `db:"last_synced"`
    IsSyncing       bool      `db:"is_syncing"`
    CreatedAt       time.Time `db:"created_at"`
    UpdatedAt       time.Time `db:"updated_at"`
}

type LeaderboardCacheItemEntity struct {
    ID              int       `db:"id"`
    LeaderboardType string    `db:"leaderboard_type"`
    UserID          int       `db:"user_id"`
    Username        string    `db:"username"`
    Score           int64     `db:"score"`
    Rank            int       `db:"rank"`
    CreatedAt       time.Time `db:"created_at"`
}
```

### Interface:

```go
type ILeaderboardRepository interface {
    // Cache metadata
    GetCacheMetadata(leaderboardType string) (*LeaderboardCacheEntity, error)
    SetSyncInProgress(leaderboardType string, inProgress bool) error
    UpdateLastSynced(leaderboardType string) error

    // Cache items
    GetTopRankings(leaderboardType string, limit int) ([]*LeaderboardCacheItemEntity, error)
    GetUserRank(leaderboardType string, userID int) (*LeaderboardCacheItemEntity, error)
    TruncateCache(leaderboardType string) error
    InsertCacheItems(items []*LeaderboardCacheItemEntity) error

    // Raw data queries (for building cache)
    GetLegendaryItemCounts() ([]*LeaderboardCacheItemEntity, error)
    GetPowerScores() ([]*LeaderboardCacheItemEntity, error)
    GetExpeditionCounts() ([]*LeaderboardCacheItemEntity, error)
}
```

### Implementation Notes:

**GetLegendaryItemCounts():**

```sql
SELECT
    u.id as user_id,
    u.username,
    COALESCE(SUM(ui.quantity), 0) as score
FROM users u
LEFT JOIN user_inventory ui ON ui.user_id = u.id
LEFT JOIN loot_items li ON li.id = ui.loot_item_id AND li.rarity = 'legendary'
GROUP BY u.id, u.username
HAVING COALESCE(SUM(ui.quantity), 0) > 0
ORDER BY score DESC
```

**GetPowerScores():**

```sql
SELECT
    u.id as user_id,
    u.username,
    COALESCE(SUM(
        ui.quantity *
        CASE li.rarity
            WHEN 'common' THEN 1
            WHEN 'uncommon' THEN 5
            WHEN 'rare' THEN 25
            WHEN 'epic' THEN 125
            WHEN 'legendary' THEN 1000
            ELSE 0
        END
    ), 0) as score
FROM users u
LEFT JOIN user_inventory ui ON ui.user_id = u.id
LEFT JOIN loot_items li ON li.id = ui.loot_item_id
GROUP BY u.id, u.username
HAVING COALESCE(SUM(ui.quantity), 0) > 0
ORDER BY score DESC
```

**GetExpeditionCounts():**

```sql
SELECT
    u.id as user_id,
    u.username,
    COUNT(e.id) as score
FROM users u
LEFT JOIN expeditions e ON e.user_id = u.id AND e.completed = true
GROUP BY u.id, u.username
HAVING COUNT(e.id) > 0
ORDER BY score DESC
```

### Tasks:

- [ ] Create repository file
- [ ] Implement all interface methods
- [ ] Add unit tests for repository methods
- [ ] Test raw data queries with sample data

### Validation:

- All tests pass: `go test ./server/database/repositories/leaderboard.repository_test.go`
- Manual query test returns expected results

### Deliverable:

`leaderboard.repository.go` with full implementation and passing tests.

---

## **Phase 3: Service Layer** ⏱️ 3 hours

### File: `server/services/leaderboard.service.go`

### DTOs (in `server/models/leaderboard.models.go`):

```go
type LeaderboardEntry struct {
    Rank     int    `json:"rank"`
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Score    int64  `json:"score"`
    IsCurrentUser bool `json:"is_current_user"`
}

type LeaderboardResponse struct {
    LeaderboardType string              `json:"leaderboard_type"`
    LastSynced      string              `json:"last_synced"` // RFC3339
    TopPlayers      []LeaderboardEntry  `json:"top_players"`
    CurrentUserRank *LeaderboardEntry   `json:"current_user_rank,omitempty"` // nil if in top 20
}
```

### Interface:

```go
type ILeaderboardService interface {
    GetLeaderboard(leaderboardType string, currentUserID int) (*LeaderboardResponse, error)
    // Internal method called by goroutine
    rebuildCache(leaderboardType string) error
}
```

### Core Logic:

**GetLeaderboard():**

1. Get cache metadata for `leaderboardType`
2. Check if `last_synced` is older than 1 hour
3. If stale AND not currently syncing:
   - Set `is_syncing = true`
   - Spawn goroutine: `go s.rebuildCache(leaderboardType)`
4. Fetch top 20 from `leaderboard_cache_items`
5. Fetch current user's rank
6. If user is in top 20, mark `is_current_user = true` on their entry
7. If user is outside top 20, populate `current_user_rank`
8. Return response

**rebuildCache():**

1. Acquire lock (use mutex to prevent concurrent rebuilds)
2. Call appropriate repo method (`GetLegendaryItemCounts()`, etc.)
3. Assign ranks (handle ties with random shuffle or deterministic sorting)
4. `TruncateCache(leaderboardType)`
5. `InsertCacheItems(items)`
6. `UpdateLastSynced(leaderboardType)`
7. Set `is_syncing = false`
8. Log any errors but don't crash (this runs in background)
9. Release lock

### Concurrency Safety:

```go
type LeaderboardService struct {
    repo  ILeaderboardRepository
    mutex sync.Mutex // Prevents concurrent cache rebuilds
}

func (s *LeaderboardService) rebuildCache(leaderboardType string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    // Set syncing flag
    if err := s.repo.SetSyncInProgress(leaderboardType, true); err != nil {
        return err
    }
    defer s.repo.SetSyncInProgress(leaderboardType, false)

    // ... rebuild logic ...
}
```

### Tie-Breaking Logic:

```go
// After getting raw scores, assign ranks
items := s.repo.GetLegendaryItemCounts()

// Shuffle items with same score for random tie-breaking
rand.Shuffle(len(items), func(i, j int) {
    items[i], items[j] = items[j], items[i]
})

// Assign ranks
currentRank := 1
for i, item := range items {
    if i > 0 && items[i].Score != items[i-1].Score {
        currentRank = i + 1
    }
    item.Rank = currentRank
}
```

**Note:** This shuffle happens on each rebuild, so tied players will see their relative positions change every hour. If this feels bad in testing, change to deterministic sorting (e.g., by `user_id`).

### Tasks:

- [ ] Create models file
- [ ] Create service file
- [ ] Implement `GetLeaderboard()`
- [ ] Implement `rebuildCache()`
- [ ] Add mutex for concurrency safety
- [ ] Write unit tests (mock repository)
- [ ] Test cache staleness detection
- [ ] Test goroutine spawning (may need integration test)

### Validation:

- Unit tests pass
- Manual test: Set `last_synced` to old date, call service, verify goroutine rebuilds cache
- Verify only one rebuild happens even with concurrent requests

### Deliverable:

`leaderboard.service.go` with full implementation and passing tests.

---

## **Phase 4: Controller & Routing** ⏱️ 1.5 hours

### File: `server/controllers/leaderboard.controller.go`

### Routes:

```
GET /leaderboards              -> Render HTML page
GET /api/leaderboards/:type    -> Return JSON for specific leaderboard
```

### Controller:

```go
type LeaderboardController struct {
    service     ILeaderboardService
    authMiddleware IAuthMiddleware
    templateService ITemplateService
}

func (c *LeaderboardController) MapController() chi.Router {
    r := chi.NewRouter()

    // UI
    r.Get("/", c.authMiddleware.RequireAuth(c.renderLeaderboardsPage))

    // API
    r.Get("/api/{type}", c.authMiddleware.RequireAuth(c.getLeaderboard))

    return r
}

func (c *LeaderboardController) getLeaderboard(w http.ResponseWriter, r *http.Request) {
    leaderboardType := chi.URLParam(r, "type")

    // Validate type
    validTypes := map[string]bool{
        "legendary": true,
        "power": true,
        "expeditions": true,
    }
    if !validTypes[leaderboardType] {
        http.Error(w, "Invalid leaderboard type", http.StatusBadRequest)
        return
    }

    user := c.authMiddleware.Authorize(r)

    result, err := c.service.GetLeaderboard(leaderboardType, user.ID)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get leaderboard")
        http.Error(w, "Leaderboard temporarily unavailable", http.StatusServiceUnavailable)
        return
    }

    json.NewEncoder(w).Encode(result)
}

func (c *LeaderboardController) renderLeaderboardsPage(w http.ResponseWriter, r *http.Request) {
    user := c.authMiddleware.Authorize(r)

    data := map[string]interface{}{
        "Title": "Leaderboards",
        "User":  user,
    }

    c.templateService.RenderTemplate(w, "leaderboards.html", data)
}
```

### Wiring (in `server/app.server.go`):

```go
// In NewAppServer.Start(), after other controllers:

leaderboardRepo := repositories.NewLeaderboardRepository(db)
leaderboardService := services.NewLeaderboardService(leaderboardRepo)
leaderboardController := controllers.NewLeaderboardController(
    leaderboardService,
    authMiddleware,
    templateService,
)

router.Mount("/leaderboards", leaderboardController.MapController())
```

### Tasks:

- [ ] Create controller file
- [ ] Implement routes
- [ ] Add validation for leaderboard types
- [ ] Wire up in `app.server.go`
- [ ] Test API endpoints with Postman/curl
- [ ] Verify auth middleware works

### Validation:

```bash
# Test legendary leaderboard
curl -H "Cookie: access_token=YOUR_TOKEN" \
  http://localhost:3000/leaderboards/api/legendary

# Should return JSON with top 20 + user rank
```

### Deliverable:

Working API endpoints that return cached leaderboard data.

---

## **Phase 5: Frontend Template** ⏱️ 3 hours

### File: `server/services/templates/pages/leaderboards.html`

### Layout:

```html
{{template "layouts/base.html" .}} {{define "content"}}
<div class="leaderboards-page">
  <div class="page-header">
    <h1>Dimensional Leaderboards</h1>
    <p class="subtitle">Compete against explorers across all parallel worlds</p>
  </div>

  <!-- Tab Navigation -->
  <div class="leaderboard-tabs">
    <button class="tab-btn active" data-type="legendary">
      <i class="fas fa-gem"></i> Legendary Items
    </button>
    <button class="tab-btn" data-type="power">
      <i class="fas fa-bolt"></i> Dimensional Power
    </button>
    <button class="tab-btn" data-type="expeditions">
      <i class="fas fa-rocket"></i> Total Expeditions
    </button>
  </div>

  <!-- Leaderboard Content (populated via JS) -->
  <div id="leaderboard-content">
    <div class="loading">
      <i class="fas fa-spinner fa-spin"></i> Loading leaderboard...
    </div>
  </div>
</div>

<style>
  .leaderboards-page {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  .page-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .page-header h1 {
    font-size: 2.5rem;
    margin-bottom: 0.5rem;
    color: #3282b8;
  }

  .subtitle {
    color: #9e9e9e;
    font-size: 1.1rem;
  }

  .leaderboard-tabs {
    display: flex;
    gap: 1rem;
    justify-content: center;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }

  .tab-btn {
    padding: 1rem 2rem;
    border: 2px solid #0f4c75;
    background: #1a1a2e;
    color: #fff;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    font-size: 1rem;
    font-weight: bold;
  }

  .tab-btn:hover {
    background: #0f4c75;
    transform: translateY(-2px);
  }

  .tab-btn.active {
    background: #3282b8;
    border-color: #3282b8;
  }

  .tab-btn i {
    margin-right: 0.5rem;
  }

  #leaderboard-content {
    background: #16213e;
    border-radius: 12px;
    padding: 2rem;
    min-height: 400px;
  }

  .loading {
    text-align: center;
    padding: 3rem;
    color: #9e9e9e;
    font-size: 1.2rem;
  }

  .leaderboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 2px solid #0f4c75;
  }

  .leaderboard-title {
    font-size: 1.5rem;
    color: #3282b8;
  }

  .last-updated {
    color: #9e9e9e;
    font-size: 0.9rem;
  }

  .leaderboard-table {
    width: 100%;
    border-collapse: collapse;
  }

  .leaderboard-table th {
    text-align: left;
    padding: 1rem;
    color: #3282b8;
    font-weight: bold;
    border-bottom: 2px solid #0f4c75;
  }

  .leaderboard-table td {
    padding: 1rem;
    border-bottom: 1px solid #0f4c75;
  }

  .leaderboard-table tr:hover {
    background: rgba(50, 130, 184, 0.1);
  }

  .leaderboard-table tr.current-user {
    background: rgba(255, 193, 7, 0.15);
    border-left: 4px solid #ffc107;
  }

  .rank-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    font-weight: bold;
    font-size: 1.1rem;
  }

  .rank-1 {
    background: linear-gradient(135deg, #ffd700, #ffed4e);
    color: #000;
  }
  .rank-2 {
    background: linear-gradient(135deg, #c0c0c0, #e8e8e8);
    color: #000;
  }
  .rank-3 {
    background: linear-gradient(135deg, #cd7f32, #daa520);
    color: #fff;
  }
  .rank-other {
    background: #0f4c75;
    color: #fff;
  }

  .score-value {
    font-size: 1.2rem;
    font-weight: bold;
    color: #3282b8;
  }

  .current-user-section {
    margin-top: 2rem;
    padding-top: 2rem;
    border-top: 2px solid #0f4c75;
  }

  .current-user-section h3 {
    color: #ffc107;
    margin-bottom: 1rem;
  }

  .error-message {
    text-align: center;
    padding: 2rem;
    color: #ff5252;
    font-size: 1.1rem;
  }

  /* Mobile responsive */
  @media (max-width: 768px) {
    .leaderboards-page {
      padding: 1rem;
    }

    .page-header h1 {
      font-size: 1.8rem;
    }

    .tab-btn {
      padding: 0.75rem 1.5rem;
      font-size: 0.9rem;
    }

    .leaderboard-table {
      font-size: 0.9rem;
    }

    .leaderboard-table th,
    .leaderboard-table td {
      padding: 0.75rem 0.5rem;
    }

    .rank-badge {
      width: 32px;
      height: 32px;
      font-size: 0.9rem;
    }
  }
</style>

<script>
  let currentType = "legendary";

  // Tab switching
  document.querySelectorAll(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      document
        .querySelectorAll(".tab-btn")
        .forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");
      currentType = btn.dataset.type;
      loadLeaderboard(currentType);
    });
  });

  // Load leaderboard data
  async function loadLeaderboard(type) {
    const content = document.getElementById("leaderboard-content");
    content.innerHTML =
      '<div class="loading"><i class="fas fa-spinner fa-spin"></i> Loading...</div>';

    try {
      const response = await fetch(`/leaderboards/api/${type}`);

      if (!response.ok) {
        throw new Error("Failed to load leaderboard");
      }

      const data = await response.json();
      renderLeaderboard(data);
    } catch (error) {
      console.error("Error loading leaderboard:", error);
      content.innerHTML =
        '<div class="error-message"><i class="fas fa-exclamation-triangle"></i> Leaderboard temporarily unavailable</div>';
    }
  }

  // Render leaderboard HTML
  function renderLeaderboard(data) {
    const typeNames = {
      legendary: "Legendary Items Collected",
      power: "Dimensional Power Score",
      expeditions: "Total Expeditions Completed",
    };

    let html = `
            <div class="leaderboard-header">
                <h2 class="leaderboard-title">${
                  typeNames[data.leaderboard_type]
                }</h2>
                <span class="last-updated">
                    <i class="fas fa-sync-alt"></i> Updated: ${formatDate(
                      data.last_synced
                    )}
                </span>
            </div>

            <table class="leaderboard-table">
                <thead>
                    <tr>
                        <th style="width: 80px;">Rank</th>
                        <th>Explorer</th>
                        <th style="text-align: right;">Score</th>
                    </tr>
                </thead>
                <tbody>
        `;

    data.top_players.forEach((player) => {
      const rankClass = player.rank <= 3 ? `rank-${player.rank}` : "rank-other";
      const rowClass = player.is_current_user ? "current-user" : "";

      html += `
                <tr class="${rowClass}">
                    <td>
                        <div class="rank-badge ${rankClass}">${
        player.rank
      }</div>
                    </td>
                    <td>
                        ${player.username}
                        ${
                          player.is_current_user
                            ? '<span style="color: #ffc107;"><i class="fas fa-user"></i> You</span>'
                            : ""
                        }
                    </td>
                    <td style="text-align: right;">
                        <span class="score-value">${player.score.toLocaleString()}</span>
                    </td>
                </tr>
            `;
    });

    html += "</tbody></table>";

    // Show current user rank if outside top 20
    if (data.current_user_rank && !data.current_user_rank.is_current_user) {
      html += `
                <div class="current-user-section">
                    <h3><i class="fas fa-user"></i> Your Ranking</h3>
                    <table class="leaderboard-table">
                        <tbody>
                            <tr class="current-user">
                                <td style="width: 80px;">
                                    <div class="rank-badge rank-other">${
                                      data.current_user_rank.rank
                                    }</div>
                                </td>
                                <td>${data.current_user_rank.username}</td>
                                <td style="text-align: right;">
                                    <span class="score-value">${data.current_user_rank.score.toLocaleString()}</span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            `;
    }

    document.getElementById("leaderboard-content").innerHTML = html;
  }

  // Format date
  function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMinutes = Math.floor((now - date) / 60000);

    if (diffMinutes < 1) return "Just now";
    if (diffMinutes < 60) return `${diffMinutes}m ago`;
    if (diffMinutes < 1440) return `${Math.floor(diffMinutes / 60)}h ago`;
    return date.toLocaleDateString();
  }

  // Load default leaderboard on page load
  loadLeaderboard(currentType);

  // Auto-refresh every 5 minutes
  setInterval(() => {
    loadLeaderboard(currentType);
  }, 300000);
</script>
{{end}}
```

### Template Service Registration:

Add to `server/services/template.service.go`:

```go
templates = append(templates,
    "server/services/templates/pages/leaderboards.html",
)
```

### Add to Navbar:

In `server/services/templates/layouts/components/navbar.html`, add:

```html
<a href="/leaderboards" class="nav-link">
  <i class="fas fa-trophy"></i> Leaderboards
</a>
```

### Tasks:

- [ ] Create template file
- [ ] Add to template service
- [ ] Add navbar link
- [ ] Test page renders
- [ ] Test tab switching
- [ ] Test mobile responsive layout
- [ ] Verify data loads from API

### Validation:

- Visit `/leaderboards` while logged in
- Click between tabs, verify data loads
- Check that current user is highlighted
- Check that user outside top 20 shows in separate section
- Test on mobile viewport

### Deliverable:

Fully functional leaderboard page with live data and tab switching.

---

## **Phase 6: Testing & Polish** ⏱️ 2 hours

### Unit Tests:

- [ ] Repository tests with mock database
- [ ] Service tests with mock repository
- [ ] Test cache staleness logic
- [ ] Test tie-breaking logic
- [ ] Test concurrent resync prevention

### Integration Tests:

- [ ] End-to-end: Create users, expeditions, inventory → verify leaderboard ranks
- [ ] Test cache rebuild with real data
- [ ] Test API returns expected JSON structure
- [ ] Test error handling (DB connection fails, etc.)

### Manual Testing Checklist:

- [ ] Create 5+ test users with different inventory/expedition counts
- [ ] Verify Legendary leaderboard shows correct counts
- [ ] Verify Power Score calculation is accurate (manually calculate expected scores)
- [ ] Verify Expeditions leaderboard matches expedition table
- [ ] Force cache stale (set `last_synced` to old date), verify rebuild happens
- [ ] Make concurrent requests, verify only one rebuild runs
- [ ] Test with user in top 20 vs outside top 20
- [ ] Test with ties (multiple users same score)
- [ ] Check mobile layout on phone/tablet
- [ ] Verify "Leaderboard temporarily unavailable" shows on DB error

### Performance Testing:

- [ ] Time the raw data queries with 10-20 users
- [ ] Verify queries complete under 500ms
- [ ] Check cache table size after multiple rebuilds

### Polish:

- [ ] Add loading spinners while fetching data
- [ ] Add smooth transitions when switching tabs
- [ ] Ensure color scheme matches rest of game
- [ ] Verify all Font Awesome icons load
- [ ] Check for console errors in browser
- [ ] Test keyboard navigation (tab key)

### Deliverable:

Production-ready leaderboard feature with comprehensive test coverage.

---

## 📋 Implementation Checklist (Copy to GitHub Issues)

### Phase 1: Database ✅

- [ ] Create migration file
- [ ] Run migration
- [ ] Verify tables exist

### Phase 2: Repository ✅

- [ ] Create repository file
- [ ] Implement interface
- [ ] Write unit tests
- [ ] All tests pass

### Phase 3: Service ✅

- [ ] Create models file
- [ ] Create service file
- [ ] Implement GetLeaderboard
- [ ] Implement rebuildCache
- [ ] Add concurrency safety
- [ ] Write unit tests
- [ ] All tests pass

### Phase 4: Controller ✅

- [ ] Create controller file
- [ ] Implement routes
- [ ] Wire in app.server.go
- [ ] Test API endpoints

### Phase 5: Frontend ✅

- [ ] Create template file
- [ ] Register in template service
- [ ] Add navbar link
- [ ] Test page rendering
- [ ] Test mobile responsive

### Phase 6: Testing ✅

- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing
- [ ] Performance check
- [ ] Polish

---

## 🚀 Estimated Timeline

| Phase               | Time      | Cumulative   |
| ------------------- | --------- | ------------ |
| Phase 1: Database   | 30 min    | 30 min       |
| Phase 2: Repository | 2 hours   | 2.5 hours    |
| Phase 3: Service    | 3 hours   | 5.5 hours    |
| Phase 4: Controller | 1.5 hours | 7 hours      |
| Phase 5: Frontend   | 3 hours   | 10 hours     |
| Phase 6: Testing    | 2 hours   | **12 hours** |

**Total: ~12 hours** (1.5 work days)

---

## 🔧 Troubleshooting Guide

### Problem: Leaderboard shows stale data even after 1 hour

**Check:**

- Is `last_synced` updating? Query `SELECT * FROM leaderboard_cache;`
- Is goroutine running? Add logging to `rebuildCache()`
- Is mutex causing deadlock? Check logs for goroutine completion

### Problem: Multiple cache rebuilds happening simultaneously

**Check:**

- Is mutex properly acquired/released?
- Is `is_syncing` flag being set correctly?
- Add logging at start/end of `rebuildCache()`

### Problem: Leaderboard always shows "temporarily unavailable"

**Check:**

- Are queries returning errors? Check Postgres logs
- Are indices created? Run `\d leaderboard_cache_items` in psql
- Is network connection stable?

### Problem: Ranks are incorrect or duplicate

**Check:**

- Is tie-breaking logic working? Verify tied scores get sequential ranks
- Are ranks being assigned correctly in loop?
- Manually query raw data and compare to cache

### Problem: Current user not highlighted in top 20

**Check:**

- Is `is_current_user` flag being set in service?
- Is `user_id` matching correctly?
- Check frontend JS `player.is_current_user` condition

---

## 🎯 Success Criteria

✅ **Functional Requirements:**

- [ ] Three leaderboards work (legendary, power, expeditions)
- [ ] Top 20 players shown
- [ ] Current user rank always visible
- [ ] Cache updates every 1 hour
- [ ] Page loads in under 2 seconds
- [ ] Mobile responsive

✅ **Technical Requirements:**

- [ ] No N+1 queries
- [ ] Concurrent requests don't cause race conditions
- [ ] Error handling returns graceful messages
- [ ] Code is tested (>80% coverage)

✅ **UX Requirements:**

- [ ] Tabs switch smoothly
- [ ] Loading states are clear
- [ ] Current user is visually highlighted
- [ ] Top 3 ranks have special badges
- [ ] Last updated time is shown

---

## 📝 Final Notes

### Potential Future Enhancements (Post-Jam):

- Add pagination for players beyond top 20
- Add filters (filter by rift type, time range)
- Add personal best tracking (peak rank achieved)
- Add historical graphs (rank over time)
- Add achievement badges next to names
- Add guild/team leaderboards

### Known Limitations:

- Random tie-breaking may feel inconsistent
- 1-hour cache TTL means slight data staleness
- No real-time updates (requires WebSocket)
- No anti-cheat mechanisms

### If You Fall Behind Schedule:

**Cut these in order:**

1. Auto-refresh every 5 minutes (just manual refresh)
2. "Your Ranking" section for users outside top 20 (just show message)
3. Fancy rank badges for top 3 (just show numbers)
4. Mobile responsive polish (functional but less pretty)

**Do NOT cut:**

- Cache system (direct queries will be too slow)
- Three leaderboard types (core feature)
- Current user highlighting (critical UX)

---

**Ready to start Phase 1?** Let me know when you want to begin implementation, or if you have any questions about the plan.
