package data

import (
	"fmt"
	"strings"

	"ktn-x.com/tft-leaderboard/tft"
)

func PrintResult(pair *tft.TftPair) string {
	summoner := pair.Summoner
	details := "skipped"

	if pair.Rank != nil {
		index := make([]string, 4)
		index[0] = fmt.Sprintf("\t%s %s", pair.Rank.Tier, pair.Rank.Rank)
		index[1] = fmt.Sprintf("\t%d LP", pair.Rank.LeaguePoints)
		index[2] = fmt.Sprintf("\t%d Wins", pair.Rank.Wins)
		index[3] = fmt.Sprintf("\t%d Losses", pair.Rank.Losses)

		details = strings.Join(index, "\n")
	}

	return fmt.Sprintf("%s [%s]:\n%s", summoner.Name, summoner.Id, details)
}

func PrintUpdateMsg(pair *tft.TftPair) string {
	summoner := pair.Summoner

	if pair.Rank != nil {
		return fmt.Sprintf("summoner %s [%s] has no rank, not updating", summoner.Name, summoner.Id)
	}

	return fmt.Sprintf("updated summoner rank of %s [%s]", summoner.Name, summoner.Id)
}
