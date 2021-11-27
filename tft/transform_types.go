package tft

import (
	"log"
	"strings"
)

const (
	Iron = iota + 1
	Bronze
	Silver
	Gold
	Platinum
	Diamond
	Masters
	Challengers
)

type RankOrder int

func lookupRank(name string) RankOrder {
	m := map[string]int{
		"":            0,
		"IRON":        Iron,
		"BRONZE":      Bronze,
		"SILVER":      Silver,
		"GOLD":        Gold,    
		"PLATINUM":    Platinum,
		"DIAMOND":     Diamond,
		"MASTERS":     Masters,
		"CHALLENGERS": Challengers,
	}

	if v, ok := m[name]; ok {
		return RankOrder(v)
	}

	return 0
}

func lookupTier(tier string) int {
	var (
		m = make(map[string]int, 4)
		tiers = []string{"I", "II", "III", "IV"}
	)

	m[""] = 0

	for i, v := range tiers {
		m[v] = i + 1
	}

	if v, ok := m[tier]; ok {
		return v
	}

	return 0
}

type WSummoner struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Level         int64  `json:"level"`
	ProfileIconId int64  `json:"profileIconId"`
}

type WTftLeague struct {
	QueueType    string `json:"queueType"`
	LeaguePoints int64  `json:"leaguePoints"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	Wins         int64  `json:"wins"`
	Losses       int64  `json:"losses"`
}

type WTftPair struct {
	Summoner *WSummoner  `json:"summoner"`
	Rank     *WTftLeague `json:"rank"`
}

type RankResults []*WTftPair

// sort.Interface

func (r RankResults) Len() int {
	return len(r)
}

func (r RankResults) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RankResults) Less(i int, j int) bool {
	pairA := r[i]
	pairB := r[j]

	// sort by rank
	var (
		rankA RankOrder = RankOrder(-1)
		rankB RankOrder = RankOrder(-1)
		tierA int
		tierB int
	)

	if (pairA.Rank != nil) {
		rankA = lookupRank(pairA.Rank.Rank)
		tierA = lookupTier(pairA.Rank.Tier)
	}

	if (pairB.Rank != nil) {
		rankB = lookupRank(pairB.Rank.Rank)
		tierB = lookupTier(pairB.Rank.Tier)
	}

	if rankA != rankB {
		return rankA > rankB
	}

	// if same rank, sort by tier
	if rankA > 0 && rankB > 0 {
		return tierA > tierB
	}

	// if no rank
	log.Printf("sort hit no rank item, using names")
	return strings.ToLower(pairA.Summoner.Name) < strings.ToLower(pairB.Summoner.Name)
}

