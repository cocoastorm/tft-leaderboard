package data

import "ktn-x.com/tft-leaderboard/tft"

type Database interface {
	// Contestants
	UpdateContestant(c *Contestant) error
	UpdateContestants(contestants []*Contestant) error
	ListContestants() ([]*Contestant, error)

	// Ranks
	UpdateContestantRanks(items []*tft.TftPair) error
	ListContestantRanks() (tft.RankResults, error)

	// Timestamps
	UpdateRankTimestamp() error
	GetRankTimestamp() (uint64, error)

	// Utility
	Close() error
}
