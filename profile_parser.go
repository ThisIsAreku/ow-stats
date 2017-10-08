package ow_stats

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"regexp"
	"strings"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
	"log"
)

var Request *gorequest.SuperAgent = gorequest.New().Timeout(10 * time.Second).Set("User-Agent", "OW-STATS/1.0")

var prestigeRegex = regexp.MustCompile(`(?P<Prestige>0x025[0-9A-F]{13})`)

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
		Stats:  pp.parseProfileStats(),
		Heroes: pp.parseHeroesStats(),
	}, nil
}

func (pp *ProfileParser) parseProfileStats() *PlayerStats {
	return &PlayerStats{
		Competitive: pp.parseGamemodeStats(pp.doc.Find("#competitive")),
		Quickplay:   pp.parseGamemodeStats(pp.doc.Find("#quickplay")),
	}
}

func (pp *ProfileParser) parseGamemodeStats(selection *goquery.Selection) *PlayerGamemodeStats {
	if emptyGamemodeData(selection) {
		return nil
	}

	gameStats, rollingAverageStats, averageStats := pp.parseGameStats(selection)
	return &PlayerGamemodeStats{
		Competitive:    selection.Is("#competitive"),
		Overall:        pp.parseOverallStats(selection),
		Game:           gameStats,
		RollingAverage: rollingAverageStats,
		Average:        averageStats,
	}
}

func (pp *ProfileParser) parseOverallStats(selection *goquery.Selection) *PlayerOverallStats {
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

	overallStats := &PlayerOverallStats{
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
			styleAttr := playerLevel.AttrOr("style", "")
			prestigeKey := prestigeRegex.FindString(styleAttr)
			if prestige, ok := PRESTIGES[prestigeKey]; ok {
				return prestige
			}

			log.Printf("prestige '%s' was not found (tag is : '%s')\n", prestigeKey, styleAttr)

			return 0
		}(),

		Games:  fnFindStat("Games Played"),
		Wins:   fnFindStat("Games Won"),
		Losses: fnFindStat("Games Lost"),
		Ties:   fnFindStat("Games Tied"),
	}

	overallStats.WinRate = nil
	if overallStats.Games != 0 {
		wr := float32((float32(overallStats.Wins) / float32(overallStats.Games-overallStats.Ties)) * 100.0)
		overallStats.WinRate = &wr
	}

	overallStats.FullLevel = overallStats.Prestige*100 + overallStats.Level

	return overallStats
}

func (pp *ProfileParser) parseGameStats(selection *goquery.Selection) (*PlayerGameStats, *PlayerRollingAverageStats, *PlayerAverageStats) {
	gameStats := make(PlayerGameStats)
	rollingAverageStats := make(PlayerRollingAverageStats)
	averageStats := make(PlayerAverageStats)

	statsDiv := selection.Find(`div[data-group-id="stats"][data-category-id="0x02E00000FFFFFFFF"]`).First()
	if statsDiv.Length() == 0 {
		return nil, nil, nil
	}

	selection.Find(`table.data-table tbody > tr`).Each(func(i int, row *goquery.Selection) {
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

	gameStats["kpd"] = 0
	if deaths, ok := gameStats["deaths"]; ok && deaths > 0 {
		eliminations := gameStats["eliminations"]
		gameStats["kpd"] = eliminations / deaths
	}

	return &gameStats, &rollingAverageStats, &averageStats
}

func (pp *ProfileParser) parseHeroesStats() *HeroesStats {
	return &HeroesStats{
		Playtime: pp.parsePlaytimeStats(),
		Stats:    pp.parseHeroesStatsData(),
	}
}

func (pp *ProfileParser) parsePlaytimeStats() *HeroesPlaytimeStats {
	return &HeroesPlaytimeStats{
		Competitive: pp.parseHeroesPlaytimeStats(pp.doc.Find("#competitive")),
		Quickplay:   pp.parseHeroesPlaytimeStats(pp.doc.Find("#quickplay")),
	}
}

func (pp *ProfileParser) parseHeroesPlaytimeStats(selection *goquery.Selection) *HeroPlaytimeStats {
	if emptyGamemodeData(selection) {
		return nil
	}

	heroPlaytimeStats := make(HeroPlaytimeStats)

	statsDiv := selection.Find(`div[data-group-id="comparisons"][data-category-id="overwatch.guid.0x0860000000000021"]`).First()
	if statsDiv.Length() == 0 {
		return nil
	}

	statsDiv.Find(` div.bar-text`).Each(func(i int, row *goquery.Selection) {
		key := row.Find("div.title").Text()
		value := row.Find("div.description").Text()

		if sanitizedValue := SanitizeValue(value); sanitizedValue != 0 {
			heroPlaytimeStats[SanitizeKey(key)] = sanitizedValue
		}
	})

	return &heroPlaytimeStats
}

func (pp *ProfileParser) parseHeroesStatsData() *HeroesStatsData {
	return &HeroesStatsData{
		Competitive: pp.parseHeroesGamemodeStats(pp.doc.Find("#competitive")),
		Quickplay:   pp.parseHeroesGamemodeStats(pp.doc.Find("#quickplay")),
	}
}

func (pp *ProfileParser) parseHeroesGamemodeStats(selection *goquery.Selection) *HeroesGamemodeStats {
	if emptyGamemodeData(selection) {
		return nil
	}

	gamemodeStats := make(HeroesGamemodeStats)

	for heroName, heroId := range HEROES {
		heroPanel := selection.Find(fmt.Sprintf(`div[data-group-id="stats"][data-category-id="%s"]`, heroId)).First()
		if heroPanel.Length() == 1 {
			gamemodeStats[heroName] = pp.parseHeroGamemodeStats(heroPanel)
		}
	}

	return &gamemodeStats
}

func (pp *ProfileParser) parseHeroGamemodeStats(selection *goquery.Selection) (*HeroGamemodeStats) {
	generalStats := make(HeroGeneralStats)
	specificStats := make(HeroSpecificStats)
	rollingAverageStats := make(HeroRollingAverageStats)
	averageStats := make(HeroAverageStats)
	selection.Find(`table.data-table`).Each(func(i int, box *goquery.Selection) {
		boxTitle := box.Find("tr h5.stat-title").Text()
		rows := box.Find("tbody > tr")
		isSpecific := boxTitle == "Hero Specific"

		rows.Each(func(i int, row *goquery.Selection) {
			key := row.Children().Eq(0).Text()
			value := row.Children().Eq(1).Text()

			sanitizedKey := SanitizeKey(key)
			sanitizedValue := SanitizeValue(value)
			// Why is it different from Profile parsing ?
			if strings.HasSuffix(sanitizedKey, "_average") || strings.HasPrefix(sanitizedKey, "average_") {
				averageStats[sanitizedKey] = sanitizedValue
			} else if strings.HasSuffix(sanitizedKey, "_avg_per_10_min") {
				rollingAverageStats[sanitizedKey[:len(sanitizedKey)-15]] = sanitizedValue
			} else if isSpecific {
				specificStats[sanitizedKey] = sanitizedValue
			} else {
				generalStats[sanitizedKey] = sanitizedValue
			}
		})
	})

	return &HeroGamemodeStats{
		Average:        &averageStats,
		RollingAverage: &rollingAverageStats,
		General:        &generalStats,
		Specific:       &specificStats,
	}
}

func emptyGamemodeData(selection *goquery.Selection) bool {
	return strings.TrimSpace(selection.Find("h6.u-align-center").Text()) == "We don't have any data for this account in this mode yet."
}
