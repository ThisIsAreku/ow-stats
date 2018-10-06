package ow_stats

type Profile struct {
	Stats       *PlayerStats `json:"stats,omitempty"`
	Heroes      *HeroesStats `json:"heroes,omitempty"`
	UserProfile *ApiPlatform `json:"user_profile"`
}

// Player
type PlayerGameStats map[string]float32
type PlayerAverageStats map[string]float32
type PlayerRollingAverageStats map[string]float32

type PlayerStats struct {
	Competitive *PlayerGamemodeStats `json:"competitive,omitempty"`
	Quickplay   *PlayerGamemodeStats `json:"quickplay,omitempty"`
}

type PlayerGamemodeStats struct {
	Competitive    bool                       `json:"competitive,omitempty"`
	Average        *PlayerAverageStats        `json:"average_stats,omitempty"`
	RollingAverage *PlayerRollingAverageStats `json:"rolling_average_stats,omitempty"`
	Overall        *PlayerOverallStats        `json:"overall_stats,omitempty"`
	Game           *PlayerGameStats           `json:"game_stats,omitempty"`
}

type PlayerOverallStats struct {
	Level     int      `json:"level"`
	Comprank  int      `json:"comprank,omitempty"`
	Games     int      `json:"games"`
	WinRate   *float32 `json:"win_rate"`
	Losses    int      `json:"losses"`
	Wins      int      `json:"wins"`
	Ties      int      `json:"ties"`
	Prestige  int      `json:"prestige"`
	FullLevel int      `json:"full_level"`
}

// Heroes
type HeroesGamemodeStats map[string]*HeroGamemodeStats
type HeroPlaytimeStats map[string]float32
type HeroAverageStats map[string]float32
type HeroRollingAverageStats map[string]float32
type HeroSpecificStats map[string]float32
type HeroGeneralStats map[string]float32

type HeroesStats struct {
	Playtime *HeroesPlaytimeStats `json:"playtime,omitempty"`
	Stats    *HeroesStatsData     `json:"stats,omitempty"`
}

type HeroesStatsData struct {
	Competitive *HeroesGamemodeStats `json:"competitive,omitempty"`
	Quickplay   *HeroesGamemodeStats `json:"quickplay,omitempty"`
}

type HeroesPlaytimeStats struct {
	Competitive *HeroPlaytimeStats `json:"competitive,omitempty"`
	Quickplay   *HeroPlaytimeStats `json:"quickplay,omitempty"`
}

type HeroGamemodeStats struct {
	Average        *HeroAverageStats        `json:"average_stats,omitempty"`
	RollingAverage *HeroRollingAverageStats `json:"rolling_average_stats,omitempty"`
	Specific       *HeroSpecificStats       `json:"hero_stats,omitempty"`
	General        *HeroGeneralStats        `json:"general_stats,omitempty"`
}

type ApiPlatform struct {
	Platform    string `json:"platform"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	UrlName     string `json:"urlName"`
	PlayerLevel int    `json:"playerLevel"`
	Portrait    string `json:"portrait"`
	IsPublic    bool   `json:"isPublic"`
}
