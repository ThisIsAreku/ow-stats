package ow

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"regexp"
	"strings"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

var Request *gorequest.SuperAgent = gorequest.New().Timeout(10 * time.Second).Set("User-Agent", "OW-STATS/1.0")

var prestigeRegex = regexp.MustCompile(`(?P<Prestige>0x025[0-9]{13})`)

type ProfileParser struct {
	doc       *goquery.Document
	platform  string
	region    string
	battleTag string
}

func (pp *ProfileParser) buildProfileUrl() string {
	return fmt.Sprintf("https://playoverwatch.com/en-us/career/%s/%s/%s",
		pp.platform,
		pp.region,
		pp.battleTag,
	)
}

func (pp *ProfileParser) fetchDocument() error {
	resp, _, err := Request.Get(pp.buildProfileUrl()).End()
	if len(err) != 0 {
		return err[0]
	}

	doc, parseErr := goquery.NewDocumentFromResponse(resp)
	if parseErr != nil {
		return parseErr
	}

	pp.doc = doc

	return nil
}

func NewProfileParser(region string, battleTag string) *ProfileParser {
	return &ProfileParser{
		platform:  "pc",
		region:    region,
		battleTag: battleTag,
	}
}

func (pp *ProfileParser) Parse() (*Profile, error) {
	if err := pp.fetchDocument(); err != nil {
		return nil, err
	}

	return &Profile{
		Stats: pp.parseProfileStats(),
	}, nil
}

func (pp *ProfileParser) parseProfileStats() *ProfileStats {
	return &ProfileStats{
		Competitive: pp.parseGamemodeStats(pp.doc.Find("#competitive")),
		Quickplay:   pp.parseGamemodeStats(pp.doc.Find("#quickplay")),
	}
}

func (pp *ProfileParser) parseGamemodeStats(selection *goquery.Selection) *GamemodeStats {
	gameStats, rollingAverageStats, averageStats := pp.parseGameStats(selection)
	return &GamemodeStats{
		Competitive:    selection.Is("#competitive"),
		Overall:        pp.parseOverallStats(selection),
		Game:           gameStats,
		RollingAverage: rollingAverageStats,
		Average:        averageStats,
	}
}

func (pp *ProfileParser) parseOverallStats(selection *goquery.Selection) *OverallStats {
	masthead := pp.doc.Find("div.masthead-player").First()
	playerLevel := masthead.Find("div.player-level").First()
	statsBoxRows := selection.Find(`div[data-group-id="stats"][data-category-id="0x02E00000FFFFFFFF"] table.data-table`).FilterFunction(func(i int, selection *goquery.Selection) bool {
		boxTitle := selection.Find("thead h5.stat-title").First().Text()
		return boxTitle == "Game" || boxTitle == "Miscellaneous"
	}).Find("tbody > tr")

	fnFindStat := func(statName string) int {
		r := 0
		statsBoxRows.EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if statName == selection.Children().Eq(0).Text() {
				if v, err := strconv.Atoi(selection.Children().Eq(1).Text()); err == nil {
					r = v

					return false
				}
			}

			return true
		})

		return r
	}

	overallStats := &OverallStats{
		Comprank: func() int {
			if comprank := masthead.Find("div.competitive-rank > div").First(); comprank.Length() == 1 {
				if v, err := strconv.Atoi(comprank.Text()); err == nil {
					return v
				}
			}

			return 0
		}(),

		Level: func() int {
			if level, err := strconv.Atoi(playerLevel.Find("div.u-vertical-center").First().Text()); err == nil {
				return level
			}

			return 0
		}(),

		Prestige: func() int {
			prestigeKey := prestigeRegex.FindString(playerLevel.AttrOr("style", ""))
			if prestige, ok := PRESTIGES[prestigeKey]; ok {
				return prestige
			}

			return 0
		}(),

		Games:  fnFindStat("Games Played"),
		Wins:   fnFindStat("Games Won"),
		Losses: fnFindStat("Games Lost"),
		Ties:   fnFindStat("Games Tied"),
	}

	if overallStats.Games != 0 {
		overallStats.WinRate = float32((float32(overallStats.Wins) / float32(overallStats.Games-overallStats.Ties)) * 100.0)
	}

	return overallStats
}

func (pp *ProfileParser) parseGameStats(selection *goquery.Selection) (*GameStats, *RollingAverageStats, *AverageStats) {
	gameStats := make(GameStats)
	rollingAverageStats := make(RollingAverageStats)
	averageStats := make(AverageStats)
	selection.Find(`div[data-group-id="stats"][data-category-id="0x02E00000FFFFFFFF"] table.data-table tbody > tr`).Each(func(i int, row *goquery.Selection) {
		key := row.Children().Eq(0).Text()
		value := row.Children().Eq(1).Text()

		sanitizedKey := SanitizeKey(key)
		sanitizedValue := SanitizeValue(value)
		if strings.HasSuffix(sanitizedKey, "_average") {
			averageStats[sanitizedKey[:len(sanitizedKey)-8]+"_avg"] = sanitizedValue
		} else if strings.HasSuffix(sanitizedKey, "_avg_per_10_min") {
			rollingAverageStats[sanitizedKey[:len(sanitizedKey)-15]] = sanitizedValue
		} else {
			gameStats[sanitizedKey] = sanitizedValue
		}
	})

	return &gameStats, &rollingAverageStats, &averageStats
}
