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

func lookupTier(tier string) int {
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

	if v, ok := m[tier]; ok {
		return v
	}

	return 0
}

func lookupRank(rank string) int {
	var (
		m     = make(map[string]int, 5)
		tiers = []string{"I", "II", "III", "IV"}
	)

	m[""] = 0

	for i, v := range tiers {
		m[v] = i + 1
	}

	if v, ok := m[rank]; ok {
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
		tierA int = -1
		tierB int = -1
		rankA int
		rankB int
	)

	if pairA.Rank != nil {
		tierA = lookupTier(pairA.Rank.Tier)
		rankA = lookupRank(pairA.Rank.Rank)
	}

	if pairB.Rank != nil {
		tierB = lookupTier(pairB.Rank.Tier)
		rankB = lookupRank(pairB.Rank.Rank)
	}

	// sort by names if unranked
	if rankA == -1 && rankB == -1 {
		return strings.ToLower(pairA.Summoner.Name) < strings.ToLower(pairB.Summoner.Name)
	}

	// if A has tier, but B doesn't
	if tierA > -1 && tierB == -1 {
		return false
	}

	// if B has tier, but A doesn't
	if tierB > -1 && tierA == -1 {
		return true
	}

	// sort by tier
	if tierA != tierB {
		return tierA < tierB
	}

	// if not, sort by rank
	if tierA == tierB {
		return rankA < rankB
	}

	// if no rank
	log.Printf("sort hit no rank item, using names")
	return strings.ToLower(pairA.Summoner.Name) < strings.ToLower(pairB.Summoner.Name)
}
