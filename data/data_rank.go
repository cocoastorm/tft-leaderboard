package data

import (
	"bytes"
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"ktn-x.com/tft-leaderboard/tft"
)

func (s *Store) ListContestantRanks(rankQueueType string) (tft.RankResults, error) {
	collection := make(tft.RankResults, 0, 7)
	err := s.storage.View(func(tx *bolt.Tx) error {
		contestants := tx.Bucket([]byte("contestants")).Cursor()
		leagueBucket := tx.Bucket([]byte("tft-leagues"))

		for k, encodedCt := contestants.First(); k != nil; k, encodedCt = contestants.Next() {
			contestant := &Contestant{}
			rank := &tft.TftLeague{}

			// decode contestant
			err := json.Unmarshal(encodedCt, contestant)
			if err != nil {
				return fmt.Errorf("failed to decode contestant %s: %w", k, err)
			}

			// decode rank
			pieces := [][]byte{[]byte(contestant.Summoner.Id), []byte(rankQueueType)}
			rankKey := bytes.Join(pieces, []byte("."))
			encodedRank := leagueBucket.Get(rankKey)

			if encodedRank != nil {
				err = json.Unmarshal(encodedRank, rank)
				if err != nil {
					return fmt.Errorf("failed to decode rank for contestant %s: %w", contestant.Summoner.Name, err)
				}
			}

			// pair result
			pair := &tft.TftPair{
				Summoner: contestant.Summoner,
				Rank:     rank,
			}
			collection = append(collection, pair.Transform())
		}

		return nil
	})

	return collection, err
}

func upsertRankLeague(tx *bolt.Tx, rank *tft.TftLeague) error {
	b := tx.Bucket([]byte("tft-leagues"))

	pieces := [][]byte{[]byte(rank.SummonerId), []byte(rank.QueueType)}
	key := bytes.Join(pieces, []byte("."))

	encoded, err := json.Marshal(rank)
	if err != nil {
		return err
	}

	return b.Put(key, encoded)
}

func (s *Store) UpdateContestantRanks(items []*tft.TftRanks) error {
	return s.storage.Batch(func(tx *bolt.Tx) error {
		for _, pair := range items {
			summoner := pair.Summoner

			// rank can be empty, if so skip it
			if len(pair.Ranks) == 0 {
				err := fmt.Errorf("rank list empty, skipping %s [%s]", summoner.Name, summoner.Id)
				fmt.Println(err)
				continue
			}

			for queueType, rank := range pair.Ranks {
				err := upsertRankLeague(tx, rank)
				if err != nil {
					return fmt.Errorf("failed to update %s rank for %s [%s]: %s", queueType, summoner.Name, summoner.Id, err)
				}
			}
		}

		// last-modified timestamp for "caching"
		upsertRankTimestamp(tx)

		return nil
	})
}
