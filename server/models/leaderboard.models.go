package models

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Rank          int    `json:"rank"`
	UserID        int    `json:"user_id"`
	Username      string `json:"username"`
	Score         int64  `json:"score"`
	IsCurrentUser bool   `json:"is_current_user"`
}

// LeaderboardResponse represents the full leaderboard response
type LeaderboardResponse struct {
	LeaderboardType string             `json:"leaderboard_type"`
	LastSynced      string             `json:"last_synced"` // RFC3339
	TopPlayers      []LeaderboardEntry `json:"top_players"`
	CurrentUserRank *LeaderboardEntry  `json:"current_user_rank,omitempty"` // nil if in top 20
}
