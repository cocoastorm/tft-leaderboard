package data

import (
	"fmt"
	"strings"

	"ktn-x.com/tft-leaderboard/tft"
)

func PrintResult(pair *tft.TftPair) string {
	summoner := pair.Summoner
	title := fmt.Sprintf("%s [%s]", summoner.Name, summoner.Id)

	// nothing to see here, move along
	if len(pair.Rank) <= 0 {
		return fmt.Sprintf("%s:\n%s", title, "skipped")
	}

	results := make([]string, len(pair.Rank))

	for _, rank := range pair.Rank {
		index := []string{
			fmt.Sprintf("\t\t%s %s", rank.Tier, rank.Rank),
			fmt.Sprintf("\t\t%d LP", rank.LeaguePoints),
			fmt.Sprintf("\t\t%d Wins", rank.Wins),
			fmt.Sprintf("\t\t%d Losses", rank.Losses),
		}

		line := strings.Join(index, "\n")
		results = append(results, fmt.Sprintf("\t[%s]:\n%s", rank.QueueType, line))
	}

	return fmt.Sprintf("%s:\n%s", title, strings.Join(results, "\n"))
}

func PrintUpdateMsg(pair *tft.TftPair) string {
	summoner := pair.Summoner

	if pair.Rank != nil {
		return fmt.Sprintf("summoner %s [%s] has no rank, not updating", summoner.Name, summoner.Id)
	}

	return fmt.Sprintf("updated summoner rank of %s [%s]", summoner.Name, summoner.Id)
}
