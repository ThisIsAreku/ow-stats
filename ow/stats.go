package ow

type Profile struct {
	Stats  *PlayerStats `json:"stats"`
	Heroes *HeroesStats `json:"heroes"`
}

// Player
type PlayerGameStats map[string]float32
type PlayerAverageStats map[string]float32
type PlayerRollingAverageStats map[string]float32

type PlayerStats struct {
	Competitive *PlayerGamemodeStats `json:"competitive"`
	Quickplay   *PlayerGamemodeStats `json:"quickplay"`
}

type PlayerGamemodeStats struct {
	Competitive    bool                       `json:"competitive"`
	Average        *PlayerAverageStats        `json:"average_stats"`
	RollingAverage *PlayerRollingAverageStats `json:"rolling_average_stats"`
	Overall        *PlayerOverallStats        `json:"overall_stats"`
	Game           *PlayerGameStats           `json:"game_stats"`
}

type PlayerOverallStats struct {
	Level    int     `json:"level"`
	Comprank int     `json:"comprank"`
	Games    int     `json:"games"`
	WinRate  float32 `json:"win_rate"`
	Losses   int     `json:"losses"`
	Wins     int     `json:"wins"`
	Ties     int     `json:"ties"`
	Prestige int     `json:"prestige"`
}

// Heroes
type HeroesGamemodeStats map[string]*HeroGamemodeStats
type HeroPlaytimeStats map[string]float32
type HeroAverageStats map[string]float32
type HeroRollingAverageStats map[string]float32
type HeroSpecificStats map[string]float32
type HeroGeneralStats map[string]float32

type HeroesStats struct {
	Playtime *HeroesPlaytimeStats `json:"playtime"`
	Stats    *HeroesStatsData     `json:"stats"`
}

type HeroesStatsData struct {
	Competitive *HeroesGamemodeStats `json:"competitive"`
	Quickplay   *HeroesGamemodeStats `json:"quickplay"`
}

type HeroesPlaytimeStats struct {
	Competitive *HeroPlaytimeStats `json:"competitive"`
	Quickplay   *HeroPlaytimeStats `json:"quickplay"`
}

type HeroGamemodeStats struct {
	Average        *HeroAverageStats        `json:"average_stats"`
	RollingAverage *HeroRollingAverageStats `json:"rolling_average_stats"`
	Specific       *HeroSpecificStats       `json:"hero_stats"`
	General        *HeroGeneralStats        `json:"general_stats"`
}
