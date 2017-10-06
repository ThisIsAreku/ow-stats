package ow

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"regexp"
)

var prestigeRegex = regexp.MustCompile(`(?P<Prestige>0x025[0-9]{13})`)

type ProfileParser struct {
	doc      *goquery.Document
	prestige int
}

func NewProfileParser(doc *goquery.Document) *ProfileParser {
	return &ProfileParser{
		doc: doc,
	}
}

func (pp *ProfileParser) Parse() (*Profile, error) {
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
	return &GamemodeStats{
		Competitive: selection.Is("#competitive"),
		Overall:     pp.parseOverallStats(selection),
	}
}

func (pp *ProfileParser) parseOverallStats(selection *goquery.Selection) *OverallStats {
	masthead := pp.doc.Find("div.masthead-player").First()
	playerLevel := masthead.Find("div.player-level").First()
	statsBoxRows := selection.Find(`div[data-group-id="stats"][data-category-id="0x02E00000FFFFFFFF"] table.data-table`).FilterFunction(func(i int, selection *goquery.Selection) bool {
		boxTitle := selection.Find("thead h5.stat-title").First().Text()
		return boxTitle == "Game" || boxTitle == "Miscellaneous"
	}).Find("tbody > tr")

	//statsBoxRows.Each(func(i int, selection *goquery.Selection) {
	//	fmt.Printf("%s => %s\n",
	//		selection.Children().Eq(0).Text(),
	//		selection.Children().Eq(1).Text(),
	//	)
	//})
	//fmt.Println()

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
