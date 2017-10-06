package ow

type Profile struct {
	Stats *ProfileStats `json:"stats"`
}

type ProfileStats struct {
	Competitive *GamemodeStats `json:"competitive"`
	Quickplay   *GamemodeStats `json:"quickplay"`
}

type GamemodeStats struct {
	Competitive    bool                 `json:"competitive"`
	Average        *AverageStats        `json:"average_stats"`
	RollingAverage *RollingAverageStats `json:"rolling_average_stats"`
	Overall        *OverallStats        `json:"overall_stats"`
	Game           *GameStats           `json:"game_stats"`
}

type GameStats map[string]float32
type AverageStats map[string]float32
type RollingAverageStats map[string]float32

type OverallStats struct {
	Level    int     `json:"level"`
	Comprank int     `json:"comprank"`
	Games    int     `json:"games"`
	WinRate  float32 `json:"win_rate"`
	Losses   int     `json:"losses"`
	Wins     int     `json:"wins"`
	Ties     int     `json:"ties"`
	Prestige int     `json:"prestige"`
}
